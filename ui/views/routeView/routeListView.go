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

package routeView

import (
	"fmt"
	"log"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/jroimartin/gocui"
)

type RouteListView struct {
	*dataView.DataListView
}

func NewRouteListView(masterUI masterUIInterface.MasterUIInterface,
	name string, bottomMargin int,
	eventProcessor *eventdata.EventProcessor) *RouteListView {

	asUI := &RouteListView{}

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("TOT-REQ", true),
		uiCommon.NewSortColumn("DOMAIN", false),
		uiCommon.NewSortColumn("HOST", false),
		uiCommon.NewSortColumn("PATH", false),
		//uiCommon.NewSortColumn("CELL_IP", false),
	}

	dataListView := dataView.NewDataListView(masterUI, nil,
		name, 0, bottomMargin,
		eventProcessor, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.UpdateHeaderCallback = asUI.updateHeader
	dataListView.GetListData = asUI.GetListData

	dataListView.SetTitle("Route List")
	dataListView.HelpText = helpText
	dataListView.HelpTextTips = appView.HelpTextTips

	asUI.DataListView = dataListView

	return asUI

}

func (asUI *RouteListView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, columnHost())
	columns = append(columns, columnDomain())
	columns = append(columns, columnPath())

	columns = append(columns, columnRoutedAppCount())

	columns = append(columns, columnTotalRequests())
	columns = append(columns, column2xx())
	columns = append(columns, column3xx())
	columns = append(columns, column4xx())
	columns = append(columns, column5xx())

	columns = append(columns, columnResponseContentLength())

	columns = append(columns, columnMethodGet())
	columns = append(columns, columnMethodPost())
	columns = append(columns, columnMethodPut())
	columns = append(columns, columnMethodDelete())
	columns = append(columns, columnLastAccess())

	return columns
}

func (asUI *RouteListView) initializeCallback(g *gocui.Gui, viewName string) error {

	if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone, asUI.enterAction); err != nil {
		log.Panicln(err)
	}

	return nil
}

func (asUI *RouteListView) enterAction(g *gocui.Gui, v *gocui.View) error {
	/*
		highlightKey := asUI.GetListWidget().HighlightKey()
		if asUI.GetListWidget().HighlightKey() != "" {
			topMargin, bottomMargin := asUI.GetMargins()

			detailView := routeListView.NewCellDetailView(asUI.GetMasterUI(), asUI,
				"routeListView",
				topMargin, bottomMargin,
				asUI.GetEventProcessor(),
				highlightKey)

			asUI.SetDetailView(detailView)
			asUI.GetMasterUI().OpenView(g, detailView)
		}
	*/
	return nil
}

func (asUI *RouteListView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *RouteListView) postProcessData() []*displaydata.DisplayRouteStats {

	domainMap := asUI.GetDisplayedEventData().DomainMap
	displayRouteArray := make([]*displaydata.DisplayRouteStats, 0)

	for domainName, domainStats := range domainMap {
		for hostName, hostStats := range domainStats.HostStatsMap {
			for pathName, routeStats := range hostStats.RouteStatsMap {

				displayRouteStat := displaydata.NewDisplayRouteStats(routeStats)
				displayRouteArray = append(displayRouteArray, displayRouteStat)
				displayRouteStat.RouteName = fmt.Sprintf("%v.%v%v", hostName, domainName, pathName)
				displayRouteStat.Host = hostName
				displayRouteStat.Domain = domainName
				displayRouteStat.Path = pathName
				displayRouteStat.RoutedAppCount = len(routeStats.AppRouteStatsMap)

				for _, appRouteStats := range routeStats.AppRouteStatsMap {

					displayRouteStat.ResponseContentLength = displayRouteStat.ResponseContentLength + appRouteStats.ResponseContentLength
					if displayRouteStat.LastAccess.Before(appRouteStats.LastAccess) {
						displayRouteStat.LastAccess = appRouteStats.LastAccess
					}

					displayRouteStat.HttpMethodGetCount = appRouteStats.HttpMethod[events.Method_GET]
					displayRouteStat.HttpMethodPostCount = appRouteStats.HttpMethod[events.Method_POST]
					displayRouteStat.HttpMethodPutCount = appRouteStats.HttpMethod[events.Method_PUT]
					displayRouteStat.HttpMethodDeleteCount = appRouteStats.HttpMethod[events.Method_DELETE]
					//displayRouteStat.HttpMethodOtherCount = ???

					for statusCode, responseCount := range appRouteStats.HttpStatusCode {
						displayRouteStat.HttpAllCount = displayRouteStat.HttpAllCount + responseCount
						switch {
						case statusCode >= 200 && statusCode < 300:
							displayRouteStat.Http2xxCount = responseCount
						case statusCode >= 300 && statusCode < 400:
							displayRouteStat.Http3xxCount = responseCount
						case statusCode >= 400 && statusCode < 500:
							displayRouteStat.Http4xxCount = responseCount
						case statusCode >= 500 && statusCode < 600:
							displayRouteStat.Http5xxCount = responseCount
						default:
							displayRouteStat.HttpOtherCount = responseCount
						}

					}
				}
			}
		}

	}

	return displayRouteArray
}

func (asUI *RouteListView) convertToListData(displayRouteArray []*displaydata.DisplayRouteStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(displayRouteArray))
	for _, d := range displayRouteArray {
		listData = append(listData, d)
	}
	return listData
}

func (asUI *RouteListView) PreRowDisplay(data uiCommon.IData, isSelected bool) string {
	return ""
}

func (asUI *RouteListView) updateHeader(g *gocui.Gui, v *gocui.View) (int, error) {
	fmt.Fprintf(v, "\nTODO: Show summary stats")
	return 3, nil
}
