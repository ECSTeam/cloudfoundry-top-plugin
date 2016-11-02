package appStats

import (

  "math"
  //"math/rand"
  //"time"
  "strconv"
	"fmt"
	"github.com/cloudfoundry/sonde-go/events"
  "encoding/binary"
  "github.com/mohae/deepcopy"
  //"github.com/paulbellamy/ratecounter"  // Uses a goroutine per call - not memory frendly
  "github.com/kkellner/cloudfoundry-top-plugin/debug"
)

type AppStatsEventProcessor struct {
  AppMap      map[string]*AppStats
  TotalEvents uint64
}

func NewAppStatsEventProcessor() *AppStatsEventProcessor {
  return &AppStatsEventProcessor {
    AppMap:  make(map[string]*AppStats),
    TotalEvents: 0,
  }
}

func (ap *AppStatsEventProcessor) Clone() *AppStatsEventProcessor {
  clone := deepcopy.Copy(ap).(*AppStatsEventProcessor)
  for _, stat := range ap.AppMap {

    clone.AppMap[stat.AppId].EventL60Rate = stat.responseL60Time.Rate()
    clone.AppMap[stat.AppId].AvgResponseL60Time = stat.responseL60Time.Avg()
    clone.AppMap[stat.AppId].EventL10Rate = stat.responseL10Time.Rate()
    clone.AppMap[stat.AppId].AvgResponseL10Time = stat.responseL10Time.Avg()
    clone.AppMap[stat.AppId].EventL1Rate = stat.responseL1Time.Rate()
    clone.AppMap[stat.AppId].AvgResponseL1Time = stat.responseL1Time.Avg()

  }
  return clone
}

func (ap *AppStatsEventProcessor) GetTotalEvents() uint64 {
  return ap.TotalEvents
}

func (ap *AppStatsEventProcessor) Clear() {
  ap.AppMap = make(map[string]*AppStats)
	ap.TotalEvents = 0
}

func (ap *AppStatsEventProcessor) Process(msg *events.Envelope) {

  eventType := msg.GetEventType()
switch eventType {
  case events.Envelope_HttpStartStop:
    ap.httpStartStopEvent(msg)
  case events.Envelope_ContainerMetric:
    ap.containerMetricEvent(msg)
  case events.Envelope_LogMessage:
    ap.logMessageEvent(msg)
  }


}

func (ap *AppStatsEventProcessor) logMessageEvent(msg *events.Envelope) {
  logMessage := msg.GetLogMessage()
  appId := logMessage.GetAppId()
  appStats := ap.AppMap[appId]
  if appStats == nil {
    // New app we haven't seen yet
    appStats = NewAppStats(appId)
    ap.AppMap[appId] = appStats
  }

  if logMessage.SourceInstance != nil {
    instanceIndex, err := strconv.Atoi(*logMessage.SourceInstance)
    if err==nil {
      // Save the metrics -- by instance id
      if len(appStats.LogMetric) <= instanceIndex {
        metricArray := make([]*LogMetric, instanceIndex+1)
        for i, metric := range appStats.LogMetric {
          metricArray[i] = metric
        }
        appStats.LogMetric = metricArray
      }


      logMetric := appStats.LogMetric[instanceIndex]
      if (logMetric == nil) {
        logMetric = &LogMetric {}
        appStats.LogMetric[instanceIndex] = logMetric
      }
      switch *logMessage.MessageType {
      case events.LogMessage_OUT:
        logMetric.OutCount++
      case events.LogMessage_ERR:
        logMetric.ErrCount++
      }
    }

  } else {
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
  appStats := ap.AppMap[appId]
  if appStats == nil {
    // New app we haven't seen yet
    appStats = NewAppStats(appId)
    ap.AppMap[appId] = appStats
  }

  // Save the container metrics -- by instance id
  if int32(len(appStats.ContainerMetric)) <= *containerMetric.InstanceIndex {
    cmArray := make([]*events.ContainerMetric, *containerMetric.InstanceIndex+1)
    for i, cm := range appStats.ContainerMetric {
      cmArray[i] = cm
    }
    appStats.ContainerMetric = cmArray
  }
  appStats.ContainerMetric[*containerMetric.InstanceIndex] = containerMetric

}


func (ap *AppStatsEventProcessor) httpStartStopEvent(msg *events.Envelope) {

  appUUID := msg.GetHttpStartStop().GetApplicationId()
  instId := msg.GetHttpStartStop().GetInstanceId()
  httpStartStopEvent := msg.GetHttpStartStop()
  if httpStartStopEvent.GetPeerType() == events.PeerType_Client &&
      appUUID != nil &&
      instId != "" {

    ap.TotalEvents++

    appId := formatUUID(appUUID)
    //c.ui.Say("**** appId:%v ****", appId)

    appStats := ap.AppMap[appId]
    if appStats == nil {
      // New app we haven't seen yet
      appStats = NewAppStats(appId)
      appStats.AppUUID = appUUID
      ap.AppMap[appId] = appStats
    }


    responseTimeMillis := *httpStartStopEvent.StopTimestamp - *httpStartStopEvent.StartTimestamp
    appStats.HttpAllCount++
    appStats.responseL60Time.Track(responseTimeMillis)
    appStats.responseL10Time.Track(responseTimeMillis)
    appStats.responseL1Time.Track(responseTimeMillis)

    statusCode := httpStartStopEvent.GetStatusCode()
    switch {
    case statusCode >= 200 && statusCode < 300:
      appStats.Http2xxCount++
    case statusCode >= 300 && statusCode < 400:
      appStats.Http3xxCount++
    case statusCode >= 400 && statusCode < 500:
      appStats.Http4xxCount++
    case statusCode >= 500 && statusCode < 600:
      appStats.Http5xxCount++
    }

  } else {
    statusCode := httpStartStopEvent.GetStatusCode()
    if statusCode == 4040 {
      debug.Debug(fmt.Sprintf("event:%v\n",msg))
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
	r := alpha * value + (1.0 - alpha) * oldValue
	return r
}
