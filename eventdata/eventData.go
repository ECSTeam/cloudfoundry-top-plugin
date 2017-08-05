// Copyright (c) 2016 ECS Team, Inc. - All Rights Reserved
// https://github.com/ECSTeam/cloudfoundry-top-plugin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eventdata

import (
	"encoding/binary"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/config"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventApp"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventCell"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventEventType"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventRoute"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/space"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/mohae/deepcopy"
)

type EventData struct {

	// This time it set at clone time
	StatsTime time.Time
	// Key: appId
	AppMap  map[string]*eventApp.AppStats
	CellMap map[string]*eventCell.CellStats
	// Domain name: Both shared + private
	DomainMap    map[string]*eventRoute.DomainStats
	EventTypeMap map[events.Envelope_EventType]*eventEventType.EventTypeStats

	EnableRouteTracking bool
	TotalEvents         int64
	mu                  *sync.Mutex
	logHttpAccess       *EventLogHttpAccess
	eventProcessor      *EventProcessor
	apiUrl              string
	apiUrlRegexp        *regexp.Regexp
	urlPCF17Regexp      *regexp.Regexp
	routeUrlRegexp      *regexp.Regexp
}

func NewEventData(mu *sync.Mutex, eventProcessor *EventProcessor) *EventData {

	logHttpAccess := NewEventLogHttpAccess()
	apiUrl := util.GetApiEndpointNoProtocol(eventProcessor.cliConnection)

	apiUrlRegexpStr := `^\/v[0-9]\/([^\/]*)\/([^\/]*)`
	apiUrlRegexp := regexp.MustCompile(apiUrlRegexpStr)

	urlPCF17RegexpStr := `^http[s]?:\/\/[^\/]*(\/v[0-9]\/.*)$`
	urlPCF17Regexp := regexp.MustCompile(urlPCF17RegexpStr)

	//routeUrlRegexpStr := `^(?:http[s]?:\/\/)?([^\.]+)\.([^\/^:]*)(?::[0-9]+)?(.*)$`

	routeUrlRegexpStr := `^(?:http[s]?:\/\/)?(?:((?:[0-9]{1,3}\.){3}[0-9]{1,3})|([^\.\:\/]+)(?:\.([^\/^:]*))?)(?::([0-9]+))?(.*)$`
	routeUrlRegexp := regexp.MustCompile(routeUrlRegexpStr)

	return &EventData{

		AppMap:         make(map[string]*eventApp.AppStats),
		CellMap:        make(map[string]*eventCell.CellStats),
		DomainMap:      make(map[string]*eventRoute.DomainStats),
		EventTypeMap:   make(map[events.Envelope_EventType]*eventEventType.EventTypeStats),
		TotalEvents:    0,
		mu:             mu,
		logHttpAccess:  logHttpAccess,
		eventProcessor: eventProcessor,
		apiUrl:         apiUrl,
		apiUrlRegexp:   apiUrlRegexp,
		urlPCF17Regexp: urlPCF17Regexp,
		routeUrlRegexp: routeUrlRegexp,
	}

}

func (ed *EventData) Process(instanceId int, msg *events.Envelope) {

	ed.mu.Lock()
	defer ed.mu.Unlock()

	ed.UpdateEventStats(msg)

	eventType := msg.GetEventType()
	switch eventType {
	case events.Envelope_HttpStartStop:
		ed.httpStartStopEvent(msg)
	case events.Envelope_ContainerMetric:
		ed.containerMetricEvent(msg)
	case events.Envelope_LogMessage:
		ed.logMessageEvent(msg)
	case events.Envelope_ValueMetric:
		ed.valueMetricEvent(msg)
	case events.Envelope_CounterEvent:
		// Message that is sent on nozzle when its not keeping up.
		// https://docs.cloudfoundry.org/loggregator/log-ops-guide.html#slow-noz
		if msg.CounterEvent.GetName() == "TruncatingBuffer.DroppedMessages" &&
			(msg.GetOrigin() == "DopplerServer" || msg.GetOrigin() == "doppler") {
			ed.droppedMessages(instanceId, msg)
		}
	case events.Envelope_Error:
	default:
	}

}

func (ed *EventData) UpdateEventStats(msg *events.Envelope) {
	eventTypeStats := ed.FindEventTypeStats(msg.GetEventType())
	originStats := eventTypeStats.FindEventOriginStats(msg.GetOrigin())
	eventDetailStats := originStats.FindEventDetailStats(msg)
	eventDetailStats.EventCount = eventDetailStats.EventCount + 1
	eventDetailStats.LastEventTime = time.Now()
	// TODO: Replace above time.Now() with time from msg
	// eventDetailStats.LastEventTime = msg.GetTimestamp
}

func (ed *EventData) FindEventTypeStats(eventType events.Envelope_EventType) *eventEventType.EventTypeStats {
	eventTypeStats := ed.EventTypeMap[eventType]
	if eventTypeStats == nil {
		eventTypeStats = eventEventType.NewEventTypeStats(eventType)
		ed.EventTypeMap[eventType] = eventTypeStats
	}
	return eventTypeStats
}

// Yuck -- need to figure out a getter way to clone the current data object
// Specifically around RateCounter and AvgTracker object
func (ed *EventData) Clone() *EventData {

	ed.mu.Lock()
	defer ed.mu.Unlock()

	clone := deepcopy.Copy(ed).(*EventData)
	now := time.Now()
	clone.StatsTime = now
	clone.eventProcessor = ed.eventProcessor

	for _, appStat := range ed.AppMap {

		// Check if this app has been deleted
		if ed.eventProcessor.GetMetadataManager().IsAppDeleted(appStat.AppId) {
			ed.eventProcessor.GetMetadataManager().RemoveAppFromDeletedQueue(appStat.AppId)
			delete(ed.AppMap, appStat.AppId)
			delete(clone.AppMap, appStat.AppId)
			continue
		}

		clonedAppStat := clone.AppMap[appStat.AppId]
		/*
			httpAllCount := int64(0)
			http2xxCount := int64(0)
			http3xxCount := int64(0)
			http4xxCount := int64(0)
			http5xxCount := int64(0)
		*/

		responseL60TimeArray := make([]*util.AvgTracker, 0)
		responseL10TimeArray := make([]*util.AvgTracker, 0)
		responseL1TimeArray := make([]*util.AvgTracker, 0)
		totalTraffic := eventApp.NewTrafficStats()

		for instanceId, containerTraffic := range appStat.ContainerTrafficMap {

			rate60 := containerTraffic.ResponseL60Time.Rate()
			clone.AppMap[appStat.AppId].ContainerTrafficMap[instanceId].EventL60Rate = rate60
			totalTraffic.EventL60Rate = totalTraffic.EventL60Rate + rate60

			clone.AppMap[appStat.AppId].ContainerTrafficMap[instanceId].AvgResponseL60Time = containerTraffic.ResponseL60Time.Avg()

			rate10 := containerTraffic.ResponseL10Time.Rate()
			clone.AppMap[appStat.AppId].ContainerTrafficMap[instanceId].EventL10Rate = rate10
			totalTraffic.EventL10Rate = totalTraffic.EventL10Rate + rate10

			clone.AppMap[appStat.AppId].ContainerTrafficMap[instanceId].AvgResponseL10Time = containerTraffic.ResponseL10Time.Avg()

			rate1 := containerTraffic.ResponseL1Time.Rate()
			clone.AppMap[appStat.AppId].ContainerTrafficMap[instanceId].EventL1Rate = rate1
			totalTraffic.EventL1Rate = totalTraffic.EventL1Rate + rate1

			clone.AppMap[appStat.AppId].ContainerTrafficMap[instanceId].AvgResponseL1Time = containerTraffic.ResponseL1Time.Avg()

			/*
				httpAllCount = httpAllCount + containerTraffic.HttpAllCount
				http2xxCount = http2xxCount + containerTraffic.Http2xxCount
				http3xxCount = http3xxCount + containerTraffic.Http3xxCount
				http4xxCount = http4xxCount + containerTraffic.Http4xxCount
				http5xxCount = http5xxCount + containerTraffic.Http5xxCount
			*/

			responseL60TimeArray = append(responseL60TimeArray, containerTraffic.ResponseL60Time)
			responseL10TimeArray = append(responseL10TimeArray, containerTraffic.ResponseL10Time)
			responseL1TimeArray = append(responseL1TimeArray, containerTraffic.ResponseL1Time)

		}

		totalTraffic.AvgResponseL60Time = util.AvgMultipleTrackers(responseL60TimeArray)
		totalTraffic.AvgResponseL10Time = util.AvgMultipleTrackers(responseL10TimeArray)
		totalTraffic.AvgResponseL1Time = util.AvgMultipleTrackers(responseL1TimeArray)

		/*
			totalTraffic.HttpAllCount = httpAllCount
			totalTraffic.Http2xxCount = http2xxCount
			totalTraffic.Http3xxCount = http3xxCount
			totalTraffic.Http4xxCount = http4xxCount
			totalTraffic.Http5xxCount = http5xxCount
		*/
		clonedAppStat.TotalTraffic = totalTraffic

		// Check if we need to remove data for any old containers -- containers that haven't reported in for awhile
		for containerIndex, cs := range appStat.ContainerArray {
			if cs != nil {
				// If we haven't gotten a container update in DeadContainerSeconds then remove the entire container data
				if cs.LastUpdateTime == nil || now.Sub(*cs.LastUpdateTime) > time.Second*config.DeadContainerSeconds {
					// Remove container stats for realtime data and the clone
					appStat.ContainerArray[containerIndex] = nil
					clonedAppStat.ContainerArray[containerIndex] = nil
				} else {
					if cs.ContainerMetric != nil {
						// If we haven't gotten a container update in StaleContainerSeconds then remove just the
						// container metrics not the entire container data
						if cs.LastUpdateTime == nil || now.Sub(*cs.LastUpdateTime) > time.Second*config.StaleContainerSeconds {
							// Remove container stats for realtime data and the clone
							appStat.ContainerArray[containerIndex].ContainerMetric = nil
							clonedAppStat.ContainerArray[containerIndex].ContainerMetric = nil
						}
					}

				}

			}
		}

	}

	return clone
}

func (ed *EventData) GetTotalEvents() int64 {
	return ed.TotalEvents
}

func (ed *EventData) Clear() {

	ed.mu.Lock()
	defer ed.mu.Unlock()

	ed.AppMap = make(map[string]*eventApp.AppStats)
	ed.CellMap = make(map[string]*eventCell.CellStats)
	ed.DomainMap = make(map[string]*eventRoute.DomainStats)
	ed.TotalEvents = 0
}

func (ed *EventData) droppedMessages(instanceId int, msg *events.Envelope) {
	delta := msg.GetCounterEvent().GetDelta()
	total := msg.GetCounterEvent().GetTotal()
	text := fmt.Sprintf("Nozzle #%v - Upstream message indicates the nozzle or the TrafficController is not keeping up. Dropped delta: %v, total: %v",
		instanceId, delta, total)
	toplog.Error(text)
}

func (ed *EventData) handleHttpAccessLogLine(logLine string) {
	ed.logHttpAccess.parseHttpAccessLogLine(logLine)
}

func (ed *EventData) getAppStats(appId string) *eventApp.AppStats {

	appStats := ed.AppMap[appId]
	if appStats == nil {
		// New app we haven't seen yet
		appStats = eventApp.NewAppStats(appId)
		ed.AppMap[appId] = appStats
	}
	return appStats
}

func (ed *EventData) getCellStats(cellIp string) *eventCell.CellStats {
	cellStats := ed.CellMap[cellIp]
	if cellStats == nil {
		// New cell we haven't seen yet
		cellStats = eventCell.NewCellStats(cellIp)
		ed.CellMap[cellIp] = cellStats
	}

	// TODO: Is this the best place for this??
	if cellStats.StackId == "" {
		ed.assignStackId(cellStats)
	}
	if cellStats.IsolationSegmentGuid == "" {
		ed.AssignIsolationSegment(cellStats)
	}
	return cellStats
}

func (ed *EventData) assignStackId(cellStats *eventCell.CellStats) {
	// Look for a container running on cell to determine which stack the cell is running
	// TODO: This is not very efficient -- if end up with a cell that has no containers yet
	// this loop will run every time the cell metric comes in.
	for _, appStats := range ed.AppMap {
		for _, containerStats := range appStats.ContainerArray {
			//if containerStats != nil && cellStats.Ip == containerStats.Ip && space.All() != nil && len(space.All()) > 0 {
			if containerStats != nil && cellStats.Ip == containerStats.Ip {
				appMetadata := ed.eventProcessor.GetMetadataManager().GetAppMdManager().FindAppMetadata(appStats.AppId)
				cellStats.StackId = appMetadata.StackGuid
				return
			}
		}
	}
}

func (ed *EventData) AssignIsolationSegment(cellStats *eventCell.CellStats) {
	for _, appStats := range ed.AppMap {
		for _, containerStats := range appStats.ContainerArray {
			if containerStats != nil && cellStats.Ip == containerStats.Ip {

				//appMetadata := ed.eventProcessor.GetMetadataManager().GetAppMdManager().FindAppMetadata(appStats.AppId)
				ep := ed.eventProcessor
				mm := ep.GetMetadataManager()
				amdm := mm.GetAppMdManager()
				appMetadata := amdm.FindAppMetadata(appStats.AppId)
				spaceMetadata := space.FindSpaceMetadata(appMetadata.SpaceGuid)
				cellStats.IsolationSegmentGuid = spaceMetadata.IsolationSegmentGuid
				return
			}
		}
	}
}

func (ed *EventData) getContainerStats(appStats *eventApp.AppStats, instIndex int) *eventApp.ContainerStats {

	// Save the container data -- by instance id
	if len(appStats.ContainerArray) <= instIndex {
		caArray := make([]*eventApp.ContainerStats, instIndex+1)
		for i, ca := range appStats.ContainerArray {
			caArray[i] = ca
		}
		appStats.ContainerArray = caArray
	}

	containerStats := appStats.ContainerArray[instIndex]

	if containerStats == nil {
		// New app instance (container) we haven't seen yet
		containerStats = eventApp.NewContainerStats(instIndex)
		appStats.ContainerArray[instIndex] = containerStats

	}
	return containerStats
}

func formatUUID(uuid *events.UUID) string {
	if uuid == nil {
		return ""
	}
	var uuidBytes [16]byte
	binary.LittleEndian.PutUint64(uuidBytes[:8], uuid.GetLow())
	binary.LittleEndian.PutUint64(uuidBytes[8:], uuid.GetHigh())
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuidBytes[0:4], uuidBytes[4:6], uuidBytes[6:8], uuidBytes[8:10], uuidBytes[10:])
}
