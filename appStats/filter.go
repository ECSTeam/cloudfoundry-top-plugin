package appStats

import (
	"fmt"
  //"strings"
  "log"
  "github.com/jroimartin/gocui"
  "github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
)

type FilterWidget struct {
  masterUI masterUIInterface.MasterUIInterface
	name string
  width int
  height int
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
			return err
		}
    v.Title = "Filter (press ENTER to close)"
    v.Frame = true
    fmt.Fprintln(v, "Future home of filter screen")
    if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeFilterWidget); err != nil {
      return err
    }

    if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
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

/*
func (asUI *AppStatsUI) showFilter(g *gocui.Gui, v *gocui.View) error {
	 return asUI.openFilterView(g, v )
}


func (asUI *AppStatsUI) openFilterView(g *gocui.Gui, v *gocui.View) error {

  filterViewName := "filterView"
  maxX, maxY := g.Size()


  if v, err := g.SetView(filterViewName, maxX/2-32, maxY/5, maxX/2+32, maxY/2+5); err != nil {
      if err != gocui.ErrUnknownView {
        return err
      }
			v.Title = "Filter"
			v.Frame = true
      v.Autoscroll = false
      v.Wrap = false

      //v.Highlight = true
      //v.SelBgColor = gocui.ColorGreen
      //v.SelFgColor = gocui.ColorBlack

			fmt.Fprintf(v, "Filter Window\n")
      fmt.Fprintf(v, "Line two\n")

      fieldText := asUI.filterAppName
      fieldLabel := "App Name"
      fmt.Fprintf(v, "\r%v:\033[32;7m%-10.10v\033[0m", fieldLabel, fieldText)
      //fmt.Fprintf(v, "Application name: Hello \033[32;7m      \033[0m")




      if err := g.SetKeybinding(filterViewName, gocui.KeyEnter, gocui.ModNone, asUI.closeFilterView); err != nil {
        return err
      }

      //fieldText := ""
      //fieldLabel := "App Name"
      fieldMax := 8
      fieldLen := 0
      for i := 32 ; i < 127 ; i++ {
          //keyPress := ch
          keyPress := rune(i)
          if err := g.SetKeybinding(filterViewName, keyPress, gocui.ModNone,
            func(g *gocui.Gui, v *gocui.View) error {
                   if (fieldLen >= fieldMax) {
                     return nil
                   }
                   c := string(keyPress)
                   fieldLen++
                   fieldText += c
                   //fmt.Fprintf(v, "%v", c)
                   v.Clear()
                   fmt.Fprintf(v, "Filter Window\n")
                   fmt.Fprintf(v, "Line two\n")

                   fmt.Fprintf(v, "%v:\033[32;7m%-10.10v\033[0m\n", fieldLabel, fieldText)
                   fmt.Fprintf(v, "%v:\033[32;7m%-10.10v\033[0m\n", "Field2", "")
                   asUI.filterAppName = fieldText
			             return nil
		        }); err != nil {
            return err
          }
      }

      if _, err := asUI.setCurrentViewOnTop(g, filterViewName); err != nil {
        return err
      }
	}
  return nil
}

func (asUI *AppStatsUI) closeFilterView(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView("filterView")
  g.DeleteKeybindings("filterView")
  if _, err := asUI.setCurrentViewOnTop(g, "detailView"); err != nil {
    return err
  }
	return nil
}
*/
