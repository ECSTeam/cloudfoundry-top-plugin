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

package appDetailView

import (
	"fmt"
	"log"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/org"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/space"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

type AppDetailView struct {
	*dataView.DataListView
	appId              string
	requestsInfoWidget *RequestsInfoWidget
}

func NewAppDetailView(masterUI masterUIInterface.MasterUIInterface,
	parentView dataView.DataListViewInterface,
	name string, bottomMargin int,
	eventProcessor *eventdata.EventProcessor,
	appId string) *AppDetailView {

	asUI := &AppDetailView{appId: appId}
	requestViewHeight := 5
	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("CPU_PERCENT", true),
		uiCommon.NewSortColumn("IDX", false),
	}

	dataListView := dataView.NewDataListView(masterUI, parentView,
		name, requestViewHeight+1, bottomMargin,
		eventProcessor, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.GetListData = asUI.GetListData
	dataListView.RefreshDisplayCallback = asUI.refreshDisplay
	dataListView.UpdateHeaderCallback = asUI.updateHeader

	dataListView.SetTitle("Container List")
	dataListView.HelpText = HelpText
	dataListView.HelpTextTips = HelpTextTips

	asUI.DataListView = dataListView

	asUI.requestsInfoWidget = NewRequestsInfoWidget(masterUI, "requestsInfoWidget", requestViewHeight, asUI)
	masterUI.LayoutManager().Add(asUI.requestsInfoWidget)

	return asUI
}

func (asUI *AppDetailView) initializeCallback(g *gocui.Gui, viewName string) error {
	if err := g.SetKeybinding(viewName, 'x', gocui.ModNone, asUI.closeAppDetailView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, gocui.KeyEsc, gocui.ModNone, asUI.closeAppDetailView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, 'i', gocui.ModNone, asUI.openInfoAction); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (asUI *AppDetailView) openInfoAction(g *gocui.Gui, v *gocui.View) error {
	infoWidgetName := "appInfoWidget"
	appInfoWidget := NewAppInfoWidget(asUI.GetMasterUI(), infoWidgetName, 70, 18, asUI)
	asUI.GetMasterUI().LayoutManager().Add(appInfoWidget)
	asUI.GetMasterUI().SetCurrentViewOnTop(g)
	asUI.GetMasterUI().AddCommonDataViewKeybindings(g, infoWidgetName)
	return nil
}

func (asUI *AppDetailView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, ColumnContainerIndex())
	columns = append(columns, ColumnTotalCpuPercentage())
	columns = append(columns, ColumnMemoryUsed())
	columns = append(columns, ColumnMemoryFree())
	columns = append(columns, ColumnDiskUsed())
	columns = append(columns, ColumnDiskFree())
	columns = append(columns, ColumnLogStdout())
	columns = append(columns, ColumnLogStderr())
	columns = append(columns, ColumnCellIp())
	return columns
}

func (asUI *AppDetailView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *AppDetailView) postProcessData() []*displaydata.DisplayContainerStats {

	displayStatsArray := make([]*displaydata.DisplayContainerStats, 0)

	appMap := asUI.GetDisplayedEventData().AppMap
	appStats := appMap[asUI.appId]
	if appStats == nil {
		return displayStatsArray
	}

	appMetadata := asUI.GetAppMdMgr().FindAppMetadata(appStats.AppId)

	for _, containerStats := range appStats.ContainerArray {
		if containerStats != nil {
			displayContainerStats := displaydata.NewDisplayContainerStats(containerStats, appStats)

			/*
				appMetadata := appMdMgr.FindAppMetadata(appStats.AppId)
				appName := appMetadata.Name
				if appName == "" {
					appName = appStats.AppId
				}
			*/

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
			displayStatsArray = append(displayStatsArray, displayContainerStats)
		}
	}
	return displayStatsArray
}

func (asUI *AppDetailView) convertToListData(containerStatsArray []*displaydata.DisplayContainerStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(containerStatsArray))
	for _, d := range containerStatsArray {
		listData = append(listData, d)
	}
	return listData
}

func (asUI *AppDetailView) closeAppDetailView(g *gocui.Gui, v *gocui.View) error {
	if err := asUI.GetMasterUI().CloseView(asUI); err != nil {
		return err
	}
	if err := asUI.GetMasterUI().CloseView(asUI.requestsInfoWidget); err != nil {
		return err
	}
	return nil
}

func (w *AppDetailView) refreshDisplay(g *gocui.Gui) error {

	// HTTP request stats -- These stands are also on the appListView so we need them in a detail view??
	/*
		fmt.Fprintf(v, "\n")
		fmt.Fprintf(v, "HTTP(S) status code:\n")
		fmt.Fprintf(v, "  2xx: %12v\n", util.Format(appStats.TotalTraffic.Http2xxCount))
		fmt.Fprintf(v, "  3xx: %12v\n", util.Format(appStats.TotalTraffic.Http3xxCount))
		fmt.Fprintf(v, "  4xx: %12v\n", util.Format(appStats.TotalTraffic.Http4xxCount))
		fmt.Fprintf(v, "  5xx: %12v\n", util.Format(appStats.TotalTraffic.Http5xxCount))
		fmt.Fprintf(v, "%v", util.BRIGHT_WHITE)
		fmt.Fprintf(v, "  All: %12v\n", util.Format(appStats.TotalTraffic.HttpAllCount))
		fmt.Fprintf(v, "%v", util.CLEAR)
	*/

	/*
		totalLogCount = totalLogCount + appStats.NonContainerOutCount + appStats.NonContainerErrCount
		fmt.Fprintf(v, "Non container logs - Stdout: %-12v ", util.Format(appStats.NonContainerOutCount))
		fmt.Fprintf(v, "Stderr: %-12v\n", util.Format(appStats.NonContainerErrCount))
		fmt.Fprintf(v, "Total log events: %12v\n", util.Format(totalLogCount))
	*/
	return nil
}

func (asUI *AppDetailView) updateHeader(g *gocui.Gui, v *gocui.View) (int, error) {
	fmt.Fprintf(v, "\nTODO: Show summary stats")
	return 3, nil
}
