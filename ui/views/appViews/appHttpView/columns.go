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

package appHttpView

import (
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func ColumnMethod() *uiCommon.ListColumn {
	defaultColSize := 8
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayHttpInfo).HttpMethod < (c2.(*DisplayHttpInfo).HttpMethod)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		info := data.(*DisplayHttpInfo)
		return fmt.Sprintf("%-8v", info.HttpMethod)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		info := data.(*DisplayHttpInfo)
		return fmt.Sprintf("%v", info.HttpMethod)
	}
	c := uiCommon.NewListColumn("METHOD", "METHOD", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func ColumnStatusCode() *uiCommon.ListColumn {
	defaultColSize := 5
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayHttpInfo).HttpStatusCode < c2.(*DisplayHttpInfo).HttpStatusCode)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayHttpInfo)
		return fmt.Sprintf("%5v", stats.HttpStatusCode)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayHttpInfo)
		return fmt.Sprintf("%v", stats.HttpStatusCode)
	}
	c := uiCommon.NewListColumn("CODE", "CODE", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func ColumnLastAcivity() *uiCommon.ListColumn {
	defaultColSize := 20
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayHttpInfo).LastAcivity.Before(*c2.(*DisplayHttpInfo).LastAcivity))
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayHttpInfo)
		return fmt.Sprintf("%20v", stats.LastAcivityFormatted)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayHttpInfo)
		return fmt.Sprintf("%v", stats.LastAcivity.UnixNano())
	}
	c := uiCommon.NewListColumn("LAST_ACTIVITY", "LAST_ACTIVITY", defaultColSize,
		uiCommon.TIMESTAMP, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func ColumnCount() *uiCommon.ListColumn {
	defaultColSize := 11
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayHttpInfo).HttpCount < c2.(*DisplayHttpInfo).HttpCount)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayHttpInfo)
		return fmt.Sprintf("%11v", util.Format(stats.HttpCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayHttpInfo)
		return fmt.Sprintf("%v", stats.HttpCount)
	}
	c := uiCommon.NewListColumn("COUNT", "COUNT", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}
