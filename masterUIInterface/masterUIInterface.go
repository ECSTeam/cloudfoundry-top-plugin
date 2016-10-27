package masterUIInterface

import (
	"github.com/jroimartin/gocui"
)
type MasterUIInterface interface {
	SetCurrentViewOnTop(*gocui.Gui, string) error
	CloseView(gocui.Manager, string ) error
	LayoutManager() LayoutManagerInterface
}

type LayoutManagerInterface interface {
	Contains(gocui.Manager) bool
	Add(gocui.Manager)
	Remove(gocui.Manager)
}
