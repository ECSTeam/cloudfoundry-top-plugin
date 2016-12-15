package displaydata

import "github.com/kkellner/cloudfoundry-top-plugin/eventdata"

type DisplayContainerStats struct {
	*eventdata.ContainerStats
	*eventdata.AppStats
	FreeMemory uint64
	FreeDisk   uint64
	key        string
}

func NewDisplayContainerStats(containerStats *eventdata.ContainerStats, appStats *eventdata.AppStats) *DisplayContainerStats {
	stats := &DisplayContainerStats{}
	stats.ContainerStats = containerStats
	stats.AppStats = appStats
	return stats
}

func (cs *DisplayContainerStats) Id() string {
	if cs.key == "" {
		cs.key = cs.AppId + string(cs.ContainerIndex)
	}
	return cs.key
}
