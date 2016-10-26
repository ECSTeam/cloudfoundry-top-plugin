package appStats

import (
	"fmt"
  //"log"
	//"github.com/Sirupsen/logrus"
	//"os"
  "strconv"
	"strings"
  "sort"
	"sync"
	//"time"
  "encoding/json"
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
  orgsMetadata []metadata.Org
}


func NewAppStatsUI(cliConnection plugin.CliConnection ) *AppStatsUI {
  processor := NewAppStatsEventProcessor()
  return &AppStatsUI {
    processor:  processor,
    cliConnection: cliConnection,
  }
}

func (asUI *AppStatsUI) Start() {
  go asUI.reloadMetadata()
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
		fmt.Fprintf(v, "%-50v %-10v %-10v %6v %6v %6v %6v %6v\n", "APPLICATION","SPACE","ORG", "2XX","3XX","4XX","5XX","TOTAL")

    sortedStatList := asUI.getStats2(m)

    for _, appStats := range sortedStatList {

      row++
			appMetadata := asUI.findAppMetadata(appStats.AppId)
      appName := appMetadata.Name
      if appName == "" {
        appName = appStats.AppId
        //appName = appStats.AppUUID.String()
      }

      spaceMetadata := asUI.findSpaceMetadata(appMetadata.SpaceGuid)
      spaceName := spaceMetadata.Name
      if spaceName == "" {
        spaceName = "unknown"
      }

      orgMetadata := asUI.findOrgMetadata(spaceMetadata.OrgGuid)
      orgName := orgMetadata.Name
      if orgName == "" {
        orgName = "unknown"
      }
      fmt.Fprintf(v, "%-50.50v %-10.10v %-10.10v %6d %6d %6d %6d %6d\n",
          appName, spaceName, orgName,
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


func (asUI *AppStatsUI) findAppMetadata(appId string) metadata.App {
	for _, app := range asUI.appsMetadata {
		if app.Guid == appId {
			return app;
		}
	}
	return metadata.App{}
}

func (asUI *AppStatsUI) findSpaceMetadata(spaceGuid string) metadata.Space {
	for _, space := range asUI.spacesMetadata {
		if space.Guid == spaceGuid {
			return space;
		}
	}
	return metadata.Space{}
}

func (asUI *AppStatsUI) findOrgMetadata(orgGuid string) metadata.Org {
	for _, org := range asUI.orgsMetadata {
		if org.Guid == orgGuid {
			return org;
		}
	}
	return metadata.Org{}
}

func (asUI *AppStatsUI) reloadMetadata() {
  asUI.getAppMetadata()
  asUI.getSpaceMetadata()
  asUI.getOrgMetadata()
}

func (asUI *AppStatsUI) getAppMetadata() {

    // Clear cache of any p
    appsMetadata := []metadata.App{ }

		//requestUrl := "/v2/apps?inline-relations-depth=2"
    baseRequestUrl := "/v2/apps"
    totalPages := 1
    for pageCount := 1; pageCount<=totalPages ; pageCount++ {
      requestUrl := baseRequestUrl+"?page="+strconv.FormatInt(int64(pageCount), 10)
      //requestUrl := baseRequestUrl+"?results-per-page=1&page="+strconv.FormatInt(int64(pageCount), 10)
      //fmt.Printf("url: %v  pageCount: %v  totalPages: %v\n", requestUrl, pageCount, totalPages)
      //debug.Debug(fmt.Sprintf("url: %v\n", requestUrl))
  		reponseJSON, err := asUI.cliConnection.CliCommandWithoutTerminalOutput("curl", requestUrl)
  		if err != nil {
  			fmt.Printf("app error: %v\n", err.Error())
  			return
  		}

  		var appResp metadata.AppResponse
  		// joining since it's an array of strings
  		outputStr := strings.Join(reponseJSON, "")
  		outputBytes := []byte(outputStr)
  		err2 := json.Unmarshal(outputBytes, &appResp)
  		if err2 != nil {
  					fmt.Printf("app error: %v\n", err2.Error())
  		}

  		for _, app := range appResp.Resources {
  			app.Entity.Guid = app.Meta.Guid
  			//app.Entity.SpaceData.Entity.Guid = app.Entity.SpaceData.Meta.Guid
  			//app.Entity.SpaceData.Entity.OrgData.Entity.Guid = app.Entity.SpaceData.Entity.OrgData.Meta.Guid
  			appsMetadata = append(appsMetadata, app.Entity)
  		}
      totalPages = appResp.Pages
    }

    asUI.appsMetadata = appsMetadata

    /*
		for _, app := range asUI.appsMetadata {
			fmt.Printf("appName: %v  appGuid:%v spaceGuid:%v\n", app.Name, app.Guid, app.SpaceGuid)
		}
    */
}


func (asUI *AppStatsUI) getSpaceMetadata() {
  // Clear cache of any p
  spacesMetadata := []metadata.Space{ }

  //requestUrl := "/v2/apps?inline-relations-depth=2"
  baseRequestUrl := "/v2/spaces"
  totalPages := 1
  for pageCount := 1; pageCount<=totalPages ; pageCount++ {
    requestUrl := baseRequestUrl+"?page="+strconv.FormatInt(int64(pageCount), 10)
    //requestUrl := baseRequestUrl+"?results-per-page=1&page="+strconv.FormatInt(int64(pageCount), 10)
    //fmt.Printf("url: %v  pageCount: %v  totalPages: %v\n", requestUrl, pageCount, totalPages)
    reponseJSON, err := asUI.cliConnection.CliCommandWithoutTerminalOutput("curl", requestUrl)
    if err != nil {
      fmt.Printf("space error: %v\n", err.Error())
      return
    }

    var spaceResp metadata.SpaceResponse
    outputStr := strings.Join(reponseJSON, "")
    outputBytes := []byte(outputStr)
    err2 := json.Unmarshal(outputBytes, &spaceResp)
    if err2 != nil {
          fmt.Printf("space error: %v\n", err2.Error())
    }

    for _, space := range spaceResp.Resources {
      space.Entity.Guid = space.Meta.Guid
      //space.Entity.OrgGuid = space.Entity.OrgData.Meta.Guid
      spacesMetadata = append(spacesMetadata, space.Entity)
    }
    totalPages = spaceResp.Pages
  }
  asUI.spacesMetadata = spacesMetadata

  /*
  for _, space := range spacesMetadata {
    fmt.Printf("spaceName: %v  spaceGuid: %v  orgGuid: %v\n", space.Name, space.Guid, space.OrgGuid)
  }
  */

}

func (asUI *AppStatsUI) getOrgMetadata() {

  orgsMetadata := []metadata.Org{ }

  baseRequestUrl := "/v2/organizations"
  totalPages := 1
  for pageCount := 1; pageCount<=totalPages ; pageCount++ {
    requestUrl := baseRequestUrl+"?page="+strconv.FormatInt(int64(pageCount), 10)
    //requestUrl := baseRequestUrl+"?results-per-page=1&page="+strconv.FormatInt(int64(pageCount), 10)
    //fmt.Printf("url: %v  pageCount: %v  totalPages: %v\n", requestUrl, pageCount, totalPages)
    reponseJSON, err := asUI.cliConnection.CliCommandWithoutTerminalOutput("curl", requestUrl)
    if err != nil {
      fmt.Printf("org error: %v\n", err.Error())
      return
    }

    var orgResp metadata.OrgResponse
    outputStr := strings.Join(reponseJSON, "")
    outputBytes := []byte(outputStr)
    err2 := json.Unmarshal(outputBytes, &orgResp)
    if err2 != nil {
          fmt.Printf("org error: %v\n", err2.Error())
    }

    for _, org := range orgResp.Resources {
      org.Entity.Guid = org.Meta.Guid
      //space.Entity.OrgGuid = space.Entity.OrgData.Meta.Guid
      orgsMetadata = append(orgsMetadata, org.Entity)
    }
    totalPages = orgResp.Pages
  }
  asUI.orgsMetadata = orgsMetadata

  /*
  for _, org := range orgsMetadata {
    fmt.Printf("orgName: %v  orgGuid: %v\n", org.Name, org.Guid)
  }
  */

  //os.Exit(1)
}
