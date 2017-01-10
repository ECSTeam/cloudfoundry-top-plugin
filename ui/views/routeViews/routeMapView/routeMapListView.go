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

package routeMapView

import (
	"fmt"
	"log"
	"strings"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventRoute"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/domain"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/org"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/route"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/space"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/cellViews/cellDetailView"
	"github.com/jroimartin/gocui"
)

type RouteMapListView struct {
	*dataView.DataListView
	routeId string

	routeMapDetailWidget *RouteMapDetailWidget
}

func NewRouteMapListView(masterUI masterUIInterface.MasterUIInterface,
	parentView dataView.DataListViewInterface,
	name string, topMargin, bottomMargin int,
	eventProcessor *eventdata.EventProcessor, routeId string) *RouteMapListView {

	asUI := &RouteMapListView{}
	asUI.routeId = routeId
	detailWidgetViewHeight := 5

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("TOTREQ", true),
		uiCommon.NewSortColumn("appName", false),
	}

	dataListView := dataView.NewDataListView(masterUI, parentView,
		name, topMargin+detailWidgetViewHeight+1, bottomMargin,
		eventProcessor, asUI, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.UpdateHeaderCallback = asUI.updateHeader
	dataListView.GetListData = asUI.GetListData

	dataListView.SetTitle("Route Map List")
	dataListView.HelpText = helpText
	dataListView.HelpTextTips = cellDetailView.HelpTextTips

	asUI.DataListView = dataListView

	asUI.routeMapDetailWidget = NewRouteMapDetailWidget(masterUI, "routeMapDetailWidget", detailWidgetViewHeight, asUI)
	masterUI.LayoutManager().Add(asUI.routeMapDetailWidget)
	return asUI

}

func (asUI *RouteMapListView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, columnAppName())

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

func (asUI *RouteMapListView) initializeCallback(g *gocui.Gui, viewName string) error {

	if err := g.SetKeybinding(viewName, 'x', gocui.ModNone, asUI.closeAppDetailView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, gocui.KeyEsc, gocui.ModNone, asUI.closeAppDetailView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone, asUI.enterAction); err != nil {
		log.Panicln(err)
	}

	return nil
}

func (asUI *RouteMapListView) closeAppDetailView(g *gocui.Gui, v *gocui.View) error {
	if err := asUI.GetMasterUI().CloseView(asUI); err != nil {
		return err
	}
	if err := asUI.GetMasterUI().CloseView(asUI.routeMapDetailWidget); err != nil {
		return err
	}
	return nil
}

func (asUI *RouteMapListView) enterAction(g *gocui.Gui, v *gocui.View) error {
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

func (asUI *RouteMapListView) seedAppsWithNoTraffic(routeStats *eventRoute.RouteStats) {

	routeMd := route.FindRouteMetadata(asUI.routeId)
	if routeMd.InternalGenerated {
		return
	}

	appIds := route.FindAppIdsForRouteMetadata(asUI.GetEventProcessor().GetCliConnection(), asUI.routeId)
	for _, appId := range appIds {
		appRouteStats := routeStats.FindAppRouteStats(appId)
		if appRouteStats == nil {
			/*
				toplog.Info("seedAppsWithNoTraffic adding appId %v to routeId %v",
					appId, asUI.routeId)
			*/
			appRouteStats = eventRoute.NewAppRouteStats(appId)
			routeStats.AppRouteStatsMap[appId] = appRouteStats
		}
	}
}

func (asUI *RouteMapListView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *RouteMapListView) postProcessData() []*DisplayRouteMapStats {

	domainMap := asUI.GetDisplayedEventData().DomainMap
	displayRouteArray := make([]*DisplayRouteMapStats, 0)

	routeMd := route.FindRouteMetadata(asUI.routeId)
	domainMd := domain.FindDomainMetadata(routeMd.DomainGuid)
	domainName := strings.ToLower(domainMd.Name)
	hostName := strings.ToLower(routeMd.Host)
	pathName := routeMd.Path
	port := routeMd.Port

	domainStats := domainMap[domainName]
	hostStats := domainStats.HostStatsMap[hostName]

	if port == 0 {
		routeStats := hostStats.RouteStatsMap[pathName]

		asUI.seedAppsWithNoTraffic(routeStats)

		for appId, appRouteStats := range routeStats.AppRouteStatsMap {

			appMetadata := asUI.GetAppMdMgr().FindAppMetadata(appId)
			appName := appMetadata.Name
			spaceName := space.FindSpaceName(appMetadata.SpaceGuid)
			orgName := org.FindOrgNameBySpaceGuid(appMetadata.SpaceGuid)

			displayRouteStat := NewDisplayRouteMapStats(routeStats, appId, appName, spaceName, orgName)
			displayRouteArray = append(displayRouteArray, displayRouteStat)
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
	} else {
		// TCP routes
		if port != 0 {
			tcpRouteStats := hostStats.TcpRouteStatsMap[port]
			for appId, _ := range tcpRouteStats.AppRouteStatsMap {
				displayRouteStat := NewDisplayRouteMapStats(tcpRouteStats, appId, "", "", "")
				displayRouteArray = append(displayRouteArray, displayRouteStat)
			}
		}
	}

	return displayRouteArray
}

func (asUI *RouteMapListView) convertToListData(displayRouteArray []*DisplayRouteMapStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(displayRouteArray))
	for _, d := range displayRouteArray {
		listData = append(listData, d)
	}
	return listData
}

func (asUI *RouteMapListView) PreRowDisplay(data uiCommon.IData, isSelected bool) string {
	return ""
}

func (asUI *RouteMapListView) updateHeader(g *gocui.Gui, v *gocui.View) (int, error) {
	fmt.Fprintf(v, "\nTODO: Show summary stats")
	return 3, nil
}
