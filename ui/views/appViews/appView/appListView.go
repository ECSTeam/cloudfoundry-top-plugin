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

package appView

import (
	"fmt"
	"log"
	"time"

	"github.com/atotto/clipboard"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventApp"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/org"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/space"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/stack"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/dataView"
	"github.com/jroimartin/gocui"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appViews/appDetailView"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

const StaleContainerSeconds = 80

type AppListView struct {
	*dataView.DataListView
	displayAppStats  []*DisplayAppStats
	isWarmupComplete bool
	// This is a count of the number of apps that do not have
	// the correct number of containers running based on app
	// instance setting
	appsNotInDesiredState int
}

func NewAppListView(masterUI masterUIInterface.MasterUIInterface,
	name string, bottomMargin int,
	eventProcessor *eventdata.EventProcessor) *AppListView {

	asUI := &AppListView{}

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("CPU", true),
		uiCommon.NewSortColumn("REQ60", true),
		uiCommon.NewSortColumn("appName", false),
		uiCommon.NewSortColumn("spaceName", false),
		uiCommon.NewSortColumn("orgName", false),
	}

	dataListView := dataView.NewDataListView(masterUI, nil,
		name, 0, bottomMargin,
		eventProcessor, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.PreRowDisplayCallback = asUI.preRowDisplay
	dataListView.UpdateHeaderCallback = asUI.updateHeader
	dataListView.GetListData = asUI.GetListData

	dataListView.SetTitle("App List")
	dataListView.HelpText = HelpText
	dataListView.HelpTextTips = HelpTextTips

	asUI.DataListView = dataListView

	return asUI

}

func (asUI *AppListView) initializeCallback(g *gocui.Gui, viewName string) error {

	if err := g.SetKeybinding(viewName, 'c', gocui.ModNone, asUI.copyAction); err != nil {
		log.Panicln(err)
	}

	// TODO: Testing -- remove later
	if err := g.SetKeybinding(viewName, 'z', gocui.ModNone, asUI.testShowUserMessage); err != nil {
		log.Panicln(err)
	}
	// TODO: Testing -- remove later
	if err := g.SetKeybinding(viewName, 'a', gocui.ModNone, asUI.testClearUserMessage); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone, asUI.enterAction); err != nil {
		log.Panicln(err)
	}

	return nil
}

func (asUI *AppListView) enterAction(g *gocui.Gui, v *gocui.View) error {
	highlightKey := asUI.GetListWidget().HighlightKey()
	if asUI.GetListWidget().HighlightKey() != "" {
		_, bottomMargin := asUI.GetMargins()

		detailView := appDetailView.NewAppDetailView(asUI.GetMasterUI(), asUI, "appDetailView",
			bottomMargin,
			asUI.GetEventProcessor(),
			highlightKey)
		asUI.SetDetailView(detailView)
		asUI.GetMasterUI().OpenView(g, detailView)
	}
	return nil
}

func (asUI *AppListView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, asUI.columnAppName())
	columns = append(columns, asUI.columnSpaceName())
	columns = append(columns, asUI.columnOrgName())

	columns = append(columns, asUI.columnDesiredInstances())
	columns = append(columns, asUI.columnReportingContainers())

	columns = append(columns, asUI.columnTotalCpu())
	columns = append(columns, asUI.columnTotalMemoryUsed())
	columns = append(columns, asUI.columnTotalDiskUsed())

	columns = append(columns, asUI.columnAvgResponseTimeL60Info())
	columns = append(columns, asUI.columnLogStdout())
	columns = append(columns, asUI.columnLogStderr())

	columns = append(columns, asUI.columnReq1())
	columns = append(columns, asUI.columnReq10())
	columns = append(columns, asUI.columnReq60())

	columns = append(columns, asUI.columnTotalReq())
	columns = append(columns, asUI.column2XX())
	columns = append(columns, asUI.column3XX())
	columns = append(columns, asUI.column4XX())
	columns = append(columns, asUI.column5XX())

	columns = append(columns, asUI.columnStackName())

	return columns
}

func (asUI *AppListView) isUserMessageOpen(g *gocui.Gui) bool {
	alertViewName := "alertView"
	view := asUI.GetMasterUI().LayoutManager().GetManagerByViewName(alertViewName)
	if view != nil {
		return true
	} else {
		return false
	}
}

func (asUI *AppListView) clearUserMessage(g *gocui.Gui) error {
	alertViewName := "alertView"
	view := asUI.GetMasterUI().LayoutManager().GetManagerByViewName(alertViewName)
	if view != nil {
		asUI.GetMasterUI().CloseView(view)
	}
	asUI.SetAlertSize(0)
	return nil
}

// TODO: Have message levels which will colorize differently
func (asUI *AppListView) showUserMessage(g *gocui.Gui, message string) error {
	alertViewName := "alertView"
	alertHeight := 1

	asUI.SetAlertSize(alertHeight)

	var alertView *uiCommon.AlertWidget
	view := asUI.GetMasterUI().LayoutManager().GetManagerByViewName(alertViewName)
	if view == nil {
		alertView = uiCommon.NewAlertWidget(asUI.GetMasterUI(), alertViewName, alertHeight)
		asUI.GetMasterUI().LayoutManager().AddToBack(alertView)
		asUI.GetMasterUI().SetCurrentViewOnTop(g)
	} else {
		// This check is to prevent alert from showing on top of the log window
		if asUI.GetMasterUI().GetCurrentView(g).Name() == asUI.Name() {
			if _, err := g.SetViewOnTop(alertViewName); err != nil {
				return err
			}
		}
		alertView = view.(*uiCommon.AlertWidget)
		alertView.SetHeight(alertHeight)
	}
	alertView.SetMessage(message)
	return nil
}

func (asUI *AppListView) testShowUserMessage(g *gocui.Gui, v *gocui.View) error {
	return asUI.showUserMessage(g, "ALERT: 1 application(s) not in desired state (EXAMPLE) ")
	//asUI.GetEventProcessor().GetMetadataManager().RequestRefreshAppMetadata("9d82ef1b-4bba-4a49-9768-4ccd817edf9c")
	//return nil
}

func (asUI *AppListView) testClearUserMessage(g *gocui.Gui, v *gocui.View) error {
	return asUI.clearUserMessage(g)
}

func (asUI *AppListView) copyAction(g *gocui.Gui, v *gocui.View) error {

	selectedAppId := asUI.GetListWidget().HighlightKey()
	if selectedAppId == "" {
		// Nothing selected
		return nil
	}
	menuItems := make([]*uiCommon.MenuItem, 0, 5)
	menuItems = append(menuItems, uiCommon.NewMenuItem("cftarget", "cf target"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("cfapp", "cf app"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("cfscale", "cf scale"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("appguid", "app guid"))
	masterUI := asUI.GetMasterUI()
	clipboardView := uiCommon.NewSelectMenuWidget(masterUI, "clipboardView", "Copy to Clipboard", menuItems, asUI.clipboardCallback)

	masterUI.LayoutManager().Add(clipboardView)
	masterUI.SetCurrentViewOnTop(g)
	return nil
}

func (asUI *AppListView) clipboardCallback(g *gocui.Gui, v *gocui.View, menuId string) error {

	clipboardValue := ""

	selectedAppId := asUI.GetListWidget().HighlightKey()
	statsMap := asUI.GetDisplayedEventData().AppMap
	appStats := statsMap[selectedAppId]
	if appStats == nil {
		// Nothing selected
		return nil
	}
	appMetadata := asUI.GetAppMdMgr().FindAppMetadata(selectedAppId)
	appName := appMetadata.Name
	spaceName := space.FindSpaceName(appMetadata.SpaceGuid)
	orgName := org.FindOrgNameBySpaceGuid(appMetadata.SpaceGuid)

	switch menuId {
	case "cftarget":
		clipboardValue = fmt.Sprintf("cf target -o %v -s %v", orgName, spaceName)
	case "cfapp":
		clipboardValue = fmt.Sprintf("cf app %v", appName)
	case "cfscale":
		clipboardValue = fmt.Sprintf("cf scale %v ", appName)
	case "appguid":
		clipboardValue = selectedAppId
	}
	err := clipboard.WriteAll(clipboardValue)
	if err != nil {
		toplog.Error("Copy into Clipboard error: " + err.Error())
	}
	return nil
}

func (asUI *AppListView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *AppListView) postProcessData() []*DisplayAppStats {

	displayStatsArray := make([]*DisplayAppStats, 0)
	appMap := asUI.GetDisplayedEventData().AppMap
	appStatsArray := eventApp.ConvertFromMap(appMap, asUI.GetAppMdMgr())
	appsNotInDesiredState := 0

	for _, appStats := range appStatsArray {
		displayAppStats := NewDisplayAppStats(appStats)
		displayStatsArray = append(displayStatsArray, displayAppStats)
		appMetadata := asUI.GetAppMdMgr().FindAppMetadata(appStats.AppId)

		displayAppStats.AppName = appMetadata.Name
		displayAppStats.SpaceName = space.FindSpaceName(appMetadata.SpaceGuid)
		displayAppStats.OrgName = org.FindOrgNameBySpaceGuid(appMetadata.SpaceGuid)

		totalCpuPercentage := 0.0
		totalUsedMemory := uint64(0)
		totalUsedDisk := uint64(0)
		totalReportingContainers := 0

		if appMetadata.State == "STARTED" {
			displayAppStats.DesiredContainers = int(appMetadata.Instances)
		}

		stack := stack.FindStackMetadata(appMetadata.StackGuid)
		displayAppStats.StackId = appMetadata.StackGuid
		displayAppStats.StackName = stack.Name

		now := time.Now()
		for containerIndex, cs := range appStats.ContainerArray {
			if cs != nil && cs.ContainerMetric != nil {
				// If we haven't gotten a container update recently, ignore the old value
				if now.Sub(cs.LastUpdate) > time.Second*StaleContainerSeconds {
					appStats.ContainerArray[containerIndex] = nil
					continue
				}

				totalCpuPercentage = totalCpuPercentage + *cs.ContainerMetric.CpuPercentage
				totalUsedMemory = totalUsedMemory + *cs.ContainerMetric.MemoryBytes
				totalUsedDisk = totalUsedDisk + *cs.ContainerMetric.DiskBytes
				totalReportingContainers++
			}
		}
		if totalReportingContainers < displayAppStats.DesiredContainers {
			appsNotInDesiredState = appsNotInDesiredState + 1
		}

		displayAppStats.TotalCpuPercentage = totalCpuPercentage
		displayAppStats.TotalUsedMemory = totalUsedMemory
		displayAppStats.TotalUsedDisk = totalUsedDisk
		displayAppStats.TotalReportingContainers = totalReportingContainers

		logStdoutCount := int64(0)
		logStderrCount := int64(0)
		for _, cs := range appStats.ContainerArray {
			if cs != nil {
				logStdoutCount = logStdoutCount + cs.OutCount
				logStderrCount = logStderrCount + cs.ErrCount
			}
		}
		displayAppStats.TotalLogStdout = logStdoutCount + appStats.NonContainerStdout
		displayAppStats.TotalLogStderr = logStderrCount + appStats.NonContainerStderr

	}
	asUI.displayAppStats = displayStatsArray
	asUI.isWarmupComplete = asUI.GetMasterUI().IsWarmupComplete()
	asUI.appsNotInDesiredState = appsNotInDesiredState
	return displayStatsArray
}

func (asUI *AppListView) convertToListData(statsArray []*DisplayAppStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, len(statsArray))
	for i, d := range statsArray {
		listData[i] = d
	}
	return listData
}

func (asUI *AppListView) detailViewClosed(g *gocui.Gui) error {
	asUI.DataListView.RefreshDisplayCallback = nil
	return asUI.RefreshDisplay(g)
}

func (asUI *AppListView) preRowDisplay(data uiCommon.IData, isSelected bool) string {
	appStats := data.(*DisplayAppStats)
	colorString := ""
	if asUI.isWarmupComplete && appStats.DesiredContainers > appStats.TotalReportingContainers {
		if isSelected {
			colorString = util.RED_TEXT_GREEN_BG
		} else {
			colorString = util.BRIGHT_RED
		}
	} else if !isSelected && appStats.TotalTraffic.EventL10Rate > 0 {
		colorString = util.BRIGHT_WHITE
	}
	return colorString
}

func (asUI *AppListView) checkForAlerts(g *gocui.Gui) error {

	// TODO: We can't alert if we are not monitoring all the apps
	// Update his this alert only on monitored apps if non-privileged
	// for now we just don't alert
	if !asUI.GetMasterUI().IsPrivileged() {
		return nil
	}

	if asUI.isWarmupComplete && asUI.appsNotInDesiredState > 0 {
		plural := ""
		if asUI.appsNotInDesiredState > 1 {
			plural = "s"
		}
		msg := fmt.Sprintf("ALERT: %v application%v not in desired state (row%v colored red) ",
			asUI.appsNotInDesiredState, plural, plural)
		asUI.showUserMessage(g, msg)
	} else if asUI.isUserMessageOpen(g) {
		asUI.clearUserMessage(g)
	}
	return nil
}
