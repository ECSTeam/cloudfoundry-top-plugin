package appStats

import (
	"fmt"
  //"strings"
  "log"
  "github.com/jroimartin/gocui"
  "github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
  "github.com/kkellner/cloudfoundry-top-plugin/util"
)

type AppDetailView struct {
  masterUI masterUIInterface.MasterUIInterface
	name string
  width int
  height int
  appId string
  appListView *AppListView
}

func NewAppDetailView(masterUI masterUIInterface.MasterUIInterface, name string, appId string, appListView *AppListView) *AppDetailView {
	return &AppDetailView{masterUI: masterUI, name: name, appId: appId, appListView: appListView}
}

func (w *AppDetailView) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
  //topMargin := w.appListView.topMargin+2
  //bottomMargin := maxY - (w.appListView.bottomMargin+2)
  //leftMargin := 4
  //rightMargin := maxX-5
  topMargin := w.appListView.topMargin
  bottomMargin := maxY - w.appListView.bottomMargin
  leftMargin := 0
  rightMargin := maxX-1
  v, err := g.SetView(w.name, leftMargin, topMargin, rightMargin, bottomMargin)


	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
    v.Title = "App Details (press ENTER to close)"
    v.Frame = true
    if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeAppDetailView); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, 'q', gocui.ModNone, w.closeAppDetailView); err != nil {
      return err
    }

    if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
      log.Panicln(err)
    }

    w.refreshDisplay(g)

	}
	return nil
}

func (w *AppDetailView) closeAppDetailView(g *gocui.Gui, v *gocui.View) error {
  if err := w.masterUI.CloseView(w, w.name); err != nil {
    return err
  }
	return nil
}

func (w *AppDetailView) refreshDisplay(g *gocui.Gui) error {

  v, err := g.View(w.name)
  if err != nil {
		return err
	}

  v.Clear()

  if (w.appId == "") {
    fmt.Fprintln(v, "No application selected")
    return nil
  }

  m := w.appListView.displayedProcessor.AppMap
  appStats := m[w.appId]


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


  fmt.Fprintf(v, " \n")
  fmt.Fprintf(v, "AppName: %v%v%v\n", util.BRIGHT_WHITE, appStats.AppName, util.CLEAR)
  fmt.Fprintf(v, "AppId:   %v\n", appStats.AppId)
  fmt.Fprintf(v, "AppUUID: %v\n", appStats.AppUUID)
  fmt.Fprintf(v, "Space:   %v\n", appStats.SpaceName)
  fmt.Fprintf(v, "Org:     %v\n", appStats.OrgName)
  fmt.Fprintf(v, "\n")

  fmt.Fprintf(v, "%22v", "")
  fmt.Fprintf(v, "    1sec   10sec   60sec\n")

  fmt.Fprintf(v, "%22v", "HTTP Event Rate:")
  fmt.Fprintf(v, "%8v", appStats.TotalTraffic.EventL1Rate)
  fmt.Fprintf(v, "%8v", appStats.TotalTraffic.EventL10Rate)
  fmt.Fprintf(v, "%8v\n", appStats.TotalTraffic.EventL60Rate)

  fmt.Fprintf(v, "%22v", "Avg Response Time(ms):")
  fmt.Fprintf(v, "%8v", avgResponseTimeL1Info)
  fmt.Fprintf(v, "%8v", avgResponseTimeL10Info)
  fmt.Fprintf(v, "%8v\n", avgResponseTimeL60Info)

  fmt.Fprintf(v, "\nHTTP Requests:\n")
  fmt.Fprintf(v, "2xx: %12v\n", util.Format(appStats.TotalTraffic.Http2xxCount))
  fmt.Fprintf(v, "3xx: %12v\n", util.Format(appStats.TotalTraffic.Http3xxCount))
  fmt.Fprintf(v, "4xx: %12v\n", util.Format(appStats.TotalTraffic.Http4xxCount))
  fmt.Fprintf(v, "5xx: %12v\n", util.Format(appStats.TotalTraffic.Http5xxCount))
  fmt.Fprintf(v, "All: %12v\n", util.Format(appStats.TotalTraffic.HttpAllCount))

  fmt.Fprintf(v, "\nContainer Info:\n")
  fmt.Fprintf(v, "%5v", "Inst")
  fmt.Fprintf(v, "%8v", "CPU%")
  fmt.Fprintf(v, "%12v", "Mem(bytes)")
  fmt.Fprintf(v, "%12v", "Disk(bytes)")
  fmt.Fprintf(v, "%12v", "Stdout")
  fmt.Fprintf(v, "%12v", "Stderr")
  fmt.Fprintf(v, "\n")

  totalCpuPercentage := 0.0
  reportingAppInstances := 0
  totalLogCount := int64(0)
  for i, ca := range appStats.ContainerArray {
    if ca != nil {

      fmt.Fprintf(v, "%5v", i)
      if ca.ContainerMetric != nil {
        cpuPercentage := *ca.ContainerMetric.CpuPercentage
        totalCpuPercentage = totalCpuPercentage + cpuPercentage
        fmt.Fprintf(v, "%8.2f", cpuPercentage)
        fmt.Fprintf(v, "%12v",  util.ByteSize(*ca.ContainerMetric.MemoryBytes))
        fmt.Fprintf(v, "%12v",  util.ByteSize(*ca.ContainerMetric.DiskBytes))
      } else {
        fmt.Fprintf(v, "%8v", "--")
        fmt.Fprintf(v, "%12v", "--")
        fmt.Fprintf(v, "%12v", "--")
      }

      fmt.Fprintf(v, "%12v", util.Format(ca.OutCount))
      fmt.Fprintf(v, "%12v", util.Format(ca.ErrCount))
      totalLogCount = totalLogCount + (ca.OutCount + ca.ErrCount)

      fmt.Fprintf(v, "\n")
      reportingAppInstances++
    }
  }
  fmt.Fprintf(v, "\n")
  if reportingAppInstances==0 {
    fmt.Fprintf(v, "%6v", " Waiting for container metrics...\n\n")
  } else {
    fmt.Fprintf(v, "Total CPU%: %6.2f\n", totalCpuPercentage)
  }
  totalLogCount = totalLogCount + appStats.NonContainerOutCount + appStats.NonContainerErrCount
  fmt.Fprintf(v, "Non container logs - Stdout: %-12v ",util.Format(appStats.NonContainerOutCount))
  fmt.Fprintf(v, "Stderr: %-12v\n",util.Format(appStats.NonContainerErrCount))
  fmt.Fprintf(v, "Total log events: %12v\n", util.Format(totalLogCount))

  //return w.appListView.updateHeader(g)
  return nil
}
