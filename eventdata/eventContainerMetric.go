// Copyright (c) 2017 ECS Team, Inc. - All Rights Reserved
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

package eventdata

import (
	"time"

	"github.com/cloudfoundry/sonde-go/events"
)

func (ed *EventData) containerMetricEvent(msg *events.Envelope) {

	containerMetric := msg.GetContainerMetric()

	appId := containerMetric.GetApplicationId()

	appStats := ed.getAppStats(appId)
	instNum := int(*containerMetric.InstanceIndex)
	containerStats := ed.getContainerStats(appStats, instNum)
	containerStats.LastUpdate = time.Now()
	containerStats.Ip = msg.GetIp()
	containerStats.ContainerMetric = containerMetric

}
