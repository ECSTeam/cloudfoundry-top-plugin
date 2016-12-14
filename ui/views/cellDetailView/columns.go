package cellDetailView

import (
	"fmt"

	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

func (asUI *CellDetailView) columnAppName() *uiCommon.ListColumn {
	defaultColSize := 50
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
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
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellDetailView) columnSpaceName() *uiCommon.ListColumn {
	defaultColSize := 10
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
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
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellDetailView) columnOrgName() *uiCommon.ListColumn {
	defaultColSize := 10
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
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
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellDetailView) columnTotalCpuPercentage() *uiCommon.ListColumn {
	defaultColSize := 6
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).TotalContainerCpuPercentage < c2.(*displaydata.DisplayCellStats).TotalContainerCpuPercentage)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		totalCpuInfo := ""
		if cellStats.TotalReportingContainers == 0 {
			totalCpuInfo = fmt.Sprintf("%6v", "--")
		} else {
			if cellStats.TotalContainerCpuPercentage >= 100.0 {
				totalCpuInfo = fmt.Sprintf("%6.0f", cellStats.TotalContainerCpuPercentage)
			} else if cellStats.TotalContainerCpuPercentage >= 10.0 {
				totalCpuInfo = fmt.Sprintf("%6.1f", cellStats.TotalContainerCpuPercentage)
			} else {
				totalCpuInfo = fmt.Sprintf("%6.2f", cellStats.TotalContainerCpuPercentage)
			}
		}
		return fmt.Sprintf("%6v", totalCpuInfo)

	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.TotalContainerCpuPercentage)
	}
	c := uiCommon.NewListColumn("CPU_PERCENT", "CPU%", defaultColSize,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)

	return c
}

func (asUI *CellDetailView) columnColA() *uiCommon.ListColumn {
	defaultColSize := 45
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return true
		//return util.CaseInsensitiveLess(c1.(*displaydata.DisplayCellStats).JobName, c2.(*displaydata.DisplayCellStats).JobName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		//cellStats := data.(*displaydata.DisplayCellStats)
		//return util.FormatDisplayData(cellStats.JobName, defaultColSize)
		return "Hello"
	}
	rawValueFunc := func(data uiCommon.IData) string {
		///cellStats := data.(*displaydata.DisplayCellStats)
		//return cellStats.JobName
		return "Hello"
	}
	c := uiCommon.NewListColumn("COL_A", "COL_A", defaultColSize,
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}
