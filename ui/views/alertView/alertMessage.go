// Copyright (c) 2017 ECS Team, Inc. - All Rights Reserved
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

type AlertLevel string

const (
	HighLevel   AlertLevel = "H"
	MediumLevel            = "M"
	LowLevel               = "L"
)

type MessageType string

const (
	AlertType MessageType = "ALERT"
	WarnType              = "WARN"
	InfoType              = "INFO"
)

var MessageCatalog = make(map[string]*AlertMessage)
var APPS_NOT_IN_DESIRED_STATE = NewAlertMessage("ANIDS", AlertType, "%v application%v not in desired state")
var AutoOpenOnErrorDisabled = NewAlertMessage("AOOE", InfoType, "Auto show errors is off. Errors since viewed: %v")

func init() {
	MessageCatalog[APPS_NOT_IN_DESIRED_STATE.Id] = APPS_NOT_IN_DESIRED_STATE
}

// Example:
type AlertMessage struct {
	Id   string
	Type MessageType
	Text string
}

func NewAlertMessage(
	id string,
	msgType MessageType,
	text string) *AlertMessage {
	alertMsg := &AlertMessage{
		Id:   id,
		Type: msgType,
		Text: text}
	return alertMsg
}
