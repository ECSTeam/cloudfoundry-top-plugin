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

package eventdata

import (
	"strconv"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

func (ed *EventData) logMessageEvent(msg *events.Envelope) {

	logMessage := msg.GetLogMessage()
	appId := logMessage.GetAppId()

	appStats := ed.getAppStats(appId) // Thread here at crash

	switch logMessage.GetSourceType() {
	case "APP":
		instNum, err := strconv.Atoi(*logMessage.SourceInstance)
		if err == nil {
			containerStats := ed.getContainerStats(appStats, instNum)
			switch *logMessage.MessageType {
			case events.LogMessage_OUT:
				containerStats.OutCount++
			case events.LogMessage_ERR:
				containerStats.ErrCount++
			}
		}
	case "API":
		// This is our notification that the state of an application may have changed
		// e.g., App was marked as STARTED or STOPPED (by a user) or
		// new app deployed or existing app deleted
		appMetadata := ed.eventProcessor.GetMetadataManager().GetAppMdManager().FindAppMetadata(appId)
		logText := string(logMessage.GetMessage())
		toplog.Debug("API event occured for app:%v name:%v msg: %v", appId, appMetadata.Name, logText)
		ed.eventProcessor.GetMetadataManager().RequestRefreshAppMetadata(appId)

	case "RTR":
		// Ignore router log messages
		// Turns out there is nothing useful in this message
		//logMsg := logMessage.GetMessage()
		//ed.handleHttpAccessLogLine(string(logMsg))
	case "HEALTH":
		// Ignore health check messages (TODO: Check sourceType of "crashed" messages)
	default:
		// Non-container log -- staging logs, router logs, etc
		switch *logMessage.MessageType {
		case events.LogMessage_OUT:
			appStats.NonContainerStdout++
		case events.LogMessage_ERR:
			appStats.NonContainerStderr++
		}
	}

}
