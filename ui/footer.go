package ui

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/jroimartin/gocui"
)

type FooterWidget struct {
	name            string
	height          int
	formatTextRegex *regexp.Regexp
}

func NewFooterWidget(name string, height int) *FooterWidget {
	w := &FooterWidget{name: name, height: height}
	w.formatTextRegex = regexp.MustCompile(`\*\*([^\*]*)*\*\*`)
	return w
}

func (w *FooterWidget) Name() string {
	return w.name
}

func (w *FooterWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(w.name, 0, maxY-w.height, maxX-1, maxY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Frame = false
		w.quickHelp(g, v)
	}
	return nil
}

func (w *FooterWidget) quickHelp(g *gocui.Gui, v *gocui.View) error {

	fmt.Fprint(v, w.formatText("**d**:display "))
	fmt.Fprint(v, w.formatText("**q**:quit "))

	fmt.Fprint(v, w.formatText("**x**:exit detail view "))
	fmt.Fprint(v, w.formatText("**h**:help "))
	fmt.Fprintln(v, w.formatText("**UP**/**DOWN** arrow to highlight row"))
	fmt.Fprint(v, w.formatText("**ENTER** to select highlighted row, "))
	fmt.Fprint(v, w.formatText(`**LEFT**/**RIGHT** arrow to scroll columns`))
	return nil
}

func (w *FooterWidget) formatText(text string) string {
	return w.formatTextRegex.ReplaceAllString(text, "\033[37;1m${1}\033[0m")
}
