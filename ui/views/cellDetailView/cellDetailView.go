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

package cellDetailView

import (
	"fmt"
	"log"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventApp"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/org"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/space"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appDetailView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
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

	dataListView.SetTitle(fmt.Sprintf("Cell IP:%v Detail - Container List", cellIp))
	dataListView.HelpText = HelpText
	dataListView.HelpTextTips = HelpTextTips

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

	columns = append(columns, appDetailView.ColumnMemoryReserved())
	columns = append(columns, appDetailView.ColumnMemoryUsed())
	columns = append(columns, appDetailView.ColumnMemoryFree())

	columns = append(columns, appDetailView.ColumnDiskReserved())
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
	appStatsArray := eventApp.ConvertFromMap(appMap, asUI.GetAppMdMgr())
	for _, appStats := range appStatsArray {
		appMetadata := asUI.GetAppMdMgr().FindAppMetadata(appStats.AppId)
		for _, containerStats := range appStats.ContainerArray {
			if containerStats != nil {
				if containerStats.Ip == asUI.cellIp {
					// This is a container on the selected cell
					displayContainerStats := displaydata.NewDisplayContainerStats(containerStats, appStats)
					displayContainerStats.AppName = appMetadata.Name
					displayContainerStats.SpaceName = space.FindSpaceName(appMetadata.SpaceGuid)
					displayContainerStats.OrgName = org.FindOrgNameBySpaceGuid(appMetadata.SpaceGuid)

					usedMemory := containerStats.ContainerMetric.GetMemoryBytes()
					reservedMemory := uint64(appMetadata.MemoryMB) * util.MEGABYTE
					freeMemory := reservedMemory - usedMemory
					displayContainerStats.FreeMemory = freeMemory
					displayContainerStats.ReservedMemory = reservedMemory

					usedDisk := containerStats.ContainerMetric.GetDiskBytes()
					reservedDisk := uint64(appMetadata.DiskQuotaMB) * util.MEGABYTE
					freeDisk := reservedDisk - usedDisk
					displayContainerStats.FreeDisk = freeDisk
					displayContainerStats.ReservedDisk = reservedDisk

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

func (asUI *CellDetailView) updateHeader(g *gocui.Gui, v *gocui.View) (int, error) {
	fmt.Fprintf(v, "\nTODO: Show summary stats")
	return 3, nil
}
