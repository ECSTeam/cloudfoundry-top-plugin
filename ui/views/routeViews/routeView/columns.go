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

package routeView

import (
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func columnRouteId() *uiCommon.ListColumn {
	defaultColSize := 16
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).RouteId) < (c2.(*DisplayRouteStats).RouteId)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		return util.FormatDisplayData(stats.RouteId, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return stats.RouteId
	}
	c := uiCommon.NewListColumn("ROUTE_ID", "ROUTE_ID", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnRoutedAppCount() *uiCommon.ListColumn {
	defaultColSize := 6
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).RoutedAppCount) < (c2.(*DisplayRouteStats).RoutedAppCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		value := ""
		if stats.Port != 0 {
			value = fmt.Sprintf("%6v", "n/a")
		} else if stats.HttpAllCount == 0 {
			value = fmt.Sprintf("%6v", "--")
		} else {
			value = fmt.Sprintf("%6v", stats.RoutedAppCount)
		}
		return value

	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return fmt.Sprintf("%v", stats.RoutedAppCount)
	}
	c := uiCommon.NewListColumn("R_APPS", "R_APPS", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnRouteName() *uiCommon.ListColumn {
	defaultColSize := 50
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).RouteName) < (c2.(*DisplayRouteStats).RouteName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		return util.FormatDisplayData(stats.RouteName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return stats.RouteName
	}
	c := uiCommon.NewListColumn("ROUTE_NAME", "ROUTE_NAME", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnHost() *uiCommon.ListColumn {
	defaultColSize := 30
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).Host) < (c2.(*DisplayRouteStats).Host)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		hostName := stats.Host
		if hostName == "" {
			hostName = "--"
		}
		return util.FormatDisplayData(hostName, defaultColSize)

	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return stats.Host
	}
	c := uiCommon.NewListColumn("HOST", "HOST", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnDomain() *uiCommon.ListColumn {
	defaultColSize := 25
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).Domain) < (c2.(*DisplayRouteStats).Domain)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		return util.FormatDisplayData(stats.Domain, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return stats.Domain
	}
	c := uiCommon.NewListColumn("DOMAIN", "DOMAIN", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnPath() *uiCommon.ListColumn {
	defaultColSize := 25
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).Path) < (c2.(*DisplayRouteStats).Path)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		return util.FormatDisplayData(stats.Path, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return stats.Path
	}
	c := uiCommon.NewListColumn("PATH", "PATH", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnPort() *uiCommon.ListColumn {
	defaultColSize := 5
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).Port < c2.(*DisplayRouteStats).Port)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		if stats.Port != 0 {
			return fmt.Sprintf("%5v", stats.Port)
		} else {
			return "   --"
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return fmt.Sprintf("%v", stats.Port)
	}
	c := uiCommon.NewListColumn("PORT", "PORT", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalRequests() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).HttpAllCount < c2.(*DisplayRouteStats).HttpAllCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		value := ""
		if stats.Port == 0 {
			value = util.Format(stats.HttpAllCount)
		} else {
			value = "n/a"
		}
		return fmt.Sprintf("%10v", value)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return fmt.Sprintf("%v", stats.HttpAllCount)
	}
	c := uiCommon.NewListColumn("TOTREQ", "TOT_REQ", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func column2xx() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).Http2xxCount < c2.(*DisplayRouteStats).Http2xxCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.Http2xxCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return fmt.Sprintf("%v", stats.Http2xxCount)
	}
	c := uiCommon.NewListColumn("2XX", "2XX", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func column3xx() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).Http3xxCount < c2.(*DisplayRouteStats).Http3xxCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.Http3xxCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return fmt.Sprintf("%v", stats.Http3xxCount)
	}
	c := uiCommon.NewListColumn("3XX", "3XX", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func column4xx() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).Http4xxCount < c2.(*DisplayRouteStats).Http4xxCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.Http4xxCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return fmt.Sprintf("%v", stats.Http4xxCount)
	}
	c := uiCommon.NewListColumn("4XX", "4XX", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func column5xx() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).Http5xxCount < c2.(*DisplayRouteStats).Http5xxCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.Http5xxCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return fmt.Sprintf("%v", stats.Http5xxCount)
	}
	c := uiCommon.NewListColumn("5XX", "5XX", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnResponseContentLength() *uiCommon.ListColumn {
	defaultColSize := 9
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).ResponseContentLength < c2.(*DisplayRouteStats).ResponseContentLength)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%9v", "--")
		} else {
			value := fmt.Sprintf("%9v", util.ByteSize(stats.ResponseContentLength).StringWithPrecision(1))
			return value
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return fmt.Sprintf("%v", stats.ResponseContentLength)
	}
	c := uiCommon.NewListColumn("RESP_DATA", "RESP_DATA", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)

	return c
}

func columnMethodGet() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).HttpMethodGetCount < c2.(*DisplayRouteStats).HttpMethodGetCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.HttpMethodGetCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return fmt.Sprintf("%v", stats.HttpMethodGetCount)
	}
	c := uiCommon.NewListColumn("M_GET", "M_GET", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnMethodPost() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).HttpMethodPostCount < c2.(*DisplayRouteStats).HttpMethodPostCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.HttpMethodPostCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return fmt.Sprintf("%v", stats.HttpMethodPostCount)
	}
	c := uiCommon.NewListColumn("M_POST", "M_POST", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnMethodPut() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).HttpMethodPutCount < c2.(*DisplayRouteStats).HttpMethodPutCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.HttpMethodPutCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return fmt.Sprintf("%v", stats.HttpMethodPutCount)
	}
	c := uiCommon.NewListColumn("M_PUT", "M_PUT", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnMethodDelete() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayRouteStats).HttpMethodDeleteCount < c2.(*DisplayRouteStats).HttpMethodDeleteCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%10v", "--")
		} else {
			return fmt.Sprintf("%10v", util.Format(stats.HttpMethodDeleteCount))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return fmt.Sprintf("%v", stats.HttpMethodDeleteCount)
	}
	c := uiCommon.NewListColumn("M_DELETE", "M_DELETE", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnLastAccess() *uiCommon.ListColumn {
	defaultColSize := 19
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayRouteStats).LastAccess.Before(c2.(*DisplayRouteStats).LastAccess)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*DisplayRouteStats)
		if stats.HttpAllCount == 0 {
			return fmt.Sprintf("%19v", "--")
		} else {
			return fmt.Sprintf("%19v", stats.LastAccess.Format("01-02-2006 15:04:05"))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayRouteStats)
		return fmt.Sprintf("%v", stats.LastAccess)
	}
	c := uiCommon.NewListColumn("LAST_ACCESS", "LAST_ACCESS", defaultColSize,
		uiCommon.TIMESTAMP, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}
