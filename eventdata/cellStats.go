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

package eventdata

type CellStats struct {
	Ip             string
	StackId        string
	DeploymentName string
	JobName        string
	JobIndex       string

	NumOfCpus                   int
	CapacityTotalMemory         int64
	CapacityRemainingMemory     int64
	CapacityTotalDisk           int64
	CapacityRemainingDisk       int64
	CapacityTotalContainers     int
	CapacityRemainingContainers int
	ContainerCount              int
}

func NewCellStats(cellIp string) *CellStats {
	stats := &CellStats{}
	stats.Ip = cellIp
	return stats
}

func (cs *CellStats) Id() string {
	return cs.Ip
}
