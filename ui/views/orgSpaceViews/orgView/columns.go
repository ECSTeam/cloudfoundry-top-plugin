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

package orgView

import (
	"fmt"
	"strconv"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func columnOrgName() *uiCommon.ListColumn {
	defaultColSize := 25
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*DisplayOrg).Name, c2.(*DisplayOrg).Name)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayOrg)
		return util.FormatDisplayData(stats.Name, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayOrg)
		return stats.Name
	}
	c := uiCommon.NewListColumn("orgName", "ORG", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnStatus() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*DisplayOrg).Status, c2.(*DisplayOrg).Status)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayOrg)
		return util.FormatDisplayData(stats.Status, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayOrg)
		return stats.Status
	}
	c := uiCommon.NewListColumn("STATUS", "STATUS", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnQuotaName() *uiCommon.ListColumn {
	defaultColSize := 11
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*DisplayOrg).QuotaName, c2.(*DisplayOrg).QuotaName)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayOrg)
		return util.FormatDisplayData(stats.QuotaName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayOrg)
		return stats.QuotaName
	}
	c := uiCommon.NewListColumn("QUOTA_NAME", "QUOTA_NAME", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnNumberOfSpaces() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).NumberOfSpaces < c2.(*DisplayOrg).NumberOfSpaces
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayOrg)
		return fmt.Sprintf("%7v", stats.NumberOfSpaces)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayOrg)
		return strconv.Itoa(stats.NumberOfSpaces)
	}
	c := uiCommon.NewListColumn("SPACES", "SPACES", 7,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnNumberOfApps() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).NumberOfApps < c2.(*DisplayOrg).NumberOfApps
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayOrg)
		return fmt.Sprintf("%7v", stats.NumberOfApps)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayOrg)
		return strconv.Itoa(stats.NumberOfApps)
	}
	c := uiCommon.NewListColumn("APPS", "APPS", 7,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnReportingContainers() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).TotalReportingContainers < c2.(*DisplayOrg).TotalReportingContainers
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%7v", appStats.TotalReportingContainers)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return strconv.Itoa(appStats.TotalReportingContainers)
	}
	c := uiCommon.NewListColumn("reportingContainers", "RCR", 7,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalCpu() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).TotalCpuPercentage < c2.(*DisplayOrg).TotalCpuPercentage
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		totalCpuInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalCpuInfo = fmt.Sprintf("%6v", "--")
		} else {
			if appStats.TotalCpuPercentage >= 100.0 {
				totalCpuInfo = fmt.Sprintf("%6.0f", appStats.TotalCpuPercentage)
			} else if appStats.TotalCpuPercentage >= 10.0 {
				totalCpuInfo = fmt.Sprintf("%6.1f", appStats.TotalCpuPercentage)
			} else {
				totalCpuInfo = fmt.Sprintf("%6.2f", appStats.TotalCpuPercentage)
			}
		}
		return fmt.Sprintf("%6v", totalCpuInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%.2f", appStats.TotalCpuPercentage)
	}
	c := uiCommon.NewListColumn("CPU", "CPU%", 6,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}
func columnMemoryLimit() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).MemoryLimitInBytes < c2.(*DisplayOrg).MemoryLimitInBytes
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		totalMemInfo := ""
		//if appStats.MemoryLimitInBytes == 0 {
		//	totalMemInfo = fmt.Sprintf("%9v", "--")
		//} else {
		totalMemInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.MemoryLimitInBytes).StringWithPrecision(1))
		//}
		return fmt.Sprintf("%9v", totalMemInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%v", appStats.MemoryLimitInBytes)
	}
	c := uiCommon.NewListColumn("MAX_MEM", "MAX_MEM", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalMemoryUsed() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).TotalUsedMemory < c2.(*DisplayOrg).TotalUsedMemory
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		totalMemInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalMemInfo = fmt.Sprintf("%9v", "--")
		} else {
			totalMemInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.TotalUsedMemory).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", totalMemInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%v", appStats.TotalUsedMemory)
	}
	c := uiCommon.NewListColumn("USED_MEM", "USED_MEM", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalUsedMemoryPercentOfQuota() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).TotalUsedMemoryPercentOfQuota < c2.(*DisplayOrg).TotalUsedMemoryPercentOfQuota
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		totalMemInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalMemInfo = fmt.Sprintf("%5v", "--")
		} else {
			totalMemInfo = fmt.Sprintf("%5.1f", appStats.TotalUsedMemoryPercentOfQuota)
		}
		return fmt.Sprintf("%5v", totalMemInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%v", appStats.TotalUsedMemoryPercentOfQuota)
	}
	c := uiCommon.NewListColumn("QUOTA_MEM_PER", "MEM%", 5,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalDiskUsed() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).TotalUsedDisk < c2.(*DisplayOrg).TotalUsedDisk
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		totalDiskInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalDiskInfo = fmt.Sprintf("%10v", "--")
		} else {
			totalDiskInfo = fmt.Sprintf("%10v", util.ByteSize(appStats.TotalUsedDisk).StringWithPrecision(1))
		}
		return fmt.Sprintf("%10v", totalDiskInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%v", appStats.TotalUsedDisk)
	}
	c := uiCommon.NewListColumn("USED_DISK", "USED_DISK", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnLogStdout() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).TotalLogStdout < c2.(*DisplayOrg).TotalLogStdout
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%11v", util.Format(appStats.TotalLogStdout))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%v", appStats.TotalLogStdout)
	}
	c := uiCommon.NewListColumn("TotalLogStdout", "LOG_OUT", 11,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnLogStderr() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).TotalLogStderr < c2.(*DisplayOrg).TotalLogStderr
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%11v", util.Format(appStats.TotalLogStderr))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%v", appStats.TotalLogStderr)
	}
	c := uiCommon.NewListColumn("TotalLogStderr", "LOG_ERR", 11,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalReq() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).HttpAllCount < c2.(*DisplayOrg).HttpAllCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%11v", util.Format(appStats.HttpAllCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%v", appStats.HttpAllCount)
	}
	c := uiCommon.NewListColumn("TOTREQ", "TOT_REQ", 11,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}
