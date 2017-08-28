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

package dataView

import (
	"log"
	"sync"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/helpView"
	"github.com/jroimartin/gocui"
)

type updateHeaderCallback func(g *gocui.Gui, v *gocui.View) (int, error)
type actionCallback func(g *gocui.Gui, v *gocui.View) error
type initializeCallback func(g *gocui.Gui, viewName string) error
type preRowDisplayCallback func(data uiCommon.IData, isSelected bool) string
type refreshDisplayCallback func(g *gocui.Gui) error
type IColumnOwner interface{}

type GetListData func() []uiCommon.IData

type DataListView struct {
	masterUI              masterUIInterface.MasterUIInterface
	parentView            DataListViewInterface
	detailView            DataListViewInterface
	name                  string
	topMargin             int
	bottomMargin          int
	eventProcessor        *eventdata.EventProcessor
	mdGlobalMgr           *metadata.GlobalManager
	mu                    sync.Mutex
	listWidget            *uiCommon.ListWidget
	initialized           bool
	Title                 string
	HelpText              string
	HelpTextTips          string
	InitializeCallback    initializeCallback
	UpdateHeaderCallback  updateHeaderCallback
	PreRowDisplayCallback preRowDisplayCallback
	columnOwner           IColumnOwner

	RefreshDisplayCallback refreshDisplayCallback
	GetListData            GetListData
}

func NewDataListView(masterUI masterUIInterface.MasterUIInterface,
	parentView DataListViewInterface,
	name string, topMargin, bottomMargin int,
	eventProcessor *eventdata.EventProcessor,
	columnOwner IColumnOwner,
	columnDefinitions []*uiCommon.ListColumn,
	defaultSortColumns []*uiCommon.SortColumn) *DataListView {

	asUI := &DataListView{
		masterUI:       masterUI,
		parentView:     parentView,
		name:           name,
		topMargin:      topMargin,
		bottomMargin:   bottomMargin,
		eventProcessor: eventProcessor,
		columnOwner:    columnOwner,
	}

	asUI.mdGlobalMgr = eventProcessor.GetMetadataManager()

	listWidget := uiCommon.NewListWidget(asUI.masterUI, asUI.name,
		asUI.bottomMargin, asUI, columnDefinitions, columnOwner, defaultSortColumns)
	listWidget.PreRowDisplayFunc = asUI.PreRowDisplay

	asUI.listWidget = listWidget

	return asUI
}

// Get the top offset where the data view should open
func (asUI *DataListView) GetTopOffset() int {
	size := asUI.masterUI.GetTopMargin() + 1
	if !asUI.listWidget.IsSelectColumnMode() {
		size = size + asUI.topMargin
	}
	return size
}

func (asUI *DataListView) Name() string {
	return asUI.name
}

func (asUI *DataListView) GetTopMargin() int {
	return asUI.topMargin
}

func (asUI *DataListView) SetTitle(title string) {
	asUI.listWidget.Title = title
}

func (asUI *DataListView) GetMargins() (int, int) {
	return asUI.topMargin, asUI.bottomMargin
}

func (asUI *DataListView) GetMasterUI() masterUIInterface.MasterUIInterface {
	return asUI.masterUI
}

func (asUI *DataListView) GetParentView() DataListViewInterface {
	return asUI.parentView
}

func (asUI *DataListView) GetDetailView() DataListViewInterface {
	return asUI.detailView
}

func (asUI *DataListView) SetDetailView(detailView DataListViewInterface) {
	asUI.detailView = detailView
}

func (asUI *DataListView) GetListWidget() *uiCommon.ListWidget {
	return asUI.listWidget
}

func (asUI *DataListView) GetEventProcessor() *eventdata.EventProcessor {
	return asUI.eventProcessor
}

func (asUI *DataListView) GetMdGlobalMgr() *metadata.GlobalManager {
	return asUI.mdGlobalMgr
}

func (asUI *DataListView) Layout(g *gocui.Gui) error {
	if !asUI.initialized {
		asUI.initialized = true
		asUI.initialize(g)
	}
	asUI.masterUI.SetHelpTextTips(g, asUI.HelpTextTips)
	//topOffset := asUI.getTopOffset()
	// TODO: Is this correct?
	//asUI.listWidget.SetTopMargin(topOffset)
	return asUI.listWidget.Layout(g)
}

func (asUI *DataListView) initialize(g *gocui.Gui) {
	if err := g.SetKeybinding(asUI.name, 'h', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			helpView := helpView.NewHelpView(asUI.masterUI, "helpView", 75, 17, asUI.HelpText)
			asUI.masterUI.LayoutManager().Add(helpView)
			asUI.masterUI.SetCurrentViewOnTop(g)
			return nil
		}); err != nil {
		log.Panicln(err)
	}
	if asUI.InitializeCallback != nil {
		asUI.InitializeCallback(g, asUI.name)
	}
}

func (asUI *DataListView) GetCurrentEventData() *eventdata.EventData {
	return asUI.eventProcessor.GetCurrentEventData()
}

func (asUI *DataListView) GetDisplayedEventData() *eventdata.EventData {
	return asUI.eventProcessor.GetDisplayedEventData()
}

func (asUI *DataListView) GetDisplayedListData() []uiCommon.IData {
	return asUI.listWidget.GetListData()
}

func (asUI *DataListView) RefreshDisplay(g *gocui.Gui) error {
	var err error

	/*
		TODO: Need to figure out framework for allowing data views to add rows to header
		headerSize, err2 := asUI.updateHeader(g)
		if err2 != nil {
			return err2
		}
		asUI.masterUI.SetStatsSummarySize(headerSize)
	*/

	if asUI.RefreshDisplayCallback != nil {
		err = asUI.RefreshDisplayCallback(g)
		if err != nil {
			return err
		}
	}

	err = asUI.refreshListDisplay(g)
	if err != nil {
		return err
	}

	return nil
}

func (asUI *DataListView) refreshListDisplay(g *gocui.Gui) error {
	err := asUI.listWidget.RefreshDisplay(g)
	if err != nil {
		return err
	}
	return err
}

func (asUI *DataListView) UpdateDisplay(g *gocui.Gui) error {
	asUI.updateData()
	return asUI.RefreshDisplay(g)
}

func (asUI *DataListView) CloseDetailView(g *gocui.Gui, v *gocui.View) error {
	if err := asUI.GetMasterUI().CloseView(asUI); err != nil {
		return err
	}
	return nil
}

// XXX
func (asUI *DataListView) updateData() {
	//asUI.eventProcessor.UpdateData()
	listData := asUI.GetListData()
	asUI.listWidget.SetListData(listData)
}

func (asUI *DataListView) PreRowDisplay(data uiCommon.IData, isSelected bool) string {
	if asUI.PreRowDisplayCallback != nil {
		return asUI.PreRowDisplayCallback(data, isSelected)
	}
	return ""
}

func (asUI *DataListView) updateHeader(g *gocui.Gui) (int, error) {

	v, err := g.View("headerView")
	if err != nil {
		return 0, err
	}

	if asUI.UpdateHeaderCallback != nil {
		return asUI.UpdateHeaderCallback(g, v)
	}
	return 0, nil
}
