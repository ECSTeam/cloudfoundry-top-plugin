package eventdata

import "github.com/kkellner/cloudfoundry-top-plugin/util"

type TrafficStats struct {
	responseL60Time    *util.AvgTracker
	AvgResponseL60Time float64 // updated after a clone of this object
	EventL60Rate       int     // updated after a clone of this object

	responseL10Time    *util.AvgTracker
	AvgResponseL10Time float64 // updated after a clone of this object
	EventL10Rate       int     // updated after a clone of this object

	responseL1Time    *util.AvgTracker
	AvgResponseL1Time float64 // updated after a clone of this object
	EventL1Rate       int     // updated after a clone of this object

	HttpAllCount int64
	Http2xxCount int64
	Http3xxCount int64
	Http4xxCount int64
	Http5xxCount int64
}

func NewTrafficStats() *TrafficStats {
	stats := &TrafficStats{}
	return stats
}
