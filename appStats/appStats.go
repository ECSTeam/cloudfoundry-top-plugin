package appStats

//package main

import (

	//"time"
	//"sort"
	//"strings"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/kkellner/cloudfoundry-top-plugin/metadata"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

const (
	UnknownName = "unknown"
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

func populateNamesIfNeeded(statsMap map[string]*AppStats) []*AppStats {

	s := make([]*AppStats, 0, len(statsMap))
	for _, d := range statsMap {
		appMetadata := metadata.FindAppMetadata(d.AppId)
		appName := appMetadata.Name
		if appName == "" {
			appName = d.AppId
		}
		d.AppName = appName

		var spaceMetadata metadata.Space
		spaceName := d.SpaceName
		if spaceName == "" || spaceName == UnknownName {
			spaceMetadata = metadata.FindSpaceMetadata(appMetadata.SpaceGuid)
			spaceName = spaceMetadata.Name
			if spaceName == "" {
				spaceName = UnknownName
			}
			d.SpaceName = spaceName
		}

		orgName := d.OrgName
		if orgName == "" || orgName == UnknownName {
			if &spaceMetadata == nil {
				spaceMetadata = metadata.FindSpaceMetadata(appMetadata.SpaceGuid)
			}
			orgMetadata := metadata.FindOrgMetadata(spaceMetadata.OrgGuid)
			orgName = orgMetadata.Name
			if orgName == "" {
				orgName = UnknownName
			}
			d.OrgName = orgName
		}
		s = append(s, d)
	}
	return s
}

// Take the stats map and generated a reverse sorted list base on attribute X
func getSortedStats(stats []*AppStats, sortFunctions []util.LessFunc) []*AppStats {

	sortStats := make([]util.Sortable, 0, len(stats))
	//debug.Debug(fmt.Sprintf("sortStats size before:%v", len(sortStats)))
	for _, s := range stats {
		sortStats = append(sortStats, s)
	}
	//debug.Debug(fmt.Sprintf("sortStats size after:%v", len(sortStats)))
	util.OrderedBy(sortFunctions).Sort(sortStats)

	s2 := make([]*AppStats, len(sortStats))
	for i, d := range sortStats {
		s2[i] = d.(*AppStats)
	}
	return s2
}
