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

package routeView

import (
	"fmt"
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventRoute"
)

type DisplayRouteStats struct {
	*eventRoute.RouteStats

	// Includes [host].[domain]/[path]   what about port??
	RouteName string
	Host      string
	Domain    string
	Path      string
	Port      int

	RoutedAppCount int

	LastAccess            time.Time
	ResponseContentLength int64

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

func NewDisplayRouteStats(routeStats *eventRoute.RouteStats, hostName string, domainName string, pathName string, port int) *DisplayRouteStats {
	stats := &DisplayRouteStats{}
	stats.RouteStats = routeStats

	stats.RouteName = fmt.Sprintf("%v.%v%v", hostName, domainName, pathName)
	stats.Host = hostName
	stats.Domain = domainName
	stats.Path = pathName
	stats.Port = port

	return stats
}
