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

import (
	"time"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

type HttpMethodSlice []*HttpMethodStats

type HttpMethodStats struct {
	Method     events.Method
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

	RequestCount int64

	// E.g., 200, 404, 500
	HttpStatusCode map[int32]int64

	// NOTE: is this realistic?? There could be unlimited number of RemoteAddresses
	// PCF 1.6 - All we have it remoteAddresses
	// PCF 1.7 - All we have it remoteAddresses as format of uri is messed up
	// PCF 1.8: "forwarded" (array) - pick the first/top forwarded value from HttpStartStop event array
	// NOTE: PCF 1.8 - remote address includes a port
	Forwarder map[string]int64

	ResponseContentLength int64
	// Not currently used
	RequestContentLength int64
}

func NewHttpMethodStats(httpMethod events.Method) *HttpMethodStats {
	stats := &HttpMethodStats{}
	stats.Method = httpMethod
	stats.HttpStatusCode = make(map[int32]int64)
	stats.Forwarder = make(map[string]int64)
	return stats
}
