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
	"bytes"
	"time"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

type EventDetailMapKey struct {
	Deployment string
	Job        string
	Index      string
	Ip         string
}

//  EventDetailStats --> OrginStats --> EventDetailStats
type EventDetailStats struct {
	EventType events.Envelope_EventType
	Origin    string

	DeploymentName string
	JobName        string
	JobIndex       string
	Ip             string
	MyId           string

	LastEventTime time.Time
	EventCount    int64
}

func NewEventDetailStats(msg *events.Envelope) *EventDetailStats {

	stats := &EventDetailStats{
		EventType:      msg.GetEventType(),
		Origin:         msg.GetOrigin(),
		DeploymentName: msg.GetDeployment(),
		JobName:        msg.GetJob(),
		JobIndex:       msg.GetIndex(),
		Ip:             msg.GetIp(),
	}

	var buffer bytes.Buffer
	buffer.WriteString(stats.DeploymentName)
	buffer.WriteString(stats.JobName)
	buffer.WriteString(stats.JobIndex)
	buffer.WriteString(stats.Ip)
	id := buffer.String()
	stats.MyId = id

	toplog.Debug("NewEventDetailStats - %+v", stats)

	return stats
}

func (ds *EventDetailStats) Id() string {
	return ds.MyId
}
