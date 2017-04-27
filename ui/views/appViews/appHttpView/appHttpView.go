// Copyright (c) 2017 ECS Team, Inc. - All Rights Reserved
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

package appHttpView

import (
	"fmt"
	"log"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventApp"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/app"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/crashData"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/dataView"
	"github.com/jroimartin/gocui"
)

type AppHttpView struct {
	*dataView.DataListView
	appId         string
	displayMenuId string
	appMdMgr      *app.AppMetadataManager
}

func NewAppHttpView(masterUI masterUIInterface.MasterUIInterface,
	parentView dataView.DataListViewInterface,
	name string, bottomMargin int,
	eventProcessor *eventdata.EventProcessor,
	appId string) *AppHttpView {

	appMdMgr := eventProcessor.GetMetadataManager().GetAppMdManager()

	asUI := &AppHttpView{appId: appId, appMdMgr: appMdMgr}
	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("METHOD", false),
		uiCommon.NewSortColumn("CODE", false),
	}

	dataListView := dataView.NewDataListView(masterUI, parentView,
		name, 0, bottomMargin,
		eventProcessor, asUI, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.GetListData = asUI.GetListData

	dataListView.SetTitle(fmt.Sprintf("App: %v - HTTP(S) Response Info", asUI.getAppName()))

	dataListView.HelpText = HelpText
	dataListView.HelpTextTips = HelpTextTips

	asUI.DataListView = dataListView

	return asUI
}

func (asUI *AppHttpView) initializeCallback(g *gocui.Gui, viewName string) error {
	if err := g.SetKeybinding(viewName, 'x', gocui.ModNone, asUI.closeAppHttpView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, gocui.KeyEsc, gocui.ModNone, asUI.closeAppHttpView); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (asUI *AppHttpView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, ColumnMethod())
	columns = append(columns, ColumnStatusCode())
	columns = append(columns, ColumnLastAcivity())
	columns = append(columns, ColumnCount())

	return columns
}

func (asUI *AppHttpView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *AppHttpView) postProcessData() []*DisplayHttpInfo {

	displayInfoList := make([]*DisplayHttpInfo, 0)

	appMap := asUI.GetDisplayedEventData().AppMap
	appStats := appMap[asUI.appId]
	if appStats == nil {
		return displayInfoList
	}

	// TODO: This is wrong -- need to deal with multiple containers worth of info not overwriting data
	for _, containerTraffic := range appStats.ContainerTrafficMap {
		for _, httpStatusCodeMap := range containerTraffic.HttpInfoMap {
			for _, httpInfo := range httpStatusCodeMap {
				if httpInfo != nil {
					displayHttpInfo := NewDisplayHttpInfo(httpInfo)
					displayInfoList = append(displayInfoList, displayHttpInfo)
					displayHttpInfo.LastAcivityFormatted = httpInfo.LastAcivity.Local().Format("01-02-2006 15:04:05")
				}
			}
		}
	}

	return displayInfoList
}

func (asUI *AppHttpView) FindLastCrash(appStats *eventApp.AppStats) *crashData.ContainerCrashInfo {
	if appStats.ContainerCrashInfo != nil && len(appStats.ContainerCrashInfo) > 0 {
		last := len(appStats.ContainerCrashInfo) - 1
		return appStats.ContainerCrashInfo[last]
	}
	return nil
}

func (asUI *AppHttpView) convertToListData(containerStatsArray []*DisplayHttpInfo) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(containerStatsArray))
	for _, d := range containerStatsArray {
		listData = append(listData, d)
	}
	return listData
}

func (asUI *AppHttpView) closeAppHttpView(g *gocui.Gui, v *gocui.View) error {
	if err := asUI.GetMasterUI().CloseView(asUI); err != nil {
		return err
	}
	return nil
}

func (asUI *AppHttpView) getAppName() string {
	appMetadata := asUI.appMdMgr.FindAppMetadata(asUI.appId)
	appName := appMetadata.Name
	return appName
}
