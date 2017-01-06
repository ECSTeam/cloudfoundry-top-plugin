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

package appView

import (
	"errors"
	"fmt"
	"log"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/jroimartin/gocui"
)

type FilterWidget struct {
	masterUI masterUIInterface.MasterUIInterface
	name     string
	width    int
	height   int
}

func NewFilterWidget(masterUI masterUIInterface.MasterUIInterface, name string, width, height int) *FilterWidget {
	return &FilterWidget{masterUI: masterUI, name: name, width: width, height: height}
}

func (w *FilterWidget) Name() string {
	return w.name
}

func (w *FilterWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(w.name, maxX/2-(w.width/2), maxY/2-(w.height/2), maxX/2+(w.width/2), maxY/2+(w.height/2))
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Title = "Filter (press ENTER to close)"
		v.Frame = true
		fmt.Fprintln(v, "Future home of filter screen")
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeFilterWidget); err != nil {
			return err
		}

		if err := w.masterUI.SetCurrentViewOnTop(g); err != nil {
			log.Panicln(err)
		}

	}
	return nil
}

func (w *FilterWidget) closeFilterWidget(g *gocui.Gui, v *gocui.View) error {
	if err := w.masterUI.CloseView(w); err != nil {
		return err
	}
	return nil
}
