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
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func (asUI *CellListView) columnIp() *uiCommon.ListColumn {
	defaultColSize := 16
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.Ip2long(c1.(*displaydata.DisplayCellStats).Ip) < util.Ip2long(c2.(*displaydata.DisplayCellStats).Ip)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return util.FormatDisplayData(cellStats.Ip, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return cellStats.Ip
	}
	c := uiCommon.NewListColumn("IP", "IP", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnNumOfCpus() *uiCommon.ListColumn {
	defaultColSize := 4
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).NumOfCpus < c2.(*displaydata.DisplayCellStats).NumOfCpus)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if cellStats.NumOfCpus == 0 {
			display = fmt.Sprintf("%4v", "--")
		} else {
			display = fmt.Sprintf("%4v", cellStats.NumOfCpus)
		}
		return fmt.Sprintf("%4v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.NumOfCpus)
	}
	c := uiCommon.NewListColumn("CPUS", "CPUS", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnCapacityTotalMemory() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayCellStats).CapacityTotalMemory < c2.(*displaydata.DisplayCellStats).CapacityTotalMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if CellStats.CapacityTotalMemory == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.CapacityTotalMemory).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.CapacityTotalMemory)
	}
	c := uiCommon.NewListColumn("TOT_MEM", "TOT_MEM", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnCapacityRemainingMemory() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayCellStats).CapacityRemainingMemory < c2.(*displaydata.DisplayCellStats).CapacityRemainingMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if CellStats.CapacityRemainingMemory == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.CapacityRemainingMemory).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.CapacityRemainingMemory)
	}
	c := uiCommon.NewListColumn("FREE_MEM", "FREE_MEM", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnCapacityTotalDisk() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayCellStats).CapacityTotalDisk < c2.(*displaydata.DisplayCellStats).CapacityTotalDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if CellStats.CapacityTotalDisk == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.CapacityTotalDisk).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.CapacityTotalDisk)
	}
	c := uiCommon.NewListColumn("TOT_DISK", "TOT_DISK", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnCapacityRemainingDisk() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayCellStats).CapacityRemainingDisk < c2.(*displaydata.DisplayCellStats).CapacityRemainingDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if CellStats.CapacityRemainingDisk == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.CapacityRemainingDisk).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.CapacityRemainingDisk)
	}
	c := uiCommon.NewListColumn("FREE_DISK", "FREE_DISK", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnCapacityTotalContainers() *uiCommon.ListColumn {
	defaultColSize := 8
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).CapacityTotalContainers < c2.(*displaydata.DisplayCellStats).CapacityTotalContainers)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if cellStats.CapacityTotalContainers == 0 {
			display = fmt.Sprintf("%8v", "--")
		} else {
			display = fmt.Sprintf("%8v", cellStats.CapacityTotalContainers)
		}
		return fmt.Sprintf("%8v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityTotalContainers)
	}
	c := uiCommon.NewListColumn("MAX_CNTR", "MAX_CNTR", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnContainerCount() *uiCommon.ListColumn {
	defaultColSize := 5
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).ContainerCount < c2.(*displaydata.DisplayCellStats).ContainerCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if cellStats.ContainerCount == 0 {
			display = fmt.Sprintf("%5v", "--")
		} else {
			display = fmt.Sprintf("%5v", cellStats.ContainerCount)
		}
		return fmt.Sprintf("%5v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.ContainerCount)
	}
	c := uiCommon.NewListColumn("CNTRS", "CNTRS", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnDeploymentName() *uiCommon.ListColumn {
	defaultColSize := 10
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*displaydata.DisplayCellStats).DeploymentName, c2.(*displaydata.DisplayCellStats).DeploymentName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return util.FormatDisplayData(cellStats.DeploymentName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return cellStats.DeploymentName
	}
	c := uiCommon.NewListColumn("DNAME", "DNAME", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnJobName() *uiCommon.ListColumn {
	defaultColSize := 45
	sortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*displaydata.DisplayCellStats).JobName, c2.(*displaydata.DisplayCellStats).JobName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return util.FormatDisplayData(cellStats.JobName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return cellStats.JobName
	}
	c := uiCommon.NewListColumn("JOB_NAME", "JOB_NAME", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnJobIndex() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).JobIndex < c2.(*displaydata.DisplayCellStats).JobIndex)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		display := fmt.Sprintf("%7v", cellStats.JobIndex)
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.JobIndex)
	}
	c := uiCommon.NewListColumn("JOB_IDX", "JOB_IDX", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnTotalCpuPercentage() *uiCommon.ListColumn {
	defaultColSize := 6
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).TotalContainerCpuPercentage < c2.(*displaydata.DisplayCellStats).TotalContainerCpuPercentage)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		totalCpuInfo := ""
		if cellStats.TotalReportingContainers == 0 {
			totalCpuInfo = fmt.Sprintf("%6v", "--")
		} else {
			if cellStats.TotalContainerCpuPercentage >= 100.0 {
				totalCpuInfo = fmt.Sprintf("%6.0f", cellStats.TotalContainerCpuPercentage)
			} else if cellStats.TotalContainerCpuPercentage >= 10.0 {
				totalCpuInfo = fmt.Sprintf("%6.1f", cellStats.TotalContainerCpuPercentage)
			} else {
				totalCpuInfo = fmt.Sprintf("%6.2f", cellStats.TotalContainerCpuPercentage)
			}
		}
		return fmt.Sprintf("%6v", totalCpuInfo)

	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.TotalContainerCpuPercentage)
	}
	c := uiCommon.NewListColumn("CPU_PERCENT", "CPU%", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)

	return c
}

func (asUI *CellListView) columnTotalReportingContainers() *uiCommon.ListColumn {
	defaultColSize := 4
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).TotalReportingContainers < c2.(*displaydata.DisplayCellStats).TotalReportingContainers)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if cellStats.TotalReportingContainers == 0 {
			display = fmt.Sprintf("%4v", "--")
		} else {
			display = fmt.Sprintf("%4v", cellStats.TotalReportingContainers)
		}
		return fmt.Sprintf("%4v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.TotalReportingContainers)
	}
	c := uiCommon.NewListColumn("RCR", "RCR", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnTotalContainerReservedMemory() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayCellStats).TotalContainerReservedMemory < c2.(*displaydata.DisplayCellStats).TotalContainerReservedMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if CellStats.TotalContainerReservedMemory == 0 {
			display = fmt.Sprintf("%10v", "--")
		} else {
			display = fmt.Sprintf("%10v", util.ByteSize(CellStats.TotalContainerReservedMemory).StringWithPrecision(1))
		}
		return fmt.Sprintf("%10v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.TotalContainerReservedMemory)
	}
	c := uiCommon.NewListColumn("C_RSVD_MEM", "C_RSVD_MEM", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnTotalContainerUsedMemory() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayCellStats).TotalContainerUsedMemory < c2.(*displaydata.DisplayCellStats).TotalContainerUsedMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if CellStats.TotalContainerUsedMemory == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.TotalContainerUsedMemory).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.TotalContainerUsedMemory)
	}
	c := uiCommon.NewListColumn("C_USD_MEM", "C_USD_MEM", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnTotalContainerReservedDisk() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayCellStats).TotalContainerReservedDisk < c2.(*displaydata.DisplayCellStats).TotalContainerReservedDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if CellStats.TotalContainerReservedDisk == 0 {
			display = fmt.Sprintf("%10v", "--")
		} else {
			display = fmt.Sprintf("%10v", util.ByteSize(CellStats.TotalContainerReservedDisk).StringWithPrecision(1))
		}
		return fmt.Sprintf("%10v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.TotalContainerReservedDisk)
	}
	c := uiCommon.NewListColumn("C_RSVD_DSK", "C_RSVD_DSK", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnTotalContainerUsedDisk() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*displaydata.DisplayCellStats).TotalContainerUsedDisk < c2.(*displaydata.DisplayCellStats).TotalContainerUsedDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if CellStats.TotalContainerUsedDisk == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.TotalContainerUsedDisk).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.TotalContainerUsedDisk)
	}
	c := uiCommon.NewListColumn("C_USD_DSK", "C_USD_DSK", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc)
	return c
}
