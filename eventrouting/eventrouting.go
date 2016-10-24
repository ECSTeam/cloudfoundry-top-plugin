package eventrouting

import (
	//"fmt"
	//"github.com/Sirupsen/logrus"
	"github.com/cloudfoundry/sonde-go/events"
	//"os"
	//"sort"
	//"strings"
	"sync"
	//"time"
	"github.com/kkellner/cloudfoundry-top-plugin/appStats"
)

type EventRouting struct {
	//CachingClient       caching.Caching
	//selectedEvents      map[string]bool
	selectedEventsCount map[string]uint64
	mutex               *sync.Mutex
	//log                 logging.Logging
	//ExtraFields         map[string]string
	processor 					*appStats.AppStatsEventProcessor

}


func (e *EventRouting) RouteEvent(processor *appStats.AppStatsEventProcessor, msg *events.Envelope) {

	eventType := msg.GetEventType()

	// Check if this is an HttpStartStop event
	if (int)(eventType) == 4 {
		//fmt.Printf("event: %v\n", msg)
		processor.Process(msg)
	}

}

// func NewEventRouting(caching caching.Caching, logging logging.Logging) *EventRouting {
func NewEventRouting() *EventRouting {
	return &EventRouting{
		//CachingClient:       caching,
		//selectedEvents:      make(map[string]bool),
		selectedEventsCount: make(map[string]uint64),
		//log:                 logging,
		mutex:               &sync.Mutex{},
		//ExtraFields:         make(map[string]string),
	}
}
