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

import "github.com/cloudfoundry/sonde-go/events"

func (ed *EventData) valueMetricEvent(msg *events.Envelope) {

	// Can we assume that all rep orgins are cflinuxfs2 diego cells? Might be a bad idea
	// 12/28/2017 - Added garden-linux to support PCF 2.0 small footprint
	// Can't use "diego_cell" in job name as isolation segment job name could be custom job name
	if msg.GetOrigin() == "rep" || msg.GetOrigin() == "garden-linux" {
		ip := msg.GetIp()
		cellStats := ed.getCellStats(ip)

		cellStats.DeploymentName = msg.GetDeployment()
		cellStats.JobName = msg.GetJob()

		cellStats.JobIndex = msg.GetIndex()

		valueMetric := msg.GetValueMetric()
		value := ed.getMetricValue(valueMetric)
		name := valueMetric.GetName()
		switch name {
		case "numCPUS":
			// ISSUE: numCPUS - PCF 2.6 no longer sending metric "numCPUS"
			//cellStats.NumOfCpus = int(value)
		case "CapacityTotalMemory":
			cellStats.CapacityMemoryTotal = int64(value)
		case "CapacityRemainingMemory":
			cellStats.CapacityMemoryRemaining = int64(value)
		case "CapacityTotalDisk":
			cellStats.CapacityDiskTotal = int64(value)
		case "CapacityRemainingDisk":
			cellStats.CapacityDiskRemaining = int64(value)
		case "CapacityTotalContainers":
			cellStats.CapacityTotalContainers = int(value)
		case "CapacityRemainingContainers":
			cellStats.CapacityRemainingContainers = int(value)
		case "ContainerCount":
			cellStats.ContainerCount = int(value)
		}
	}

}

func (ed *EventData) getMetricValue(valueMetric *events.ValueMetric) float64 {

	value := valueMetric.GetValue()
	switch valueMetric.GetUnit() {
	case "KiB":
		value = value * 1024
	case "MiB":
		value = value * 1024 * 1024
	case "GiB":
		value = value * 1024 * 1024 * 1024
	case "TiB":
		value = value * 1024 * 1024 * 1024 * 1024
	}

	return value
}
