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

package eventRoute

import "github.com/cloudfoundry/sonde-go/events"

// Used as an overflow key when too many values are in map
// E.g., if too many values in UserAgent map, use OTHER bucket
const OTHER = "OTHER"

const MaxUserAgentBucket = 100

type AppRouteSlice []*AppRouteStats

type AppRouteStats struct {
	AppId string

	// E.g., GET, PUT, POST, DELETE
	HttpMethodStatsMap map[events.Method]*HttpMethodStats

	// Good idea??
	UserAgentMap map[string]int64
}

func NewAppRouteStats(appId string) *AppRouteStats {
	stats := &AppRouteStats{}
	stats.AppId = appId
	stats.HttpMethodStatsMap = make(map[events.Method]*HttpMethodStats)
	stats.UserAgentMap = make(map[string]int64)
	return stats
}

func (ars *AppRouteStats) FindHttpMethodStats(httpMethod events.Method) *HttpMethodStats {
	return ars.HttpMethodStatsMap[httpMethod]
}
