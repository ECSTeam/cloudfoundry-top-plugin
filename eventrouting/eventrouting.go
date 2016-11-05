package eventrouting

import (
	//"fmt"
	"github.com/cloudfoundry/sonde-go/events"
	"sync"
	"time"
	"github.com/kkellner/cloudfoundry-top-plugin/appStats"
)

type EventRouter struct {
	eventCount 				int64
	startTime					time.Time
	mutex             *sync.Mutex
	processor 			  *appStats.AppStatsEventProcessor
	errors 						<-chan error
	messages 					<-chan *events.Envelope
}

func NewEventRouter(processor *appStats.AppStatsEventProcessor) *EventRouter {
	return &EventRouter {
		processor:					processor,
		startTime:				time.Now(),
	}

}

func (er *EventRouter) GetEventCount() int64 {
	return er.eventCount
}

func (er *EventRouter) GetStartTime() time.Time {
	return er.startTime
}

func (er *EventRouter) Clear() {
	er.eventCount = 0
	er.startTime = time.Now()
}


func (er *EventRouter) Route(msg *events.Envelope) {
	er.eventCount++
	//eventType := msg.GetEventType()
	er.processor.Process(msg)
}

func (er *EventRouter) GetTotalEventCount() int64 {
	return er.eventCount
}
