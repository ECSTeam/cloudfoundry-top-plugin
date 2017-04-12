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

package dataCommon

import (
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/config"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventrouting"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/app"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/isolationSegment"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/org"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/space"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/stack"
)

type CommonData struct {
	//masterUI       masterUIInterface.MasterUIInterface
	router *eventrouting.EventRouter
	//eventProcessor *eventdata.EventProcessor
	appMdMgr *app.AppMetadataManager

	displayAppStatsMap map[string]*DisplayAppStats

	isWarmupComplete bool
	// This is a count of the number of apps that do not have
	// the correct number of containers running based on app
	// instance setting
	appsNotInDesiredState int
	totalCrashCount       int
}

// TODO:  Create a common data struct -- which needs access to masterUI
// to get GetDisplayedEventData and GetAppMdMgr.  Also will allow
// appView to access this data through  masterUI so we don't process
// the same data twice

func NewCommonData(router *eventrouting.EventRouter) *CommonData {
	cd := &CommonData{router: router}

	cd.appMdMgr = router.GetProcessor().GetMetadataManager().GetAppMdManager()

	return cd
}

func (cd *CommonData) GetDisplayAppStatsMap() map[string]*DisplayAppStats {
	return cd.displayAppStatsMap
}

func (cd *CommonData) IsWarmupComplete() bool {
	return cd.isWarmupComplete
}

func (cd *CommonData) AppsNotInDesiredState() int {
	return cd.appsNotInDesiredState
}

func (cd *CommonData) TotalCrashCount() int {
	return cd.totalCrashCount
}

func (cd *CommonData) PostProcessData() map[string]*DisplayAppStats {

	eventData := cd.router.GetProcessor().GetDisplayedEventData()
	statsTime := eventData.StatsTime
	runtimeSeconds := statsTime.Sub(cd.router.GetStartTime())
	cd.isWarmupComplete = runtimeSeconds > time.Second*config.WarmUpSeconds

	//displayStatsArray := make([]*DisplayAppStats, 0)
	displayStatsMap := make(map[string]*DisplayAppStats)

	appMap := cd.router.GetProcessor().GetDisplayedEventData().AppMap
	//appStatsArray := eventApp.ConvertFromMap(appMap, cd.appMdMgr)
	appsNotInDesiredState := 0
	totalCrashCount := 0
	//now := time.Now()

	for appId, appStats := range appMap {
		displayAppStats := NewDisplayAppStats(appStats)

		displayStatsMap[appId] = displayAppStats
		//displayStatsArray = append(displayStatsArray, displayAppStats)
		appMetadata := cd.appMdMgr.FindAppMetadata(appStats.AppId)

		displayAppStats.AppName = appMetadata.Name
		displayAppStats.SpaceId = appMetadata.SpaceGuid

		spaceMetadata := space.FindSpaceMetadata(appMetadata.SpaceGuid)
		displayAppStats.SpaceName = spaceMetadata.Name

		displayAppStats.OrgId, displayAppStats.OrgName = org.FindBySpaceGuid(appMetadata.SpaceGuid)
		totalCpuPercentage := 0.0
		totalMemoryUsed := int64(0)
		totalDiskUsed := int64(0)
		totalReportingContainers := 0

		if appMetadata.State == "STARTED" {
			displayAppStats.DesiredContainers = int(appMetadata.Instances)
		}

		stack := stack.FindStackMetadata(appMetadata.StackGuid)
		displayAppStats.StackId = appMetadata.StackGuid
		displayAppStats.StackName = stack.Name

		isoSeg := isolationSegment.FindMetadata(spaceMetadata.IsolationSegmentGuid)
		displayAppStats.IsolationSegmentGuid = isoSeg.Guid
		displayAppStats.IsolationSegmentName = isoSeg.Name

		crashCount := 0
		var lastCrashTime *time.Time

		for containerIndex, cs := range appStats.ContainerArray {
			if cs != nil && cs.ContainerMetric != nil {

				crashCount = crashCount + cs.CrashCount()
				if cs.LastCrashTime() != nil {
					if lastCrashTime == nil || lastCrashTime.Before(*cs.LastCrashTime()) {
						lastCrashTime = cs.LastCrashTime()
					}
				}

				// If we haven't gotten a container update recently, ignore the old value
				if statsTime.Sub(cs.LastUpdate) > time.Second*config.StaleContainerSeconds {
					appStats.ContainerArray[containerIndex] = nil
					continue
				}
				totalCpuPercentage = totalCpuPercentage + *cs.ContainerMetric.CpuPercentage
				totalMemoryUsed = totalMemoryUsed + int64(*cs.ContainerMetric.MemoryBytes)
				totalDiskUsed = totalDiskUsed + int64(*cs.ContainerMetric.DiskBytes)
				totalReportingContainers++

			}
		}
		if totalReportingContainers < displayAppStats.DesiredContainers {
			appsNotInDesiredState = appsNotInDesiredState + 1
		}
		if totalReportingContainers > 0 {
			displayAppStats.TotalCpuPercentage = totalCpuPercentage
		} else {
			// In PCF 1.9 running containers can report 0.00 CPU percent usage
			// To help distiquish between a container with 0 CPU and no container
			// at all we set this to a very small negative number to help sort
			// no-container apps to the bottom when sorting by CPU%
			displayAppStats.TotalCpuPercentage = -0.0001
		}
		displayAppStats.TotalMemoryUsed = totalMemoryUsed
		displayAppStats.TotalDiskUsed = totalDiskUsed
		displayAppStats.TotalReportingContainers = totalReportingContainers
		displayAppStats.CrashCount = crashCount
		totalCrashCount = totalCrashCount + crashCount
		/*
			logStdoutCount := int64(0)
			logStderrCount := int64(0)
			for _, cs := range appStats.ContainerArray {
				if cs != nil {
					logStdoutCount = logStdoutCount + cs.OutCount
					logStderrCount = logStderrCount + cs.ErrCount
				}
			}
			displayAppStats.TotalLogStdout = logStdoutCount + appStats.NonContainerStdout
			displayAppStats.TotalLogStderr = logStderrCount + appStats.NonContainerStderr
		*/
	}
	cd.displayAppStatsMap = displayStatsMap
	cd.appsNotInDesiredState = appsNotInDesiredState
	cd.totalCrashCount = totalCrashCount
	return displayStatsMap
}
