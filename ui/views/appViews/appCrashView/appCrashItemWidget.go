// Copyright (c) 2017 ECS Team, Inc. - All Rights Reserved
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

package appCrashView

import (
	"errors"
	"fmt"
	"log"

	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/app"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/jroimartin/gocui"
)

type AppCrashItemWidget struct {
	masterUI   masterUIInterface.MasterUIInterface
	name       string
	width      int
	height     int
	detailView *AppCrashView
	crashInfo  *DisplayContainerCrashInfo
	appMdMgr   *app.AppMetadataManager
}

func NewAppCrashItemWidget(masterUI masterUIInterface.MasterUIInterface, name string, width, height int,
	detailView *AppCrashView, crashInfo *DisplayContainerCrashInfo) *AppCrashItemWidget {

	appMdMgr := detailView.GetEventProcessor().GetMetadataManager().GetAppMdManager()
	return &AppCrashItemWidget{masterUI: masterUI, name: name, width: width, height: height,
		detailView: detailView, crashInfo: crashInfo, appMdMgr: appMdMgr}
}

func (w *AppCrashItemWidget) Name() string {
	return w.name
}

func (w *AppCrashItemWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	sideMargin := 5
	left := sideMargin
	right := maxX - sideMargin
	if right <= left+1 {
		right = left + 2
	}
	top := maxY/2 - (w.height / 2)
	if top < 1 {
		top = 1
	}
	bottom := maxY - 2
	if bottom <= top {
		bottom = top + 1
	}
	//v, err := g.SetView(w.name, maxX/2-(w.width/2), maxY/2-(w.height/2), maxX/2+(w.width/2), maxY/2+(w.height/2))
	v, err := g.SetView(w.name, left, top, right, bottom)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Title = "App CRASH Item Detail"
		v.Frame = true
		v.Wrap = true
		if err := g.SetKeybinding(w.name, 'x', gocui.ModNone, w.closeAppCrashItemWidget); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.closeAppCrashItemWidget); err != nil {
			return err
		}
		if err := w.masterUI.SetCurrentViewOnTop(g); err != nil {
			log.Panicln(err)
		}

	}
	w.RefreshDisplay(g)
	return nil
}

func (w *AppCrashItemWidget) closeAppCrashItemWidget(g *gocui.Gui, v *gocui.View) error {
	if err := w.masterUI.CloseView(w); err != nil {
		return err
	}
	return nil
}

func (w *AppCrashItemWidget) UpdateDisplay(g *gocui.Gui) error {
	return w.RefreshDisplay(g)
}

func (w *AppCrashItemWidget) RefreshDisplay(g *gocui.Gui) error {

	v, err := g.View(w.name)
	if err != nil {
		return err
	}

	v.Clear()

	m := w.detailView.GetDisplayedEventData().AppMap
	appStats := m[w.detailView.appId]
	if appStats == nil {
		return nil
	}
	//appMetadata := w.appMdMgr.FindItem(appStats.AppId)

	//fmt.Fprintf(v, "Crash Details - ")
	//fmt.Fprintf(v, "%vx%v:exit view", "\033[37;1m", "\033[0m")
	//fmt.Fprintf(v, "\n")

	maxX, _ := v.Size()

	fmt.Fprintf(v, " \n")
	fmt.Fprintf(v, "   App Crash Time: %v\n", w.crashInfo.CrashTimeFormatted)
	fmt.Fprintf(v, "  Container Index: %v\n", w.crashInfo.ContainerIndex)
	exitDescLabel := " Exit Description: "
	lineBreakIfNeeded := ""
	if len(w.crashInfo.ExitDescription) > (maxX - len(exitDescLabel) - 1) {
		lineBreakIfNeeded = "\n\n"
	}
	fmt.Fprintf(v, "%v%v%v", exitDescLabel, lineBreakIfNeeded, w.crashInfo.ExitDescription)

	//fmt.Fprintf(v, "%v", util.BRIGHT_WHITE)
	//fmt.Fprintf(v, "%v", util.CLEAR)
	fmt.Fprintf(v, "\n")

	return nil
}
