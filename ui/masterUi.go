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
	termbox "github.com/nsf/termbox-go"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventrouting"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/dataCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/interfaces/managerUI"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/aboutView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/alertView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/appViews/appView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/capacityPlanView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/cellViews/cellView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/eventRateHistoryView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/eventViews/eventView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/headerView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/orgSpaceViews/orgView"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/routeViews/routeView"
	"github.com/jroimartin/gocui"
)

const DefaultRefreshInternalMS = 1000
const HELP_TEXT_VIEW_NAME = "helpTextTipsView"
const STATUS_VIEW_NAME = "statusView"

type MasterUI struct {
	layoutManager  *uiCommon.LayoutManager
	gui            *gocui.Gui
	cliConnection  plugin.CliConnection
	pluginMetadata *plugin.PluginMetadata
	privileged     bool
	username       string
	targetDisplay  string

	headerView      *headerView.HeaderWidget
	alertManager    *alertView.AlertManager
	currentDataView masterUIInterface.UpdatableView

	router            *eventrouting.EventRouter
	refreshNow        chan bool
	refreshIntervalMS time.Duration
	displayPaused     bool
	editColumnMode    bool
	headerMinimized   bool
	commonData        *dataCommon.CommonData

	statusMsg chan string

	//baseHeaderSize       int
	//headerSize           int
	helpTextTipsViewSize int

	displayMenuId string
}

func NewMasterUI(cliConnection plugin.CliConnection, pluginMetadata *plugin.PluginMetadata, privileged bool) *MasterUI {

	mui := &MasterUI{
		cliConnection:  cliConnection,
		pluginMetadata: pluginMetadata,
		privileged:     privileged,
		refreshNow:     make(chan bool),
		statusMsg:      make(chan string),
	}

	eventProcessor := eventdata.NewEventProcessor(mui.cliConnection, privileged, mui.statusMsg)
	mui.router = eventrouting.NewEventRouter(eventProcessor)

	username, err := cliConnection.Username()
	if err != nil {
		toplog.Info("Unable to get username, error: %v", err)
		mui.username = "UNKNOWN"
	} else {
		mui.username = username
	}
	apiEndpoint, err := mui.cliConnection.ApiEndpoint()
	if err != nil {
		toplog.Info("Unable to get apiEndpoint, error: %v", err)
	}

	url, err := url.Parse(apiEndpoint)
	if err != nil {
		toplog.Info("Unable to get Parse apiEndpoint, error: %v", err)
	} else {
		usernameDisplay := fmt.Sprintf("(%v)", username)
		if privileged {
			usernameDisplay = fmt.Sprintf("[%v]", username)
		}
		mui.targetDisplay = fmt.Sprintf("%v%v", usernameDisplay, url.Host)
	}
	return mui
}

func (mui *MasterUI) CliConnection() plugin.CliConnection {
	return mui.cliConnection
}

func (mui *MasterUI) IsPrivileged() bool {
	return mui.privileged
}

func (mui *MasterUI) LayoutManager() managerUI.LayoutManagerInterface {
	return mui.layoutManager
}

func (mui *MasterUI) GetRouter() *eventrouting.EventRouter {
	return mui.router
}

func (mui *MasterUI) GetCommonData() *dataCommon.CommonData {
	return mui.commonData
}

func (mui *MasterUI) GetTargetDisplay() string {
	return mui.targetDisplay
}

func (mui *MasterUI) Start(monitoredAppGuids map[string]bool) {
	mui.router.GetProcessor().Start()
	mui.initGui(monitoredAppGuids)
}

func (mui *MasterUI) initGui(monitoredAppGuids map[string]bool) {

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

	statusView := NewStatusWidget(mui, STATUS_VIEW_NAME)
	mui.layoutManager.Add(statusView)

	mui.commonData = dataCommon.NewCommonData(mui.router, monitoredAppGuids)

	mui.alertManager = alertView.NewAlertManager(mui, mui.commonData)
	//mui.baseHeaderSize = 3
	//mui.headerSize = 6
	mui.headerView = headerView.NewHeaderWidget(mui, "headerView", mui.router, mui.commonData)
	mui.layoutManager.Add(mui.headerView)

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
	go mui.refreshStatusThread(g)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		m := merry.Details(err)
		log.Panicln(m)
	}

}

func (mui *MasterUI) flushKeyboardBuffer() {

	go func() {
		loop := true
		for loop {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
			case termbox.EventError:
				loop = false
			case termbox.EventInterrupt:
				loop = false
			}
		}
	}()
	time.Sleep(10 * time.Microsecond)
	termbox.Interrupt()
}

func (mui *MasterUI) GetHeaderSize() int {
	return mui.headerView.HeaderSize
}

func (mui *MasterUI) GetAlertSize() int {
	return mui.alertManager.AlertSize
}

func (mui *MasterUI) GetTopMargin() int {
	size := mui.headerView.HeaderSize
	if !mui.IsEditColumnMode() {
		size = size + mui.alertManager.AlertSize
	}
	return size
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
	if err := g.SetKeybinding(viewName, 'p', gocui.ModNone, mui.toggleDisplayPauseAction); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, 'H', gocui.ModNone, mui.toggleHeaderMinimizeAction); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(viewName, 'E', gocui.ModNone, mui.logTestError); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(viewName, 'Z', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			toplog.Debug("Top: %v", mui.layoutManager.Top().Name())
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

	/*
		// ******** BEGIN: Test List scrolling *************
		if err := g.SetKeybinding(viewName, 'z', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			mui.router.GetProcessor().GetMetadataManager().GetAppMdManager().CreateTestData(10)
			mui.router.GetProcessor().SeedStatsFromMetadata()
			mui.router.Clear()
			mui.updateDisplay(g)
			return nil
		}); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding(viewName, 'x', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			mui.router.GetProcessor().GetMetadataManager().GetAppMdManager().CreateTestData(50)
			mui.router.GetProcessor().SeedStatsFromMetadata()
			mui.router.Clear()
			mui.updateDisplay(g)
			return nil
		}); err != nil {
			log.Panicln(err)
		}
		// ******** END: Test List scrolling *************
	*/

	/*
		// ******** BEGIN: Test Alert Messages *************
		if err := g.SetKeybinding(viewName, 'z', gocui.ModNone, mui.testShowUserMessage); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(viewName, 'x', gocui.ModNone, mui.test2ShowUserMessage); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(viewName, 'a', gocui.ModNone, mui.testClearUserMessage); err != nil {
			log.Panicln(err)
		}
		// ******** END: Test Alert Messages *************
	*/
	return nil
}

func (mui *MasterUI) testShowUserMessage(g *gocui.Gui, v *gocui.View) error {
	return mui.alertManager.ShowMessage(g, alertView.APPS_NOT_IN_DESIRED_STATE, 99, "s")
}

func (mui *MasterUI) test2ShowUserMessage(g *gocui.Gui, v *gocui.View) error {
	return mui.alertManager.ShowMessage(g, alertView.TestMessage)
}

func (mui *MasterUI) testClearUserMessage(g *gocui.Gui, v *gocui.View) error {
	return mui.alertManager.ClearUserMessage(g, alertView.TestMessage)
}

func (mui *MasterUI) CloseView(m managerUI.Manager) error {

	mui.gui.DeleteView(m.Name())
	mui.gui.DeleteKeybindings(m.Name())
	nextForFocus := mui.layoutManager.Remove(m)
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

	viewMgr := mui.layoutManager.Top()
	viewMgr.Layout(g)

	topName := viewMgr.Name()
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

func (mui *MasterUI) SetStatus(g *gocui.Gui, status string) error {
	statusMgr := mui.layoutManager.GetManagerByViewName(STATUS_VIEW_NAME)
	if statusMgr != nil {
		statusWidget := statusMgr.(*StatusWidget)
		statusWidget.SetStatus(g, status)
	}
	return nil
}

func (mui *MasterUI) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (mui *MasterUI) selectDisplayAction(g *gocui.Gui, v *gocui.View) error {

	menuItems := make([]*uiCommon.MenuItem, 0, 5)
	menuItems = append(menuItems, uiCommon.NewMenuItem("appListView", "App Stats"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("orgListView", "Org Stats"))
	if mui.privileged {
		menuItems = append(menuItems, uiCommon.NewMenuItem("cellListView", "Cell Stats"))
	}
	menuItems = append(menuItems, uiCommon.NewMenuItem("routeListView", "Route Stats"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("eventRateHistoryListView", "Event Rate History"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("eventListView", "Event Stats"))
	if mui.privileged {
		menuItems = append(menuItems, uiCommon.NewMenuItem("capacityPlanView", "Capacity Plan (memory)"))
	}
	menuItems = append(menuItems, uiCommon.NewMenuItem("aboutView", "About Top"))

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
		dataView = appView.NewAppListView(mui, nil, "appListView", mui.helpTextTipsViewSize, ep, "")
	case "orgListView":
		dataView = orgView.NewOrgListView(mui, "orgListView", mui.helpTextTipsViewSize, ep)
	case "cellListView":
		dataView = cellView.NewCellListView(mui, "cellListView", mui.helpTextTipsViewSize, ep)
	case "routeListView":
		dataView = routeView.NewRouteListView(mui, "routeListView", mui.helpTextTipsViewSize, ep)
	case "eventListView":
		dataView = eventView.NewEventListView(mui, "eventListView", mui.helpTextTipsViewSize, ep)
	case "capacityPlanView":
		dataView = capacityPlanView.NewCapacityPlanView(mui, "capacityPlanView", mui.helpTextTipsViewSize, ep)
	case "eventRateHistoryListView":
		dataView = eventRateHistoryView.NewEventRateHistoryView(mui, "eventRateHistoryListView", mui.helpTextTipsViewSize, ep)
	case "aboutView":
		dataView = aboutView.NewTopView(mui, "aboutView", mui.helpTextTipsViewSize, ep, mui.pluginMetadata)

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

	valueText := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", mui.refreshIntervalMS.Seconds()), "0"), ".")

	applyCallbackFunc := func(g *gocui.Gui, v *gocui.View, w managerUI.Manager, inputValue string) error {
		f, err := strconv.ParseFloat(inputValue, 64)
		if err != nil {
			return err
		}
		// Minimum refresh time
		if f < float64(0.01) {
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

func (mui *MasterUI) toggleDisplayPauseAction(g *gocui.Gui, v *gocui.View) error {
	mui.SetDisplayPaused(!mui.GetDisplayPaused())
	//mui.updateHeaderDisplay(mui.gui)
	return mui.currentDataView.RefreshDisplay(mui.gui)
}

func (mui *MasterUI) toggleHeaderMinimizeAction(g *gocui.Gui, v *gocui.View) error {
	mui.SetHeaderMinimize(g, !mui.IsHeaderMinimized())
	return mui.currentDataView.RefreshDisplay(mui.gui)
}

func (mui *MasterUI) logTestError(g *gocui.Gui, v *gocui.View) error {
	toplog.Error("test error")
	return nil
}

func (mui *MasterUI) GetDisplayPaused() bool {
	return mui.displayPaused
}

func (mui *MasterUI) SetDisplayPaused(paused bool) {
	mui.displayPaused = paused
	mui.router.GetProcessor().GetCurrentEventRateHistory().SetFreezeData(paused)

	if !paused {
		mui.snapshotLiveData()
	}
	// Moved updateHeaderDisplay here from toggleDisplayPauseAction
	mui.updateHeaderDisplay(mui.gui)
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

func (mui *MasterUI) refreshStatusThread(g *gocui.Gui) {

	for {
		status := <-mui.statusMsg
		mui.SetStatus(g, status)
		g.Execute(func(g *gocui.Gui) error {
			statusMgr := mui.layoutManager.GetManagerByViewName(STATUS_VIEW_NAME)
			if statusMgr != nil {
				statusWidget := statusMgr.(*StatusWidget)
				statusWidget.ShowStatus(g)
			}
			return nil
		})
	}
}

func (mui *MasterUI) snapshotLiveData() {
	mui.router.GetProcessor().UpdateData()
}

func (mui *MasterUI) updateDisplay(g *gocui.Gui) {
	g.Execute(func(g *gocui.Gui) error {

		if !mui.displayPaused {
			// This takes a snapshot of the live data
			mui.snapshotLiveData()
			mui.commonData.PostProcessData()
		}

		mui.updateHeaderDisplay(g)
		mui.currentDataView.UpdateDisplay(g)
		return nil
	})
}

func (mui *MasterUI) refreshMetadata(g *gocui.Gui, v *gocui.View) error {
	go mui.router.GetProcessor().FlushCache()
	return nil
}

func (mui *MasterUI) IsWarmupComplete() bool {
	return mui.commonData.IsWarmupComplete()
}

func (mui *MasterUI) IsEditColumnMode() bool {
	return mui.editColumnMode
}

func (mui *MasterUI) SetEditColumnMode(g *gocui.Gui, editColumnMode bool) {
	mui.editColumnMode = editColumnMode
	mui.updateHeaderDisplay(g)
}

func (mui *MasterUI) IsHeaderMinimized() bool {
	return mui.editColumnMode || mui.headerMinimized
}

func (mui *MasterUI) SetHeaderMinimize(g *gocui.Gui, minimizeHeader bool) {
	toplog.Debug("SetMinimizeHeader:%v", minimizeHeader)
	mui.headerMinimized = minimizeHeader
	//mui.RefeshNow()
	mui.updateHeaderDisplay(g)
}

func (mui *MasterUI) updateHeaderDisplay(g *gocui.Gui) error {
	err := mui.headerView.UpdateDisplay(g)
	if err != nil {
		return err
	}

	return mui.alertManager.CheckForAlerts(g)
}
