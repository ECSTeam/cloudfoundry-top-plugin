package appView

import (
	"fmt"
	"strconv"

	"github.com/kkellner/cloudfoundry-top-plugin/eventdata"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

func (asUI *AppListView) columnAppName() *uiCommon.ListColumn {
	defaultColSize := 50
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*eventdata.AppStats).AppName, c2.(*eventdata.AppStats).AppName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		return formatDisplayData(appStats.AppName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return appStats.AppName
	}
	c := uiCommon.NewListColumn("appName", "APPLICATION", defaultColSize,
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnSpaceName() *uiCommon.ListColumn {
	defaultColSize := 10
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*eventdata.AppStats).SpaceName, c2.(*eventdata.AppStats).SpaceName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		return formatDisplayData(appStats.SpaceName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return appStats.SpaceName
	}
	c := uiCommon.NewListColumn("spaceName", "SPACE", defaultColSize,
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnOrgName() *uiCommon.ListColumn {
	defaultColSize := 10
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*eventdata.AppStats).OrgName, c2.(*eventdata.AppStats).OrgName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		return formatDisplayData(appStats.OrgName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return appStats.OrgName
	}
	c := uiCommon.NewListColumn("orgName", "ORG", defaultColSize,
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnReportingContainers() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.AppStats).TotalReportingContainers < c2.(*eventdata.AppStats).TotalReportingContainers
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%3v", appStats.TotalReportingContainers)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return strconv.Itoa(appStats.TotalReportingContainers)
	}
	c := uiCommon.NewListColumn("reportingContainers", "RCR", 3,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnTotalCpu() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.AppStats).TotalCpuPercentage < c2.(*eventdata.AppStats).TotalCpuPercentage
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
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
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%.2f", appStats.TotalCpuPercentage)
	}
	c := uiCommon.NewListColumn("CPU", "CPU%", 6,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnTotalMemoryUsed() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.AppStats).TotalUsedMemory < c2.(*eventdata.AppStats).TotalUsedMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		totalMemInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalMemInfo = fmt.Sprintf("%9v", "--")
		} else {
			totalMemInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.TotalUsedMemory).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", totalMemInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%v", appStats.TotalUsedMemory)
	}
	c := uiCommon.NewListColumn("MEM", "MEM", 9,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnTotalDiskUsed() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.AppStats).TotalUsedDisk < c2.(*eventdata.AppStats).TotalUsedDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		totalDiskInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalDiskInfo = fmt.Sprintf("%9v", "--")
		} else {
			totalDiskInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.TotalUsedDisk).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", totalDiskInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%v", appStats.TotalUsedDisk)
	}
	c := uiCommon.NewListColumn("DISK", "DISK", 9,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnAvgResponseTimeL60Info() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.AppStats).TotalTraffic.AvgResponseL60Time < c2.(*eventdata.AppStats).TotalTraffic.AvgResponseL60Time
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		avgResponseTimeL60Info := "--"
		if appStats.TotalTraffic.AvgResponseL60Time >= 0 {
			avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL60Time / 1000000
			avgResponseTimeL60Info = fmt.Sprintf("%6.0f", avgResponseTimeMs)
		}
		return fmt.Sprintf("%6v", avgResponseTimeL60Info)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.AvgResponseL60Time)
	}
	c := uiCommon.NewListColumn("avgResponseTimeL60", "RESP", 6,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnLogCount() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.AppStats).TotalLogCount < c2.(*eventdata.AppStats).TotalLogCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%11v", util.Format(appStats.TotalLogCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%v", appStats.TotalLogCount)
	}
	c := uiCommon.NewListColumn("totalLogCount", "LOGS", 11,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnReq1() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.AppStats).TotalTraffic.EventL1Rate < c2.(*eventdata.AppStats).TotalTraffic.EventL1Rate
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%6v", util.Format(int64(appStats.TotalTraffic.EventL1Rate)))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.EventL1Rate)
	}
	c := uiCommon.NewListColumn("REQ1", "REQ/1", 6,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnReq10() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.AppStats).TotalTraffic.EventL10Rate < c2.(*eventdata.AppStats).TotalTraffic.EventL10Rate
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%7v", util.Format(int64(appStats.TotalTraffic.EventL10Rate)))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.EventL10Rate)
	}
	c := uiCommon.NewListColumn("REQ10", "REQ/10", 7,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnReq60() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.AppStats).TotalTraffic.EventL60Rate < c2.(*eventdata.AppStats).TotalTraffic.EventL60Rate
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%7v", util.Format(int64(appStats.TotalTraffic.EventL60Rate)))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.EventL60Rate)
	}
	c := uiCommon.NewListColumn("REQ60", "REQ/60", 7,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnTotalReq() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.AppStats).TotalTraffic.HttpAllCount < c2.(*eventdata.AppStats).TotalTraffic.HttpAllCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.HttpAllCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.HttpAllCount)
	}
	c := uiCommon.NewListColumn("TOTREQ", "TOT-REQ", 10,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}
func (asUI *AppListView) column2XX() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.AppStats).TotalTraffic.Http2xxCount < c2.(*eventdata.AppStats).TotalTraffic.Http2xxCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.Http2xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http2xxCount)
	}
	c := uiCommon.NewListColumn("2XX", "2XX", 10,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}
func (asUI *AppListView) column3XX() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.AppStats).TotalTraffic.Http3xxCount < c2.(*eventdata.AppStats).TotalTraffic.Http3xxCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.Http3xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http3xxCount)
	}
	c := uiCommon.NewListColumn("3XX", "3XX", 10,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) column4XX() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.AppStats).TotalTraffic.Http4xxCount < c2.(*eventdata.AppStats).TotalTraffic.Http4xxCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.Http4xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http4xxCount)
	}
	c := uiCommon.NewListColumn("4XX", "4XX", 10,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) column5XX() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.AppStats).TotalTraffic.Http5xxCount < c2.(*eventdata.AppStats).TotalTraffic.Http5xxCount
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%10v", util.Format(appStats.TotalTraffic.Http5xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*eventdata.AppStats)
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http5xxCount)
	}
	c := uiCommon.NewListColumn("5XX", "5XX", 10,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}
