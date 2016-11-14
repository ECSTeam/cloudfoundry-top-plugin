package masterUIInterface

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type EditFilterView struct {
	*EditColumnViewAbs
}

func NewEditFilterView(masterUI MasterUIInterface, name string, listWidget *ListWidget) *EditFilterView {
	w := &EditFilterView{EditColumnViewAbs: NewEditColumnViewAbs(masterUI, name, listWidget)}
	w.width = 55
	w.height = 14
	w.title = "Edit Filter"

	w.refreshDisplayCallbackFunc = func(g *gocui.Gui, v *gocui.View) error {
		return w.refreshDisplayCallback(g, v)
	}

	w.initialLayoutCallbackFunc = func(g *gocui.Gui, v *gocui.View) error {
		return w.initialLayoutCallback(g, v)
	}

	w.applyActionCallbackFunc = func(g *gocui.Gui, v *gocui.View) error {
		return w.applyActionCallback(g, v)
	}

	w.cancelActionCallbackFunc = func(g *gocui.Gui, v *gocui.View) error {
		return w.cancelActionCallback(g, v)
	}

	return w
}

func (w *EditFilterView) initialLayoutCallback(g *gocui.Gui, v *gocui.View) error {

	if err := g.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone, w.keySpaceAction); err != nil {
		return err
	}

	return nil
}

func (w *EditFilterView) refreshDisplayCallback(g *gocui.Gui, v *gocui.View) error {

	v.Clear()
	fmt.Fprintln(v, " ")
	fmt.Fprintln(v, "  RIGHT or LEFT arrow - select column")
	fmt.Fprintln(v, "  SPACE - select column to edit")
	fmt.Fprintln(v, "  ENTER - apply filter")
	fmt.Fprintln(v, "")

	return nil
}

func (w *EditFilterView) keySpaceAction(g *gocui.Gui, v *gocui.View) error {

	return nil
}

func (w *EditFilterView) applyActionCallback(g *gocui.Gui, v *gocui.View) error {

	return nil
}

func (w *EditFilterView) cancelActionCallback(g *gocui.Gui, v *gocui.View) error {

	return nil
}
