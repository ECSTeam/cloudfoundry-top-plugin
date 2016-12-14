package cellDetailView

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/eventdata"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/dataView"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/displaydata"
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
		//uiCommon.NewSortColumn("CPU", true),
		//uiCommon.NewSortColumn("REQ60", true),
		uiCommon.NewSortColumn("appName", false),
		uiCommon.NewSortColumn("spaceName", false),
		uiCommon.NewSortColumn("orgName", false)}

	dataListView := dataView.NewDataListView(masterUI, parentView,
		name, topMargin, bottomMargin,
		eventProcessor, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.UpdateHeaderCallback = asUI.updateHeader
	dataListView.GetListData = asUI.GetListData

	dataListView.SetTitle("Cell Detail - Container List")
	dataListView.HelpText = helpText

	asUI.DataListView = dataListView

	return asUI

}

func (asUI *CellDetailView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, asUI.columnAppName())
	columns = append(columns, asUI.columnSpaceName())
	columns = append(columns, asUI.columnOrgName())
	return columns
}

func (asUI *CellDetailView) initializeCallback(g *gocui.Gui, viewName string) error {
	// TODO: This needs to be handled in dataListView someplace for child (detailed) views as all of them will need a back action
	if err := g.SetKeybinding(viewName, 'q', gocui.ModNone, asUI.closeAppDetailView); err != nil {
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
	appStatsArray := eventdata.PopulateNamesIfNeeded(appMap)
	for _, appStats := range appStatsArray {
		for _, containerStats := range appStats.ContainerArray {
			if containerStats != nil {
				if containerStats.Ip == asUI.cellIp {
					// This is a container on the selected cell
					containerStats := displaydata.NewDisplayContainerStats(containerStats, appStats)
					containerStatsArray = append(containerStatsArray, containerStats)
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
