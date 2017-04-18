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
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

type RequestsInfoWidget struct {
	masterUI   masterUIInterface.MasterUIInterface
	name       string
	height     int
	detailView *AppDetailView
	appMdMgr   *app.AppMetadataManager
}

func NewRequestsInfoWidget(masterUI masterUIInterface.MasterUIInterface, name string, height int, detailView *AppDetailView) *RequestsInfoWidget {
	appMdMgr := detailView.GetEventProcessor().GetMetadataManager().GetAppMdManager()
	return &RequestsInfoWidget{masterUI: masterUI, name: name, height: height, detailView: detailView, appMdMgr: appMdMgr}
}

func (w *RequestsInfoWidget) Name() string {
	return w.name
}

func (w *RequestsInfoWidget) Layout(g *gocui.Gui) error {

	topOffset := w.detailView.GetTopOffset()
	if w.masterUI.IsMinimizeHeader() {
		// This will hide this view by displaying it off-view (negative top)
		topOffset = 0
	}

	maxX, _ := g.Size()
	top := topOffset - w.height - 1
	width := int(math.Floor((float64(maxX) / 2) + 2))

	v, err := g.SetView(w.name, 0, top, width, top+w.height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Frame = true
	}
	//v.Title = "App Request Info"
	v.Title = fmt.Sprintf("App: %v", w.getAppName())
	w.refreshDisplay(g)
	return nil
}

func (w *RequestsInfoWidget) getAppName() string {
	appMetadata := w.appMdMgr.FindAppMetadata(w.detailView.appId)
	appName := appMetadata.Name
	return appName
}

func (w *RequestsInfoWidget) refreshDisplay(g *gocui.Gui) error {

	v, err := g.View(w.name)
	if err != nil {
		return err
	}

	v.Clear()

	if w.detailView.appId == "" {
		fmt.Fprintln(v, "No application selected")
		return nil
	}

	m := w.detailView.GetDisplayedEventData().AppMap
	appStats := m[w.detailView.appId]
	if appStats == nil {
		return nil
	}

	avgResponseTimeL60Info := "--"
	if appStats.TotalTraffic.AvgResponseL60Time >= 0 {
		avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL60Time / 1000000
		avgResponseTimeL60Info = fmt.Sprintf("%8.1f", avgResponseTimeMs)
	}

	avgResponseTimeL10Info := "--"
	if appStats.TotalTraffic.AvgResponseL10Time >= 0 {
		avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL10Time / 1000000
		avgResponseTimeL10Info = fmt.Sprintf("%8.1f", avgResponseTimeMs)
	}

	avgResponseTimeL1Info := "--"
	if appStats.TotalTraffic.AvgResponseL1Time >= 0 {
		avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL1Time / 1000000
		avgResponseTimeL1Info = fmt.Sprintf("%8.1f", avgResponseTimeMs)
	}

	fmt.Fprintf(v, "%16v", "")
	fmt.Fprintf(v, "    1sec   10sec   60sec\n")

	fmt.Fprintf(v, "%16v", "  HTTP(S) Rate:")
	fmt.Fprintf(v, "%8v", appStats.TotalTraffic.EventL1Rate)
	fmt.Fprintf(v, "%8v", appStats.TotalTraffic.EventL10Rate)
	fmt.Fprintf(v, "%8v\n", appStats.TotalTraffic.EventL60Rate)

	fmt.Fprintf(v, "%16v", "Rspns Time(ms):")
	fmt.Fprintf(v, "%8v", avgResponseTimeL1Info)
	fmt.Fprintf(v, "%8v", avgResponseTimeL10Info)
	fmt.Fprintf(v, "%8v\n", avgResponseTimeL60Info)
	fmt.Fprintf(v, "%v", util.BRIGHT_WHITE)
	fmt.Fprintf(v, "%v", util.CLEAR)
	return nil
}
