package eventrouting

import (
	//"fmt"
	"sync"
	"time"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/kkellner/cloudfoundry-top-plugin/appStats"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

type EventRouter struct {
	eventCount       int64
	eventRateCounter *util.RateCounter

	startTime time.Time
	mu        sync.Mutex
	processor *appStats.AppStatsEventProcessor
}

func NewEventRouter(processor *appStats.AppStatsEventProcessor) *EventRouter {
	return &EventRouter{
		processor:        processor,
		startTime:        time.Now(),
		eventRateCounter: util.NewRateCounter(time.Second),
	}
}

func (er *EventRouter) GetEventCount() int64 {
	return er.eventCount
}

func (er *EventRouter) GetEventRate() int {
	return er.eventRateCounter.Rate()
}

func (er *EventRouter) GetStartTime() time.Time {
	return er.startTime
}

func (er *EventRouter) Clear() {
	er.eventCount = 0
	er.startTime = time.Now()
}

func (er *EventRouter) Route(instanceId int, msg *events.Envelope) {
	er.mu.Lock()
	defer er.mu.Unlock()
	er.eventCount++
	er.eventRateCounter.Incr()
	er.processor.Process(instanceId, msg)
}
