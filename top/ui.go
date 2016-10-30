package top


import (
	"fmt"
  "log"
	//"github.com/Sirupsen/logrus"
	//"os"

	//"strings"
	"sync"
	"time"
  "net/url"
  //"encoding/json"
  "github.com/jroimartin/gocui"
  "github.com/cloudfoundry/cli/plugin"
  "github.com/kkellner/cloudfoundry-top-plugin/appStats"
  "github.com/kkellner/cloudfoundry-top-plugin/eventrouting"
  "github.com/kkellner/cloudfoundry-top-plugin/debug"
  "github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
)


type MasterUI struct {
  layoutManager   *LayoutManager
  gui             *gocui.Gui
  cliConnection   plugin.CliConnection

  // TODO: Need to combind detailUI and detailView
  //detailUI        *appStats.AppStatsUI
  detailView      *appStats.DetailView

  mu  sync.Mutex // protects ctr
  router *eventrouting.EventRouter
  refreshNow  chan bool
}


func NewMasterUI(cliConnection plugin.CliConnection ) *MasterUI {
  //detailUI := appStats.NewAppStatsUI(cliConnection)

  ui := &MasterUI {
    //detailUI:      detailUI,
    cliConnection: cliConnection,
    refreshNow:   make(chan bool),
  }

  headerView := NewHeaderWidget(ui, "summaryView", 4)
  footerView := NewFooterWidget("footerView", 4)
  //helpView := NewHelpWidget(ui, "helpView", 60,10)

  detailView := appStats.NewDetailView(ui, "detailView", 5, 4, ui.cliConnection)
  ui.detailView = detailView
  ui.router = eventrouting.NewEventRouter(detailView.GetCurrentProcessor())

  ui.layoutManager = NewLayoutManager()
  ui.layoutManager.Add(headerView)
  ui.layoutManager.Add(footerView)
  //ui.layoutManager.Add(helpView)
  ui.layoutManager.Add(detailView)

  return ui
}

func (ui *MasterUI) CliConnection() plugin.CliConnection {
  return ui.cliConnection
}

func (ui *MasterUI) LayoutManager() masterUIInterface.LayoutManagerInterface {
  return ui.layoutManager
}
func (ui *MasterUI) GetRouter() *eventrouting.EventRouter {
  return ui.router
}

func (ui *MasterUI) Start() {
  ui.detailView.Start()
  ui.initGui()
}


func (ui *MasterUI) initGui() {

  g, err := gocui.NewGui()
	if err != nil {
		log.Panicln(err)
	}
  ui.gui = g

	defer g.Close()
  debug.Init(g)

  //g.SetManagerFunc(ui.layout)
  //filter := appStats.NewFilterWidget("filterWidget", 10, 10, "Example text")


  g.SetManager(ui.layoutManager)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("detailView", 'q', gocui.ModNone, ui.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("detailView", 'Q', gocui.ModNone, ui.quit); err != nil {
		log.Panicln(err)
	}

  /*
	if err := g.SetKeybinding("", 'h', gocui.ModNone,
    func(g *gocui.Gui, v *gocui.View) error {
         if !ui.layoutManager.Contains(helpView) {
           ui.layoutManager.Add(helpView)
         }
         ui.SetCurrentViewOnTop(g,"helpView")
         return nil
    }); err != nil {
		log.Panicln(err)
	}
  */

  if err := g.SetKeybinding("detailView", 'c', gocui.ModNone, ui.clearStats); err != nil {
    log.Panicln(err)
  }
  if err := g.SetKeybinding("detailView", gocui.KeySpace, gocui.ModNone, ui.refeshNow); err != nil {
    log.Panicln(err)
  }

  go ui.counter(g)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func (ui *MasterUI) CloseView(m gocui.Manager, name string ) error {

  	ui.gui.DeleteView(name)
    ui.gui.DeleteKeybindings(name)
    ui.layoutManager.Remove(m)

    if err := ui.SetCurrentViewOnTop(ui.gui, "detailView"); err != nil {
      return err
    }
  	return nil
}


func (ui *MasterUI) SetCurrentViewOnTop(g *gocui.Gui, name string) (error) {
  //log.Panicln(fmt.Sprintf("DEBUG: %v", name))
  if _, err := g.SetCurrentView(name); err != nil {
		return err
	}
  //log.Panicln(fmt.Sprintf("DEBUG2: %v", name))

  if _, err := g.SetViewOnTop(name); err != nil {
    return err
  }
  return nil
}

func (ui *MasterUI) quit(g *gocui.Gui, v *gocui.View) error {
  //TODO: Where should this close go?
	//dopplerConnection.Close()
	return gocui.ErrQuit
}

func (ui *MasterUI) clearStats(g *gocui.Gui, v *gocui.View) error {
  ui.router.Clear()
  ui.detailView.ClearStats(g, v)
	ui.updateDisplay(g)
	return nil
}

func (ui *MasterUI) refeshNow(g *gocui.Gui, v *gocui.View) error {
  ui.refreshNow <- true
  return nil
}


func (ui *MasterUI) counter(g *gocui.Gui) {

  ui.updateDisplay(g)
	for {
		select {
		case <-ui.refreshNow:
      ui.updateDisplay(g)
      //ui.refeshDisplay(g)
		case <-time.After(1000 * time.Millisecond):
      ui.updateDisplay(g)
		}
	}
}

func (ui *MasterUI) updateDisplay(g *gocui.Gui) {
	g.Execute(func(g *gocui.Gui) error {
    ui.updateHeaderDisplay(g)
    ui.detailView.UpdateDisplay(g)
		return nil
	})
}

func (ui *MasterUI) refeshDisplay(g *gocui.Gui) {
	g.Execute(func(g *gocui.Gui) error {
    ui.updateHeaderDisplay(g)
    ui.detailView.RefreshDisplay(g)
		return nil
	})
}

func (ui *MasterUI) updateHeaderDisplay(g *gocui.Gui) error {

  v, err := g.View("summaryView")
  if err != nil {
    return err
  }
  v.Clear()

  fmt.Fprintf(v, "Total events: %-11v", ui.router.GetEventCount())
  fmt.Fprintf(v, "Stats duration: %-10v ", Round(time.Now().Sub(ui.router.GetStartTime()), time.Second) )
  fmt.Fprintf(v, "%v\n", time.Now().Format("2006-01-02 15:04:05"))

  apiEndpoint, err := ui.cliConnection.ApiEndpoint()
  if err != nil {
    return err
  }

  url, err  := url.Parse(apiEndpoint)
  if err != nil {
    return err
  }

  username, err := ui.cliConnection.Username()
  if err != nil {
    return err
  }

  fmt.Fprintf(v, "Target: %v@%v\n", username, url.Host)
  // TODO: this should be info that parent UI has / displays
  //fmt.Fprintf(v, "API EP:%v", apiEndpoint)

  //fmt.Fprintf(v, "Total Apps: %-11v", len(ui.appsMetadata))
  //fmt.Fprintf(v, "Unique Apps: %-11v", len(m))
  return nil
}

func Round(d, r time.Duration) time.Duration {
	if r <= 0 {
		return d
	}
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d = d - m
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}
