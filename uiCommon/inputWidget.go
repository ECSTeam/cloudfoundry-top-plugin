package uiCommon

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
)

type applyCallbackFunc func(g *gocui.Gui, v *gocui.View, w masterUIInterface.Manager, inputValue string) error
type cancelCallbackFunc func(g *gocui.Gui, v *gocui.View) error

type Label struct {
	parentUI         masterUIInterface.Manager
	name             string
	offsetX, offsetY int
	width, height    int
	labelText        string
}

func NewLabel(parentUI masterUIInterface.Manager, name string, offsetX, offsetY int, labelText string) *Label {
	//lines := strings.Split(body, "\n")
	width := len(labelText) + 1
	return &Label{parentUI: parentUI, name: name,
		offsetX: offsetX, offsetY: offsetY, width: width, height: 3, labelText: labelText}
}

func (l *Label) Name() string {
	return l.name
}

func (l *Label) Layout(g *gocui.Gui) error {

	l.parentUI.Layout(g)
	x0, y0, _, _, err := g.ViewPosition(l.parentUI.Name())
	if err != nil {
		return errors.New(l.name + " layout error:" + err.Error())
	}

	v, err := g.SetView(l.name, x0+l.offsetX, y0+l.offsetY, x0+l.offsetX+l.width, y0+l.offsetY+l.height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(l.name + " layout error:" + err.Error())
		}
		v.Frame = false
		fmt.Fprint(v, l.labelText)
	}
	return nil
}

type Input struct {
	parentUI           masterUIInterface.Manager
	name               string
	offsetX, offsetY   int
	width, height      int
	maxLength          int
	inputValue         string
	baseEditor         gocui.Editor
	applyCallbackFunc  applyCallbackFunc
	cancelCallbackFunc cancelCallbackFunc
}

func NewInput(
	parentUI masterUIInterface.Manager,
	name string,
	offsetX, offsetY,
	width,
	maxLength int,
	inputValue string,
	applyCallbackFunc applyCallbackFunc,
	cancelCallbackFunc cancelCallbackFunc,
) *Input {
	return &Input{parentUI: parentUI, name: name, offsetX: offsetX,
		offsetY: offsetY, width: width, height: 2, maxLength: maxLength,
		inputValue:         inputValue,
		baseEditor:         gocui.DefaultEditor,
		applyCallbackFunc:  applyCallbackFunc,
		cancelCallbackFunc: cancelCallbackFunc,
	}
}

func (i *Input) Name() string {
	return i.name
}

func (i *Input) Layout(g *gocui.Gui) error {

	i.parentUI.Layout(g)
	x0, y0, _, _, err := g.ViewPosition(i.parentUI.Name())
	if err != nil {
		return err
	}

	v, err := g.SetView(i.name, x0+i.offsetX, y0+i.offsetY, x0+i.offsetX+i.width, y0+i.offsetY+i.height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editor = i
		v.Editable = true
		v.BgColor = gocui.ColorWhite
		v.FgColor = gocui.ColorBlack
		v.Frame = false
		g.Cursor = true

		fmt.Fprintf(v, i.inputValue)
		//i.baseEditor.MoveCursor(1, 0, false)
		//i.Edit(v, gocui.KeyArrowRight, nil, nil)
		v.MoveCursor(len(i.inputValue), 0, true)
		//v.SetOrigin(2,0)

		if err := g.SetKeybinding(i.name, gocui.KeyEnter, gocui.ModNone, i.applyValueAction); err != nil {
			return err
		}
		if err := g.SetKeybinding(i.name, gocui.KeyEsc, gocui.ModNone, i.cancelValueAction); err != nil {
			return err
		}

	}
	return nil
}

func (w *Input) cancelValueAction(g *gocui.Gui, v *gocui.View) error {
	return w.cancelCallbackFunc(g, v)
}

func (w *Input) applyValueAction(g *gocui.Gui, v *gocui.View) error {
	err := w.applyCallbackFunc(g, v, w.parentUI, w.getValue(v))
	if err != nil {
		// TODO: Display error
	}
	return nil
}

func (w *Input) getValue(v *gocui.View) string {
	lineValue, _ := v.Line(0)
	lineValue = strings.Replace(lineValue, "\x00", "", -1)
	return lineValue
}

func (w *Input) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {

	lineValue := w.getValue(v)
	currentSize := len(lineValue)
	atMax := false
	if currentSize >= w.maxLength {
		atMax = true
	}
	//fmt.Printf("\n[%v], currentSize:%v atMax:%v ", lineValue, currentSize, atMax)

	switch {
	case key == gocui.KeyArrowRight:
		x, _ := v.Cursor()
		//fmt.Printf("\nx:%v", x)
		if x < currentSize {
			w.baseEditor.Edit(v, key, ch, mod)
		}
	case key == gocui.KeyArrowLeft:
		fallthrough
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		fallthrough
	case key == gocui.KeyDelete:
		fallthrough
	case key == gocui.KeyInsert:
		w.baseEditor.Edit(v, key, ch, mod)
	case key == gocui.KeyEnter:
	//	v.EditNewLine()
	case key == gocui.KeyArrowDown:
	//	v.MoveCursor(0, 1, false)
	case key == gocui.KeyArrowUp:
		//	v.MoveCursor(0, -1, false)
	default:
		if !atMax {
			//fmt.Printf("key:[%v] ch:[%v]", key, ch )
			w.baseEditor.Edit(v, key, ch, mod)
		}
	}

}
