package top

import (
	//"fmt"
  "github.com/jroimartin/gocui"
)

type HeaderWidget struct {
  masterUI *MasterUI
	name string
  height int
}

func NewHeaderWidget(masterUI *MasterUI, name string, height int) *HeaderWidget {
	return &HeaderWidget{masterUI: masterUI, name: name, height: height}
}

func (w *HeaderWidget) Layout(g *gocui.Gui) error {
  maxX, _ := g.Size()
	_, err := g.SetView(w.name, 0, 0, maxX-1, w.height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//fmt.Fprint(v, w.body)
	}
	return nil
}
