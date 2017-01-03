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

package eventView

import (
	"fmt"
	"log"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/eventViews/eventOriginView"
	"github.com/jroimartin/gocui"
)

type EventListView struct {
	*dataView.DataListView
}

func NewEventListView(masterUI masterUIInterface.MasterUIInterface,
	name string, bottomMargin int,
	eventProcessor *eventdata.EventProcessor) *EventListView {

	asUI := &EventListView{}

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("COUNT", true),
		uiCommon.NewSortColumn("EVENT_TYPE", true),
	}

	dataListView := dataView.NewDataListView(masterUI, nil,
		name, 0, bottomMargin,
		eventProcessor, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.UpdateHeaderCallback = asUI.updateHeader
	dataListView.GetListData = asUI.GetListData

	dataListView.SetTitle("Event List")
	dataListView.HelpText = helpText
	dataListView.HelpTextTips = appView.HelpTextTips

	asUI.DataListView = dataListView

	return asUI

}

func (asUI *EventListView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)

	columns = append(columns, columnEventType())
	columns = append(columns, columnEventCount())

	return columns
}

func (asUI *EventListView) initializeCallback(g *gocui.Gui, viewName string) error {

	if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone, asUI.enterAction); err != nil {
		log.Panicln(err)
	}

	return nil
}

func (asUI *EventListView) enterAction(g *gocui.Gui, v *gocui.View) error {

	highlightKey := asUI.GetListWidget().HighlightKey()
	if asUI.GetListWidget().HighlightKey() != "" {
		topMargin, bottomMargin := asUI.GetMargins()

		detailView := eventOriginView.NewEventOriginListView(asUI.GetMasterUI(), asUI,
			"eventOriginView",
			topMargin, bottomMargin,
			asUI.GetEventProcessor(),
			highlightKey)

		asUI.SetDetailView(detailView)
		asUI.GetMasterUI().OpenView(g, detailView)
	}

	return nil
}

func (asUI *EventListView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *EventListView) postProcessData() []*displaydata.DisplayEventStats {

	eventTypeMap := asUI.GetDisplayedEventData().EventTypeMap
	displayEventList := make([]*displaydata.DisplayEventStats, 0, len(eventTypeMap))

	for _, eventTypeStats := range eventTypeMap {
		displayEventStat := displaydata.NewDisplayEventStats(eventTypeStats)
		displayEventList = append(displayEventList, displayEventStat)
		for _, eventOriginStats := range eventTypeStats.EventOriginStatsMap {
			for _, eventDetailStats := range eventOriginStats.EventDetailStatsMap {
				displayEventStat.EventCount = displayEventStat.EventCount + eventDetailStats.EventCount
			}
		}
	}
	return displayEventList
}

func (asUI *EventListView) convertToListData(displayEventList []*displaydata.DisplayEventStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(displayEventList))
	for _, d := range displayEventList {
		listData = append(listData, d)
	}
	return listData
}

func (asUI *EventListView) PreRowDisplay(data uiCommon.IData, isSelected bool) string {
	return ""
}

func (asUI *EventListView) updateHeader(g *gocui.Gui, v *gocui.View) (int, error) {
	fmt.Fprintf(v, "\nTODO: Show summary stats")
	return 3, nil
}
