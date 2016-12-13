package cellView

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/eventdata"
	"github.com/kkellner/cloudfoundry-top-plugin/metadata"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/dataView"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/displaydata"
)

const MEGABYTE = (1024 * 1024)

type CellListView struct {
	*dataView.DataListView
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

	dataListView.Title = "Cell List"
	dataListView.HelpText = helpText

	asUI.DataListView = dataListView

	return asUI

}

func (asUI *CellListView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, asUI.columnIp())

	columns = append(columns, asUI.columnTotalCpuPercentage())
	columns = append(columns, asUI.columnTotalReportingContainers())

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

	columns = append(columns, asUI.columnTotalContainerReservedMemory())
	columns = append(columns, asUI.columnTotalContainerUsedMemory())

	columns = append(columns, asUI.columnTotalContainerReservedDisk())
	columns = append(columns, asUI.columnTotalContainerUsedDisk())

	return columns
}

func (asUI *CellListView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *CellListView) postProcessData() map[string]*displaydata.DisplayCellStats {
	cellMap := asUI.GetDisplayedEventData().CellMap

	displayCellMap := make(map[string]*displaydata.DisplayCellStats)
	for ip, cellStats := range cellMap {
		displayStat := displaydata.NewDisplayCellStats(cellStats)
		displayCellMap[ip] = displayStat
	}

	appMap := asUI.GetDisplayedEventData().AppMap
	for _, appStats := range appMap {
		for _, containerStats := range appStats.ContainerArray {
			if containerStats != nil {
				cellStats := displayCellMap[containerStats.Ip]

				if cellStats != nil {
					logOutCount := containerStats.OutCount
					cellStats.TotalLogOutCount = cellStats.TotalLogOutCount + logOutCount

					logErrCount := containerStats.ErrCount
					cellStats.TotalLogErrCount = cellStats.TotalLogErrCount + logErrCount

					if containerStats.ContainerMetric != nil {

						appMetadata := metadata.FindAppMetadata(appStats.AppId)

						cellStats.TotalReportingContainers = cellStats.TotalReportingContainers + 1

						cpuValue := containerStats.ContainerMetric.GetCpuPercentage()
						cellStats.TotalContainerCpuPercentage = cellStats.TotalContainerCpuPercentage + cpuValue

						//reservedMemoryValue := containerStats.ContainerMetric.GetMemoryBytesQuota()
						cellStats.TotalContainerReservedMemory = cellStats.TotalContainerReservedMemory + uint64(appMetadata.MemoryMB*MEGABYTE)

						usedMemoryValue := containerStats.ContainerMetric.GetMemoryBytes()
						cellStats.TotalContainerUsedMemory = cellStats.TotalContainerUsedMemory + usedMemoryValue

						//reservedDiskValue := containerStats.ContainerMetric.GetDiskBytesQuota()
						cellStats.TotalContainerReservedDisk = cellStats.TotalContainerReservedDisk + uint64(appMetadata.DiskQuotaMB*MEGABYTE)

						usedDiskValue := containerStats.ContainerMetric.GetDiskBytes()
						cellStats.TotalContainerUsedDisk = cellStats.TotalContainerUsedDisk + usedDiskValue
					}
				}
			}
		}
	}

	return displayCellMap
}

func (asUI *CellListView) convertToListData(displayCellMap map[string]*displaydata.DisplayCellStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(displayCellMap))
	for _, d := range displayCellMap {
		listData = append(listData, d)
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
