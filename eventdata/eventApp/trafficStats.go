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

package eventApp

import "github.com/ecsteam/cloudfoundry-top-plugin/util"

type TrafficStats struct {
	ResponseL60Time    *util.AvgTracker
	AvgResponseL60Time float64 // updated after a clone of this object
	EventL60Rate       int     // updated after a clone of this object

	ResponseL10Time    *util.AvgTracker
	AvgResponseL10Time float64 // updated after a clone of this object
	EventL10Rate       int     // updated after a clone of this object

	ResponseL1Time    *util.AvgTracker
	AvgResponseL1Time float64 // updated after a clone of this object
	EventL1Rate       int     // updated after a clone of this object

	HttpAllCount int64
	Http2xxCount int64
	Http3xxCount int64
	Http4xxCount int64
	Http5xxCount int64
}

func NewTrafficStats() *TrafficStats {
	stats := &TrafficStats{}
	return stats
}
