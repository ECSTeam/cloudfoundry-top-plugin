package cellView

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/eventdata"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/dataView"
)

type CellListView struct {
	dataListView *dataView.DataListView
}

func NewCellListView(masterUI masterUIInterface.MasterUIInterface,
	name string, topMargin, bottomMargin int,
	eventProcessor *eventdata.EventProcessor) *CellListView {

	asUI := &CellListView{}

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("IP", false),
		//uiCommon.NewSortColumn("colB", true),
	}

	dataListView := dataView.NewDataListView(masterUI,
		name, topMargin, bottomMargin,
		eventProcessor, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.UpdateHeaderCallback = asUI.updateHeader
	dataListView.GetListData = asUI.GetListData

	asUI.dataListView = dataListView

	return asUI

}

func (asUI *CellListView) Layout(g *gocui.Gui) error {
	return asUI.dataListView.Layout(g)
}
func (asUI *CellListView) Name() string {
	return asUI.dataListView.Name()
}

func (asUI *CellListView) UpdateDisplay(g *gocui.Gui) error {
	return asUI.dataListView.UpdateDisplay(g)
}

func (asUI *CellListView) GetCurrentEventData() *eventdata.EventData {
	return asUI.dataListView.GetCurrentEventData()
}

func (asUI *CellListView) GetDisplayedEventData() *eventdata.EventData {
	return asUI.dataListView.GetDisplayedEventData()
}

func (asUI *CellListView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, asUI.columnIp())
	columns = append(columns, asUI.columnNumOfCpus())
	columns = append(columns, asUI.columnCapacityTotalMemory())
	columns = append(columns, asUI.columnCapacityRemainingMemory())
	columns = append(columns, asUI.columnCapacityTotalDisk())
	columns = append(columns, asUI.columnCapacityRemainingDisk())
	columns = append(columns, asUI.columnCapacityTotalContainers())
	columns = append(columns, asUI.columnContainerCount())

	columns = append(columns, asUI.columnDeploymentName())
	columns = append(columns, asUI.columnJobName())
	columns = append(columns, asUI.columnJobIndex())

	return columns
}

func (asUI *CellListView) GetListData() []uiCommon.IData {
	cellList := asUI.postProcessData()
	listData := asUI.convertToListData(cellList)
	return listData
}

func (asUI *CellListView) postProcessData() []*eventdata.CellStats {
	cellMap := asUI.GetDisplayedEventData().CellMap
	cellList := make([]*eventdata.CellStats, 0, len(cellMap))
	for _, cellStats := range cellMap {
		cellList = append(cellList, cellStats)
	}
	return cellList
}

func (asUI *CellListView) convertToListData(statsList []*eventdata.CellStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, len(statsList))
	for i, d := range statsList {
		listData[i] = d
	}
	return listData
}

func (asUI *CellListView) PreRowDisplay(data uiCommon.IData, isSelected bool) string {
	return ""
}

func (asUI *CellListView) updateHeader(g *gocui.Gui, v *gocui.View) error {
	fmt.Fprintf(v, "\nTODO: Show summary Cell stats")
	return nil
}
