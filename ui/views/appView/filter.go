package appView

import (
	"errors"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
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
