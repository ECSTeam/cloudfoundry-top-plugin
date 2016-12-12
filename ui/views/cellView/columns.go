package cellView

import (
	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

func (asUI *CellListView) columnA() *uiCommon.ListColumn {
	defaultColSize := 20
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		//return util.CaseInsensitiveLess(c1.(*eventdata.AppStats).AppName, c2.(*eventdata.AppStats).AppName)
		return true
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		//appStats := data.(*eventdata.AppStats)
		//return formatDisplayData(appStats.AppName, defaultColSize)
		return "Data A"
	}
	rawValueFunc := func(data uiCommon.IData) string {
		//appStats := data.(*eventdata.AppStats)
		//return appStats.AppName
		return "Data A"
	}
	c := uiCommon.NewListColumn("colA", "COL_A", defaultColSize,
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnB() *uiCommon.ListColumn {
	defaultColSize := 20
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		//return util.CaseInsensitiveLess(c1.(*eventdata.AppStats).AppName, c2.(*eventdata.AppStats).AppName)
		return true
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		//appStats := data.(*eventdata.AppStats)
		//return formatDisplayData(appStats.AppName, defaultColSize)
		return "Data B"
	}
	rawValueFunc := func(data uiCommon.IData) string {
		//appStats := data.(*eventdata.AppStats)
		//return appStats.AppName
		return "Data B"
	}
	c := uiCommon.NewListColumn("colB", "COL_B", defaultColSize,
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}
