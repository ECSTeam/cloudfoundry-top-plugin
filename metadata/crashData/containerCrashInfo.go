// Copyright (c) 2017 ECS Team, Inc. - All Rights Reserved
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

package crashData

import (
	"regexp"
	"time"
)

var (
	regexExtractExitStatus = regexp.MustCompile(`Exited with status (.*)`)
	regexErrorOccured      = regexp.MustCompile(`[0-9]* error\(s\) occurred:\n\n\*[ ]?`)
	regexCancelled         = regexp.MustCompile(`\* cancelled`)
	regexLineFeed          = regexp.MustCompile(`\n`)
)

type ContainerCrashInfo struct {
	ContainerIndex  int
	CrashTime       *time.Time
	ExitDescription string
}

func NewContainerCrashInfo(containerIndex int, crashTime *time.Time, exitDescription string) *ContainerCrashInfo {
	exitDescriptionClean := CleanupExitDescription(exitDescription)
	//exitStatus := ExtractExitStatusFromExitDescription(exitDescription)
	info := &ContainerCrashInfo{ContainerIndex: containerIndex, CrashTime: crashTime,
		ExitDescription: exitDescriptionClean}
	return info
}

type ContainerCrashInfoSlice []*ContainerCrashInfo

func (p ContainerCrashInfoSlice) Len() int {
	return len(p)
}

func (p ContainerCrashInfoSlice) Less(i, j int) bool {
	return p[i].CrashTime.Before(*p[j].CrashTime)
}

func (p ContainerCrashInfoSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func ExtractExitStatusFromExitDescription(exitDescription string) string {
	parsedData := regexExtractExitStatus.FindAllStringSubmatch(exitDescription, -1)
	if parsedData != nil && len(parsedData) > 0 {
		dataArray := parsedData[0]
		if dataArray != nil && len(dataArray) > 1 {
			return dataArray[1]
		}
	}
	return exitDescription
}

func CleanupExitDescription(exitDescription string) string {
	filteredExitDesc := regexErrorOccured.ReplaceAllString(exitDescription, "")
	filteredExitDesc = regexCancelled.ReplaceAllString(filteredExitDesc, "")
	filteredExitDesc = regexLineFeed.ReplaceAllString(filteredExitDesc, "")
	return filteredExitDesc
}
