package helpView

import (
	"fmt"
  "log"
  "github.com/jroimartin/gocui"
  "github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
)

type HelpView struct {
  masterUI masterUIInterface.MasterUIInterface
	name string
  width int
  height int
  helpText string
}

func NewHelpView(masterUI masterUIInterface.MasterUIInterface, name string, width, height int, helpText string) *HelpView {
	return &HelpView{masterUI: masterUI, name: name, width: width, height: height, helpText: helpText}
}

func (w *HelpView) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
	v, err := g.SetView(w.name, maxX/2-(w.width/2), maxY/2-(w.height/2), maxX/2+(w.width/2), maxY/2+(w.height/2))
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
    v.Title = "Help (press ENTER to close)"
    v.Frame = true
    fmt.Fprintln(v, w.helpText)
    if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeHelpView); err != nil {
      return err
    }

    if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
      log.Panicln(err)
    }

	}
	return nil
}

func (w *HelpView) closeHelpView(g *gocui.Gui, v *gocui.View) error {
  if err := w.masterUI.CloseView(w, w.name); err != nil {
    return err
  }
	return nil
}
