package appStats

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/uiCommon"
	//"github.com/mohae/deepcopy"
	"github.com/kkellner/cloudfoundry-top-plugin/debug"
	"github.com/kkellner/cloudfoundry-top-plugin/helpView"
	"github.com/kkellner/cloudfoundry-top-plugin/metadata"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

type AppListView struct {
	masterUI     masterUIInterface.MasterUIInterface
	name         string
	topMargin    int
	bottomMargin int

	currentProcessor        *AppStatsEventProcessor
	displayedProcessor      *AppStatsEventProcessor
	displayedSortedStatList []*AppStats

	cliConnection plugin.CliConnection
	mu            sync.Mutex
	filterAppName string

	appDetailView *AppDetailView
	appListWidget *uiCommon.ListWidget

	displayPaused bool
}

func NewAppListView(masterUI masterUIInterface.MasterUIInterface, name string, topMargin, bottomMargin int,
	cliConnection plugin.CliConnection) *AppListView {

	currentProcessor := NewAppStatsEventProcessor()
	displayedProcessor := NewAppStatsEventProcessor()

	return &AppListView{
		masterUI:           masterUI,
		name:               name,
		topMargin:          topMargin,
		bottomMargin:       bottomMargin,
		cliConnection:      cliConnection,
		currentProcessor:   currentProcessor,
		displayedProcessor: displayedProcessor}
}

func (asUI *AppListView) Name() string {
	return asUI.name
}

func (asUI *AppListView) Layout(g *gocui.Gui) error {

	if asUI.appListWidget == nil {

		appListWidget := uiCommon.NewListWidget(asUI.masterUI, asUI.name,
			asUI.topMargin, asUI.bottomMargin, asUI, asUI.columnDefinitions())
		appListWidget.Title = "App List"
		appListWidget.PreRowDisplayFunc = asUI.PreRowDisplay
		appListWidget.GetListSize = asUI.GetListSize
		appListWidget.GetUnfilteredListSize = asUI.GetUnfilteredListSize
		appListWidget.GetRowKey = asUI.GetRowKey

		defaultSortColums := []*uiCommon.SortColumn{
			uiCommon.NewSortColumn("CPU", true),
			uiCommon.NewSortColumn("L60", true),
			uiCommon.NewSortColumn("appName", false),
			uiCommon.NewSortColumn("spaceName", false),
			uiCommon.NewSortColumn("orgName", false),
		}
		appListWidget.SetSortColumns(defaultSortColums)

		asUI.appListWidget = appListWidget

		/*
			if err := g.SetKeybinding(asUI.name, 'f', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
				filter := NewFilterWidget(asUI.masterUI, "filterWidget", 30, 10)
				asUI.masterUI.LayoutManager().Add(filter)
				asUI.masterUI.SetCurrentViewOnTop(g, "filterWidget")
				return nil
			}); err != nil {
				log.Panicln(err)
			}
		*/

		if err := g.SetKeybinding(asUI.name, 'h', gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				helpView := helpView.NewHelpView(asUI.masterUI, "helpView", 75, 17, helpText)
				asUI.masterUI.LayoutManager().Add(helpView)
				asUI.masterUI.SetCurrentViewOnTop(g, "helpView")
				return nil
			}); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(asUI.name, 'c', gocui.ModNone, asUI.copyAction); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(asUI.name, 'r', gocui.ModNone, asUI.refreshMetadata); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding(asUI.name, 'D', gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				//msg := "Test debug message 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 "
				/*
				   msg := "test"
				   for i:=0;i<50;i++ {
				     debug.Info(msg)
				   }
				*/
				debug.Open()
				return nil
			}); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding(asUI.name, gocui.KeyEnter, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				if asUI.appListWidget.HighlightKey() != "" {
					asUI.appDetailView = NewAppDetailView(asUI.masterUI, "appDetailView", asUI.appListWidget.HighlightKey(), asUI)
					asUI.masterUI.LayoutManager().Add(asUI.appDetailView)
					asUI.masterUI.SetCurrentViewOnTop(g, "appDetailView")
				}
				return nil
			}); err != nil {
			log.Panicln(err)
		}

	}
	return asUI.appListWidget.Layout(g)
}

func (asUI *AppListView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, asUI.columnAppName())
	columns = append(columns, asUI.columnSpaceName())
	columns = append(columns, asUI.columnOrgName())

	columns = append(columns, asUI.columnReportingContainers())
	columns = append(columns, asUI.columnTotalCpu())
	columns = append(columns, asUI.columnTotalMemoryUsed())
	columns = append(columns, asUI.columnTotalDiskUsed())

	columns = append(columns, asUI.columnAvgResponseTimeL60Info())
	columns = append(columns, asUI.columnLogCount())

	columns = append(columns, asUI.columnL1())
	columns = append(columns, asUI.columnL10())
	columns = append(columns, asUI.columnL60())

	columns = append(columns, asUI.columnHTTP())
	columns = append(columns, asUI.column2XX())
	columns = append(columns, asUI.column3XX())
	columns = append(columns, asUI.column4XX())
	columns = append(columns, asUI.column5XX())

	return columns
}

func formatDisplayData(value string, size int) string {
	if len(value) > size {
		value = value[0:size-1] + uiCommon.Ellipsis
	}
	format := fmt.Sprintf("%%-%v.%vv", size, size)
	return fmt.Sprintf(format, value)
}

func (asUI *AppListView) columnAppName() *uiCommon.ListColumn {
	defaultColSize := 50
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*AppStats).AppName, c2.(*AppStats).AppName)
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		//return fmt.Sprintf("%-50.50v", appStats.AppName)
		return formatDisplayData(appStats.AppName, defaultColSize)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return appStats.AppName
	}
	c := uiCommon.NewListColumn("appName", "APPLICATION", defaultColSize,
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnSpaceName() *uiCommon.ListColumn {
	defaultColSize := 10
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*AppStats).SpaceName, c2.(*AppStats).SpaceName)
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		//return fmt.Sprintf("%-10.10v", appStats.SpaceName)
		return formatDisplayData(appStats.SpaceName, defaultColSize)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return appStats.SpaceName
	}
	c := uiCommon.NewListColumn("spaceName", "SPACE", defaultColSize,
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnOrgName() *uiCommon.ListColumn {
	defaultColSize := 10
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*AppStats).OrgName, c2.(*AppStats).OrgName)
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		//return fmt.Sprintf("%-10.10v", appStats.OrgName)
		return formatDisplayData(appStats.OrgName, defaultColSize)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return appStats.OrgName
	}
	c := uiCommon.NewListColumn("orgName", "ORG", defaultColSize,
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnReportingContainers() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*AppStats).TotalReportingContainers < c2.(*AppStats).TotalReportingContainers
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%3v", appStats.TotalReportingContainers)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return strconv.Itoa(appStats.TotalReportingContainers)
	}
	c := uiCommon.NewListColumn("reportingContainers", "RCR", 3,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnTotalCpu() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*AppStats).TotalCpuPercentage < c2.(*AppStats).TotalCpuPercentage
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		totalCpuInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalCpuInfo = fmt.Sprintf("%6v", "--")
		} else {
			if appStats.TotalCpuPercentage >= 100.0 {
				totalCpuInfo = fmt.Sprintf("%6.0f", appStats.TotalCpuPercentage)
			} else if appStats.TotalCpuPercentage >= 10.0 {
				totalCpuInfo = fmt.Sprintf("%6.1f", appStats.TotalCpuPercentage)
			} else {
				totalCpuInfo = fmt.Sprintf("%6.2f", appStats.TotalCpuPercentage)
			}
		}
		return fmt.Sprintf("%6v", totalCpuInfo)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%.2f", appStats.TotalCpuPercentage)
	}
	c := uiCommon.NewListColumn("CPU", "CPU%", 6,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnTotalMemoryUsed() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*AppStats).TotalUsedMemory < c2.(*AppStats).TotalUsedMemory
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		totalMemInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalMemInfo = fmt.Sprintf("%9v", "--")
		} else {
			totalMemInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.TotalUsedMemory))
		}
		return fmt.Sprintf("%9v", totalMemInfo)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%v", appStats.TotalUsedMemory)
	}
	c := uiCommon.NewListColumn("MEM", "MEM", 9,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnTotalDiskUsed() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*AppStats).TotalUsedDisk < c2.(*AppStats).TotalUsedDisk
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		totalDiskInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalDiskInfo = fmt.Sprintf("%9v", "--")
		} else {
			totalDiskInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.TotalUsedDisk))
		}
		return fmt.Sprintf("%9v", totalDiskInfo)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%v", appStats.TotalUsedDisk)
	}
	c := uiCommon.NewListColumn("DISK", "DISK", 9,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnAvgResponseTimeL60Info() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*AppStats).TotalTraffic.AvgResponseL60Time < c2.(*AppStats).TotalTraffic.AvgResponseL60Time
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		avgResponseTimeL60Info := "--"
		if appStats.TotalTraffic.AvgResponseL60Time >= 0 {
			avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL60Time / 1000000
			avgResponseTimeL60Info = fmt.Sprintf("%6.0f", avgResponseTimeMs)
		}
		return fmt.Sprintf("%6v", avgResponseTimeL60Info)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%v", appStats.TotalTraffic.AvgResponseL60Time)
	}
	c := uiCommon.NewListColumn("avgResponseTimeL60", "RESP", 6,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnLogCount() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*AppStats).TotalLogCount < c2.(*AppStats).TotalLogCount
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%11v", util.Format(appStats.TotalLogCount))
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%v", appStats.TotalLogCount)
	}
	c := uiCommon.NewListColumn("totalLogCount", "LOGS", 11,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnL1() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*AppStats).TotalTraffic.EventL1Rate < c2.(*AppStats).TotalTraffic.EventL1Rate
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%5d", appStats.TotalTraffic.EventL1Rate)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%v", appStats.TotalTraffic.EventL1Rate)
	}
	c := uiCommon.NewListColumn("L1", "L1", 5,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnL10() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*AppStats).TotalTraffic.EventL10Rate < c2.(*AppStats).TotalTraffic.EventL10Rate
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%5d", appStats.TotalTraffic.EventL10Rate)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%v", appStats.TotalTraffic.EventL10Rate)
	}
	c := uiCommon.NewListColumn("L10", "L10", 5,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnL60() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*AppStats).TotalTraffic.EventL60Rate < c2.(*AppStats).TotalTraffic.EventL60Rate
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%5d", appStats.TotalTraffic.EventL60Rate)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%v", appStats.TotalTraffic.EventL60Rate)
	}
	c := uiCommon.NewListColumn("L60", "L60", 5,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) columnHTTP() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*AppStats).TotalTraffic.HttpAllCount < c2.(*AppStats).TotalTraffic.HttpAllCount
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%8d", appStats.TotalTraffic.HttpAllCount)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%v", appStats.TotalTraffic.HttpAllCount)
	}
	c := uiCommon.NewListColumn("http", "HTTP", 8,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}
func (asUI *AppListView) column2XX() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*AppStats).TotalTraffic.Http2xxCount < c2.(*AppStats).TotalTraffic.Http2xxCount
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%8d", appStats.TotalTraffic.Http2xxCount)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http2xxCount)
	}
	c := uiCommon.NewListColumn("2XX", "2XX", 8,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}
func (asUI *AppListView) column3XX() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*AppStats).TotalTraffic.Http3xxCount < c2.(*AppStats).TotalTraffic.Http3xxCount
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%8d", appStats.TotalTraffic.Http3xxCount)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http3xxCount)
	}
	c := uiCommon.NewListColumn("3XX", "3XX", 8,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) column4XX() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*AppStats).TotalTraffic.Http4xxCount < c2.(*AppStats).TotalTraffic.Http4xxCount
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%8d", appStats.TotalTraffic.Http4xxCount)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http4xxCount)
	}
	c := uiCommon.NewListColumn("4XX", "4XX", 8,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) column5XX() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*AppStats).TotalTraffic.Http5xxCount < c2.(*AppStats).TotalTraffic.Http5xxCount
	}
	displayFunc := func(statIndex int, isSelected bool) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%8d", appStats.TotalTraffic.Http5xxCount)
	}
	rawValueFunc := func(statIndex int) string {
		appStats := asUI.displayedSortedStatList[statIndex]
		return fmt.Sprintf("%v", appStats.TotalTraffic.Http5xxCount)
	}
	c := uiCommon.NewListColumn("5XX", "5XX", 8,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *AppListView) Start() {
	go asUI.loadCacheAtStartup()
}

func (asUI *AppListView) loadCacheAtStartup() {
	asUI.loadMetadata()
	asUI.seedStatsFromMetadata()
}

func (asUI *AppListView) refreshMetadata(g *gocui.Gui, v *gocui.View) error {
	go asUI.loadCacheAtStartup()
	return nil
}

func (asUI *AppListView) seedStatsFromMetadata() {

	debug.Info("appListView>seedStatsFromMetadata")

	asUI.currentProcessor.mu.Lock()
	defer asUI.currentProcessor.mu.Unlock()

	currentStatsMap := asUI.currentProcessor.AppMap
	for _, app := range metadata.AllApps() {
		appId := app.Guid
		appStats := currentStatsMap[appId]
		if appStats == nil {
			// New app we haven't seen yet
			appStats = NewAppStats(appId)
			currentStatsMap[appId] = appStats
		}
	}
}

func (asUI *AppListView) GetCurrentProcessor() *AppStatsEventProcessor {
	return asUI.currentProcessor
}

func (asUI *AppListView) SetDisplayPaused(paused bool) {
	asUI.displayPaused = paused
	if !paused {
		asUI.updateData()
	}
}

func (asUI *AppListView) GetDisplayPaused() bool {
	return asUI.displayPaused
}

func (w *AppListView) copyAction(g *gocui.Gui, v *gocui.View) error {

	selectedAppId := w.appListWidget.HighlightKey()
	if selectedAppId == "" {
		// Nothing selected
		return nil
	}
	menuItems := make([]*uiCommon.MenuItem, 0, 5)
	menuItems = append(menuItems, uiCommon.NewMenuItem("cftarget", "cf target"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("cfapp", "cf app"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("cfscale", "cf scale"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("appguid", "app guid"))
	clipboardView := uiCommon.NewSelectMenuWidget(w.masterUI, "clipboardView", "Copy to Clipboard", menuItems, w.clipboardCallback)

	w.masterUI.LayoutManager().Add(clipboardView)
	w.masterUI.SetCurrentViewOnTop(g, "clipboardView")

	//return asUI.displayView.RefreshDisplay(g)
	return nil
}

func (w *AppListView) clipboardCallback(g *gocui.Gui, v *gocui.View, menuId string) error {

	clipboardValue := "hello from clipboard" + time.Now().Format("01-02-2006 15:04:05")

	selectedAppId := w.appListWidget.HighlightKey()
	statsMap := w.displayedProcessor.AppMap
	appStats := statsMap[selectedAppId]
	if appStats == nil {
		// Nothing selected
		return nil
	}
	switch menuId {
	case "cftarget":
		clipboardValue = fmt.Sprintf("cf target -o %v -s %v", appStats.OrgName, appStats.SpaceName)
	case "cfapp":
		clipboardValue = fmt.Sprintf("cf app %v", appStats.AppName)
	case "cfscale":
		clipboardValue = fmt.Sprintf("cf scale %v ", appStats.AppName)
	case "appguid":
		clipboardValue = selectedAppId
	}
	err := clipboard.WriteAll(clipboardValue)
	if err != nil {
		debug.Error("Copy into Clipboard error: " + err.Error())
	}
	return nil
}

func (asUI *AppListView) ClearStats(g *gocui.Gui, v *gocui.View) error {
	debug.Info("appListView>ClearStats")
	asUI.currentProcessor.Clear()
	asUI.updateData()
	asUI.seedStatsFromMetadata()
	return nil
}

func (asUI *AppListView) UpdateDisplay(g *gocui.Gui) error {
	if !asUI.displayPaused {
		asUI.updateData()
	}
	return asUI.RefreshDisplay(g)
}

func (asUI *AppListView) updateData() {
	asUI.mu.Lock()
	processorCopy := asUI.currentProcessor.Clone()
	asUI.displayedProcessor = processorCopy
	asUI.FilterAndSortData()
	asUI.mu.Unlock()
}

func (asUI *AppListView) FilterAndSortData() {
	if len(asUI.displayedProcessor.AppMap) > 0 {
		statsMap := asUI.displayedProcessor.AppMap
		stats := populateNamesIfNeeded(statsMap)
		asUI.displayedSortedStatList = stats
		stats = asUI.filterData(stats)
		stats = asUI.sortData(stats)
		asUI.displayedSortedStatList = stats
	} else {
		asUI.displayedSortedStatList = nil
	}

}

func (asUI *AppListView) filterData(stats []*AppStats) []*AppStats {
	filteredStats := make([]*AppStats, 0, len(stats))
	for rowIndex, s := range stats {
		if asUI.appListWidget.FilterRow(rowIndex) {
			filteredStats = append(filteredStats, s)
		}
	}
	return filteredStats
}

func (asUI *AppListView) sortData(stats []*AppStats) []*AppStats {
	sortFunctions := asUI.appListWidget.GetSortFunctions()
	stats = getSortedStats(stats, sortFunctions)
	return stats
}

func (asUI *AppListView) RefreshDisplay(g *gocui.Gui) error {

	currentView := asUI.masterUI.GetCurrentView(g)
	currentName := currentView.Name()
	if strings.HasPrefix(currentName, asUI.name) {
		err := asUI.refreshListDisplay(g)
		if err != nil {
			return err
		}
	} else if asUI.appDetailView != nil && strings.HasPrefix(currentName, asUI.appDetailView.name) {
		err := asUI.appDetailView.refreshDisplay(g)
		if err != nil {
			return err
		}
	}
	return asUI.updateHeader(g)
}

func (asUI *AppListView) GetListSize() int {
	return len(asUI.displayedSortedStatList)
}

func (asUI *AppListView) GetUnfilteredListSize() int {
	return len(asUI.displayedProcessor.AppMap)
}

func (asUI *AppListView) GetRowKey(statIndex int) string {
	return asUI.displayedSortedStatList[statIndex].AppId
}

func (asUI *AppListView) PreRowDisplay(statIndex int, isSelected bool) string {
	appStats := asUI.displayedSortedStatList[statIndex]
	v := bytes.NewBufferString("")
	if !isSelected && appStats.TotalTraffic.EventL10Rate > 0 {
		fmt.Fprintf(v, util.BRIGHT_WHITE)
	}
	return v.String()
}

func (asUI *AppListView) refreshListDisplay(g *gocui.Gui) error {
	err := asUI.appListWidget.RefreshDisplay(g)
	if err != nil {
		return err
	}
	//return asUI.updateHeader(g)
	return err
}

func (asUI *AppListView) updateHeader(g *gocui.Gui) error {

	v, err := g.View("summaryView")
	if err != nil {
		return err
	}
	if asUI.displayPaused {
		fmt.Fprintf(v, util.REVERSE_GREEN)
		fmt.Fprintf(v, "\r Display update paused ")
		fmt.Fprintf(v, util.CLEAR)
		return nil
	}

	totalReportingAppInstances := 0
	totalActiveApps := 0
	totalUsedMemoryAppInstances := uint64(0)
	totalUsedDiskAppInstances := uint64(0)
	totalCpuPercentage := float64(0)
	for _, appStats := range asUI.displayedSortedStatList {
		for _, cs := range appStats.ContainerArray {
			if cs != nil && cs.ContainerMetric != nil {
				totalReportingAppInstances++
				totalUsedMemoryAppInstances = totalUsedMemoryAppInstances + *cs.ContainerMetric.MemoryBytes
				totalUsedDiskAppInstances = totalUsedDiskAppInstances + *cs.ContainerMetric.DiskBytes
			}
		}
		totalCpuPercentage = totalCpuPercentage + appStats.TotalCpuPercentage
		if appStats.TotalTraffic.EventL60Rate > 0 {
			totalActiveApps++
		}
	}

	totalUsedMemoryAppInstancesDisplay := "--"
	totalUsedDiskAppInstancesDisplay := "--"
	totalCpuPercentageDisplay := "--"
	if totalReportingAppInstances > 0 {
		totalUsedMemoryAppInstancesDisplay = util.ByteSize(totalUsedMemoryAppInstances).String()
		totalUsedDiskAppInstancesDisplay = util.ByteSize(totalUsedDiskAppInstances).String()
		if totalCpuPercentage >= 100 {
			totalCpuPercentageDisplay = fmt.Sprintf("%.0f%%", totalCpuPercentage)
		} else {
			totalCpuPercentageDisplay = fmt.Sprintf("%.1f%%", totalCpuPercentage)
		}
	}

	cellTotalCPUs := 0
	capacityTotalMemory := int64(0)
	capacityTotalDisk := int64(0)
	for _, cellStats := range asUI.displayedProcessor.CellMap {
		cellTotalCPUs = cellTotalCPUs + cellStats.NumOfCpus
		capacityTotalMemory = capacityTotalMemory + cellStats.CapacityTotalMemory
		capacityTotalDisk = capacityTotalDisk + cellStats.CapacityTotalDisk
	}

	cellTotalCapacityDisplay := "--"
	if cellTotalCPUs > 0 {
		cellTotalCapacityDisplay = fmt.Sprintf("%v%%", (cellTotalCPUs * 100))
	}

	capacityTotalMemoryDisplay := "--"
	if capacityTotalMemory > 0 {
		capacityTotalMemoryDisplay = fmt.Sprintf("%v", util.ByteSize(capacityTotalMemory))
	}
	capacityTotalDiskDisplay := "--"
	if capacityTotalDisk > 0 {
		capacityTotalDiskDisplay = fmt.Sprintf("%v", util.ByteSize(capacityTotalDisk))
	}
	fmt.Fprintf(v, "\r")

	// Active apps are apps that have had go-rounter traffic in last 60 seconds
	// Reporting containers are containers that reported metrics in last 90 seconds
	fmt.Fprintf(v, "CPU:%9v Used,%9v Max,   Apps:%5v Total,%5v Actv,   Cntrs:%5v\n",
		totalCpuPercentageDisplay, cellTotalCapacityDisplay, len(asUI.displayedSortedStatList), totalActiveApps, totalReportingAppInstances)

	displayTotalMem := "--"
	totalMem := metadata.GetTotalMemoryAllStartedApps()
	if totalMem > 0 {
		displayTotalMem = util.ByteSize(totalMem).String()
	}
	fmt.Fprintf(v, "Mem:%9v Used,", totalUsedMemoryAppInstancesDisplay)
	// Total quota memory of all running app instances
	fmt.Fprintf(v, "%9v Max,%9v Rsrvd\n", capacityTotalMemoryDisplay, displayTotalMem)

	displayTotalDisk := "--"
	totalDisk := metadata.GetTotalDiskAllStartedApps()
	if totalMem > 0 {
		displayTotalDisk = util.ByteSize(totalDisk).String()
	}

	fmt.Fprintf(v, "Dsk:%9v Used,", totalUsedDiskAppInstancesDisplay)
	fmt.Fprintf(v, "%9v Max,%9v Rsrvd\n", capacityTotalDiskDisplay, displayTotalDisk)

	return nil
}

func (asUI *AppListView) loadMetadata() {
	debug.Info("appListView>loadMetadata")
	metadata.LoadAppCache(asUI.cliConnection)
	metadata.LoadSpaceCache(asUI.cliConnection)
	metadata.LoadOrgCache(asUI.cliConnection)
}
