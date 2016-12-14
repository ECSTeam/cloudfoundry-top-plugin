package appDetailView

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/eventdata"
	"github.com/kkellner/cloudfoundry-top-plugin/metadata"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/dataView"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

const MEGABYTE = (1024 * 1024)

type AppDetailView struct {
	*dataView.DataListView
	appId              string
	requestsInfoWidget *RequestsInfoWidget
}

func NewAppDetailView(masterUI masterUIInterface.MasterUIInterface,
	parentView dataView.DataListViewInterface,
	name string,
	topMargin, bottomMargin int,
	eventProcessor *eventdata.EventProcessor,
	appId string) *AppDetailView {

	asUI := &AppDetailView{appId: appId}

	asUI.requestsInfoWidget = NewRequestsInfoWidget(masterUI, "requestsInfoWidget")
	masterUI.LayoutManager().Add(asUI.requestsInfoWidget)

	defaultSortColumns := []*uiCommon.SortColumn{
	//uiCommon.NewSortColumn("CPU_PERCENT", true),
	//uiCommon.NewSortColumn("IP", false),
	}

	dataListView := dataView.NewDataListView(masterUI, parentView,
		name, topMargin+16, bottomMargin,
		eventProcessor, asUI.columnDefinitions(),
		defaultSortColumns)

	//dataListView.UpdateHeaderCallback = asUI.updateHeader
	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.GetListData = asUI.GetListData
	dataListView.RefreshDisplayCallback = asUI.refreshDisplayX

	//dataListView.SetTitle("App Details (press 'q' to quit view)")
	dataListView.HelpText = helpText

	asUI.DataListView = dataListView

	return asUI
}

func (asUI *AppDetailView) initializeCallback(g *gocui.Gui, viewName string) error {
	if err := g.SetKeybinding(viewName, 'q', gocui.ModNone, asUI.closeAppDetailView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, gocui.KeyEsc, gocui.ModNone, asUI.closeAppDetailView); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (asUI *AppDetailView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	// TODO: Replace this with container specific columns
	columns = append(columns, asUI.columnAppName())
	columns = append(columns, asUI.columnSpaceName())
	columns = append(columns, asUI.columnOrgName())
	return columns
}

func (asUI *AppDetailView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *AppDetailView) postProcessData() []*displaydata.DisplayContainerStats {

	containerStatsArray := make([]*displaydata.DisplayContainerStats, 0)

	appMap := asUI.GetDisplayedEventData().AppMap
	appStats := appMap[asUI.appId]
	for _, containerStats := range appStats.ContainerArray {
		if containerStats != nil {
			containerStats := displaydata.NewDisplayContainerStats(containerStats, appStats)
			containerStatsArray = append(containerStatsArray, containerStats)
		}
	}

	return containerStatsArray
}

func (asUI *AppDetailView) convertToListData(containerStatsArray []*displaydata.DisplayContainerStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, 0, len(containerStatsArray))
	for _, d := range containerStatsArray {
		listData = append(listData, d)
	}
	return listData
}

/*
func (asUI *AppDetailView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *AppDetailView) postProcessData() []*eventdata.AppStats {

	// TODO: Move most of the clone() code here

	appMap := asUI.GetDisplayedEventData().AppMap
	if len(appMap) > 0 {
		stats := eventdata.PopulateNamesIfNeeded(appMap)
		return stats
	} else {
		return nil
	}
}

func (asUI *AppDetailView) convertToListData(statsList []*eventdata.AppStats) []uiCommon.IData {
	listData := make([]uiCommon.IData, len(statsList))
	for i, d := range statsList {
		listData[i] = d
	}
	return listData
}
*/

/*
func (w *AppDetailView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	bottom := maxY - w.bottomMargin
	if w.topMargin >= bottom {
		bottom = w.topMargin + 1
	}

	leftMargin := 0
	rightMargin := maxX - 1
	v, err := g.SetView(w.name, leftMargin, w.topMargin, rightMargin, bottom)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return merry.Wrap(err).Appendf("viewName:[%v] left:%v, top:%v, right:%v, bottom: %v",
				w.name, leftMargin, w.topMargin, rightMargin, bottom)
		}
		v.Title = "App Details (press 'q' to quit view)"
		v.Frame = true

		if err := g.SetKeybinding(w.name, 'q', gocui.ModNone, w.closeAppDetailView); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.closeAppDetailView); err != nil {
			return err
		}

		if err := g.SetKeybinding(w.name, 'i', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			appInfoWidget := NewAppInfoWidget(w.masterUI, "appInfoWidget", 70, 20)
			w.masterUI.LayoutManager().Add(appInfoWidget)
			w.masterUI.SetCurrentViewOnTop(g)
			return nil
		}); err != nil {
			log.Panicln(err)
		}

		if err := w.masterUI.SetCurrentViewOnTop(g); err != nil {
			log.Panicln(err)
		}

		w.refreshDisplay(g)

	}
	return nil
}
*/

func (asUI *AppDetailView) closeAppDetailView(g *gocui.Gui, v *gocui.View) error {
	if err := asUI.GetMasterUI().CloseView(asUI); err != nil {
		return err
	}

	if err := asUI.GetMasterUI().CloseView(asUI.requestsInfoWidget); err != nil {
		return err
	}
	return nil
}

func (w *AppDetailView) refreshDisplayX(g *gocui.Gui) error {

	v, err := g.View("requestsInfoWidget")
	if err != nil {
		return err
	}

	v.Clear()

	if w.appId == "" {
		fmt.Fprintln(v, "No application selected")
		return nil
	}

	m := w.GetDisplayedEventData().AppMap

	appStats := m[w.appId]

	avgResponseTimeL60Info := "--"
	if appStats.TotalTraffic.AvgResponseL60Time >= 0 {
		avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL60Time / 1000000
		avgResponseTimeL60Info = fmt.Sprintf("%8.1f", avgResponseTimeMs)
	}

	avgResponseTimeL10Info := "--"
	if appStats.TotalTraffic.AvgResponseL10Time >= 0 {
		avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL10Time / 1000000
		avgResponseTimeL10Info = fmt.Sprintf("%8.1f", avgResponseTimeMs)
	}

	avgResponseTimeL1Info := "--"
	if appStats.TotalTraffic.AvgResponseL1Time >= 0 {
		avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL1Time / 1000000
		avgResponseTimeL1Info = fmt.Sprintf("%8.1f", avgResponseTimeMs)
	}

	appMetadata := metadata.FindAppMetadata(appStats.AppId)

	memoryDisplay := "--"
	totalMemoryDisplay := "--"
	totalDiskDisplay := "--"
	instancesDisplay := "--"
	diskQuotaDisplay := "--"
	if appMetadata.Guid != "" {
		memoryDisplay = util.ByteSize(appMetadata.MemoryMB * MEGABYTE).String()
		diskQuotaDisplay = util.ByteSize(appMetadata.DiskQuotaMB * MEGABYTE).String()
		instancesDisplay = fmt.Sprintf("%v", appMetadata.Instances)
		totalMemoryDisplay = util.ByteSize((appMetadata.MemoryMB * MEGABYTE) * appMetadata.Instances).String()
		totalDiskDisplay = util.ByteSize((appMetadata.DiskQuotaMB * MEGABYTE) * appMetadata.Instances).String()
	}

	fmt.Fprintf(v, " \n")
	fmt.Fprintf(v, "App Name:        %v%v%v\n", util.BRIGHT_WHITE, appStats.AppName, util.CLEAR)
	fmt.Fprintf(v, "AppId:           %v\n", appStats.AppId)
	fmt.Fprintf(v, "AppUUID:         %v\n", appStats.AppUUID)
	fmt.Fprintf(v, "Space:           %v\n", appStats.SpaceName)
	fmt.Fprintf(v, "Organization:    %v\n", appStats.OrgName)
	fmt.Fprintf(v, "Desired insts:   %v\n", instancesDisplay)
	fmt.Fprintf(v, "Rsrvd mem per (all):  %8v (%8v)\n", memoryDisplay, totalMemoryDisplay)
	fmt.Fprintf(v, "Rsrvd disk per (all): %8v (%8v)\n", diskQuotaDisplay, totalDiskDisplay)
	//fmt.Fprintf(v, "%v", util.BRIGHT_WHITE)
	//fmt.Fprintf(v, "%v", util.CLEAR)
	fmt.Fprintf(v, "\n")

	fmt.Fprintf(v, "%22v", "")
	fmt.Fprintf(v, "    1sec   10sec   60sec\n")

	fmt.Fprintf(v, "%22v", "HTTP(S) Event Rate:")
	fmt.Fprintf(v, "%8v", appStats.TotalTraffic.EventL1Rate)
	fmt.Fprintf(v, "%8v", appStats.TotalTraffic.EventL10Rate)
	fmt.Fprintf(v, "%8v\n", appStats.TotalTraffic.EventL60Rate)

	fmt.Fprintf(v, "%22v", "Avg Rspnse Time(ms):")
	fmt.Fprintf(v, "%8v", avgResponseTimeL1Info)
	fmt.Fprintf(v, "%8v", avgResponseTimeL10Info)
	fmt.Fprintf(v, "%8v\n", avgResponseTimeL60Info)

	fmt.Fprintf(v, "\n")
	fmt.Fprintf(v, "HTTP(S) status code:\n")
	fmt.Fprintf(v, "  2xx: %12v\n", util.Format(appStats.TotalTraffic.Http2xxCount))
	fmt.Fprintf(v, "  3xx: %12v\n", util.Format(appStats.TotalTraffic.Http3xxCount))
	fmt.Fprintf(v, "  4xx: %12v\n", util.Format(appStats.TotalTraffic.Http4xxCount))
	fmt.Fprintf(v, "  5xx: %12v\n", util.Format(appStats.TotalTraffic.Http5xxCount))
	fmt.Fprintf(v, "%v", util.BRIGHT_WHITE)
	fmt.Fprintf(v, "  All: %12v\n", util.Format(appStats.TotalTraffic.HttpAllCount))
	fmt.Fprintf(v, "%v", util.CLEAR)

	fmt.Fprintf(v, "\nContainer Info:\n")
	fmt.Fprintf(v, "%5v", "INST")
	fmt.Fprintf(v, "%8v", "CPU%")
	fmt.Fprintf(v, "%12v", "MEM USED")
	fmt.Fprintf(v, "%12v", "MEM FREE")
	fmt.Fprintf(v, "%12v", "DISK USED")
	fmt.Fprintf(v, "%12v", "DISK FREE")
	fmt.Fprintf(v, "%12v", "STDOUT")
	fmt.Fprintf(v, "%12v", "STDERR")
	fmt.Fprintf(v, "\n")

	totalCpuPercentage := 0.0
	totalMemory := uint64(0)
	totalDisk := uint64(0)
	reportingAppInstances := 0
	totalLogCount := int64(0)
	for i, ca := range appStats.ContainerArray {
		if ca != nil {

			fmt.Fprintf(v, "%5v", i)
			if ca.ContainerMetric != nil {
				cpuPercentage := *ca.ContainerMetric.CpuPercentage
				totalCpuPercentage = totalCpuPercentage + cpuPercentage
				memory := *ca.ContainerMetric.MemoryBytes
				totalMemory = totalMemory + memory
				disk := *ca.ContainerMetric.DiskBytes
				totalDisk = totalDisk + disk
				fmt.Fprintf(v, "%8.2f", cpuPercentage)
				fmt.Fprintf(v, "%12v", util.ByteSize(memory))
				fmt.Fprintf(v, "%12v", util.ByteSize((uint64(appMetadata.MemoryMB)*MEGABYTE)-memory))
				fmt.Fprintf(v, "%12v", util.ByteSize(disk))
				fmt.Fprintf(v, "%12v", util.ByteSize((uint64(appMetadata.DiskQuotaMB)*MEGABYTE)-disk))
			} else {
				fmt.Fprintf(v, "%8v", "--")
				fmt.Fprintf(v, "%12v", "--")
				fmt.Fprintf(v, "%12v", "--")
				fmt.Fprintf(v, "%12v", "--")
				fmt.Fprintf(v, "%12v", "--")
			}

			fmt.Fprintf(v, "%12v", util.Format(ca.OutCount))
			fmt.Fprintf(v, "%12v", util.Format(ca.ErrCount))
			totalLogCount = totalLogCount + (ca.OutCount + ca.ErrCount)

			fmt.Fprintf(v, "\n")
			reportingAppInstances++
		}
	}
	//fmt.Fprintf(v, "\n")
	if reportingAppInstances == 0 {
		fmt.Fprintf(v, "%6v", "\n Waiting for container metrics...")
	} else {
		if totalMemory > 0 {
			fmt.Fprintf(v, "%v", util.BRIGHT_WHITE)
			fmt.Fprintf(v, "Total")
			fmt.Fprintf(v, "%8.2f", totalCpuPercentage)
			fmt.Fprintf(v, "%12v", util.ByteSize(totalMemory))
			fmt.Fprintf(v, "%12v", util.ByteSize(totalDisk))
			fmt.Fprintf(v, "%v", util.CLEAR)
		}
	}
	fmt.Fprintf(v, "\n\n")
	totalLogCount = totalLogCount + appStats.NonContainerOutCount + appStats.NonContainerErrCount
	fmt.Fprintf(v, "Non container logs - Stdout: %-12v ", util.Format(appStats.NonContainerOutCount))
	fmt.Fprintf(v, "Stderr: %-12v\n", util.Format(appStats.NonContainerErrCount))
	fmt.Fprintf(v, "Total log events: %12v\n", util.Format(totalLogCount))

	return nil
}
