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

import (
	"time"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/app"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/crashData"
)

const (
	UnknownName = "unknown"
)

type dataSlice []*AppStats

type AppStats struct {
	AppUUID *events.UUID
	AppId   string

	NonContainerStdout int64
	NonContainerStderr int64

	ContainerArray []*ContainerStats
	// Key: instanceId
	ContainerTrafficMap map[string]*TrafficStats

	// ISSUE: Must do this at clone time because of AvgTracker counter
	TotalTraffic *TrafficStats

	// We save container info in AppStats and not in ContainerStats
	// because container stats entries and come and go (scale up/down)
	// but we want to keep all crash info regardless of container status
	ContainerCrashInfo []*crashData.ContainerCrashInfo
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

func (as *AppStats) AddCrashInfo(containerIndex int, crashTime *time.Time, exitDescription string) {
	crashInfo := crashData.NewContainerCrashInfo(containerIndex, crashTime, exitDescription)
	if as.ContainerCrashInfo == nil {
		as.ContainerCrashInfo = make([]*crashData.ContainerCrashInfo, 0, 10)
	}
	as.ContainerCrashInfo = append(as.ContainerCrashInfo, crashInfo)
}

func (as *AppStats) CrashSince(since time.Duration) []*crashData.ContainerCrashInfo {

	crashInfoList := as.ContainerCrashInfo
	if crashInfoList != nil {
		sinceTime := time.Now().Add(since)
		crashInfoSize := len(crashInfoList)
		filteredCrashInfoList := make([]*crashData.ContainerCrashInfo, 0)
		for i := range crashInfoList {
			// Reverse loop through array
			crashInfo := crashInfoList[crashInfoSize-i-1]
			if crashInfo.CrashTime == nil || crashInfo.CrashTime.Before(sinceTime) {
				break
			}
			filteredCrashInfoList = append(filteredCrashInfoList, crashInfo)
		}
		return filteredCrashInfoList
	}
	return nil
}

// Crash count in last duration
func (as *AppStats) CrashCountSince(since time.Duration) int {
	crashInfoList := as.CrashSince(since)
	if crashInfoList != nil {
		return len(crashInfoList)
	}
	return 0
}

// Crash count in last 1 hour recorded since top started
func (as *AppStats) Crash1hCount() int {
	//return as.CrashCount(-1 * time.Minute)
	return as.CrashCountSince(-1 * time.Hour)
}

// Crash count in last 24 hours recorded since top started
func (as *AppStats) Crash24hCount() int {
	//return as.CrashCount(-2 * time.Minute)
	return as.CrashCountSince(-24 * time.Hour)
}
