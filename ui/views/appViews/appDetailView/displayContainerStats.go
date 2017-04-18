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
	"strconv"

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
	key            string
}

func NewDisplayContainerStats(containerStats *eventApp.ContainerStats, appStats *eventApp.AppStats) *DisplayContainerStats {
	stats := &DisplayContainerStats{}
	stats.ContainerStats = containerStats
	stats.AppStats = appStats
	return stats
}

func (cs *DisplayContainerStats) Id() string {
	if cs.key == "" {
		//cs.key = cs.AppId + "-" + strconv.FormatInt(int64(cs.ContainerIndex), 10)
		cs.key = strconv.FormatInt(int64(cs.ContainerIndex), 10)
	}
	return cs.key
}
