package appView

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/atotto/clipboard"
	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/eventdata"
	"github.com/kkellner/cloudfoundry-top-plugin/toplog"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"

	"github.com/kkellner/cloudfoundry-top-plugin/metadata"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/appDetailView"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/dataView"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

const StaleContainerSeconds = 80

type AppListView struct {
	*dataView.DataListView
	displayAppStats  []*displaydata.DisplayAppStats
	isWarmupComplete bool
	// This is a count of the number of apps that do not have
	// the correct number of containers running based on app
	// instance setting
	appsNotInDesiredState int
}

func NewAppListView(masterUI masterUIInterface.MasterUIInterface,
	name string, topMargin, bottomMargin int,
	eventProcessor *eventdata.EventProcessor) *AppListView {

	asUI := &AppListView{}

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("CPU", true),
		uiCommon.NewSortColumn("REQ60", true),
		uiCommon.NewSortColumn("appName", false),
		uiCommon.NewSortColumn("spaceName", false),
		uiCommon.NewSortColumn("orgName", false),
	}

	dataListView := dataView.NewDataListView(masterUI, nil,
		name, topMargin, bottomMargin,
		eventProcessor, asUI.columnDefinitions(),
		defaultSortColumns)

	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.PreRowDisplayCallback = asUI.preRowDisplay
	dataListView.UpdateHeaderCallback = asUI.updateHeader
	dataListView.GetListData = asUI.GetListData

	dataListView.SetTitle("App List")
	dataListView.HelpText = helpText

	asUI.DataListView = dataListView

	return asUI

}

func (asUI *AppListView) initializeCallback(g *gocui.Gui, viewName string) error {

	if err := g.SetKeybinding(viewName, 'c', gocui.ModNone, asUI.copyAction); err != nil {
		log.Panicln(err)
	}

	// TODO: Testing -- remove later
	if err := g.SetKeybinding(viewName, 'z', gocui.ModNone, asUI.testShowUserMessage); err != nil {
		log.Panicln(err)
	}
	// TODO: Testing -- remove later
	if err := g.SetKeybinding(viewName, 'a', gocui.ModNone, asUI.testClearUserMessage); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone, asUI.enterAction); err != nil {
		log.Panicln(err)
	}

	return nil
}

func (asUI *AppListView) enterAction(g *gocui.Gui, v *gocui.View) error {
	highlightKey := asUI.GetListWidget().HighlightKey()
	if asUI.GetListWidget().HighlightKey() != "" {
		topMargin, bottomMargin := asUI.GetMargins()

		detailView := appDetailView.NewAppDetailView(asUI.GetMasterUI(), asUI, "appDetailView",
			topMargin, bottomMargin,
			asUI.GetEventProcessor(),
			highlightKey)
		asUI.SetDetailView(detailView)
		asUI.GetMasterUI().OpenView(g, detailView)
	}
	return nil
}

func (asUI *AppListView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, asUI.columnAppName())
	columns = append(columns, asUI.columnSpaceName())
	columns = append(columns, asUI.columnOrgName())

	columns = append(columns, asUI.columnDesiredInstances())
	columns = append(columns, asUI.columnReportingContainers())

	columns = append(columns, asUI.columnTotalCpu())
	columns = append(columns, asUI.columnTotalMemoryUsed())
	columns = append(columns, asUI.columnTotalDiskUsed())

	columns = append(columns, asUI.columnAvgResponseTimeL60Info())
	columns = append(columns, asUI.columnLogStdout())
	columns = append(columns, asUI.columnLogStderr())

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

func (asUI *AppListView) isUserMessageOpen(g *gocui.Gui) bool {
	alertViewName := "alertView"
	view := asUI.GetMasterUI().LayoutManager().GetManagerByViewName(alertViewName)
	if view != nil {
		return true
	} else {
		return false
	}
}

func (asUI *AppListView) clearUserMessage(g *gocui.Gui) error {
	alertViewName := "alertView"
	view := asUI.GetMasterUI().LayoutManager().GetManagerByViewName(alertViewName)
	if view != nil {
		asUI.GetMasterUI().CloseView(view)
	}
	topMargin := asUI.GetTopMargin()
	asUI.SetTopMarginOnListWidget(topMargin)
	return nil
}

// TODO: Have message levels which will colorize differently
func (asUI *AppListView) showUserMessage(g *gocui.Gui, message string) error {
	alertViewName := "alertView"
	alertHeight := 1

	topMargin := asUI.GetTopMargin()
	asUI.SetTopMarginOnListWidget(topMargin + alertHeight)

	var alertView *uiCommon.AlertWidget
	view := asUI.GetMasterUI().LayoutManager().GetManagerByViewName(alertViewName)
	if view == nil {
		alertView = uiCommon.NewAlertWidget(alertViewName, topMargin, alertHeight)
		asUI.GetMasterUI().LayoutManager().Add(alertView)
	} else {
		alertView = view.(*uiCommon.AlertWidget)
		alertView.SetHeight(alertHeight)
	}
	alertView.SetMessage(message)
	return nil
}

func (asUI *AppListView) testShowUserMessage(g *gocui.Gui, v *gocui.View) error {
	return asUI.showUserMessage(g, "ALERT: 1 application(s) not in desired state (row colored red) ")
}

func (asUI *AppListView) testClearUserMessage(g *gocui.Gui, v *gocui.View) error {
	return asUI.clearUserMessage(g)
}

func (asUI *AppListView) copyAction(g *gocui.Gui, v *gocui.View) error {

	selectedAppId := asUI.GetListWidget().HighlightKey()
	if selectedAppId == "" {
		// Nothing selected
		return nil
	}
	menuItems := make([]*uiCommon.MenuItem, 0, 5)
	menuItems = append(menuItems, uiCommon.NewMenuItem("cftarget", "cf target"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("cfapp", "cf app"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("cfscale", "cf scale"))
	menuItems = append(menuItems, uiCommon.NewMenuItem("appguid", "app guid"))
	masterUI := asUI.GetMasterUI()
	clipboardView := uiCommon.NewSelectMenuWidget(masterUI, "clipboardView", "Copy to Clipboard", menuItems, asUI.clipboardCallback)

	masterUI.LayoutManager().Add(clipboardView)
	masterUI.SetCurrentViewOnTop(g)
	return nil
}

func (asUI *AppListView) clipboardCallback(g *gocui.Gui, v *gocui.View, menuId string) error {

	clipboardValue := ""

	selectedAppId := asUI.GetListWidget().HighlightKey()
	statsMap := asUI.GetDisplayedEventData().AppMap
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

func (asUI *AppListView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *AppListView) postProcessData() []*displaydata.DisplayAppStats {

	displayStatsArray := make([]*displaydata.DisplayAppStats, 0)
	appMap := asUI.GetDisplayedEventData().AppMap
	appStatsArray := eventdata.PopulateNamesFromMap(appMap)
	appsNotInDesiredState := 0

	for _, appStats := range appStatsArray {
		displayAppStats := displaydata.NewDisplayAppStats(appStats)
		displayStatsArray = append(displayStatsArray, displayAppStats)
		appMetadata := metadata.FindAppMetadata(appStats.AppId)

		totalCpuPercentage := 0.0
		totalUsedMemory := uint64(0)
		totalUsedDisk := uint64(0)
		totalReportingContainers := 0

		if appMetadata.State == "STARTED" {
			displayAppStats.DesiredContainers = int(appMetadata.Instances)
		}

		now := time.Now()
		for containerIndex, cs := range appStats.ContainerArray {
			if cs != nil && cs.ContainerMetric != nil {
				// If we haven't gotten a container update recently, ignore the old value
				if now.Sub(cs.LastUpdate) > time.Second*StaleContainerSeconds {
					appStats.ContainerArray[containerIndex] = nil
					continue
				}

				totalCpuPercentage = totalCpuPercentage + *cs.ContainerMetric.CpuPercentage
				totalUsedMemory = totalUsedMemory + *cs.ContainerMetric.MemoryBytes
				totalUsedDisk = totalUsedDisk + *cs.ContainerMetric.DiskBytes
				totalReportingContainers++
			}
		}
		if totalReportingContainers < displayAppStats.DesiredContainers {
			appsNotInDesiredState = appsNotInDesiredState + 1
		}

		displayAppStats.TotalCpuPercentage = totalCpuPercentage
		displayAppStats.TotalUsedMemory = totalUsedMemory
		displayAppStats.TotalUsedDisk = totalUsedDisk
		displayAppStats.TotalReportingContainers = totalReportingContainers

		logStdoutCount := int64(0)
		logStderrCount := int64(0)
		for _, cs := range appStats.ContainerArray {
			if cs != nil {
				logStdoutCount = logStdoutCount + cs.OutCount
				logStderrCount = logStderrCount + cs.ErrCount
			}
		}
		displayAppStats.TotalLogStdout = logStdoutCount + appStats.NonContainerStdout
		displayAppStats.TotalLogStderr = logStderrCount + appStats.NonContainerStderr

	}
	asUI.displayAppStats = displayStatsArray
	asUI.isWarmupComplete = asUI.GetMasterUI().IsWarmupComplete()
	asUI.appsNotInDesiredState = appsNotInDesiredState
	return displayStatsArray
}

func (asUI *AppListView) convertToListData(statsArray []*displaydata.DisplayAppStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, len(statsArray))
	for i, d := range statsArray {
		listData[i] = d
	}
	return listData
}

func (asUI *AppListView) detailViewClosed(g *gocui.Gui) error {
	asUI.DataListView.RefreshDisplayCallback = nil
	return asUI.RefreshDisplay(g)
}

func (asUI *AppListView) preRowDisplay(data uiCommon.IData, isSelected bool) string {
	appStats := data.(*displaydata.DisplayAppStats)
	v := bytes.NewBufferString("")
	if asUI.isWarmupComplete && appStats.DesiredContainers > appStats.TotalReportingContainers {
		if isSelected {
			fmt.Fprintf(v, util.RED_TEXT_GREEN_BG)
		} else {
			fmt.Fprintf(v, util.BRIGHT_RED)
		}
	} else if !isSelected && appStats.TotalTraffic.EventL10Rate > 0 {
		fmt.Fprintf(v, util.BRIGHT_WHITE)
	}
	return v.String()
}

func (asUI *AppListView) checkForAlerts(g *gocui.Gui) error {
	if asUI.isWarmupComplete && asUI.appsNotInDesiredState > 0 {
		plural := ""
		if asUI.appsNotInDesiredState > 1 {
			plural = "s"
		}
		msg := fmt.Sprintf("ALERT: %v application%v not in desired state (row%v colored red) ",
			asUI.appsNotInDesiredState, plural, plural)
		asUI.showUserMessage(g, msg)
	} else if asUI.isUserMessageOpen(g) {
		asUI.clearUserMessage(g)
	}
	return nil
}

func (asUI *AppListView) updateHeader(g *gocui.Gui, v *gocui.View) error {

	// TODO: Is this the best stop to check for alerts?? Seems out of place in the updateHeader method
	asUI.checkForAlerts(g)

	totalReportingAppInstances := 0
	totalActiveApps := 0
	totalUsedMemoryAppInstances := uint64(0)
	totalUsedDiskAppInstances := uint64(0)
	totalCpuPercentage := float64(0)
	for _, appStats := range asUI.displayAppStats {
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
