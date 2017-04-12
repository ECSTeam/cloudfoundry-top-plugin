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

package eventApp

import (
	"time"

	"github.com/cloudfoundry/sonde-go/events"
)

type ContainerStats struct {
	ContainerIndex  int
	Ip              string
	ContainerMetric *events.ContainerMetric
	LastUpdate      time.Time
	OutCount        int64
	ErrCount        int64
	//CrashCount         int
	//LastCrashTime      *time.Time
	ContainerCrashInfo []*ContainerCrashInfo
}

func NewContainerStats(containerIndex int) *ContainerStats {
	stats := &ContainerStats{ContainerIndex: containerIndex}
	return stats
}

func (cs *ContainerStats) AddCrashInfo(crashTime *time.Time, exitDescription string) {
	crashInfo := NewContainerCrashInfo(crashTime, exitDescription)
	if cs.ContainerCrashInfo == nil {
		cs.ContainerCrashInfo = make([]*ContainerCrashInfo, 0, 10)
	}
	cs.ContainerCrashInfo = append(cs.ContainerCrashInfo, crashInfo)
}

func (cs *ContainerStats) CrashCount() int {
	return len(cs.ContainerCrashInfo)
}

func (cs *ContainerStats) LastCrashTime() *time.Time {
	if cs.ContainerCrashInfo != nil && len(cs.ContainerCrashInfo) > 0 {
		last := len(cs.ContainerCrashInfo) - 1
		return cs.ContainerCrashInfo[last].CrashTime
	}
	return nil
}

type ContainerCrashInfo struct {
	CrashTime       *time.Time
	ExitDescription string
}

func NewContainerCrashInfo(crashTime *time.Time, exitDescription string) *ContainerCrashInfo {
	info := &ContainerCrashInfo{CrashTime: crashTime, ExitDescription: exitDescription}
	return info
}
