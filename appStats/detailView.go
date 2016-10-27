package appStats

import (
	"fmt"
  "log"
  "github.com/jroimartin/gocui"
  "github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
)

type DetailView struct {
  masterUI masterUIInterface.MasterUIInterface
	name string
  topMargin int
  bottomMargin int
}

func NewDetailView(masterUI masterUIInterface.MasterUIInterface,name string, topMargin, bottomMargin int) *DetailView {
	return &DetailView{masterUI: masterUI,name: name, topMargin: topMargin, bottomMargin: bottomMargin}
}

func (w *DetailView) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
  v, err := g.SetView(w.name, 0, w.topMargin, maxX-1, maxY-w.bottomMargin)
  if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
    fmt.Fprintln(v, "")
    filter := NewFilterWidget(w.masterUI, "filterWidget", 30, 10)
    if err := g.SetKeybinding(w.name, 'f', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
         if !w.masterUI.LayoutManager().Contains(filter) {
           w.masterUI.LayoutManager().Add(filter)
         }
         w.masterUI.SetCurrentViewOnTop(g,"filterWidget")
         return nil
    }); err != nil {
      log.Panicln(err)
    }

    if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
      log.Panicln(err)
    }

	}
  return nil
}
