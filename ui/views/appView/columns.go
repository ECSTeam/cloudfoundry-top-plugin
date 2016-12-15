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

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func (asUI *AppListView) columnAppName() *uiCommon.ListColumn {
	defaultColSize := 50
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*displaydata.DisplayAppStats).AppName, c2.(*displaydata.DisplayAppStats).AppName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return util.FormatDisplayData(appStats.AppName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return appStats.AppName
	}
	c := uiCommon.NewListColumn("appName", "APPLICATION", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnSpaceName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*displaydata.DisplayAppStats).SpaceName, c2.(*displaydata.DisplayAppStats).SpaceName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return util.FormatDisplayData(appStats.SpaceName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return appStats.SpaceName
	}
	c := uiCommon.NewListColumn("spaceName", "SPACE", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnOrgName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*displaydata.DisplayAppStats).OrgName, c2.(*displaydata.DisplayAppStats).OrgName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return util.FormatDisplayData(appStats.OrgName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return appStats.OrgName
	}
	c := uiCommon.NewListColumn("orgName", "ORG", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnReportingContainers() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalReportingContainers < c2.(*displaydata.DisplayAppStats).TotalReportingContainers
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%3v", appStats.TotalReportingContainers)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return strconv.Itoa(appStats.TotalReportingContainers)
	}
	c := uiCommon.NewListColumn("reportingContainers", "RCR", 3,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnDesiredInstances() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).DesiredContainers < c2.(*displaydata.DisplayAppStats).DesiredContainers
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%3v", appStats.DesiredContainers)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return strconv.Itoa(appStats.DesiredContainers)
	}
	c := uiCommon.NewListColumn("desiredInstances", "DCR", 3,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnTotalCpu() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalCpuPercentage < c2.(*displaydata.DisplayAppStats).TotalCpuPercentage
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
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
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%.2f", appStats.TotalCpuPercentage)
	}
	c := uiCommon.NewListColumn("CPU", "CPU%", 6,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnTotalMemoryUsed() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalUsedMemory < c2.(*displaydata.DisplayAppStats).TotalUsedMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		totalMemInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalMemInfo = fmt.Sprintf("%9v", "--")
		} else {
			totalMemInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.TotalUsedMemory).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", totalMemInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalUsedMemory)
	}
	c := uiCommon.NewListColumn("MEM", "MEM", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnTotalDiskUsed() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalUsedDisk < c2.(*displaydata.DisplayAppStats).TotalUsedDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		totalDiskInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalDiskInfo = fmt.Sprintf("%9v", "--")
		} else {
			totalDiskInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.TotalUsedDisk).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", totalDiskInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalUsedDisk)
	}
	c := uiCommon.NewListColumn("DISK", "DISK", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnAvgResponseTimeL60Info() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalTraffic.AvgResponseL60Time < c2.(*displaydata.DisplayAppStats).TotalTraffic.AvgResponseL60Time
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		avgResponseTimeL60Info := "--"
		if appStats.TotalTraffic.AvgResponseL60Time >= 0 {
			avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL60Time / 1000000
			avgResponseTimeL60Info = fmt.Sprintf("%6.0f", avgResponseTimeMs)
		}
		return fmt.Sprintf("%6v", avgResponseTimeL60Info)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.AvgResponseL60Time)
	}
	c := uiCommon.NewListColumn("avgResponseTimeL60", "RESP", 6,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnLogStdout() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalLogStdout < c2.(*displaydata.DisplayAppStats).TotalLogStdout
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%11v", util.Format(appStats.TotalLogStdout))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalLogStdout)
	}
	c := uiCommon.NewListColumn("TotalLogStdout", "LOG_OUT", 11,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnLogStderr() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalLogStderr < c2.(*displaydata.DisplayAppStats).TotalLogStderr
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%11v", util.Format(appStats.TotalLogStderr))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalLogStderr)
	}
	c := uiCommon.NewListColumn("TotalLogStderr", "LOG_ERR", 11,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnReq1() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalTraffic.EventL1Rate < c2.(*displaydata.DisplayAppStats).TotalTraffic.EventL1Rate
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%6v", util.Format(int64(appStats.TotalTraffic.EventL1Rate)))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.EventL1Rate)
	}
	c := uiCommon.NewListColumn("REQ1", "REQ/1", 6,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnReq10() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalTraffic.EventL10Rate < c2.(*displaydata.DisplayAppStats).TotalTraffic.EventL10Rate
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%7v", util.Format(int64(appStats.TotalTraffic.EventL10Rate)))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.EventL10Rate)
	}
	c := uiCommon.NewListColumn("REQ10", "REQ/10", 7,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnReq60() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalTraffic.EventL60Rate < c2.(*displaydata.DisplayAppStats).TotalTraffic.EventL60Rate
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%7v", util.Format(int64(appStats.TotalTraffic.EventL60Rate)))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.EventL60Rate)
	}
	c := uiCommon.NewListColumn("REQ60", "REQ/60", 7,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnTotalReq() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalTraffic.HttpAllCount < c2.(*displaydata.DisplayAppStats).TotalTraffic.HttpAllCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.HttpAllCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.HttpAllCount)
	}
	c := uiCommon.NewListColumn("TOTREQ", "TOT-REQ", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}
func (asUI *AppListView) column2XX() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalTraffic.Http2xxCount < c2.(*displaydata.DisplayAppStats).TotalTraffic.Http2xxCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.Http2xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http2xxCount)
	}
	c := uiCommon.NewListColumn("2XX", "2XX", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}
func (asUI *AppListView) column3XX() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalTraffic.Http3xxCount < c2.(*displaydata.DisplayAppStats).TotalTraffic.Http3xxCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.Http3xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http3xxCount)
	}
	c := uiCommon.NewListColumn("3XX", "3XX", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) column4XX() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalTraffic.Http4xxCount < c2.(*displaydata.DisplayAppStats).TotalTraffic.Http4xxCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.Http4xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http4xxCount)
	}
	c := uiCommon.NewListColumn("4XX", "4XX", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) column5XX() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayAppStats).TotalTraffic.Http5xxCount < c2.(*displaydata.DisplayAppStats).TotalTraffic.Http5xxCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.Http5xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayAppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http5xxCount)
	}
	c := uiCommon.NewListColumn("5XX", "5XX", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}
