package eventrouting

import (
	"fmt"
	"sync"
	"time"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/kkellner/cloudfoundry-top-plugin/eventdata"
	"github.com/kkellner/cloudfoundry-top-plugin/toplog"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

type EventRouter struct {
	eventCount       int64
	eventRateCounter *util.RateCounter
	eventRatePeak    int
	startTime        time.Time
	mu               sync.Mutex
	processor        *eventdata.EventProcessor
}

func NewEventRouter(processor *eventdata.EventProcessor) *EventRouter {
	return &EventRouter{
		processor:        processor,
		startTime:        time.Now(),
		eventRateCounter: util.NewRateCounter(time.Second),
	}
}

func (er *EventRouter) GetProcessor() *eventdata.EventProcessor {
	return er.processor
}

func (er *EventRouter) GetEventCount() int64 {
	return er.eventCount
}

func (er *EventRouter) GetEventRatePeak() int {
	return er.eventRatePeak
}

func (er *EventRouter) GetEventRate() int {
	rate := er.eventRateCounter.Rate()
	if rate > er.eventRatePeak {
		er.eventRatePeak = rate
		toplog.Info(fmt.Sprintf("New event rate per second peak: %v", rate))
	}
	return rate
}

func (er *EventRouter) GetStartTime() time.Time {
	return er.startTime
}

func (er *EventRouter) Clear() {
	er.eventCount = 0
	er.eventRatePeak = 0
	er.startTime = time.Now()
}

func (er *EventRouter) Route(instanceId int, msg *events.Envelope) {
	er.mu.Lock()
	defer er.mu.Unlock()
	er.eventCount++
	er.eventRateCounter.Incr()
	er.processor.Process(instanceId, msg)
}
