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

package eventrouting

import (
	"fmt"
	"sync"
	"time"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
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
	er.processor.ClearStats()
}

func (er *EventRouter) Route(instanceId int, msg *events.Envelope) {
	er.mu.Lock()
	defer er.mu.Unlock()
	er.eventCount++
	er.eventRateCounter.Incr()
	er.processor.Process(instanceId, msg)
}
