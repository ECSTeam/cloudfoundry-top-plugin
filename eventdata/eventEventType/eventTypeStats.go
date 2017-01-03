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

package eventEventType

import (
	"github.com/cloudfoundry/sonde-go/events"
)

//type domainSlice []*EventTypeStats

//  EventTypeStats --> OrginStats --> EventDetailStats
type EventTypeStats struct {
	EventType events.Envelope_EventType
	// Key: Envelope_EventType
	EventOriginStatsMap map[string]*EventOriginStats
}

func NewEventTypeStats(eventType events.Envelope_EventType) *EventTypeStats {
	stats := &EventTypeStats{}
	stats.EventType = eventType
	stats.EventOriginStatsMap = make(map[string]*EventOriginStats)
	return stats
}

func (ets *EventTypeStats) Id() string {
	return ets.EventType.String()
}

func (ets *EventTypeStats) FindEventOriginStats(origin string) *EventOriginStats {
	originStats := ets.EventOriginStatsMap[origin]
	if originStats == nil {
		originStats = NewEventOriginStats(origin)
		ets.EventOriginStatsMap[origin] = originStats
	}
	return originStats
}
