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

package eventDetailView

import (
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func columnDeploymentName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*displaydata.DisplayEventDetailStats).DeploymentName, c2.(*displaydata.DisplayEventDetailStats).DeploymentName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayEventDetailStats)
		return util.FormatDisplayData(cellStats.DeploymentName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayEventDetailStats)
		return cellStats.DeploymentName
	}
	c := uiCommon.NewListColumn("DNAME", "DNAME", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func columnJobName() *uiCommon.ListColumn {
	defaultColSize := 45
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*displaydata.DisplayEventDetailStats).JobName, c2.(*displaydata.DisplayEventDetailStats).JobName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayEventDetailStats)
		return util.FormatDisplayData(cellStats.JobName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayEventDetailStats)
		return cellStats.JobName
	}
	c := uiCommon.NewListColumn("JOB_NAME", "JOB_NAME", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

// Job Index in PCF 1.8 is now a GUID not a integer number
func columnJobIndex() *uiCommon.ListColumn {
	defaultColSize := 36
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayEventDetailStats).JobIndex < c2.(*displaydata.DisplayEventDetailStats).JobIndex
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayEventDetailStats)
		return util.FormatDisplayData(cellStats.JobIndex, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayEventDetailStats)
		return cellStats.JobIndex
	}
	c := uiCommon.NewListColumn("JOB_IDX", "JOB_IDX", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func columnIp() *uiCommon.ListColumn {
	defaultColSize := 16
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.Ip2long(c1.(*displaydata.DisplayEventDetailStats).Ip) < util.Ip2long(c2.(*displaydata.DisplayEventDetailStats).Ip)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayEventDetailStats)
		return util.FormatDisplayData(cellStats.Ip, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayEventDetailStats)
		return cellStats.Ip
	}
	c := uiCommon.NewListColumn("IP", "IP", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func columnEventCount() *uiCommon.ListColumn {
	defaultColSize := 12
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayEventDetailStats).EventCount < c2.(*displaydata.DisplayEventDetailStats).EventCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		stats := data.(*displaydata.DisplayEventDetailStats)
		return fmt.Sprintf("%12v", util.Format(stats.EventCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*displaydata.DisplayEventDetailStats)
		return fmt.Sprintf("%v", stats.EventCount)
	}
	c := uiCommon.NewListColumn("COUNT", "COUNT", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}
