package top

import (
	"errors"
	"fmt"

	"github.com/jroimartin/gocui"
)

type FooterWidget struct {
	name   string
	height int
}

func NewFooterWidget(name string, height int) *FooterWidget {
	return &FooterWidget{name: name, height: height}
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
		v.Title = "Footer"
		fmt.Fprintln(v, "h:help UP/DOWN arrow to highlight row, ENTER to select highlighted row")
		//fmt.Fprint(v, "h:help c:clear q:quit space:refresh o:order p:pause s:sleep f:filter(todo)")
		fmt.Fprint(v, "LEFT/RIGHT arrow to scroll columns h:help")
	}
	return nil
}
