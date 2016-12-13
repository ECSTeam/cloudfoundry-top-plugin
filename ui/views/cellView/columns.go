package cellView

import (
	"fmt"

	"github.com/kkellner/cloudfoundry-top-plugin/eventdata"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

func (asUI *CellListView) columnIp() *uiCommon.ListColumn {
	defaultColSize := 16
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*eventdata.CellStats).Ip, c2.(*eventdata.CellStats).Ip)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*eventdata.CellStats)
		return util.FormatDisplayData(cellStats.Ip, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*eventdata.CellStats)
		return cellStats.Ip
	}
	c := uiCommon.NewListColumn("IP", "IP", defaultColSize,
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnNumOfCpus() *uiCommon.ListColumn {
	defaultColSize := 4
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*eventdata.CellStats).NumOfCpus > c2.(*eventdata.CellStats).NumOfCpus)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*eventdata.CellStats)
		display := ""
		if cellStats.NumOfCpus == 0 {
			display = fmt.Sprintf("%4v", "--")
		} else {
			display = fmt.Sprintf("%4v", cellStats.NumOfCpus)
		}
		return fmt.Sprintf("%4v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*eventdata.CellStats)
		return fmt.Sprintf("%v", cellStats.NumOfCpus)
	}
	c := uiCommon.NewListColumn("CPUS", "CPUS", defaultColSize,
		uiCommon.NUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnCapacityTotalMemory() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.CellStats).CapacityTotalMemory < c2.(*eventdata.CellStats).CapacityTotalMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*eventdata.CellStats)
		display := ""
		if CellStats.CapacityTotalMemory == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.CapacityTotalMemory).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*eventdata.CellStats)
		return fmt.Sprintf("%v", CellStats.CapacityTotalMemory)
	}
	c := uiCommon.NewListColumn("TOT_MEM", "TOT_MEM", 9,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnCapacityRemainingMemory() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.CellStats).CapacityRemainingMemory < c2.(*eventdata.CellStats).CapacityRemainingMemory
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*eventdata.CellStats)
		display := ""
		if CellStats.CapacityRemainingMemory == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.CapacityRemainingMemory).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*eventdata.CellStats)
		return fmt.Sprintf("%v", CellStats.CapacityRemainingMemory)
	}
	c := uiCommon.NewListColumn("FREE_MEM", "FREE_MEM", 9,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnCapacityTotalDisk() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.CellStats).CapacityTotalDisk < c2.(*eventdata.CellStats).CapacityTotalDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*eventdata.CellStats)
		display := ""
		if CellStats.CapacityTotalDisk == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.CapacityTotalDisk).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*eventdata.CellStats)
		return fmt.Sprintf("%v", CellStats.CapacityTotalDisk)
	}
	c := uiCommon.NewListColumn("TOT_DISK", "TOT_DISK", 9,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnCapacityRemainingDisk() *uiCommon.ListColumn {
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*eventdata.CellStats).CapacityRemainingDisk < c2.(*eventdata.CellStats).CapacityRemainingDisk
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		CellStats := data.(*eventdata.CellStats)
		display := ""
		if CellStats.CapacityRemainingDisk == 0 {
			display = fmt.Sprintf("%9v", "--")
		} else {
			display = fmt.Sprintf("%9v", util.ByteSize(CellStats.CapacityRemainingDisk).StringWithPrecision(1))
		}
		return fmt.Sprintf("%9v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		CellStats := data.(*eventdata.CellStats)
		return fmt.Sprintf("%v", CellStats.CapacityRemainingDisk)
	}
	c := uiCommon.NewListColumn("FREE_DISK", "FREE_DISK", 9,
		uiCommon.NUMERIC, false, appNameSortFunc, true, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnCapacityTotalContainers() *uiCommon.ListColumn {
	defaultColSize := 8
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*eventdata.CellStats).CapacityTotalContainers > c2.(*eventdata.CellStats).CapacityTotalContainers)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*eventdata.CellStats)
		display := ""
		if cellStats.CapacityTotalContainers == 0 {
			display = fmt.Sprintf("%8v", "--")
		} else {
			display = fmt.Sprintf("%8v", cellStats.CapacityTotalContainers)
		}
		return fmt.Sprintf("%8v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*eventdata.CellStats)
		return fmt.Sprintf("%v", cellStats.CapacityTotalContainers)
	}
	c := uiCommon.NewListColumn("MAX_CNTR", "MAX_CNTR", defaultColSize,
		uiCommon.NUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnContainerCount() *uiCommon.ListColumn {
	defaultColSize := 5
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*eventdata.CellStats).ContainerCount > c2.(*eventdata.CellStats).ContainerCount)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*eventdata.CellStats)
		display := ""
		if cellStats.ContainerCount == 0 {
			display = fmt.Sprintf("%5v", "--")
		} else {
			display = fmt.Sprintf("%5v", cellStats.ContainerCount)
		}
		return fmt.Sprintf("%5v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*eventdata.CellStats)
		return fmt.Sprintf("%v", cellStats.ContainerCount)
	}
	c := uiCommon.NewListColumn("CNTRS", "CNTRS", defaultColSize,
		uiCommon.NUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnDeploymentName() *uiCommon.ListColumn {
	defaultColSize := 5
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*eventdata.CellStats).DeploymentName, c2.(*eventdata.CellStats).DeploymentName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*eventdata.CellStats)
		return util.FormatDisplayData(cellStats.DeploymentName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*eventdata.CellStats)
		return cellStats.DeploymentName
	}
	c := uiCommon.NewListColumn("DNAME", "DNAME", defaultColSize,
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnJobName() *uiCommon.ListColumn {
	defaultColSize := 20
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*eventdata.CellStats).JobName, c2.(*eventdata.CellStats).JobName)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*eventdata.CellStats)
		return util.FormatDisplayData(cellStats.JobName, defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*eventdata.CellStats)
		return cellStats.JobName
	}
	c := uiCommon.NewListColumn("JOB_NAME", "JOB_NAME", defaultColSize,
		uiCommon.ALPHANUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}

func (asUI *CellListView) columnJobIndex() *uiCommon.ListColumn {
	defaultColSize := 7
	appNameSortFunc := func(c1, c2 util.Sortable) bool {
		return util.CaseInsensitiveLess(c1.(*eventdata.CellStats).JobIndex, c2.(*eventdata.CellStats).JobIndex)
	}
	displayFunc := func(data uiCommon.IData, isSelected bool) string {
		cellStats := data.(*eventdata.CellStats)
		display := fmt.Sprintf("%7v", cellStats.JobIndex)
		return display
	}
	rawValueFunc := func(data uiCommon.IData) string {
		cellStats := data.(*eventdata.CellStats)
		return cellStats.JobIndex
	}
	c := uiCommon.NewListColumn("JOB_IDX", "JOB_IDX", defaultColSize,
		uiCommon.NUMERIC, true, appNameSortFunc, false, displayFunc, rawValueFunc)
	return c
}
