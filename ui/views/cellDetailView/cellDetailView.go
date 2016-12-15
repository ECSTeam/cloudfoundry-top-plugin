package cellDetailView

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/eventdata"
	"github.com/kkellner/cloudfoundry-top-plugin/metadata"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/appDetailView"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/dataView"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

type CellDetailView struct {
	*dataView.DataListView
	cellIp string
}

func NewCellDetailView(masterUI masterUIInterface.MasterUIInterface,
	parentView dataView.DataListViewInterface,
	name string, topMargin, bottomMargin int,
	eventProcessor *eventdata.EventProcessor, cellIp string) *CellDetailView {

	asUI := &CellDetailView{
		cellIp: cellIp,
	}

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("CPU_PERCENT", true),
		uiCommon.NewSortColumn("appName", false),
		uiCommon.NewSortColumn("spaceName", false),
		uiCommon.NewSortColumn("orgName", false),
		uiCommon.NewSortColumn("IDX", false),
	}

	dataListView := dataView.NewDataListView(masterUI, parentView,
		name, topMargin, bottomMargin,
		eventProcessor, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.UpdateHeaderCallback = asUI.updateHeader
	dataListView.GetListData = asUI.GetListData

	dataListView.SetTitle(fmt.Sprintf("Cell %v Detail - Container List", cellIp))
	dataListView.HelpText = helpText

	asUI.DataListView = dataListView

	return asUI

}

func (asUI *CellDetailView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, appDetailView.ColumnAppName())
	columns = append(columns, appDetailView.ColumnContainerIndex())
	columns = append(columns, appDetailView.ColumnSpaceName())
	columns = append(columns, appDetailView.ColumnOrgName())
	columns = append(columns, appDetailView.ColumnTotalCpuPercentage())
	columns = append(columns, appDetailView.ColumnMemoryUsed())
	columns = append(columns, appDetailView.ColumnMemoryFree())
	columns = append(columns, appDetailView.ColumnDiskUsed())
	columns = append(columns, appDetailView.ColumnDiskFree())
	columns = append(columns, appDetailView.ColumnLogStdout())
	columns = append(columns, appDetailView.ColumnLogStderr())

	return columns
}

func (asUI *CellDetailView) initializeCallback(g *gocui.Gui, viewName string) error {
	// TODO: This needs to be handled in dataListView someplace for child (detailed) views as all of them will need a back action
	if err := g.SetKeybinding(viewName, 'x', gocui.ModNone, asUI.closeAppDetailView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, gocui.KeyEsc, gocui.ModNone, asUI.closeAppDetailView); err != nil {
		log.Panicln(err)
	}
	return nil
}

// TODO: Need to put this in common dataListView - but allow for callback to do special close processing (as needed by appDetailView to close other views)
func (asUI *CellDetailView) closeAppDetailView(g *gocui.Gui, v *gocui.View) error {
	if err := asUI.GetMasterUI().CloseView(asUI); err != nil {
		return err
	}
	return nil
}

func (asUI *CellDetailView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *CellDetailView) postProcessData() []*displaydata.DisplayContainerStats {

	containerStatsArray := make([]*displaydata.DisplayContainerStats, 0)

	appMap := asUI.GetDisplayedEventData().AppMap
	appStatsArray := eventdata.PopulateNamesFromMap(appMap)
	for _, appStats := range appStatsArray {
		appMetadata := metadata.FindAppMetadata(appStats.AppId)
		for _, containerStats := range appStats.ContainerArray {
			if containerStats != nil {
				if containerStats.Ip == asUI.cellIp {
					// This is a container on the selected cell
					displayContainerStats := displaydata.NewDisplayContainerStats(containerStats, appStats)
					usedMemory := containerStats.ContainerMetric.GetMemoryBytes()
					freeMemory := (uint64(appMetadata.MemoryMB) * util.MEGABYTE) - usedMemory
					displayContainerStats.FreeMemory = freeMemory
					containerStatsArray = append(containerStatsArray, displayContainerStats)
				}
			}
		}
	}

	return containerStatsArray
}

func (asUI *CellDetailView) convertToListData(containerStatsArray []*displaydata.DisplayContainerStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(containerStatsArray))
	for _, d := range containerStatsArray {
		listData = append(listData, d)
	}
	return listData
}

func (asUI *CellDetailView) PreRowDisplay(data uiCommon.IData, isSelected bool) string {
	return ""
}

func (asUI *CellDetailView) updateHeader(g *gocui.Gui, v *gocui.View) error {
	fmt.Fprintf(v, "\nTODO: Show summary Cell stats")
	return nil
}
