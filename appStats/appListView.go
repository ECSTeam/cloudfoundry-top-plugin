package appStats

import (
	"fmt"
  "log"
  "sync"
  "bytes"
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

  //highlightAppId string
  //displayIndexOffset int

  currentProcessor         *AppStatsEventProcessor
  displayedProcessor       *AppStatsEventProcessor
  displayedSortedStatList  []*AppStats

  cliConnection   plugin.CliConnection
  mu  sync.Mutex // protects ctr
  filterAppName string

  appDetailView *AppDetailView
  appListWidget *masterUIInterface.ListWidget

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

func (asUI *AppListView) Layout(g *gocui.Gui) error {
  /*
  maxX, maxY := g.Size()
  v, err := g.SetView(asUI.name, 0, asUI.topMargin, maxX-1, maxY-asUI.bottomMargin)
  if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
    v.Title = "App List"
    fmt.Fprintln(v, "")
    */
  if asUI.appListWidget == nil {

    // START list widget
    appListWidget := masterUIInterface.NewListWidget(asUI.masterUI, asUI.name, asUI.topMargin, asUI.bottomMargin, asUI)
    appListWidget.Title = "Testing..."
    appListWidget.GetRowDisplay = asUI.GetRowDisplay
    appListWidget.GetListSize = asUI.GetListSize
    appListWidget.GetRowKey = asUI.GetRowKey
    appListWidget.GetDisplayHeader = asUI.GetDisplayHeader
    asUI.appListWidget = appListWidget
    //asUI.masterUI.LayoutManager().Add(appListWidget)
    // END

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
           helpView := helpView.NewHelpView(asUI.masterUI, "helpView", 60,10, "This is the appStats help text")
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
  //return nil
  /*
  if err := asUI.masterUI.SetCurrentViewOnTop(g, asUI.name); err != nil {
    log.Panicln(err)
  }
  */



}


// This is for debugging -- remove it later
func writeFooter(g *gocui.Gui, msg string) {
  v, _ := g.View("footerView")
  fmt.Fprint(v, msg)

}

func (asUI *AppListView) Start() {
  go asUI.loadCacheAtStartup()
}

func (asUI *AppListView) loadCacheAtStartup() {
  asUI.loadMetadata()
  asUI.seedStatsFromMetadata()
}

func (asUI *AppListView) seedStatsFromMetadata() {

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


func (asUI *AppListView) ClearStats(g *gocui.Gui, v *gocui.View) error {
  // TODO: I think this needs to be in a sync/mutex
  asUI.currentProcessor.Clear()
  asUI.displayedProcessor.Clear()
  asUI.seedStatsFromMetadata()
	return nil
}

func (asUI *AppListView) UpdateDisplay(g *gocui.Gui) error {
	asUI.mu.Lock()
  processorCopy := asUI.currentProcessor.Clone()

  asUI.displayedProcessor = processorCopy
  if len(processorCopy.AppMap) > 0 {
    asUI.displayedSortedStatList = getStats(processorCopy.AppMap)
  }
	asUI.mu.Unlock()
  return asUI.RefreshDisplay(g)
}

func (asUI *AppListView) RefreshDisplay(g *gocui.Gui) error {

  currentView := asUI.masterUI.GetCurrentView(g)
  currentName := currentView.Name()
  if currentName == asUI.name {
    return asUI.refreshListDisplay(g)
  } else if asUI.appDetailView != nil && currentName == asUI.appDetailView.name {
    return asUI.appDetailView.refreshDisplay(g)
  } else {
    return nil
  }
}

func (asUI *AppListView) GetListSize() int {
  return len(asUI.displayedSortedStatList)
}

func (asUI *AppListView) GetRowKey(statIndex int) string  {
  return asUI.displayedSortedStatList[statIndex].AppId
}

func (asUI *AppListView) GetRowDisplay(statIndex int, isSelected bool) string {
  appStats := asUI.displayedSortedStatList[statIndex]
  v := bytes.NewBufferString("")

  if (!isSelected && appStats.TotalTraffic.EventL10Rate > 0) {
    fmt.Fprintf(v, util.BRIGHT_WHITE)
  }

  totalCpuPercentage := 0.0
  reportingAppInstances := 0
  for _, cs := range appStats.ContainerArray {
    if cs != nil && cs.ContainerMetric != nil {
      cpuPercentage := *cs.ContainerMetric.CpuPercentage
      totalCpuPercentage = totalCpuPercentage + cpuPercentage
      reportingAppInstances++
    }
  }

  totalCpuInfo := ""
  if reportingAppInstances==0 {
    totalCpuInfo = fmt.Sprintf("%6v", "--")
  } else {
    totalCpuInfo = fmt.Sprintf("%6.2f", totalCpuPercentage)
  }

  logCount := int64(0)
  for _, cs := range appStats.ContainerArray {
    if cs != nil {
      logCount = logCount + (cs.OutCount + cs.ErrCount)
    }
  }

  avgResponseTimeL60Info := "--"
  if appStats.TotalTraffic.AvgResponseL60Time >= 0 {
    avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL60Time / 1000000
    avgResponseTimeL60Info = fmt.Sprintf("%6.0f", avgResponseTimeMs)
  }

  fmt.Fprintf(v, "%-50.50v ", appStats.AppName)
  fmt.Fprintf(v, "%-10.10v ", appStats.SpaceName)
  fmt.Fprintf(v, "%-10.10v ",appStats.OrgName )
  fmt.Fprintf(v, "%8d ", appStats.TotalTraffic.Http2xxCount)
  fmt.Fprintf(v, "%8d ", appStats.TotalTraffic.Http3xxCount)
  fmt.Fprintf(v, "%8d ", appStats.TotalTraffic.Http4xxCount)
  fmt.Fprintf(v, "%8d ", appStats.TotalTraffic.Http5xxCount)
  fmt.Fprintf(v, "%8d ", appStats.TotalTraffic.HttpAllCount)

  // Last 1 second
  fmt.Fprintf(v, "%5d ", appStats.TotalTraffic.EventL1Rate);
  // Last 10 seconds
  fmt.Fprintf(v, "%5d ", appStats.TotalTraffic.EventL10Rate);
  // Last 60 seconds
  fmt.Fprintf(v, "%5d ", appStats.TotalTraffic.EventL60Rate);

  fmt.Fprintf(v, "%6v ", totalCpuInfo)
  fmt.Fprintf(v, "%3v ", reportingAppInstances)
  fmt.Fprintf(v, "%6v ", avgResponseTimeL60Info)
  fmt.Fprintf(v, "%11v", util.Format(logCount))

  return v.String()
}

func (asUI *AppListView) GetDisplayHeader() string {
  v := bytes.NewBufferString("")
  fmt.Fprintf(v, "%-50v ","APPLICATION")
  fmt.Fprintf(v, "%-10v ","SPACE")
  fmt.Fprintf(v, "%-10v ","ORG")

  fmt.Fprintf(v, "%8v ","2XX")
  fmt.Fprintf(v, "%8v ","3XX")
  fmt.Fprintf(v, "%8v ","4XX")
  fmt.Fprintf(v, "%8v ","5XX")
  fmt.Fprintf(v, "%8v ","TOTAL")

  fmt.Fprintf(v, "%5v ","L1")
  fmt.Fprintf(v, "%5v ","L10")
  fmt.Fprintf(v, "%5v ","L60")

  fmt.Fprintf(v, "%6v ","CPU%")
  fmt.Fprintf(v, "%3v ","RCR")
  fmt.Fprintf(v, "%6v ","RESP")
  fmt.Fprintf(v, "%11v","LOGS")
  return v.String()
}

func (asUI *AppListView) refreshListDisplay(g *gocui.Gui) error {
  err := asUI.appListWidget.RefreshDisplay(g)
  if err != nil {
    return err
  }
  return asUI.updateHeader(g)
}

func (asUI *AppListView) updateHeader(g *gocui.Gui) error {

  v, err := g.View("summaryView")
  if err != nil {
    return err
  }

  totalReportingAppInstances := 0
  totalActiveApps := 0
  for _, appStats := range asUI.displayedSortedStatList {
    for _, cs := range appStats.ContainerArray {
      if cs != nil && cs.ContainerMetric != nil {
        totalReportingAppInstances++
      }
    }
    if appStats.TotalTraffic.EventL60Rate > 0 {
      totalActiveApps++
    }
  }

  fmt.Fprintf(v, "\r")
  fmt.Fprintf(v, "Total Apps: %-11v", metadata.AppMetadataSize())
  // Active apps are apps that have had go-rounter traffic in last 60 seconds
  fmt.Fprintf(v, "Active Apps: %-4v", totalActiveApps)
  // Reporting containers are containers that reported metrics in last 90 seconds
  fmt.Fprintf(v, "Rprt Cntnrs: %-4v", totalReportingAppInstances)

  return nil
}

func (asUI *AppListView) loadMetadata() {
  metadata.LoadAppCache(asUI.cliConnection)
  metadata.LoadSpaceCache(asUI.cliConnection)
  metadata.LoadOrgCache(asUI.cliConnection)
}
