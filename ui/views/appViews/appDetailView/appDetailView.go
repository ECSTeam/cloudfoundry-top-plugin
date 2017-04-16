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
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/crashData"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/org"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/space"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

type AppDetailView struct {
	*dataView.DataListView
	appId              string
	requestsInfoWidget *RequestsInfoWidget
	displayMenuId      string
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
		eventProcessor, asUI, asUI.columnDefinitions(),
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
	if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone, asUI.enterAction); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (asUI *AppDetailView) enterAction(g *gocui.Gui, v *gocui.View) error {

	highlightKey := asUI.GetListWidget().HighlightKey()
	if highlightKey != "" {
		/*
			_, bottomMargin := asUI.GetMargins()

			detailView := appDetailView.NewAppDetailView(asUI.GetMasterUI(), asUI, "appDetailView",
				bottomMargin,
				asUI.GetEventProcessor(),
				highlightKey)
			asUI.SetDetailView(detailView)
			asUI.GetMasterUI().OpenView(g, detailView)
		*/
		toplog.Info("highlightKey: [%v]", highlightKey)

		menuItems := make([]*uiCommon.MenuItem, 0, 5)
		menuItems = append(menuItems, uiCommon.NewMenuItem("infoView", "View CRASH info"))
		menuItems = append(menuItems, uiCommon.NewMenuItem("infoView", "View Logs"))
		menuItems = append(menuItems, uiCommon.NewMenuItem("infoView", "Todo"))

		windowTitle := fmt.Sprintf("Select View for %v", highlightKey)
		selectDisplayView := uiCommon.NewSelectMenuWidget(asUI.GetMasterUI(), "selectDisplayView", windowTitle, menuItems, asUI.selectDisplayCallback)
		selectDisplayView.SetMenuId(asUI.displayMenuId)

		asUI.GetMasterUI().LayoutManager().Add(selectDisplayView)
		asUI.GetMasterUI().SetCurrentViewOnTop(g)

	}
	return nil
}

func (asUI *AppDetailView) selectDisplayCallback(g *gocui.Gui, v *gocui.View, menuId string) error {
	asUI.displayMenuId = menuId
	asUI.createAndOpenView(g, menuId)
	return nil
}

func (asUI *AppDetailView) createAndOpenView(g *gocui.Gui, viewName string) error {

	var view masterUIInterface.UpdatableView
	switch viewName {
	case "infoView":
		infoWidgetName := "appInfoWidget"
		view = NewAppInfoWidget(asUI.GetMasterUI(), infoWidgetName, 70, 18, asUI)
	case "xxx":
		//view = orgView.NewOrgListView(mui, "orgListView", mui.helpTextTipsViewSize, ep)
	default:
		return errors.New("Unable to find view " + viewName)
	}
	return asUI.GetMasterUI().OpenView(g, view)
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
	columns = append(columns, ColumnCrash1hCount())
	columns = append(columns, ColumnCrash24hCount())
	columns = append(columns, ColumnLastCrash())

	columns = append(columns, ColumnCellIp())
	return columns
}

func (asUI *AppDetailView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *AppDetailView) postProcessData() []*DisplayContainerStats {

	displayStatsArray := make([]*DisplayContainerStats, 0)

	appMap := asUI.GetDisplayedEventData().AppMap
	appStats := appMap[asUI.appId]
	if appStats == nil {
		return displayStatsArray
	}

	appMetadata := asUI.GetAppMdMgr().FindAppMetadata(appStats.AppId)

	for _, containerStats := range appStats.ContainerArray {
		if containerStats != nil {
			displayContainerStats := NewDisplayContainerStats(containerStats, appStats)

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

			displayContainerStats.Crash1hCount = containerStats.Crash1hCount()
			displayContainerStats.Crash1hCount = displayContainerStats.Crash1hCount +
				crashData.FindCountSinceByAppAndInstance(appStats.AppId, containerStats.ContainerIndex, -1*time.Hour)

			displayContainerStats.Crash24hCount = containerStats.Crash24hCount()
			displayContainerStats.Crash24hCount = displayContainerStats.Crash24hCount +
				crashData.FindCountSinceByAppAndInstance(appStats.AppId, containerStats.ContainerIndex, -24*time.Hour)

			if displayContainerStats.Crash24hCount > 0 {
				// Lookup crash time from container stats
				displayContainerStats.LastCrashTime = containerStats.LastCrashTime()
				if displayContainerStats.LastCrashTime == nil {
					// If we don't find last crash in container stats, last crash must have occured
					// before top was started.  Look for last crash time in metadata (/v2/event data)
					displayContainerStats.LastCrashTime = crashData.FindLastCrashByAppAndInstance(appStats.AppId, containerStats.ContainerIndex)
				}
			}

		}
	}
	return displayStatsArray
}

func (asUI *AppDetailView) convertToListData(containerStatsArray []*DisplayContainerStats) []uiCommon.IData {
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
