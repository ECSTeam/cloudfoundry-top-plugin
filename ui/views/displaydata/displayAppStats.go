package displaydata

import "github.com/kkellner/cloudfoundry-top-plugin/eventdata"

type DisplayAppStats struct {
	*eventdata.AppStats

	DesiredContainers int
	//TotalTraffic *eventdata.TrafficStats

	TotalCpuPercentage float64
	TotalUsedMemory    uint64
	TotalUsedDisk      uint64

	TotalReportingContainers int
	TotalLogStdout           int64
	TotalLogStderr           int64
}

func NewDisplayAppStats(appStats *eventdata.AppStats) *DisplayAppStats {
	stats := &DisplayAppStats{}
	stats.AppStats = appStats
	return stats
}
