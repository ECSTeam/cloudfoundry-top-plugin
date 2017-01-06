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

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func ColumnTotalCpuPercentage() *uiCommon.ListColumn {
	defaultColSize := 6
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayContainerStats).ContainerMetric.GetCpuPercentage() < c2.(*DisplayContainerStats).ContainerMetric.GetCpuPercentage())
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayContainerStats)
		totalCpuInfo := ""
		cpuPercentage := stats.ContainerMetric.GetCpuPercentage()
		if cpuPercentage == 0 {
			totalCpuInfo = fmt.Sprintf("%6v", "--")
		} else {
			if cpuPercentage >= 100.0 {
				totalCpuInfo = fmt.Sprintf("%6.0f", cpuPercentage)
			} else if cpuPercentage >= 10.0 {
				totalCpuInfo = fmt.Sprintf("%6.1f", cpuPercentage)
			} else {
				totalCpuInfo = fmt.Sprintf("%6.2f", cpuPercentage)
			}
		}
		return fmt.Sprintf("%6v", totalCpuInfo)

	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", stats.ContainerMetric.GetCpuPercentage())
	}
	c := uiCommon.NewListColumn("CPU_PERCENT", "CPU%", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)

	return c
}

func ColumnContainerIndex() *uiCommon.ListColumn {
	defaultColSize := 4
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayContainerStats).ContainerIndex < c2.(*DisplayContainerStats).ContainerIndex)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayContainerStats)
		display := fmt.Sprintf("%4v", stats.ContainerIndex)
		return fmt.Sprintf("%4v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", stats.ContainerIndex)
	}
	c := uiCommon.NewListColumn("IDX", "IDX", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnAppName() *uiCommon.ListColumn {
	defaultColSize := 50
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*DisplayContainerStats).AppName, c2.(*DisplayContainerStats).AppName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*DisplayContainerStats)
		return util.FormatDisplayData(appStats.AppName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return appStats.AppName
	}
	c := uiCommon.NewListColumn("appName", "APPLICATION", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func ColumnSpaceName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*DisplayContainerStats).SpaceName, c2.(*DisplayContainerStats).SpaceName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*DisplayContainerStats)
		return util.FormatDisplayData(appStats.SpaceName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return appStats.SpaceName
	}
	c := uiCommon.NewListColumn("spaceName", "SPACE", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func ColumnOrgName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*DisplayContainerStats).OrgName, c2.(*DisplayContainerStats).OrgName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*DisplayContainerStats)
		return util.FormatDisplayData(appStats.OrgName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return appStats.OrgName
	}
	c := uiCommon.NewListColumn("orgName", "ORG", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func ColumnMemoryUsed() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).ContainerMetric.GetMemoryBytes() < c2.(*DisplayContainerStats).ContainerMetric.GetMemoryBytes()
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayContainerStats)
		memInfo := fmt.Sprintf("%9v", util.ByteSize(stats.ContainerMetric.GetMemoryBytes()).StringWithPrecision(1))
		return fmt.Sprintf("%9v", memInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.ContainerMetric.GetMemoryBytes())
	}
	c := uiCommon.NewListColumn("MEM_USED", "MEM_USED", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnMemoryFree() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).FreeMemory < c2.(*DisplayContainerStats).FreeMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayContainerStats)
		memInfo := fmt.Sprintf("%9v", util.ByteSize(stats.FreeMemory).StringWithPrecision(1))
		return fmt.Sprintf("%9v", memInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.FreeMemory)
	}
	c := uiCommon.NewListColumn("MEM_FREE", "MEM_FREE", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnMemoryReserved() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).ReservedMemory < c2.(*DisplayContainerStats).ReservedMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayContainerStats)
		memInfo := fmt.Sprintf("%9v", util.ByteSize(stats.ReservedMemory).StringWithPrecision(1))
		return fmt.Sprintf("%9v", memInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.ReservedMemory)
	}
	c := uiCommon.NewListColumn("MEM_RSVD", "MEM_RSVD", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnDiskUsed() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).ContainerMetric.GetDiskBytes() < c2.(*DisplayContainerStats).ContainerMetric.GetDiskBytes()
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*DisplayContainerStats)
		diskUsed := fmt.Sprintf("%9v", util.ByteSize(appStats.ContainerMetric.GetDiskBytes()).StringWithPrecision(1))
		return diskUsed
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.ContainerMetric.GetDiskBytes())
	}
	c := uiCommon.NewListColumn("DISK_USED", "DISK_USED", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnDiskFree() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).FreeDisk < c2.(*DisplayContainerStats).FreeDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayContainerStats)
		memInfo := fmt.Sprintf("%9v", util.ByteSize(stats.FreeDisk).StringWithPrecision(1))
		return fmt.Sprintf("%9v", memInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.FreeDisk)
	}
	c := uiCommon.NewListColumn("DISK_FREE", "DISK_FREE", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnDiskReserved() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).ReservedDisk < c2.(*DisplayContainerStats).ReservedDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayContainerStats)
		memInfo := fmt.Sprintf("%9v", util.ByteSize(stats.ReservedDisk).StringWithPrecision(1))
		return fmt.Sprintf("%9v", memInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.ReservedDisk)
	}
	c := uiCommon.NewListColumn("DISK_RSVD", "DISK_RSVD", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnLogStdout() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).OutCount < c2.(*DisplayContainerStats).OutCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayContainerStats)
		display := fmt.Sprintf("%9v", util.Format(stats.OutCount))
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.OutCount)
	}
	c := uiCommon.NewListColumn("LOG_OUT", "LOG_OUT", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnLogStderr() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).ErrCount < c2.(*DisplayContainerStats).ErrCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayContainerStats)
		display := fmt.Sprintf("%9v", util.Format(stats.ErrCount))
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.ErrCount)
	}
	c := uiCommon.NewListColumn("LOG_ERR", "LOG_ERR", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnCellIp() *uiCommon.ListColumn {
	defaultColSize := 16
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.Ip2long(c1.(*DisplayContainerStats).Ip) < util.Ip2long(c2.(*DisplayContainerStats).Ip)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*DisplayContainerStats)
		if appStats.Ip != "" {
			return util.FormatDisplayData(appStats.Ip, defaultColSize)
		} else {
			return "--"
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return appStats.Ip
	}
	c := uiCommon.NewListColumn("CELL_IP", "CELL_IP", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}
