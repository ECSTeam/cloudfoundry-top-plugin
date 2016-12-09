package appStats

import (
	"fmt"
	//"strings"
	"errors"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/metadata"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

const MEGABYTE = (1024 * 1024)

type AppDetailView struct {
	masterUI    masterUIInterface.MasterUIInterface
	name        string
	width       int
	height      int
	appId       string
	appListView *AppListView
}

func NewAppDetailView(masterUI masterUIInterface.MasterUIInterface, name string, appId string, appListView *AppListView) *AppDetailView {
	return &AppDetailView{masterUI: masterUI, name: name, appId: appId, appListView: appListView}
}

func (w *AppDetailView) Name() string {
	return w.name
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
	rightMargin := maxX - 1
	v, err := g.SetView(w.name, leftMargin, topMargin, rightMargin, bottomMargin)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Title = "App Details (press 'q' to quit view)"
		v.Frame = false

		if err := g.SetKeybinding(w.name, 'q', gocui.ModNone, w.closeAppDetailView); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.closeAppDetailView); err != nil {
			return err
		}

		if err := g.SetKeybinding(w.name, 'i', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			appInfoWidget := NewAppInfoWidget(w.masterUI, "appInfoWidget", 70, 20)
			w.masterUI.LayoutManager().Add(appInfoWidget)
			w.masterUI.SetCurrentViewOnTop(g, "appInfoWidget")
			return nil
		}); err != nil {
			log.Panicln(err)
		}

		if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
			log.Panicln(err)
		}

		w.refreshDisplay(g)

	}
	return nil
}

func (w *AppDetailView) closeAppDetailView(g *gocui.Gui, v *gocui.View) error {
	if err := w.masterUI.CloseView(w); err != nil {
		return err
	}
	w.appListView.RefreshDisplay(g)
	return nil
}

func (w *AppDetailView) refreshDisplay(g *gocui.Gui) error {

	v, err := g.View(w.name)
	if err != nil {
		return err
	}

	v.Clear()

	if w.appId == "" {
		fmt.Fprintln(v, "No application selected")
		return nil
	}

	m := w.appListView.GetDisplayedEventData().AppMap
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

	appMetadata := metadata.FindAppMetadata(appStats.AppId)

	memoryDisplay := "--"
	totalMemoryDisplay := "--"
	totalDiskDisplay := "--"
	instancesDisplay := "--"
	diskQuotaDisplay := "--"
	if appMetadata.Guid != "" {
		memoryDisplay = util.ByteSize(appMetadata.MemoryMB * MEGABYTE).String()
		diskQuotaDisplay = util.ByteSize(appMetadata.DiskQuotaMB * MEGABYTE).String()
		instancesDisplay = fmt.Sprintf("%v", appMetadata.Instances)
		totalMemoryDisplay = util.ByteSize((appMetadata.MemoryMB * MEGABYTE) * appMetadata.Instances).String()
		totalDiskDisplay = util.ByteSize((appMetadata.DiskQuotaMB * MEGABYTE) * appMetadata.Instances).String()
	}

	fmt.Fprintf(v, " \n")
	fmt.Fprintf(v, "App Name:        %v%v%v\n", util.BRIGHT_WHITE, appStats.AppName, util.CLEAR)
	fmt.Fprintf(v, "AppId:           %v\n", appStats.AppId)
	fmt.Fprintf(v, "AppUUID:         %v\n", appStats.AppUUID)
	fmt.Fprintf(v, "Space:           %v\n", appStats.SpaceName)
	fmt.Fprintf(v, "Organization:    %v\n", appStats.OrgName)
	fmt.Fprintf(v, "Desired insts:   %v\n", instancesDisplay)
	fmt.Fprintf(v, "Rsrvd mem per (all):  %8v (%8v)\n", memoryDisplay, totalMemoryDisplay)
	fmt.Fprintf(v, "Rsrvd disk per (all): %8v (%8v)\n", diskQuotaDisplay, totalDiskDisplay)
	//fmt.Fprintf(v, "%v", util.BRIGHT_WHITE)
	//fmt.Fprintf(v, "%v", util.CLEAR)
	fmt.Fprintf(v, "\n")

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

	fmt.Fprintf(v, "\n")
	fmt.Fprintf(v, "HTTP(S) status code:\n")
	fmt.Fprintf(v, "  2xx: %12v\n", util.Format(appStats.TotalTraffic.Http2xxCount))
	fmt.Fprintf(v, "  3xx: %12v\n", util.Format(appStats.TotalTraffic.Http3xxCount))
	fmt.Fprintf(v, "  4xx: %12v\n", util.Format(appStats.TotalTraffic.Http4xxCount))
	fmt.Fprintf(v, "  5xx: %12v\n", util.Format(appStats.TotalTraffic.Http5xxCount))
	fmt.Fprintf(v, "%v", util.BRIGHT_WHITE)
	fmt.Fprintf(v, "  All: %12v\n", util.Format(appStats.TotalTraffic.HttpAllCount))
	fmt.Fprintf(v, "%v", util.CLEAR)

	fmt.Fprintf(v, "\nContainer Info:\n")
	fmt.Fprintf(v, "%5v", "INST")
	fmt.Fprintf(v, "%8v", "CPU%")
	fmt.Fprintf(v, "%12v", "MEM USED")
	fmt.Fprintf(v, "%12v", "MEM FREE")
	fmt.Fprintf(v, "%12v", "DISK USED")
	fmt.Fprintf(v, "%12v", "DISK FREE")
	fmt.Fprintf(v, "%12v", "STDOUT")
	fmt.Fprintf(v, "%12v", "STDERR")
	fmt.Fprintf(v, "\n")

	totalCpuPercentage := 0.0
	totalMemory := uint64(0)
	totalDisk := uint64(0)
	reportingAppInstances := 0
	totalLogCount := int64(0)
	for i, ca := range appStats.ContainerArray {
		if ca != nil {

			fmt.Fprintf(v, "%5v", i)
			if ca.ContainerMetric != nil {
				cpuPercentage := *ca.ContainerMetric.CpuPercentage
				totalCpuPercentage = totalCpuPercentage + cpuPercentage
				memory := *ca.ContainerMetric.MemoryBytes
				totalMemory = totalMemory + memory
				disk := *ca.ContainerMetric.DiskBytes
				totalDisk = totalDisk + disk
				fmt.Fprintf(v, "%8.2f", cpuPercentage)
				fmt.Fprintf(v, "%12v", util.ByteSize(memory))
				fmt.Fprintf(v, "%12v", util.ByteSize((uint64(appMetadata.MemoryMB)*MEGABYTE)-memory))
				fmt.Fprintf(v, "%12v", util.ByteSize(disk))
				fmt.Fprintf(v, "%12v", util.ByteSize((uint64(appMetadata.DiskQuotaMB)*MEGABYTE)-disk))
			} else {
				fmt.Fprintf(v, "%8v", "--")
				fmt.Fprintf(v, "%12v", "--")
				fmt.Fprintf(v, "%12v", "--")
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
	//fmt.Fprintf(v, "\n")
	if reportingAppInstances == 0 {
		fmt.Fprintf(v, "%6v", "\n Waiting for container metrics...")
	} else {
		if totalMemory > 0 {
			fmt.Fprintf(v, "%v", util.BRIGHT_WHITE)
			fmt.Fprintf(v, "Total")
			fmt.Fprintf(v, "%8.2f", totalCpuPercentage)
			fmt.Fprintf(v, "%12v", util.ByteSize(totalMemory))
			fmt.Fprintf(v, "%12v", util.ByteSize(totalDisk))
			fmt.Fprintf(v, "%v", util.CLEAR)
		}
	}
	fmt.Fprintf(v, "\n\n")
	totalLogCount = totalLogCount + appStats.NonContainerOutCount + appStats.NonContainerErrCount
	fmt.Fprintf(v, "Non container logs - Stdout: %-12v ", util.Format(appStats.NonContainerOutCount))
	fmt.Fprintf(v, "Stderr: %-12v\n", util.Format(appStats.NonContainerErrCount))
	fmt.Fprintf(v, "Total log events: %12v\n", util.Format(totalLogCount))

	//return w.appListView.updateHeader(g)
	return nil
}
