package debug

import (
	"fmt"
  "time"
  "github.com/jroimartin/gocui"
)

var (
	 gui *gocui.Gui
   //view *gocui.View
)

func Init(g *gocui.Gui) {
  gui = g
}


func Debug(msg string)  {

  gui.Execute(func(gui *gocui.Gui) error {
    openView()
    v, err := gui.View("debugView")
    if err != nil {
      return err
    }
    fmt.Fprintf(v, "%v %v\n", time.Now().Format("15:04:05"), msg)
		return nil
	})


}

func closeView(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView("debugView")
	return nil
}


func openView() error {

  debugViewName := "debugView"
  maxX, maxY := gui.Size()
  left := 5
  right := maxX - 5
  top := 4
  bottom := maxY - 2

  if v, err := gui.SetView(debugViewName, left, top, right, bottom); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "DEBUG"
			v.Frame = true
      v.Autoscroll = true
      v.Wrap = true
      v.BgColor = gocui.ColorRed
			//fmt.Fprintln(v, "Debug Window")
      if err := gui.SetKeybinding(debugViewName, gocui.KeyEnter, gocui.ModNone, closeView); err != nil {
        return err
      }
      if _, err := gui.SetCurrentView(debugViewName); err != nil {
    		return err
    	}
      if _, err :=  gui.SetViewOnTop(debugViewName); err != nil {
        return err
      }
	}
  return nil
}
