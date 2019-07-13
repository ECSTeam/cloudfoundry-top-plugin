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

package eventCell

type CellStats struct {
	Ip string
	// TODO: SG - We can't assign a single stackId to a cell -- change to StackGroupId
	StackGroupId         string
	IsolationSegmentGuid string
	DeploymentName       string
	JobName              string
	JobIndex             string

	NumOfCpus                   int
	CapacityMemoryTotal         int64
	CapacityMemoryRemaining     int64
	CapacityDiskTotal           int64
	CapacityDiskRemaining       int64
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
