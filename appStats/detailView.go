package appStats

import (
	"fmt"
  "log"
  "sync"
  "sort"
  "github.com/jroimartin/gocui"
  "github.com/cloudfoundry/cli/plugin"
  "github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
  "github.com/mohae/deepcopy"
  "github.com/kkellner/cloudfoundry-top-plugin/metadata"
)

type DetailView struct {
  masterUI masterUIInterface.MasterUIInterface
	name string
  topMargin int
  bottomMargin int
  highlightRow int

  currentProcessor    *AppStatsEventProcessor
  displayedProcessor  *AppStatsEventProcessor
  lastProcessor       *AppStatsEventProcessor


  cliConnection   plugin.CliConnection
  mu  sync.Mutex // protects ctr
  filterAppName string
  //lastRefreshAppMap      map[string]*AppStats
}

func NewDetailView(masterUI masterUIInterface.MasterUIInterface,name string, topMargin, bottomMargin int,
    cliConnection plugin.CliConnection ) *DetailView {

  currentProcessor := NewAppStatsEventProcessor()
  displayedProcessor := NewAppStatsEventProcessor()
  lastProcessor := NewAppStatsEventProcessor()

	return &DetailView{
    masterUI: masterUI,
    name: name,
    topMargin: topMargin,
    bottomMargin: bottomMargin,
    cliConnection: cliConnection,
    currentProcessor:  currentProcessor,
    displayedProcessor: displayedProcessor,
    lastProcessor: lastProcessor}
}

func (w *DetailView) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
  v, err := g.SetView(w.name, 0, w.topMargin, maxX-1, maxY-w.bottomMargin)
  if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
    fmt.Fprintln(v, "")
    filter := NewFilterWidget(w.masterUI, "filterWidget", 30, 10)
    if err := g.SetKeybinding(w.name, 'f', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
         if !w.masterUI.LayoutManager().Contains(filter) {
           w.masterUI.LayoutManager().Add(filter)
         }
         w.masterUI.SetCurrentViewOnTop(g,"filterWidget")
         return nil
    }); err != nil {
      log.Panicln(err)
    }

    if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
        // TODO: We want to scroll beyond what is visable
         _, viewY := v.Size()
         if w.highlightRow < viewY-2 {
           w.highlightRow++
         }
         return nil
    }); err != nil {
      log.Panicln(err)
    }

    if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
        // TODO: We want to scroll beyond what is visable
         if w.highlightRow > 0 {
           w.highlightRow--
         }
         return nil
    }); err != nil {
      log.Panicln(err)
    }



    if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
      log.Panicln(err)
    }

	}
  return nil
}


func (asUI *DetailView) Start() {
  go asUI.loadCacheAtStartup()
}

func (asUI *DetailView) loadCacheAtStartup() {

  asUI.loadMetadata()

  currentStatsMap := asUI.currentProcessor.AppMap

  for _, app := range metadata.AllApps() {
    appId := app.Guid
    appStats := currentStatsMap[appId]
    if appStats == nil {
      // New app we haven't seen yet
      appStats = &AppStats {
        AppId: appId,
      }
      currentStatsMap[appId] = appStats
    }

  }

}

func (asUI *DetailView) GetCurrentProcessor() *AppStatsEventProcessor {
    return asUI.currentProcessor
}


func (asUI *DetailView) ClearStats(g *gocui.Gui, v *gocui.View) error {
  asUI.currentProcessor.Clear()
  asUI.displayedProcessor.Clear()
  asUI.lastProcessor.Clear()
	return nil
}


func (asUI *DetailView) UpdateDisplay(g *gocui.Gui) error {
	asUI.mu.Lock()
  //asUI.lastProcessor = asUI.displayedProcessor
  asUI.lastProcessor = deepcopy.Copy(asUI.displayedProcessor).(*AppStatsEventProcessor)
  processorCopy := deepcopy.Copy(asUI.currentProcessor).(*AppStatsEventProcessor)
  asUI.displayedProcessor = processorCopy
  //asUI.displayedProcessor = asUI.currentProcessor
	asUI.mu.Unlock()
  return asUI.RefeshDisplay(g)
}


func (asUI *DetailView) RefeshDisplay(g *gocui.Gui) error {

  m := asUI.displayedProcessor.AppMap
  lastRefreshAppMap := asUI.lastProcessor.AppMap

  asUI.updateHeader(g)

  v, err := g.View("detailView")
  if err != nil {
		return err
	}

  //maxX, maxY := v.Size()
  _, maxY := v.Size()
  maxRows := maxY - 1
	if len(m) > 0 {
		v.Clear()

		fmt.Fprintf(v, "%-50v %-10v %-10v %6v %6v %6v %6v %6v %4v %9v\n",
      "APPLICATION","SPACE","ORG", "2XX","3XX","4XX","5XX","TOTAL", "INTR", "CPU%/I ")

    sortedStatList := asUI.getStats2(m)


    for row, appStats := range sortedStatList {

      appId := appStats.AppId
      lastEventCount := uint64(0)
      if lastRefreshAppMap[appId] != nil {
        lastEventCount = lastRefreshAppMap[appId].EventCount
      }
      eventsPerRefresh := appStats.EventCount - lastEventCount



      maxCpuPercentage := -1.0
      maxCpuAppInstance := -1
      for i, cm := range appStats.ContainerMetric {
        if (cm != nil) {
          cpuPercentage := *cm.CpuPercentage
          if (cpuPercentage > maxCpuPercentage) {
            maxCpuPercentage = cpuPercentage
            maxCpuAppInstance = i
          }
        }
      }

      maxCpuInfo := ""
      if maxCpuPercentage==-1 {
        maxCpuInfo = fmt.Sprintf("%9v", "-- ")
      } else {
        maxCpuInfo = fmt.Sprintf("%6.2f/%-2v", maxCpuPercentage, maxCpuAppInstance)
      }

      if row == asUI.highlightRow {
        fmt.Fprintf(v, "\033[32;7m")
      }

      fmt.Fprintf(v, "%-50.50v %-10.10v %-10.10v %6d %6d %6d %6d %6d %4d %9v [%v]\n",
          appStats.AppName,
          appStats.SpaceName,
          appStats.OrgName,
          appStats.Event2xxCount,
          appStats.Event3xxCount,
          appStats.Event4xxCount,
          appStats.Event5xxCount,
          appStats.EventCount, eventsPerRefresh,
          maxCpuInfo,
          asUI.highlightRow)

      if row == asUI.highlightRow {
        fmt.Fprintf(v, "\033[0m")
      }

      if row+1 == maxRows {
        break
      }
		}
	} else {
		v.Clear()
		fmt.Fprintln(v, "No data yet...")
	}
	return nil

}

func (asUI *DetailView) getStats2(statsMap map[string]*AppStats) []*AppStats {
  s := make(dataSlice, 0, len(statsMap))
  for _, d := range statsMap {

      appMetadata := metadata.FindAppMetadata(d.AppId)
      appName := appMetadata.Name
      if appName == "" {
        appName = d.AppId
        //appName = appStats.AppUUID.String()
      }
      d.AppName = appName

      spaceMetadata := metadata.FindSpaceMetadata(appMetadata.SpaceGuid)
      spaceName := spaceMetadata.Name
      if spaceName == "" {
        spaceName = "unknown"
      }
      d.SpaceName = spaceName

      orgMetadata := metadata.FindOrgMetadata(spaceMetadata.OrgGuid)
      orgName := orgMetadata.Name
      if orgName == "" {
        orgName = "unknown"
      }
      d.OrgName = orgName

      s = append(s, d)
  }
  sort.Sort(sort.Reverse(s))
  return s
}

func (asUI *DetailView) updateHeader(g *gocui.Gui) error {


  v, err := g.View("summaryView")
  if err != nil {
    return err
  }

  fmt.Fprintf(v, "Total Apps: %-11v", metadata.AppMetadataSize())
  // TODO: Active apps are apps that have had go-rounter traffic in last 60 seconds
  fmt.Fprintf(v, "Active Apps: %-11v", "--")
  // TODO: Reporting containers are containers that reported metrics in last 90 seconds
  fmt.Fprintf(v, "Rprt Cntnrs: %-11v", "--")


  /*
  for key, _ := range asUI.lastProcessor.AppMap {
    fmt.Fprintf(v, "%v,", key)
  }
  */

  return nil
}

func (asUI *DetailView) loadMetadata() {
  metadata.LoadAppCache(asUI.cliConnection)
  metadata.LoadSpaceCache(asUI.cliConnection)
  metadata.LoadOrgCache(asUI.cliConnection)
}
