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
	"strconv"
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventApp"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/crashData"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appViews/appCrashView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appViews/appHttpView"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

type AppDetailView struct {
	*dataView.DataListView
	appId              string
	requestsInfoWidget *RequestsInfoWidget
	crashInfoWidget    *CrashInfoWidget
	displayMenuId      string

	Crash10mCount int
	Crash1hCount  int
	Crash24hCount int
	LastCrashInfo *crashData.ContainerCrashInfo
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
	//dataListView.RefreshDisplayCallback = asUI.refreshDisplay

	dataListView.SetTitle(func() string { return "Container List" })

	dataListView.HelpText = HelpText
	dataListView.HelpTextTips = HelpTextTips

	asUI.DataListView = dataListView

	asUI.requestsInfoWidget = NewRequestsInfoWidget(masterUI, dataListView, "requestsInfoWidget", requestViewHeight, asUI)
	masterUI.LayoutManager().Add(asUI.requestsInfoWidget)

	asUI.crashInfoWidget = NewCrashInfoWidget(masterUI, dataListView, "crashInfoWidget", requestViewHeight, asUI)
	masterUI.LayoutManager().Add(asUI.crashInfoWidget)

	return asUI
}

func (asUI *AppDetailView) initializeCallback(g *gocui.Gui, viewName string) error {
	if err := g.SetKeybinding(viewName, 'x', gocui.ModNone, asUI.closeAppDetailView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, gocui.KeyEsc, gocui.ModNone, asUI.closeAppDetailView); err != nil {
		log.Panicln(err)
	}

	copyMenu := NewCopyMenu(asUI.GetMasterUI(), asUI)
	if err := g.SetKeybinding(viewName, 'c', gocui.ModNone, copyMenu.CopyAction); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(viewName, 'd', gocui.ModNone, asUI.selectDisplayAction); err != nil {
		log.Panicln(err)
	}
	/*
		if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone, asUI.enterAction); err != nil {
			log.Panicln(err)
		}
	*/
	return nil
}

func (asUI *AppDetailView) GetAppId() string {
	return asUI.appId
}

func (asUI *AppDetailView) selectDisplayAction(g *gocui.Gui, v *gocui.View) error {

	menuItems := make([]*uiCommon.MenuItem, 0, 5)
	menuItems = append(menuItems, uiCommon.NewMenuItem("infoView", "App Info"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("crashInfoView", "View CRASH List"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("appHttpView", "HTTP Response Info"))

	windowTitle := fmt.Sprintf("Select App Detail View")
	selectDisplayView := uiCommon.NewSelectMenuWidget(asUI.GetMasterUI(), "selectDisplayView", windowTitle, menuItems, asUI.selectDisplayCallback)
	selectDisplayView.SetMenuId(asUI.displayMenuId)

	asUI.GetMasterUI().LayoutManager().Add(selectDisplayView)
	asUI.GetMasterUI().SetCurrentViewOnTop(g)
	return nil
}

func (asUI *AppDetailView) enterAction(g *gocui.Gui, v *gocui.View) error {

	highlightKey := asUI.GetListWidget().HighlightKey()
	if highlightKey != "" {
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
		view = NewAppInfoWidget(asUI.GetMasterUI(), asUI, infoWidgetName, 70, 20, asUI)
	case "crashInfoView":
		_, bottomMargin := asUI.GetMargins()
		view = appCrashView.NewAppCrashView(asUI.GetMasterUI(), asUI, "crashInfoView", bottomMargin,
			asUI.GetEventProcessor(),
			asUI.appId)
	case "appHttpView":
		_, bottomMargin := asUI.GetMargins()
		view = appHttpView.NewAppHttpView(asUI.GetMasterUI(), asUI, "appHttpView", bottomMargin,
			asUI.GetEventProcessor(),
			asUI.appId)
	default:
		return errors.New("Unable to find view " + viewName)
	}
	return asUI.GetMasterUI().OpenView(g, view)
}

func (asUI *AppDetailView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, ColumnContainerIndex())

	columns = append(columns, ColumnState())
	columns = append(columns, ColumnStateDuration())

	columns = append(columns, ColumnTotalCpuPercentage())
	columns = append(columns, ColumnMemoryUsed())
	columns = append(columns, ColumnMemoryFree())
	columns = append(columns, ColumnDiskUsed())
	columns = append(columns, ColumnDiskFree())
	columns = append(columns, columnAvgResponseTimeL60Info())
	columns = append(columns, ColumnLogStdout())
	columns = append(columns, ColumnLogStderr())

	columns = append(columns, columnReq1())
	columns = append(columns, columnReq10())
	columns = append(columns, columnReq60())

	columns = append(columns, columnTotalReq())
	columns = append(columns, column2XX())
	columns = append(columns, column3XX())
	columns = append(columns, column4XX())
	columns = append(columns, column5XX())

	columns = append(columns, ColumnCellIp())

	columns = append(columns, ColumnStartupDuration())
	columns = append(columns, ColumnCreateCount())
	columns = append(columns, ColumnStateTime())

	columns = append(columns, ColumnCellLastStartMsgText())
	columns = append(columns, ColumnCellLastStartMsgTime())
	return columns
}

func (asUI *AppDetailView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *AppDetailView) postProcessData() []*DisplayContainerStats {

	displayStatsArray := make([]*DisplayContainerStats, 0)
	displayContainerStatsMap := make(map[int]*DisplayContainerStats)

	now := time.Now().Truncate(time.Second)

	appMap := asUI.GetDisplayedEventData().AppMap

	appId := asUI.appId
	mdMgr := asUI.GetMdGlobalMgr()
	mdAppMgr := mdMgr.GetAppMdManager()
	appMetadata := mdAppMgr.FindItem(appId)
	if mdAppMgr.IsPendingDeleteFromCache(appId) || mdAppMgr.IsDeletedFromCache(appId) {
		asUI.DataListView.RefreshDisplayCallback = asUI.refreshDisplay
		return nil
	}

	appStats := appMap[appId]
	if appStats == nil {
		return displayStatsArray
	}
	asUI.GetEventProcessor().GetMetadataManager().MonitorAppDetails(appId, &now)

	appInsts := asUI.GetEventProcessor().GetMetadataManager().GetAppInstMdManager().FindItem(appId)
	if appInsts == nil && appMetadata.State == "STARTED" {
		// Update the app instance statistics
		asUI.GetEventProcessor().GetMetadataManager().RequestRefreshAppInstancesMetadata(appId)
		toplog.Debug("No app inst data loaded yet")
	}

	// If we have app instance list from metadata -- populate map with current state / uptime
	if appInsts != nil && appInsts.Data != nil {
		for containerIndexStr, appInstStats := range appInsts.Data {
			if appInstStats != nil {
				containerIndex, err := strconv.Atoi(containerIndexStr)
				if err != nil {
					// TODO
				}
				placeHolder := eventApp.NewContainerStats(containerIndex)
				displayContainerStats := NewDisplayContainerStats(placeHolder, appStats)

				displayContainerStatsMap[containerIndex] = displayContainerStats
				displayContainerStats.State = appInstStats.State
				stateTime := appInstStats.StartTime
				displayContainerStats.StateTime = stateTime
				if stateTime != nil {
					stateDuration := now.Sub(*stateTime)
					displayContainerStats.StateDuration = &stateDuration
				}
			}
		}
	}

	// Loop through any reporting containers and populate cpu/memory/disk info
	for _, containerStats := range appStats.ContainerArray {

		if containerStats != nil {
			displayContainerStats := displayContainerStatsMap[containerStats.ContainerIndex]
			if displayContainerStats == nil {
				// We have container stats but /v2/app/<GUID>/stats doesn't know about it
				// So either the app was scaled down (container terminated) or we've scaled up
				// and we got container messages and we're working with stale data from /v2/app/<GUID>/stats
				displayContainerStats = NewDisplayContainerStats(containerStats, appStats)
				displayContainerStatsMap[containerStats.ContainerIndex] = displayContainerStats
				if appInsts != nil {
					if containerStats.CellLastCreatingMsgTime == nil || (containerStats.CellLastCreatingMsgTime != nil &&
						appInsts.GetCacheTime().After(*containerStats.CellLastCreatingMsgTime)) {
						displayContainerStats.State = "TERM"
					} else {
						displayContainerStats.State = "UNKNOWN"
					}
				} else {
					displayContainerStats.State = "CHECKING"
				}
			} else {
				displayContainerStats.ContainerStats = containerStats
			}

			displayContainerStats.AppName = appMetadata.Name
			spaceMd := mdMgr.GetSpaceMdManager().FindItem(appMetadata.SpaceGuid)
			displayContainerStats.SpaceName = spaceMd.Name
			orgMd := mdMgr.GetOrgMdManager().FindItem(spaceMd.OrgGuid)
			displayContainerStats.OrgName = orgMd.Name

			if containerStats.CellCreatedMsgTime != nil && containerStats.CellHealthyMsgTime != nil {
				startupDuration := containerStats.CellHealthyMsgTime.Sub(*containerStats.CellCreatedMsgTime)
				if startupDuration > 0 {
					displayContainerStats.StartupDuration = &startupDuration
				}
			}

			if displayContainerStats.State == "DOWN" {
				containerStats.ContainerMetric = nil
				containerStats.Ip = ""
			}

			if containerStats.ContainerMetric != nil {
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
			}

		}
	}

	for appInstId, containerTraffic := range appStats.ContainerTrafficMap {

		// TODO: Since we collect container info by appInst GUID and not index
		// we could have multiple entries in the map for the same index in
		// the case where a container restarted.  We only want to display
		// states for the currently "alive" container -- how do we determine that?
		appInstId = appInstId

		//toplog.Info("appInstId: %v index: %v", appInstId, containerTraffic.InstanceIndex)

		displayContainerStats := displayContainerStatsMap[int(containerTraffic.InstanceIndex)]
		if displayContainerStats != nil {

			displayContainerStats.AvgResponseL60Time = containerTraffic.AvgResponseL60Time
			displayContainerStats.EventL1Rate = containerTraffic.EventL1Rate
			displayContainerStats.EventL10Rate = containerTraffic.EventL10Rate
			displayContainerStats.EventL60Rate = containerTraffic.EventL60Rate

			for _, httpStatusCodeMap := range containerTraffic.HttpInfoMap {
				for statusCode, httpCountInfo := range httpStatusCodeMap {
					if httpCountInfo != nil {
						displayContainerStats.HttpAllCount += httpCountInfo.HttpCount
						switch {
						case statusCode >= 200 && statusCode < 300:
							displayContainerStats.Http2xxCount += httpCountInfo.HttpCount
						case statusCode >= 300 && statusCode < 400:
							displayContainerStats.Http3xxCount += httpCountInfo.HttpCount
						case statusCode >= 400 && statusCode < 500:
							displayContainerStats.Http4xxCount += httpCountInfo.HttpCount
						case statusCode >= 500 && statusCode < 600:
							displayContainerStats.Http5xxCount += httpCountInfo.HttpCount
						}
					}
				}
			}
		}

	}

	for containerIndex, displayContainerStats := range displayContainerStatsMap {

		// Populate the DOWN reason
		if displayContainerStats.State == "DOWN" {
			if appInsts != nil && appInsts.Data != nil {
				instance := appInsts.Data[strconv.Itoa(containerIndex)]
				if instance != nil {
					displayContainerStats.CellLastStartMsgText = instance.Details
					displayContainerStats.CellLastStartMsgTime = nil
				}
			}
		}

		displayStatsArray = append(displayStatsArray, displayContainerStats)
	}

	displayAppStatsMap := asUI.GetMasterUI().GetCommonData().GetDisplayAppStatsMap()
	displayAppStats := displayAppStatsMap[asUI.appId]

	asUI.Crash1hCount = displayAppStats.Crash1hCount
	asUI.Crash24hCount = displayAppStats.Crash24hCount

	crash10mCount := crashData.FindCountSinceByApp(appStats.AppId, -10*time.Minute)
	crash10mCount = crash10mCount + appStats.CrashCountSince(-10*time.Minute)
	asUI.Crash10mCount = crash10mCount

	if displayAppStats.Crash24hCount > 0 {
		// Lookup crash time from container stats
		asUI.LastCrashInfo = asUI.FindLastCrash(appStats)
		if asUI.LastCrashInfo == nil {
			// If we don't find last crash in container stats, last crash must have occured
			// before top was started.  Look for last crash time in metadata (/v2/event data)
			asUI.LastCrashInfo = crashData.FindLastCrashByApp(appStats.AppId)
		}
	}
	return displayStatsArray
}

func (asUI *AppDetailView) FindLastCrash(appStats *eventApp.AppStats) *crashData.ContainerCrashInfo {
	if appStats.ContainerCrashInfo != nil && len(appStats.ContainerCrashInfo) > 0 {
		last := len(appStats.ContainerCrashInfo) - 1
		return appStats.ContainerCrashInfo[last]
	}
	return nil
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
	if err := asUI.GetMasterUI().CloseView(asUI.crashInfoWidget); err != nil {
		return err
	}
	return nil
}

func (w *AppDetailView) refreshDisplay(g *gocui.Gui) (propagateRefresh bool, err error) {

	name := w.DataListView.Name()

	v, err := g.View(name)
	if err != nil {
		return false, err
	}
	AppDeletedMsg(g, v)

	return false, nil
}

func AppDeletedMsg(g *gocui.Gui, v *gocui.View) {
	v.Clear()
	fmt.Fprintf(v, " \n")
	fmt.Fprintf(v, "%v", util.BRIGHT_RED)
	fmt.Fprintf(v, " Application has been deleted")
	fmt.Fprintf(v, "%v", util.CLEAR)
}
