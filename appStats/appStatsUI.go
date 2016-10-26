package appStats

import (
	"fmt"
  //"log"
	//"github.com/Sirupsen/logrus"
	//"os"
  "sort"
	"sync"
	//"time"
  "github.com/jroimartin/gocui"
  "github.com/cloudfoundry/cli/plugin"
  //cfclient "github.com/cloudfoundry-community/go-cfclient"
  //"github.com/kkellner/cloudfoundry-top-plugin/debug"
  "github.com/kkellner/cloudfoundry-top-plugin/metadata"
)


type AppStatsUI struct {
  processor     *AppStatsEventProcessor
  cliConnection   plugin.CliConnection
  mu  sync.Mutex // protects ctr

  appsMetadata []metadata.App
  spacesMetadata []metadata.Space

}


func NewAppStatsUI(cliConnection plugin.CliConnection ) *AppStatsUI {
  processor := NewAppStatsEventProcessor()
  return &AppStatsUI {
    processor:  processor,
    cliConnection: cliConnection,
  }
}

func (asUI *AppStatsUI) Start() {
  go asUI.loadMetadata()
}

func (asUI *AppStatsUI) GetProcessor() *AppStatsEventProcessor {
    return asUI.processor
}


func (asUI *AppStatsUI) InitGui(g *gocui.Gui) error {
  /*
  if err := g.SetKeybinding("", 'c', gocui.ModNone, asUI.clearStats); err != nil {
    log.Panicln(err)
  }
  */
  return nil
}


func (asUI *AppStatsUI) Layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView("detailView", 0, 5, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//fmt.Fprintln(v, "Hello world!")
		fmt.Fprintln(v, "")
	}

	return nil
}


func (asUI *AppStatsUI) ClearStats(g *gocui.Gui, v *gocui.View) error {
  asUI.processor.Clear()
	return nil
}

func (asUI *AppStatsUI) UpdateDisplay(g *gocui.Gui) error {
	asUI.mu.Lock()
	m := asUI.processor.GetAppMap()
	asUI.mu.Unlock()

  asUI.updateHeader(g, m)

  v, err := g.View("detailView")
  if err != nil {
		return err
	}

  //maxX, maxY := v.Size()
  _, maxY := v.Size()
	if len(m) > 0 {
		v.Clear()
    row := 1
		fmt.Fprintf(v, "%-50v %-10v %-10v %6v %6v %6v %6v %6v\n",
      "APPLICATION","SPACE","ORG", "2XX","3XX","4XX","5XX","TOTAL")

    sortedStatList := asUI.getStats2(m)

    for _, appStats := range sortedStatList {

      row++
      fmt.Fprintf(v, "%-50.50v %-10.10v %-10.10v %6d %6d %6d %6d %6d\n",
          appStats.AppName,
          appStats.SpaceName,
          appStats.OrgName,
          appStats.Event2xxCount,
          appStats.Event3xxCount,
          appStats.Event4xxCount,
          appStats.Event5xxCount,
          appStats.EventCount)
      if row == maxY {
        break
      }
		}
	} else {
		v.Clear()
		fmt.Fprintln(v, "No data yet...")
	}
	return nil

}

func (asUI *AppStatsUI) getStats2(statsMap map[string]*AppStats) []*AppStats {
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

func (asUI *AppStatsUI) updateHeader(g *gocui.Gui, appStatsMap map[string]*AppStats) error {
  v, err := g.View("summaryView")
  if err != nil {
    return err
  }
  fmt.Fprintf(v, "Total Apps: %-11v", len(asUI.appsMetadata))
  fmt.Fprintf(v, "Unique Apps: %-11v", len(appStatsMap))
  return nil
}

func (asUI *AppStatsUI) loadMetadata() {
  metadata.LoadAppCache(asUI.cliConnection)
  metadata.LoadSpaceCache(asUI.cliConnection)
  metadata.LoadOrgCache(asUI.cliConnection)
}
