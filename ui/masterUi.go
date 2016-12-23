// Copyright (c) 2016 ECS Team, Inc. - All Rights Reserved
// https://github.com/ECSTeam/cloudfoundry-top-plugin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ui

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ansel1/merry"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventrouting"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/capacityPlanView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/cellView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/routeView"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

const WarmUpSeconds = 60
const DefaultRefreshInternalMS = 1000
const HELP_TEXT_VIEW_NAME = "helpTextTipsView"

type MasterUI struct {
	layoutManager *uiCommon.LayoutManager
	gui           *gocui.Gui
	cliConnection plugin.CliConnection

	currentDataView masterUIInterface.UpdatableView

	router            *eventrouting.EventRouter
	refreshNow        chan bool
	refreshIntervalMS time.Duration

	baseHeaderSize       int
	headerSize           int
	helpTextTipsViewSize int

	displayMenuId string
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

	mui.helpTextTipsViewSize = 4
	helpTextTipsView := NewHelpTextTipsWidget(mui, HELP_TEXT_VIEW_NAME, mui.helpTextTipsViewSize)
	mui.layoutManager.Add(helpTextTipsView)

	mui.baseHeaderSize = 3
	mui.headerSize = 6
	headerView := NewHeaderWidget(mui, "headerView")
	mui.layoutManager.Add(headerView)
	// We add the common keybindings to the header view in the event
	// that no DataView is open
	mui.AddCommonDataViewKeybindings(g, "headerView")

	mui.createAndOpenView(g, "appListView")

	// default refresh to 1 second
	mui.refreshIntervalMS = DefaultRefreshInternalMS * time.Millisecond

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, mui.quit); err != nil {
		log.Panicln(err)
	}

	go mui.refreshDataAndDisplayThread(g)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		m := merry.Details(err)
		log.Panicln(m)
	}

}

func (mui *MasterUI) SetStatsSummarySize(statSummarySize int) {
	mui.headerSize = mui.baseHeaderSize + statSummarySize
	//toplog.Info(fmt.Sprintf("headerSize set to: %v  statSummarySize: %v", mui.headerSize, statSummarySize))
}

func (mui *MasterUI) GetHeaderSize() int {
	return mui.headerSize
}

// Add keybindings for top level data views -- note must also call addCommonDataViewKeybindings
// to get a full set of keybindings
func (mui *MasterUI) addTopLevelDataViewKeybindings(g *gocui.Gui, viewName string) error {
	if err := g.SetKeybinding(viewName, 'q', gocui.ModNone, mui.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, 'd', gocui.ModNone, mui.selectDisplayAction); err != nil {
		log.Panicln(err)
	}
	return nil
}

// Add common keybindings for all data views -- note that this does not include
// keybindings for "top level" data views which are ones that are selectable from
// the "select view" menu ('d' command)
func (mui *MasterUI) AddCommonDataViewKeybindings(g *gocui.Gui, viewName string) error {
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
	if err := g.SetKeybinding(viewName, 'r', gocui.ModNone, mui.refreshMetadata); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(viewName, 'Z', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			toplog.Debug(fmt.Sprintf("Top: %v", mui.layoutManager.Top().Name()))
			return nil
		}); err != nil {
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

	//toplog.Debug(fmt.Sprintf("type:%v", checkIfUpdatableView(nextForFocus)))

	nextViewName := nextForFocus.Name()
	if err := mui.SetCurrentViewOnTop(mui.gui); err != nil {
		return merry.Wrap(err).Appendf("SetCurrentViewOnTop viewName:[%v]", nextViewName)
	}

	if mui.checkIfUpdatableView(nextForFocus) {
		mui.currentDataView = nextForFocus.(masterUIInterface.UpdatableView)
		mui.updateHeaderDisplay(mui.gui)
		mui.currentDataView.RefreshDisplay(mui.gui)
	}

	return nil
}

func (mui *MasterUI) checkIfUpdatableView(x interface{}) bool {
	// Declare a type object representing UpdatableView
	updatableView := reflect.TypeOf((*masterUIInterface.UpdatableView)(nil)).Elem()
	// Get a type object of the pointer on the object represented by the parameter
	// and see if it implements UpdatableView
	return reflect.TypeOf(x).Implements(updatableView)
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

func (mui *MasterUI) SetHelpTextTips(g *gocui.Gui, helpTextTips string) error {
	helpMgr := mui.layoutManager.GetManagerByViewName(HELP_TEXT_VIEW_NAME)
	if helpMgr != nil {
		helpTextTipsWidget := helpMgr.(*HelpTextTipsWidget)
		helpTextTipsWidget.SetHelpTextTips(g, helpTextTips)
	}
	return nil
}

func (mui *MasterUI) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (mui *MasterUI) selectDisplayAction(g *gocui.Gui, v *gocui.View) error {

	menuItems := make([]*uiCommon.MenuItem, 0, 5)
	menuItems = append(menuItems, uiCommon.NewMenuItem("appListView", "App Stats"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("cellListView", "Cell Stats"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("routeListView", "Route Stats"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("capacityPlanView", "Capacity Plan (memory)"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("eventstats", "TODO: Event Stats"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("eventhistory", "TODO: Event Rate History"))
	selectDisplayView := uiCommon.NewSelectMenuWidget(mui, "selectDisplayView", "Select Display", menuItems, mui.selectDisplayCallback)
	selectDisplayView.SetMenuId(mui.displayMenuId)

	mui.LayoutManager().Add(selectDisplayView)
	mui.SetCurrentViewOnTop(g)
	return nil
}

func (mui *MasterUI) selectDisplayCallback(g *gocui.Gui, v *gocui.View, menuId string) error {
	mui.displayMenuId = menuId
	mui.createAndOpenView(g, menuId)
	return nil
}

func (mui *MasterUI) createAndOpenView(g *gocui.Gui, viewName string) error {

	if mui.layoutManager.ContainsViewName(viewName) {
		mui.layoutManager.SetCurrentView(viewName)
		mui.SetCurrentViewOnTop(g)
		mui.updateDisplay(g)
		return nil
	}

	ep := mui.router.GetProcessor()
	var dataView masterUIInterface.UpdatableView
	switch viewName {
	case "appListView":
		dataView = appView.NewAppListView(mui, "appListView", mui.helpTextTipsViewSize, ep)
	case "cellListView":
		dataView = cellView.NewCellListView(mui, "cellListView", mui.helpTextTipsViewSize, ep)
	case "routeListView":
		dataView = routeView.NewRouteListView(mui, "routeListView", mui.helpTextTipsViewSize, ep)
	case "capacityPlanView":
		dataView = capacityPlanView.NewCapacityPlanView(mui, "capacityPlanView", mui.helpTextTipsViewSize, ep)

	default:
		return errors.New("Unable to find view " + viewName)
	}
	mui.OpenView(g, dataView)
	mui.addTopLevelDataViewKeybindings(g, dataView.Name())
	return nil
}

func (mui *MasterUI) OpenView(g *gocui.Gui, dataView masterUIInterface.UpdatableView) error {
	mui.currentDataView = dataView
	mui.layoutManager.Add(dataView)
	dataView.Layout(g)
	mui.AddCommonDataViewKeybindings(g, dataView.Name())
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

func (mui *MasterUI) refreshDataAndDisplayThread(g *gocui.Gui) {

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
	go mui.router.GetProcessor().LoadCacheAndSeedData()
	return nil
}

func (mui *MasterUI) IsWarmupComplete() bool {
	runtimeSeconds := Round(time.Now().Sub(mui.router.GetStartTime()), time.Second)
	return runtimeSeconds > time.Second*WarmUpSeconds
}

func (mui *MasterUI) SetMinimizeHeader(g *gocui.Gui, minimizeHeader bool) {
	// TODO: The header needs to be the same accross data display types -- so need
	// to move the responsability of displaying the summary stats to a common location
	mui.RefeshNow()
	//mui.updateHeaderDisplay(g)
}

func (mui *MasterUI) updateHeaderDisplay(g *gocui.Gui) error {

	v, err := g.View("headerView")
	if err != nil {
		return err
	}
	v.Clear()

	fmt.Fprintf(v, "Events: ")
	eventsText := fmt.Sprintf("%v (%v/sec)", util.Format(mui.router.GetEventCount()), mui.router.GetEventRate())
	fmt.Fprintf(v, "%-27v", eventsText)

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
