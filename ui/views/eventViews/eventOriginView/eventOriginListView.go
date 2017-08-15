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

package eventOriginView

import (
	"fmt"
	"log"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/eventViews/eventDetailView"
	"github.com/jroimartin/gocui"
)

type EventOriginListView struct {
	*dataView.DataListView
	eventType events.Envelope_EventType
}

func NewEventOriginListView(masterUI masterUIInterface.MasterUIInterface,
	parentView dataView.DataListViewInterface,
	name string, topMargin, bottomMargin int,
	eventProcessor *eventdata.EventProcessor, eventTypeName string) *EventOriginListView {

	eventType := events.Envelope_EventType(events.Envelope_EventType_value[eventTypeName])

	asUI := &EventOriginListView{
		eventType: eventType,
	}

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("COUNT", true),
		uiCommon.NewSortColumn("ORIGIN", true),
	}

	dataListView := dataView.NewDataListView(masterUI, parentView,
		name, topMargin, bottomMargin,
		eventProcessor, asUI, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.UpdateHeaderCallback = asUI.updateHeader
	dataListView.GetListData = asUI.GetListData

	dataListView.SetTitle(fmt.Sprintf("Event Origin List - Event Type: %v", eventTypeName))
	dataListView.HelpText = HelpText
	dataListView.HelpTextTips = HelpTextTips

	asUI.DataListView = dataListView

	return asUI

}

func (asUI *EventOriginListView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)

	columns = append(columns, columnEventOrigin())
	columns = append(columns, columnEventCount())

	return columns
}

func (asUI *EventOriginListView) initializeCallback(g *gocui.Gui, viewName string) error {

	// TODO: This needs to be handled in dataListView someplace for child (detailed) views as all of them will need a back action
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

// TODO: Need to put this in common dataListView - but allow for callback to do special close processing (as needed by appDetailView to close other views)
func (asUI *EventOriginListView) closeAppDetailView(g *gocui.Gui, v *gocui.View) error {
	if err := asUI.GetMasterUI().CloseView(asUI); err != nil {
		return err
	}
	return nil
}

func (asUI *EventOriginListView) enterAction(g *gocui.Gui, v *gocui.View) error {

	highlightKey := asUI.GetListWidget().HighlightKey()
	if highlightKey != "" {
		topMargin, bottomMargin := asUI.GetMargins()

		detailView := eventDetailView.NewEventDetailListView(asUI.GetMasterUI(), asUI,
			"eventDetailView",
			topMargin, bottomMargin,
			asUI.GetEventProcessor(),
			asUI.eventType.String(),
			highlightKey)

		asUI.SetDetailView(detailView)
		asUI.GetMasterUI().OpenView(g, detailView)
	}

	return nil
}

func (asUI *EventOriginListView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *EventOriginListView) postProcessData() []*DisplayEventOriginStats {

	eventTypeMap := asUI.GetDisplayedEventData().EventTypeMap

	eventStats := eventTypeMap[asUI.eventType]
	eventOriginStatsMap := eventStats.EventOriginStatsMap
	displayEventOriginList := make([]*DisplayEventOriginStats, 0, len(eventOriginStatsMap))

	for _, eventOriginStats := range eventOriginStatsMap {
		displayEventOriginStat := NewDisplayEventOriginStats(eventOriginStats)
		displayEventOriginList = append(displayEventOriginList, displayEventOriginStat)
		for _, eventDetailStats := range eventOriginStats.EventDetailStatsMap {
			displayEventOriginStat.EventCount = displayEventOriginStat.EventCount + eventDetailStats.EventCount
		}
	}
	return displayEventOriginList
}

func (asUI *EventOriginListView) convertToListData(displayEventList []*DisplayEventOriginStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(displayEventList))
	for _, d := range displayEventList {
		listData = append(listData, d)
	}
	return listData
}

func (asUI *EventOriginListView) PreRowDisplay(data uiCommon.IData, isSelected bool) string {
	return ""
}

func (asUI *EventOriginListView) updateHeader(g *gocui.Gui, v *gocui.View) (int, error) {
	fmt.Fprintf(v, "\nTODO: Show summary stats")
	return 3, nil
}
