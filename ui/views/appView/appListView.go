package appView

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/eventdata"
	"github.com/kkellner/cloudfoundry-top-plugin/toplog"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"

	"github.com/kkellner/cloudfoundry-top-plugin/metadata"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/helpView"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

type AppListView struct {
	masterUI       masterUIInterface.MasterUIInterface
	name           string
	topMargin      int
	bottomMargin   int
	eventProcessor *eventdata.EventProcessor
	mu             sync.Mutex
	appDetailView  *AppDetailView
	appListWidget  *uiCommon.ListWidget
	displayPaused  bool
}

func NewAppListView(masterUI masterUIInterface.MasterUIInterface,
	name string, topMargin, bottomMargin int,
	eventProcessor *eventdata.EventProcessor) *AppListView {

	return &AppListView{
		masterUI:       masterUI,
		name:           name,
		topMargin:      topMargin,
		bottomMargin:   bottomMargin,
		eventProcessor: eventProcessor,
	}
}

func (asUI *AppListView) Name() string {
	return asUI.name
}

func (asUI *AppListView) Layout(g *gocui.Gui) error {

	if asUI.appListWidget == nil {

		statList := asUI.postProcessData(asUI.GetDisplayedEventData().AppMap)
		listData := asUI.convertToListData(statList)

		appListWidget := uiCommon.NewListWidget(asUI.masterUI, asUI.name,
			asUI.topMargin, asUI.bottomMargin, asUI, asUI.columnDefinitions(),
			listData)
		appListWidget.Title = "App List"
		appListWidget.PreRowDisplayFunc = asUI.PreRowDisplay

		defaultSortColums := []*uiCommon.SortColumn{
			uiCommon.NewSortColumn("CPU", true),
			uiCommon.NewSortColumn("REQ60", true),
			uiCommon.NewSortColumn("appName", false),
			uiCommon.NewSortColumn("spaceName", false),
			uiCommon.NewSortColumn("orgName", false),
		}
		appListWidget.SetSortColumns(defaultSortColums)

		asUI.appListWidget = appListWidget
		if err := g.SetKeybinding(asUI.name, 'h', gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				helpView := helpView.NewHelpView(asUI.masterUI, "helpView", 75, 17, helpText)
				asUI.masterUI.LayoutManager().Add(helpView)
				asUI.masterUI.SetCurrentViewOnTop(g)
				return nil
			}); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(asUI.name, 'c', gocui.ModNone, asUI.copyAction); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding(asUI.name, gocui.KeyEnter, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				if asUI.appListWidget.HighlightKey() != "" {
					asUI.appDetailView = NewAppDetailView(asUI.masterUI, "appDetailView", asUI.appListWidget.HighlightKey(), asUI)
					asUI.masterUI.LayoutManager().Add(asUI.appDetailView)
					asUI.masterUI.SetCurrentViewOnTop(g)
				}
				return nil
			}); err != nil {
			log.Panicln(err)
		}

	}
	return asUI.appListWidget.Layout(g)
}

func (asUI *AppListView) convertToListData(statsList []*eventdata.AppStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, len(statsList))
	for i, d := range statsList {
		listData[i] = d
	}
	return listData
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

	columns = append(columns, asUI.columnReq1())
	columns = append(columns, asUI.columnReq10())
	columns = append(columns, asUI.columnReq60())

	columns = append(columns, asUI.columnTotalReq())
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

func (asUI *AppListView) refreshMetadata(g *gocui.Gui, v *gocui.View) error {
	go asUI.eventProcessor.LoadCacheAndSeeData()
	return nil
}

func (asUI *AppListView) GetCurrentEventData() *eventdata.EventData {
	return asUI.eventProcessor.GetCurrentEventData()
}

func (asUI *AppListView) GetDisplayedEventData() *eventdata.EventData {
	return asUI.eventProcessor.GetDisplayedEventData()
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
	w.masterUI.SetCurrentViewOnTop(g)
	return nil
}

func (w *AppListView) clipboardCallback(g *gocui.Gui, v *gocui.View, menuId string) error {

	clipboardValue := "hello from clipboard" + time.Now().Format("01-02-2006 15:04:05")

	selectedAppId := w.appListWidget.HighlightKey()
	statsMap := w.GetDisplayedEventData().AppMap
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
		toplog.Error("Copy into Clipboard error: " + err.Error())
	}
	return nil
}

func (asUI *AppListView) UpdateDisplay(g *gocui.Gui) error {
	if !asUI.displayPaused {
		asUI.updateData()
	}
	return asUI.RefreshDisplay(g)
}

// XXX
func (asUI *AppListView) updateData() {

	asUI.eventProcessor.UpdateData()
	processor := asUI.GetDisplayedEventData()
	statList := asUI.postProcessData(processor.AppMap)
	listData := asUI.convertToListData(statList)
	asUI.appListWidget.SetListData(listData)

}

func (asUI *AppListView) postProcessData(statsMap map[string]*eventdata.AppStats) []*eventdata.AppStats {
	if len(statsMap) > 0 {
		stats := eventdata.PopulateNamesIfNeeded(statsMap)
		return stats
	} else {
		return nil
	}
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

func (asUI *AppListView) PreRowDisplay(data uiCommon.IData, isSelected bool) string {
	appStats := data.(*eventdata.AppStats)
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
	return err
}

func (asUI *AppListView) updateHeader(g *gocui.Gui) error {

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

	totalReportingAppInstances := 0
	totalActiveApps := 0
	totalUsedMemoryAppInstances := uint64(0)
	totalUsedDiskAppInstances := uint64(0)
	totalCpuPercentage := float64(0)
	for _, appStats := range asUI.GetDisplayedEventData().AppMap {
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
		totalUsedMemoryAppInstancesDisplay = util.ByteSize(totalUsedMemoryAppInstances).StringWithPrecision(0)
		totalUsedDiskAppInstancesDisplay = util.ByteSize(totalUsedDiskAppInstances).StringWithPrecision(0)
		if totalCpuPercentage >= 100 {
			totalCpuPercentageDisplay = fmt.Sprintf("%.0f%%", totalCpuPercentage)
		} else {
			totalCpuPercentageDisplay = fmt.Sprintf("%.1f%%", totalCpuPercentage)
		}
	}

	cellTotalCPUs := 0
	capacityTotalMemory := int64(0)
	capacityTotalDisk := int64(0)
	for _, cellStats := range asUI.GetDisplayedEventData().CellMap {
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
		capacityTotalMemoryDisplay = fmt.Sprintf("%v", util.ByteSize(capacityTotalMemory).StringWithPrecision(0))
	}
	capacityTotalDiskDisplay := "--"
	if capacityTotalDisk > 0 {
		capacityTotalDiskDisplay = fmt.Sprintf("%v", util.ByteSize(capacityTotalDisk).StringWithPrecision(0))
	}
	fmt.Fprintf(v, "\r")

	// Active apps are apps that have had go-rounter traffic in last 60 seconds
	// Reporting containers are containers that reported metrics in last 90 seconds
	fmt.Fprintf(v, "CPU:%6v Used,%6v Max,         Apps:%5v Total,%5v Actv,   Cntrs:%5v\n",
		totalCpuPercentageDisplay, cellTotalCapacityDisplay,
		len(asUI.GetDisplayedEventData().AppMap),
		totalActiveApps, totalReportingAppInstances)

	displayTotalMem := "--"
	totalMem := metadata.GetTotalMemoryAllStartedApps()
	if totalMem > 0 {
		displayTotalMem = util.ByteSize(totalMem).StringWithPrecision(0)
	}
	fmt.Fprintf(v, "Mem:%6v Used,", totalUsedMemoryAppInstancesDisplay)
	// Total quota memory of all running app instances
	fmt.Fprintf(v, "%6v Max,%6v Rsrvd\n", capacityTotalMemoryDisplay, displayTotalMem)

	displayTotalDisk := "--"
	totalDisk := metadata.GetTotalDiskAllStartedApps()
	if totalMem > 0 {
		displayTotalDisk = util.ByteSize(totalDisk).StringWithPrecision(0)
	}

	fmt.Fprintf(v, "Dsk:%6v Used,", totalUsedDiskAppInstancesDisplay)
	fmt.Fprintf(v, "%6v Max,%6v Rsrvd\n", capacityTotalDiskDisplay, displayTotalDisk)

	return nil
}
