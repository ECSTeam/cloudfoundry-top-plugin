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

package appView

import (
	"fmt"
	"sort"

	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/stack"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

type StackSummaryStats struct {
	StackId                     string
	StackName                   string
	TotalApps                   int
	TotalReportingAppInstances  int
	TotalActiveApps             int
	TotalUsedMemoryAppInstances uint64
	TotalUsedDiskAppInstances   uint64
	TotalCpuPercentage          float64
	TotalCellCPUs               int
	TotalCapacityMemory         int64
	TotalCapacityDisk           int64
	ReservedMem                 float64
	ReservedDisk                float64
}

type StackSummaryStatsArray []*StackSummaryStats

func (slice StackSummaryStatsArray) Len() int {
	return len(slice)
}

func (slice StackSummaryStatsArray) Less(i, j int) bool {
	return slice[i].StackName < slice[j].StackName
}

func (slice StackSummaryStatsArray) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Output header stats by stack
// Returns the number of rows (lines) written to header

func (asUI *AppListView) updateHeader(g *gocui.Gui, v *gocui.View) (int, error) {

	// TODO: Is this the best spot to check for alerts?? Seems out of place in the updateHeader method
	asUI.checkForAlerts(g)

	stacks := stack.AllStacks()
	if len(stacks) == 0 {
		fmt.Fprintf(v, "\n Waiting for more data...")
		return 3, nil
	}
	summaryStatsByStack := make(map[string]*StackSummaryStats)
	for _, stack := range stacks {
		summaryStatsByStack[stack.Guid] = &StackSummaryStats{StackId: stack.Guid, StackName: stack.Name}
	}

	for _, appStats := range asUI.displayAppStats {
		sumStats := summaryStatsByStack[appStats.StackId]
		if sumStats == nil {
			// This appStats has no stackId -- This could be caused by not having
			// the app metadata in cache.  Either because it hasn't been loaded yet
			// or it was deleted because the app has been deleted but we still
			// have stats from when the app was deployed.

			//log.Panic(fmt.Sprintf("We didn't find the stack data, StackId: %v", appStats.StackId))
			//fmt.Fprintf(v, "\n Waiting for more data...")
			//return 3, nil
			continue
		}
		for _, cs := range appStats.ContainerArray {
			if cs != nil && cs.ContainerMetric != nil {
				sumStats.TotalReportingAppInstances++
				sumStats.TotalUsedMemoryAppInstances = sumStats.TotalUsedMemoryAppInstances + *cs.ContainerMetric.MemoryBytes
				sumStats.TotalUsedDiskAppInstances = sumStats.TotalUsedDiskAppInstances + *cs.ContainerMetric.DiskBytes
			}
		}
		sumStats.TotalCpuPercentage = sumStats.TotalCpuPercentage + appStats.TotalCpuPercentage
		if appStats.TotalTraffic.EventL60Rate > 0 {
			sumStats.TotalActiveApps++
		}
		sumStats.TotalApps++
	}

	appMdMgr := asUI.GetEventProcessor().GetMetadataManager().GetAppMdManager()
	for _, app := range appMdMgr.AllApps() {
		sumStats := summaryStatsByStack[app.StackGuid]
		if sumStats != nil {
			if app.State == "STARTED" {
				sumStats.ReservedMem = sumStats.ReservedMem + ((app.MemoryMB * util.MEGABYTE) * app.Instances)
				sumStats.ReservedDisk = sumStats.ReservedDisk + ((app.DiskQuotaMB * util.MEGABYTE) * app.Instances)
			}
		}
	}

	for _, cellStats := range asUI.GetDisplayedEventData().CellMap {
		//toplog.Info("cellStats.StackId:%v", cellStats.StackId)
		sumStats := summaryStatsByStack[cellStats.StackId]
		if sumStats != nil {
			sumStats.TotalCellCPUs = sumStats.TotalCellCPUs + cellStats.NumOfCpus
			sumStats.TotalCapacityMemory = sumStats.TotalCapacityMemory + cellStats.CapacityTotalMemory
			sumStats.TotalCapacityDisk = sumStats.TotalCapacityDisk + cellStats.CapacityTotalDisk
		}
	}

	// Output stack information by stack name sort order
	stackSummaryStatsArray := make(StackSummaryStatsArray, 0, len(summaryStatsByStack))
	for _, stackSummaryStats := range summaryStatsByStack {
		stackSummaryStatsArray = append(stackSummaryStatsArray, stackSummaryStats)
	}
	sort.Sort(stackSummaryStatsArray)
	linesWritten := 0
	for _, stackSummaryStats := range stackSummaryStatsArray {
		if stackSummaryStats.TotalApps > 0 || stackSummaryStats.TotalCellCPUs > 0 {
			linesWritten += asUI.outputHeaderForStack(g, v, stackSummaryStats)
		}
	}

	if linesWritten == 0 {
		// Likely the "stacks" metadata loaded but not the "apps" metadata
		fmt.Fprintf(v, "\n Waiting for even more data...")
		return 3, nil
	}

	return linesWritten, nil
}

// Called for each stack - this should output 3 lines:
//
//  Stack: cflinuxfs2
//     CPU:  8.4% Used,  800% Max,       Mem:   7GB Used,  63GB Max,  22GB Rsrvd
//     Apps:  122 Total, Cntrs:  127     Dsk:   7GB Used, 190GB Max,  27GB Rsrvd
//
func (asUI *AppListView) outputHeaderForStack(g *gocui.Gui, v *gocui.View, stackSummaryStats *StackSummaryStats) int {

	totalUsedMemoryAppInstancesDisplay := "--"
	totalUsedDiskAppInstancesDisplay := "--"
	totalCpuPercentageDisplay := "--"
	if stackSummaryStats.TotalReportingAppInstances > 0 {
		totalUsedMemoryAppInstancesDisplay = util.ByteSize(stackSummaryStats.TotalUsedMemoryAppInstances).StringWithPrecision(0)
		totalUsedDiskAppInstancesDisplay = util.ByteSize(stackSummaryStats.TotalUsedDiskAppInstances).StringWithPrecision(0)
		if stackSummaryStats.TotalCpuPercentage >= 100 {
			totalCpuPercentageDisplay = fmt.Sprintf("%.0f%%", stackSummaryStats.TotalCpuPercentage)
		} else {
			totalCpuPercentageDisplay = fmt.Sprintf("%.1f%%", stackSummaryStats.TotalCpuPercentage)
		}
	}

	cellTotalCPUCapacityDisplay := "--"
	if stackSummaryStats.TotalCellCPUs > 0 {
		cellTotalCPUCapacityDisplay = fmt.Sprintf("%v%%", (stackSummaryStats.TotalCellCPUs * 100))
	}

	capacityTotalMemoryDisplay := "--"
	if stackSummaryStats.TotalCapacityMemory > 0 {
		capacityTotalMemoryDisplay = fmt.Sprintf("%v", util.ByteSize(stackSummaryStats.TotalCapacityMemory).StringWithPrecision(0))
	}
	capacityTotalDiskDisplay := "--"
	if stackSummaryStats.TotalCapacityDisk > 0 {
		capacityTotalDiskDisplay = fmt.Sprintf("%v", util.ByteSize(stackSummaryStats.TotalCapacityDisk).StringWithPrecision(0))
	}
	//fmt.Fprintf(v, "\r")

	fmt.Fprintf(v, "Stack: %v\n", stackSummaryStats.StackName)

	// Active apps are apps that have had go-rounter traffic in last 60 seconds
	// Reporting containers are containers that reported metrics in last 90 seconds
	fmt.Fprintf(v, "   CPU:%6v Used,%6v Max,       ", totalCpuPercentageDisplay, cellTotalCPUCapacityDisplay)

	displayTotalMem := "--"
	totalMem := stackSummaryStats.ReservedMem
	if totalMem > 0 {
		displayTotalMem = util.ByteSize(totalMem).StringWithPrecision(0)
	}
	fmt.Fprintf(v, "Mem:%6v Used,", totalUsedMemoryAppInstancesDisplay)
	// Total quota memory of all running app instances
	fmt.Fprintf(v, "%6v Max,%6v Rsrvd\n", capacityTotalMemoryDisplay, displayTotalMem)

	fmt.Fprintf(v, "   Apps:%5v Total, Cntrs:%5v     ",
		stackSummaryStats.TotalApps,
		stackSummaryStats.TotalReportingAppInstances)

	displayTotalDisk := "--"
	totalDisk := stackSummaryStats.ReservedDisk
	if totalMem > 0 {
		displayTotalDisk = util.ByteSize(totalDisk).StringWithPrecision(0)
	}

	fmt.Fprintf(v, "Dsk:%6v Used,", totalUsedDiskAppInstancesDisplay)
	fmt.Fprintf(v, "%6v Max,%6v Rsrvd\n", capacityTotalDiskDisplay, displayTotalDisk)

	// Number of lines written
	return 3

}
