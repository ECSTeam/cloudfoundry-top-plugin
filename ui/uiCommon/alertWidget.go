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

package uiCommon

import (
	"errors"
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

type AlertWidget struct {
	masterUI masterUIInterface.MasterUIInterface
	name     string
	height   int
	message  string
}

func NewAlertWidget(masterUI masterUIInterface.MasterUIInterface, name string, height int) *AlertWidget {
	return &AlertWidget{masterUI: masterUI, name: name, height: height}
}

func (w *AlertWidget) Name() string {
	return w.name
}

func (w *AlertWidget) SetHeight(height int) {
	w.height = height
}

func (w *AlertWidget) SetMessage(msg string) {
	w.message = msg
}

func (w *AlertWidget) Layout(g *gocui.Gui) error {
	maxX, _ := g.Size()

	top := w.masterUI.GetHeaderSize() + 1
	v, err := g.SetView(w.name, -1, top-1, maxX, top+w.height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		} else {
			// We need to ensure the alert line does not popup over another window (e.g., log window)
			w.masterUI.SetCurrentViewOnTop(g)
		}
		v.Frame = false
	} else {

	}

	v.Clear()
	fmt.Fprintf(v, " %v", util.WHITE_TEXT_RED_BG)
	if w.message != "" {
		fmt.Fprintln(v, w.message)
	} else {
		fmt.Fprintln(v, "No ALERT message")
	}
	fmt.Fprintf(v, "%v", util.CLEAR)

	return nil
}
