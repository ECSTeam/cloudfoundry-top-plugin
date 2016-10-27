package top

import (
	"fmt"
  "github.com/jroimartin/gocui"
)

type FooterWidget struct {
	name string
  height int
}

func NewFooterWidget(name string, height int) *FooterWidget {
	return &FooterWidget{name: name, height: height}
}

func (w *FooterWidget) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
	v, err := g.SetView(w.name, 0, maxY-w.height, maxX-1, maxY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
    v.Frame = false
    v.Title = "Footer"
    fmt.Fprintln(v, "c:clear q:quit space:fresh")
    fmt.Fprintln(v, "s:sleep(todo) f:filter(todo) o:order(todo)")
    fmt.Fprintln(v, "h:help")
	}
	return nil
}
