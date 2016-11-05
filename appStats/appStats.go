package appStats
//package main

import (
    //"fmt"
    //"time"
    //"sort"
    //"strings"
    "github.com/cloudfoundry/sonde-go/events"
    "github.com/kkellner/cloudfoundry-top-plugin/util"
    "github.com/kkellner/cloudfoundry-top-plugin/metadata"
)



type Traffic struct {

  responseL60Time    *util.AvgTracker
  AvgResponseL60Time float64
  EventL60Rate       int

  responseL10Time    *util.AvgTracker
  AvgResponseL10Time float64
  EventL10Rate       int

  responseL1Time    *util.AvgTracker
  AvgResponseL1Time float64
  EventL1Rate       int

  HttpAllCount int64
  Http2xxCount int64
  Http3xxCount int64
  Http4xxCount int64
  Http5xxCount int64

}

type ContainerStats struct {
  ContainerMetric *events.ContainerMetric
  OutCount int64
  ErrCount int64
}

type dataSlice []*AppStats

type AppStats struct {
  AppUUID     *events.UUID
  AppId       string
  AppName     string
  SpaceName   string
  OrgName     string

  NonContainerOutCount int64
  NonContainerErrCount int64

  ContainerArray []*ContainerStats
  ContainerTrafficMap map[string]*Traffic
  TotalTraffic *Traffic
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
func getStats(statsMap map[string]*AppStats) []*AppStats {


  /*
	eventCount := func(c1, c2 util.Sortable) bool {
    d1 := c1.(*AppStats)
    d2 := c2.(*AppStats)
		return d1.HttpAllCount < d2.HttpAllCount
	}
  */
  /*
  eventCountRev := func(c1, c2 util.Sortable) bool {
    d1 := c1.(*AppStats)
    d2 := c2.(*AppStats)
    return d1.HttpAllCount > d2.HttpAllCount
  }
  */

	appName := func(c1, c2 util.Sortable) bool {
    d1 := c1.(*AppStats)
    d2 := c2.(*AppStats)
		return util.CaseInsensitiveLess(d1.AppName, d2.AppName)
	}
  /*
  appNameRev := func(c1, c2 util.Sortable) bool {
    d1 := c1.(*AppStats)
    d2 := c2.(*AppStats)
		return util.CaseInsensitiveLess(d2.AppName, d1.AppName)
	}
  */

  s := make([]util.Sortable, 0, len(statsMap))
  for _, d := range statsMap {
    appMetadata := metadata.FindAppMetadata(d.AppId)
    appName := appMetadata.Name
    if appName == "" {
      appName = d.AppId
      //appName = appStats.AppUUID.String()
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

  //util.OrderedBy(eventCountRev, appName).Sort(s)
  util.OrderedBy(appName).Sort(s)

  s2 := make([]*AppStats, 0, len(s))
  for _, d := range s {
      s2 = append(s2, d.(*AppStats))
  }

  //sort.Sort(sort.Reverse(s))
  return s2
}



/*
// Len is part of sort.Interface.
func (d dataSlice) Len() int {
    return len(d)
}

// Swap is part of sort.Interface.
func (d dataSlice) Swap(i, j int) {
    d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d dataSlice) Less(i, j int) bool {
    if (d[i].EventCount == d[j].EventCount) {
      return d[i].AppName >= d[j].AppName
    }
    return d[i].EventCount <= d[j].EventCount
}
*/





/*

func main() {

    m := map[string]*AppStats {}

    as := NewAppStats("a")
    as.AppName = "A"
    as.SpaceName = "AA"
    as.OrgName = "ORGA"
    as.EventCount = 5
    as.EventL60Rate = 5
    m["a"] = as

    as = NewAppStats("b")
    as.AppName = "B"
    as.SpaceName = "AA"
    as.OrgName = "ORGA"
    as.EventCount = 15
    as.EventL60Rate = 15
    m["b"] = as

    as = NewAppStats("c")
    as.AppName = "C"
    as.SpaceName = "AA"
    as.OrgName = "ORGA"
    as.EventCount = 2
    as.EventL60Rate = 2
    m["c"] = as

    as = NewAppStats("d")
    as.AppName = "D"
    as.SpaceName = "AA"
    as.OrgName = "ORGA"
    as.EventCount = 8
    as.EventL60Rate = 8
    m["d"] = as

    s := getStats(m)

    for _, d := range s {
        //fmt.Printf("%+v\n", *d)
        //d2 := d.(*AppStats)
        fmt.Printf("%v %v %v %v %v\n", d.AppId, d.AppName, d.SpaceName, d.OrgName, d.EventCount)
    }
}
*/
