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

package dataCommon

import (
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventApp"
)

type DisplayAppStats struct {
	*eventApp.AppStats

	AppName              string
	AppNameForSort       string
	SpaceId              string
	SpaceName            string
	OrgId                string
	OrgName              string
	DesiredContainers    int
	StackId              string
	StackName            string
	IsolationSegmentGuid string
	IsolationSegmentName string

	// If the app is in the pending delete from cache, then its been deleted
	// but we keep it around for a n seconds to show newly deleted app on the UI
	IsDeleted             bool
	IsPackageStatePending bool

	// Indicate if this app is monitored.  For privileged users
	// this should always be true.
	Monitored bool

	//TotalTraffic *eventdata.TrafficStats

	TotalCpuPercentage float64
	TotalMemoryUsed    int64
	TotalDiskUsed      int64

	TotalReportingContainers int
	TotalLogStdout           int64
	TotalLogStderr           int64
	Crash1hCount             int
	Crash24hCount            int
	LastCrashTime            *time.Time

	// Summerize HTTP response codes
	HttpAllCount int64
	Http2xxCount int64
	Http3xxCount int64
	Http4xxCount int64
	Http5xxCount int64
}

func NewDisplayAppStats(appStats *eventApp.AppStats) *DisplayAppStats {
	stats := &DisplayAppStats{}
	stats.AppStats = appStats
	return stats
}
