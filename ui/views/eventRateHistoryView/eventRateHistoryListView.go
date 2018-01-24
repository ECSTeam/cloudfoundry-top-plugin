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

package eventRateHistoryView

import (
	"log"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appViews/appView"
	"github.com/jroimartin/gocui"
)

type EventRateHistoryView struct {
	*dataView.DataListView
}

func NewEventRateHistoryView(masterUI masterUIInterface.MasterUIInterface,
	name string, bottomMargin int,
	eventProcessor *eventdata.EventProcessor) *EventRateHistoryView {

	asUI := &EventRateHistoryView{}

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("TOTAL", true),
		uiCommon.NewSortColumn("BEGIN_TIME", true),
	}

	dataListView := dataView.NewDataListView(masterUI, nil,
		name, 0, bottomMargin,
		eventProcessor, asUI, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.GetListData = asUI.GetListData

	dataListView.SetTitle(func() string { return "Event Rate Peak History" })
	dataListView.HelpText = HelpText
	dataListView.HelpTextTips = appView.HelpTextTips

	asUI.DataListView = dataListView

	return asUI

}

func (asUI *EventRateHistoryView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)

	columns = append(columns, columnEventBeginTime())
	columns = append(columns, columnEventEndTime())
	columns = append(columns, columnInterval())

	columns = append(columns, columnTotalRateHigh())
	//columns = append(columns, columnTotalRateLow())

	columns = append(columns, columnHttpStartStopEventRateHigh())
	//columns = append(columns, columnHttpStartStopEventRateLow())

	columns = append(columns, columnContainerMetricEventRateHigh())
	//columns = append(columns, columnContainerMetricEventRateLow())

	columns = append(columns, columnLogMessageEventRateHigh())
	//columns = append(columns, columnLogMessageEventRateLow())

	columns = append(columns, columnValueMetricEventRateHigh())
	//columns = append(columns, columnValueMetricEventRateLow())

	columns = append(columns, columnCounterEventRateHigh())
	//columns = append(columns, columnCounterEventRateLow())

	columns = append(columns, columnErrorEventRateHigh())
	//columns = append(columns, columnErrorEventRateLow())

	//columns = append(columns, columnOtherEventRateHigh())
	//columns = append(columns, columnOtherEventRateLow())

	return columns
}

func (asUI *EventRateHistoryView) initializeCallback(g *gocui.Gui, viewName string) error {

	keys := [...]gocui.Key{gocui.KeyArrowUp, gocui.KeyArrowDown, gocui.KeyPgdn, gocui.KeyPgup, gocui.KeyArrowRight, gocui.KeyArrowLeft}
	for _, key := range keys {
		if err := g.SetKeybinding(viewName, key, gocui.ModNone, asUI.highlightNavigationAction); err != nil {
			log.Panicln(err)
		}
	}
	if err := g.SetKeybinding(viewName, gocui.KeyEsc, gocui.ModNone, asUI.highlightNavigationEscAction); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (asUI *EventRateHistoryView) highlightNavigationEscAction(g *gocui.Gui, v *gocui.View) error {
	return asUI.highlightNavigationPauseIfNeeded(g, v, true)

}

func (asUI *EventRateHistoryView) highlightNavigationAction(g *gocui.Gui, v *gocui.View) error {
	return asUI.highlightNavigationPauseIfNeeded(g, v, false)
}

// To make the user's experience with scrolling through rows a better experience
// we auto-pause the display update if we are sorted by BEGIN_TIME or END_TIME in reverse order
// This is because new records are coming in every second and with reverse order sort on
// these time fields, its difficult to navigate as the highlighted row keeps shifting down.
// The ESC key will (like on all screens) de-select row and reset to top of list
func (asUI *EventRateHistoryView) highlightNavigationPauseIfNeeded(g *gocui.Gui, v *gocui.View, isEscKey bool) error {

	sortColumns := asUI.GetListWidget().GetSortColumns()
	if len(sortColumns) > 0 {
		sortColumn := sortColumns[0]
		if (sortColumn.Id == "BEGIN_TIME" || sortColumn.Id == "END_TIME") && sortColumn.ReverseSort {
			if isEscKey {
				if asUI.GetMasterUI().GetDisplayPaused() {
					asUI.GetMasterUI().SetDisplayPaused(false)
				}
			} else {
				if !asUI.GetMasterUI().GetDisplayPaused() {
					asUI.GetMasterUI().SetDisplayPaused(true)
				}
			}
		}
	}
	return nil
}

func (asUI *EventRateHistoryView) GetListData() []uiCommon.IData {

	// This is a special case for this datatype -- Since the data we're viewing
	// is not in the normal "refresh interval" snapshot, we need to freeze
	// the data ourselves.
	isPaused := asUI.GetMasterUI().GetDisplayPaused()
	if isPaused {
		listData := asUI.GetDisplayedListData()
		if listData != nil {
			return asUI.GetDisplayedListData()
		}
	}

	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *EventRateHistoryView) postProcessData() []*DisplayEventRateHistoryStats {

	ep := asUI.GetEventProcessor()
	erh := ep.GetCurrentEventRateHistory()
	eventRateList := erh.GetDisplayedHistory()

	displayList := make([]*DisplayEventRateHistoryStats, 0, len(eventRateList))
	for _, eventRate := range eventRateList {
		displayEventRate := NewDisplayEventRateHistoryStats(eventRate)
		displayEventRate.Duration = eventRate.EndTime.Sub(eventRate.BeginTime)
		displayList = append(displayList, displayEventRate)
	}
	return displayList
}

func (asUI *EventRateHistoryView) convertToListData(displayEventList []*DisplayEventRateHistoryStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(displayEventList))
	for _, d := range displayEventList {
		listData = append(listData, d)
	}
	return listData
}
