package appStats

import (
	"fmt"
  "log"
	//"github.com/Sirupsen/logrus"
	//"os"
	//"sort"
	"strings"
	"sync"
	"time"
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
  asUI.initGui()
}

func (asUI *AppStatsUI) GetProcessor() *AppStatsEventProcessor {
    return asUI.processor
}


func (asUI *AppStatsUI) initGui() {

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, asUI.quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", 'q', gocui.ModNone, asUI.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'Q', gocui.ModNone, asUI.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'c', gocui.ModNone, asUI.clearStats); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'h', gocui.ModNone, asUI.showHelp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("helpView", gocui.KeyEnter, gocui.ModNone, asUI.closeHelp); err != nil {
		log.Panicln(err)
	}

	go asUI.counter(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("helpView", maxX/2-32, maxY/5, maxX/2+32, maxY/2+5); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Help"
			v.Frame = true
			fmt.Fprintln(v, "Future home of help text")
	}

	if v, err := g.SetView("summaryView", 0, 0, maxX-1, 4); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Summary"
			v.Frame = true
			fmt.Fprintln(v, "")
	}

	if v, err := g.SetView("detailView", 0, 5, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//fmt.Fprintln(v, "Hello world!")
		fmt.Fprintln(v, "")
	}

	if v, err := g.SetView("footerView", 0, maxY-4, maxX-1, maxY); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Footer"
			v.Frame = false
			fmt.Fprintln(v, "c:clear q:quit")
			fmt.Fprintln(v, "s:sleep")
			fmt.Fprintln(v, "h:help")
	}
	return nil
}

func (asUI *AppStatsUI) quit(g *gocui.Gui, v *gocui.View) error {
  //TODO: Where should this close go?
	//dopplerConnection.Close()
	return gocui.ErrQuit
}

func (asUI *AppStatsUI) clearStats(g *gocui.Gui, v *gocui.View) error {
  asUI.processor.Clear()
	asUI.updateDisplay(g)
	return nil
}

func (asUI *AppStatsUI) setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func (asUI *AppStatsUI) showHelp(g *gocui.Gui, v *gocui.View) error {
	 _, err := asUI.setCurrentViewOnTop(g, "helpView")
	 return err
}

func (asUI *AppStatsUI) closeHelp(g *gocui.Gui, v *gocui.View) error {
	_, err := asUI.setCurrentViewOnTop(g, "detailView")
	return err
}


func (asUI *AppStatsUI) counter(g *gocui.Gui) {

  // TODO: What is doneX used for and how is it set?
  doneX := make(chan bool)

	for {
		select {
		case <-doneX:
			return
		case <-time.After(1000 * time.Millisecond):
			asUI.updateDisplay(g)
		}
	}
}

func (asUI *AppStatsUI) updateDisplay(g *gocui.Gui) {
	asUI.mu.Lock()
	m := asUI.processor.GetAppMap()
	asUI.mu.Unlock()

	//maxX, maxY := g.Size()

	g.Execute(func(g *gocui.Gui) error {
		v, err := g.View("detailView")
		if err != nil {
			return err
		}
		if len(m) > 0 {
			v.Clear()
			fmt.Fprintf(v, "%-40v %10v %6v %6v %6v %6v %6v\n", "APPLICATION","SPACE","2XX","3XX","4XX","5XX","TOTAL")
			for appId, count := range m {
				appMetadata := asUI.findAppMetadata(appId)
				fmt.Fprintf(v, "%-40v %10v %6d %6d %6d %6d %6d\n", appMetadata.Name, appMetadata.SpaceData.Entity.Name,0,0,0 ,0,count)
			}
		} else {
			v.Clear()
			fmt.Fprintln(v, "No data yet...")
		}

		v, err = g.View("summaryView")
		if err != nil {
			return err
		}
		v.Clear()
		/*
		for i := 0; i < 200000; i++ {
			v.Clear()
		}
		*/
		fmt.Fprintf(v, "Total events: %-11v", asUI.processor.GetTotalEvents())
		fmt.Fprintf(v, "Total Apps: %-11v", len(asUI.appsMetadata))
		fmt.Fprintf(v, "Unique Apps: %-11v", len(m))
		fmt.Fprintf(v, "%v\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Fprintf(v, "Stats duration: %v\n", 0)
    // TODO: this should be info that parent UI has / displays
		//fmt.Fprintf(v, "API EP:%v", apiEndpoint)

		return nil
	})
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
