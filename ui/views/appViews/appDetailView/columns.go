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

package appDetailView

import (
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func stateAttentionFunc(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
	stats := data.(*DisplayContainerStats)
	attentionType := uiCommon.ATTENTION_NORMAL

	// "FLAPPING"
	switch stats.State {
	case "DOWN":
		attentionType = uiCommon.ATTENTION_STATE_DOWN
	case "TERM":
		attentionType = uiCommon.ATTENTION_STATE_TERM
	case "STARTING":
		attentionType = uiCommon.ATTENTION_STATE_STARTING
	case "UNKNOWN":
		attentionType = uiCommon.ATTENTION_STATE_UNKNOWN
	case "CRASHED":
		attentionType = uiCommon.ATTENTION_STATE_CRASHED
	case "RUNNING":
		attentionType = uiCommon.ATTENTION_NORMAL
	case "":
		attentionType = uiCommon.ATTENTION_NORMAL
	default:
		attentionType = uiCommon.ATTENTION_NORMAL
	}

	return attentionType
}

func uptimeAttentionFunc(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {

	attentionType := stateAttentionFunc(data, columnOwner)
	if attentionType == uiCommon.ATTENTION_NORMAL {
		stats := data.(*DisplayContainerStats)
		if stats.StateDuration != nil && stats.StateDuration.Seconds() < 60 {
			attentionType = uiCommon.ATTENTION_CONTAINER_SHORT_UPTIME
		}
	}
	return attentionType
}

func ColumnContainerIndex() *uiCommon.ListColumn {
	defaultColSize := 4
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayContainerStats).ContainerIndex < c2.(*DisplayContainerStats).ContainerIndex)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerStats)
		display := fmt.Sprintf("%4v", stats.ContainerIndex)
		return fmt.Sprintf("%4v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", stats.ContainerIndex)
	}
	c := uiCommon.NewListColumn("IDX", "IDX", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnTotalCpuPercentage() *uiCommon.ListColumn {
	defaultColSize := 6
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayContainerStats).ContainerMetric.GetCpuPercentage() < c2.(*DisplayContainerStats).ContainerMetric.GetCpuPercentage())
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerStats)
		totalCpuInfo := ""
		// We use MemoryBytes instead of CPU% as CPU% can be zero and if that is the case we want
		// to display 0.
		if stats.ContainerMetric.GetMemoryBytes() == 0 {
			totalCpuInfo = fmt.Sprintf("%6v", "--")
		} else {
			cpuPercentage := stats.ContainerMetric.GetCpuPercentage()
			if cpuPercentage >= 100.0 {
				totalCpuInfo = fmt.Sprintf("%6.0f", cpuPercentage)
			} else if cpuPercentage >= 10.0 {
				totalCpuInfo = fmt.Sprintf("%6.1f", cpuPercentage)
			} else {
				totalCpuInfo = fmt.Sprintf("%6.2f", cpuPercentage)
			}
		}
		return fmt.Sprintf("%6v", totalCpuInfo)

	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", stats.ContainerMetric.GetCpuPercentage())
	}
	c := uiCommon.NewListColumn("CPU_PERCENT", "CPU%", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)

	return c
}

func ColumnAppName() *uiCommon.ListColumn {
	defaultColSize := 50
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*DisplayContainerStats).AppName, c2.(*DisplayContainerStats).AppName)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayContainerStats)
		return util.FormatDisplayData(appStats.AppName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return appStats.AppName
	}
	c := uiCommon.NewListColumn("appName", "APPLICATION", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnSpaceName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*DisplayContainerStats).SpaceName, c2.(*DisplayContainerStats).SpaceName)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayContainerStats)
		return util.FormatDisplayData(appStats.SpaceName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return appStats.SpaceName
	}
	c := uiCommon.NewListColumn("spaceName", "SPACE", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnOrgName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*DisplayContainerStats).OrgName, c2.(*DisplayContainerStats).OrgName)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		appStats := data.(*DisplayContainerStats)
		return util.FormatDisplayData(appStats.OrgName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerStats)
		return appStats.OrgName
	}
	c := uiCommon.NewListColumn("orgName", "ORG", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnMemoryUsed() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).ContainerMetric.GetMemoryBytes() < c2.(*DisplayContainerStats).ContainerMetric.GetMemoryBytes()
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerStats)
		if stats.ContainerMetric.GetMemoryBytes() == 0 {
			return fmt.Sprintf("%9v", "--")
		} else {
			return fmt.Sprintf("%9v", util.ByteSize(stats.ContainerMetric.GetMemoryBytes()).StringWithPrecision(1))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.ContainerMetric.GetMemoryBytes())
	}
	c := uiCommon.NewListColumn("MEM_USED", "MEM_USED", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnMemoryFree() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).FreeMemory < c2.(*DisplayContainerStats).FreeMemory
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerStats)
		if stats.ContainerMetric.GetMemoryBytes() == 0 {
			return fmt.Sprintf("%9v", "--")
		} else {
			return fmt.Sprintf("%9v", util.ByteSize(stats.FreeMemory).StringWithPrecision(1))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.FreeMemory)
	}
	c := uiCommon.NewListColumn("MEM_FREE", "MEM_FREE", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnMemoryReserved() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).ReservedMemory < c2.(*DisplayContainerStats).ReservedMemory
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerStats)
		if stats.ContainerMetric.GetMemoryBytes() == 0 {
			return fmt.Sprintf("%9v", "--")
		} else {
			return fmt.Sprintf("%9v", util.ByteSize(stats.ReservedMemory).StringWithPrecision(1))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.ReservedMemory)
	}
	c := uiCommon.NewListColumn("MEM_RSVD", "MEM_RSVD", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnDiskUsed() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).ContainerMetric.GetDiskBytes() < c2.(*DisplayContainerStats).ContainerMetric.GetDiskBytes()
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerStats)
		if stats.ContainerMetric.GetMemoryBytes() == 0 {
			return fmt.Sprintf("%9v", "--")
		} else {
			return fmt.Sprintf("%9v", util.ByteSize(stats.ContainerMetric.GetDiskBytes()).StringWithPrecision(1))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.ContainerMetric.GetDiskBytes())
	}
	c := uiCommon.NewListColumn("DISK_USED", "DISK_USED", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnDiskFree() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).FreeDisk < c2.(*DisplayContainerStats).FreeDisk
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerStats)
		if stats.ContainerMetric.GetMemoryBytes() == 0 {
			return fmt.Sprintf("%9v", "--")
		} else {
			return fmt.Sprintf("%9v", util.ByteSize(stats.FreeDisk).StringWithPrecision(1))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.FreeDisk)
	}
	c := uiCommon.NewListColumn("DISK_FREE", "DISK_FREE", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnDiskReserved() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).ReservedDisk < c2.(*DisplayContainerStats).ReservedDisk
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerStats)
		if stats.ContainerMetric.GetMemoryBytes() == 0 {
			return fmt.Sprintf("%9v", "--")
		} else {
			return fmt.Sprintf("%9v", util.ByteSize(stats.ReservedDisk).StringWithPrecision(1))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.ReservedDisk)
	}
	c := uiCommon.NewListColumn("DISK_RSVD", "DISK_RSVD", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func columnAvgResponseTimeL60Info() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).AvgResponseL60Time < c2.(*DisplayContainerStats).AvgResponseL60Time
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		containerStats := data.(*DisplayContainerStats)
		avgResponseTimeL60Info := "--"
		if containerStats.AvgResponseL60Time > 0 {
			avgResponseTimeMs := containerStats.AvgResponseL60Time / 1000000
			if avgResponseTimeMs >= 10 {
				avgResponseTimeL60Info = fmt.Sprintf("%6.0f", avgResponseTimeMs)
			} else if avgResponseTimeMs >= 1 {
				avgResponseTimeL60Info = fmt.Sprintf("%6.1f", avgResponseTimeMs)
			} else {
				avgResponseTimeL60Info = fmt.Sprintf("%6.2f", avgResponseTimeMs)
			}
		}
		return fmt.Sprintf("%6v", avgResponseTimeL60Info)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.AvgResponseL60Time)
	}
	c := uiCommon.NewListColumn("RESP", "RESP", 6,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnLogStdout() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).OutCount < c2.(*DisplayContainerStats).OutCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%11v", util.Format(stats.OutCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.OutCount)
	}
	c := uiCommon.NewListColumn("LOG_OUT", "LOG_OUT", 11,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnLogStderr() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).ErrCount < c2.(*DisplayContainerStats).ErrCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%11v", util.Format(stats.ErrCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.ErrCount)
	}
	c := uiCommon.NewListColumn("LOG_ERR", "LOG_ERR", 11,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func columnReq1() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).EventL1Rate < c2.(*DisplayContainerStats).EventL1Rate
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%6v", util.Format(int64(containerStats.EventL1Rate)))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.EventL1Rate)
	}
	attentionFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
		containerStats := data.(*DisplayContainerStats)
		attentionType := stateAttentionFunc(data, columnOwner)
		if attentionType == uiCommon.ATTENTION_NORMAL && containerStats.EventL1Rate > 0 {
			attentionType = uiCommon.ATTENTION_ACTIVITY
		}
		return attentionType
	}
	c := uiCommon.NewListColumn("REQ1", "REQ/1", 6,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, attentionFunc)
	return c
}

func columnReq10() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).EventL10Rate < c2.(*DisplayContainerStats).EventL10Rate
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%7v", util.Format(int64(containerStats.EventL10Rate)))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.EventL10Rate)
	}
	attentionFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
		containerStats := data.(*DisplayContainerStats)
		attentionType := stateAttentionFunc(data, columnOwner)
		if attentionType == uiCommon.ATTENTION_NORMAL && containerStats.EventL10Rate > 0 {
			attentionType = uiCommon.ATTENTION_ACTIVITY
		}
		return attentionType
	}
	c := uiCommon.NewListColumn("REQ10", "REQ/10", 7,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, attentionFunc)
	return c
}

func columnReq60() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).EventL60Rate < c2.(*DisplayContainerStats).EventL60Rate
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%7v", util.Format(int64(containerStats.EventL60Rate)))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.EventL60Rate)
	}
	attentionFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
		containerStats := data.(*DisplayContainerStats)
		attentionType := stateAttentionFunc(data, columnOwner)
		if attentionType == uiCommon.ATTENTION_NORMAL && containerStats.EventL60Rate > 0 {
			attentionType = uiCommon.ATTENTION_ACTIVITY
		}
		return attentionType
	}
	c := uiCommon.NewListColumn("REQ60", "REQ/60", 7,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, attentionFunc)
	return c
}

func columnTotalReq() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).HttpAllCount < c2.(*DisplayContainerStats).HttpAllCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%10v", util.Format(containerStats.HttpAllCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.HttpAllCount)
	}
	c := uiCommon.NewListColumn("TOT_REQ", "TOT_REQ", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}
func column2XX() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).Http2xxCount < c2.(*DisplayContainerStats).Http2xxCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%10v", util.Format(containerStats.Http2xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.Http2xxCount)
	}
	c := uiCommon.NewListColumn("2XX", "2XX", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}
func column3XX() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).Http3xxCount < c2.(*DisplayContainerStats).Http3xxCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%10v", util.Format(containerStats.Http3xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.Http3xxCount)
	}
	c := uiCommon.NewListColumn("3XX", "3XX", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func column4XX() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).Http4xxCount < c2.(*DisplayContainerStats).Http4xxCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%10v", util.Format(containerStats.Http4xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.Http4xxCount)
	}
	c := uiCommon.NewListColumn("4XX", "4XX", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func column5XX() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).Http5xxCount < c2.(*DisplayContainerStats).Http5xxCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%10v", util.Format(containerStats.Http5xxCount))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.Http5xxCount)
	}
	c := uiCommon.NewListColumn("5XX", "5XX", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnCellIp() *uiCommon.ListColumn {
	defaultColSize := 16
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.Ip2long(c1.(*DisplayContainerStats).Ip) < util.Ip2long(c2.(*DisplayContainerStats).Ip)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		containerStats := data.(*DisplayContainerStats)
		ip := containerStats.Ip
		if ip == "" {
			ip = "--"
		}
		return util.FormatDisplayData(ip, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return containerStats.Ip
	}
	c := uiCommon.NewListColumn("CELL_IP", "CELL_IP", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnState() *uiCommon.ListColumn {
	defaultColSize := 9
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayContainerStats).State) < (c2.(*DisplayContainerStats).State)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		containerStats := data.(*DisplayContainerStats)
		msgText := containerStats.State
		if msgText == "" {
			msgText = "--"
		}
		return util.FormatDisplayData(msgText, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return containerStats.State
	}
	c := uiCommon.NewListColumn("STATE", "STATE", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnStateDuration() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {

		uptime1 := c1.(*DisplayContainerStats).StateDuration
		uptime2 := c2.(*DisplayContainerStats).StateDuration
		if uptime1 == nil {
			return true
		}
		if uptime2 == nil {
			return false
		}
		return uptime1.Seconds() < uptime2.Seconds()
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerStats)
		if stats.StateDuration == nil {
			return fmt.Sprintf("%11v", "--")
		} else {
			alwaysShowSeconds := stats.State == "STARTING"
			return fmt.Sprintf("%11v", util.FormatDuration(stats.StateDuration, alwaysShowSeconds))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.StateDuration.Seconds())
	}
	c := uiCommon.NewListColumn("STATE_DUR", "STATE_DUR", 11,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, uptimeAttentionFunc)
	return c
}

func ColumnStartupDuration() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {

		uptime1 := c1.(*DisplayContainerStats).StartupDuration
		uptime2 := c2.(*DisplayContainerStats).StartupDuration
		if uptime1 == nil {
			return true
		}
		if uptime2 == nil {
			return false
		}
		return uptime1.Seconds() < uptime2.Seconds()
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerStats)
		if stats.StartupDuration == nil {
			return fmt.Sprintf("%8v", "--")
		} else {
			return fmt.Sprintf("%8v", util.FormatDuration(stats.StartupDuration, true))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.StartupDuration.Seconds())
	}
	c := uiCommon.NewListColumn("STRT_DUR", "STRT_DUR", 8,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnCreateCount() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerStats).CreateCount < c2.(*DisplayContainerStats).CreateCount
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%4v", stats.CreateCount)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", containerStats.CreateCount)
	}
	c := uiCommon.NewListColumn("CCR", "CCR", 4,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnStateTime() *uiCommon.ListColumn {
	defaultColSize := 19
	sortFunc := func(c1, c2 util.Sortable) bool {
		t1 := c1.(*DisplayContainerStats).StateTime
		t2 := c2.(*DisplayContainerStats).StateTime
		if t1 == nil {
			return true
		}
		if t2 == nil {
			return false
		}
		return t1.Before(*t2)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		containerStats := data.(*DisplayContainerStats)
		if containerStats.StateTime == nil {
			return fmt.Sprintf("%-19v", "--")
		} else {
			return fmt.Sprintf("%-19v", containerStats.StateTime.Format("01-02-2006 15:04:05"))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", stats.StateTime)
	}
	c := uiCommon.NewListColumn("STATE_TIME", "STATE_TIME", defaultColSize,
		uiCommon.TIMESTAMP, true, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnCellLastStartMsgText() *uiCommon.ListColumn {
	defaultColSize := 57
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayContainerStats).CellLastStartMsgText) < (c2.(*DisplayContainerStats).CellLastStartMsgText)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		containerStats := data.(*DisplayContainerStats)
		msgText := containerStats.CellLastStartMsgText
		if msgText == "" {
			msgText = "--"
		}
		return util.FormatDisplayData(msgText, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		containerStats := data.(*DisplayContainerStats)
		return containerStats.CellLastStartMsgText
	}
	c := uiCommon.NewListColumn("CNTR_START_MSG", "CNTR_START_MSG", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}

func ColumnCellLastStartMsgTime() *uiCommon.ListColumn {
	defaultColSize := 19
	sortFunc := func(c1, c2 util.Sortable) bool {
		t1 := c1.(*DisplayContainerStats).CellLastStartMsgTime
		t2 := c2.(*DisplayContainerStats).CellLastStartMsgTime
		if t1 == nil {
			return true
		}
		if t2 == nil {
			return false
		}
		if t1 == nil && t2 == nil {
			return false
		}
		return t1.Before(*t2)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		containerStats := data.(*DisplayContainerStats)
		msgText := containerStats.CellLastStartMsgText
		if msgText == "" || containerStats.CellLastStartMsgTime == nil {
			return fmt.Sprintf("%-19v", "--")
		} else {
			return fmt.Sprintf("%-19v", containerStats.CellLastStartMsgTime.Format("01-02-2006 15:04:05"))
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayContainerStats)
		return fmt.Sprintf("%v", stats.CellLastStartMsgTime)
	}
	c := uiCommon.NewListColumn("CNTR_START_MSG_TM", "CNTR_START_MSG_TM", defaultColSize,
		uiCommon.TIMESTAMP, true, sortFunc, true, displayFunc, rawValueFunc, stateAttentionFunc)
	return c
}
