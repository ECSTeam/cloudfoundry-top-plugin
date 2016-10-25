package top


import (
	"fmt"
  "log"
	//"github.com/Sirupsen/logrus"
	//"os"

	//"strings"
	"sync"
	"time"
  //"encoding/json"
  "github.com/jroimartin/gocui"
  "github.com/cloudfoundry/cli/plugin"
  //cfclient "github.com/cloudfoundry-community/go-cfclient"
  "github.com/kkellner/cloudfoundry-top-plugin/appStats"
  "github.com/kkellner/cloudfoundry-top-plugin/eventrouting"
  "github.com/kkellner/cloudfoundry-top-plugin/debug"
)


type UI struct {
  cliConnection   plugin.CliConnection
  detailUI        *appStats.AppStatsUI
  mu  sync.Mutex // protects ctr
  router *eventrouting.EventRouter
}


func NewUI(cliConnection plugin.CliConnection ) *UI {
  detailUI := appStats.NewAppStatsUI(cliConnection)
  router := eventrouting.NewEventRouter(detailUI.GetProcessor())
  return &UI {
    detailUI:      detailUI,
    cliConnection: cliConnection,
    router: router,
  }
}

func (ui *UI) GetRouter() *eventrouting.EventRouter {
  return ui.router
}

func (ui *UI) Start() {
  ui.detailUI.Start()
  ui.initGui()
  //go c.routeEvent()
}


func (ui *UI) initGui() {

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()
  debug.Init(g)

	g.SetLayout(ui.layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", 'q', gocui.ModNone, ui.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'Q', gocui.ModNone, ui.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'h', gocui.ModNone, ui.showHelp); err != nil {
		log.Panicln(err)
	}
  if err := g.SetKeybinding("", 'c', gocui.ModNone, ui.clearStats); err != nil {
    log.Panicln(err)
  }
	if err := g.SetKeybinding("helpView", gocui.KeyEnter, gocui.ModNone, ui.closeHelp); err != nil {
		log.Panicln(err)
	}

  ui.detailUI.InitGui(g)

  go ui.counter(g)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func (ui *UI) layout(g *gocui.Gui) error {
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

  ui.detailUI.Layout(g)

	return nil
}

func (ui *UI) setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func (ui *UI) quit(g *gocui.Gui, v *gocui.View) error {
  //TODO: Where should this close go?
	//dopplerConnection.Close()
	return gocui.ErrQuit
}

func (ui *UI) showHelp(g *gocui.Gui, v *gocui.View) error {
	 _, err := ui.setCurrentViewOnTop(g, "helpView")
	 return err
}

func (ui *UI) closeHelp(g *gocui.Gui, v *gocui.View) error {
	_, err := ui.setCurrentViewOnTop(g, "detailView")
	return err
}

func (ui *UI) clearStats(g *gocui.Gui, v *gocui.View) error {
  ui.router.Clear()
  ui.detailUI.ClearStats(g, v)
	ui.updateDisplay(g)
	return nil
}


func (ui *UI) counter(g *gocui.Gui) {

  // TODO: What is doneX used for and how is it set?
  doneX := make(chan bool)

	for {
		select {
		case <-doneX:
			return
		case <-time.After(1000 * time.Millisecond):
			ui.updateDisplay(g)
		}
	}
}

func (ui *UI) updateDisplay(g *gocui.Gui) {

	g.Execute(func(g *gocui.Gui) error {
    ui.updateHeaderDisplay(g)
    ui.detailUI.UpdateDisplay(g)
		return nil
	})
}

func (ui *UI) updateHeaderDisplay(g *gocui.Gui) error {

  v, err := g.View("summaryView")
  if err != nil {
    return err
  }
  v.Clear()

  fmt.Fprintf(v, "Total events: %-11v", ui.router.GetEventCount())
  fmt.Fprintf(v, "Stats duration: %-10v ", Round(time.Now().Sub(ui.router.GetStartTime()), time.Second) )
  fmt.Fprintf(v, "%v\n", time.Now().Format("2006-01-02 15:04:05"))
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
