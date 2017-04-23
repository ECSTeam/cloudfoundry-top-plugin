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

package appView

import (
	"fmt"
	"strconv"

	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/crashData"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/isolationSegment"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/dataCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func columnAppName() *uiCommon.ListColumn {
	defaultColSize := 50
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*dataCommon.DisplayAppStats).AppName, c2.(*dataCommon.DisplayAppStats).AppName)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return util.FormatDisplayData(appStats.AppName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return appStats.AppName
	}
	attentionFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
		appStats := data.(*dataCommon.DisplayAppStats)
		if !appStats.Monitored {
			return uiCommon.ATTENTION_NOT_MONITORED
		}
		attentionType := notInDesiredStateAttentionFunc(data, columnOwner)
		if attentionType == uiCommon.ATTENTION_NORMAL {
			attentionType = activityAttentionFunc(data, columnOwner)
		}
		return attentionType
	}
	c := uiCommon.NewListColumn("APPLICATION", "APPLICATION", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, attentionFunc)
	return c
}

func notInDesiredStateAttentionFunc(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
	appStats := data.(*dataCommon.DisplayAppStats)
	appListView := columnOwner.(*AppListView)
	if !appStats.Monitored {
		return uiCommon.ATTENTION_NOT_MONITORED
	}
	attentionType := uiCommon.ATTENTION_NORMAL
	if !appStats.Monitored {
		return attentionType
	}
	if appListView.isWarmupComplete && appStats.DesiredContainers > appStats.TotalReportingContainers {
		attentionType = uiCommon.ATTENTION_NOT_DESIRED_STATE
	}
	return attentionType
}

func activityAttentionFunc(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
	appStats := data.(*dataCommon.DisplayAppStats)
	if !appStats.Monitored {
		return uiCommon.ATTENTION_NOT_MONITORED
	}
	attentionType := uiCommon.ATTENTION_NORMAL
	if appStats.TotalTraffic.EventL10Rate > 0 {
		attentionType = uiCommon.ATTENTION_ACTIVITY
	}
	return attentionType
}

func notMonitoredAttentionFunc(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
	appStats := data.(*dataCommon.DisplayAppStats)
	if !appStats.Monitored {
		return uiCommon.ATTENTION_NOT_MONITORED
	}
	return uiCommon.ATTENTION_NORMAL
}

func columnSpaceName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*dataCommon.DisplayAppStats).SpaceName, c2.(*dataCommon.DisplayAppStats).SpaceName)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		spaceNameDisplay := "--"
		if appStats.SpaceName != "" {
			spaceNameDisplay = appStats.SpaceName
		}
		return util.FormatDisplayData(spaceNameDisplay, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return appStats.SpaceName
	}
	c := uiCommon.NewListColumn("SPACE", "SPACE", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}

func columnOrgName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*dataCommon.DisplayAppStats).OrgName, c2.(*dataCommon.DisplayAppStats).OrgName)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		orgName := "--"
		if appStats.OrgId != "" {
			orgName = appStats.OrgName
		}
		return util.FormatDisplayData(orgName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return appStats.OrgName
	}
	c := uiCommon.NewListColumn("ORG", "ORG", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}

func columnReportingContainers() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalReportingContainers < c2.(*dataCommon.DisplayAppStats).TotalReportingContainers
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		if appStats.Monitored {
			return fmt.Sprintf("%3v", appStats.TotalReportingContainers)
		} else {
			return fmt.Sprintf("%3v", "--")
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return strconv.Itoa(appStats.TotalReportingContainers)
	}
	attentionFunc := notInDesiredStateAttentionFunc
	c := uiCommon.NewListColumn("RCR", "RCR", 3,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, attentionFunc)
	return c
}

func columnDesiredInstances() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).DesiredContainers < c2.(*dataCommon.DisplayAppStats).DesiredContainers
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%3v", appStats.DesiredContainers)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return strconv.Itoa(appStats.DesiredContainers)
	}
	attentionFunc := notInDesiredStateAttentionFunc
	c := uiCommon.NewListColumn("DCR", "DCR", 3,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, attentionFunc)
	return c
}

func columnTotalCpu() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalCpuPercentage < c2.(*dataCommon.DisplayAppStats).TotalCpuPercentage
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		totalCpuInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalCpuInfo = fmt.Sprintf("%6v", "--")
		} else {
			if appStats.TotalCpuPercentage >= 100.0 {
				totalCpuInfo = fmt.Sprintf("%6.0f", appStats.TotalCpuPercentage)
			} else if appStats.TotalCpuPercentage >= 10.0 {
				totalCpuInfo = fmt.Sprintf("%6.1f", appStats.TotalCpuPercentage)
			} else {
				totalCpuInfo = fmt.Sprintf("%6.2f", appStats.TotalCpuPercentage)
			}
		}
		return fmt.Sprintf("%6v", totalCpuInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%.2f", appStats.TotalCpuPercentage)
	}
	c := uiCommon.NewListColumn("CPU_PER", "CPU%", 6,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}

func columnTotalMemoryUsed() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalMemoryUsed < c2.(*dataCommon.DisplayAppStats).TotalMemoryUsed
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		totalMemInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalMemInfo = fmt.Sprintf("%9v", "--")
		} else {
			totalMemInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.TotalMemoryUsed).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", totalMemInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalMemoryUsed)
	}
	c := uiCommon.NewListColumn("MEM_USED", "MEM_USED", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}

func columnTotalDiskUsed() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalDiskUsed < c2.(*dataCommon.DisplayAppStats).TotalDiskUsed
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		totalDiskInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalDiskInfo = fmt.Sprintf("%9v", "--")
		} else {
			totalDiskInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.TotalDiskUsed).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", totalDiskInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalDiskUsed)
	}
	c := uiCommon.NewListColumn("DSK_USED", "DSK_USED", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}

func columnAvgResponseTimeL60Info() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalTraffic.AvgResponseL60Time < c2.(*dataCommon.DisplayAppStats).TotalTraffic.AvgResponseL60Time
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		avgResponseTimeL60Info := "--"
		if appStats.TotalTraffic.AvgResponseL60Time >= 0 {
			avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL60Time / 1000000
			avgResponseTimeL60Info = fmt.Sprintf("%6.0f", avgResponseTimeMs)
		}
		return fmt.Sprintf("%6v", avgResponseTimeL60Info)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.AvgResponseL60Time)
	}
	c := uiCommon.NewListColumn("RESP", "RESP", 6,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}

func columnLogStdout() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalLogStdout < c2.(*dataCommon.DisplayAppStats).TotalLogStdout
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%11v", util.Format(appStats.TotalLogStdout))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalLogStdout)
	}
	c := uiCommon.NewListColumn("LOG_OUT", "LOG_OUT", 11,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}

func columnLogStderr() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalLogStderr < c2.(*dataCommon.DisplayAppStats).TotalLogStderr
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%11v", util.Format(appStats.TotalLogStderr))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalLogStderr)
	}
	c := uiCommon.NewListColumn("LOG_ERR", "LOG_ERR", 11,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}

func columnReq1() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalTraffic.EventL1Rate < c2.(*dataCommon.DisplayAppStats).TotalTraffic.EventL1Rate
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%6v", util.Format(int64(appStats.TotalTraffic.EventL1Rate)))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.EventL1Rate)
	}
	attentionFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
		appStats := data.(*dataCommon.DisplayAppStats)
		if !appStats.Monitored {
			return uiCommon.ATTENTION_NOT_MONITORED
		}
		attentionType := uiCommon.ATTENTION_NORMAL
		if appStats.TotalTraffic.EventL1Rate > 0 {
			attentionType = uiCommon.ATTENTION_ACTIVITY
		}
		return attentionType
	}
	c := uiCommon.NewListColumn("REQ1", "REQ/1", 6,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, attentionFunc)
	return c
}

func columnReq10() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalTraffic.EventL10Rate < c2.(*dataCommon.DisplayAppStats).TotalTraffic.EventL10Rate
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%7v", util.Format(int64(appStats.TotalTraffic.EventL10Rate)))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.EventL10Rate)
	}
	attentionFunc := activityAttentionFunc
	c := uiCommon.NewListColumn("REQ10", "REQ/10", 7,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, attentionFunc)
	return c
}

func columnReq60() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalTraffic.EventL60Rate < c2.(*dataCommon.DisplayAppStats).TotalTraffic.EventL60Rate
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%7v", util.Format(int64(appStats.TotalTraffic.EventL60Rate)))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.EventL60Rate)
	}
	attentionFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
		appStats := data.(*dataCommon.DisplayAppStats)
		if !appStats.Monitored {
			return uiCommon.ATTENTION_NOT_MONITORED
		}
		attentionType := uiCommon.ATTENTION_NORMAL
		if appStats.TotalTraffic.EventL60Rate > 0 {
			attentionType = uiCommon.ATTENTION_ACTIVITY
		}
		return attentionType
	}
	c := uiCommon.NewListColumn("REQ60", "REQ/60", 7,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, attentionFunc)
	return c
}

func columnTotalReq() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalTraffic.HttpAllCount < c2.(*dataCommon.DisplayAppStats).TotalTraffic.HttpAllCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.HttpAllCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.HttpAllCount)
	}
	c := uiCommon.NewListColumn("TOT_REQ", "TOT_REQ", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}
func column2XX() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalTraffic.Http2xxCount < c2.(*dataCommon.DisplayAppStats).TotalTraffic.Http2xxCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.Http2xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http2xxCount)
	}
	c := uiCommon.NewListColumn("2XX", "2XX", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}
func column3XX() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalTraffic.Http3xxCount < c2.(*dataCommon.DisplayAppStats).TotalTraffic.Http3xxCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.Http3xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http3xxCount)
	}
	c := uiCommon.NewListColumn("3XX", "3XX", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}

func column4XX() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalTraffic.Http4xxCount < c2.(*dataCommon.DisplayAppStats).TotalTraffic.Http4xxCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.Http4xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http4xxCount)
	}
	c := uiCommon.NewListColumn("4XX", "4XX", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}

func column5XX() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).TotalTraffic.Http5xxCount < c2.(*dataCommon.DisplayAppStats).TotalTraffic.Http5xxCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.Http5xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http5xxCount)
	}
	c := uiCommon.NewListColumn("5XX", "5XX", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}

func columnStackName() *uiCommon.ListColumn {
	defaultColSize := 15
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*dataCommon.DisplayAppStats).StackName, c2.(*dataCommon.DisplayAppStats).StackName)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return util.FormatDisplayData(appStats.StackName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return appStats.StackName
	}
	c := uiCommon.NewListColumn("STACK", "STACK", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}

func columnIsolationSegmentName() *uiCommon.ListColumn {
	defaultColSize := 15
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*dataCommon.DisplayAppStats).IsolationSegmentName, c2.(*dataCommon.DisplayAppStats).IsolationSegmentName)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		isolationSegmentName := "--"
		if appStats.IsolationSegmentGuid != isolationSegment.UnknownIsolationSegmentGuid {
			isolationSegmentName = appStats.IsolationSegmentName
		}
		return util.FormatDisplayData(isolationSegmentName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return appStats.IsolationSegmentName
	}
	c := uiCommon.NewListColumn("ISO_SEG", "ISO_SEG", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, notMonitoredAttentionFunc)
	return c
}

func columnCrashCount() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*dataCommon.DisplayAppStats).Crash24hCount < c2.(*dataCommon.DisplayAppStats).Crash24hCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*dataCommon.DisplayAppStats)
		crashCount := stats.Crash24hCount
		if crashCount > 0 || crashData.IsCacheLoaded() {
			display := fmt.Sprintf("%4v", util.Format(int64(crashCount)))
			return display
		} else {
			return fmt.Sprintf("%4v", "--")
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*dataCommon.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.Crash24hCount)
	}
	attentionFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
		appStats := data.(*dataCommon.DisplayAppStats)
		if !appStats.Monitored {
			return uiCommon.ATTENTION_NOT_MONITORED
		}
		attentionType := uiCommon.ATTENTION_NORMAL
		if appStats.Crash24hCount > 0 {
			attentionType = uiCommon.ATTENTION_WARM
		}
		return attentionType
	}
	c := uiCommon.NewListColumn("CRH", "CRH", 4,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, attentionFunc)
	return c
}
