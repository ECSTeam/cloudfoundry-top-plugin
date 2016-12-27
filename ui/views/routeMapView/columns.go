// Copyright (c) 2016 ECS Team, Inc. - All Rights Reserved
// https://github.com/ECSTeam/cloudfoundry-top-plugin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package routeMapView

import (
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func columnAppName() *uiCommon.ListColumn {
	defaultColSize := 50
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*displaydata.DisplayRouteMapStats).AppName, c2.(*displaydata.DisplayRouteMapStats).AppName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayRouteMapStats)
		value := appStats.AppName
		if len(value) == 0 {
			value = "n/a"
		}
		return util.FormatDisplayData(value, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayRouteMapStats)
		return appStats.AppName
	}
	c := uiCommon.NewListColumn("appName", "APPLICATION", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func columnSpaceName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*displaydata.DisplayRouteMapStats).SpaceName, c2.(*displaydata.DisplayRouteMapStats).SpaceName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayRouteMapStats)
		return util.FormatDisplayData(appStats.SpaceName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayRouteMapStats)
		return appStats.SpaceName
	}
	c := uiCommon.NewListColumn("spaceName", "SPACE", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func columnOrgName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*displaydata.DisplayRouteMapStats).OrgName, c2.(*displaydata.DisplayRouteMapStats).OrgName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*displaydata.DisplayRouteMapStats)
		return util.FormatDisplayData(appStats.OrgName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*displaydata.DisplayRouteMapStats)
		return appStats.OrgName
	}
	c := uiCommon.NewListColumn("orgName", "ORG", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func columnTotalRequests() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayRouteMapStats).HttpAllCount < c2.(*displaydata.DisplayRouteMapStats).HttpAllCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		value := util.Format(stats.HttpAllCount)
		return fmt.Sprintf("%10v", value)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		return fmt.Sprintf("%v", stats.HttpAllCount)
	}
	c := uiCommon.NewListColumn("TOTREQ", "TOT_REQ", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func column2xx() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayRouteMapStats).Http2xxCount < c2.(*displaydata.DisplayRouteMapStats).Http2xxCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.Http2xxCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		return fmt.Sprintf("%v", stats.Http2xxCount)
	}
	c := uiCommon.NewListColumn("2XX", "2XX", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func column3xx() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayRouteMapStats).Http3xxCount < c2.(*displaydata.DisplayRouteMapStats).Http3xxCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.Http3xxCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		return fmt.Sprintf("%v", stats.Http3xxCount)
	}
	c := uiCommon.NewListColumn("3XX", "3XX", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func column4xx() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayRouteMapStats).Http4xxCount < c2.(*displaydata.DisplayRouteMapStats).Http4xxCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.Http4xxCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		return fmt.Sprintf("%v", stats.Http4xxCount)
	}
	c := uiCommon.NewListColumn("4XX", "4XX", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func column5xx() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayRouteMapStats).Http5xxCount < c2.(*displaydata.DisplayRouteMapStats).Http5xxCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.Http5xxCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		return fmt.Sprintf("%v", stats.Http5xxCount)
	}
	c := uiCommon.NewListColumn("5XX", "5XX", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func columnResponseContentLength() *uiCommon.ListColumn {
	defaultColSize := 9
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayRouteMapStats).ResponseContentLength < c2.(*displaydata.DisplayRouteMapStats).ResponseContentLength)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%9v", "--")
		} else {
			value := fmt.Sprintf("%9v", util.ByteSize(stats.ResponseContentLength).StringWithPrecision(1))
			return value
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		return fmt.Sprintf("%v", stats.ResponseContentLength)
	}
	c := uiCommon.NewListColumn("RESP_DATA", "RESP_DATA", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)

	return c
}

func columnMethodGet() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayRouteMapStats).HttpMethodGetCount < c2.(*displaydata.DisplayRouteMapStats).HttpMethodGetCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.HttpMethodGetCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		return fmt.Sprintf("%v", stats.HttpMethodGetCount)
	}
	c := uiCommon.NewListColumn("M_GET", "M_GET", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func columnMethodPost() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayRouteMapStats).HttpMethodPostCount < c2.(*displaydata.DisplayRouteMapStats).HttpMethodPostCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.HttpMethodPostCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		return fmt.Sprintf("%v", stats.HttpMethodPostCount)
	}
	c := uiCommon.NewListColumn("M_POST", "M_POST", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func columnMethodPut() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayRouteMapStats).HttpMethodPutCount < c2.(*displaydata.DisplayRouteMapStats).HttpMethodPutCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.HttpMethodPutCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		return fmt.Sprintf("%v", stats.HttpMethodPutCount)
	}
	c := uiCommon.NewListColumn("M_PUT", "M_PUT", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func columnMethodDelete() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayRouteMapStats).HttpMethodDeleteCount < c2.(*displaydata.DisplayRouteMapStats).HttpMethodDeleteCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.HttpMethodDeleteCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		return fmt.Sprintf("%v", stats.HttpMethodDeleteCount)
	}
	c := uiCommon.NewListColumn("M_DELETE", "M_DELETE", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func columnLastAccess() *uiCommon.ListColumn {
	defaultColSize := 19
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayRouteMapStats).LastAccess.Before(c2.(*displaydata.DisplayRouteMapStats).LastAccess)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%19v", "--")
		} else {
			return fmt.Sprintf("%19v", stats.LastAccess.Format("01-02-2006 15:04:05"))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*displaydata.DisplayRouteMapStats)
		return fmt.Sprintf("%v", stats.LastAccess)
	}
	c := uiCommon.NewListColumn("LAST_ACCESS", "LAST_ACCESS", defaultColSize,
		uiCommon.TIMESTAMP, true, sortFunc, true, displayFunc, rawValueFunc)
	return c
}
