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
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/views/cellViews/cellView"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func columnIp() *uiCommon.ListColumn {
	defaultColSize := 16
	sortFunc := func(c1, c2 util.Sortable) bool {
		ip1 := c1.(*cellView.DisplayCellStats).Ip
		// Always sort TOTAL to the top
		if ip1 == "TOTAL" {
			return true
		}
		ip2 := c2.(*cellView.DisplayCellStats).Ip
		if ip2 == "TOTAL" {
			return false
		}
		return util.Ip2long(ip1) < util.Ip2long(ip2)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		cellStats := data.(*cellView.DisplayCellStats)
		return util.FormatDisplayData(cellStats.Ip, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*cellView.DisplayCellStats)
		return cellStats.Ip
	}
	c := uiCommon.NewListColumn("CELL_IP", "CELL_IP", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnNumOfCpus() *uiCommon.ListColumn {
	defaultColSize := 4
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*cellView.DisplayCellStats).NumOfCpus < c2.(*cellView.DisplayCellStats).NumOfCpus)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		cellStats := data.(*cellView.DisplayCellStats)
		display := ""
		if cellStats.NumOfCpus == 0 {
			display = fmt.Sprintf("%4v", "--")
		} else {
			display = fmt.Sprintf("%4v", cellStats.NumOfCpus)
		}
		return fmt.Sprintf("%4v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*cellView.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.NumOfCpus)
	}
	c := uiCommon.NewListColumn("CPUS", "CPUS", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}

func columnCapacityMemoryTotal() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*cellView.DisplayCellStats).CapacityMemoryTotal < c2.(*cellView.DisplayCellStats).CapacityMemoryTotal
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		CellStats := data.(*cellView.DisplayCellStats)
		display := ""
		if CellStats.CapacityMemoryTotal == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.CapacityMemoryTotal).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*cellView.DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.CapacityMemoryTotal)
	}
	c := uiCommon.NewListColumn("MEM_TOT", "MEM_TOT", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnCapacityMemoryRemaining() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*cellView.DisplayCellStats).CapacityMemoryRemaining < c2.(*cellView.DisplayCellStats).CapacityMemoryRemaining
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		CellStats := data.(*cellView.DisplayCellStats)
		display := ""
		if CellStats.CapacityMemoryRemaining == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.CapacityMemoryRemaining).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*cellView.DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.CapacityMemoryRemaining)
	}
	c := uiCommon.NewListColumn("MEM_FREE", "MEM_FREE", 9,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnCapacityTotalContainers() *uiCommon.ListColumn {
	defaultColSize := 8
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*cellView.DisplayCellStats).CapacityTotalContainers < c2.(*cellView.DisplayCellStats).CapacityTotalContainers)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		cellStats := data.(*cellView.DisplayCellStats)
		display := ""
		if cellStats.CapacityTotalContainers == 0 {
			display = fmt.Sprintf("%8v", "--")
		} else {
			display = fmt.Sprintf("%8v", cellStats.CapacityTotalContainers)
		}
		return fmt.Sprintf("%8v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*cellView.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityTotalContainers)
	}
	c := uiCommon.NewListColumn("MAX_CNTR", "MAX_CNTR", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnContainerCount() *uiCommon.ListColumn {
	defaultColSize := 5
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*cellView.DisplayCellStats).ContainerCount < c2.(*cellView.DisplayCellStats).ContainerCount)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		cellStats := data.(*cellView.DisplayCellStats)
		display := ""
		if cellStats.ContainerCount == 0 {
			display = fmt.Sprintf("%5v", "--")
		} else {
			display = fmt.Sprintf("%5v", cellStats.ContainerCount)
		}
		return fmt.Sprintf("%5v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*cellView.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.ContainerCount)
	}
	c := uiCommon.NewListColumn("CNTRS", "CNTRS", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnTotalContainerMemoryReserved() *uiCommon.ListColumn {
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*cellView.DisplayCellStats).TotalContainerMemoryReserved < c2.(*cellView.DisplayCellStats).TotalContainerMemoryReserved
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		CellStats := data.(*cellView.DisplayCellStats)
		display := ""
		if CellStats.TotalContainerMemoryReserved == 0 {
			display = fmt.Sprintf("%10v", "--")
		} else {
			display = fmt.Sprintf("%10v", util.ByteSize(CellStats.TotalContainerMemoryReserved).StringWithPrecision(1))
		}
		return fmt.Sprintf("%10v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*cellView.DisplayCellStats)
		return fmt.Sprintf("%v", CellStats.TotalContainerMemoryReserved)
	}
	c := uiCommon.NewListColumn("C_MEM_RSVD", "C_MEM_RSVD", 10,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnCapacityPlan0_5GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*cellView.DisplayCellStats).CapacityPlan0_5GMem < c2.(*cellView.DisplayCellStats).CapacityPlan0_5GMem)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		cellStats := data.(*cellView.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan0_5GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan0_5GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*cellView.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan0_5GMem)
	}
	c := uiCommon.NewListColumn("0.5GB", "  0.5GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnCapacityPlan1_0GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*cellView.DisplayCellStats).CapacityPlan1_0GMem < c2.(*cellView.DisplayCellStats).CapacityPlan1_0GMem)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		cellStats := data.(*cellView.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan1_0GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan1_0GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*cellView.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan1_0GMem)
	}
	c := uiCommon.NewListColumn("1.0GB", "  1.0GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnCapacityPlan1_5GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*cellView.DisplayCellStats).CapacityPlan1_5GMem < c2.(*cellView.DisplayCellStats).CapacityPlan1_5GMem)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		cellStats := data.(*cellView.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan1_5GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan1_5GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*cellView.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan1_5GMem)
	}
	c := uiCommon.NewListColumn("1.5GB", "  1.5GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnCapacityPlan2_0GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*cellView.DisplayCellStats).CapacityPlan2_0GMem < c2.(*cellView.DisplayCellStats).CapacityPlan2_0GMem)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		cellStats := data.(*cellView.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan2_0GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan2_0GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*cellView.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan2_0GMem)
	}
	c := uiCommon.NewListColumn("2.0GB", "  2.0GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}
func columnCapacityPlan2_5GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*cellView.DisplayCellStats).CapacityPlan2_5GMem < c2.(*cellView.DisplayCellStats).CapacityPlan2_5GMem)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		cellStats := data.(*cellView.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan2_5GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan2_5GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*cellView.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan2_5GMem)
	}
	c := uiCommon.NewListColumn("2.5GB", "  2.5GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnCapacityPlan3_0GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*cellView.DisplayCellStats).CapacityPlan3_0GMem < c2.(*cellView.DisplayCellStats).CapacityPlan3_0GMem)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		cellStats := data.(*cellView.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan3_0GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan3_0GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*cellView.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan3_0GMem)
	}
	c := uiCommon.NewListColumn("3.0GB", "  3.0GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnCapacityPlan3_5GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*cellView.DisplayCellStats).CapacityPlan3_5GMem < c2.(*cellView.DisplayCellStats).CapacityPlan3_5GMem)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		cellStats := data.(*cellView.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan3_5GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan3_5GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*cellView.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan3_5GMem)
	}
	c := uiCommon.NewListColumn("3.5GB", "  3.5GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnCapacityPlan4_0GMem() *uiCommon.ListColumn {
	defaultColSize := 7
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*cellView.DisplayCellStats).CapacityPlan4_0GMem < c2.(*cellView.DisplayCellStats).CapacityPlan4_0GMem)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		cellStats := data.(*cellView.DisplayCellStats)
		display := ""
		if cellStats.CapacityPlan4_0GMem == UNKNOWN {
			display = fmt.Sprintf("%7v", "--")
		} else {
			display = fmt.Sprintf("%7v", cellStats.CapacityPlan4_0GMem)
		}
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*cellView.DisplayCellStats)
		return fmt.Sprintf("%v", cellStats.CapacityPlan4_0GMem)
	}
	c := uiCommon.NewListColumn("4.0GB", "  4.0GB", defaultColSize,
		uiCommon.NUMERIC, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}
