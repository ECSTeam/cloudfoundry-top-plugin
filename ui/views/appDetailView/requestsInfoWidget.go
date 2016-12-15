package appDetailView

import (
	"errors"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

type RequestsInfoWidget struct {
	masterUI   masterUIInterface.MasterUIInterface
	name       string
	height     int
	detailView *AppDetailView
}

func NewRequestsInfoWidget(masterUI masterUIInterface.MasterUIInterface, name string, height int, detailView *AppDetailView) *RequestsInfoWidget {
	return &RequestsInfoWidget{masterUI: masterUI, name: name, height: height, detailView: detailView}
}

func (w *RequestsInfoWidget) Name() string {
	return w.name
}

func (w *RequestsInfoWidget) Layout(g *gocui.Gui) error {
	maxX, _ := g.Size()

	topMargin := 7
	width := maxX - 1

	//v, err := g.SetView(w.name, maxX/2-(w.width/2), maxY/2-(w.height/2), maxX/2+(w.width/2), maxY/2+(w.height/2))
	v, err := g.SetView(w.name, 0, topMargin, width, topMargin+w.height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Frame = true
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeRequestsInfoWidget); err != nil {
			return err
		}

		if err := w.masterUI.SetCurrentViewOnTop(g); err != nil {
			log.Panicln(err)
		}

	}
	v.Title = "App Request Info for: " + w.getAppName()
	w.refreshDisplay(g)
	return nil
}

func (w *RequestsInfoWidget) closeRequestsInfoWidget(g *gocui.Gui, v *gocui.View) error {
	if err := w.masterUI.CloseView(w); err != nil {
		return err
	}
	return nil
}

func (w *RequestsInfoWidget) getAppName() string {
	m := w.detailView.GetDisplayedEventData().AppMap
	appStats := m[w.detailView.appId]
	return appStats.AppName
}

func (w *RequestsInfoWidget) refreshDisplay(g *gocui.Gui) error {

	v, err := g.View("requestsInfoWidget")
	if err != nil {
		return err
	}

	v.Clear()

	if w.detailView.appId == "" {
		fmt.Fprintln(v, "No application selected")
		return nil
	}

	m := w.detailView.GetDisplayedEventData().AppMap
	appStats := m[w.detailView.appId]

	avgResponseTimeL60Info := "--"
	if appStats.TotalTraffic.AvgResponseL60Time >= 0 {
		avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL60Time / 1000000
		avgResponseTimeL60Info = fmt.Sprintf("%8.1f", avgResponseTimeMs)
	}

	avgResponseTimeL10Info := "--"
	if appStats.TotalTraffic.AvgResponseL10Time >= 0 {
		avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL10Time / 1000000
		avgResponseTimeL10Info = fmt.Sprintf("%8.1f", avgResponseTimeMs)
	}

	avgResponseTimeL1Info := "--"
	if appStats.TotalTraffic.AvgResponseL1Time >= 0 {
		avgResponseTimeMs := appStats.TotalTraffic.AvgResponseL1Time / 1000000
		avgResponseTimeL1Info = fmt.Sprintf("%8.1f", avgResponseTimeMs)
	}

	fmt.Fprintf(v, "%22v", "")
	fmt.Fprintf(v, "    1sec   10sec   60sec\n")

	fmt.Fprintf(v, "%22v", "HTTP(S) Event Rate:")
	fmt.Fprintf(v, "%8v", appStats.TotalTraffic.EventL1Rate)
	fmt.Fprintf(v, "%8v", appStats.TotalTraffic.EventL10Rate)
	fmt.Fprintf(v, "%8v\n", appStats.TotalTraffic.EventL60Rate)

	fmt.Fprintf(v, "%22v", "Avg Rspnse Time(ms):")
	fmt.Fprintf(v, "%8v", avgResponseTimeL1Info)
	fmt.Fprintf(v, "%8v", avgResponseTimeL10Info)
	fmt.Fprintf(v, "%8v\n", avgResponseTimeL60Info)
	fmt.Fprintf(v, "%v", util.BRIGHT_WHITE)
	fmt.Fprintf(v, "  Press 'i' for more app info")
	fmt.Fprintf(v, "%v", util.CLEAR)
	return nil
}
