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
	"regexp"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/jroimartin/gocui"
)

type HelpTextTipsWidget struct {
	masterUI        masterUIInterface.MasterUIInterface
	name            string
	height          int
	formatTextRegex *regexp.Regexp
	helpTextTips    string
}

func NewHelpTextTipsWidget(masterUI masterUIInterface.MasterUIInterface, name string, height int) *HelpTextTipsWidget {
	w := &HelpTextTipsWidget{masterUI: masterUI, name: name, height: height}
	w.formatTextRegex = regexp.MustCompile(`\*\*([^\*]*)*\*\*`)
	return w
}

func (w *HelpTextTipsWidget) Name() string {
	return w.name
}

func (w *HelpTextTipsWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	top := maxY - w.height
	if top < 0 {
		top = 0
	}
	right := maxX - 1
	if right < 1 {
		right = 1
	}

	v, err := g.SetView(w.name, 0, top, right, maxY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Frame = false
		return w.showHelpTextTips(g, v)
	}
	return nil
}

func (w *HelpTextTipsWidget) SetHelpTextTips(g *gocui.Gui, helpTextTips string) error {

	w.helpTextTips = w.formatText(helpTextTips)
	v, err := g.View(w.name)
	if err != nil {
		return err
	}
	return w.showHelpTextTips(g, v)
}

func (w *HelpTextTipsWidget) showHelpTextTips(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
	fmt.Fprint(v, w.helpTextTips)
	return nil
}

func (w *HelpTextTipsWidget) formatText(text string) string {
	return w.formatTextRegex.ReplaceAllString(text, "\033[37;1m${1}\033[0m")
}
