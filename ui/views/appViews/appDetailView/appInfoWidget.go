// Copyright (c) 2016 ECS Team, Inc. - All Rights Reserved
// https://github.com/ECSTeam/cloudfoundry-top-plugin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package appDetailView

import (
	"errors"
	"fmt"
	"log"

	"github.com/ecsteam/cloudfoundry-top-plugin/metadata"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/dataView"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

type AppInfoWidget struct {
	masterUI   masterUIInterface.MasterUIInterface
	parentView dataView.DataListViewInterface
	name       string
	width      int
	height     int
	detailView *AppDetailView
	mdMgr      *metadata.GlobalManager
}

func NewAppInfoWidget(masterUI masterUIInterface.MasterUIInterface, parentView dataView.DataListViewInterface, name string, width, height int, detailView *AppDetailView) *AppInfoWidget {
	mdMgr := detailView.GetEventProcessor().GetMetadataManager()
	return &AppInfoWidget{masterUI: masterUI, parentView: parentView, name: name, width: width, height: height, detailView: detailView, mdMgr: mdMgr}
}

func (w *AppInfoWidget) Name() string {
	return w.name
}

func (w *AppInfoWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(w.name, maxX/2-(w.width/2), maxY/2-(w.height/2), maxX/2+(w.width/2), maxY/2+(w.height/2))
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Title = "App Information"
		v.Frame = true
		if err := g.SetKeybinding(w.name, 'x', gocui.ModNone, w.closeAppInfoWidget); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.closeAppInfoWidget); err != nil {
			return err
		}
		if err := w.masterUI.SetCurrentViewOnTop(g); err != nil {
			log.Panicln(err)
		}

	}
	w.RefreshDisplay(g)
	return nil
}

func (w *AppInfoWidget) closeAppInfoWidget(g *gocui.Gui, v *gocui.View) error {
	if err := w.masterUI.CloseView(w); err != nil {
		return err
	}
	return nil
}

func (w *AppInfoWidget) UpdateDisplay(g *gocui.Gui) error {

	// Refresh the background (parent) view
	w.parentView.UpdateDisplay(g)

	return w.RefreshDisplay(g)
}

func (w *AppInfoWidget) RefreshDisplay(g *gocui.Gui) error {

	v, err := g.View(w.name)
	if err != nil {
		return err
	}

	// Refresh the background (parent) view
	w.parentView.RefreshDisplay(g)

	v.Clear()

	m := w.detailView.GetDisplayedEventData().AppMap
	appStats := m[w.detailView.appId]
	if appStats == nil {
		return nil
	}
	appMetadata := w.mdMgr.GetAppMdManager().FindItem(appStats.AppId)

	if appMetadata.Guid != "" {
		memoryDisplay := util.ByteSize(appMetadata.MemoryMB * util.MEGABYTE).String()
		diskQuotaDisplay := util.ByteSize(appMetadata.DiskQuotaMB * util.MEGABYTE).String()
		instancesDisplay := fmt.Sprintf("%v", appMetadata.Instances)
		totalMemoryDisplay := util.ByteSize((appMetadata.MemoryMB * util.MEGABYTE) * appMetadata.Instances).String()
		totalDiskDisplay := util.ByteSize((appMetadata.DiskQuotaMB * util.MEGABYTE) * appMetadata.Instances).String()
		state := appMetadata.State
		packageState := appMetadata.PackageState
		buildpack := appMetadata.Buildpack
		if buildpack == "" {
			buildpack = appMetadata.DetectedBuildpack
		}
		packageUpdated := appMetadata.PackageUpdatedAt
		dockerImage := appMetadata.DockerImage

		appName := appMetadata.Name
		spaceMd := w.mdMgr.GetSpaceMdManager().FindItem(appMetadata.SpaceGuid)
		orgMdMgr := w.mdMgr.GetOrgMdManager()
		org := orgMdMgr.FindItem(spaceMd.OrgGuid)
		orgName := org.Name

		spaceName := spaceMd.Name
		isoSegName := w.mdMgr.GetIsoSegMdManager().FindItem(spaceMd.IsolationSegmentGuid).Name

		stackMd := w.mdMgr.GetStackMdManager().FindItem(appMetadata.StackGuid)
		stackName := stackMd.Name

		fmt.Fprintf(v, " \n")
		fmt.Fprintf(v, " App Name:        %v%v%v\n", util.BRIGHT_WHITE, appName, util.CLEAR)
		fmt.Fprintf(v, " AppId:           %v\n", appStats.AppId)
		fmt.Fprintf(v, " AppUUID:         %v\n", appStats.AppUUID)
		fmt.Fprintf(v, " Organization:    %v\n", orgName)
		fmt.Fprintf(v, " Space:           %v\n", spaceName)

		fmt.Fprintf(v, " Stack:           %v\n", stackName)
		fmt.Fprintf(v, " Isolation Seg:   %v\n", isoSegName)

		fmt.Fprintf(v, " Desired insts:   %v\n", instancesDisplay)
		fmt.Fprintf(v, " Package State:   %v\n", packageState)
		fmt.Fprintf(v, " State:           %v\n", state)
		if dockerImage != "" {
			fmt.Fprintf(v, " Docker Image:    %v\n", dockerImage)
		} else {
			fmt.Fprintf(v, " Buildpack:       %v\n", buildpack)
		}
		fmt.Fprintf(v, " Package Updated: %v\n", packageUpdated)
		fmt.Fprintf(v, "\n Reserved:\n")

		fmt.Fprintf(v, "   Mem per (total):  %8v (%8v)\n", memoryDisplay, totalMemoryDisplay)
		fmt.Fprintf(v, "   Disk per (total): %8v (%8v)\n", diskQuotaDisplay, totalDiskDisplay)

	} else {
		fmt.Fprintf(v, " \n Metadata not loaded yet...\n")
	}

	//fmt.Fprintf(v, "%v", util.BRIGHT_WHITE)
	//fmt.Fprintf(v, "%v", util.CLEAR)
	//fmt.Fprintf(v, "\n")

	fmt.Fprintf(v, "\n %vx%v:exit view", "\033[37;1m", "\033[0m")

	return nil
}
