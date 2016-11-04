package appStats

import (
	"fmt"
  "log"
  "sync"
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

  highlightAppId string
  displayIndexOffset int

  currentProcessor         *AppStatsEventProcessor
  displayedProcessor       *AppStatsEventProcessor
  displayedSortedStatList  []*AppStats

  cliConnection   plugin.CliConnection
  mu  sync.Mutex // protects ctr
  filterAppName string

  appDetailView *AppDetailView

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
  maxX, maxY := g.Size()
  v, err := g.SetView(asUI.name, 0, asUI.topMargin, maxX-1, maxY-asUI.bottomMargin)
  if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
    fmt.Fprintln(v, "")
    if err := g.SetKeybinding(asUI.name, 'f', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
         filter := NewFilterWidget(asUI.masterUI, "filterWidget", 30, 10)
         asUI.masterUI.LayoutManager().Add(filter)
         asUI.masterUI.SetCurrentViewOnTop(g,"filterWidget")
         return nil
    }); err != nil {
      log.Panicln(err)
    }

    if err := g.SetKeybinding(asUI.name, gocui.KeyArrowUp, gocui.ModNone, asUI.arrowUp); err != nil {
      log.Panicln(err)
    }
    if err := g.SetKeybinding(asUI.name, gocui.KeyArrowDown, gocui.ModNone, asUI.arrowDown); err != nil {
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
           asUI.appDetailView = NewAppDetailView(asUI.masterUI, "appDetailView", asUI.highlightAppId, asUI)
           asUI.masterUI.LayoutManager().Add(asUI.appDetailView)
           asUI.masterUI.SetCurrentViewOnTop(g,"appDetailView")
           //asUI.refreshDisplay(g)
           return nil
      }); err != nil {
  		log.Panicln(err)
  	}

    if err := asUI.masterUI.SetCurrentViewOnTop(g, asUI.name); err != nil {
      log.Panicln(err)
    }

	}
  return nil
}


func (asUI *AppListView) arrowUp(g *gocui.Gui, v *gocui.View) error {

  statsList := asUI.displayedSortedStatList
  if asUI.highlightAppId == "" {
    if len(statsList) > 0 {
      asUI.highlightAppId = statsList[0].AppId
    }
  } else {
    lastAppId := ""
    for row, appStats := range statsList {
      if appStats.AppId == asUI.highlightAppId {
        if row > 0 {
          asUI.highlightAppId = lastAppId
          offset := row-1
          //writeFooter(g,"\r row["+strconv.Itoa(row)+"]")
          //writeFooter(g,"o["+strconv.Itoa(offset)+"]")
          //writeFooter(g,"rowOff["+strconv.Itoa(asUI.displayIndexOffset)+"]")
          if (asUI.displayIndexOffset > offset) {
            asUI.displayIndexOffset = offset
          }
          break
        }
      }
      lastAppId = appStats.AppId
    }
  }

   asUI.refreshDisplay(g)
   return nil
}

func (asUI *AppListView) arrowDown(g *gocui.Gui, v *gocui.View) error {

  statsList := asUI.displayedSortedStatList
  if asUI.highlightAppId == "" {
    if len(statsList) > 0 {
      asUI.highlightAppId = statsList[0].AppId
    }
  } else {
    for row, appStats := range statsList {
      if appStats.AppId == asUI.highlightAppId {
        if row+1 < len(statsList) {
          asUI.highlightAppId = statsList[row+1].AppId
          _, viewY := v.Size()
          offset := (row+2) - (viewY-1)
          if (offset>asUI.displayIndexOffset) {
            asUI.displayIndexOffset = offset
          }
          //writeFooter(g,"\r row["+strconv.Itoa(row)+"]")
          //writeFooter(g,"viewY["+strconv.Itoa(viewY)+"]")
          //writeFooter(g,"o["+strconv.Itoa(offset)+"]")
          //writeFooter(g,"rowOff["+strconv.Itoa(asUI.displayIndexOffset)+"]")
          break
        }
      }
    }
  }

   asUI.refreshDisplay(g)
   return nil
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


func (asUI *AppListView) UpdateDisplay(g *gocui.Gui) {
  //g.Execute(func(g *gocui.Gui) error {
  //  return asUI.updateDisplay(g)
	//})
  asUI.updateDisplay(g)
}

func (asUI *AppListView) updateDisplay(g *gocui.Gui) error {
	asUI.mu.Lock()
  processorCopy := asUI.currentProcessor.Clone()

  asUI.displayedProcessor = processorCopy
  if len(processorCopy.AppMap) > 0 {
    asUI.displayedSortedStatList = getStats(processorCopy.AppMap)
  }
	asUI.mu.Unlock()
  return asUI.refreshDisplay(g)
}

func (asUI *AppListView) RefreshDisplay(g *gocui.Gui) {
	//g.Execute(func(g *gocui.Gui) error {
  //  return asUI.refreshDisplay(g)
	//})
  asUI.refreshDisplay(g)
}

func (asUI *AppListView) refreshDisplay(g *gocui.Gui) error {

  currentView := asUI.masterUI.GetCurrentView(g)
  if currentView.Name() == asUI.name {
    return asUI.refreshListDisplay(g)
  } else if currentView.Name() == asUI.appDetailView.name {
    return asUI.appDetailView.refreshDisplay(g)
  } else {
    return nil
  }
}

func (asUI *AppListView) refreshListDisplay(g *gocui.Gui) error {

  m := asUI.displayedProcessor.AppMap

  v, err := g.View("appListView")
  if err != nil {
		return err
	}

	if len(m) > 0 {
    //maxX, maxY := v.Size()
    _, maxY := v.Size()
    maxRows := maxY - 1

		v.Clear()

		fmt.Fprintf(v, "%-50v %-10v %-10v %8v %8v %8v %8v %8v %5v %5v %5v %6v %3v %6v %9v\n",
      "APPLICATION","SPACE","ORG", "2XX","3XX","4XX","5XX","TOTAL", "L1", "L10", "L60", "CPU%", "RCR", "RESP", "LOGS")

    totalActiveApps := 0
    totalReportingAppInstances := 0
    row := 0
    for statIndex, appStats := range asUI.displayedSortedStatList {

      if statIndex < asUI.displayIndexOffset {
        continue
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

      logCount := uint64(0)
      for _, cs := range appStats.ContainerArray {
        if cs != nil {
          logCount = logCount + (cs.OutCount + cs.ErrCount)
        }
      }


      if appStats.AppId == asUI.highlightAppId {
        fmt.Fprintf(v, util.GREEN + util.REVERSE)
      }

      avgResponseTimeL60Info := "--"
      if appStats.TotalTraffic.AvgResponseL60Time >= 0 {
        avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL60Time / 1000000
        avgResponseTimeL60Info = fmt.Sprintf("%6.0f", avgResponseTimeMs)
      }

      if appStats.TotalTraffic.EventL60Rate > 0 {
        totalActiveApps++
      }
      totalReportingAppInstances = totalReportingAppInstances + reportingAppInstances


      fmt.Fprintf(v, "%-50.50v %-10.10v %-10.10v %8d %8d %8d %8d %8d %5d %5d %5d %6v %3v %6v %9v\n",
          appStats.AppName,
          appStats.SpaceName,
          appStats.OrgName,
          appStats.TotalTraffic.Http2xxCount,
          appStats.TotalTraffic.Http3xxCount,
          appStats.TotalTraffic.Http4xxCount,
          appStats.TotalTraffic.Http5xxCount,
          appStats.TotalTraffic.HttpAllCount,
          appStats.TotalTraffic.EventL1Rate, // Last 1 second
          appStats.TotalTraffic.EventL10Rate, // Last 10 seconds
          appStats.TotalTraffic.EventL60Rate, // Last 60 seconds
          totalCpuInfo,
          reportingAppInstances,
          avgResponseTimeL60Info,
          logCount)

      if appStats.AppId == asUI.highlightAppId {
        fmt.Fprintf(v, util.CLEAR)
      }

      row++
      if row == maxRows {
        break
      }
		}

    asUI.updateHeader(g, totalActiveApps, totalReportingAppInstances)

	} else {
		v.Clear()
		fmt.Fprintln(v, "No data yet...")
	}
	return nil

}

func (asUI *AppListView) updateHeader(g *gocui.Gui, totalActiveApps int, totalReportingAppInstances int) error {


  v, err := g.View("summaryView")
  if err != nil {
    return err
  }

  fmt.Fprintf(v, "\r")
  fmt.Fprintf(v, "Total Apps: %-11v", metadata.AppMetadataSize())
  // TODO: Active apps are apps that have had go-rounter traffic in last 60 seconds
  fmt.Fprintf(v, "Active Apps: %-4v", totalActiveApps)
  // TODO: Reporting containers are containers that reported metrics in last 90 seconds
  fmt.Fprintf(v, "Rprt Cntnrs: %-4v", totalReportingAppInstances)

  return nil
}

func (asUI *AppListView) loadMetadata() {
  metadata.LoadAppCache(asUI.cliConnection)
  metadata.LoadSpaceCache(asUI.cliConnection)
  metadata.LoadOrgCache(asUI.cliConnection)
}
