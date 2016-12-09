package uiCommon

import (
	"errors"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

type menuItemSelectedCallbackFunc func(g *gocui.Gui, v *gocui.View, menuId string) error

type MenuItem struct {
	id    string
	label string
}

func NewMenuItem(id, label string) *MenuItem {
	return &MenuItem{id: id, label: label}
}

type SelectMenuWidget struct {
	masterUI masterUIInterface.MasterUIInterface
	name     string
	width    int
	height   int
	title    string

	menuItemSelectedCallback menuItemSelectedCallbackFunc

	menuPosition int
	menuItems    []*MenuItem
}

func NewSelectMenuWidget(
	masterUI masterUIInterface.MasterUIInterface,
	name string,
	title string,
	menuItems []*MenuItem,
	menuItemSelectedCallback menuItemSelectedCallbackFunc) *SelectMenuWidget {

	w := &SelectMenuWidget{
		masterUI:                 masterUI,
		name:                     name,
		title:                    title,
		menuItems:                menuItems,
		menuItemSelectedCallback: menuItemSelectedCallback,
	}

	w.width = w.getMaxMenuLabelSize() + 14
	w.height = len(menuItems) + 3

	return w
}

func (w *SelectMenuWidget) Name() string {
	return w.name
}

func (w *SelectMenuWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	right := maxX/2 - (w.width / 2)
	top := maxY/2 - (w.height / 2)
	v, err := g.SetView(w.name, right, top, right+w.width, top+w.height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Title = w.title
		v.Frame = true
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.menuItemSelectedAction); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.closeSelectMenuWidget); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.keyArrowDownAction); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.keyArrowUpAction); err != nil {
			return err
		}

		if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
			log.Panicln(err)
		}
	}
	w.RefreshDisplay(g)
	return nil
}

func (w *SelectMenuWidget) RefreshDisplay(g *gocui.Gui) error {
	v, err := g.View(w.name)
	if err != nil {
		return err
	}

	v.Clear()

	fmt.Fprintln(v, " ")
	if len(w.menuItems) == 0 {
		fmt.Fprintln(v, "--empty menu--")
	}
	for i, menuItem := range w.menuItems {
		fmt.Fprintf(v, "    ")
		if w.menuPosition == i {
			fmt.Fprintf(v, util.REVERSE_WHITE)
		}
		fmt.Fprintf(v, "  %v  \n", menuItem.label)
		if w.menuPosition == i {
			fmt.Fprintf(v, util.CLEAR)
		}
	}

	return nil
}

func (w *SelectMenuWidget) getMaxMenuLabelSize() int {
	maxSize := 0
	for _, menuItem := range w.menuItems {
		size := len(menuItem.label)
		if size > maxSize {
			maxSize = size
		}
	}
	return maxSize
}

func (w *SelectMenuWidget) GetMenuSelection() *MenuItem {
	return w.menuItems[w.menuPosition]
}

func (w *SelectMenuWidget) menuItemSelectedAction(g *gocui.Gui, v *gocui.View) error {
	// TODO: Callback
	if w.menuItemSelectedCallback != nil {
		w.menuItemSelectedCallback(g, v, w.GetMenuSelection().id)
	}
	return w.closeSelectMenuWidget(g, v)
}

func (w *SelectMenuWidget) closeSelectMenuWidget(g *gocui.Gui, v *gocui.View) error {
	if err := w.masterUI.CloseView(w); err != nil {
		return err
	}
	return nil
}

func (w *SelectMenuWidget) keyArrowDownAction(g *gocui.Gui, v *gocui.View) error {
	if w.menuPosition+1 < len(w.menuItems) {
		w.menuPosition++
	}
	return w.RefreshDisplay(g)
}

func (w *SelectMenuWidget) keyArrowUpAction(g *gocui.Gui, v *gocui.View) error {
	if w.menuPosition > 0 {
		w.menuPosition--
	}
	return w.RefreshDisplay(g)
}
