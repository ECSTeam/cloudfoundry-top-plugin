package top

import (
	"fmt"
	"log"
	//"github.com/Sirupsen/logrus"
	"net/url"
	"strconv"
	"strings"
	"time"
	//"github.com/go-errors/errors"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/appStats"
	"github.com/kkellner/cloudfoundry-top-plugin/debug"
	"github.com/kkellner/cloudfoundry-top-plugin/eventrouting"
	"github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/metadata"
	"github.com/kkellner/cloudfoundry-top-plugin/uiCommon"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

type MasterUI struct {
	layoutManager *LayoutManager
	gui           *gocui.Gui
	cliConnection plugin.CliConnection

	appListView *appStats.AppListView

	router            *eventrouting.EventRouter
	refreshNow        chan bool
	refreshIntervalMS time.Duration
}

func NewMasterUI(cliConnection plugin.CliConnection) *MasterUI {

	ui := &MasterUI{
		cliConnection: cliConnection,
		refreshNow:    make(chan bool),
	}

	headerView := NewHeaderWidget(ui, "summaryView", 4)
	footerView := NewFooterWidget("footerView", 4)

	appListView := appStats.NewAppListView(ui, "appListView", 5, 4, ui.cliConnection)
	ui.appListView = appListView
	ui.router = eventrouting.NewEventRouter(appListView.GetCurrentProcessor())

	ui.layoutManager = NewLayoutManager()
	ui.layoutManager.Add(footerView)
	ui.layoutManager.Add(headerView)
	ui.layoutManager.Add(appListView)

	return ui
}

func (ui *MasterUI) CliConnection() plugin.CliConnection {
	return ui.cliConnection
}

func (ui *MasterUI) LayoutManager() masterUIInterface.LayoutManagerInterface {
	return ui.layoutManager
}
func (ui *MasterUI) GetRouter() *eventrouting.EventRouter {
	return ui.router
}

func (ui *MasterUI) Start() {
	ui.appListView.Start()
	ui.initGui()
}

func (ui *MasterUI) initGui() {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	ui.gui = g
	g.InputEsc = true
	defer g.Close()
	g.SetManager(ui.layoutManager)

	debug.InitDebug(g, ui)

	// default refresh to 1 second
	ui.refreshIntervalMS = 1000 * time.Millisecond

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("appListView", 'q', gocui.ModNone, ui.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("appListView", 'Q', gocui.ModNone, ui.quit); err != nil {
		log.Panicln(err)
	}
	/*
	  if err := g.SetKeybinding("appListView", gocui.KeyEsc, gocui.ModNone, ui.quit); err != nil {
	    log.Panicln(err)
	  }
	*/
	if err := g.SetKeybinding("appListView", 'c', gocui.ModNone, ui.clearStats); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		ui.RefeshNow()
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("appListView", 's', gocui.ModNone, ui.editUpdateInterval); err != nil {
		log.Panicln(err)
	}

	go ui.counter(g)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		//log.Panicln(err.(*errors.Error).ErrorStack())
		log.Panicln(err)
	}

}

func (ui *MasterUI) CloseView(m masterUIInterface.Manager) error {

	ui.gui.DeleteView(m.Name())
	ui.gui.DeleteKeybindings(m.Name())
	nextForFocus := ui.layoutManager.Remove(m)

	if err := ui.SetCurrentViewOnTop(ui.gui, nextForFocus.Name()); err != nil {
		return err
	}
	return nil
}

func (ui *MasterUI) CloseViewByName(viewName string) error {
	m := ui.layoutManager.GetManagerByViewName(viewName)
	return ui.CloseView(m)
}

func (ui *MasterUI) SetCurrentViewOnTop(g *gocui.Gui, name string) error {
	if _, err := g.SetCurrentView(name); err != nil {
		return err
	}
	if _, err := g.SetViewOnTop(name); err != nil {
		return err
	}
	return nil
}

func (ui *MasterUI) GetCurrentView(g *gocui.Gui) *gocui.View {
	views := g.Views()
	view := views[len(views)-1]
	return view
}

func (ui *MasterUI) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (ui *MasterUI) editUpdateInterval(g *gocui.Gui, v *gocui.View) error {

	labelText := "Seconds:"
	maxLength := 4
	titleText := "Update refresh interval"
	helpText := "no help"

	valueText := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", ui.refreshIntervalMS.Seconds()), "0"), ".")

	applyCallbackFunc := func(g *gocui.Gui, v *gocui.View, w masterUIInterface.Manager, inputValue string) error {
		f, err := strconv.ParseFloat(inputValue, 64)
		if err != nil {
			return err
		}
		if f < float64(0.1) {
			return nil
		}
		ui.refreshIntervalMS = time.Duration(f*1000) * time.Millisecond
		ui.RefeshNow()

		return w.(*uiCommon.InputDialogWidget).CloseWidget(g, v)
	}

	intervalWidget := uiCommon.NewInputDialogWidget(ui,
		"editIntervalWidget", 30, 6, labelText, maxLength, titleText, helpText,
		valueText, applyCallbackFunc)
	return intervalWidget.Init(g)
}

func (ui *MasterUI) clearStats(g *gocui.Gui, v *gocui.View) error {
	ui.router.Clear()
	ui.appListView.ClearStats(g, v)
	ui.updateDisplay(g)
	return nil

}

func (ui *MasterUI) RefeshNow() {
	ui.refreshNow <- true
}

func (ui *MasterUI) counter(g *gocui.Gui) {

	ui.updateDisplay(g)
	for {
		select {
		case <-ui.refreshNow:
			ui.updateDisplay(g)
		case <-time.After(ui.refreshIntervalMS):
			ui.updateDisplay(g)
		}
	}
}

func (ui *MasterUI) updateDisplay(g *gocui.Gui) {
	g.Execute(func(g *gocui.Gui) error {
		ui.updateHeaderDisplay(g)
		ui.appListView.UpdateDisplay(g)
		return nil
	})
}

func (ui *MasterUI) refeshDisplayX(g *gocui.Gui) {
	g.Execute(func(g *gocui.Gui) error {
		ui.updateHeaderDisplay(g)
		ui.appListView.RefreshDisplay(g)
		return nil
	})
}

func (ui *MasterUI) updateHeaderDisplay(g *gocui.Gui) error {

	v, err := g.View("summaryView")
	if err != nil {
		return err
	}
	v.Clear()

	fmt.Fprintf(v, "Total events: %-13v", util.Format(ui.router.GetEventCount()))

	runtimeSeconds := Round(time.Now().Sub(ui.router.GetStartTime()), time.Second)
	if runtimeSeconds < time.Second*30 {
		warmUpTimeRemaining := (time.Second * 30) - runtimeSeconds
		fmt.Fprintf(v, util.DIM_GREEN)
		fmt.Fprintf(v, "Warm-up period: %-10v ", warmUpTimeRemaining)
		fmt.Fprintf(v, util.CLEAR)
	} else {
		fmt.Fprintf(v, "Stats duration: %-10v ", runtimeSeconds)
	}

	fmt.Fprintf(v, "%v\n", time.Now().Format("01-02-2006 15:04:05"))

	apiEndpoint, err := ui.cliConnection.ApiEndpoint()
	if err != nil {
		return err
	}

	url, err := url.Parse(apiEndpoint)
	if err != nil {
		return err
	}

	username, err := ui.cliConnection.Username()
	if err != nil {
		return err
	}

	targetDisplay := fmt.Sprintf("%v@%v", username, url.Host)
	fmt.Fprintf(v, "Target: %-45.45v ", targetDisplay)

	displayTotalMem := "--"
	totalMem := metadata.GetTotalMemoryAllStartedApps()
	if totalMem > 0 {
		displayTotalMem = util.ByteSize(totalMem).String()
	}
	// Total quota memory of all running app instances
	fmt.Fprintf(v, "Rsrvd Mem: %8v ", displayTotalMem)

	displayTotalDisk := "--"
	totalDisk := metadata.GetTotalDiskAllStartedApps()
	if totalMem > 0 {
		displayTotalDisk = util.ByteSize(totalDisk).String()
	}

	fmt.Fprintf(v, "Rsrvd Disk: %8v ", displayTotalDisk)
	fmt.Fprintf(v, "\n")
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
