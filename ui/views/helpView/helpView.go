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

package helpView

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/jroimartin/gocui"
)

type HelpView struct {
	masterUI      masterUIInterface.MasterUIInterface
	name          string
	width         int
	height        int
	helpText      string
	displayText   string
	helpTextLines int

	viewOffset int
}

func NewHelpView(masterUI masterUIInterface.MasterUIInterface, name string, width, height int, helpText string) *HelpView {
	hv := &HelpView{masterUI: masterUI, name: name, width: width, height: height, helpText: helpText}
	hv.helpTextLines = strings.Count(helpText, "\n")
	return hv
}

func (w *HelpView) Name() string {
	return w.name
}

func (w *HelpView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	top := maxY/2 - (w.height / 2)
	v, err := g.SetView(w.name, maxX/2-(w.width/2), top, maxX/2+(w.width/2), top+w.height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Help (press ENTER to close, DOWN/UP arrow to scroll)"
		v.Frame = true

		if w.displayText == "" {
			re := regexp.MustCompile(`\*\*([^\*]*)*\*\*`)
			w.displayText = re.ReplaceAllString(w.helpText, "\033[37;1m${1}\033[0m")
		}

		fmt.Fprintf(v, w.displayText)
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeHelpView); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.closeHelpView); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, 'x', gocui.ModNone, w.closeHelpView); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.arrowUp); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.arrowDown); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyPgup, gocui.ModNone, w.pageUp); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyPgdn, gocui.ModNone, w.pageDown); err != nil {
			log.Panicln(err)
		}

		g.Highlight = true

		fgColor := gocui.ColorWhite | gocui.AttrBold
		//v.FgColor = fgColor
		g.SelFgColor = fgColor

		//bgColor := gocui.ColorWhite
		//v.BgColor = bgColor
		//g.SelBgColor = bgColor

		if err := w.masterUI.SetCurrentViewOnTop(g); err != nil {
			log.Panicln(err)
		}

	}
	return nil
}

func (w *HelpView) closeHelpView(g *gocui.Gui, v *gocui.View) error {
	g.Highlight = false
	g.SelBgColor = gocui.ColorBlack
	g.SelFgColor = gocui.ColorWhite
	if err := w.masterUI.CloseView(w); err != nil {
		return err
	}
	return nil
}

func (w *HelpView) arrowUp(g *gocui.Gui, v *gocui.View) error {
	if w.viewOffset > 0 {
		w.viewOffset--
		v.SetOrigin(0, w.viewOffset)
	}
	return nil
}

func (w *HelpView) arrowDown(g *gocui.Gui, v *gocui.View) error {
	if w.viewOffset <= (w.helpTextLines - w.height) {
		w.viewOffset++
		v.SetOrigin(0, w.viewOffset)
	}
	return nil
}

func (w *HelpView) pageUp(g *gocui.Gui, v *gocui.View) error {
	realHeight := w.height - 1
	if w.viewOffset > 0 {
		w.viewOffset = w.viewOffset - realHeight
		if w.viewOffset < 0 {
			w.viewOffset = 0
		}
		v.SetOrigin(0, w.viewOffset)
	}
	return nil
}

func (w *HelpView) pageDown(g *gocui.Gui, v *gocui.View) error {
	h := w.height - 1
	textLines := w.helpTextLines

	w.viewOffset = w.viewOffset + h
	if !(w.viewOffset < textLines && (textLines-h) > w.viewOffset) {
		w.viewOffset = textLines - h
	}
	v.SetOrigin(0, w.viewOffset)
	return nil
}
