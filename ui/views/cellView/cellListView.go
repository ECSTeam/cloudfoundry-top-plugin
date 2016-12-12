package cellView

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

type CellListView struct {
	masterUI       masterUIInterface.MasterUIInterface
	name           string
	topMargin      int
	bottomMargin   int
	eventProcessor *eventdata.EventProcessor
	mu             sync.Mutex
	listWidget     *uiCommon.ListWidget
	displayPaused  bool
}

func NewCellListView(masterUI masterUIInterface.MasterUIInterface,
	name string, topMargin, bottomMargin int,
	eventProcessor *eventdata.EventProcessor) *CellListView {

	return &CellListView{
		masterUI:       masterUI,
		name:           name,
		topMargin:      topMargin,
		bottomMargin:   bottomMargin,
		eventProcessor: eventProcessor,
	}
}

func (asUI *CellListView) Name() string {
	return asUI.name
}

func (asUI *CellListView) Layout(g *gocui.Gui) error {

	if asUI.listWidget == nil {

		statList := asUI.postProcessData(asUI.GetDisplayedEventData().AppMap)
		listData := asUI.convertToListData(statList)

		listWidget := uiCommon.NewListWidget(asUI.masterUI, asUI.name,
			asUI.topMargin, asUI.bottomMargin, asUI, asUI.columnDefinitions(),
			listData)
		listWidget.Title = "Cell List"
		listWidget.PreRowDisplayFunc = asUI.PreRowDisplay

		defaultSortColums := []*uiCommon.SortColumn{
			uiCommon.NewSortColumn("colA", true),
			uiCommon.NewSortColumn("colB", true),
		}
		listWidget.SetSortColumns(defaultSortColums)

		asUI.listWidget = listWidget
		if err := g.SetKeybinding(asUI.name, 'h', gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				helpView := helpView.NewHelpView(asUI.masterUI, "helpView", 75, 17, helpText)
				asUI.masterUI.LayoutManager().Add(helpView)
				asUI.masterUI.SetCurrentViewOnTop(g)
				return nil
			}); err != nil {
			log.Panicln(err)
		}

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

func (asUI *CellListView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, asUI.columnA())
	columns = append(columns, asUI.columnB())

	return columns
}

func (asUI *CellListView) GetDisplayPaused() bool {
	return asUI.displayPaused
}

func (asUI *CellListView) SetDisplayPaused(paused bool) {
	asUI.displayPaused = paused
	if !paused {
		asUI.updateData()
	}
}

func (asUI *CellListView) GetCurrentEventData() *eventdata.EventData {
	return asUI.eventProcessor.GetCurrentEventData()
}

func (asUI *CellListView) GetDisplayedEventData() *eventdata.EventData {
	return asUI.eventProcessor.GetDisplayedEventData()
}

func (asUI *CellListView) postProcessData(statsMap map[string]*eventdata.AppStats) []*eventdata.AppStats {
	if len(statsMap) > 0 {
		stats := eventdata.PopulateNamesIfNeeded(statsMap)
		return stats
	} else {
		return nil
	}
}

func (asUI *CellListView) convertToListData(statsList []*eventdata.AppStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, len(statsList))
	for i, d := range statsList {
		listData[i] = d
	}
	return listData
}

func (asUI *CellListView) RefreshDisplay(g *gocui.Gui) error {
	err := asUI.refreshListDisplay(g)
	if err != nil {
		return err
	}
	return asUI.updateHeader(g)
}

func (asUI *CellListView) refreshListDisplay(g *gocui.Gui) error {
	err := asUI.listWidget.RefreshDisplay(g)
	if err != nil {
		return err
	}
	return err
}

func (asUI *CellListView) UpdateDisplay(g *gocui.Gui) error {
	if !asUI.displayPaused {
		asUI.updateData()
	}
	return asUI.RefreshDisplay(g)
}

// XXX
func (asUI *CellListView) updateData() {
	asUI.eventProcessor.UpdateData()
	processor := asUI.GetDisplayedEventData()
	statList := asUI.postProcessData(processor.AppMap)
	listData := asUI.convertToListData(statList)
	asUI.listWidget.SetListData(listData)
}

func (asUI *CellListView) PreRowDisplay(data uiCommon.IData, isSelected bool) string {
	return ""
}

func (asUI *CellListView) updateHeader(g *gocui.Gui) error {

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
	return nil
}
