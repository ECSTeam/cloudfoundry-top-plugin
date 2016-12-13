package displaydata

import "github.com/kkellner/cloudfoundry-top-plugin/eventdata"

type DisplayCellStats struct {
	*eventdata.CellStats
	TotalContainerCpuPercentage  float64 // updated when snapshot
	TotalContainerReservedMemory uint64  // updated when snapshot
	TotalContainerUsedMemory     uint64  // updated when snapshot
	TotalContainerReservedDisk   uint64  // updated when snapshot
	TotalContainerUsedDisk       uint64  // updated when snapshot
	TotalReportingContainers     int     // updated when snapshot
	TotalLogOutCount             int64   // updated when snapshot
	TotalLogErrCount             int64   // updated when snapshot

}

func NewDisplayCellStats(cellStats *eventdata.CellStats) *DisplayCellStats {
	stats := &DisplayCellStats{}
	stats.CellStats = cellStats
	return stats
}
