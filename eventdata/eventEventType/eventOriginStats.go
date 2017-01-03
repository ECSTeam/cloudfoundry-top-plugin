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

//  EventOriginStats --> OrginStats --> EventDetailStats
type EventOriginStats struct {
	Origin string
	// Key: EventDetailMapKey
	EventDetailStatsMap map[EventDetailMapKey]*EventDetailStats
}

func NewEventOriginStats(origin string) *EventOriginStats {
	stats := &EventOriginStats{}
	stats.Origin = origin
	stats.EventDetailStatsMap = make(map[EventDetailMapKey]*EventDetailStats)
	return stats
}

func (os *EventOriginStats) Id() string {
	return os.Origin
}

func (os *EventOriginStats) FindEventDetailStats(msg *events.Envelope) *EventDetailStats {

	edKey := &EventDetailMapKey{
		Deployment: msg.GetDeployment(),
		Job:        msg.GetJob(),
		Index:      msg.GetIndex(),
		Ip:         msg.GetIp(),
	}

	eventDetailStats := os.EventDetailStatsMap[*edKey]
	if eventDetailStats == nil {
		eventDetailStats = NewEventDetailStats(msg)
		os.EventDetailStatsMap[*edKey] = eventDetailStats
	}
	return eventDetailStats
}
