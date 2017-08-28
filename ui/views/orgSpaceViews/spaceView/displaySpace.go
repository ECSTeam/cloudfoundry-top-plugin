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

package spaceView

import "github.com/ecsteam/cloudfoundry-top-plugin/metadata/space"

type DisplaySpace struct {
	*space.SpaceMetadata

	QuotaName          string
	MemoryLimitInBytes int64

	NumberOfSpaces int
	NumberOfApps   int

	TotalCpuPercentage float64

	TotalMemoryReserved                    int64
	TotalMemoryUsed                        int64
	TotalMemoryReservedPercentOfSpaceQuota float64
	TotalMemoryReservedPercentOfOrgQuota   float64

	TotalDiskReserved int64
	TotalDiskUsed     int64

	DesiredContainers        int
	TotalReportingContainers int
	TotalLogStdout           int64
	TotalLogStderr           int64

	HttpAllCount int64

	IsolationSegmentName string
}

func NewDisplaySpace(spaceMetadata *space.SpaceMetadata) *DisplaySpace {
	stats := &DisplaySpace{}
	stats.SpaceMetadata = spaceMetadata
	return stats
}

func (do *DisplaySpace) Id() string {
	return do.Guid
}
