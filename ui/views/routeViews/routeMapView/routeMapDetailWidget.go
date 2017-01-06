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

package routeMapView

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/app"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/domain"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/org"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/route"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/space"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/jroimartin/gocui"
)

type RouteMapDetailWidget struct {
	masterUI         masterUIInterface.MasterUIInterface
	name             string
	height           int
	routeMapListView *RouteMapListView
	appMdMgr         *app.AppMetadataManager
}

func NewRouteMapDetailWidget(masterUI masterUIInterface.MasterUIInterface, name string, height int, routeMapListView *RouteMapListView) *RouteMapDetailWidget {
	appMdMgr := routeMapListView.GetEventProcessor().GetMetadataManager().GetAppMdManager()
	return &RouteMapDetailWidget{masterUI: masterUI, name: name, height: height, routeMapListView: routeMapListView, appMdMgr: appMdMgr}
}

func (w *RouteMapDetailWidget) Name() string {
	return w.name
}

func (w *RouteMapDetailWidget) Layout(g *gocui.Gui) error {
	maxX, _ := g.Size()
	top := w.routeMapListView.GetTopOffset() - w.height - 1
	width := maxX - 1

	//v, err := g.SetView(w.name, maxX/2-(w.width/2), maxY/2-(w.height/2), maxX/2+(w.width/2), maxY/2+(w.height/2))
	v, err := g.SetView(w.name, 0, top, width, top+w.height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Frame = true
	}
	v.Title = "Route Map Details"
	w.refreshDisplay(g)
	return nil
}

func (w *RouteMapDetailWidget) refreshDisplay(g *gocui.Gui) error {

	v, err := g.View(w.name)
	if err != nil {
		return err
	}
	routeId := w.routeMapListView.routeId
	routeMd := route.FindRouteMetadata(routeId)
	domainMd := domain.FindDomainMetadata(routeMd.DomainGuid)

	var urlBuffer bytes.Buffer
	urlBuffer.WriteString(routeMd.Host)
	if len(routeMd.Host) > 0 {
		urlBuffer.WriteString(".")
	}
	urlBuffer.WriteString(domainMd.Name)
	if routeMd.Port != 0 {
		urlBuffer.WriteString(fmt.Sprintf(":%v", routeMd.Port))
	}
	urlBuffer.WriteString(routeMd.Path)
	url := urlBuffer.String()

	domainType := "Private"
	if domainMd.SharedDomain {
		domainType = "Shared"
	}

	spaceMd := space.FindSpaceMetadata(routeMd.SpaceGuid)
	spaceName := spaceMd.Name
	orgMd := org.FindOrgMetadata(spaceMd.OrgGuid)
	orgName := orgMd.Name

	v.Clear()
	fmt.Fprintf(v, "Route:       %v\n", url)
	fmt.Fprintf(v, "Domain type: %v\n", domainType)
	fmt.Fprintf(v, "Org:         %v\n", orgName)
	fmt.Fprintf(v, "Space:       %v\n", spaceName)

	return nil
}
