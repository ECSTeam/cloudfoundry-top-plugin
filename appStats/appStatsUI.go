package appStats

import (
	"fmt"
  //"log"
	//"github.com/Sirupsen/logrus"
	//"os"

	"strings"
  "sort"
	"sync"
	//"time"
  "encoding/json"
  "github.com/jroimartin/gocui"
  "github.com/cloudfoundry/cli/plugin"
  cfclient "github.com/cloudfoundry-community/go-cfclient"
)


type AppStatsUI struct {
  processor     *AppStatsEventProcessor
  cliConnection   plugin.CliConnection
  mu  sync.Mutex // protects ctr
  appsMetadata []cfclient.App
}


func NewAppStatsUI(cliConnection plugin.CliConnection ) *AppStatsUI {
  processor := NewAppStatsEventProcessor()
  return &AppStatsUI {
    processor:  processor,
    cliConnection: cliConnection,
  }
}

func (asUI *AppStatsUI) Start() {
  go asUI.getAppMetadata()
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
		fmt.Fprintf(v, "%-40v %-10v %-10v %6v %6v %6v %6v %6v\n", "APPLICATION","SPACE","ORG", "2XX","3XX","4XX","5XX","TOTAL")

    sortedStatList := asUI.getStats2(m)

    for _, appStats := range sortedStatList {

      row++
			appMetadata := asUI.findAppMetadata(appStats.AppId)
      appName := appMetadata.Name
      if appName == "" {
        appName = appStats.AppId
      }
      spaceName := appMetadata.SpaceData.Entity.Name
      if spaceName == "" {
        spaceName = "unknown"
      }
      orgName := appMetadata.SpaceData.Entity.OrgData.Entity.Name
      if orgName == "" {
        orgName = "unknown"
      }
      fmt.Fprintf(v, "%-40.40v %-10.10v %-10.10v %6d %6d %6d %6d %6d\n", appName, spaceName, orgName,
          appStats.Event2xxCount,appStats.Event3xxCount,appStats.Event4xxCount,appStats.Event5xxCount,appStats.EventCount)
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

      if d.AppName == "" {
        d.AppName = asUI.findAppMetadata(d.AppId).Name
      }

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


func (asUI *AppStatsUI) findAppMetadata(appId string) cfclient.App {
	for _, app := range asUI.appsMetadata {
		if app.Guid == appId {
			return app;
		}
	}
	return cfclient.App{}
}

func (asUI *AppStatsUI) getAppMetadata() {

		requestUrl := "/v2/apps?inline-relations-depth=2"
		//requestUrl := "/v2/apps"
		reponseJSON, err := asUI.cliConnection.CliCommandWithoutTerminalOutput("curl", requestUrl)
		if err != nil {
			fmt.Printf("error: %v\n", err.Error())
			return
		}

		var appResp cfclient.AppResponse
		// joining since it's an array of strings
		outputStr := strings.Join(reponseJSON, "")
		//fmt.Printf("Response Size: %v\n", len(outputStr))
		outputBytes := []byte(outputStr)
		err2 := json.Unmarshal(outputBytes, &appResp)
		if err2 != nil {
					fmt.Printf("error: %v\n", err.Error())
		}

		//var apps []cfclient.App
		for _, app := range appResp.Resources {
			app.Entity.Guid = app.Meta.Guid
			app.Entity.SpaceData.Entity.Guid = app.Entity.SpaceData.Meta.Guid
			app.Entity.SpaceData.Entity.OrgData.Entity.Guid = app.Entity.SpaceData.Entity.OrgData.Meta.Guid
			asUI.appsMetadata = append(asUI.appsMetadata, app.Entity)
		}

    /*
		for _, app := range asUI.appsMetadata {
			fmt.Printf("appName: %v  %v\n", app.Name, app.Guid)
		}
    */


}
