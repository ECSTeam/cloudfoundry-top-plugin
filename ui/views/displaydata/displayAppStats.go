package displaydata

import "github.com/kkellner/cloudfoundry-top-plugin/eventdata"

type DisplayAppStats struct {
	*eventdata.AppStats

	//TotalTraffic *eventdata.TrafficStats

	TotalCpuPercentage float64 // updated after a clone of this object
	TotalUsedMemory    uint64  // updated after a clone of this object
	TotalUsedDisk      uint64  // updated after a clone of this object

	TotalReportingContainers int   //updated after a clone of this object
	TotalLogStdout           int64 //updated after a clone of this object
	TotalLogStderr           int64 //updated after a clone of this object

}

func NewDisplayAppStats(appStats *eventdata.AppStats) *DisplayAppStats {
	stats := &DisplayAppStats{}
	stats.AppStats = appStats
	return stats
}
