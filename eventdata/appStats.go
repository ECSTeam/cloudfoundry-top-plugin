package eventdata

import (
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/kkellner/cloudfoundry-top-plugin/metadata"
)

const (
	UnknownName = "unknown"
)

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
	ContainerTrafficMap map[string]*TrafficStats
	TotalTraffic        *TrafficStats

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

func (as *AppStats) Id() string {
	return as.AppId
}

func PopulateNamesIfNeeded(statsMap map[string]*AppStats) []*AppStats {

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
