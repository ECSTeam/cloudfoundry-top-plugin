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

package headerView

import (
	//"fmt"
	"errors"
	"fmt"
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/config"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventrouting"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/dataCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

type HeaderWidget struct {
	masterUI   masterUIInterface.MasterUIInterface
	router     *eventrouting.EventRouter
	commonData *dataCommon.CommonData
	name       string

	HeaderSize int
}

func NewHeaderWidget(masterUI masterUIInterface.MasterUIInterface,
	name string,
	router *eventrouting.EventRouter,
	commonData *dataCommon.CommonData) *HeaderWidget {

	headerSize := 6
	return &HeaderWidget{masterUI: masterUI, name: name, router: router, commonData: commonData, HeaderSize: headerSize}
}

func (w *HeaderWidget) Name() string {
	return w.name
}

func (w *HeaderWidget) Layout(g *gocui.Gui) error {
	maxX, _ := g.Size()
	_, err := g.SetView(w.name, 0, 0, maxX-1, w.HeaderSize)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		//fmt.Fprint(v, w.body)
	}
	return nil
}

func (w *HeaderWidget) UpdateDisplay(g *gocui.Gui) error {

	v, err := g.View("headerView")
	if err != nil {
		return err
	}

	router := w.router
	processor := router.GetProcessor()
	eventData := processor.GetDisplayedEventData()
	statsTime := eventData.StatsTime

	currentEventRate := processor.GetCurrentEventRateHistory().GetCurrentRate()
	eventsText := fmt.Sprintf("%v (%v/sec)", util.FormatUint64(router.GetEventCount()), currentEventRate)
	runtimeSeconds := Round(statsTime.Sub(router.GetStartTime()), time.Second)
	v.Clear()

	fmt.Fprintf(v, "Events: ")
	fmt.Fprintf(v, "%-27v", eventsText)
	if runtimeSeconds < time.Second*config.WarmUpSeconds {
		warmUpTimeRemaining := (time.Second * config.WarmUpSeconds) - runtimeSeconds
		fmt.Fprintf(v, util.DIM_GREEN)
		fmt.Fprintf(v, " Warm-up: %-10v ", warmUpTimeRemaining)
		fmt.Fprintf(v, util.CLEAR)
	} else {
		fmt.Fprintf(v, "Duration: %-10v ", runtimeSeconds)
	}

	fmt.Fprintf(v, "   %v\n", statsTime.Format("01-02-2006 15:04:05"))

	if w.masterUI.GetDisplayPaused() {
		fmt.Fprintf(v, util.REVERSE_GREEN)
		fmt.Fprintf(v, " Display update paused \n")
		fmt.Fprintf(v, util.CLEAR)
	} else {
		fmt.Fprintf(v, "Target: %-78.78v\n", w.masterUI.GetTargetDisplay())
	}

	headerStackLines, err := w.updateHeaderStack(g, v)
	if err != nil {
		return err
	}

	// Base header is 2 rows plus 1 for border
	headerLines := 3 + headerStackLines
	w.HeaderSize = headerLines

	return nil
}

func Round(d, r time.Duration) time.Duration {
	if r <= 0 {
		return d
	}
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d = d - m
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}
