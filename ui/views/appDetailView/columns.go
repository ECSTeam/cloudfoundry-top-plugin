package appDetailView

import (
	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

func (asUI *AppDetailView) columnAppName() *uiCommon.ListColumn {
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

func (asUI *AppDetailView) columnSpaceName() *uiCommon.ListColumn {
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

func (asUI *AppDetailView) columnOrgName() *uiCommon.ListColumn {
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
