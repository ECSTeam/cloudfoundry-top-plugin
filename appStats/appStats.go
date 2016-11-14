package appStats

//package main

import (
	//"fmt"
	//"time"
	//"sort"
	//"strings"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/kkellner/cloudfoundry-top-plugin/metadata"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

type Traffic struct {
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

type ContainerStats struct {
	ContainerMetric *events.ContainerMetric
	OutCount        int64
	ErrCount        int64
}

type dataSlice []*AppStats

type AppStats struct {
	AppUUID   *events.UUID
	AppId     string
	AppName   string
	SpaceName string
	OrgName   string

	NonContainerOutCount int64
	NonContainerErrCount int64

	ContainerArray      []*ContainerStats
	ContainerTrafficMap map[string]*Traffic
	TotalTraffic        *Traffic

	TotalCpuPercentage float64 // updated after a clone of this object
	TotalUsedMemory    uint64  // updated after a clone of this object
	TotalUsedDisk      uint64  // updated after a clone of this object

	TotalReportingContainers int   //updated after a clone of this object
	TotalLogCount            int64 //updated after a clone of this object
}

func NewAppStats(appId string) *AppStats {
	stats := &AppStats{}
	stats.AppId = appId
	return stats
}

func NewContainerStats() *ContainerStats {
	stats := &ContainerStats{}
	return stats
}

func NewTraffic() *Traffic {
	stats := &Traffic{}
	return stats
}

// Take the stats map and generated a reverse sorted list base on attribute X
func getSortedStats(statsMap map[string]*AppStats, sortFunctions []util.LessFunc) []*AppStats {

	s := make([]util.Sortable, 0, len(statsMap))
	for _, d := range statsMap {
		appMetadata := metadata.FindAppMetadata(d.AppId)
		appName := appMetadata.Name
		if appName == "" {
			appName = d.AppId
		}
		d.AppName = appName

		spaceMetadata := metadata.FindSpaceMetadata(appMetadata.SpaceGuid)
		spaceName := spaceMetadata.Name
		if spaceName == "" {
			spaceName = "unknown"
		}
		d.SpaceName = spaceName

		orgMetadata := metadata.FindOrgMetadata(spaceMetadata.OrgGuid)
		orgName := orgMetadata.Name
		if orgName == "" {
			orgName = "unknown"
		}
		d.OrgName = orgName

		s = append(s, d)
	}

	util.OrderedBy(sortFunctions).Sort(s)

	s2 := make([]*AppStats, len(s))
	for i, d := range s {
		s2[i] = d.(*AppStats)
	}

	return s2
}
