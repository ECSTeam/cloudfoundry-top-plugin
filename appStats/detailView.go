package appStats

import (
	"fmt"
  "log"
  "sync"
  "sort"
  //"strconv"
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

  highlightAppId string
  displayRowOffset int

  currentProcessor         *AppStatsEventProcessor
  displayedProcessor       *AppStatsEventProcessor
  displayedSortedStatList  []*AppStats
  lastProcessor            *AppStatsEventProcessor


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

func (asUI *DetailView) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
  v, err := g.SetView(asUI.name, 0, asUI.topMargin, maxX-1, maxY-asUI.bottomMargin)
  if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
    fmt.Fprintln(v, "")
    filter := NewFilterWidget(asUI.masterUI, "filterWidget", 30, 10)
    if err := g.SetKeybinding(asUI.name, 'f', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
         if !asUI.masterUI.LayoutManager().Contains(filter) {
           asUI.masterUI.LayoutManager().Add(filter)
         }
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

    if err := asUI.masterUI.SetCurrentViewOnTop(g, asUI.name); err != nil {
      log.Panicln(err)
    }

	}
  return nil
}


func (asUI *DetailView) arrowUp(g *gocui.Gui, v *gocui.View) error {

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
          //writeFooter(g,"rowOff["+strconv.Itoa(asUI.displayRowOffset)+"]")

          if (asUI.displayRowOffset > offset) {
            asUI.displayRowOffset = offset
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

func (asUI *DetailView) arrowDown(g *gocui.Gui, v *gocui.View) error {

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
          if (offset>asUI.displayRowOffset) {
            asUI.displayRowOffset = offset
          }
          //writeFooter(g,"\r row["+strconv.Itoa(row)+"]")
          //writeFooter(g,"viewY["+strconv.Itoa(viewY)+"]")
          //writeFooter(g,"o["+strconv.Itoa(offset)+"]")
          //writeFooter(g,"rowOff["+strconv.Itoa(asUI.displayRowOffset)+"]")
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

func (asUI *DetailView) Start() {
  go asUI.loadCacheAtStartup()
}

func (asUI *DetailView) loadCacheAtStartup() {

  asUI.loadMetadata()
  asUI.seedStatsFromMetadata()

}

func (asUI *DetailView) seedStatsFromMetadata() {

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
  // TODO: I think this needs to be in a sync/mutex
  asUI.currentProcessor.Clear()
  asUI.displayedProcessor.Clear()
  asUI.lastProcessor.Clear()
  asUI.seedStatsFromMetadata()
	return nil
}


func (asUI *DetailView) UpdateDisplay(g *gocui.Gui) {
  //g.Execute(func(g *gocui.Gui) error {
  //  return asUI.updateDisplay(g)
	//})
  asUI.updateDisplay(g)
}

func (asUI *DetailView) updateDisplay(g *gocui.Gui) error {
	asUI.mu.Lock()
  //asUI.lastProcessor = asUI.displayedProcessor
  asUI.lastProcessor = deepcopy.Copy(asUI.displayedProcessor).(*AppStatsEventProcessor)
  processorCopy := deepcopy.Copy(asUI.currentProcessor).(*AppStatsEventProcessor)
  asUI.displayedProcessor = processorCopy
  if len(processorCopy.AppMap) > 0 {
    asUI.displayedSortedStatList = asUI.getStats2(processorCopy.AppMap)
  }

  //asUI.displayedProcessor = asUI.currentProcessor
	asUI.mu.Unlock()
  return asUI.refreshDisplay(g)
}

func (asUI *DetailView) RefreshDisplay(g *gocui.Gui) {
	//g.Execute(func(g *gocui.Gui) error {
  //  return asUI.refreshDisplay(g)
	//})
  asUI.refreshDisplay(g)
}

func (asUI *DetailView) refreshDisplay(g *gocui.Gui) error {

  m := asUI.displayedProcessor.AppMap
  lastRefreshAppMap := asUI.lastProcessor.AppMap

  asUI.updateHeader(g)

  v, err := g.View("detailView")
  if err != nil {
		return err
	}

	if len(m) > 0 {
    //maxX, maxY := v.Size()
    //_, maxY := v.Size()
    //maxRows := maxY - 1

		v.Clear()

		fmt.Fprintf(v, "%-50v %-10v %-10v %6v %6v %6v %6v %6v %4v %9v\n",
      "APPLICATION","SPACE","ORG", "2XX","3XX","4XX","5XX","TOTAL", "INTR", "CPU%/I ")

    //sortedStatList := asUI.getStats2(m)


    for row, appStats := range asUI.displayedSortedStatList {

      if row < asUI.displayRowOffset {
        continue
      }

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

      //if row == asUI.highlightRow {
      if appStats.AppId == asUI.highlightAppId {
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
          asUI.displayRowOffset)

      if appStats.AppId == asUI.highlightAppId {
        fmt.Fprintf(v, "\033[0m")
      }

      /*
      if row+1 == maxRows {
        break
      }
      */

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

  fmt.Fprintf(v, "\r")
  fmt.Fprintf(v, "Total Apps: %-11v", metadata.AppMetadataSize())
  // TODO: Active apps are apps that have had go-rounter traffic in last 60 seconds
  fmt.Fprintf(v, "Active Apps: %-11v", "--")
  // TODO: Reporting containers are containers that reported metrics in last 90 seconds
  fmt.Fprintf(v, "Rprt Cntnrs: %-11v", "--")

  return nil
}

func (asUI *DetailView) loadMetadata() {
  metadata.LoadAppCache(asUI.cliConnection)
  metadata.LoadSpaceCache(asUI.cliConnection)
  metadata.LoadOrgCache(asUI.cliConnection)
}
