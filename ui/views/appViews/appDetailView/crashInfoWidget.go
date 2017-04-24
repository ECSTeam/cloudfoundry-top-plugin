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

package appDetailView

import (
	"errors"
	"fmt"
	"math"

	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/app"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/crashData"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

type CrashInfoWidget struct {
	masterUI   masterUIInterface.MasterUIInterface
	name       string
	height     int
	detailView *AppDetailView
	appMdMgr   *app.AppMetadataManager
}

func NewCrashInfoWidget(masterUI masterUIInterface.MasterUIInterface, name string, height int, detailView *AppDetailView) *CrashInfoWidget {
	appMdMgr := detailView.GetEventProcessor().GetMetadataManager().GetAppMdManager()
	return &CrashInfoWidget{masterUI: masterUI, name: name, height: height, detailView: detailView, appMdMgr: appMdMgr}
}

func (w *CrashInfoWidget) Name() string {
	return w.name
}

func (w *CrashInfoWidget) Layout(g *gocui.Gui) error {

	topOffset := w.detailView.GetTopOffset()
	if w.masterUI.IsHeaderMinimized() {
		// This will hide this view by displaying it off-view (negative top)
		topOffset = 0
	}

	maxX, _ := g.Size()
	top := topOffset - w.height - 1
	//width := maxX - 1
	left := int(math.Floor((float64(maxX) / 2) + 3))
	width := int(math.Ceil((float64(maxX) / 2) - 4))
	if width < 1 {
		width = 1
	}
	right := left + width

	v, err := g.SetView(w.name, left, top, right, top+w.height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Frame = true
	}
	v.Title = "Crash Info"
	w.refreshDisplay(g)
	return nil
}

func (w *CrashInfoWidget) getAppName() string {
	//appMdMgr := w.detailView.GetEventProcessor().GetMetadataManager().GetAppMdManager()
	appMetadata := w.appMdMgr.FindAppMetadata(w.detailView.appId)
	appName := appMetadata.Name
	return appName
}

func (w *CrashInfoWidget) refreshDisplay(g *gocui.Gui) error {

	v, err := g.View(w.name)
	if err != nil {
		return err
	}

	v.Clear()

	if w.detailView.appId == "" {
		fmt.Fprintln(v, "No application selected")
		return nil
	}

	lastCrashInfo := w.detailView.LastCrashInfo
	lastCrashTimeDisplay := "--"
	if lastCrashInfo != nil {
		lastCrashTimeDisplay = fmt.Sprintf("%v%v", util.DIM_YELLOW, lastCrashInfo.CrashTime.Local().Format("01-02-2006 15:04:05"))
	}

	fmt.Fprintf(v, "%11v", "")
	fmt.Fprintf(v, "    10min   1hr  24hr\n")

	fmt.Fprintf(v, "%11v", "    Crashes:  ")

	fmt.Fprintf(v, "%v%6v", w.getCrashCountColor(w.detailView.Crash10mCount), w.getCrashCount(w.detailView.Crash10mCount))
	fmt.Fprintf(v, "%v%6v", w.getCrashCountColor(w.detailView.Crash1hCount), w.getCrashCount(w.detailView.Crash1hCount))
	fmt.Fprintf(v, "%v%6v", w.getCrashCountColor(w.detailView.Crash24hCount), w.getCrashCount(w.detailView.Crash24hCount))
	fmt.Fprintf(v, "%v\n", util.CLEAR)
	fmt.Fprintf(v, "%11v", " Last crash:")
	fmt.Fprintf(v, " %v", lastCrashTimeDisplay)
	fmt.Fprintf(v, "%v", util.CLEAR)
	return nil
}

func (w *CrashInfoWidget) getCrashCount(crashCount int) string {
	if crashCount > 0 || crashData.IsCacheLoaded() {
		return fmt.Sprintf("%v", crashCount)
	}
	return "--"
}

func (w *CrashInfoWidget) getCrashCountColor(crashCount int) string {
	color := ""
	if crashCount > 0 {
		color = util.DIM_YELLOW
	} else {
		color = util.DIM_WHITE
	}
	return color
}
