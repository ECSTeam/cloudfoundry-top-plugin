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

package capacityPlanView

import (
	"log"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventCell"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appViews/appView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/cellViews/cellDetailView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/cellViews/cellView"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

const UNKNOWN = -1
const DUMMY_CELL_NAME_FOR_TOTAL = "TOTAL"

type CapacityPlanView struct {
	*dataView.DataListView
	//detailView *cellDetailView.CellDetailView
}

func NewCapacityPlanView(masterUI masterUIInterface.MasterUIInterface,
	name string, bottomMargin int,
	eventProcessor *eventdata.EventProcessor) *CapacityPlanView {

	asUI := &CapacityPlanView{}

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("CELL_IP", false),
	}

	dataListView := dataView.NewDataListView(masterUI, nil,
		name, 0, bottomMargin,
		eventProcessor, asUI, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	//dataListView.UpdateHeaderCallback = asUI.updateHeader
	dataListView.GetListData = asUI.GetListData
	//dataListView.PreRowDisplayCallback = asUI.preRowDisplay

	dataListView.SetTitle(func() string { return "Capacity Plan (memory)" })
	dataListView.HelpText = HelpText
	dataListView.HelpTextTips = appView.HelpTextTips

	asUI.DataListView = dataListView

	return asUI

}

func (asUI *CapacityPlanView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, columnIp())

	//columns = append(columns, columnNumOfCpus())

	columns = append(columns, columnCapacityMemoryTotal())
	columns = append(columns, columnCapacityMemoryRemaining())
	columns = append(columns, columnTotalContainerMemoryReserved())

	columns = append(columns, columnCapacityTotalContainers())
	columns = append(columns, columnContainerCount())
	columns = append(columns, columnCapacityPlan0_5GMem())
	columns = append(columns, columnCapacityPlan1_0GMem())
	columns = append(columns, columnCapacityPlan1_5GMem())
	columns = append(columns, columnCapacityPlan2_0GMem())
	columns = append(columns, columnCapacityPlan2_5GMem())
	columns = append(columns, columnCapacityPlan3_0GMem())
	columns = append(columns, columnCapacityPlan3_5GMem())
	columns = append(columns, columnCapacityPlan4_0GMem())

	return columns
}

func (asUI *CapacityPlanView) initializeCallback(g *gocui.Gui, viewName string) error {

	if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone, asUI.enterAction); err != nil {
		log.Panicln(err)
	}

	return nil
}

func (asUI *CapacityPlanView) enterAction(g *gocui.Gui, v *gocui.View) error {
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

func (asUI *CapacityPlanView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *CapacityPlanView) postProcessData() map[string]*cellView.DisplayCellStats {
	cellMap := asUI.GetDisplayedEventData().CellMap

	displayCellMap := make(map[string]*cellView.DisplayCellStats)
	for ip, cellStats := range cellMap {
		displayStat := cellView.NewDisplayCellStats(cellStats)
		displayCellMap[ip] = displayStat

		if cellStats.CapacityMemoryTotal > 0 {
			displayStat.CapacityPlan0_5GMem = int(cellStats.CapacityMemoryRemaining / (util.GIGABYTE * 0.5))
			displayStat.CapacityPlan1_0GMem = int(cellStats.CapacityMemoryRemaining / (util.GIGABYTE * 1))
			displayStat.CapacityPlan1_5GMem = int(cellStats.CapacityMemoryRemaining / (util.GIGABYTE * 1.5))
			displayStat.CapacityPlan2_0GMem = int(cellStats.CapacityMemoryRemaining / (util.GIGABYTE * 2))
			displayStat.CapacityPlan2_5GMem = int(cellStats.CapacityMemoryRemaining / (util.GIGABYTE * 2.5))
			displayStat.CapacityPlan3_0GMem = int(cellStats.CapacityMemoryRemaining / (util.GIGABYTE * 3))
			displayStat.CapacityPlan3_5GMem = int(cellStats.CapacityMemoryRemaining / (util.GIGABYTE * 3.5))
			displayStat.CapacityPlan4_0GMem = int(cellStats.CapacityMemoryRemaining / (util.GIGABYTE * 4))
		} else {
			displayStat.CapacityPlan0_5GMem = UNKNOWN
			displayStat.CapacityPlan1_0GMem = UNKNOWN
			displayStat.CapacityPlan1_5GMem = UNKNOWN
			displayStat.CapacityPlan2_0GMem = UNKNOWN
			displayStat.CapacityPlan2_5GMem = UNKNOWN
			displayStat.CapacityPlan3_0GMem = UNKNOWN
			displayStat.CapacityPlan3_5GMem = UNKNOWN
			displayStat.CapacityPlan4_0GMem = UNKNOWN
		}

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

						appMetadata := asUI.GetMdGlobalMgr().GetAppMdManager().FindItem(appStats.AppId)
						cellStats.TotalContainerMemoryReserved = cellStats.TotalContainerMemoryReserved + uint64(appMetadata.MemoryMB*util.MEGABYTE)

						usedMemoryValue := containerStats.ContainerMetric.GetMemoryBytes()
						cellStats.TotalContainerMemoryUsed = cellStats.TotalContainerMemoryUsed + usedMemoryValue

					}
				}
			}
		}
	}

	asUI.addTotalRow(displayCellMap)

	return displayCellMap
}

func (asUI *CapacityPlanView) addTotalRow(displayCellMap map[string]*cellView.DisplayCellStats) {

	totalLabel := DUMMY_CELL_NAME_FOR_TOTAL
	totalCellStats := eventCell.NewCellStats(totalLabel)
	totalDisplayStat := cellView.NewDisplayCellStats(totalCellStats)
	displayCellMap[totalLabel] = totalDisplayStat

	totalDisplayStat.CapacityPlan0_5GMem = UNKNOWN
	totalDisplayStat.CapacityPlan1_0GMem = UNKNOWN
	totalDisplayStat.CapacityPlan1_5GMem = UNKNOWN
	totalDisplayStat.CapacityPlan2_0GMem = UNKNOWN
	totalDisplayStat.CapacityPlan2_5GMem = UNKNOWN
	totalDisplayStat.CapacityPlan3_0GMem = UNKNOWN
	totalDisplayStat.CapacityPlan3_5GMem = UNKNOWN
	totalDisplayStat.CapacityPlan4_0GMem = UNKNOWN

	//NumOfCpus := 0
	CapacityMemoryTotal := int64(0)
	CapacityMemoryRemaining := int64(0)
	TotalContainerMemoryReserved := uint64(0)
	CapacityTotalContainers := 0
	ContainerCount := 0
	CapacityPlan0_5GMem := 0
	CapacityPlan1_0GMem := 0
	CapacityPlan1_5GMem := 0
	CapacityPlan2_0GMem := 0
	CapacityPlan2_5GMem := 0
	CapacityPlan3_0GMem := 0
	CapacityPlan3_5GMem := 0
	CapacityPlan4_0GMem := 0

	capacityPlanHasValue := false
	for _, cellStats := range displayCellMap {

		//NumOfCpus = NumOfCpus + cellStats.NumOfCpus
		CapacityMemoryTotal = CapacityMemoryTotal + cellStats.CapacityMemoryTotal
		CapacityMemoryRemaining = CapacityMemoryRemaining + cellStats.CapacityMemoryRemaining
		TotalContainerMemoryReserved = TotalContainerMemoryReserved + cellStats.TotalContainerMemoryReserved
		CapacityTotalContainers = CapacityTotalContainers + cellStats.CapacityTotalContainers
		ContainerCount = ContainerCount + cellStats.ContainerCount
		if cellStats.CapacityMemoryTotal > 0 {
			capacityPlanHasValue = true
			CapacityPlan0_5GMem = CapacityPlan0_5GMem + cellStats.CapacityPlan0_5GMem
			CapacityPlan1_0GMem = CapacityPlan1_0GMem + cellStats.CapacityPlan1_0GMem
			CapacityPlan1_5GMem = CapacityPlan1_5GMem + cellStats.CapacityPlan1_5GMem
			CapacityPlan2_0GMem = CapacityPlan2_0GMem + cellStats.CapacityPlan2_0GMem
			CapacityPlan2_5GMem = CapacityPlan2_5GMem + cellStats.CapacityPlan2_5GMem
			CapacityPlan3_0GMem = CapacityPlan3_0GMem + cellStats.CapacityPlan3_0GMem
			CapacityPlan3_5GMem = CapacityPlan3_5GMem + cellStats.CapacityPlan3_5GMem
			CapacityPlan4_0GMem = CapacityPlan4_0GMem + cellStats.CapacityPlan4_0GMem
		}
	}
	//totalDisplayStat.NumOfCpus = NumOfCpus
	totalDisplayStat.CapacityMemoryTotal = CapacityMemoryTotal
	totalDisplayStat.CapacityMemoryRemaining = CapacityMemoryRemaining
	totalDisplayStat.TotalContainerMemoryReserved = TotalContainerMemoryReserved
	totalDisplayStat.CapacityTotalContainers = CapacityTotalContainers
	totalDisplayStat.ContainerCount = ContainerCount

	if capacityPlanHasValue {
		totalDisplayStat.CapacityPlan0_5GMem = CapacityPlan0_5GMem
		totalDisplayStat.CapacityPlan1_0GMem = CapacityPlan1_0GMem
		totalDisplayStat.CapacityPlan1_5GMem = CapacityPlan1_5GMem
		totalDisplayStat.CapacityPlan2_0GMem = CapacityPlan2_0GMem
		totalDisplayStat.CapacityPlan2_5GMem = CapacityPlan2_5GMem
		totalDisplayStat.CapacityPlan3_0GMem = CapacityPlan3_0GMem
		totalDisplayStat.CapacityPlan3_5GMem = CapacityPlan3_5GMem
		totalDisplayStat.CapacityPlan4_0GMem = CapacityPlan4_0GMem
	}

}

func (asUI *CapacityPlanView) convertToListData(displayCellMap map[string]*cellView.DisplayCellStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(displayCellMap))
	for _, d := range displayCellMap {
		listData = append(listData, d)
	}
	return listData
}

/*
func (asUI *CapacityPlanView) preRowDisplay(data uiCommon.IData, isSelected bool) string {
	cellStats := data.(*cellView.DisplayCellStats)
	v := bytes.NewBufferString("")
	if !isSelected && cellStats.Ip == DUMMY_CELL_NAME_FOR_TOTAL {
		fmt.Fprintf(v, util.BRIGHT_WHITE)
	}
	return v.String()
}
*/
