package eventdata

import (
	"fmt"
	"regexp"

	"github.com/kkellner/cloudfoundry-top-plugin/toplog"
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
