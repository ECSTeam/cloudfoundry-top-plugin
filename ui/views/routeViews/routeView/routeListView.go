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
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appViews/appView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/routeViews/routeMapView"
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
		uiCommon.NewSortColumn("TOTREQ", true),
		uiCommon.NewSortColumn("DOMAIN", false),
		uiCommon.NewSortColumn("HOST", false),
		uiCommon.NewSortColumn("PATH", false),
		uiCommon.NewSortColumn("PORT", false),
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
	columns = append(columns, columnPort())

	columns = append(columns, columnRoutedAppCount())

	columns = append(columns, columnTotalRequests())
	columns = append(columns, column2xx())
	columns = append(columns, column3xx())
	columns = append(columns, column4xx())
	columns = append(columns, column5xx())

	columns = append(columns, columnResponseContentLength())

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
	highlightKey := asUI.GetListWidget().HighlightKey()
	if asUI.GetListWidget().HighlightKey() != "" {
		topMargin, bottomMargin := asUI.GetMargins()

		detailView := routeMapView.NewRouteMapListView(asUI.GetMasterUI(), asUI,
			"routeMapListView",
			topMargin, bottomMargin,
			asUI.GetEventProcessor(),
			highlightKey)

		asUI.SetDetailView(detailView)
		asUI.GetMasterUI().OpenView(g, detailView)
	}

	return nil
}

func (asUI *RouteListView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *RouteListView) postProcessData() []*DisplayRouteStats {

	domainMap := asUI.GetDisplayedEventData().DomainMap
	displayRouteArray := make([]*DisplayRouteStats, 0)

	for domainName, domainStats := range domainMap {
		for hostName, hostStats := range domainStats.HostStatsMap {
			// HTTP routes
			for pathName, routeStats := range hostStats.RouteStatsMap {

				displayRouteStat := NewDisplayRouteStats(routeStats, hostName, domainName, pathName, 0)
				displayRouteArray = append(displayRouteArray, displayRouteStat)

				for appId, appRouteStats := range routeStats.AppRouteStatsMap {

					if appId != "" {
						displayRouteStat.RoutedAppCount = displayRouteStat.RoutedAppCount + 1
					}

					for method, httpMethodStats := range appRouteStats.HttpMethodStatsMap {

						displayRouteStat.ResponseContentLength = displayRouteStat.ResponseContentLength + httpMethodStats.ResponseContentLength
						if displayRouteStat.LastAccess.Before(httpMethodStats.LastAccess) {
							displayRouteStat.LastAccess = httpMethodStats.LastAccess
						}

						switch method {
						case events.Method_GET:
							displayRouteStat.HttpMethodGetCount = httpMethodStats.RequestCount
						case events.Method_POST:
							displayRouteStat.HttpMethodPostCount = httpMethodStats.RequestCount
						case events.Method_PUT:
							displayRouteStat.HttpMethodPutCount = httpMethodStats.RequestCount
						case events.Method_DELETE:
							displayRouteStat.HttpMethodDeleteCount = httpMethodStats.RequestCount
						}

						for statusCode, responseCount := range httpMethodStats.HttpStatusCode {
							displayRouteStat.HttpAllCount = displayRouteStat.HttpAllCount + responseCount
							switch {
							case statusCode >= 200 && statusCode < 300:
								displayRouteStat.Http2xxCount = displayRouteStat.Http2xxCount + responseCount
							case statusCode >= 300 && statusCode < 400:
								displayRouteStat.Http3xxCount = displayRouteStat.Http3xxCount + responseCount
							case statusCode >= 400 && statusCode < 500:
								displayRouteStat.Http4xxCount = displayRouteStat.Http4xxCount + responseCount
							case statusCode >= 500 && statusCode < 600:
								displayRouteStat.Http5xxCount = displayRouteStat.Http5xxCount + responseCount
							default:
								displayRouteStat.HttpOtherCount = displayRouteStat.HttpOtherCount + responseCount
							}

						}
					}
				}
			}

			// TCP routes
			for port, routeStats := range hostStats.TcpRouteStatsMap {
				displayRouteStat := NewDisplayRouteStats(routeStats, hostName, domainName, "", port)
				displayRouteArray = append(displayRouteArray, displayRouteStat)
				//for appId, appRouteStats := range routeStats.AppRouteStatsMap {
				for appId, _ := range routeStats.AppRouteStatsMap {
					if appId != "" {
						displayRouteStat.RoutedAppCount = displayRouteStat.RoutedAppCount + 1
					}
				}
			}

		}

	}

	return displayRouteArray
}

func (asUI *RouteListView) convertToListData(displayRouteArray []*DisplayRouteStats) []uiCommon.IData {
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
