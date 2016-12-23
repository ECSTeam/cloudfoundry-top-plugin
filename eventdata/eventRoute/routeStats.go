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

type RouteSlice []*RouteStats

type RouteStats struct {
	RouteId string

	// Key: appId
	AppRouteStatsMap map[string]*AppRouteStats
}

func NewRouteStats(routeId string) *RouteStats {
	stats := &RouteStats{}
	stats.RouteId = routeId
	stats.AppRouteStatsMap = make(map[string]*AppRouteStats)
	return stats
}

func (rs *RouteStats) Id() string {
	return rs.RouteId
}

func (rs *RouteStats) FindAppRouteStats(appId string) *AppRouteStats {
	return rs.AppRouteStatsMap[appId]
}
