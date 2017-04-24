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

package masterUIInterface

import (
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/dataCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/interfaces/managerUI"
	"github.com/jroimartin/gocui"
)

type MasterUIInterface interface {
	SetCurrentViewOnTop(*gocui.Gui) error
	GetCurrentView(g *gocui.Gui) *gocui.View
	CloseView(managerUI.Manager) error
	CloseViewByName(viewName string) error
	LayoutManager() managerUI.LayoutManagerInterface
	OpenView(g *gocui.Gui, dataView UpdatableView) error
	IsWarmupComplete() bool
	SetHelpTextTips(g *gocui.Gui, helpTextTips string) error
	AddCommonDataViewKeybindings(g *gocui.Gui, viewName string) error
	GetHeaderSize() int
	GetAlertSize() int
	GetTopMargin() int
	//SetStatsSummarySize(statSummarySize int)
	SetHeaderMinimize(g *gocui.Gui, minimizeHeader bool)
	IsHeaderMinimized() bool
	SetEditColumnMode(g *gocui.Gui, editColumnMode bool)
	IsEditColumnMode() bool
	IsPrivileged() bool
	GetCommonData() *dataCommon.CommonData
	GetDisplayPaused() bool
	SetDisplayPaused(paused bool)
	GetTargetDisplay() string
}

type UpdatableView interface {
	Layout(*gocui.Gui) error
	Name() string
	UpdateDisplay(g *gocui.Gui) error
	RefreshDisplay(g *gocui.Gui) error
}
