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

import (
	"time"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

// Used as an overflow key when too many values are in map
// E.g., if too many values in UserAgent map, use OTHER bucket
const OTHER = "OTHER"

const MaxUserAgentBucket = 100
const MaxRemoteAddressBucket = 100

type RouteSlice []*RouteStats

type RouteStats struct {
	RouteId string

	LastAccess time.Time

	responseL60Time    *util.AvgTracker
	AvgResponseL60Time float64 // updated after a clone of this object
	EventL60Rate       int     // updated after a clone of this object

	responseL10Time    *util.AvgTracker
	AvgResponseL10Time float64 // updated after a clone of this object
	EventL10Rate       int     // updated after a clone of this object

	responseL1Time    *util.AvgTracker
	AvgResponseL1Time float64 // updated after a clone of this object
	EventL1Rate       int     // updated after a clone of this object

	// E.g., 200, 404, 500
	HttpStatusCode map[int32]int64

	// E.g., GET, PUT, POST, DELETE
	HttpMethod map[events.Method]int64

	// NOTE: is this realistic?? There could be unlimited number of RemoteAddresses
	// PCF 1.7 x-forward???
	// PCF 1.8: "forwarded" (array)
	// NOTE: PCF 1.8 - remote address includes a port
	RemoteAddress map[string]int64

	// Good idea??
	UserAgent map[string]int64

	ResponseContentLength int64
	// Not currently used
	RequestContentLength int64

	// Key: GUID of application
	ApplicationId map[string]int64
}

func NewRouteStats(routeId string) *RouteStats {
	stats := &RouteStats{}
	stats.RouteId = routeId
	stats.HttpStatusCode = make(map[int32]int64)
	stats.HttpMethod = make(map[events.Method]int64)
	stats.RemoteAddress = make(map[string]int64)
	stats.UserAgent = make(map[string]int64)
	stats.ApplicationId = make(map[string]int64)
	return stats
}

func (ds *RouteStats) Id() string {
	return ds.RouteId
}
