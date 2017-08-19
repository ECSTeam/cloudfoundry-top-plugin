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

package orgView

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
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appViews/appView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/orgSpaceViews/spaceView"
)

type OrgListView struct {
	*dataView.DataListView
	displayAppStatsMap map[string]*dataCommon.DisplayAppStats
	isWarmupComplete   bool
}

func NewOrgListView(masterUI masterUIInterface.MasterUIInterface,
	name string, bottomMargin int,
	eventProcessor *eventdata.EventProcessor) *OrgListView {

	asUI := &OrgListView{}

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("CPU_PER", true),
		uiCommon.NewSortColumn("ORG", false),
	}

	dataListView := dataView.NewDataListView(masterUI, nil,
		name, 0, bottomMargin,
		eventProcessor, asUI, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.GetListData = asUI.GetListData

	dataListView.SetTitle("Org List")
	dataListView.HelpText = HelpText
	dataListView.HelpTextTips = appView.HelpTextTips

	asUI.DataListView = dataListView

	return asUI

}

func (asUI *OrgListView) initializeCallback(g *gocui.Gui, viewName string) error {

	if err := g.SetKeybinding(viewName, 'c', gocui.ModNone, asUI.copyAction); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone, asUI.enterAction); err != nil {
		log.Panicln(err)
	}

	return nil
}

func (asUI *OrgListView) enterAction(g *gocui.Gui, v *gocui.View) error {
	highlightKey := asUI.GetListWidget().HighlightKey()
	if asUI.GetListWidget().HighlightKey() != "" {
		_, bottomMargin := asUI.GetMargins()

		view := spaceView.NewSpaceListView(asUI.GetMasterUI(), asUI, "spaceListView",
			bottomMargin,
			asUI.GetEventProcessor(),
			highlightKey)
		asUI.SetDetailView(view)
		asUI.GetMasterUI().OpenView(g, view)
	}
	return nil
}

func (asUI *OrgListView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, columnOrgName())
	columns = append(columns, columnStatus())

	columns = append(columns, columnQuotaName())

	columns = append(columns, columnNumberOfSpaces())
	columns = append(columns, columnNumberOfApps())

	columns = append(columns, columnDesiredContainers())
	columns = append(columns, columnReportingContainers())

	columns = append(columns, columnTotalCpu())

	columns = append(columns, columnMemoryLimit())
	columns = append(columns, columnTotalMemoryReserved())
	columns = append(columns, columnTotalMemoryReservedPercentOfQuota())
	columns = append(columns, columnTotalMemoryUsed())

	columns = append(columns, columnTotalDiskReserved())
	columns = append(columns, columnTotalDiskUsed())

	columns = append(columns, columnLogStdout())
	columns = append(columns, columnLogStderr())

	columns = append(columns, columnTotalReq())

	return columns
}

func (asUI *OrgListView) copyAction(g *gocui.Gui, v *gocui.View) error {

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

func (asUI *OrgListView) clipboardCallback(g *gocui.Gui, v *gocui.View, menuId string) error {

	clipboardValue := ""

	selectedAppId := asUI.GetListWidget().HighlightKey()
	statsMap := asUI.GetDisplayedEventData().AppMap
	appStats := statsMap[selectedAppId]
	if appStats == nil {
		// Nothing selected
		return nil
	}
	appMetadata := asUI.GetAppMdMgr().FindItem(selectedAppId)
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

func (asUI *OrgListView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *OrgListView) postProcessData() map[string]*DisplayOrg {

	orgQuotaMdMgr := asUI.GetEventProcessor().GetMetadataManager().GetOrgQuotaMdManager()
	appMdMgr := asUI.GetEventProcessor().GetMetadataManager().GetAppMdManager()

	// Build map of all spaces by Org
	spaces := space.All()
	spaceByOrgMap := make(map[string][]space.Space)
	for _, spaceMetadata := range spaces {
		spacesList := spaceByOrgMap[spaceMetadata.OrgGuid]
		if spacesList == nil {
			spacesList = make([]space.Space, 0)
		}
		spacesList = append(spacesList, spaceMetadata)
		spaceByOrgMap[spaceMetadata.OrgGuid] = spacesList
	}

	// Build map of all apps by Org
	displayStatsMap := asUI.GetMasterUI().GetCommonData().GetDisplayAppStatsMap()
	appsByOrgMap := make(map[string][]*dataCommon.DisplayAppStats)
	for _, appStats := range displayStatsMap {
		appsList := appsByOrgMap[appStats.OrgId]
		if appsList == nil {
			appsList = make([]*dataCommon.DisplayAppStats, 0)
		}
		appsList = append(appsList, appStats)
		appsByOrgMap[appStats.OrgId] = appsList
	}

	// Build Map of Orgs
	orgs := org.All()
	displayOrgMap := make(map[string]*DisplayOrg)
	for _, org := range orgs {
		anOrg := org
		displayOrg := NewDisplayOrg(&anOrg)
		displayOrgMap[org.Guid] = displayOrg
		displayOrg.NumberOfSpaces = len(spaceByOrgMap[org.Guid])
		displayOrg.NumberOfApps = len(appsByOrgMap[org.Guid])

		orgQuotaMd := orgQuotaMdMgr.Find(org.QuotaGuid)
		displayOrg.QuotaName = orgQuotaMd.Name
		displayOrg.MemoryLimitInBytes = int64(orgQuotaMd.MemoryLimit) * util.MEGABYTE

		for _, appStats := range appsByOrgMap[org.Guid] {
			displayOrg.TotalCpuPercentage += appStats.TotalCpuPercentage
			displayOrg.TotalMemoryUsed += appStats.TotalMemoryUsed
			displayOrg.TotalDiskUsed += appStats.TotalDiskUsed
			displayOrg.TotalReportingContainers += appStats.TotalReportingContainers
			displayOrg.DesiredContainers += appStats.DesiredContainers

			appMetadata := appMdMgr.FindItem(appStats.AppId)
			displayOrg.TotalMemoryReserved += (int64(appMetadata.MemoryMB) * util.MEGABYTE) * int64(appMetadata.Instances)
			displayOrg.TotalDiskReserved += (int64(appMetadata.DiskQuotaMB) * util.MEGABYTE) * int64(appMetadata.Instances)

			if appStats.TotalTraffic != nil {
				displayOrg.HttpAllCount += appStats.HttpAllCount
			}

			for _, cs := range appStats.ContainerArray {
				if cs != nil {
					displayOrg.TotalLogStdout += cs.OutCount
					displayOrg.TotalLogStderr += cs.ErrCount
				}
			}

		}

		if displayOrg.MemoryLimitInBytes > 0 {
			displayOrg.TotalMemoryReservedPercentOfQuota = (float64(displayOrg.TotalMemoryReserved) / float64(displayOrg.MemoryLimitInBytes)) * 100
		}

	}
	asUI.isWarmupComplete = asUI.GetMasterUI().IsWarmupComplete()
	return displayOrgMap
}

func (asUI *OrgListView) convertToListData(statsMap map[string]*DisplayOrg) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(statsMap))
	for _, d := range statsMap {
		listData = append(listData, d)
	}
	return listData
}
