package masterUIInterface

import (
	"github.com/jroimartin/gocui"
)

type MasterUIInterface interface {
	SetCurrentViewOnTop(*gocui.Gui, string) error
	GetCurrentView(g *gocui.Gui) *gocui.View
	CloseView(Manager) error
	LayoutManager() LayoutManagerInterface
}

type LayoutManagerInterface interface {
	Contains(Manager) bool
	Add(Manager)
	Remove(Manager) Manager
}

type Manager interface {
	Layout(*gocui.Gui) error
	Name() string
}
