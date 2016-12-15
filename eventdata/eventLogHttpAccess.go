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

package eventdata

import (
	"fmt"
	"regexp"

	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

// Index of flds:  1         2            3         4       5        6        7        8        9         10        11      12                13                14              15              16
// Def of fields:  FROM  -   DATE_TIME    METHOD    PATH    PROTO    CODE     RCV      SEND     REFER     AGENT     SRC_IP  x_forwarded_for   x_forwarded_proto vcap_request_id response_time   app_id
const regexLogStr = `^(\S+) \S+ \[([^\]]+)\] "([A-Z]+) ([^ ]*) ([^"]*)" ([0-9]+) ([0-9]+) ([0-9]+) "([^"]+)" "([^"]+)" ([^ ]*) ([^:]*):"([^"]*)" ([^:]*):"([^ ]*)" ([^:]*):([^ ]*) ([^:]*):([^ ]*) ([^:]*):([^\n]*)`

type EventLogHttpAccess struct {
	regexHttpLog *regexp.Regexp
}

func NewEventLogHttpAccess() *EventLogHttpAccess {

	regexHttpLog := regexp.MustCompile(regexLogStr)
	return &EventLogHttpAccess{
		regexHttpLog: regexHttpLog,
	}
}

func (ha *EventLogHttpAccess) parseHttpAccessLogLine(logLine string) {
	//toplog.Debug(logLine)
	parsedData := ha.regexHttpLog.FindAllStringSubmatch(logLine, -1)
	dataArray := parsedData[0]
	toplog.Debug(fmt.Sprintf("method:%v code:%v", dataArray[3], dataArray[6]))
}
