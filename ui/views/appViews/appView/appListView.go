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

	"github.com/atotto/clipboard"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/org"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/space"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/dataCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/dataView"
	"github.com/jroimartin/gocui"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appViews/appDetailView"
)

type AppListView struct {
	*dataView.DataListView
	displayAppStatsMap map[string]*dataCommon.DisplayAppStats
	isWarmupComplete   bool
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
		eventProcessor, asUI, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	//dataListView.SetColumnOwner(asUI)

	//dataListView.PreRowDisplayCallback = asUI.preRowDisplay
	// TODO: Add additional header rows such as "active apps"
	//dataListView.UpdateHeaderCallback = asUI.updateHeader
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
	columns = append(columns, columnAppName())
	columns = append(columns, columnSpaceName())
	columns = append(columns, columnOrgName())

	columns = append(columns, columnDesiredInstances())
	columns = append(columns, columnReportingContainers())

	columns = append(columns, columnTotalCpu())
	columns = append(columns, columnTotalMemoryUsed())
	columns = append(columns, columnTotalDiskUsed())

	columns = append(columns, columnAvgResponseTimeL60Info())
	columns = append(columns, columnLogStdout())
	columns = append(columns, columnLogStderr())

	columns = append(columns, columnReq1())
	columns = append(columns, columnReq10())
	columns = append(columns, columnReq60())

	columns = append(columns, columnTotalReq())
	columns = append(columns, column2XX())
	columns = append(columns, column3XX())
	columns = append(columns, column4XX())
	columns = append(columns, column5XX())

	columns = append(columns, columnStackName())

	return columns
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

func (asUI *AppListView) postProcessData() map[string]*dataCommon.DisplayAppStats {

	displayStatsMap := asUI.GetMasterUI().GetCommonData().GetDisplayAppStatsMap()

	for appId, appStats := range displayStatsMap {
		logStdoutCount := int64(0)
		logStderrCount := int64(0)
		for _, cs := range appStats.ContainerArray {
			if cs != nil {
				logStdoutCount = logStdoutCount + cs.OutCount
				logStderrCount = logStderrCount + cs.ErrCount
			}
		}
		displayAppStats := displayStatsMap[appId]
		displayAppStats.TotalLogStdout = logStdoutCount + appStats.NonContainerStdout
		displayAppStats.TotalLogStderr = logStderrCount + appStats.NonContainerStderr
	}

	asUI.displayAppStatsMap = displayStatsMap
	asUI.isWarmupComplete = asUI.GetMasterUI().IsWarmupComplete()
	return displayStatsMap
}

func (asUI *AppListView) convertToListData(statsMap map[string]*dataCommon.DisplayAppStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(statsMap))
	for _, d := range statsMap {
		listData = append(listData, d)
	}
	return listData
}

func (asUI *AppListView) detailViewClosed(g *gocui.Gui) error {
	asUI.DataListView.RefreshDisplayCallback = nil
	return asUI.RefreshDisplay(g)
}

/*
func (asUI *AppListView) preRowDisplay(data uiCommon.IData, isSelected bool) string {
	appStats := data.(*dataCommon.DisplayAppStats)
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
*/
