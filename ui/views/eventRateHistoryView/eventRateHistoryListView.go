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
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appViews/appView"
)

type EventRateHistoryView struct {
	*dataView.DataListView
}

func NewEventRateHistoryView(masterUI masterUIInterface.MasterUIInterface,
	name string, bottomMargin int,
	eventProcessor *eventdata.EventProcessor) *EventRateHistoryView {

	asUI := &EventRateHistoryView{}

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("BEGIN_TIME", true),
		//uiCommon.NewSortColumn("EVENT_TYPE", true),
	}

	dataListView := dataView.NewDataListView(masterUI, nil,
		name, 0, bottomMargin,
		eventProcessor, asUI, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.GetListData = asUI.GetListData

	dataListView.SetTitle("Event Rate History")
	dataListView.HelpText = HelpText
	dataListView.HelpTextTips = appView.HelpTextTips

	asUI.DataListView = dataListView

	return asUI

}

func (asUI *EventRateHistoryView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)

	columns = append(columns, columnEventBeginTime())
	columns = append(columns, columnEventEndTime())
	columns = append(columns, columnDuration())

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
	erh := ep.GetEventRateHistory()
	eventRateList := erh.GetHistory()

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
