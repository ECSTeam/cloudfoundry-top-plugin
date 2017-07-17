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
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventApp"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

func (ed *EventData) logMessageEvent(msg *events.Envelope) {

	logMessage := msg.GetLogMessage()
	appId := logMessage.GetAppId()
	appStats := ed.getAppStats(appId)
	sourceType := logMessage.GetSourceType()
	switch {
	case sourceType == "CELL":
		// The Diego cell emits CELL logs when it starts or stops the app. These actions implement the
		// desired state requested by the user. The Diego cell also emits messages when an app crashes.
		// PCF 1.10 has "CELL" for non-app related logging: e.g. "Container became healthy"
		ed.logCellMsg(msg, logMessage, appStats)
		fallthrough
	case strings.HasPrefix(sourceType, "APP"):
		// PCF 1.6 - 1.9 used "APP" but 1.10 changed to "APP/PROC/WEB/0" and "APP/TASK/f7e79060/0"
		// TODO: Is this wrong? Can a TASK stdout/stderr output be attributed to instance 0 of real app?
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
	case sourceType == "API":
		// This is our notification that the state of an application may have changed
		// e.g., App was marked as STARTED or STOPPED (by a user) or
		// new app deployed or existing app deleted
		ed.logApiCall(msg)

	case sourceType == "RTR":
		// Ignore router log messages
		// Turns out there is nothing useful in this message
		//logMsg := logMessage.GetMessage()
		//ed.handleHttpAccessLogLine(string(logMsg))
	case sourceType == "HEALTH":
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

func (ed *EventData) logCellMsg(msg *events.Envelope, logMessage *events.LogMessage, appStats *eventApp.AppStats) {
	instNum, err := strconv.Atoi(*logMessage.SourceInstance)
	if err != nil {
		return
	}
	containerStats := ed.getContainerStats(appStats, instNum)
	msgTime := time.Unix(0, msg.GetTimestamp())

	msgBytes := logMessage.GetMessage()
	msgText := ""
	if msgBytes != nil {
		msgText = string(logMessage.GetMessage())
	}

	// We need the last "Creating container" time so we can ignore "Destroying container" and "Successfully destroyed container"
	// messages that occur after this time.
	// We need to do tall this because the destroying of a container is async so a new container can be created while the old
	// container is still being destroyed.   We don't want to capture the last cell message of "Successfully destroyed container"
	// When a new container is running hand healthy.
	switch {
	case strings.HasPrefix(msgText, "Exit status"):
		// This message comes out just before the container begins to get destroyed.  The real message is "Exit status 0" but
		// not sure if there could be other exit status values so we just look for the begging of ths string
		containerStats.CellLastExitStatusMsgTime = &msgTime
	case strings.HasPrefix(msgText, "Creating container"):
		containerStats.CellLastCreatingContainerMsgTime = &msgTime
		containerStats.Ip = msg.GetIp()
	case strings.Contains(msgText, "estroying"): // Removed the leading "D" since we don't want to deal with upper/lower case
		fallthrough
	case strings.Contains(msgText, "estroyed"): // Removed the leading "D" since we don't want to deal with upper/lower case
		if containerStats.CellLastExitStatusMsgTime == nil ||
			(containerStats.CellLastCreatingContainerMsgTime != nil &&
				containerStats.CellLastExitStatusMsgTime != nil &&
				containerStats.CellLastCreatingContainerMsgTime.After(*containerStats.CellLastExitStatusMsgTime)) {
			// ignore this event message (see comment above switch statement)
			return
		}
		// Since the container is about to die, we clear out the container metrics since they are no longer valid
		containerStats.ContainerMetric = nil
	}

	if containerStats.CellLastMsgTime == nil || containerStats.CellLastMsgTime.Before(msgTime) {
		containerStats.CellLastMsgText = msgText
		containerStats.CellLastMsgTime = &msgTime
		containerStats.LastUpdateTime = &msgTime
	}
}

func (ed *EventData) logApiCall(msg *events.Envelope) {

	logMessage := msg.GetLogMessage()
	appId := logMessage.GetAppId()
	appStats := ed.getAppStats(appId)

	appMetadata := ed.eventProcessor.GetMetadataManager().GetAppMdManager().FindAppMetadata(appId)
	logText := string(logMessage.GetMessage())
	toplog.Debug("API event occured for app:%v name:%v msg: %v", appId, appMetadata.Name, logText)
	ed.eventProcessor.GetMetadataManager().RequestRefreshAppMetadata(appId)

	if !strings.HasPrefix(logText, "App instance exited") {
		return
	}

	// Message logged:
	// <message:"App instance exited with guid 18507fe2-a67c-4a56-815b-47c9ce195692 payload: {\"instance\"=>\"\", \"index\"=>0, \"reason\"=>\"CRASHED\",
	// 	\"exit_description\"=>\"2 error(s) occurred:\\n\\n* 2 error(s) occurred:\\n\\n* Exited with status 1\\n* cancelled\\n* cancelled\",
	//  \"crash_count\"=>4, \"crash_timestamp\"=>1491877875272064482, \"version\"=>\"ee3ecb12-d1f3-489e-afa4-1ab78fe59381\"}"
	payloadFieldName := "payload:"
	payloadIndex := strings.Index(logText, payloadFieldName)
	if payloadIndex < 0 {
		return
	}

	payload := logText[payloadIndex+len(payloadFieldName) : len(logText)]
	toplog.Debug("payload: %v", payload)
	payload = strings.Replace(payload, "=>", ":", -1)
	jsonParsed, err := gabs.ParseJSON([]byte(payload))
	if err != nil {
		toplog.Error("ParseJSON err: %v payload: %v", err, payload)
		return
	}

	fields, err := jsonParsed.ChildrenMap()
	if err != nil {
		toplog.Error("ParseJSON err: %v payload: %v", err, payload)
		return
	}

	reasonField := fields["reason"]
	if reasonField == nil {
		return
	}
	reason := reasonField.Data().(string)
	toplog.Debug("API app event event occured for app:%v name:%v reason: %v", appId, appMetadata.Name, reason)
	if reason == "CRASHED" {

		index := fields["index"]
		if index == nil {
			return
		}
		instNum := (int)(index.Data().(float64))
		if err != nil {
			return
		}

		exitDescriptionField := fields["exit_description"]
		exitDescription := ""
		if exitDescriptionField != nil {
			exitDescription = exitDescriptionField.Data().(string)
		}

		crashTimestampField := fields["crash_timestamp"]
		if crashTimestampField != nil {
			timestamp64 := crashTimestampField.Data().(float64)
			timestamp := time.Unix(0, int64(timestamp64))
			appStats.AddCrashInfo(instNum, &timestamp, exitDescription)
		}
		toplog.Info("CRASH of app %v exit desc: %v", appMetadata.Name, exitDescription)

	}

}
