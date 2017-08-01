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
	"strconv"
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/config"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventrouting"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/app"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/appInstances"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/crashData"
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
	monitoredAppGuids  map[string]bool

	isWarmupComplete bool
	// This is a count of the number of apps that do not have
	// the correct number of containers running based on app
	// instance setting
	appsNotInDesiredState int
	totalCrash1hCount     int
	totalCrash24hCount    int
}

// TODO:  Create a common data struct -- which needs access to masterUI
// to get GetDisplayedEventData and GetAppMdMgr.  Also will allow
// appView to access this data through  masterUI so we don't process
// the same data twice

func NewCommonData(router *eventrouting.EventRouter, monitoredAppGuids map[string]bool) *CommonData {
	cd := &CommonData{router: router}

	cd.appMdMgr = router.GetProcessor().GetMetadataManager().GetAppMdManager()
	cd.monitoredAppGuids = monitoredAppGuids
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

func (cd *CommonData) TotalCrash1hCount() int {
	return cd.totalCrash1hCount
}

func (cd *CommonData) TotalCrash24hCount() int {
	return cd.totalCrash24hCount
}

func (cd *CommonData) SetMonitoredAppGuids(monitoredAppGuids map[string]bool) {
	cd.monitoredAppGuids = monitoredAppGuids
}

func (cd *CommonData) GetMonitoredAppGuids() map[string]bool {
	return cd.monitoredAppGuids
}

func (cd *CommonData) IsMonitoredAppGuid(appGuid string) bool {
	// Check if we're only monitoring a subset of apps and if this app is one of them
	if cd.monitoredAppGuids != nil {
		if cd.monitoredAppGuids[appGuid] {
			return true
		}
		return false
	}
	return true
}

func (cd *CommonData) PostProcessData() map[string]*DisplayAppStats {

	eventData := cd.router.GetProcessor().GetDisplayedEventData()
	statsTime := eventData.StatsTime
	runtimeSeconds := statsTime.Sub(cd.router.GetStartTime())
	cd.isWarmupComplete = runtimeSeconds > time.Second*config.WarmUpSeconds

	mdMgr := cd.router.GetProcessor().GetMetadataManager()
	displayStatsMap := make(map[string]*DisplayAppStats)

	appMap := cd.router.GetProcessor().GetDisplayedEventData().AppMap
	appsNotInDesiredState := 0
	totalCrash1hCount := 0
	totalCrash24hCount := 0

	for appId, appStats := range appMap {
		displayAppStats := NewDisplayAppStats(appStats)

		displayAppStats.Monitored = cd.IsMonitoredAppGuid(appId)

		displayStatsMap[appId] = displayAppStats
		appMetadata := cd.appMdMgr.FindAppMetadata(appStats.AppId)

		displayAppStats.AppName = appMetadata.Name
		if mdMgr.IsMonitorAppDetails(appId) {
			displayAppStats.AppName += "*"
		}
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

		// Crash count in last 1 hour (from call to /v2/events)
		crash1hCount := crashData.FindCountSinceByApp(appId, -1*time.Hour)
		crash1hCount = crash1hCount + appStats.Crash1hCount()

		// Crash count in last 24 hours (from call to /v2/events)
		crash24hCount := crashData.FindCountSinceByApp(appId, -24*time.Hour)
		crash24hCount = crash24hCount + appStats.Crash24hCount()

		for _, containerTraffic := range appStats.ContainerTrafficMap {
			for _, httpStatusCodeMap := range containerTraffic.HttpInfoMap {
				for statusCode, httpCountInfo := range httpStatusCodeMap {
					if httpCountInfo != nil {
						displayAppStats.HttpAllCount += httpCountInfo.HttpCount
						switch {
						case statusCode >= 200 && statusCode < 300:
							displayAppStats.Http2xxCount += httpCountInfo.HttpCount
						case statusCode >= 300 && statusCode < 400:
							displayAppStats.Http3xxCount += httpCountInfo.HttpCount
						case statusCode >= 400 && statusCode < 500:
							displayAppStats.Http4xxCount += httpCountInfo.HttpCount
						case statusCode >= 500 && statusCode < 600:
							displayAppStats.Http5xxCount += httpCountInfo.HttpCount
						}
					}
				}
			}
		}

		for _, cs := range appStats.ContainerArray {
			if cs != nil {

				appInsts := appInstances.FindAppInstancesMetadata(appId)
				if appInsts != nil {
					// If we have app instance metadata, lets check if the app is in a good state
					appInst := appInsts.Data[strconv.Itoa(cs.ContainerIndex)]
					if appInst == nil || appInst.State == "DOWN" || appInst.State == "CRASHED" {
						continue
					}
				}
				if cs.ContainerMetric != nil {
					totalReportingContainers++

					totalCpuPercentage = totalCpuPercentage + *cs.ContainerMetric.CpuPercentage
					totalMemoryUsed = totalMemoryUsed + int64(*cs.ContainerMetric.MemoryBytes)
					totalDiskUsed = totalDiskUsed + int64(*cs.ContainerMetric.DiskBytes)
				}
			}
		}
		if displayAppStats.Monitored && totalReportingContainers < displayAppStats.DesiredContainers {
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
		displayAppStats.Crash1hCount = crash1hCount
		displayAppStats.Crash24hCount = crash24hCount
		totalCrash1hCount = totalCrash1hCount + crash1hCount
		totalCrash24hCount = totalCrash24hCount + crash24hCount
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
	cd.totalCrash1hCount = totalCrash1hCount
	cd.totalCrash24hCount = totalCrash24hCount
	return displayStatsMap
}
