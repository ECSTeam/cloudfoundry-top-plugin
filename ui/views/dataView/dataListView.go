package dataView

import (
	"fmt"
	"log"
	"sync"

	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/eventdata"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/helpView"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

type actionCallback func(g *gocui.Gui, v *gocui.View) error
type initializeCallback func(g *gocui.Gui, viewName string) error
type preRowDisplayCallback func(data uiCommon.IData, isSelected bool) string
type refreshDisplayCallback func(g *gocui.Gui) error

type GetListData func() []uiCommon.IData

type DataListView struct {
	masterUI               masterUIInterface.MasterUIInterface
	parentView             DataListViewInterface
	detailView             DataListViewInterface
	name                   string
	topMargin              int
	bottomMargin           int
	eventProcessor         *eventdata.EventProcessor
	mu                     sync.Mutex
	listWidget             *uiCommon.ListWidget
	displayPaused          bool
	initialized            bool
	Title                  string
	HelpText               string
	InitializeCallback     initializeCallback
	UpdateHeaderCallback   actionCallback
	PreRowDisplayCallback  preRowDisplayCallback
	RefreshDisplayCallback refreshDisplayCallback
	GetListData            GetListData
}

func NewDataListView(masterUI masterUIInterface.MasterUIInterface,
	parentView DataListViewInterface,
	name string, topMargin, bottomMargin int,
	eventProcessor *eventdata.EventProcessor,
	columnDefinitions []*uiCommon.ListColumn,
	defaultSortColumns []*uiCommon.SortColumn) *DataListView {

	asUI := &DataListView{
		masterUI:       masterUI,
		parentView:     parentView,
		name:           name,
		topMargin:      topMargin,
		bottomMargin:   bottomMargin,
		eventProcessor: eventProcessor,
	}

	listWidget := uiCommon.NewListWidget(asUI.masterUI, asUI.name,
		asUI.topMargin, asUI.bottomMargin, asUI, columnDefinitions)
	listWidget.PreRowDisplayFunc = asUI.PreRowDisplay

	listWidget.SetSortColumns(defaultSortColumns)

	asUI.listWidget = listWidget

	return asUI

}

func (asUI *DataListView) Name() string {
	return asUI.name
}

func (asUI *DataListView) SetTitle(title string) {
	asUI.listWidget.Title = title
}

func (asUI *DataListView) GetMargins() (int, int) {
	return asUI.topMargin, asUI.bottomMargin
}

func (asUI *DataListView) GetMasterUI() masterUIInterface.MasterUIInterface {
	return asUI.masterUI
}

func (asUI *DataListView) GetParentView() DataListViewInterface {
	return asUI.parentView
}

func (asUI *DataListView) GetDetailView() DataListViewInterface {
	return asUI.detailView
}

func (asUI *DataListView) SetDetailView(detailView DataListViewInterface) {
	asUI.detailView = detailView
}

func (asUI *DataListView) GetListWidget() *uiCommon.ListWidget {
	return asUI.listWidget
}

func (asUI *DataListView) GetEventProcessor() *eventdata.EventProcessor {
	return asUI.eventProcessor
}

func (asUI *DataListView) Layout(g *gocui.Gui) error {

	if !asUI.initialized {

		asUI.initialized = true
		asUI.initialize(g)

		// TODO
		/*
			if err := g.SetKeybinding(asUI.name, gocui.KeyEnter, gocui.ModNone,
				func(g *gocui.Gui, v *gocui.View) error {
					if asUI.listWidget.HighlightKey() != "" {
						asUI.appDetailView = NewAppDetailView(asUI.masterUI, "appDetailView", asUI.listWidget.HighlightKey(), asUI)
						asUI.masterUI.LayoutManager().Add(asUI.appDetailView)
						asUI.masterUI.SetCurrentViewOnTop(g)
					}
					return nil
				}); err != nil {
				log.Panicln(err)
			}
		*/

	}
	return asUI.listWidget.Layout(g)
}

func (asUI *DataListView) initialize(g *gocui.Gui) {
	if err := g.SetKeybinding(asUI.name, 'h', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			helpView := helpView.NewHelpView(asUI.masterUI, "helpView", 75, 17, asUI.HelpText)
			asUI.masterUI.LayoutManager().Add(helpView)
			asUI.masterUI.SetCurrentViewOnTop(g)
			return nil
		}); err != nil {
		log.Panicln(err)
	}
	if asUI.InitializeCallback != nil {
		asUI.InitializeCallback(g, asUI.name)
	}
}

func (asUI *DataListView) GetDisplayPaused() bool {
	return asUI.displayPaused
}

func (asUI *DataListView) SetDisplayPaused(paused bool) {
	asUI.displayPaused = paused
	if !paused {
		asUI.updateData()
	}
}

func (asUI *DataListView) GetCurrentEventData() *eventdata.EventData {
	return asUI.eventProcessor.GetCurrentEventData()
}

func (asUI *DataListView) GetDisplayedEventData() *eventdata.EventData {
	return asUI.eventProcessor.GetDisplayedEventData()
}

func (asUI *DataListView) RefreshDisplay(g *gocui.Gui) error {
	var err error

	if asUI.RefreshDisplayCallback != nil {
		err = asUI.RefreshDisplayCallback(g)
	} else {
		err = asUI.refreshListDisplay(g)
	}
	if err != nil {
		return err
	}
	return asUI.updateHeader(g)
}

func (asUI *DataListView) refreshListDisplay(g *gocui.Gui) error {
	err := asUI.listWidget.RefreshDisplay(g)
	if err != nil {
		return err
	}
	return err
}

func (asUI *DataListView) UpdateDisplay(g *gocui.Gui) error {
	if !asUI.displayPaused {
		asUI.updateData()
	}
	return asUI.RefreshDisplay(g)
}

// XXX
func (asUI *DataListView) updateData() {
	asUI.eventProcessor.UpdateData()
	listData := asUI.GetListData()
	asUI.listWidget.SetListData(listData)
}

func (asUI *DataListView) PreRowDisplay(data uiCommon.IData, isSelected bool) string {
	if asUI.PreRowDisplayCallback != nil {
		return asUI.PreRowDisplayCallback(data, isSelected)
	}
	return ""
}

func (asUI *DataListView) updateHeader(g *gocui.Gui) error {

	v, err := g.View("headerView")
	if err != nil {
		return err
	}
	if asUI.displayPaused {
		fmt.Fprintf(v, util.REVERSE_GREEN)
		fmt.Fprintf(v, "\r Display update paused ")
		fmt.Fprintf(v, util.CLEAR)
		return nil
	}

	if asUI.UpdateHeaderCallback != nil {
		return asUI.UpdateHeaderCallback(g, v)
	}
	return nil
}
