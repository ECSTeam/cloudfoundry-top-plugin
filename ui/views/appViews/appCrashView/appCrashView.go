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

package appCrashView

import (
	"fmt"
	"log"
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventApp"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/app"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/crashData"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

type AppCrashView struct {
	*dataView.DataListView
	appId         string
	displayMenuId string
	appMdMgr      *app.AppMetadataManager
}

func NewAppCrashView(masterUI masterUIInterface.MasterUIInterface,
	parentView dataView.DataListViewInterface,
	name string, bottomMargin int,
	eventProcessor *eventdata.EventProcessor,
	appId string) *AppCrashView {

	appMdMgr := eventProcessor.GetMetadataManager().GetAppMdManager()

	asUI := &AppCrashView{appId: appId, appMdMgr: appMdMgr}
	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("CRASH_TIME", true),
		uiCommon.NewSortColumn("IDX", false),
	}

	dataListView := dataView.NewDataListView(masterUI, parentView,
		name, 0, bottomMargin,
		eventProcessor, asUI, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.GetListData = asUI.GetListData

	titleFunc := func() string {
		return fmt.Sprintf("App: %v - Container CRASH List (last 24 hours)", asUI.getAppName())
	}
	//dataListView.SetTitle(fmt.Sprintf("App: %v - Container CRASH List (last 24 hours)", asUI.getAppName()))
	dataListView.SetTitle(titleFunc)
	dataListView.RefreshDisplayCallback = asUI.refreshDisplay

	dataListView.HelpText = HelpText
	dataListView.HelpTextTips = HelpTextTips

	asUI.DataListView = dataListView

	return asUI
}

func (asUI *AppCrashView) initializeCallback(g *gocui.Gui, viewName string) error {
	if err := g.SetKeybinding(viewName, 'x', gocui.ModNone, asUI.closeAppCrashView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, gocui.KeyEsc, gocui.ModNone, asUI.closeAppCrashView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone, asUI.enterAction); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (asUI *AppCrashView) enterAction(g *gocui.Gui, v *gocui.View) error {

	highlightKey := asUI.GetListWidget().HighlightKey()
	if highlightKey != "" {
		widgetName := "crashItemView"
		idata := asUI.GetListWidget().HighlightData()
		if idata != nil {
			crashInfo := idata.(*DisplayContainerCrashInfo)
			view := NewAppCrashItemWidget(asUI.GetMasterUI(), widgetName, 70, 18, asUI, crashInfo)
			return asUI.GetMasterUI().OpenView(g, view)
		}
	}
	return nil
}

func (asUI *AppCrashView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, ColumnCrashTime())
	columns = append(columns, ColumnContainerIndex())
	columns = append(columns, ColumnExitDescription())
	return columns
}

func (asUI *AppCrashView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *AppCrashView) postProcessData() []*DisplayContainerCrashInfo {

	displayCrashInfoList := make([]*DisplayContainerCrashInfo, 0)

	appMap := asUI.GetDisplayedEventData().AppMap
	appId := asUI.appId
	appStats := appMap[appId]
	if appStats == nil {
		return displayCrashInfoList
	}

	crashInfoListFromMetadata := crashData.FindSinceByApp(appStats.AppId, -24*time.Hour)
	displayCrashInfoList = append(displayCrashInfoList, asUI.createDisplayContainerCrashInfo(crashInfoListFromMetadata)...)

	crashInfoListFromLiveCapture := appStats.CrashSince(-24 * time.Hour)
	displayCrashInfoList = append(displayCrashInfoList, asUI.createDisplayContainerCrashInfo(crashInfoListFromLiveCapture)...)

	return displayCrashInfoList
}

func (asUI *AppCrashView) createDisplayContainerCrashInfo(crashInfoList []*crashData.ContainerCrashInfo) []*DisplayContainerCrashInfo {

	displayCrashInfoList := make([]*DisplayContainerCrashInfo, 0)
	for _, crashInfo := range crashInfoList {
		displayCrashInfo := NewDisplayContainerCrashInfo(crashInfo)
		displayCrashInfo.CrashTimeFormatted = crashInfo.CrashTime.Local().Format("01-02-2006 15:04:05")
		displayCrashInfoList = append(displayCrashInfoList, displayCrashInfo)
	}
	return displayCrashInfoList
}

func (asUI *AppCrashView) FindLastCrash(appStats *eventApp.AppStats) *crashData.ContainerCrashInfo {
	if appStats.ContainerCrashInfo != nil && len(appStats.ContainerCrashInfo) > 0 {
		last := len(appStats.ContainerCrashInfo) - 1
		return appStats.ContainerCrashInfo[last]
	}
	return nil
}

func (asUI *AppCrashView) convertToListData(containerStatsArray []*DisplayContainerCrashInfo) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(containerStatsArray))
	for _, d := range containerStatsArray {
		listData = append(listData, d)
	}
	return listData
}

func (asUI *AppCrashView) closeAppCrashView(g *gocui.Gui, v *gocui.View) error {
	if err := asUI.GetMasterUI().CloseView(asUI); err != nil {
		return err
	}
	return nil
}

func (asUI *AppCrashView) getAppName() string {
	appMetadata := asUI.appMdMgr.FindItem(asUI.appId)
	appName := appMetadata.Name
	return appName
}

func (asUI *AppCrashView) refreshDisplay(g *gocui.Gui) (propagateRefresh bool, err error) {
	v, err := g.View(asUI.DataListView.Name())
	if err != nil {
		return false, err
	}
	appId := asUI.appId
	mdAppMgr := asUI.appMdMgr
	if mdAppMgr.IsPendingDeleteFromCache(appId) || mdAppMgr.IsDeletedFromCache(appId) {
		v.Clear()
		fmt.Fprintf(v, " \n")
		fmt.Fprintf(v, "%v", util.BRIGHT_RED)
		fmt.Fprintf(v, " Application has been deleted")
		fmt.Fprintf(v, "%v", util.CLEAR)
		return false, nil
	}
	return true, nil
}
