package ui

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ansel1/merry"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/eventdata"
	"github.com/kkellner/cloudfoundry-top-plugin/eventrouting"
	"github.com/kkellner/cloudfoundry-top-plugin/toplog"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/appView"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/cellView"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

const WarmUpSeconds = 60

type MasterUI struct {
	layoutManager *uiCommon.LayoutManager
	gui           *gocui.Gui
	cliConnection plugin.CliConnection

	currentDataView masterUIInterface.UpdatableView

	router            *eventrouting.EventRouter
	refreshNow        chan bool
	refreshIntervalMS time.Duration

	headerSize int
	footerSize int
}

func NewMasterUI(cliConnection plugin.CliConnection) *MasterUI {

	mui := &MasterUI{
		cliConnection: cliConnection,
		refreshNow:    make(chan bool),
	}

	eventProcessor := eventdata.NewEventProcessor(mui.cliConnection)
	mui.router = eventrouting.NewEventRouter(eventProcessor)

	return mui
}

func (mui *MasterUI) CliConnection() plugin.CliConnection {
	return mui.cliConnection
}

func (mui *MasterUI) LayoutManager() masterUIInterface.LayoutManagerInterface {
	return mui.layoutManager
}
func (mui *MasterUI) GetRouter() *eventrouting.EventRouter {
	return mui.router
}

func (mui *MasterUI) Start() {
	mui.router.GetProcessor().Start()
	mui.initGui()
}

func (mui *MasterUI) initGui() {

	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	mui.gui = g
	g.InputEsc = true
	defer g.Close()

	mui.layoutManager = uiCommon.NewLayoutManager()
	g.SetManager(mui.layoutManager)

	toplog.InitDebug(g, mui)

	mui.footerSize = 4
	footerView := NewFooterWidget("footerView", mui.footerSize)
	mui.layoutManager.Add(footerView)

	mui.headerSize = 6
	headerView := NewHeaderWidget(mui, "headerView", mui.headerSize)
	mui.layoutManager.Add(headerView)
	// We add the common keybindings to the header view in the event
	// that no DataView is open
	mui.addCommonDataViewKeybindings(g, "headerView")

	mui.openView(g, "appListView")

	// default refresh to 1 second
	mui.refreshIntervalMS = 1000 * time.Millisecond

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, mui.quit); err != nil {
		log.Panicln(err)
	}

	go mui.counter(g)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		m := merry.Details(err)
		log.Panicln(m)
	}

}

func (mui *MasterUI) addCommonDataViewKeybindings(g *gocui.Gui, viewName string) error {
	if err := g.SetKeybinding(viewName, 'q', gocui.ModNone, mui.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, 'C', gocui.ModNone, mui.clearStats); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, gocui.KeySpace, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		mui.RefeshNow()
		return nil
	}); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, 's', gocui.ModNone, mui.editUpdateInterval); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, 'd', gocui.ModNone, mui.selectDisplayAction); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, 'r', gocui.ModNone, mui.refreshMetadata); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(viewName, 'D', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			toplog.Open()
			return nil
		}); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (mui *MasterUI) CloseView(m masterUIInterface.Manager) error {

	mui.gui.DeleteView(m.Name())
	mui.gui.DeleteKeybindings(m.Name())
	nextForFocus := mui.layoutManager.Remove(m)
	nextViewName := nextForFocus.Name()
	if err := mui.SetCurrentViewOnTop(mui.gui); err != nil {
		return merry.Wrap(err).Appendf("SetCurrentViewOnTop viewName:[%v]", nextViewName)
	}
	return nil
}

func (mui *MasterUI) CloseViewByName(viewName string) error {
	m := mui.layoutManager.GetManagerByViewName(viewName)
	return mui.CloseView(m)
}

func (mui *MasterUI) SetCurrentViewOnTop(g *gocui.Gui) error {
	topName := mui.layoutManager.Top().Name()
	if _, err := g.SetCurrentView(topName); err != nil {
		return merry.Wrap(err).Appendf("SetCurrentView viewName:[%v]", topName)
	}
	if _, err := g.SetViewOnTop(topName); err != nil {
		return merry.Wrap(err).Appendf("SetViewOnTop viewName:[%v]", topName)
	}
	return nil
}

func (mui *MasterUI) GetCurrentView(g *gocui.Gui) *gocui.View {
	views := g.Views()
	view := views[len(views)-1]
	return view
}

func (mui *MasterUI) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (mui *MasterUI) selectDisplayAction(g *gocui.Gui, v *gocui.View) error {

	menuItems := make([]*uiCommon.MenuItem, 0, 5)
	menuItems = append(menuItems, uiCommon.NewMenuItem("appListView", "App Stats"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("cellListView", "Cell Stats"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("eventstats", "TODO: Event Stats"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("eventhistory", "TODO: Event Rate History"))
	selectDisplayView := uiCommon.NewSelectMenuWidget(mui, "selectDisplayView", "Select Display", menuItems, mui.selectDisplayCallback)
	mui.LayoutManager().Add(selectDisplayView)
	mui.SetCurrentViewOnTop(g)
	return nil
}

func (mui *MasterUI) selectDisplayCallback(g *gocui.Gui, v *gocui.View, menuId string) error {
	mui.CloseView(mui.currentDataView)
	mui.openView(g, menuId)
	return nil
}

func (mui *MasterUI) openView(g *gocui.Gui, viewName string) error {
	ep := mui.router.GetProcessor()
	var dataView masterUIInterface.UpdatableView
	switch viewName {
	case "appListView":
		dataView = appView.NewAppListView(mui, "appListView", mui.headerSize+1, mui.footerSize, ep)
	case "cellListView":
		dataView = cellView.NewCellListView(mui, "cellListView", mui.headerSize+1, mui.footerSize, ep)
	default:
		return errors.New("Unable to find view " + viewName)
	}
	mui.currentDataView = dataView
	mui.layoutManager.Add(dataView)
	dataView.Layout(g)
	mui.addCommonDataViewKeybindings(g, dataView.Name())
	mui.updateDisplay(g)
	return nil
}

func (mui *MasterUI) editUpdateInterval(g *gocui.Gui, v *gocui.View) error {

	labelText := "Seconds:"
	maxLength := 4
	titleText := "Update refresh interval"
	helpText := "no help"

	valueText := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", mui.refreshIntervalMS.Seconds()), "0"), ".")

	applyCallbackFunc := func(g *gocui.Gui, v *gocui.View, w masterUIInterface.Manager, inputValue string) error {
		f, err := strconv.ParseFloat(inputValue, 64)
		if err != nil {
			return err
		}
		if f < float64(0.1) {
			return nil
		}
		mui.refreshIntervalMS = time.Duration(f*1000) * time.Millisecond
		mui.RefeshNow()

		return w.(*uiCommon.InputDialogWidget).CloseWidget(g, v)
	}

	intervalWidget := uiCommon.NewInputDialogWidget(mui,
		"editIntervalWidget", 30, 6, labelText, maxLength, titleText, helpText,
		valueText, applyCallbackFunc)

	return intervalWidget.Init(g)
}

func (mui *MasterUI) clearStats(g *gocui.Gui, v *gocui.View) error {
	mui.router.Clear()
	mui.updateDisplay(g)
	return nil
}

func (mui *MasterUI) RefeshNow() {
	mui.refreshNow <- true
}

func (mui *MasterUI) counter(g *gocui.Gui) {

	mui.updateDisplay(g)
	for {
		select {
		case <-mui.refreshNow:
			mui.updateDisplay(g)
		case <-time.After(mui.refreshIntervalMS):
			mui.updateDisplay(g)
		}
	}
}

func (mui *MasterUI) updateDisplay(g *gocui.Gui) {
	g.Execute(func(g *gocui.Gui) error {
		mui.updateHeaderDisplay(g)
		mui.currentDataView.UpdateDisplay(g)
		return nil
	})
}

func (mui *MasterUI) refreshMetadata(g *gocui.Gui, v *gocui.View) error {
	go mui.router.GetProcessor().LoadCacheAndSeeData()
	return nil
}

func (mui *MasterUI) updateHeaderDisplay(g *gocui.Gui) error {

	v, err := g.View("headerView")
	if err != nil {
		return err
	}
	v.Clear()

	fmt.Fprintf(v, "Evnts: ")
	eventsText := fmt.Sprintf("%v (%v/sec)", util.Format(mui.router.GetEventCount()), mui.router.GetEventRate())
	fmt.Fprintf(v, "%-28v", eventsText)

	runtimeSeconds := Round(time.Now().Sub(mui.router.GetStartTime()), time.Second)
	if runtimeSeconds < time.Second*WarmUpSeconds {
		warmUpTimeRemaining := (time.Second * WarmUpSeconds) - runtimeSeconds
		fmt.Fprintf(v, util.DIM_GREEN)
		fmt.Fprintf(v, " Warm-up: %-10v ", warmUpTimeRemaining)
		fmt.Fprintf(v, util.CLEAR)
	} else {
		fmt.Fprintf(v, "Duration: %-10v ", runtimeSeconds)
	}

	fmt.Fprintf(v, "   %v\n", time.Now().Format("01-02-2006 15:04:05"))

	apiEndpoint, err := mui.cliConnection.ApiEndpoint()
	if err != nil {
		return err
	}

	url, err := url.Parse(apiEndpoint)
	if err != nil {
		return err
	}

	username, err := mui.cliConnection.Username()
	if err != nil {
		return err
	}

	targetDisplay := fmt.Sprintf("%v@%v", username, url.Host)
	fmt.Fprintf(v, "Target: %-78.78v\n", targetDisplay)

	return nil
}

func Round(d, r time.Duration) time.Duration {
	if r <= 0 {
		return d
	}
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d = d - m
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}
