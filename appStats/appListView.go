package appStats

import (
	"fmt"
  "log"
  "sync"
  "bytes"
  "strings"
  //"sort"
  //"strconv"
  "github.com/jroimartin/gocui"
  "github.com/cloudfoundry/cli/plugin"
  "github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
  //"github.com/mohae/deepcopy"
  "github.com/kkellner/cloudfoundry-top-plugin/metadata"
  "github.com/kkellner/cloudfoundry-top-plugin/helpView"
  "github.com/kkellner/cloudfoundry-top-plugin/util"
)

type AppListView struct {
  masterUI masterUIInterface.MasterUIInterface
	name string
  topMargin int
  bottomMargin int

  currentProcessor         *AppStatsEventProcessor
  displayedProcessor       *AppStatsEventProcessor
  displayedSortedStatList  []*AppStats

  cliConnection   plugin.CliConnection
  mu  sync.Mutex
  filterAppName string

  appDetailView *AppDetailView
  appListWidget *masterUIInterface.ListWidget

  displayPaused bool

}

func NewAppListView(masterUI masterUIInterface.MasterUIInterface,name string, topMargin, bottomMargin int,
    cliConnection plugin.CliConnection ) *AppListView {

  currentProcessor := NewAppStatsEventProcessor()
  displayedProcessor := NewAppStatsEventProcessor()

	return &AppListView{
    masterUI: masterUI,
    name: name,
    topMargin: topMargin,
    bottomMargin: bottomMargin,
    cliConnection: cliConnection,
    currentProcessor:  currentProcessor,
    displayedProcessor: displayedProcessor,}
}

func (asUI *AppListView) Name() string {
  return asUI.name
}


func (asUI *AppListView) Layout(g *gocui.Gui) error {

  if asUI.appListWidget == nil {

    appListWidget := masterUIInterface.NewListWidget(asUI.masterUI, asUI.name,
      asUI.topMargin, asUI.bottomMargin, asUI, asUI.columnDefinitions())
    appListWidget.Title = "App List"
    appListWidget.PreRowDisplayFunc = asUI.PreRowDisplay
    appListWidget.GetListSize = asUI.GetListSize
    appListWidget.GetRowKey = asUI.GetRowKey

    defaultSortColums := []*masterUIInterface.SortColumn {
      masterUIInterface.NewSortColumn("CPU", true),
      masterUIInterface.NewSortColumn("L60", true),
      masterUIInterface.NewSortColumn("appName", false),
      masterUIInterface.NewSortColumn("spaceName", false),
      masterUIInterface.NewSortColumn("orgName", false),
    }
    appListWidget.SetSortColumns(defaultSortColums)

    asUI.appListWidget = appListWidget

    if err := g.SetKeybinding(asUI.name, 'f', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
         filter := NewFilterWidget(asUI.masterUI, "filterWidget", 30, 10)
         asUI.masterUI.LayoutManager().Add(filter)
         asUI.masterUI.SetCurrentViewOnTop(g,"filterWidget")
         return nil
    }); err != nil {
      log.Panicln(err)
    }

  	if err := g.SetKeybinding(asUI.name, 'h', gocui.ModNone,
      func(g *gocui.Gui, v *gocui.View) error {
           helpView := helpView.NewHelpView(asUI.masterUI, "helpView", 70,15, helpText)
           asUI.masterUI.LayoutManager().Add(helpView)
           asUI.masterUI.SetCurrentViewOnTop(g,"helpView")
           return nil
      }); err != nil {
  		log.Panicln(err)
  	}

    if err := g.SetKeybinding(asUI.name, gocui.KeyEnter, gocui.ModNone,
      func(g *gocui.Gui, v *gocui.View) error {
           asUI.appDetailView = NewAppDetailView(asUI.masterUI, "appDetailView", asUI.appListWidget.HighlightKey(), asUI)
           asUI.masterUI.LayoutManager().Add(asUI.appDetailView)
           asUI.masterUI.SetCurrentViewOnTop(g,"appDetailView")
           //asUI.refreshDisplay(g)
           return nil
      }); err != nil {
  		log.Panicln(err)
  	}

  }

  return asUI.appListWidget.Layout(g)

}

func (asUI *AppListView) columnDefinitions() []*masterUIInterface.ListColumn {
  columns := make([]*masterUIInterface.ListColumn, 0)
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

func (asUI *AppListView) columnAppName() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return util.CaseInsensitiveLess(c1.(*AppStats).AppName, c2.(*AppStats).AppName)
  }

  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    return fmt.Sprintf("%-50.50v", appStats.AppName)
  }
  c := masterUIInterface.NewListColumn("appName", "APPLICATION", 50, true, appNameSortFunc, false, displayFunc)
  return c
}

func (asUI *AppListView) columnSpaceName() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return util.CaseInsensitiveLess(c1.(*AppStats).SpaceName, c2.(*AppStats).SpaceName)
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    return fmt.Sprintf("%-10.10v", appStats.SpaceName)
  }
  c := masterUIInterface.NewListColumn("spaceName", "SPACE", 10, true, appNameSortFunc, false, displayFunc)
  return c
}

func (asUI *AppListView) columnOrgName() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return util.CaseInsensitiveLess(c1.(*AppStats).OrgName, c2.(*AppStats).OrgName)
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    return fmt.Sprintf("%-10.10v", appStats.OrgName)
  }
  c := masterUIInterface.NewListColumn("orgName", "ORG", 10, true, appNameSortFunc, false, displayFunc)
  return c
}

func (asUI *AppListView) columnReportingContainers() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return c1.(*AppStats).TotalReportingContainers < c2.(*AppStats).TotalReportingContainers
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    return fmt.Sprintf("%3v", appStats.TotalReportingContainers)
  }
  c := masterUIInterface.NewListColumn("reportingContainers", "RCR", 3, false, appNameSortFunc, true, displayFunc)
  return c
}

func (asUI *AppListView) columnTotalCpu() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return c1.(*AppStats).TotalCpuPercentage < c2.(*AppStats).TotalCpuPercentage
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    totalCpuInfo := ""
    if appStats.TotalReportingContainers==0 {
      totalCpuInfo = fmt.Sprintf("%6v", "--")
    } else {
      totalCpuInfo = fmt.Sprintf("%6.2f", appStats.TotalCpuPercentage)
    }
    return fmt.Sprintf("%6v", totalCpuInfo)
  }
  c := masterUIInterface.NewListColumn("CPU", "CPU%", 6, false, appNameSortFunc, true, displayFunc)
  return c
}

func (asUI *AppListView) columnTotalMemoryUsed() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return c1.(*AppStats).TotalUsedMemory < c2.(*AppStats).TotalUsedMemory
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    totalMemInfo := ""
    if appStats.TotalReportingContainers==0 {
      totalMemInfo = fmt.Sprintf("%9v", "--")
    } else {
      totalMemInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.TotalUsedMemory))
    }
    return fmt.Sprintf("%9v", totalMemInfo)
  }
  c := masterUIInterface.NewListColumn("MEM", "MEM", 9, false, appNameSortFunc, true, displayFunc)
  return c
}

func (asUI *AppListView) columnTotalDiskUsed() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return c1.(*AppStats).TotalUsedDisk < c2.(*AppStats).TotalUsedDisk
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    totalDiskInfo := ""
    if appStats.TotalReportingContainers==0 {
      totalDiskInfo = fmt.Sprintf("%9v", "--")
    } else {
      totalDiskInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.TotalUsedDisk))
    }
    return fmt.Sprintf("%9v", totalDiskInfo)
  }
  c := masterUIInterface.NewListColumn("DISK", "DISK", 9, false, appNameSortFunc, true, displayFunc)
  return c
}

func (asUI *AppListView) columnAvgResponseTimeL60Info() *masterUIInterface.ListColumn {
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
  c := masterUIInterface.NewListColumn("avgResponseTimeL60", "RESP", 6, false, appNameSortFunc, true, displayFunc)
  return c
}

func (asUI *AppListView) columnLogCount() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return c1.(*AppStats).TotalLogCount < c2.(*AppStats).TotalLogCount
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    return fmt.Sprintf("%11v", util.Format(appStats.TotalLogCount))
  }
  c := masterUIInterface.NewListColumn("totalLogCount", "LOGS", 11, false, appNameSortFunc, true, displayFunc)
  return c
}

func (asUI *AppListView) columnL1() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return c1.(*AppStats).TotalTraffic.EventL1Rate < c2.(*AppStats).TotalTraffic.EventL1Rate
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    return fmt.Sprintf("%5d", appStats.TotalTraffic.EventL1Rate)
  }
  c := masterUIInterface.NewListColumn("L1", "L1", 5, false, appNameSortFunc, true, displayFunc)
  return c
}

func (asUI *AppListView) columnL10() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return c1.(*AppStats).TotalTraffic.EventL10Rate < c2.(*AppStats).TotalTraffic.EventL10Rate
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    return fmt.Sprintf("%5d", appStats.TotalTraffic.EventL10Rate)
  }
  c := masterUIInterface.NewListColumn("L10", "L10", 5, false, appNameSortFunc, true, displayFunc)
  return c
}

func (asUI *AppListView) columnL60() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return c1.(*AppStats).TotalTraffic.EventL60Rate < c2.(*AppStats).TotalTraffic.EventL60Rate
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    return fmt.Sprintf("%5d", appStats.TotalTraffic.EventL60Rate)
  }
  c := masterUIInterface.NewListColumn("L60", "L60", 5, false, appNameSortFunc, true, displayFunc)
  return c
}


func (asUI *AppListView) columnHTTP() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return c1.(*AppStats).TotalTraffic.HttpAllCount < c2.(*AppStats).TotalTraffic.HttpAllCount
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    return fmt.Sprintf("%8d", appStats.TotalTraffic.HttpAllCount)
  }
  c := masterUIInterface.NewListColumn("http", "HTTP", 8, false, appNameSortFunc, true, displayFunc)
  return c
}
func (asUI *AppListView) column2XX() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return c1.(*AppStats).TotalTraffic.Http2xxCount < c2.(*AppStats).TotalTraffic.Http2xxCount
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    return fmt.Sprintf("%8d", appStats.TotalTraffic.Http2xxCount)
  }
  c := masterUIInterface.NewListColumn("2XX", "2XX", 8, false, appNameSortFunc, true, displayFunc)
  return c
}
func (asUI *AppListView) column3XX() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return c1.(*AppStats).TotalTraffic.Http3xxCount < c2.(*AppStats).TotalTraffic.Http3xxCount
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    return fmt.Sprintf("%8d", appStats.TotalTraffic.Http3xxCount)
  }
  c := masterUIInterface.NewListColumn("3XX", "3XX", 8, false, appNameSortFunc, true, displayFunc)
  return c
}

func (asUI *AppListView) column4XX() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return c1.(*AppStats).TotalTraffic.Http4xxCount < c2.(*AppStats).TotalTraffic.Http4xxCount
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    return fmt.Sprintf("%8d", appStats.TotalTraffic.Http4xxCount)
  }
  c := masterUIInterface.NewListColumn("4XX", "4XX", 8, false, appNameSortFunc, true, displayFunc)
  return c
}

func (asUI *AppListView) column5XX() *masterUIInterface.ListColumn {
  appNameSortFunc := func(c1, c2 util.Sortable) bool {
    return c1.(*AppStats).TotalTraffic.Http5xxCount < c2.(*AppStats).TotalTraffic.Http5xxCount
  }
  displayFunc := func(statIndex int, isSelected bool) string {
    appStats := asUI.displayedSortedStatList[statIndex]
    return fmt.Sprintf("%8d", appStats.TotalTraffic.Http5xxCount)
  }
  c := masterUIInterface.NewListColumn("5XX", "5XX", 8, false, appNameSortFunc, true, displayFunc)
  return c
}



func (asUI *AppListView) Start() {
  go asUI.loadCacheAtStartup()
}

func (asUI *AppListView) loadCacheAtStartup() {
  asUI.loadMetadata()
  asUI.seedStatsFromMetadata()
}

func (asUI *AppListView) seedStatsFromMetadata() {

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

func (asUI *AppListView) ClearStats(g *gocui.Gui, v *gocui.View) error {
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
  asUI.SortData()
	asUI.mu.Unlock()
}


func (asUI *AppListView) SortData() {
  if len(asUI.displayedProcessor.AppMap) > 0 {
    sortFunctions := asUI.appListWidget.GetSortFunctions()
    asUI.displayedSortedStatList = getSortedStats(asUI.displayedProcessor.AppMap, sortFunctions)
  } else {
    asUI.displayedSortedStatList = nil
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

func (asUI *AppListView) GetListSize() int {
  return len(asUI.displayedSortedStatList)
}

func (asUI *AppListView) GetRowKey(statIndex int) string  {
  return asUI.displayedSortedStatList[statIndex].AppId
}

func (asUI *AppListView) PreRowDisplay(statIndex int, isSelected bool) string {
  appStats := asUI.displayedSortedStatList[statIndex]
  v := bytes.NewBufferString("")
  if (!isSelected && appStats.TotalTraffic.EventL10Rate > 0) {
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
  for _, appStats := range asUI.displayedSortedStatList {
    for _, cs := range appStats.ContainerArray {
      if cs != nil && cs.ContainerMetric != nil {
        totalReportingAppInstances++
        totalUsedMemoryAppInstances = totalUsedMemoryAppInstances + *cs.ContainerMetric.MemoryBytes
        totalUsedDiskAppInstances = totalUsedDiskAppInstances + *cs.ContainerMetric.DiskBytes
      }
    }
    if appStats.TotalTraffic.EventL60Rate > 0 {
      totalActiveApps++
    }
  }

  fmt.Fprintf(v, "\r")
  fmt.Fprintf(v, "Total Apps: %-5v ", metadata.AppMetadataSize())
  // Active apps are apps that have had go-rounter traffic in last 60 seconds
  fmt.Fprintf(v, "Active Apps: %-4v ", totalActiveApps)
  // Reporting containers are containers that reported metrics in last 90 seconds
  fmt.Fprintf(v, "Rprt Cntnrs: %-4v ", totalReportingAppInstances)

  totalUsedMemoryAppInstancesDisplay := "--"
  totalUsedDiskAppInstancesDisplay := "--"
  if totalReportingAppInstances > 0 {
    totalUsedMemoryAppInstancesDisplay = util.ByteSize(totalUsedMemoryAppInstances).String()
    totalUsedDiskAppInstancesDisplay = util.ByteSize(totalUsedDiskAppInstances).String()
  }
  fmt.Fprintf(v, "Used Mem:  %8v ", totalUsedMemoryAppInstancesDisplay)
  fmt.Fprintf(v, "Used Disk:  %8v ", totalUsedDiskAppInstancesDisplay)


  return nil
}

func (asUI *AppListView) loadMetadata() {
  metadata.LoadAppCache(asUI.cliConnection)
  metadata.LoadSpaceCache(asUI.cliConnection)
  metadata.LoadOrgCache(asUI.cliConnection)
}
