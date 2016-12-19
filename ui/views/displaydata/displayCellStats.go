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

package displaydata

import "github.com/ecsteam/cloudfoundry-top-plugin/eventdata"

type DisplayCellStats struct {
	*eventdata.CellStats
	StackName                    string
	TotalContainerCpuPercentage  float64
	TotalContainerReservedMemory uint64
	TotalContainerUsedMemory     uint64
	TotalContainerReservedDisk   uint64
	TotalContainerUsedDisk       uint64
	TotalReportingContainers     int
	TotalLogOutCount             int64
	TotalLogErrCount             int64

	CapacityPlan0_5GMem int
	CapacityPlan1_0GMem int
	CapacityPlan1_5GMem int
	CapacityPlan2_0GMem int
	CapacityPlan2_5GMem int
	CapacityPlan3_0GMem int
	CapacityPlan3_5GMem int
	CapacityPlan4_0GMem int
}

func NewDisplayCellStats(cellStats *eventdata.CellStats) *DisplayCellStats {
	stats := &DisplayCellStats{}
	stats.CellStats = cellStats
	return stats
}
