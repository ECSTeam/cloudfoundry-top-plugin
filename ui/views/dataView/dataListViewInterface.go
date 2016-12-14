package dataView

import (
	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/eventdata"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"
)

type DataListViewInterface interface {
	Name() string
	SetTitle(title string)
	GetMargins() (int, int)
	GetMasterUI() masterUIInterface.MasterUIInterface
	GetParentView() DataListViewInterface
	GetDetailView() DataListViewInterface
	SetDetailView(detailView DataListViewInterface)
	GetListWidget() *uiCommon.ListWidget
	GetEventProcessor() *eventdata.EventProcessor
	Layout(g *gocui.Gui) error
	GetDisplayPaused() bool
	SetDisplayPaused(paused bool)
	GetCurrentEventData() *eventdata.EventData
	GetDisplayedEventData() *eventdata.EventData
	RefreshDisplay(g *gocui.Gui) error
	UpdateDisplay(g *gocui.Gui) error
}
