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

package ui

import (
	"errors"
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/jroimartin/gocui"
)

type StatusWidget struct {
	masterUI masterUIInterface.MasterUIInterface
	name     string
	status   string
}

func NewStatusWidget(masterUI masterUIInterface.MasterUIInterface, name string) *StatusWidget {
	w := &StatusWidget{masterUI: masterUI, name: name}
	w.status = " "
	return w
}

func (w *StatusWidget) Name() string {
	return w.name
}

func (w *StatusWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	top := maxY - 2
	if top < 0 {
		top = 0
	}
	right := maxX - 1
	if right < 1 {
		right = 1
	}

	v, err := g.SetView(w.name, 0, top, right, maxY+1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Frame = false
		return w.showStatus(g, v)
	}
	return nil
}

// NOTE: To update/refresh UI with status, the function MasterUI.SetStatus() should be called
func (w *StatusWidget) SetStatus(g *gocui.Gui, status string) error {
	w.status = status
	return nil
}

// Called by thread running within the gocui
func (w *StatusWidget) ShowStatus(g *gocui.Gui) error {
	v, err := g.View(w.name)
	if err != nil {
		return err
	}
	return w.showStatus(g, v)
}

func (w *StatusWidget) showStatus(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
	fmt.Fprint(v, w.status)
	return nil
}
