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

	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

type AlertWidget struct {
	name      string
	topMargin int
	height    int
	message   string
}

func NewAlertWidget(name string, topMargin, height int) *AlertWidget {
	return &AlertWidget{name: name, topMargin: topMargin, height: height}
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
	v, err := g.SetView(w.name, 0, w.topMargin-1, maxX-1, w.topMargin+w.height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Frame = false

		v.Clear()
		fmt.Fprintf(v, "%v", util.REVERSE_RED)

		if w.message != "" {
			fmt.Fprintln(v, w.message)
		} else {
			fmt.Fprintln(v, "No ALERT message specified")
		}
		fmt.Fprintf(v, "%v", util.CLEAR)
		fmt.Fprintln(v, "line 2")
		fmt.Fprintln(v, "line 3")
		fmt.Fprintln(v, "line 4")
	}
	return nil
}
