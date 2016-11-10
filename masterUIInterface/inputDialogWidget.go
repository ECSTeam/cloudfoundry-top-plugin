package masterUIInterface

import (
	"fmt"
  "errors"
  //"strings"
  //"log"
  "github.com/jroimartin/gocui"
)

type InputDialogWidget struct {
  masterUI MasterUIInterface
	name string
  width int
  height int
  titleText string
  helpText string

  labelWidget Manager
  inputWidget Manager

}

func NewInputDialogWidget(masterUI MasterUIInterface,
    name string,
    width, height int,
    labelText string, maxLength int,
    titleText string, helpText string,
    valueText string,
    applyValueCallback applyCallbackFunc) *InputDialogWidget {


  w := &InputDialogWidget{
    masterUI: masterUI,
    name: name,
    width: width,
    height: height,
    titleText: titleText,
    helpText: helpText,
  }


  w.labelWidget = NewLabel(w, "label", 1, 2, labelText)

  /*
  applyCallbackFunc := func(g *gocui.Gui, v *gocui.View, w InputDialogWidget, inputValue string) error {
    fmt.Printf("\n**** ENTER: [%v] ****\n", inputValue)
    return w.closeWidget(g, v)
  }
  */

  cancelCallbackFunc := func(g *gocui.Gui, v *gocui.View) error {
    //fmt.Printf("\n**** CANCELED ****\n")
    return w.CloseWidget(g, v)
  }
  w.inputWidget = NewInput(w, "input", len(labelText)+2, 2, maxLength+2,
      maxLength, valueText, applyValueCallback, cancelCallbackFunc)

  return w
}

func (w *InputDialogWidget) Name() string {
  return w.name
}

func (w *InputDialogWidget) Init(g *gocui.Gui) error {
  w.masterUI.LayoutManager().Add(w)
  w.masterUI.LayoutManager().Add(w.labelWidget)
  w.masterUI.LayoutManager().Add(w.inputWidget)
  w.Layout(g)
  w.labelWidget.Layout(g)
  w.inputWidget.Layout(g)
  return w.masterUI.SetCurrentViewOnTop(g,"input")
}

func (w *InputDialogWidget) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
	v, err := g.SetView(w.name, maxX/2-(w.width/2), maxY/2-(w.height/2), maxX/2+(w.width/2), maxY/2+(w.height/2))
	if err != nil {
		if err != gocui.ErrUnknownView {
      return errors.New(w.name+" layout error:" + err.Error())
		}
    v.Title = ""
    v.Frame = true
    fmt.Fprintf(v, " %v", w.titleText)


    if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.CloseWidget); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, 'q', gocui.ModNone, w.CloseWidget); err != nil {
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

func (w *InputDialogWidget) CloseWidget(g *gocui.Gui, v *gocui.View) error {
  g.Cursor = false

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
