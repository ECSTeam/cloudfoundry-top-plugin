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

package orgView

import "github.com/ecsteam/cloudfoundry-top-plugin/metadata/org"

type DisplayOrg struct {
	*org.Org

	QuotaName          string
	MemoryLimitInBytes int64

	NumberOfSpaces int
	NumberOfApps   int

	TotalCpuPercentage float64

	TotalReservedMemory               int64
	TotalUsedMemory                   int64
	TotalReservedMemoryPercentOfQuota float64

	TotalReservedDisk int64
	TotalUsedDisk     int64

	TotalReportingContainers int
	TotalLogStdout           int64
	TotalLogStderr           int64

	HttpAllCount int64
}

func NewDisplayOrg(Org *org.Org) *DisplayOrg {
	stats := &DisplayOrg{}
	stats.Org = Org
	return stats
}

func (do *DisplayOrg) Id() string {
	return do.Guid
}
