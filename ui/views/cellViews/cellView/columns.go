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

package cellView

import (
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func columnCellIp() *uiCommon.ListColumn {
	defaultColSize := 16
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.Ip2long(c1.(*DisplayCellStats).Ip) < util.Ip2long(c2.(*DisplayCellStats).Ip)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*DisplayCellStats)
		return util.FormatDisplayData(cellStats.Ip, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*DisplayCellStats)
		return cellStats.Ip
	}
	c := uiCommon.NewListColumn("CELL_IP", "CELL_IP", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalCpuPercentage() *uiCommon.ListColumn {
	defaultColSize := 6
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayCellStats).TotalContainerCpuPercentage < c2.(*DisplayCellStats).TotalContainerCpuPercentage)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*DisplayCellStats)

		totalCpuInfo := ""
		if cellStats.TotalReportingContainers == 0 {
			totalCpuInfo = fmt.Sprintf("%6v", "--")
		} else {
			cpuPercentage := cellStats.TotalContainerCpuPercentage

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
		cellStats := data.(*DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.TotalContainerCpuPercentage)
	}
	attentionFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
		cellStats := data.(*DisplayCellStats)
		attentionType := uiCommon.ATTENTION_NORMAL

		cpuPercentage := cellStats.TotalContainerCpuPercentage
		// This is the overall percentage of CPU in use on this cell, where 100% means all CPUs are 100% consumed
		cellMaxCpuCapacity := cpuPercentage / float64(cellStats.NumOfCpus)
		switch {
		case cellMaxCpuCapacity >= 90:
			attentionType = uiCommon.ATTENTION_HOT
		case cellMaxCpuCapacity >= 80:
			attentionType = uiCommon.ATTENTION_WARM
		}
		return attentionType
	}
	c := uiCommon.NewListColumn("CPU_PERCENT", "CPU%", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, attentionFunc)

	return c
}

func columnTotalReportingContainers() *uiCommon.ListColumn {
	defaultColSize := 4
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayCellStats).TotalReportingContainers < c2.(*DisplayCellStats).TotalReportingContainers)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*DisplayCellStats)
		display := ""
		if cellStats.TotalReportingContainers == 0 {
			display = fmt.Sprintf("%4v", "--")
		} else {
			display = fmt.Sprintf("%4v", cellStats.TotalReportingContainers)
		}
		return fmt.Sprintf("%4v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.TotalReportingContainers)
	}
	c := uiCommon.NewListColumn("RCR", "RCR", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}
func columnNumOfCpus() *uiCommon.ListColumn {
	defaultColSize := 4
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayCellStats).NumOfCpus < c2.(*DisplayCellStats).NumOfCpus)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*DisplayCellStats)
		display := ""
		if cellStats.NumOfCpus == 0 {
			display = fmt.Sprintf("%4v", "--")
		} else {
			display = fmt.Sprintf("%4v", cellStats.NumOfCpus)
		}
		return fmt.Sprintf("%4v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.NumOfCpus)
	}
	c := uiCommon.NewListColumn("CPUS", "CPUS", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnCapacityTotalMemory() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayCellStats).CapacityTotalMemory < c2.(*DisplayCellStats).CapacityTotalMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*DisplayCellStats)
		display := ""
		if CellStats.CapacityTotalMemory == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.CapacityTotalMemory).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.CapacityTotalMemory)
	}
	c := uiCommon.NewListColumn("TOT_MEM", "TOT_MEM", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnCapacityRemainingMemory() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayCellStats).CapacityRemainingMemory < c2.(*DisplayCellStats).CapacityRemainingMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*DisplayCellStats)
		display := ""
		if cellStats.CapacityRemainingMemory == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(cellStats.CapacityRemainingMemory).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityRemainingMemory)
	}
	attentionFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
		cellStats := data.(*DisplayCellStats)
		attentionType := uiCommon.ATTENTION_NORMAL
		if cellStats.CapacityTotalMemory > 0 {
			cellCapacity := (1 - (float64(cellStats.CapacityRemainingMemory) / float64(cellStats.CapacityTotalMemory))) * 100
			switch {
			case cellCapacity >= 90:
				attentionType = uiCommon.ATTENTION_HOT
			case cellCapacity >= 80:
				attentionType = uiCommon.ATTENTION_WARM
			}
		}
		return attentionType
	}
	c := uiCommon.NewListColumn("FREE_MEM", "FREE_MEM", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, attentionFunc)
	return c
}

func columnCapacityTotalDisk() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayCellStats).CapacityTotalDisk < c2.(*DisplayCellStats).CapacityTotalDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*DisplayCellStats)
		display := ""
		if CellStats.CapacityTotalDisk == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.CapacityTotalDisk).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.CapacityTotalDisk)
	}
	c := uiCommon.NewListColumn("TOT_DISK", "TOT_DISK", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnCapacityRemainingDisk() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayCellStats).CapacityRemainingDisk < c2.(*DisplayCellStats).CapacityRemainingDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*DisplayCellStats)
		display := ""
		if CellStats.CapacityRemainingDisk == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.CapacityRemainingDisk).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.CapacityRemainingDisk)
	}
	attentionFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) uiCommon.AttentionType {
		cellStats := data.(*DisplayCellStats)
		attentionType := uiCommon.ATTENTION_NORMAL
		if cellStats.CapacityTotalDisk > 0 {
			cellCapacity := (1 - (float64(cellStats.CapacityRemainingDisk) / float64(cellStats.CapacityTotalDisk))) * 100
			switch {
			case cellCapacity >= 90:
				attentionType = uiCommon.ATTENTION_HOT
			case cellCapacity >= 80:
				attentionType = uiCommon.ATTENTION_WARM
			}
		}
		return attentionType
	}
	c := uiCommon.NewListColumn("FREE_DISK", "FREE_DISK", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, attentionFunc)
	return c
}

func columnCapacityTotalContainers() *uiCommon.ListColumn {
	defaultColSize := 8
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayCellStats).CapacityTotalContainers < c2.(*DisplayCellStats).CapacityTotalContainers)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*DisplayCellStats)
		display := ""
		if cellStats.CapacityTotalContainers == 0 {
			display = fmt.Sprintf("%8v", "--")
		} else {
			display = fmt.Sprintf("%8v", cellStats.CapacityTotalContainers)
		}
		return fmt.Sprintf("%8v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityTotalContainers)
	}
	c := uiCommon.NewListColumn("MAX_CNTR", "MAX_CNTR", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnContainerCount() *uiCommon.ListColumn {
	defaultColSize := 5
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayCellStats).ContainerCount < c2.(*DisplayCellStats).ContainerCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*DisplayCellStats)
		display := ""
		if cellStats.ContainerCount == 0 {
			display = fmt.Sprintf("%5v", "--")
		} else {
			display = fmt.Sprintf("%5v", cellStats.ContainerCount)
		}
		return fmt.Sprintf("%5v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.ContainerCount)
	}
	c := uiCommon.NewListColumn("CNTRS", "CNTRS", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalContainerReservedMemory() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayCellStats).TotalContainerReservedMemory < c2.(*DisplayCellStats).TotalContainerReservedMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*DisplayCellStats)
		display := ""
		if CellStats.TotalContainerReservedMemory == 0 {
			display = fmt.Sprintf("%10v", "--")
		} else {
			display = fmt.Sprintf("%10v", util.ByteSize(CellStats.TotalContainerReservedMemory).StringWithPrecision(1))
		}
		return fmt.Sprintf("%10v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.TotalContainerReservedMemory)
	}
	c := uiCommon.NewListColumn("C_RSVD_MEM", "C_RSVD_MEM", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalContainerUsedMemory() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayCellStats).TotalContainerUsedMemory < c2.(*DisplayCellStats).TotalContainerUsedMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*DisplayCellStats)
		display := ""
		if CellStats.TotalContainerUsedMemory == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.TotalContainerUsedMemory).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.TotalContainerUsedMemory)
	}
	c := uiCommon.NewListColumn("C_USD_MEM", "C_USD_MEM", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalContainerReservedDisk() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayCellStats).TotalContainerReservedDisk < c2.(*DisplayCellStats).TotalContainerReservedDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*DisplayCellStats)
		display := ""
		if CellStats.TotalContainerReservedDisk == 0 {
			display = fmt.Sprintf("%10v", "--")
		} else {
			display = fmt.Sprintf("%10v", util.ByteSize(CellStats.TotalContainerReservedDisk).StringWithPrecision(1))
		}
		return fmt.Sprintf("%10v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.TotalContainerReservedDisk)
	}
	c := uiCommon.NewListColumn("C_RSVD_DSK", "C_RSVD_DSK", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalContainerUsedDisk() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayCellStats).TotalContainerUsedDisk < c2.(*DisplayCellStats).TotalContainerUsedDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*DisplayCellStats)
		display := ""
		if CellStats.TotalContainerUsedDisk == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.TotalContainerUsedDisk).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.TotalContainerUsedDisk)
	}
	c := uiCommon.NewListColumn("C_USD_DSK", "C_USD_DSK", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnStackName() *uiCommon.ListColumn {
	defaultColSize := 15
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*DisplayCellStats).StackName, c2.(*DisplayCellStats).StackName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		appStats := data.(*DisplayCellStats)
		displayName := appStats.StackName
		if displayName == "" {
			displayName = "--"
		}
		return util.FormatDisplayData(displayName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayCellStats)
		return appStats.StackName
	}
	c := uiCommon.NewListColumn("stackName", "STACK", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnDeploymentName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*DisplayCellStats).DeploymentName, c2.(*DisplayCellStats).DeploymentName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*DisplayCellStats)
		return util.FormatDisplayData(cellStats.DeploymentName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*DisplayCellStats)
		return cellStats.DeploymentName
	}
	c := uiCommon.NewListColumn("DNAME", "DNAME", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnJobName() *uiCommon.ListColumn {
	defaultColSize := 45
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*DisplayCellStats).JobName, c2.(*DisplayCellStats).JobName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*DisplayCellStats)
		return util.FormatDisplayData(cellStats.JobName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*DisplayCellStats)
		return cellStats.JobName
	}
	c := uiCommon.NewListColumn("JOB_NAME", "JOB_NAME", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

// Job Index in PCF 1.8 is now a GUID not a integer number
func columnJobIndex() *uiCommon.ListColumn {
	defaultColSize := 36
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayCellStats).JobIndex < c2.(*DisplayCellStats).JobIndex
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*DisplayCellStats)
		return util.FormatDisplayData(cellStats.JobIndex, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*DisplayCellStats)
		return cellStats.JobIndex
	}
	c := uiCommon.NewListColumn("JOB_IDX", "JOB_IDX", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}
