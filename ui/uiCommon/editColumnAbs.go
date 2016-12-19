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
	"log"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/jroimartin/gocui"
)

type initialLayoutCallbackFunc func(g *gocui.Gui, v *gocui.View) error
type updateLayoutCallbackFunc func(g *gocui.Gui, v *gocui.View) error
type refreshDisplayCallbackFunc func(g *gocui.Gui, v *gocui.View) error
type applyActionCallbackFunc func(g *gocui.Gui, v *gocui.View) error
type cancelActionCallbackFunc func(g *gocui.Gui, v *gocui.View) error

type EditColumnViewAbs struct {
	masterUI         masterUIInterface.MasterUIInterface
	name             string
	width            int
	height           int
	title            string
	listWidget       *ListWidget
	minTopViewMargin int

	initialLayoutCallbackFunc  initialLayoutCallbackFunc
	updateLayoutCallbackFunc   updateLayoutCallbackFunc
	refreshDisplayCallbackFunc refreshDisplayCallbackFunc
	applyActionCallbackFunc    applyActionCallbackFunc
	cancelActionCallbackFunc   cancelActionCallbackFunc

	priorStateOfDisplayPaused bool
}

func NewEditColumnViewAbs(masterUI masterUIInterface.MasterUIInterface, name string, listWidget *ListWidget) *EditColumnViewAbs {
	w := &EditColumnViewAbs{masterUI: masterUI, name: name, listWidget: listWidget}
	w.priorStateOfDisplayPaused = listWidget.displayView.GetDisplayPaused()
	listWidget.displayView.SetDisplayPaused(true)
	w.minTopViewMargin = 9
	return w
}

func (w *EditColumnViewAbs) Name() string {
	return w.name
}

func (w *EditColumnViewAbs) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	top := maxY/2 - (w.height / 2)

	if top < w.minTopViewMargin {
		top = w.minTopViewMargin
	}
	//bottom := maxY/2 + (w.height / 2)

	bottom := top + w.height

	//v, err := g.SetView(w.name, maxX/2-(w.width/2), maxY/2-(w.height/2), maxX/2+(w.width/2), maxY/2+(w.height/2))
	v, err := g.SetView(w.name, maxX/2-(w.width/2), top, maxX/2+(w.width/2), bottom)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Title = w.title
		v.Frame = true
		g.Highlight = true
		//g.SelFgColor = gocui.ColorGreen
		//g.SelFgColor = gocui.ColorWhite | gocui.AttrBold
		//v.BgColor = gocui.ColorRed
		//v.FgColor = gocui.ColorGreen
		fmt.Fprintln(v, "...")
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.applyAction); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, 'x', gocui.ModNone, w.applyAction); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowRight, gocui.ModNone, w.keyArrowRightAction); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowLeft, gocui.ModNone, w.keyArrowLeftAction); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.cancelAction); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, 'q', gocui.ModNone, w.cancelAction); err != nil {
			return err
		}

		// If the current selected column is not within view, then select first column
		if !w.listWidget.isColumnVisable(g, w.listWidget.selectedColumnId) {
			w.listWidget.selectedColumnId = w.listWidget.columns[0].id
		}

		if w.initialLayoutCallbackFunc != nil {
			w.initialLayoutCallbackFunc(g, v)
		}

		if err := w.masterUI.SetCurrentViewOnTop(g); err != nil {
			log.Panicln(err)
		}
		w.RefreshDisplay(g)
	} else {
		if w.updateLayoutCallbackFunc != nil {
			w.updateLayoutCallbackFunc(g, v)
		}
	}
	return nil
}

func (w *EditColumnViewAbs) RefreshDisplay(g *gocui.Gui) error {
	v, err := g.View(w.name)
	if err != nil {
		return err
	}
	if w.refreshDisplayCallbackFunc != nil {
		w.refreshDisplayCallbackFunc(g, v)
	}
	return w.listWidget.RefreshDisplay(g)
}

func (w *EditColumnViewAbs) applyAction(g *gocui.Gui, v *gocui.View) error {
	if w.applyActionCallbackFunc != nil {
		w.applyActionCallbackFunc(g, v)
	}
	return w.closeView(g, v)
}

func (w *EditColumnViewAbs) cancelAction(g *gocui.Gui, v *gocui.View) error {

	if w.cancelActionCallbackFunc != nil {
		w.cancelActionCallbackFunc(g, v)
	}

	return w.closeView(g, v)
}

func (w *EditColumnViewAbs) closeView(g *gocui.Gui, v *gocui.View) error {

	w.listWidget.enableSelectColumnMode(false)
	if err := w.masterUI.CloseViewByName(w.name); err != nil {
		return err
	}

	// TODO: Is this the correct spot to do this?
	w.masterUI.SetMinimizeHeader(g, false)

	w.listWidget.displayView.SetDisplayPaused(w.priorStateOfDisplayPaused)
	w.listWidget.displayView.RefreshDisplay(g)
	return nil
}

func (w *EditColumnViewAbs) keyArrowRightAction(g *gocui.Gui, v *gocui.View) error {
	columnId := w.listWidget.selectedColumnId
	columns := w.listWidget.columns
	columnsLen := len(columns)
	for i, col := range columns {
		if col.id == columnId && i+1 < columnsLen {
			columnId = columns[i+1].id
			break
		}
	}
	//writeFooter(g, fmt.Sprintf("\r columnId: %v", columnId) )
	w.listWidget.selectedColumnId = columnId
	w.RefreshDisplay(g)
	return w.listWidget.scollSelectedColumnIntoView(g)
}

func (w *EditColumnViewAbs) keyArrowLeftAction(g *gocui.Gui, v *gocui.View) error {
	columnId := w.listWidget.selectedColumnId
	columns := w.listWidget.columns
	for i, col := range columns {
		if col.id == columnId && i > 0 {
			columnId = columns[i-1].id
			break
		}
	}
	//writeFooter(g, fmt.Sprintf("\r columnId: %v", columnId) )
	w.listWidget.selectedColumnId = columnId
	w.RefreshDisplay(g)
	return w.listWidget.scollSelectedColumnIntoView(g)
}
