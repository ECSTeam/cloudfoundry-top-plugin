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
	"log"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata"
)

const (
	UnknownName = "unknown"
)

type dataSlice []*AppStats

type AppStats struct {
	AppUUID   *events.UUID
	AppId     string
	AppName   string
	SpaceName string
	OrgName   string

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

func PopulateNamesIfNeeded(appStats *AppStats) {
	appMetadata := metadata.FindAppMetadata(appStats.AppId)
	appName := appMetadata.Name
	if appName == "" {
		appName = appStats.AppId
	}
	appStats.AppName = appName

	var spaceMetadata metadata.Space
	spaceName := appStats.SpaceName
	if spaceName == "" || spaceName == UnknownName {
		spaceMetadata = metadata.FindSpaceMetadata(appMetadata.SpaceGuid)
		spaceName = spaceMetadata.Name
		if spaceName == "" {
			spaceName = UnknownName
		}
		appStats.SpaceName = spaceName
	}

	orgName := appStats.OrgName
	if orgName == "" || orgName == UnknownName {
		if &spaceMetadata == nil {
			spaceMetadata = metadata.FindSpaceMetadata(appMetadata.SpaceGuid)
		}
		orgMetadata := metadata.FindOrgMetadata(spaceMetadata.OrgGuid)
		orgName = orgMetadata.Name
		if orgName == "" {
			orgName = UnknownName
		}
		appStats.OrgName = orgName
	}

}

func PopulateNamesFromMap(statsMap map[string]*AppStats) []*AppStats {

	s := make([]*AppStats, 0, len(statsMap))
	for _, d := range statsMap {
		appMetadata := metadata.FindAppMetadata(d.AppId)
		// TODO: Take this check out
		if appMetadata == nil {
			log.Panic("appMetadata shouldn't be nil")
		}
		appName := appMetadata.Name
		if appName == "" {
			appName = d.AppId
		}
		d.AppName = appName

		var spaceMetadata metadata.Space
		spaceName := d.SpaceName
		if spaceName == "" || spaceName == UnknownName {
			spaceMetadata = metadata.FindSpaceMetadata(appMetadata.SpaceGuid)
			spaceName = spaceMetadata.Name
			if spaceName == "" {
				spaceName = UnknownName
			}
			d.SpaceName = spaceName
		}

		orgName := d.OrgName
		if orgName == "" || orgName == UnknownName {
			if &spaceMetadata == nil {
				spaceMetadata = metadata.FindSpaceMetadata(appMetadata.SpaceGuid)
			}
			orgMetadata := metadata.FindOrgMetadata(spaceMetadata.OrgGuid)
			orgName = orgMetadata.Name
			if orgName == "" {
				orgName = UnknownName
			}
			d.OrgName = orgName
		}
		s = append(s, d)
	}
	return s
}
