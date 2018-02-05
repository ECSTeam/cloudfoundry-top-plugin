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

package headerView

import (
	"fmt"
	"sort"

	"strings"

	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/isolationSegment"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

const UNKNOWN_STACK_NAME = "UNKNOWN"
const UNKNOWN_ISOSEG_NAME = "UNKNOWN"

type StackSummaryStats struct {
	IsolationSegmentGuid        string
	IsolationSegmentName        string
	StackId                     string
	StackName                   string
	TotalCells                  int
	TotalApps                   int
	TotalReportingAppInstances  int
	TotalActiveApps             int
	TotalMemoryUsedAppInstances int64
	TotalDiskUsedAppInstances   int64
	TotalCpuPercentage          float64
	TotalCellCPUs               int
	// This is the hightest CPU percent of all the cells
	CellMaxCpuPercentage float64
	// This is the percent of CPU capacity of the hightest CPU percent cell
	// E.g.,  If CellMaxCpuPercentage is 160.0 and our cell has 2 vCPUs
	// then CellMaxCpuCapacity = (160 * 100) / (2 * 100) = 80 percent
	// To simplify the above equation it would be (160 / 2) = 80 percent
	CellMaxCpuCapacity  float64
	TotalCapacityMemory int64
	TotalCapacityDisk   int64
	ReservedMem         float64
	ReservedDisk        float64
}

type StackSummaryStatsArray []*StackSummaryStats

func (slice StackSummaryStatsArray) Len() int {
	return len(slice)
}

func (slice StackSummaryStatsArray) Less(i, j int) bool {

	isoSegName1 := slice[i].IsolationSegmentName
	isoSegName2 := slice[j].IsolationSegmentName

	if isoSegName1 == isoSegName2 {
		// Always sort UNKNOWN stack to bottom
		stackName1 := slice[i].StackName
		if strings.HasPrefix(stackName1, UNKNOWN_STACK_NAME) {
			return false
		}
		stackName2 := slice[j].StackName
		if strings.HasPrefix(stackName2, UNKNOWN_STACK_NAME) {
			return true
		}
		return stackName1 < stackName2
	}

	// Always sort "shared" isolation segment to top
	if strings.HasPrefix(isoSegName1, isolationSegment.SharedIsolationSegmentName) {
		return true
	}
	if strings.HasPrefix(isoSegName2, isolationSegment.SharedIsolationSegmentName) {
		return false
	}
	// Always sort UNKNOWN isolation segment to bottom
	if strings.HasPrefix(isoSegName1, UNKNOWN_ISOSEG_NAME) {
		return false
	}
	if strings.HasPrefix(isoSegName2, UNKNOWN_ISOSEG_NAME) {
		return true
	}

	return isoSegName1 < isoSegName2

}

func (slice StackSummaryStatsArray) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Output header stats by isolation-segment & stack
// Returns the number of rows (lines) written to header
func (asUI *HeaderWidget) updateHeaderStack(g *gocui.Gui, v *gocui.View) (int, error) {

	isWarmupComplete := asUI.masterUI.IsWarmupComplete()

	router := asUI.router
	processor := router.GetProcessor()
	mdMgr := processor.GetMetadataManager()
	stacks := mdMgr.GetStackMdManager().GetAll()
	if len(stacks) == 0 {
		fmt.Fprintf(v, "\n Waiting for more data...")
		return 3, nil
	}

	// Keys: [IsoSegGuid] [StackGuid]
	summaryStatsByIsoSeg := make(map[string]map[string]*StackSummaryStats)
	//summaryStatsByStack := make(map[string]*StackSummaryStats)

	isolationSegments := mdMgr.GetIsoSegMdManager().GetAll()

	if len(isolationSegments) == 0 {
		// We must be running against a Cloud Foundry foundation prior to isolation segment support
		isolationSegments = append(isolationSegments, isolationSegment.GetDefault())
	} else {
		isolationSegments = append(isolationSegments, isolationSegment.UnknownIsolationSegment)
		//isolationSegments = append(isolationSegments, isolationSegment.SharedIsolationSegment)
	}

	for _, isoSeg := range isolationSegments {
		summaryStatsByIsoSeg[isoSeg.Guid] = make(map[string]*StackSummaryStats)
		for _, stack := range stacks {
			//toplog.Info("isoSeg.Guid: %v  stack.Guid: %v", isoSeg.Guid, stack.Guid)
			summaryStatsByIsoSeg[isoSeg.Guid][stack.Guid] = &StackSummaryStats{
				IsolationSegmentGuid: isoSeg.Guid,
				IsolationSegmentName: isoSeg.Name,
				StackId:              stack.Guid,
				StackName:            stack.Name}
		}
	}

	if isWarmupComplete {
		// We add an extra StackSummaryStats with no isolation-segment and stackId to handle cells that have no containers (yet)
		if summaryStatsByIsoSeg[""] == nil {
			summaryStatsByIsoSeg[""] = make(map[string]*StackSummaryStats)
		}
		summaryStatsByIsoSeg[""][""] = &StackSummaryStats{IsolationSegmentName: UNKNOWN_ISOSEG_NAME, StackId: "", StackName: UNKNOWN_STACK_NAME}
	}

	//toplog.Info("asUI.commonData.GetDisplayAppStatsMap len: %v", len(asUI.commonData.GetDisplayAppStatsMap()))

	// Key: cellIP
	cpuByCellMap := make(map[string]float64)

	for _, appStats := range asUI.commonData.GetDisplayAppStatsMap() {
		isolationSegGuid := appStats.IsolationSegmentGuid

		if isolationSegGuid == isolationSegment.DefaultIsolationSegmentGuid && isolationSegment.SharedIsolationSegment != nil {
			isolationSegGuid = isolationSegment.SharedIsolationSegment.Guid
		}

		//toplog.Info("*** isolationSegGuid: %v", isolationSegGuid)

		sumStats := summaryStatsByIsoSeg[isolationSegGuid][appStats.StackId]
		if appStats.StackId == "" || sumStats == nil {
			// This appStats has no stackId -- This could be caused by not having
			// the app metadata in cache.  Either because it hasn't been loaded yet
			// or it was deleted because the app has been deleted but we still
			// have stats from when the app was deployed.

			//log.Panic(fmt.Sprintf("We didn't find the stack data, StackId: %v", appStats.StackId))
			//fmt.Fprintf(v, "\n Waiting for more data...")
			//return 3, nil
			//toplog.Info("******** appStats.StackId: %v sumStats:%v", appStats.StackId, sumStats)
			continue
		}

		// Track how much CPU is consumed per cell
		for _, containerStats := range appStats.ContainerArray {
			if containerStats != nil {
				cellIP := containerStats.Ip
				cpuPercent := containerStats.ContainerMetric.GetCpuPercentage()
				cpuByCellMap[cellIP] = cpuByCellMap[cellIP] + cpuPercent
			}
		}

		sumStats.TotalReportingAppInstances = sumStats.TotalReportingAppInstances + appStats.TotalReportingContainers
		sumStats.TotalMemoryUsedAppInstances = sumStats.TotalMemoryUsedAppInstances + appStats.TotalMemoryUsed
		sumStats.TotalDiskUsedAppInstances = sumStats.TotalDiskUsedAppInstances + appStats.TotalDiskUsed
		sumStats.TotalCpuPercentage = sumStats.TotalCpuPercentage + appStats.TotalCpuPercentage
		if appStats.TotalTraffic.EventL60Rate > 0 {
			sumStats.TotalActiveApps++
		}
		sumStats.TotalApps++
	}

	appMdMgr := mdMgr.GetAppMdManager()
	for _, app := range appMdMgr.AllApps() {
		spaceMetadata := mdMgr.GetSpaceMdManager().FindItem(app.SpaceGuid)
		isolationSegGuid := spaceMetadata.IsolationSegmentGuid
		if isolationSegGuid == isolationSegment.DefaultIsolationSegmentGuid && isolationSegment.SharedIsolationSegment != nil {
			isolationSegGuid = isolationSegment.SharedIsolationSegment.Guid
		}
		sumStats := summaryStatsByIsoSeg[isolationSegGuid][app.StackGuid]
		if sumStats != nil {
			if app.State == "STARTED" {
				sumStats.ReservedMem = sumStats.ReservedMem + ((app.MemoryMB * util.MEGABYTE) * app.Instances)
				sumStats.ReservedDisk = sumStats.ReservedDisk + ((app.DiskQuotaMB * util.MEGABYTE) * app.Instances)
			}
		}
	}

	for _, cellStats := range processor.GetDisplayedEventData().CellMap {
		//toplog.Info("cellStats.StackId:%v", cellStats.StackId)
		isolationSegGuid := cellStats.IsolationSegmentGuid
		if isolationSegGuid == isolationSegment.DefaultIsolationSegmentGuid && isolationSegment.SharedIsolationSegment != nil {
			isolationSegGuid = isolationSegment.SharedIsolationSegment.Guid
		} else if isolationSegGuid == isolationSegment.UnknownIsolationSegmentGuid {
			// Do another attempt to resolve unknown IsoSeg before we display header
			processor.GetDisplayedEventData().AssignIsolationSegment(cellStats)
			isolationSegGuid = cellStats.IsolationSegmentGuid
		}
		sumStats := summaryStatsByIsoSeg[isolationSegGuid][cellStats.StackId]
		// We might get nil sumStats if we are still in the warm-up period and stackId is unknown yet
		if sumStats != nil {
			sumStats.TotalCells = sumStats.TotalCells + 1
			sumStats.TotalCellCPUs = sumStats.TotalCellCPUs + cellStats.NumOfCpus
			sumStats.TotalCapacityMemory = sumStats.TotalCapacityMemory + cellStats.CapacityMemoryTotal
			sumStats.TotalCapacityDisk = sumStats.TotalCapacityDisk + cellStats.CapacityDiskTotal

			cellCpu := cpuByCellMap[cellStats.Ip]
			if cellCpu > sumStats.CellMaxCpuPercentage {
				sumStats.CellMaxCpuPercentage = cpuByCellMap[cellStats.Ip]
				if cellStats.NumOfCpus > 0 {
					// Only calc the capacity if we know the total CPU count
					sumStats.CellMaxCpuCapacity = sumStats.CellMaxCpuPercentage / float64(cellStats.NumOfCpus)
				}
			}
		}
	}

	// Output stack information by stack name sort order
	stackSummaryStatsArray := make(StackSummaryStatsArray, 0)
	for _, stackSummaryStatsMap := range summaryStatsByIsoSeg {
		for _, stackSummaryStats := range stackSummaryStatsMap {
			stackSummaryStatsArray = append(stackSummaryStatsArray, stackSummaryStats)
		}
	}
	//toplog.Info("stackSummaryStatsArray:%v", len(stackSummaryStatsArray))
	sort.Sort(stackSummaryStatsArray)

	// Set check for greater then two because we force an 'unknown' segment.  If we only have 'shared'
	// as a real segment, there is no need to show segment
	showIsolationSegment := len(isolationSegments) > 2
	linesWritten := 0

	for _, stackSummaryStats := range stackSummaryStatsArray {

		//toplog.Info("stackSummaryStats name: %v TotalApps: %v", stackSummaryStats.StackName, stackSummaryStats.TotalApps)

		if stackSummaryStats.TotalApps > 0 || stackSummaryStats.TotalCellCPUs > 0 {
			linesWritten += asUI.outputHeaderForStack(g, v, stackSummaryStats, showIsolationSegment)
		}
	}

	if linesWritten == 0 {
		// Likely the "stacks" metadata loaded but not the "apps" metadata
		fmt.Fprintf(v, "\n Waiting for even more data...")
		return 3, nil
	}

	return linesWritten, nil
}

// Called for each isolation-segment and stack combination - this should output 3 lines:
//
//  IsoSeg: shared   Stack: cflinuxfs2    Cells: 5
//     CPU:  8.4% Used,  800% Max,       Mem:   7GB Used,  63GB Max,  22GB Rsrvd
//     Apps:  122 Total, Cntrs:  127     Dsk:   7GB Used, 190GB Max,  27GB Rsrvd
//
func (asUI *HeaderWidget) outputHeaderForStack(g *gocui.Gui, v *gocui.View, stackSummaryStats *StackSummaryStats, showIsolationSegment bool) int {

	TotalMemoryUsedAppInstancesDisplay := "--"
	TotalDiskUsedAppInstancesDisplay := "--"
	totalCpuPercentageDisplay := "--"
	if stackSummaryStats.TotalReportingAppInstances > 0 {
		TotalMemoryUsedAppInstancesDisplay = util.ByteSize(stackSummaryStats.TotalMemoryUsedAppInstances).StringWithPrecision(0)
		TotalDiskUsedAppInstancesDisplay = util.ByteSize(stackSummaryStats.TotalDiskUsedAppInstances).StringWithPrecision(0)
		if stackSummaryStats.TotalCpuPercentage >= 100 {
			totalCpuPercentageDisplay = fmt.Sprintf("%.0f%%", stackSummaryStats.TotalCpuPercentage)
		} else {
			totalCpuPercentageDisplay = fmt.Sprintf("%.1f%%", stackSummaryStats.TotalCpuPercentage)
		}

		isWarmupComplete := asUI.masterUI.IsWarmupComplete()
		if isWarmupComplete {
			switch {
			case stackSummaryStats.CellMaxCpuCapacity >= 90:
				colorString := util.BRIGHT_RED
				totalCpuPercentageDisplay = fmt.Sprintf("%v%7v%v", colorString, totalCpuPercentageDisplay, util.CLEAR)
			case stackSummaryStats.CellMaxCpuCapacity >= 80:
				colorString := util.BRIGHT_YELLOW
				totalCpuPercentageDisplay = fmt.Sprintf("%v%7v%v", colorString, totalCpuPercentageDisplay, util.CLEAR)
			}
		}
	}

	cellTotalCPUCapacityDisplay := "--"
	if stackSummaryStats.TotalCellCPUs > 0 {
		cellTotalCPUCapacityDisplay = fmt.Sprintf("%v%%", (stackSummaryStats.TotalCellCPUs * 100))
	}

	CapacityMemoryTotalDisplay := "--"
	if stackSummaryStats.TotalCapacityMemory > 0 {
		CapacityMemoryTotalDisplay = fmt.Sprintf("%v", util.ByteSize(stackSummaryStats.TotalCapacityMemory).StringWithPrecision(0))
	}
	CapacityDiskTotalDisplay := "--"
	if stackSummaryStats.TotalCapacityDisk > 0 {
		CapacityDiskTotalDisplay = fmt.Sprintf("%v", util.ByteSize(stackSummaryStats.TotalCapacityDisk).StringWithPrecision(0))
	}
	TotalCellsDisplay := "--"
	if stackSummaryStats.TotalCells > 0 {
		TotalCellsDisplay = fmt.Sprintf("%v", stackSummaryStats.TotalCells)
	}

	if showIsolationSegment {
		fmt.Fprintf(v, "IsoSeg: %-12v ", stackSummaryStats.IsolationSegmentName)
	}
	notes := ""
	if stackSummaryStats.StackName == UNKNOWN_STACK_NAME {
		notes = "(cells with no containers)"
	}
	fmt.Fprintf(v, "Stack: %-13v Cells: %-3v%v\n", stackSummaryStats.StackName, TotalCellsDisplay, notes)
	fmt.Fprintf(v, "   CPU:%7v Used,%7v Max,     ", totalCpuPercentageDisplay, cellTotalCPUCapacityDisplay)

	displayTotalMem := "--"
	totalMem := stackSummaryStats.ReservedMem
	if totalMem > 0 {
		displayTotalMem = util.ByteSize(totalMem).StringWithPrecision(0)
	}
	fmt.Fprintf(v, "Mem:%6v Used,", TotalMemoryUsedAppInstancesDisplay)
	// Total quota memory of all running app instances
	fmt.Fprintf(v, "%6v Max,%6v Rsrvd\n", CapacityMemoryTotalDisplay, displayTotalMem)

	// Reporting containers are containers that reported metrics in last 'StaleContainerSeconds'
	fmt.Fprintf(v, "   Apps:%5v Total, Cntrs:%5v     ",
		stackSummaryStats.TotalApps,
		stackSummaryStats.TotalReportingAppInstances)

	displayTotalDisk := "--"
	totalDisk := stackSummaryStats.ReservedDisk
	if totalMem > 0 {
		displayTotalDisk = util.ByteSize(totalDisk).StringWithPrecision(0)
	}

	fmt.Fprintf(v, "Dsk:%6v Used,", TotalDiskUsedAppInstancesDisplay)
	fmt.Fprintf(v, "%6v Max,%6v Rsrvd\n", CapacityDiskTotalDisplay, displayTotalDisk)

	// Number of lines written
	return 3

}
