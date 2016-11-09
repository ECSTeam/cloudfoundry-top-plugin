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
  "github.com/kkellner/cloudfoundry-top-plugin/util"
  "github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
  "github.com/kkellner/cloudfoundry-top-plugin/metadata"
)


type MasterUI struct {
  layoutManager   *LayoutManager
  gui             *gocui.Gui
  cliConnection   plugin.CliConnection

  appListView      *appStats.AppListView

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


  appListView := appStats.NewAppListView(ui, "appListView", 5, 4, ui.cliConnection)
  ui.appListView = appListView
  ui.router = eventrouting.NewEventRouter(appListView.GetCurrentProcessor())

  ui.layoutManager = NewLayoutManager()
  ui.layoutManager.Add(headerView)
  ui.layoutManager.Add(footerView)
  ui.layoutManager.Add(appListView)

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
  ui.appListView.Start()
  ui.initGui()
}


func (ui *MasterUI) initGui() {

  g, err := gocui.NewGui()
	if err != nil {
		log.Panicln(err)
	}
  ui.gui = g
  g.InputEsc = true
	defer g.Close()
  debug.Init(g)


  g.SetManager(ui.layoutManager)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("appListView", 'q', gocui.ModNone, ui.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("appListView", 'Q', gocui.ModNone, ui.quit); err != nil {
		log.Panicln(err)
	}
  /*
  if err := g.SetKeybinding("appListView", gocui.KeyEsc, gocui.ModNone, ui.quit); err != nil {
    log.Panicln(err)
  }
  */
  if err := g.SetKeybinding("appListView", 'c', gocui.ModNone, ui.clearStats); err != nil {
    log.Panicln(err)
  }
  if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
       ui.RefeshNow()
       return nil
  }); err != nil {
    log.Panicln(err)
  }

	if err := g.SetKeybinding("appListView", 'i', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			 intervalWidget := NewEditIntervalWidget(ui, "editIntervalWidget", 30, 6)
			 return intervalWidget.Init(g)
	}); err != nil {
		log.Panicln(err)
	}

  go ui.counter(g)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func (ui *MasterUI) CloseView(m masterUIInterface.Manager ) error {

  	ui.gui.DeleteView(m.Name())
    ui.gui.DeleteKeybindings(m.Name())
    nextForFocus := ui.layoutManager.Remove(m)

    if err := ui.SetCurrentViewOnTop(ui.gui, nextForFocus.Name()); err != nil {
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

func (ui *MasterUI) GetCurrentView(g *gocui.Gui) *gocui.View {
  views := g.Views()
  view := views[len(views)-1]
  return view
}


func (ui *MasterUI) quit(g *gocui.Gui, v *gocui.View) error {
  //TODO: Where should this close go?
	//dopplerConnection.Close()
	return gocui.ErrQuit
}

func (ui *MasterUI) clearStats(g *gocui.Gui, v *gocui.View) error {
  ui.router.Clear()
  ui.appListView.ClearStats(g, v)
	ui.updateDisplay(g)
	return nil
}

func (ui *MasterUI) RefeshNow() {
  ui.refreshNow <- true
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
    ui.appListView.UpdateDisplay(g)
		return nil
	})
}

func (ui *MasterUI) refeshDisplay(g *gocui.Gui) {
	g.Execute(func(g *gocui.Gui) error {
    ui.updateHeaderDisplay(g)
    ui.appListView.RefreshDisplay(g)
		return nil
	})
}

func (ui *MasterUI) updateHeaderDisplay(g *gocui.Gui) error {

  v, err := g.View("summaryView")
  if err != nil {
    return err
  }
  v.Clear()

  fmt.Fprintf(v, "Total events: %-13v", util.Format(ui.router.GetEventCount()))

  runtimeSeconds := Round(time.Now().Sub(ui.router.GetStartTime()), time.Second)
  if runtimeSeconds < time.Second * 30 {
    warmUpTimeRemaining := (time.Second * 30) - runtimeSeconds
		fmt.Fprintf(v, util.DIM_GREEN)
    fmt.Fprintf(v, "Warm-up period: %-10v ", warmUpTimeRemaining)
    fmt.Fprintf(v, util.CLEAR)
  } else {
    fmt.Fprintf(v, "Stats duration: %-10v ", runtimeSeconds)
  }

  fmt.Fprintf(v, "%v\n", time.Now().Format("01-02-2006 15:04:05"))

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

  targetDisplay := fmt.Sprintf("%v@%v", username, url.Host)
  fmt.Fprintf(v, "Target: %-45.45v ", targetDisplay)

  displayTotalMem := "--"
  totalMem := metadata.GetTotalMemoryAllStartedApps()
  if totalMem > 0 {
    displayTotalMem = util.ByteSize(totalMem).String()
  }
  // Total quota memory of all running app instances
  fmt.Fprintf(v, "Rsrvd Mem: %8v ", displayTotalMem)

  displayTotalDisk := "--"
  totalDisk := metadata.GetTotalDiskAllStartedApps()
  if totalMem > 0 {
    displayTotalDisk = util.ByteSize(totalDisk).String()
  }

  fmt.Fprintf(v, "Rsrvd Disk: %8v ", displayTotalDisk)
  fmt.Fprintf(v, "\n")
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
