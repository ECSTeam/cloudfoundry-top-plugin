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

package alertView

import (
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/dataCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/jroimartin/gocui"
)

type AlertManager struct {
	masterUI   masterUIInterface.MasterUIInterface
	name       string
	commonData *dataCommon.CommonData

	AlertSize int

	// TODO: We need an Alert object that contains the messages as well as level: INFO, WARN, ERROR
	// We also need to hold onto an array of Alert objects
	message string
}

func NewAlertManager(masterUI masterUIInterface.MasterUIInterface, commonData *dataCommon.CommonData) *AlertManager {
	return &AlertManager{masterUI: masterUI, commonData: commonData}
}

func (am *AlertManager) CheckForAlerts(g *gocui.Gui) error {

	// TODO: We can't alert if we are not monitoring all the apps
	// Update this alert only on monitored apps if non-privileged
	// for now we just don't alert
	if !am.masterUI.IsPrivileged() {
		return nil
	}

	commonData := am.commonData
	appsNotInDesiredState := commonData.AppsNotInDesiredState()

	if commonData.IsWarmupComplete() && appsNotInDesiredState > 0 {
		plural := ""
		if appsNotInDesiredState > 1 {
			plural = "s"
		}
		msg := fmt.Sprintf("ALERT: %v application%v not in desired state (row%v colored red) ",
			appsNotInDesiredState, plural, plural)
		am.ShowUserMessage(g, msg)
	} else if am.isUserMessageOpen(g) {
		am.ClearUserMessage(g)
	}
	return nil
}

func (am *AlertManager) isUserMessageOpen(g *gocui.Gui) bool {
	alertViewName := "alertView"
	view := am.masterUI.LayoutManager().GetManagerByViewName(alertViewName)
	if view != nil {
		return true
	} else {
		return false
	}
}

func (am *AlertManager) ClearUserMessage(g *gocui.Gui) error {
	alertViewName := "alertView"
	view := am.masterUI.LayoutManager().GetManagerByViewName(alertViewName)
	if view != nil {
		am.masterUI.CloseView(view)
	}
	am.AlertSize = 0
	return nil
}

// TODO: Have message levels which will colorize differently
func (am *AlertManager) ShowUserMessage(g *gocui.Gui, message string) error {
	alertViewName := "alertView"
	alertHeight := 1

	am.AlertSize = alertHeight

	var alertView *AlertWidget
	view := am.masterUI.LayoutManager().GetManagerByViewName(alertViewName)
	if view == nil {
		alertView = NewAlertWidget(am.masterUI, alertViewName, alertHeight, am.commonData)
		am.masterUI.LayoutManager().AddToBack(alertView)
		am.masterUI.SetCurrentViewOnTop(g)
	} else {
		// This check is to prevent alert from showing on top of the log window
		if am.masterUI.GetCurrentView(g).Name() == alertViewName {
			if _, err := g.SetViewOnTop(alertViewName); err != nil {
				return err
			}
		}
		alertView = view.(*AlertWidget)
		alertView.SetHeight(alertHeight)
	}
	alertView.SetMessage(message)
	return nil
}
