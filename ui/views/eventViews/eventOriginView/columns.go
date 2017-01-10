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

package eventOriginView

import (
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func columnEventOrigin() *uiCommon.ListColumn {
	defaultColSize := 25
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayEventOriginStats).Origin < c2.(*DisplayEventOriginStats).Origin
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayEventOriginStats)
		return util.FormatDisplayData(stats.Origin, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayEventOriginStats)
		return stats.Origin
	}
	c := uiCommon.NewListColumn("ORIGIN", "ORIGIN", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnEventCount() *uiCommon.ListColumn {
	defaultColSize := 12
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayEventOriginStats).EventCount < c2.(*DisplayEventOriginStats).EventCount)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayEventOriginStats)
		return fmt.Sprintf("%12v", util.Format(stats.EventCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayEventOriginStats)
		return fmt.Sprintf("%v", stats.EventCount)
	}
	c := uiCommon.NewListColumn("COUNT", "COUNT", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}
