package eventdata

type CellStats struct {
	Ip                          string
	NumOfCpus                   int
	CapacityTotalMemory         int64
	CapacityRemainingMemory     int64
	CapacityTotalDisk           int64
	CapacityRemainingDisk       int64
	CapacityTotalContainers     int
	CapacityRemainingContainers int
	ContainerCount              int
}

func NewCellStats(cellIp string) *CellStats {
	stats := &CellStats{}
	stats.Ip = cellIp
	return stats
}
