package appStats

import (
	"fmt"
	"github.com/cloudfoundry/sonde-go/events"
  "encoding/binary"
)


type AppStatsEventProcessor struct {
  appMap      map[string]int64
  totalEvents int64
}


func NewAppStatsEventProcessor() *AppStatsEventProcessor {
  return &AppStatsEventProcessor {
    appMap:  make(map[string]int64),
    totalEvents: 0,
  }
}

func (ap *AppStatsEventProcessor) GetTotalEvents() int64 {
  return ap.totalEvents
}

func (ap *AppStatsEventProcessor) GetAppMap() map[string]int64 {
  return ap.appMap
}

func (ap *AppStatsEventProcessor) Clear() {
  ap.appMap = make(map[string]int64)
	ap.totalEvents = 0
}

func (ap *AppStatsEventProcessor) Process(msg *events.Envelope) {

  ap.totalEvents++

  appUUID := msg.GetHttpStartStop().GetApplicationId()
  instId := msg.GetHttpStartStop().GetInstanceId()

  // Check if this is an application event
  if appUUID != nil && instId != "" {

    appId := formatUUID(appUUID)
    //c.ui.Say("**** appId:%v ****", appId)

    count := ap.appMap[appId]
    count++
    ap.appMap[appId] = count

  }

}

func formatUUID(uuid *events.UUID) string {
	var uuidBytes [16]byte
	binary.LittleEndian.PutUint64(uuidBytes[:8], uuid.GetLow())
	binary.LittleEndian.PutUint64(uuidBytes[8:], uuid.GetHigh())
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuidBytes[0:4], uuidBytes[4:6], uuidBytes[6:8], uuidBytes[8:10], uuidBytes[10:])
}
