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

import "github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventRoute"

type DisplayRouteStats struct {
	*eventRoute.RouteStats
	// TODO: Put display fields here

	// Includes [host].[domain]/[path]   what about port??
	RouteName string
	Host      string
	Domain    string
	Path      string

	HttpAllCount   int64
	Http2xxCount   int64
	Http3xxCount   int64
	Http4xxCount   int64
	Http5xxCount   int64
	HttpOtherCount int64

	HttpMethodGetCount    int64
	HttpMethodPostCount   int64
	HttpMethodPutCount    int64
	HttpMethodDeleteCount int64
	HttpMethodOtherCount  int64
}

func NewDisplayRouteStats(routeStats *eventRoute.RouteStats) *DisplayRouteStats {
	stats := &DisplayRouteStats{}
	stats.RouteStats = routeStats
	return stats
}
