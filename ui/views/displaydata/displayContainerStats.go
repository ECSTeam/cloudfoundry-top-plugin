package displaydata

import "github.com/kkellner/cloudfoundry-top-plugin/eventdata"

type DisplayContainerStats struct {
	*eventdata.ContainerStats
	*eventdata.AppStats
}

func NewDisplayContainerStats(containerStats *eventdata.ContainerStats, appStats *eventdata.AppStats) *DisplayContainerStats {
	stats := &DisplayContainerStats{}
	stats.ContainerStats = containerStats
	stats.AppStats = appStats
	return stats
}
