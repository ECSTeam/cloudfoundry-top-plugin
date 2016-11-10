package debug

import (
	"fmt"
  "log"
  "time"
  "strings"
  "github.com/jroimartin/gocui"
  "github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
)


var (
	 gui *gocui.Gui
   debugWidget *DebugWidget
)

func Debug(msg string)  {

  if gui != nil {
    gui.Execute(func(gui *gocui.Gui) error {
      msg := strings.Replace(msg, "\n"," | ",-1)

      line := fmt.Sprintf("%v %v", time.Now().Format("15:04:05"), msg)
      debugWidget.debugLines = append(debugWidget.debugLines, line)

      _, maxY := gui.Size()
      top := 4
      bottom := maxY - 2
      height := bottom - top - 1
      //height = 10
      viewOffset := len(debugWidget.debugLines) - height
      if (viewOffset < 0) {
        viewOffset = 0
      }
      debugWidget.viewOffset = viewOffset

      openView()
      return nil
  	})
  }
}

type DebugWidget struct {
  masterUI masterUIInterface.MasterUIInterface
	name string
  height int
  width int
  viewOffset int
  horizonalOffset int
  debugLines []string
}

func InitDebug(g *gocui.Gui, masterUI masterUIInterface.MasterUIInterface) {
  debugWidget = NewDebugWidget(masterUI, "debugView")
  gui = g
}

func openView() {
  layoutMgr := debugWidget.masterUI.LayoutManager()
  if layoutMgr.Top() != debugWidget {
    layoutMgr.Add(debugWidget)
  }
  debugWidget.Layout(gui)
}

func NewDebugWidget(masterUI masterUIInterface.MasterUIInterface, name string) *DebugWidget {
	hv := &DebugWidget{masterUI: masterUI, name: name}
  return hv
}

func (w *DebugWidget) Name() string {
  return w.name
}

func (w *DebugWidget) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
  left := 5
  right := maxX - 5
  top := 4
  bottom := maxY - 2
  w.height = bottom - top - 1
  w.width = right - left

	v, err := g.SetView(w.name, left, top, right, bottom)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
    v.Title = "DEBUG (press ENTER to close, DOWN/UP arrow to scroll)"
    v.Frame = true
    v.Autoscroll = false
    v.Wrap = false
    v.BgColor = gocui.ColorRed
    g.SelBgColor = gocui.ColorRed
    g.Highlight = true

    fmt.Fprintf(v, "Debug window\n")

    if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeDebugWidget); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.closeDebugWidget); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, 'q', gocui.ModNone, w.closeDebugWidget); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.arrowUp); err != nil {
      log.Panicln(err)
    }
    if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.arrowDown); err != nil {
      log.Panicln(err)
    }
    if err := g.SetKeybinding(w.name, gocui.KeyArrowRight, gocui.ModNone, w.arrowRight); err != nil {
      log.Panicln(err)
    }
    if err := g.SetKeybinding(w.name, gocui.KeyArrowLeft, gocui.ModNone, w.arrowLeft); err != nil {
      log.Panicln(err)
    }
    if err := g.SetKeybinding(w.name, 'x', gocui.ModNone, w.testMsg); err != nil {
      log.Panicln(err)
    }

    if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
      log.Panicln(err)
    }
	} else {
    v.Clear()
    h := w.height
    for index:=w.viewOffset; (index-w.viewOffset)<(h) && index < len(w.debugLines);index++ {
      debugLine := w.debugLines[index]
      if w.horizonalOffset < len(debugLine) {
        debugLine = debugLine[w.horizonalOffset:len(debugLine)]
      } else {
        debugLine = ""
      }
      line := fmt.Sprintf("[%03v] %v\n", index, debugLine)
      fmt.Fprintf(v, line)
    }

  }

	return nil
}

func (w *DebugWidget) closeDebugWidget(g *gocui.Gui, v *gocui.View) error {
  g.Highlight = false
  if err := w.masterUI.CloseView(w); err != nil {
    return err
  }
	return nil
}

func (w *DebugWidget) testMsg(g *gocui.Gui, v *gocui.View) error {
  Debug("hello")
	return nil
}


func (w *DebugWidget) arrowRight(g *gocui.Gui, v *gocui.View) error {
  w.horizonalOffset = w.horizonalOffset + 5
	return nil
}


func (w *DebugWidget) arrowLeft(g *gocui.Gui, v *gocui.View) error {
  w.horizonalOffset = w.horizonalOffset - 5
  if w.horizonalOffset < 0 {
    w.horizonalOffset = 0
  }
	return nil
}


func (w *DebugWidget) arrowUp(g *gocui.Gui, v *gocui.View) error {
  if w.viewOffset > 0 {
    w.viewOffset--
  }
	return nil
}

func (w *DebugWidget) arrowDown(g *gocui.Gui, v *gocui.View) error {
  h := w.height
  if w.viewOffset < len(w.debugLines) && (len(w.debugLines) - h) > w.viewOffset {
    w.viewOffset++
  }
	return nil
}
