package helpView

import (
	"fmt"
  "log"
  "strings"
  "regexp"
  "github.com/jroimartin/gocui"
  "github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
)

type HelpView struct {
  masterUI masterUIInterface.MasterUIInterface
	name string
  width int
  height int
  helpText string
  displayText string
  helpTextLines int

  viewOffset int
}

func NewHelpView(masterUI masterUIInterface.MasterUIInterface, name string, width, height int, helpText string) *HelpView {
	hv := &HelpView{masterUI: masterUI, name: name, width: width, height: height, helpText: helpText}
  hv.helpTextLines = strings.Count(helpText,"\n")
  return hv
}

func (w *HelpView) Name() string {
  return w.name
}

func (w *HelpView) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
	v, err := g.SetView(w.name, maxX/2-(w.width/2), maxY/2-(w.height/2), maxX/2+(w.width/2), maxY/2+(w.height/2))
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
    v.Title = "Help (press ENTER to close, DOWN/UP arrow to scroll)"
    v.Frame = true

    if w.displayText == "" {
      re := regexp.MustCompile(`\*\*(.*)\*\*`)
      w.displayText = re.ReplaceAllString(w.helpText, "\033[37;1m${1}\033[0m")
    }

    fmt.Fprintf(v, w.displayText)
    if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeHelpView); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.closeHelpView); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, 'q', gocui.ModNone, w.closeHelpView); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.arrowUp); err != nil {
      log.Panicln(err)
    }
    if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.arrowDown); err != nil {
      log.Panicln(err)
    }

    if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
      log.Panicln(err)
    }

	}
	return nil
}

func (w *HelpView) closeHelpView(g *gocui.Gui, v *gocui.View) error {
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
