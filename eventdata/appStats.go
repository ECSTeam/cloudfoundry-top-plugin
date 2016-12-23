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
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/app"
)

const (
	UnknownName = "unknown"
)

type dataSlice []*AppStats

type AppStats struct {
	AppUUID *events.UUID
	AppId   string
	//AppName string

	NonContainerStdout int64
	NonContainerStderr int64

	ContainerArray      []*ContainerStats
	ContainerTrafficMap map[string]*TrafficStats

	// ISSUE: Must do this at clone time because of AvgTracker counter
	TotalTraffic *TrafficStats
}

func NewAppStats(appId string) *AppStats {
	stats := &AppStats{}
	stats.AppId = appId
	return stats
}

func (as *AppStats) Id() string {
	return as.AppId
}

func ConvertFromMap(statsMap map[string]*AppStats, appMdMgr *app.AppMetadataManager) []*AppStats {
	s := make([]*AppStats, 0, len(statsMap))
	for _, d := range statsMap {
		s = append(s, d)
	}
	return s
}
