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

const ATTENTION_HOT_PERCENT = 90
const ATTENTION_WARM_PERCENT = 80

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
	c := uiCommon.NewListColumn("ORG", "ORG", defaultColSize,
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

func columnDesiredContainers() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).DesiredContainers < c2.(*DisplayOrg).DesiredContainers
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%7v", appStats.DesiredContainers)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return strconv.Itoa(appStats.DesiredContainers)
	}
	c := uiCommon.NewListColumn("DCR", "DCR", 7,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, notInDesiredStateAttentionFunc)
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
	c := uiCommon.NewListColumn("RCR", "RCR", 7,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, notInDesiredStateAttentionFunc)
	return c
}

func notInDesiredStateAttentionFunc(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
	stats := data.(*DisplayOrg)
	appListView := columnOwner.(*OrgListView)
	attentionType := uiCommon.ATTENTION_NORMAL
	if appListView.isWarmupComplete && stats.DesiredContainers > stats.TotalReportingContainers {
		attentionType = uiCommon.ATTENTION_NOT_DESIRED_STATE
	}
	return attentionType
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
	c := uiCommon.NewListColumn("CPU_PER", "CPU%", 6,
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
		if appStats.MemoryLimitInBytes == 0 {
			totalMemInfo = fmt.Sprintf("%9v", "--")
		} else {
			totalMemInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.MemoryLimitInBytes).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", totalMemInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%v", appStats.MemoryLimitInBytes)
	}
	c := uiCommon.NewListColumn("MEM_MAX", "MEM_MAX", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalMemoryReserved() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).TotalMemoryReserved < c2.(*DisplayOrg).TotalMemoryReserved
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		totalMemInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalMemInfo = fmt.Sprintf("%9v", "--")
		} else {
			totalMemInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.TotalMemoryReserved).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", totalMemInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%v", appStats.TotalMemoryReserved)
	}
	c := uiCommon.NewListColumn("MEM_RSVD", "MEM_RSVD", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, closeToMemoryQuotaAttentionFunc)
	return c
}

func columnTotalMemoryReservedPercentOfQuota() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).TotalMemoryReservedPercentOfQuota < c2.(*DisplayOrg).TotalMemoryReservedPercentOfQuota
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		totalMemInfo := ""
		if appStats.TotalReportingContainers == 0 || appStats.MemoryLimitInBytes == 0 {
			totalMemInfo = fmt.Sprintf("%7v", "--")
		} else {
			totalMemInfo = fmt.Sprintf("%7.1f", appStats.TotalMemoryReservedPercentOfQuota)
		}
		return fmt.Sprintf("%7v", totalMemInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%v", appStats.TotalMemoryReservedPercentOfQuota)
	}
	c := uiCommon.NewListColumn("O_MEM_PER", "O_MEM%", 7,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, closeToMemoryQuotaAttentionFunc)
	return c
}

func closeToMemoryQuotaAttentionFunc(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
	stats := data.(*DisplayOrg)
	attentionType := uiCommon.ATTENTION_NORMAL
	if stats.TotalReportingContainers > 0 && stats.MemoryLimitInBytes > 0 {
		percentOfQuota := stats.TotalMemoryReservedPercentOfQuota
		switch {
		case percentOfQuota >= ATTENTION_HOT_PERCENT:
			attentionType = uiCommon.ATTENTION_HOT
		case percentOfQuota >= ATTENTION_WARM_PERCENT:
			attentionType = uiCommon.ATTENTION_WARM
		}
	}
	return attentionType
}

func columnTotalMemoryUsed() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).TotalMemoryUsed < c2.(*DisplayOrg).TotalMemoryUsed
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		totalMemInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalMemInfo = fmt.Sprintf("%9v", "--")
		} else {
			totalMemInfo = fmt.Sprintf("%9v", util.ByteSize(appStats.TotalMemoryUsed).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", totalMemInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%v", appStats.TotalMemoryUsed)
	}
	c := uiCommon.NewListColumn("MEM_USED", "MEM_USED", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalDiskReserved() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).TotalDiskReserved < c2.(*DisplayOrg).TotalDiskReserved
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		totalDiskInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalDiskInfo = fmt.Sprintf("%10v", "--")
		} else {
			totalDiskInfo = fmt.Sprintf("%10v", util.ByteSize(appStats.TotalDiskReserved).StringWithPrecision(1))
		}
		return fmt.Sprintf("%10v", totalDiskInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%v", appStats.TotalDiskReserved)
	}
	c := uiCommon.NewListColumn("DSK_RSVD", "DSK_RSVD", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalDiskUsed() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayOrg).TotalDiskUsed < c2.(*DisplayOrg).TotalDiskUsed
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayOrg)
		totalDiskInfo := ""
		if appStats.TotalReportingContainers == 0 {
			totalDiskInfo = fmt.Sprintf("%10v", "--")
		} else {
			totalDiskInfo = fmt.Sprintf("%10v", util.ByteSize(appStats.TotalDiskUsed).StringWithPrecision(1))
		}
		return fmt.Sprintf("%10v", totalDiskInfo)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayOrg)
		return fmt.Sprintf("%v", appStats.TotalDiskUsed)
	}
	c := uiCommon.NewListColumn("DSK_USED", "DSK_USED", 10,
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
	c := uiCommon.NewListColumn("LOG_OUT", "LOG_OUT", 11,
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
	c := uiCommon.NewListColumn("LOG_ERR", "LOG_ERR", 11,
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
	c := uiCommon.NewListColumn("TOT_REQ", "TOT_REQ", 11,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}
