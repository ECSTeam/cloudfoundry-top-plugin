package eventrouting

import (
	"github.com/cloudfoundry/sonde-go/events"
	"sync"
	"github.com/kkellner/cloudfoundry-top-plugin/appStats"
)

type EventRouter struct {
	eventCount 				uint64
	mutex             *sync.Mutex
	processor 			  *appStats.AppStatsEventProcessor
}

func NewEventRouter(processor *appStats.AppStatsEventProcessor) *EventRouter {
	return &EventRouter {
		processor:					processor,
	}
}

func (e *EventRouter) Route(msg *events.Envelope) {

	e.eventCount++
	//eventType := msg.GetEventType()
	e.processor.Process(msg)

}

func (e *EventRouter) GetTotalEventCount() uint64 {
	return e.eventCount
}
