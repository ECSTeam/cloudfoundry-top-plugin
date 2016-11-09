package top

import (
	"fmt"
  //"strings"
  //"log"
  "github.com/jroimartin/gocui"
  "github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
)

type EditIntervalWidget struct {
  masterUI masterUIInterface.MasterUIInterface
	name string
  width int
  height int
  titleText string
  helpText string

  labelWidget masterUIInterface.Manager
  inputWidget masterUIInterface.Manager

}

func NewEditIntervalWidget(masterUI masterUIInterface.MasterUIInterface, name string, width, height int) *EditIntervalWidget {


  w := &EditIntervalWidget{masterUI: masterUI, name: name, width: width, height: height}

  labelText := "Seconds:"
  maxLength := 10
  w.titleText = "Update refresh interval"
  w.helpText = "no help"

  w.labelWidget = masterUIInterface.NewLabel(w, "label", 1, 2, labelText)

  applyCallbackFunc := func(g *gocui.Gui, v *gocui.View, inputValue string) error {
    fmt.Printf("\n**** ENTER: [%v] ****\n", inputValue)
    return w.closeWidget(g, v)
  }
  cancelCallbackFunc := func(g *gocui.Gui, v *gocui.View) error {
    fmt.Printf("\n**** CANCELED ****\n")
    return w.closeWidget(g, v)
  }
  inputValue := "test"
  w.inputWidget = masterUIInterface.NewInput(w, "input", len(labelText)+2, 2, maxLength+2,
      maxLength, inputValue, applyCallbackFunc,cancelCallbackFunc)

  return w
}

func (w *EditIntervalWidget) Name() string {
  return w.name
}

func (w *EditIntervalWidget) Init(g *gocui.Gui) error {
  w.masterUI.LayoutManager().Add(w)
  w.masterUI.LayoutManager().Add(w.labelWidget)
  w.masterUI.LayoutManager().Add(w.inputWidget)
  w.Layout(g)
  w.labelWidget.Layout(g)
  w.inputWidget.Layout(g)
  return w.masterUI.SetCurrentViewOnTop(g,"input")
}

func (w *EditIntervalWidget) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
	v, err := g.SetView(w.name, maxX/2-(w.width/2), maxY/2-(w.height/2), maxX/2+(w.width/2), maxY/2+(w.height/2))
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
    v.Title = ""
    v.Frame = true
    fmt.Fprintf(v, " %v", w.titleText)


    if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.closeWidget); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, 'q', gocui.ModNone, w.closeWidget); err != nil {
      return err
    }
    /*
    if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
      log.Panicln(err)
    }
    */
    //return w.masterUI.SetCurrentViewOnTop(g,"input")
	}
	return nil
}

func (w *EditIntervalWidget) closeWidget(g *gocui.Gui, v *gocui.View) error {

  if err := w.masterUI.CloseView(w.labelWidget); err != nil {
    return err
  }

  if err := w.masterUI.CloseView(w.inputWidget); err != nil {
    return err
  }

  if err := w.masterUI.CloseView(w); err != nil {
    return err
  }
	return nil
}
