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

package ui

import (
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventApp"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/app"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/org"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/space"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/stack"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appViews/appView"
)

type CommonData struct {
	masterUI       *MasterUI
	eventProcessor *eventdata.EventProcessor
	appMdMgr       *app.AppMetadataManager

	displayAppStats []*appView.DisplayAppStats

	isWarmupComplete bool
	// This is a count of the number of apps that do not have
	// the correct number of containers running based on app
	// instance setting
	appsNotInDesiredState int
}

// TODO:  Create a common data struct -- which needs access to masterUI
// to get GetDisplayedEventData and GetAppMdMgr.  Also will allow
// appView to access this data through  masterUI so we don't process
// the same data twice

func NewCommonData(masterUI *MasterUI, eventProcessor *eventdata.EventProcessor) *CommonData {
	cd := &CommonData{masterUI: masterUI, eventProcessor: eventProcessor}

	cd.appMdMgr = eventProcessor.GetMetadataManager().GetAppMdManager()

	return cd
}

func (cd *CommonData) postProcessData() []*appView.DisplayAppStats {

	displayStatsArray := make([]*appView.DisplayAppStats, 0)
	appMap := cd.eventProcessor.GetDisplayedEventData().AppMap
	appStatsArray := eventApp.ConvertFromMap(appMap, cd.appMdMgr)
	appsNotInDesiredState := 0
	now := time.Now()

	for _, appStats := range appStatsArray {
		displayAppStats := appView.NewDisplayAppStats(appStats)
		displayStatsArray = append(displayStatsArray, displayAppStats)
		appMetadata := cd.appMdMgr.FindAppMetadata(appStats.AppId)

		displayAppStats.AppName = appMetadata.Name
		displayAppStats.SpaceName = space.FindSpaceName(appMetadata.SpaceGuid)
		displayAppStats.OrgName = org.FindOrgNameBySpaceGuid(appMetadata.SpaceGuid)

		totalCpuPercentage := 0.0
		totalUsedMemory := uint64(0)
		totalUsedDisk := uint64(0)
		totalReportingContainers := 0

		if appMetadata.State == "STARTED" {
			displayAppStats.DesiredContainers = int(appMetadata.Instances)
		}

		stack := stack.FindStackMetadata(appMetadata.StackGuid)
		displayAppStats.StackId = appMetadata.StackGuid
		displayAppStats.StackName = stack.Name

		for containerIndex, cs := range appStats.ContainerArray {
			if cs != nil && cs.ContainerMetric != nil {
				// If we haven't gotten a container update recently, ignore the old value
				if now.Sub(cs.LastUpdate) > time.Second*appView.StaleContainerSeconds {
					appStats.ContainerArray[containerIndex] = nil
					continue
				}
				totalCpuPercentage = totalCpuPercentage + *cs.ContainerMetric.CpuPercentage
				totalUsedMemory = totalUsedMemory + *cs.ContainerMetric.MemoryBytes
				totalUsedDisk = totalUsedDisk + *cs.ContainerMetric.DiskBytes
				totalReportingContainers++
			}
		}
		if totalReportingContainers < displayAppStats.DesiredContainers {
			appsNotInDesiredState = appsNotInDesiredState + 1
		}

		displayAppStats.TotalCpuPercentage = totalCpuPercentage
		displayAppStats.TotalUsedMemory = totalUsedMemory
		displayAppStats.TotalUsedDisk = totalUsedDisk
		displayAppStats.TotalReportingContainers = totalReportingContainers

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
	cd.displayAppStats = displayStatsArray
	cd.isWarmupComplete = cd.masterUI.IsWarmupComplete()
	cd.appsNotInDesiredState = appsNotInDesiredState
	return displayStatsArray
}
