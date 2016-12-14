package appDetailView

import (
	"errors"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
)

type RequestsInfoWidget struct {
	masterUI masterUIInterface.MasterUIInterface
	name     string
}

func NewRequestsInfoWidget(masterUI masterUIInterface.MasterUIInterface, name string) *RequestsInfoWidget {
	return &RequestsInfoWidget{masterUI: masterUI, name: name}
}

func (w *RequestsInfoWidget) Name() string {
	return w.name
}

func (w *RequestsInfoWidget) Layout(g *gocui.Gui) error {
	maxX, _ := g.Size()

	topMargin := 7
	width := maxX - 1
	height := 15

	//v, err := g.SetView(w.name, maxX/2-(w.width/2), maxY/2-(w.height/2), maxX/2+(w.width/2), maxY/2+(w.height/2))
	v, err := g.SetView(w.name, 0, topMargin, width, topMargin+height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Title = "App Request Info"
		v.Frame = true
		fmt.Fprintln(v, "Future home of app info")
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeRequestsInfoWidget); err != nil {
			return err
		}

		if err := w.masterUI.SetCurrentViewOnTop(g); err != nil {
			log.Panicln(err)
		}

	}
	return nil
}

func (w *RequestsInfoWidget) closeRequestsInfoWidget(g *gocui.Gui, v *gocui.View) error {
	if err := w.masterUI.CloseView(w); err != nil {
		return err
	}
	return nil
}
