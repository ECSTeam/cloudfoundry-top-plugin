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

package spaceView

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
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appViews/appView"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

type SpaceListView struct {
	*dataView.DataListView
	displayAppStatsMap map[string]*dataCommon.DisplayAppStats
	orgId              string
}

func NewSpaceListView(masterUI masterUIInterface.MasterUIInterface,
	parentView dataView.DataListViewInterface,
	name string, bottomMargin int,
	eventProcessor *eventdata.EventProcessor,
	orgId string) *SpaceListView {

	asUI := &SpaceListView{orgId: orgId}

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("CPU", true),
		uiCommon.NewSortColumn("spaceName", false),
	}

	dataListView := dataView.NewDataListView(masterUI, parentView,
		name, 0, bottomMargin,
		eventProcessor, asUI, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.GetListData = asUI.GetListData

	org := org.FindOrgMetadata(orgId)
	dataListView.SetTitle(fmt.Sprintf("Space List of Org %v", org.Name))
	dataListView.HelpText = HelpText
	dataListView.HelpTextTips = HelpTextTips

	asUI.DataListView = dataListView

	return asUI

}

func (asUI *SpaceListView) initializeCallback(g *gocui.Gui, viewName string) error {

	if err := g.SetKeybinding(viewName, 'c', gocui.ModNone, asUI.copyAction); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(viewName, 'x', gocui.ModNone, asUI.CloseDetailView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, gocui.KeyEsc, gocui.ModNone, asUI.CloseDetailView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone, asUI.enterAction); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (asUI *SpaceListView) enterAction(g *gocui.Gui, v *gocui.View) error {
	highlightKey := asUI.GetListWidget().HighlightKey()
	if asUI.GetListWidget().HighlightKey() != "" {
		_, bottomMargin := asUI.GetMargins()

		// TODO: This should be changed to space view
		detailView := appView.NewAppListView(asUI.GetMasterUI(), asUI, "appBySpaceView",
			bottomMargin,
			asUI.GetEventProcessor(),
			highlightKey)

		asUI.SetDetailView(detailView)
		asUI.GetMasterUI().OpenView(g, detailView)
	}
	return nil
}

func (asUI *SpaceListView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, columnSpaceName())

	columns = append(columns, columnQuotaName())

	columns = append(columns, columnNumberOfApps())
	columns = append(columns, columnReportingContainers())

	columns = append(columns, columnTotalCpu())

	columns = append(columns, columnMemoryLimit())
	columns = append(columns, columnTotalReservedMemory())
	columns = append(columns, columnTotalReservedMemoryPercentOfSpaceQuota())
	columns = append(columns, columnTotalReservedMemoryPercentOfOrgQuota())
	columns = append(columns, columnTotalMemoryUsed())

	columns = append(columns, columnTotalReservedDisk())
	columns = append(columns, columnTotalDiskUsed())

	columns = append(columns, columnLogStdout())
	columns = append(columns, columnLogStderr())

	columns = append(columns, columnTotalReq())

	return columns
}

func (asUI *SpaceListView) copyAction(g *gocui.Gui, v *gocui.View) error {

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

func (asUI *SpaceListView) clipboardCallback(g *gocui.Gui, v *gocui.View, menuId string) error {

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

func (asUI *SpaceListView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *SpaceListView) postProcessData() map[string]*DisplaySpace {

	orgId := asUI.orgId
	spaceQuotaMdMgr := asUI.GetEventProcessor().GetMetadataManager().GetSpaceQuotaMdManager()
	appMdMgr := asUI.GetEventProcessor().GetMetadataManager().GetAppMdManager()

	// Get list of space in selected org
	allSpaces := space.All()
	orgSpaces := make([]space.Space, 0)
	for _, spaceMetadata := range allSpaces {
		if spaceMetadata.OrgGuid == orgId {
			orgSpaces = append(orgSpaces, spaceMetadata)
		}
	}

	// Build map of all apps by Space
	displayStatsMap := asUI.GetMasterUI().GetCommonData().GetDisplayAppStatsMap()
	appsBySpaceMap := make(map[string][]*dataCommon.DisplayAppStats)
	for _, appStats := range displayStatsMap {
		if appStats.OrgId == orgId {
			appsList := appsBySpaceMap[appStats.SpaceId]
			if appsList == nil {
				appsList = make([]*dataCommon.DisplayAppStats, 0)
			}
			appsList = append(appsList, appStats)
			appsBySpaceMap[appStats.SpaceId] = appsList
		}
	}

	// Build Map of Spaces
	displaySpaceMap := make(map[string]*DisplaySpace)
	for _, spaceMetadata := range orgSpaces {
		aSpaceMetadata := spaceMetadata
		displaySpace := NewDisplaySpace(&aSpaceMetadata)
		displaySpaceMap[spaceMetadata.Guid] = displaySpace
		displaySpace.NumberOfApps = len(appsBySpaceMap[spaceMetadata.Guid])

		if spaceMetadata.QuotaGuid != "" {
			spaceQuotaMd := spaceQuotaMdMgr.Find(spaceMetadata.QuotaGuid)
			displaySpace.QuotaName = spaceQuotaMd.Name
			displaySpace.MemoryLimitInBytes = int64(spaceQuotaMd.MemoryLimit) * util.MEGABYTE
		} else {
			displaySpace.QuotaName = "-none-"
		}

		for _, appStats := range appsBySpaceMap[spaceMetadata.Guid] {
			displaySpace.TotalCpuPercentage += appStats.TotalCpuPercentage
			displaySpace.TotalUsedMemory += appStats.TotalUsedMemory
			displaySpace.TotalUsedDisk += appStats.TotalUsedDisk
			displaySpace.TotalReportingContainers += appStats.TotalReportingContainers

			appMetadata := appMdMgr.FindAppMetadata(appStats.AppId)
			displaySpace.TotalReservedMemory += (int64(appMetadata.MemoryMB) * util.MEGABYTE) * int64(appMetadata.Instances)
			displaySpace.TotalReservedDisk += (int64(appMetadata.DiskQuotaMB) * util.MEGABYTE) * int64(appMetadata.Instances)

			if appStats.TotalTraffic != nil {
				displaySpace.HttpAllCount += appStats.TotalTraffic.HttpAllCount
			}

			for _, cs := range appStats.ContainerArray {
				if cs != nil {
					displaySpace.TotalLogStdout += cs.OutCount
					displaySpace.TotalLogStderr += cs.ErrCount
				}
			}
		}
		if displaySpace.MemoryLimitInBytes > 0 {
			displaySpace.TotalReservedMemoryPercentOfSpaceQuota = (float64(displaySpace.TotalReservedMemory) / float64(displaySpace.MemoryLimitInBytes)) * 100
		}

		org := org.FindOrgMetadata(orgId)
		orgQuotaMdMgr := asUI.GetEventProcessor().GetMetadataManager().GetOrgQuotaMdManager()
		orgQuotaMd := orgQuotaMdMgr.Find(org.QuotaGuid)
		orgQuotaMemoryLimitInBytes := orgQuotaMd.MemoryLimit * util.MEGABYTE
		if orgQuotaMemoryLimitInBytes > 0 {
			displaySpace.TotalReservedMemoryPercentOfOrgQuota = (float64(displaySpace.TotalReservedMemory) / float64(orgQuotaMemoryLimitInBytes)) * 100
		}

	}
	return displaySpaceMap
}

func (asUI *SpaceListView) convertToListData(statsMap map[string]*DisplaySpace) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(statsMap))
	for _, d := range statsMap {
		listData = append(listData, d)
	}
	return listData
}
