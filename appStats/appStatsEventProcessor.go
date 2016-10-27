package appStats

import (
	"fmt"
	"github.com/cloudfoundry/sonde-go/events"
  "encoding/binary"
  "github.com/kkellner/cloudfoundry-top-plugin/debug"
)

type AppStatsEventProcessor struct {
  appMap      map[string]*AppStats
  totalEvents uint64
}

func NewAppStatsEventProcessor() *AppStatsEventProcessor {
  return &AppStatsEventProcessor {
    appMap:  make(map[string]*AppStats),
    totalEvents: 0,
  }
}

func (ap *AppStatsEventProcessor) GetTotalEvents() uint64 {
  return ap.totalEvents
}

func (ap *AppStatsEventProcessor) GetAppMap() map[string]*AppStats {
  return ap.appMap
}

func (ap *AppStatsEventProcessor) Clear() {
  ap.appMap = make(map[string]*AppStats)
	ap.totalEvents = 0
}

func (ap *AppStatsEventProcessor) Process(msg *events.Envelope) {

  eventType := msg.GetEventType()
switch eventType {
  case events.Envelope_HttpStartStop:
    ap.httpStartStopEvent(msg)
  case events.Envelope_ContainerMetric:
    ap.httpContainerMetric(msg)
  }


}

func (ap *AppStatsEventProcessor) httpContainerMetric(msg *events.Envelope) {

  containerMetric := msg.GetContainerMetric()

  appId := containerMetric.GetApplicationId()
  appStats := ap.appMap[appId]
  if appStats == nil {
    // New app we haven't seen yet
    appStats = &AppStats {
      AppId: appId,
    }
    ap.appMap[appId] = appStats
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

  if msg.GetHttpStartStop().GetPeerType() == events.PeerType_Client &&
      appUUID != nil &&
      instId != "" {

    ap.totalEvents++

    appId := formatUUID(appUUID)
    //c.ui.Say("**** appId:%v ****", appId)

    appStats := ap.appMap[appId]
    if appStats == nil {
      // New app we haven't seen yet
      appStats = &AppStats {
        AppId: appId,
        AppUUID: appUUID,
      }
      ap.appMap[appId] = appStats
    }
    appStats.EventCount++
    statusCode := msg.GetHttpStartStop().GetStatusCode()
    switch {
    case statusCode >= 200 && statusCode < 300:
      appStats.Event2xxCount++
    case statusCode >= 300 && statusCode < 400:
      appStats.Event3xxCount++
    case statusCode >= 400 && statusCode < 500:
      appStats.Event4xxCount++
    case statusCode >= 500 && statusCode < 600:
      appStats.Event5xxCount++
    }

  } else {
    statusCode := msg.GetHttpStartStop().GetStatusCode()
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
