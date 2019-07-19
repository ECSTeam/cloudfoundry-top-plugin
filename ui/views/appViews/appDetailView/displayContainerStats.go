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

package appDetailView

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventApp"
)

type DisplayContainerStats struct {
	*eventApp.ContainerStats
	*eventApp.AppStats

	AppName   string
	OrgName   string
	SpaceName string

	FreeMemory     uint64
	ReservedMemory uint64
	FreeDisk       uint64
	ReservedDisk   uint64

	State           string
	StateTime       *time.Time
	StateDuration   *time.Duration
	StartupDuration *time.Duration

	key string

	CrashCount int

	AvgResponseL60Time float64
	EventL60Rate       int
	//AvgResponseL10Time float64
	EventL10Rate int
	//AvgResponseL1Time float64
	EventL1Rate int

	HttpAllCount int64
	Http2xxCount int64
	Http3xxCount int64
	Http4xxCount int64
	Http5xxCount int64
}

func NewDisplayContainerStats(containerStats *eventApp.ContainerStats, appStats *eventApp.AppStats) *DisplayContainerStats {
	stats := &DisplayContainerStats{}
	stats.ContainerStats = containerStats
	stats.AppStats = appStats
	return stats
}

func (cs *DisplayContainerStats) Id() string {
	if cs.key == "" {
		// NOTE: Must include AppId and Index because this view is used by Diego cell view as well as App Detail view
		cs.key = fmt.Sprintf("%v-%v", cs.AppId, strconv.FormatInt(int64(cs.ContainerIndex), 10))
	}
	return cs.key
}
