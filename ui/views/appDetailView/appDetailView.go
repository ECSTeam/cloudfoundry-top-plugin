package appDetailView

import (
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

	requestViewHeight := 5

	asUI.requestsInfoWidget = NewRequestsInfoWidget(masterUI, "requestsInfoWidget", requestViewHeight, asUI)
	masterUI.LayoutManager().Add(asUI.requestsInfoWidget)

	defaultSortColumns := []*uiCommon.SortColumn{
		uiCommon.NewSortColumn("IDX", false),
	}

	dataListView := dataView.NewDataListView(masterUI, parentView,
		name, topMargin+requestViewHeight+1, bottomMargin,
		eventProcessor, asUI.columnDefinitions(),
		defaultSortColumns)

	//dataListView.UpdateHeaderCallback = asUI.updateHeader
	dataListView.InitializeCallback = asUI.initializeCallback
	dataListView.GetListData = asUI.GetListData
	dataListView.RefreshDisplayCallback = asUI.refreshDisplay

	dataListView.SetTitle("Container List")
	dataListView.HelpText = helpText

	asUI.DataListView = dataListView

	return asUI
}

func (asUI *AppDetailView) initializeCallback(g *gocui.Gui, viewName string) error {
	if err := g.SetKeybinding(viewName, 'x', gocui.ModNone, asUI.closeAppDetailView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, gocui.KeyEsc, gocui.ModNone, asUI.closeAppDetailView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(viewName, 'i', gocui.ModNone, asUI.openInfoAction); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (asUI *AppDetailView) openInfoAction(g *gocui.Gui, v *gocui.View) error {
	appInfoWidget := NewAppInfoWidget(asUI.GetMasterUI(), "appInfoWidget", 70, 20, asUI)
	asUI.GetMasterUI().LayoutManager().Add(appInfoWidget)
	asUI.GetMasterUI().SetCurrentViewOnTop(g)
	return nil
}

func (asUI *AppDetailView) columnDefinitions() []*uiCommon.ListColumn {
	columns := make([]*uiCommon.ListColumn, 0)
	columns = append(columns, ColumnContainerIndex())
	columns = append(columns, ColumnTotalCpuPercentage())
	columns = append(columns, ColumnMemoryUsed())
	columns = append(columns, ColumnMemoryFree())
	columns = append(columns, ColumnDiskUsed())
	columns = append(columns, ColumnDiskFree())
	columns = append(columns, ColumnLogStdout())
	columns = append(columns, ColumnLogStderr())

	return columns
}

func (asUI *AppDetailView) GetListData() []uiCommon.IData {
	displayDataList := asUI.postProcessData()
	listData := asUI.convertToListData(displayDataList)
	return listData
}

func (asUI *AppDetailView) postProcessData() []*displaydata.DisplayContainerStats {

	displayStatsArray := make([]*displaydata.DisplayContainerStats, 0)

	appMap := asUI.GetDisplayedEventData().AppMap
	appStats := appMap[asUI.appId]

	eventdata.PopulateNamesIfNeeded(appStats)
	appMetadata := metadata.FindAppMetadata(appStats.AppId)

	for _, containerStats := range appStats.ContainerArray {
		if containerStats != nil {
			displayContainerStats := displaydata.NewDisplayContainerStats(containerStats, appStats)
			usedMemory := containerStats.ContainerMetric.GetMemoryBytes()
			reservedMemory := uint64(appMetadata.MemoryMB) * util.MEGABYTE
			freeMemory := reservedMemory - usedMemory
			displayContainerStats.FreeMemory = freeMemory
			displayContainerStats.ReservedMemory = reservedMemory
			usedDisk := containerStats.ContainerMetric.GetDiskBytes()
			reservedDisk := uint64(appMetadata.DiskQuotaMB) * util.MEGABYTE
			freeDisk := reservedDisk - usedDisk
			displayContainerStats.FreeDisk = freeDisk
			displayContainerStats.ReservedDisk = reservedDisk
			displayStatsArray = append(displayStatsArray, displayContainerStats)
		}
	}
	return displayStatsArray
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

func (w *AppDetailView) refreshDisplay(g *gocui.Gui) error {

	// HTTP request stats -- These stands are also on the appListView so we need them in a detail view??
	/*
		fmt.Fprintf(v, "\n")
		fmt.Fprintf(v, "HTTP(S) status code:\n")
		fmt.Fprintf(v, "  2xx: %12v\n", util.Format(appStats.TotalTraffic.Http2xxCount))
		fmt.Fprintf(v, "  3xx: %12v\n", util.Format(appStats.TotalTraffic.Http3xxCount))
		fmt.Fprintf(v, "  4xx: %12v\n", util.Format(appStats.TotalTraffic.Http4xxCount))
		fmt.Fprintf(v, "  5xx: %12v\n", util.Format(appStats.TotalTraffic.Http5xxCount))
		fmt.Fprintf(v, "%v", util.BRIGHT_WHITE)
		fmt.Fprintf(v, "  All: %12v\n", util.Format(appStats.TotalTraffic.HttpAllCount))
		fmt.Fprintf(v, "%v", util.CLEAR)
	*/

	// App totals -- this is avaiable on appListView, do we need it here??
	/*
			totalCpuPercentage := 0.0
			totalMemory := uint64(0)
			totalDisk := uint64(0)
			reportingAppInstances := 0
			totalLogCount := int64(0)
			for _, ca := range appStats.ContainerArray {
				if ca != nil {
					if ca.ContainerMetric != nil {
						cpuPercentage := *ca.ContainerMetric.CpuPercentage
						totalCpuPercentage = totalCpuPercentage + cpuPercentage
						memory := *ca.ContainerMetric.MemoryBytes
						totalMemory = totalMemory + memory
						disk := *ca.ContainerMetric.DiskBytes
						totalDisk = totalDisk + disk
					}
					totalLogCount = totalLogCount + (ca.OutCount + ca.ErrCount)
					reportingAppInstances++
				}
			}

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
	*/
	return nil
}
