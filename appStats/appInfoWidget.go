package appStats

import (
	"fmt"
  //"strings"
  "log"
  "errors"
  "github.com/jroimartin/gocui"
  "github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
)

type AppInfoWidget struct {
  masterUI masterUIInterface.MasterUIInterface
	name string
  width int
  height int
}

func NewAppInfoWidget(masterUI masterUIInterface.MasterUIInterface, name string, width, height int) *AppInfoWidget {
	return &AppInfoWidget{masterUI: masterUI, name: name, width: width, height: height}
}

func (w *AppInfoWidget) Name() string {
  return w.name
}

func (w *AppInfoWidget) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
	v, err := g.SetView(w.name, maxX/2-(w.width/2), maxY/2-(w.height/2), maxX/2+(w.width/2), maxY/2+(w.height/2))
	if err != nil {
		if err != gocui.ErrUnknownView {
      return errors.New(w.name+" layout error:" + err.Error())
		}
    v.Title = "App Info (press ENTER to close)"
    v.Frame = true
    fmt.Fprintln(v, "Future home of app info")
    if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeAppInfoWidget); err != nil {
      return err
    }

    if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
      log.Panicln(err)
    }

	}
	return nil
}

func (w *AppInfoWidget) closeAppInfoWidget(g *gocui.Gui, v *gocui.View) error {
  if err := w.masterUI.CloseView(w); err != nil {
    return err
  }
	return nil
}
