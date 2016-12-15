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
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

type AppListView struct {
	*dataView.DataListView
	//appDetailView *AppDetailView
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

		//asUI.GetMasterUI().LayoutManager().Add(asUI.appDetailView)
		//asUI.GetMasterUI().SetCurrentViewOnTop(g)
		//asUI.DataListView.RefreshDisplayCallback = asUI.refreshDetailView
	}
	return nil
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

	clipboardValue := "hello from clipboard" + time.Now().Format("01-02-2006 15:04:05")

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

func (asUI *AppListView) postProcessData() []*eventdata.AppStats {

	// TODO: Move most of the clone() code here

	appMap := asUI.GetDisplayedEventData().AppMap
	if len(appMap) > 0 {
		stats := eventdata.PopulateNamesFromMap(appMap)
		return stats
	} else {
		return nil
	}
}

func (asUI *AppListView) convertToListData(statsList []*eventdata.AppStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, len(statsList))
	for i, d := range statsList {
		listData[i] = d
	}
	return listData
}

func (asUI *AppListView) detailViewClosed(g *gocui.Gui) error {
	asUI.DataListView.RefreshDisplayCallback = nil
	return asUI.RefreshDisplay(g)
}

func (asUI *AppListView) preRowDisplay(data uiCommon.IData, isSelected bool) string {
	appStats := data.(*eventdata.AppStats)
	v := bytes.NewBufferString("")
	if !isSelected && appStats.TotalTraffic.EventL10Rate > 0 {
		fmt.Fprintf(v, util.BRIGHT_WHITE)
	}
	return v.String()
}

func (asUI *AppListView) updateHeader(g *gocui.Gui, v *gocui.View) error {

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
