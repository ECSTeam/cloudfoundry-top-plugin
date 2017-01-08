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
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/jroimartin/gocui"
)

type DataListViewInterface interface {
	Name() string
	SetTitle(title string)
	GetMargins() (int, int)
	GetMasterUI() masterUIInterface.MasterUIInterface
	GetParentView() DataListViewInterface
	GetDetailView() DataListViewInterface
	SetDetailView(detailView DataListViewInterface)
	GetListWidget() *uiCommon.ListWidget
	GetEventProcessor() *eventdata.EventProcessor
	Layout(g *gocui.Gui) error
	GetCurrentEventData() *eventdata.EventData
	GetDisplayedEventData() *eventdata.EventData
	RefreshDisplay(g *gocui.Gui) error
	UpdateDisplay(g *gocui.Gui) error
	GetTopOffset() int
	SetAlertSize(alertSize int)
	GetDisplayedListData() []uiCommon.IData
}
