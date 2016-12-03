package appStats

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/kkellner/cloudfoundry-top-plugin/debug"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
	"github.com/mohae/deepcopy"
)

type AppStatsEventProcessor struct {
	AppMap      map[string]*AppStats
	TotalEvents int64
	mu          sync.Mutex
}

func NewAppStatsEventProcessor() *AppStatsEventProcessor {
	return &AppStatsEventProcessor{
		AppMap:      make(map[string]*AppStats),
		TotalEvents: 0,
	}
}

func (ap *AppStatsEventProcessor) Clone() *AppStatsEventProcessor {

	ap.mu.Lock()
	defer ap.mu.Unlock()

	clone := deepcopy.Copy(ap).(*AppStatsEventProcessor)

	for _, appStat := range ap.AppMap {

		clonedAppStat := clone.AppMap[appStat.AppId]

		httpAllCount := int64(0)
		http2xxCount := int64(0)
		http3xxCount := int64(0)
		http4xxCount := int64(0)
		http5xxCount := int64(0)

		//trafficMapSize := len(appStat.ContainerTrafficMap)
		responseL60TimeArray := make([]*util.AvgTracker, 0)
		responseL10TimeArray := make([]*util.AvgTracker, 0)
		responseL1TimeArray := make([]*util.AvgTracker, 0)
		totalTraffic := NewTraffic()

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

			//fmt.Printf("\n **** instanceId: %v\n", instanceId)

		}

		totalCpuPercentage := 0.0
		totalUsedMemory := uint64(0)
		totalUsedDisk := uint64(0)
		totalReportingContainers := 0

		now := time.Now()
		for containerIndex, cs := range appStat.ContainerArray {
			if cs != nil && cs.ContainerMetric != nil {
				// If we haven't gotten a container update recently, ignore the old value
				if now.Sub(cs.LastUpdate) > time.Second*40 {
					clonedAppStat.ContainerArray[containerIndex] = nil
					continue
				}

				totalCpuPercentage = totalCpuPercentage + *cs.ContainerMetric.CpuPercentage
				totalUsedMemory = totalUsedMemory + *cs.ContainerMetric.MemoryBytes
				totalUsedDisk = totalUsedDisk + *cs.ContainerMetric.DiskBytes
				totalReportingContainers++
			}
		}
		clonedAppStat.TotalCpuPercentage = totalCpuPercentage
		clonedAppStat.TotalUsedMemory = totalUsedMemory
		clonedAppStat.TotalUsedDisk = totalUsedDisk
		clonedAppStat.TotalReportingContainers = totalReportingContainers

		logCount := int64(0)
		for _, cs := range appStat.ContainerArray {
			if cs != nil {
				logCount = logCount + cs.OutCount + cs.ErrCount
			}
		}

		clonedAppStat.TotalLogCount = logCount + appStat.NonContainerOutCount + appStat.NonContainerErrCount

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

func (ap *AppStatsEventProcessor) GetTotalEvents() int64 {
	return ap.TotalEvents
}

func (ap *AppStatsEventProcessor) Clear() {

	ap.mu.Lock()
	defer ap.mu.Unlock()

	ap.AppMap = make(map[string]*AppStats)
	ap.TotalEvents = 0
}

func (ap *AppStatsEventProcessor) Process(instanceId int, msg *events.Envelope) {

	ap.mu.Lock()
	defer ap.mu.Unlock()

	eventType := msg.GetEventType()
	switch eventType {
	case events.Envelope_HttpStartStop:
		ap.httpStartStopEvent(msg)
	case events.Envelope_ContainerMetric:
		ap.containerMetricEvent(msg)
	case events.Envelope_LogMessage:
		ap.logMessageEvent(msg)
	case events.Envelope_CounterEvent:
		if msg.CounterEvent.GetName() == "TruncatingBuffer.DroppedMessages" && msg.GetOrigin() == "DopplerServer" {
			ap.droppedMessages(instanceId, msg)
		}
	}

}

func (ap *AppStatsEventProcessor) droppedMessages(instanceId int, msg *events.Envelope) {
	delta := msg.GetCounterEvent().GetDelta()
	total := msg.GetCounterEvent().GetTotal()
	text := fmt.Sprintf("Nozzle #%v - Upstream message indicates the nozzle or the TrafficController is not keeping up. Dropped delta: %v, total: %v",
		instanceId, delta, total)
	debug.Error(text)
}

func (ap *AppStatsEventProcessor) logMessageEvent(msg *events.Envelope) {

	logMessage := msg.GetLogMessage()
	appId := logMessage.GetAppId()

	appStats := ap.getAppStats(appId)

	switch logMessage.GetSourceType() {
	case "APP":
		instNum, err := strconv.Atoi(*logMessage.SourceInstance)
		if err == nil {
			containerStats := ap.getContainerStats(appStats, instNum)
			switch *logMessage.MessageType {
			case events.LogMessage_OUT:
				containerStats.OutCount++
			case events.LogMessage_ERR:
				containerStats.ErrCount++
			}
		}
	default:
		// Non-container -- staging logs?
		switch *logMessage.MessageType {
		case events.LogMessage_OUT:
			appStats.NonContainerOutCount++
		case events.LogMessage_ERR:
			appStats.NonContainerErrCount++
		}
	}

}

func (ap *AppStatsEventProcessor) containerMetricEvent(msg *events.Envelope) {

	containerMetric := msg.GetContainerMetric()

	appId := containerMetric.GetApplicationId()

	appStats := ap.getAppStats(appId)
	instNum := int(*containerMetric.InstanceIndex)
	containerStats := ap.getContainerStats(appStats, instNum)
	containerStats.LastUpdate = time.Now()
	containerStats.ContainerMetric = containerMetric

}

func (ap *AppStatsEventProcessor) getAppStats(appId string) *AppStats {
	appStats := ap.AppMap[appId]
	if appStats == nil {
		// New app we haven't seen yet
		appStats = NewAppStats(appId)
		ap.AppMap[appId] = appStats
	}
	return appStats
}

func (ap *AppStatsEventProcessor) getContainerStats(appStats *AppStats, instIndex int) *ContainerStats {

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
		containerStats = NewContainerStats()
		appStats.ContainerArray[instIndex] = containerStats

	}
	return containerStats
}

func (ap *AppStatsEventProcessor) getContainerTraffic(appStats *AppStats, instId string) *Traffic {

	// Save the container data -- by instance id

	if appStats.ContainerTrafficMap == nil {
		appStats.ContainerTrafficMap = make(map[string]*Traffic)
	}

	containerTraffic := appStats.ContainerTrafficMap[instId]
	if containerTraffic == nil {
		containerTraffic = NewTraffic()
		appStats.ContainerTrafficMap[instId] = containerTraffic
		containerTraffic.responseL60Time = util.NewAvgTracker(time.Minute)
		containerTraffic.responseL10Time = util.NewAvgTracker(time.Second * 10)
		containerTraffic.responseL1Time = util.NewAvgTracker(time.Second)
	}

	return containerTraffic
}

func (ap *AppStatsEventProcessor) httpStartStopEvent(msg *events.Envelope) {

	appUUID := msg.GetHttpStartStop().GetApplicationId()
	instId := msg.GetHttpStartStop().GetInstanceId()
	//instIndex := msg.GetHttpStartStop().GetInstanceIndex()
	httpStartStopEvent := msg.GetHttpStartStop()
	if httpStartStopEvent.GetPeerType() == events.PeerType_Client &&
		appUUID != nil &&
		instId != "" {

		//debug.Debug(fmt.Sprintf("index: %v\n", instIndex))
		//debug.Debug(fmt.Sprintf("index mem: %v\n", msg.GetHttpStartStop().InstanceIndex))

		//fmt.Printf("index: %v\n", instIndex)

		ap.TotalEvents++
		appId := formatUUID(appUUID)
		//c.ui.Say("**** appId:%v ****", appId)

		appStats := ap.getAppStats(appId)
		if appStats.AppUUID == nil {
			appStats.AppUUID = appUUID
		}

		containerTraffic := ap.getContainerTraffic(appStats, instId)

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

	} else {
		statusCode := httpStartStopEvent.GetStatusCode()
		if statusCode == 4040 {
			debug.Debug(fmt.Sprintf("event:%v\n", msg))
		}
	}
}

func formatUUID(uuid *events.UUID) string {
	var uuidBytes [16]byte
	binary.LittleEndian.PutUint64(uuidBytes[:8], uuid.GetLow())
	binary.LittleEndian.PutUint64(uuidBytes[8:], uuid.GetHigh())
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuidBytes[0:4], uuidBytes[4:6], uuidBytes[6:8], uuidBytes[8:10], uuidBytes[10:])
}

func MovingExpAvg(value, oldValue, fdtime, ftime float64) float64 {
	alpha := 1.0 - math.Exp(-fdtime/ftime)
	r := alpha*value + (1.0-alpha)*oldValue
	return r
}
