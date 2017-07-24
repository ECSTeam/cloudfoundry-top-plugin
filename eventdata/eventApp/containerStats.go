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
	ContainerIndex          int
	Ip                      string
	ContainerMetric         *events.ContainerMetric
	LastUpdateTime          *time.Time
	LastContainerUpdateTime *time.Time
	OutCount                int64
	ErrCount                int64
	StartTime               *time.Time
	Uptime                  *time.Duration
	CellLastStartMsgText    string
	CellLastStartMsgTime    *time.Time
}

func NewContainerStats(containerIndex int) *ContainerStats {
	stats := &ContainerStats{ContainerIndex: containerIndex}
	return stats
}
