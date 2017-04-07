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
	"sort"

	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/dataCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

type AlertManager struct {
	masterUI   masterUIInterface.MasterUIInterface
	name       string
	commonData *dataCommon.CommonData

	AlertSize int
	//visableMessages map[string]*AlertMessage
	visableMessages AlertMessages
	expandedMessage map[string]string

	//visableMessages map[string]interface{}

	// TODO: We need an Alert object that contains the messages as well as level: INFO, WARN, ERROR
	// We also need to hold onto an array of Alert objects
	message string
}

func NewAlertManager(masterUI masterUIInterface.MasterUIInterface, commonData *dataCommon.CommonData) *AlertManager {
	return &AlertManager{masterUI: masterUI, commonData: commonData,
		visableMessages: make([]*AlertMessage, 0, 10),
		expandedMessage: make(map[string]string),
	}
}

func (am *AlertManager) CheckForAlerts(g *gocui.Gui) error {
	am.checkForAppsNotInDesiredState(g)
	am.checkForErrorMsgDelta(g)
	return nil
}

func (am *AlertManager) checkForAppsNotInDesiredState(g *gocui.Gui) error {

	commonData := am.commonData

	if am.masterUI.GetDisplayPaused() || !commonData.IsWarmupComplete() {
		return nil
	}

	// TODO: We can't alert if we are not monitoring all the apps
	// Update this alert only on monitored apps if non-privileged
	// for now we just don't alert
	if !am.masterUI.IsPrivileged() {
		return nil
	}

	appsNotInDesiredState := commonData.AppsNotInDesiredState()
	if commonData.IsWarmupComplete() && appsNotInDesiredState > 0 {
		plural := ""
		if appsNotInDesiredState > 1 {
			plural = "s"
		}
		return am.ShowMessage(g, APPS_NOT_IN_DESIRED_STATE, appsNotInDesiredState, plural)

	} else if am.isUserMessageOpen(g) {
		return am.ClearUserMessage(g, APPS_NOT_IN_DESIRED_STATE)
	}
	return nil
}

func (am *AlertManager) checkForErrorMsgDelta(g *gocui.Gui) error {
	_, _, _, errorMsgDelta := toplog.GetMsgDeltas()
	if errorMsgDelta > 0 {
		return am.ShowMessage(g, ErrorsSinceViewed, errorMsgDelta)
	} else {
		return am.ClearUserMessage(g, ErrorsSinceViewed)
	}
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

func (am *AlertManager) ClearUserMessage(g *gocui.Gui, removeMessage *AlertMessage) error {

	_, found := am.expandedMessage[removeMessage.Id]
	if found {
		delete(am.expandedMessage, removeMessage.Id)
		visableMessages := make([]*AlertMessage, 0, 10)
		for _, message := range am.visableMessages {
			if message.Id != removeMessage.Id {
				visableMessages = append(visableMessages, message)
			}
		}
		am.visableMessages = visableMessages
		return am.UpdateMessageView(g)
	}
	return nil
	/*
		alertViewName := "alertView"
		view := am.masterUI.LayoutManager().GetManagerByViewName(alertViewName)
		if view != nil {
			am.masterUI.CloseView(view)
		}
		am.AlertSize = 0
		return nil
	*/
}

// TODO: Have message levels which will colorize differently
func (am *AlertManager) ShowMessage(g *gocui.Gui, message *AlertMessage, args ...interface{}) error {

	//message := MessageCatalog[messageId]
	msgText := message.Text
	expandedMessage := fmt.Sprintf(msgText, args...)
	colorizeText := ""
	switch message.Type {
	case AlertType:
		colorizeText = util.WHITE_TEXT_RED_BG
	case WarnType:
		colorizeText = util.BRIGHT_YELLOW
	case InfoType:
		colorizeText = util.BRIGHT_GREEN

	}
	expandedMessage = fmt.Sprintf(" %v %v: %v %v", colorizeText, message.Type, expandedMessage, util.CLEAR)

	if am.expandedMessage[message.Id] == "" {
		am.visableMessages = append(am.visableMessages, message)
		sort.Sort(am.visableMessages)
	}
	am.expandedMessage[message.Id] = expandedMessage

	return am.UpdateMessageView(g)
}

// TODO: Have message levels which will colorize differently
func (am *AlertManager) UpdateMessageView(g *gocui.Gui) error {
	alertViewName := "alertView"

	/*
		fmt.Fprintf(v, " %v", util.WHITE_TEXT_RED_BG)
		if w.message != "" {
			fmt.Fprintln(v, w.message)
		} else {
			fmt.Fprintln(v, "No ALERT message")
		}
		fmt.Fprintf(v, "%v", util.CLEAR)
	*/

	expandedMsg := ""
	for _, message := range am.visableMessages {
		expandedMsg = expandedMsg + am.expandedMessage[message.Id] + "\n"
	}

	alertHeight := len(am.visableMessages)
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

	alertView.SetMessage(expandedMsg)
	return nil
}
