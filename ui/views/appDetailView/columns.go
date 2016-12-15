package appDetailView

import (
	"fmt"

	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

func ColumnAppName() *uiCommon.ListColumn {
	defaultColSize := 50
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*displaydata.DisplayContainerStats).AppName, c2.(*displaydata.DisplayContainerStats).AppName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		return util.FormatDisplayData(appStats.AppName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		return appStats.AppName
	}
	c := uiCommon.NewListColumn("appName", "APPLICATION", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func ColumnSpaceName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*displaydata.DisplayContainerStats).SpaceName, c2.(*displaydata.DisplayContainerStats).SpaceName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		return util.FormatDisplayData(appStats.SpaceName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		return appStats.SpaceName
	}
	c := uiCommon.NewListColumn("spaceName", "SPACE", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func ColumnOrgName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*displaydata.DisplayContainerStats).OrgName, c2.(*displaydata.DisplayContainerStats).OrgName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		return util.FormatDisplayData(appStats.OrgName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		return appStats.OrgName
	}
	c := uiCommon.NewListColumn("orgName", "ORG", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func ColumnContainerIndex() *uiCommon.ListColumn {
	defaultColSize := 4
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayContainerStats).ContainerIndex < c2.(*displaydata.DisplayContainerStats).ContainerIndex)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayContainerStats)
		display := fmt.Sprintf("%4v", stats.ContainerIndex)
		return fmt.Sprintf("%4v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*displaydata.DisplayContainerStats)
		return fmt.Sprintf("%v", stats.ContainerIndex)
	}
	c := uiCommon.NewListColumn("IDX", "IDX", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnTotalCpuPercentage() *uiCommon.ListColumn {
	defaultColSize := 6
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayContainerStats).ContainerMetric.GetCpuPercentage() < c2.(*displaydata.DisplayContainerStats).ContainerMetric.GetCpuPercentage())
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayContainerStats)
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
		stats := data.(*displaydata.DisplayContainerStats)
		return fmt.Sprintf("%v", stats.ContainerMetric.GetCpuPercentage())
	}
	c := uiCommon.NewListColumn("CPU_PERCENT", "CPU%", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)

	return c
}

func ColumnMemoryUsed() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayContainerStats).ContainerMetric.GetMemoryBytes() < c2.(*displaydata.DisplayContainerStats).ContainerMetric.GetMemoryBytes()
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayContainerStats)
		memInfo := fmt.Sprintf("%9v", util.ByteSize(stats.ContainerMetric.GetMemoryBytes()).StringWithPrecision(1))
		return fmt.Sprintf("%9v", memInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.ContainerMetric.GetMemoryBytes())
	}
	c := uiCommon.NewListColumn("MEM_USED", "MEM_USED", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnMemoryFree() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayContainerStats).FreeMemory < c2.(*displaydata.DisplayContainerStats).FreeMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayContainerStats)
		memInfo := fmt.Sprintf("%9v", util.ByteSize(stats.FreeMemory).StringWithPrecision(1))
		return fmt.Sprintf("%9v", memInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.FreeMemory)
	}
	c := uiCommon.NewListColumn("MEM_FREE", "MEM_FREE", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnMemoryReserved() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayContainerStats).ReservedMemory < c2.(*displaydata.DisplayContainerStats).ReservedMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayContainerStats)
		memInfo := fmt.Sprintf("%9v", util.ByteSize(stats.ReservedMemory).StringWithPrecision(1))
		return fmt.Sprintf("%9v", memInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.ReservedMemory)
	}
	c := uiCommon.NewListColumn("MEM_RSVD", "MEM_RSVD", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnDiskUsed() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayContainerStats).ContainerMetric.GetDiskBytes() < c2.(*displaydata.DisplayContainerStats).ContainerMetric.GetDiskBytes()
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		diskUsed := fmt.Sprintf("%9v", util.ByteSize(appStats.ContainerMetric.GetDiskBytes()).StringWithPrecision(1))
		return diskUsed
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.ContainerMetric.GetDiskBytes())
	}
	c := uiCommon.NewListColumn("DISK_USED", "DISK_USED", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnDiskFree() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayContainerStats).FreeDisk < c2.(*displaydata.DisplayContainerStats).FreeDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayContainerStats)
		memInfo := fmt.Sprintf("%9v", util.ByteSize(stats.FreeDisk).StringWithPrecision(1))
		return fmt.Sprintf("%9v", memInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.FreeDisk)
	}
	c := uiCommon.NewListColumn("DISK_FREE", "DISK_FREE", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnDiskReserved() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayContainerStats).ReservedDisk < c2.(*displaydata.DisplayContainerStats).ReservedDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayContainerStats)
		memInfo := fmt.Sprintf("%9v", util.ByteSize(stats.ReservedDisk).StringWithPrecision(1))
		return fmt.Sprintf("%9v", memInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.ReservedDisk)
	}
	c := uiCommon.NewListColumn("DISK_RSVD", "DISK_RSVD", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnLogStdout() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayContainerStats).OutCount < c2.(*displaydata.DisplayContainerStats).OutCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayContainerStats)
		display := fmt.Sprintf("%9v", util.Format(stats.OutCount))
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.OutCount)
	}
	c := uiCommon.NewListColumn("STDOUT", "STDOUT", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func ColumnLogStderr() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayContainerStats).ErrCount < c2.(*displaydata.DisplayContainerStats).ErrCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayContainerStats)
		display := fmt.Sprintf("%9v", util.Format(stats.ErrCount))
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayContainerStats)
		return fmt.Sprintf("%v", appStats.ErrCount)
	}
	c := uiCommon.NewListColumn("STDERR", "STDERR", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}
