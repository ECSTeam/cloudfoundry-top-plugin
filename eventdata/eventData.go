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
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/mohae/deepcopy"
)

type EventData struct {
	AppMap         map[string]*AppStats
	CellMap        map[string]*CellStats
	TotalEvents    int64
	mu             *sync.Mutex
	logHttpAccess  *EventLogHttpAccess
	eventProcessor *EventProcessor
	apiUrl         string
	apiUrlRegexp   *regexp.Regexp
	urlPCF17Regexp *regexp.Regexp
}

func NewEventData(mu *sync.Mutex, eventProcessor *EventProcessor) *EventData {

	logHttpAccess := NewEventLogHttpAccess()
	apiEndpoint, err := eventProcessor.cliConnection.ApiEndpoint()
	if err != nil {
		toplog.Error("Call to ApiEndpoint failed:" + err.Error())
		apiEndpoint = "UNABLE_TO_GET_API_ENDPOINT"
	}

	url, err := url.Parse(apiEndpoint)
	if err != nil {
		toplog.Error("parse error:" + err.Error())
	}
	apiUrl := fmt.Sprintf("%v%v", url.Host, url.Path)

	apiUrlRegexpStr := `^\/v[0-9]\/([^\/]*)\/([^\/]*)`
	apiUrlRegexp := regexp.MustCompile(apiUrlRegexpStr)

	urlPCF17RegexpStr := `^http[s]?:\/\/[^\/]*(\/v[0-9]\/.*)$`
	urlPCF17Regexp := regexp.MustCompile(urlPCF17RegexpStr)

	return &EventData{
		AppMap:         make(map[string]*AppStats),
		CellMap:        make(map[string]*CellStats),
		TotalEvents:    0,
		mu:             mu,
		logHttpAccess:  logHttpAccess,
		eventProcessor: eventProcessor,
		apiUrl:         apiUrl,
		apiUrlRegexp:   apiUrlRegexp,
		urlPCF17Regexp: urlPCF17Regexp,
	}

}

func (ed *EventData) Process(instanceId int, msg *events.Envelope) {

	ed.mu.Lock()
	defer ed.mu.Unlock()

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
		if msg.CounterEvent.GetName() == "TruncatingBuffer.DroppedMessages" &&
			(msg.GetOrigin() == "DopplerServer" || msg.GetOrigin() == "doppler") {
			ed.droppedMessages(instanceId, msg)
		}
	}

}

func (ed *EventData) Clone() *EventData {

	ed.mu.Lock()
	defer ed.mu.Unlock()

	clone := deepcopy.Copy(ed).(*EventData)

	for _, appStat := range ed.AppMap {

		clonedAppStat := clone.AppMap[appStat.AppId]

		httpAllCount := int64(0)
		http2xxCount := int64(0)
		http3xxCount := int64(0)
		http4xxCount := int64(0)
		http5xxCount := int64(0)

		responseL60TimeArray := make([]*util.AvgTracker, 0)
		responseL10TimeArray := make([]*util.AvgTracker, 0)
		responseL1TimeArray := make([]*util.AvgTracker, 0)
		totalTraffic := NewTrafficStats()

		for instanceId, containerTraffic := range appStat.ContainerTrafficMap {

			rate60 := containerTraffic.responseL60Time.Rate()
			clone.AppMap[appStat.AppId].ContainerTrafficMap[instanceId].EventL60Rate = rate60
			totalTraffic.EventL60Rate = totalTraffic.EventL60Rate + rate60

			clone.AppMap[appStat.AppId].ContainerTrafficMap[instanceId].AvgResponseL60Time = containerTraffic.responseL60Time.Avg()

			rate10 := containerTraffic.responseL10Time.Rate()
			clone.AppMap[appStat.AppId].ContainerTrafficMap[instanceId].EventL10Rate = rate10
			totalTraffic.EventL10Rate = totalTraffic.EventL10Rate + rate10

			clone.AppMap[appStat.AppId].ContainerTrafficMap[instanceId].AvgResponseL10Time = containerTraffic.responseL10Time.Avg()

			rate1 := containerTraffic.responseL1Time.Rate()
			clone.AppMap[appStat.AppId].ContainerTrafficMap[instanceId].EventL1Rate = rate1
			totalTraffic.EventL1Rate = totalTraffic.EventL1Rate + rate1

			clone.AppMap[appStat.AppId].ContainerTrafficMap[instanceId].AvgResponseL1Time = containerTraffic.responseL1Time.Avg()

			httpAllCount = httpAllCount + containerTraffic.HttpAllCount
			http2xxCount = http2xxCount + containerTraffic.Http2xxCount
			http3xxCount = http3xxCount + containerTraffic.Http3xxCount
			http4xxCount = http4xxCount + containerTraffic.Http4xxCount
			http5xxCount = http5xxCount + containerTraffic.Http5xxCount

			responseL60TimeArray = append(responseL60TimeArray, containerTraffic.responseL60Time)
			responseL10TimeArray = append(responseL10TimeArray, containerTraffic.responseL10Time)
			responseL1TimeArray = append(responseL1TimeArray, containerTraffic.responseL1Time)

		}

		totalTraffic.AvgResponseL60Time = util.AvgMultipleTrackers(responseL60TimeArray)
		totalTraffic.AvgResponseL10Time = util.AvgMultipleTrackers(responseL10TimeArray)
		totalTraffic.AvgResponseL1Time = util.AvgMultipleTrackers(responseL1TimeArray)

		totalTraffic.HttpAllCount = httpAllCount
		totalTraffic.Http2xxCount = http2xxCount
		totalTraffic.Http3xxCount = http3xxCount
		totalTraffic.Http4xxCount = http4xxCount
		totalTraffic.Http5xxCount = http5xxCount
		clonedAppStat.TotalTraffic = totalTraffic

	}

	return clone
}

func (ed *EventData) GetTotalEvents() int64 {
	return ed.TotalEvents
}

func (ed *EventData) Clear() {

	ed.mu.Lock()
	defer ed.mu.Unlock()

	ed.AppMap = make(map[string]*AppStats)
	ed.CellMap = make(map[string]*CellStats)
	ed.TotalEvents = 0
}

func (ed *EventData) droppedMessages(instanceId int, msg *events.Envelope) {
	delta := msg.GetCounterEvent().GetDelta()
	total := msg.GetCounterEvent().GetTotal()
	text := fmt.Sprintf("Nozzle #%v - Upstream message indicates the nozzle or the TrafficController is not keeping up. Dropped delta: %v, total: %v",
		instanceId, delta, total)
	toplog.Error(text)
}

func (ed *EventData) logMessageEvent(msg *events.Envelope) {

	logMessage := msg.GetLogMessage()
	appId := logMessage.GetAppId()

	appStats := ed.getAppStats(appId) // Thread here at crash

	switch logMessage.GetSourceType() {
	case "APP":
		instNum, err := strconv.Atoi(*logMessage.SourceInstance)
		if err == nil {
			containerStats := ed.getContainerStats(appStats, instNum)
			switch *logMessage.MessageType {
			case events.LogMessage_OUT:
				containerStats.OutCount++
			case events.LogMessage_ERR:
				containerStats.ErrCount++
			}
		}
	case "API":
		// This is our notification that the state of an application may have changed
		// e.g., App was marked as STARTED or STOPPED (by a user)
		appMetadata := ed.eventProcessor.GetMetadataManager().GetAppMdManager().FindAppMetadata(appId)
		logText := string(logMessage.GetMessage())
		toplog.Debug(fmt.Sprintf("API event occured for app:%v name:%v msg: %v", appId, appMetadata.Name, logText))
		ed.eventProcessor.GetMetadataManager().RequestRefreshAppMetadata(appId)

	case "RTR":
		// Ignore router log messages
		// Turns out there is nothing useful in this message
		//logMsg := logMessage.GetMessage()
		//ed.handleHttpAccessLogLine(string(logMsg))
	case "HEALTH":
		// Ignore health check messages (TODO: Check sourceType of "crashed" messages)
	default:
		// Non-container log -- staging logs, router logs, etc
		switch *logMessage.MessageType {
		case events.LogMessage_OUT:
			appStats.NonContainerStdout++
		case events.LogMessage_ERR:
			appStats.NonContainerStderr++
		}
	}

}

func (ed *EventData) handleHttpAccessLogLine(logLine string) {
	ed.logHttpAccess.parseHttpAccessLogLine(logLine)
}

func (ed *EventData) valueMetricEvent(msg *events.Envelope) {

	// Can we assume that all rep orgins are cflinuxfs2 diego cells? Might be a bad idea
	if msg.GetOrigin() == "rep" {
		ip := msg.GetIp()
		cellStats := ed.getCellStats(ip)

		cellStats.DeploymentName = msg.GetDeployment()
		cellStats.JobName = msg.GetJob()

		jobIndex, err := strconv.Atoi(msg.GetIndex())
		if err == nil {
			cellStats.JobIndex = jobIndex
		} else {
			cellStats.JobIndex = -1
		}

		valueMetric := msg.GetValueMetric()
		value := ed.getMetricValue(valueMetric)
		name := valueMetric.GetName()
		switch name {
		case "numCPUS":
			cellStats.NumOfCpus = int(value)
		case "CapacityTotalMemory":
			cellStats.CapacityTotalMemory = int64(value)
		case "CapacityRemainingMemory":
			cellStats.CapacityRemainingMemory = int64(value)
		case "CapacityTotalDisk":
			cellStats.CapacityTotalDisk = int64(value)
		case "CapacityRemainingDisk":
			cellStats.CapacityRemainingDisk = int64(value)
		case "CapacityTotalContainers":
			cellStats.CapacityTotalContainers = int(value)
		case "CapacityRemainingContainers":
			cellStats.CapacityRemainingContainers = int(value)
		case "ContainerCount":
			cellStats.ContainerCount = int(value)
		}
	}

}

func (ed *EventData) getMetricValue(valueMetric *events.ValueMetric) float64 {

	value := valueMetric.GetValue()
	switch valueMetric.GetUnit() {
	case "KiB":
		value = value * 1024
	case "MiB":
		value = value * 1024 * 1024
	case "GiB":
		value = value * 1024 * 1024 * 1024
	case "TiB":
		value = value * 1024 * 1024 * 1024 * 1024
	}

	return value
}

func (ed *EventData) containerMetricEvent(msg *events.Envelope) {

	containerMetric := msg.GetContainerMetric()

	appId := containerMetric.GetApplicationId()

	appStats := ed.getAppStats(appId)
	instNum := int(*containerMetric.InstanceIndex)
	containerStats := ed.getContainerStats(appStats, instNum)
	containerStats.LastUpdate = time.Now()
	containerStats.Ip = msg.GetIp()
	containerStats.ContainerMetric = containerMetric

}

func (ed *EventData) getAppStats(appId string) *AppStats {

	appStats := ed.AppMap[appId]
	if appStats == nil {
		// New app we haven't seen yet
		appStats = NewAppStats(appId)
		ed.AppMap[appId] = appStats
	}
	return appStats
}

func (ed *EventData) getCellStats(cellIp string) *CellStats {
	cellStats := ed.CellMap[cellIp]
	if cellStats == nil {
		// New cell we haven't seen yet
		cellStats = NewCellStats(cellIp)
		ed.CellMap[cellIp] = cellStats
	}

	// TODO: Is this the best place for this??
	if cellStats.StackId == "" {
		// Look for a container running on cell to determine which stack the cell is running
		for _, appStats := range ed.AppMap {
			for _, containerStats := range appStats.ContainerArray {
				if containerStats != nil && cellStats.Ip == containerStats.Ip {
					appMetadata := ed.eventProcessor.GetMetadataManager().GetAppMdManager().FindAppMetadata(appStats.AppId)
					cellStats.StackId = appMetadata.StackGuid
					return cellStats
				}
			}
		}
	}
	return cellStats
}

func (ed *EventData) getContainerStats(appStats *AppStats, instIndex int) *ContainerStats {

	// Save the container data -- by instance id
	if len(appStats.ContainerArray) <= instIndex {
		caArray := make([]*ContainerStats, instIndex+1)
		for i, ca := range appStats.ContainerArray {
			caArray[i] = ca
		}
		appStats.ContainerArray = caArray
	}

	containerStats := appStats.ContainerArray[instIndex]

	if containerStats == nil {
		// New app we haven't seen yet
		containerStats = NewContainerStats(instIndex)
		appStats.ContainerArray[instIndex] = containerStats

	}
	return containerStats
}

func (ed *EventData) getContainerTraffic(appStats *AppStats, instId string) *TrafficStats {

	// Save the container data -- by instance id

	if appStats.ContainerTrafficMap == nil {
		appStats.ContainerTrafficMap = make(map[string]*TrafficStats)
	}

	containerTraffic := appStats.ContainerTrafficMap[instId]
	if containerTraffic == nil {
		containerTraffic = NewTrafficStats()
		appStats.ContainerTrafficMap[instId] = containerTraffic
		containerTraffic.responseL60Time = util.NewAvgTracker(time.Minute)
		containerTraffic.responseL10Time = util.NewAvgTracker(time.Second * 10)
		containerTraffic.responseL1Time = util.NewAvgTracker(time.Second)
	}

	return containerTraffic
}

func (ed *EventData) httpStartStopEvent(msg *events.Envelope) {

	appUUID := msg.GetHttpStartStop().GetApplicationId()
	instId := msg.GetHttpStartStop().GetInstanceId()
	//instIndex := msg.GetHttpStartStop().GetInstanceIndex()
	httpEvent := msg.GetHttpStartStop()
	peerType := httpEvent.GetPeerType()
	switch {
	case peerType == events.PeerType_Client:
		if appUUID != nil && instId != "" {
			ed.httpStartStopEventForApp(msg)
		} else if appUUID == nil && instId != "" {
			switch httpEvent.GetMethod() {
			case events.Method_PUT:
				fallthrough
			//case events.Method_POST:
			// TODO: Maybe this shouldn't be here as POST doesn't seem to include guid on URI
			//	fallthrough
			case events.Method_DELETE:
				// Check if we have a PCF API call
				uri := httpEvent.GetUri()

				isApiCall, apiUri := ed.checkIfApiCall_PCF1_6(uri)
				if !isApiCall {
					// PCF 1.7 needs more testing -- comment out this check for now
					//isApiCall, apiUri = ed.checkIfApiCall_PCF1_7(uri)
				}
				if isApiCall {
					ed.pcfApiHasBeenCalled(msg, apiUri)
				}
			}
		}
	default:
		// Ignore
	}

}

// Format of HttpStartStop in PCF 1.6
// origin:"router__0" eventType:HttpStartStop timestamp:1481942228218776855 deployment:"cf"
// job:"router-partition-6c9fddda6d386d1b5b54" index:"0" ip:"172.28.1.57"
// httpStartStop:<startTimestamp:1481942228172587068 stopTimestamp:1481942228218776855
// requestId:<low:8306249200620409206  high:5348636437173287548 >  peerType:Client method:PUT
// uri:"api.system.laba.ecsteam.io/v2/spaces/2a7f2b63-e3f9-4e26-a73c-d3dd0be4b77f"
// remoteAddress:"172.28.1.51:46187" userAgent:"Mozilla/5.0" statusCode:201 contentLength:1300
// instanceId:"a455303f-144b-4936-9483-a61bfa23b35b" 11:"xxxx" >
func (ed *EventData) checkIfApiCall_PCF1_6(uri string) (bool, string) {
	if strings.HasPrefix(uri, ed.apiUrl) {
		// A PCF API has been called
		return true, strings.TrimPrefix(uri, ed.apiUrl)
	}
	return false, ""
}

// Format of HttpStartStop in PCF 1.7
// origin:"gorouter" eventType:HttpStartStop timestamp:1481950482175868678 deployment:"cf"
// job:"router-partition-72c346932f9a11cd262e" index:"0" ip:"172.28.3.57"
// httpStartStop:<startTimestamp:1481950482140382655 stopTimestamp:1481950482175838683
// requestId:<low:4704681661434642403 high:14399846285434694512 >  peerType:Client method:PUT
// uri:"http://73.169.24.191, 172.28.3.51, 172.29.0.2, 172.28.3.51/v2/spaces/605f2e92-a311-4bf8-a37d-296b0a692e25"
// remoteAddress:"172.28.3.51:38274" userAgent:"Mozilla/5.0" statusCode:201 contentLength:1311
// instanceId:"66a0e28b-81d7-44d6-46a4-ed46822b9b1f" 11:"xxxx" >
//
// NOTE: the uri has lost its hostname part which was the indication we are calling an api in PCF 1.6
// so we need to make some guesses that if we don't have an applicationId property but have an instanceId
// we must be calling an API.  Is this a good assumption?  Not a very clean solution.
func (ed *EventData) checkIfApiCall_PCF1_7(uri string) (bool, string) {
	toplog.Debug(fmt.Sprintf("Check if PCF 1.7 API call:%v", uri))
	parsedData := ed.urlPCF17Regexp.FindAllStringSubmatch(uri, -1)
	if len(parsedData) != 1 {
		toplog.Debug("Not a PCF 1.7 API that we care about")
		return false, ""
	}
	dataArray := parsedData[0]
	if len(dataArray) != 2 {
		toplog.Warn(fmt.Sprintf("checkIfApiCall_PCF1_7>>Unable to parse uri: %v", uri))
		return false, ""
	}
	apiUri := dataArray[1]
	toplog.Debug(fmt.Sprintf("This is a PCF 1.7 API call:%v", apiUri))
	return true, apiUri
}

// A PCF API has been called -- use this to trigger reload of metadata if appropriate
// Example: "/v2/spaces/59cde607-2cda-4e20-ab30-cc779c4026b0"
func (ed *EventData) pcfApiHasBeenCalled(msg *events.Envelope, apiUri string) {
	toplog.Debug(fmt.Sprintf("API called: %v", apiUri))

	parsedData := ed.apiUrlRegexp.FindAllStringSubmatch(apiUri, -1)
	if len(parsedData) != 1 {
		toplog.Debug(fmt.Sprintf("pcfApiHasBeenCalled>>Unable to parse (parsedData size) apiUri: %v", apiUri))
		return
	}
	dataArray := parsedData[0]
	if len(dataArray) != 3 {
		toplog.Debug(fmt.Sprintf("pcfApiHasBeenCalled>>Unable to parse (dataArray size) apiUri: %v", apiUri))
		return
	}
	dataType := dataArray[1]
	guid := dataArray[2]
	ed.pcfApiHasBeenCalledReloadMetadata(dataType, guid)

}

func (ed *EventData) pcfApiHasBeenCalledReloadMetadata(dataType, guid string) {
	toplog.Debug(fmt.Sprintf("Data type:%v GUID:%v", dataType, guid))
	switch dataType {
	case "spaces":
		// TODO reload metadata
		toplog.Debug(fmt.Sprintf("Reload SPACE metadata for space with GUID:%v", guid))
	case "organizations":
		// TODO reload metadata
		toplog.Debug(fmt.Sprintf("Reload ORG metadata for org with GUID:%v", guid))
	default:
	}
}

func (ed *EventData) httpStartStopEventForApp(msg *events.Envelope) {

	appUUID := msg.GetHttpStartStop().GetApplicationId()
	instId := msg.GetHttpStartStop().GetInstanceId()
	//instIndex := msg.GetHttpStartStop().GetInstanceIndex()
	httpStartStopEvent := msg.GetHttpStartStop()

	//toplog.Debug(fmt.Sprintf("index: %v\n", instIndex))
	//toplog.Debug(fmt.Sprintf("index mem: %v\n", msg.GetHttpStartStop().InstanceIndex))
	//fmt.Printf("index: %v\n", instIndex)
	ed.TotalEvents++
	appId := formatUUID(appUUID)
	//c.ui.Say("**** appId:%v ****", appId)

	appStats := ed.getAppStats(appId)
	if appStats.AppUUID == nil {
		appStats.AppUUID = appUUID
	}

	containerTraffic := ed.getContainerTraffic(appStats, instId)

	responseTimeMillis := *httpStartStopEvent.StopTimestamp - *httpStartStopEvent.StartTimestamp
	containerTraffic.HttpAllCount++
	containerTraffic.responseL60Time.Track(responseTimeMillis)
	containerTraffic.responseL10Time.Track(responseTimeMillis)
	containerTraffic.responseL1Time.Track(responseTimeMillis)

	statusCode := httpStartStopEvent.GetStatusCode()
	switch {
	case statusCode >= 200 && statusCode < 300:
		containerTraffic.Http2xxCount++
	case statusCode >= 300 && statusCode < 400:
		containerTraffic.Http3xxCount++
	case statusCode >= 400 && statusCode < 500:
		containerTraffic.Http4xxCount++
	case statusCode >= 500 && statusCode < 600:
		containerTraffic.Http5xxCount++
	}

}

func formatUUID(uuid *events.UUID) string {
	var uuidBytes [16]byte
	binary.LittleEndian.PutUint64(uuidBytes[:8], uuid.GetLow())
	binary.LittleEndian.PutUint64(uuidBytes[8:], uuid.GetHigh())
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuidBytes[0:4], uuidBytes[4:6], uuidBytes[6:8], uuidBytes[8:10], uuidBytes[10:])
}
