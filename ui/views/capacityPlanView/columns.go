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

package capacityPlanView

import (
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/displaydata"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func columnIp() *uiCommon.ListColumn {
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

func columnNumOfCpus() *uiCommon.ListColumn {
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

func columnCapacityTotalMemory() *uiCommon.ListColumn {
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

func columnCapacityRemainingMemory() *uiCommon.ListColumn {
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

func columnCapacityTotalContainers() *uiCommon.ListColumn {
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

func columnContainerCount() *uiCommon.ListColumn {
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

func columnTotalContainerReservedMemory() *uiCommon.ListColumn {
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

func columnCapacityPlan0_5GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).CapacityPlan0_5GMem < c2.(*displaydata.DisplayCellStats).CapacityPlan0_5GMem)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan0_5GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan0_5GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan0_5GMem)
	}
	c := uiCommon.NewListColumn("0.5GB", "  0.5GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func columnCapacityPlan1_0GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).CapacityPlan1_0GMem < c2.(*displaydata.DisplayCellStats).CapacityPlan1_0GMem)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan1_0GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan1_0GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan1_0GMem)
	}
	c := uiCommon.NewListColumn("1.0GB", "  1.0GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func columnCapacityPlan1_5GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).CapacityPlan1_5GMem < c2.(*displaydata.DisplayCellStats).CapacityPlan1_5GMem)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan1_5GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan1_5GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan1_5GMem)
	}
	c := uiCommon.NewListColumn("1.5GB", "  1.5GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func columnCapacityPlan2_0GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).CapacityPlan2_0GMem < c2.(*displaydata.DisplayCellStats).CapacityPlan2_0GMem)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan2_0GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan2_0GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan2_0GMem)
	}
	c := uiCommon.NewListColumn("2.0GB", "  2.0GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc)
	return c
}
func columnCapacityPlan2_5GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).CapacityPlan2_5GMem < c2.(*displaydata.DisplayCellStats).CapacityPlan2_5GMem)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan2_5GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan2_5GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan2_5GMem)
	}
	c := uiCommon.NewListColumn("2.5GB", "  2.5GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func columnCapacityPlan3_0GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).CapacityPlan3_0GMem < c2.(*displaydata.DisplayCellStats).CapacityPlan3_0GMem)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan3_0GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan3_0GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan3_0GMem)
	}
	c := uiCommon.NewListColumn("3.0GB", "  3.0GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func columnCapacityPlan3_5GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).CapacityPlan3_5GMem < c2.(*displaydata.DisplayCellStats).CapacityPlan3_5GMem)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan3_5GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan3_5GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan3_5GMem)
	}
	c := uiCommon.NewListColumn("3.5GB", "  3.5GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc)
	return c
}

func columnCapacityPlan4_0GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*displaydata.DisplayCellStats).CapacityPlan4_0GMem < c2.(*displaydata.DisplayCellStats).CapacityPlan4_0GMem)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan4_0GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan4_0GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*displaydata.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan4_0GMem)
	}
	c := uiCommon.NewListColumn("4.0GB", "  4.0GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc)
	return c
}
