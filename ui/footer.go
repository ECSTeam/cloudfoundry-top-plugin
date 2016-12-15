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

	"github.com/jroimartin/gocui"
)

type FooterWidget struct {
	name            string
	height          int
	formatTextRegex *regexp.Regexp
}

func NewFooterWidget(name string, height int) *FooterWidget {
	w := &FooterWidget{name: name, height: height}
	w.formatTextRegex = regexp.MustCompile(`\*\*([^\*]*)*\*\*`)
	return w
}

func (w *FooterWidget) Name() string {
	return w.name
}

func (w *FooterWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(w.name, 0, maxY-w.height, maxX-1, maxY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Frame = false
		w.quickHelp(g, v)
	}
	return nil
}

func (w *FooterWidget) quickHelp(g *gocui.Gui, v *gocui.View) error {

	fmt.Fprint(v, w.formatText("**d**:display "))
	fmt.Fprint(v, w.formatText("**q**:quit "))

	fmt.Fprint(v, w.formatText("**x**:exit detail view "))
	fmt.Fprint(v, w.formatText("**h**:help "))
	fmt.Fprintln(v, w.formatText("**UP**/**DOWN** arrow to highlight row"))
	fmt.Fprint(v, w.formatText("**ENTER** to select highlighted row, "))
	fmt.Fprint(v, w.formatText(`**LEFT**/**RIGHT** arrow to scroll columns`))
	return nil
}

func (w *FooterWidget) formatText(text string) string {
	return w.formatTextRegex.ReplaceAllString(text, "\033[37;1m${1}\033[0m")
}
