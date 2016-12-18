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

package cellView

import (
	"fmt"
	"log"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/cellDetailView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

type CellListView struct {
	*dataView.DataListView
}

func NewCellListView(masterUI masterUIInterface.MasterUIInterface,
	name string, topMargin, bottomMargin int,
	eventProcessor *eventdata.EventProcessor) *CellListView {

	asUI := &CellListView{}

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("CPU_PERCENT", true),
		uiCommon.NewSortColumn("CELL_IP", false),
	}

	dataListView := dataView.NewDataListView(masterUI, nil,
		name, topMargin, bottomMargin,
		eventProcessor, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.UpdateHeaderCallback = asUI.updateHeader
	dataListView.GetListData = asUI.GetListData

	dataListView.SetTitle("Cell List")
	dataListView.HelpText = helpText
	dataListView.HelpTextTips = appView.HelpTextTips

	asUI.DataListView = dataListView

	return asUI

}

func (asUI *CellListView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, asUI.columnCellIp())

	columns = append(columns, asUI.columnTotalCpuPercentage())
	columns = append(columns, asUI.columnTotalReportingContainers())

	columns = append(columns, asUI.columnNumOfCpus())

	columns = append(columns, asUI.columnCapacityTotalMemory())
	columns = append(columns, asUI.columnCapacityRemainingMemory())
	columns = append(columns, asUI.columnTotalContainerReservedMemory())
	columns = append(columns, asUI.columnTotalContainerUsedMemory())

	columns = append(columns, asUI.columnCapacityTotalDisk())
	columns = append(columns, asUI.columnCapacityRemainingDisk())
	columns = append(columns, asUI.columnTotalContainerReservedDisk())
	columns = append(columns, asUI.columnTotalContainerUsedDisk())

	columns = append(columns, asUI.columnCapacityTotalContainers())
	columns = append(columns, asUI.columnContainerCount())

	columns = append(columns, asUI.columnDeploymentName())
	columns = append(columns, asUI.columnJobName())
	columns = append(columns, asUI.columnJobIndex())

	return columns
}

func (asUI *CellListView) initializeCallback(g *gocui.Gui, viewName string) error {

	if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone, asUI.enterAction); err != nil {
		log.Panicln(err)
	}

	return nil
}

func (asUI *CellListView) enterAction(g *gocui.Gui, v *gocui.View) error {
	highlightKey := asUI.GetListWidget().HighlightKey()
	if asUI.GetListWidget().HighlightKey() != "" {
		topMargin, bottomMargin := asUI.GetMargins()

		detailView := cellDetailView.NewCellDetailView(asUI.GetMasterUI(), asUI,
			"cellDetailView",
			topMargin, bottomMargin,
			asUI.GetEventProcessor(),
			highlightKey)

		asUI.SetDetailView(detailView)
		asUI.GetMasterUI().OpenView(g, detailView)
	}

	return nil
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

						appMetadata := asUI.GetAppMdMgr().FindAppMetadata(appStats.AppId)

						cellStats.TotalReportingContainers = cellStats.TotalReportingContainers + 1

						cpuValue := containerStats.ContainerMetric.GetCpuPercentage()
						cellStats.TotalContainerCpuPercentage = cellStats.TotalContainerCpuPercentage + cpuValue

						//reservedMemoryValue := containerStats.ContainerMetric.GetMemoryBytesQuota()
						cellStats.TotalContainerReservedMemory = cellStats.TotalContainerReservedMemory + uint64(appMetadata.MemoryMB*util.MEGABYTE)

						usedMemoryValue := containerStats.ContainerMetric.GetMemoryBytes()
						cellStats.TotalContainerUsedMemory = cellStats.TotalContainerUsedMemory + usedMemoryValue

						//reservedDiskValue := containerStats.ContainerMetric.GetDiskBytesQuota()
						cellStats.TotalContainerReservedDisk = cellStats.TotalContainerReservedDisk + uint64(appMetadata.DiskQuotaMB*util.MEGABYTE)

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
