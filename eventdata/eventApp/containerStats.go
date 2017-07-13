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
	CellLastMsgText string
	CellLastMsgTime *time.Time
	// We need the last "Creating container" time so we can ignore "Destroying container" and "Successfully destroyed container"
	// messages that occur after this time.
	// We need to do tall this because the destroying of a container is async so a new container can be created while the old
	// container is still being destroyed.
	CellLastCreatingContainer *time.Time
	CellLastExitStatus        *time.Time
}

func NewContainerStats(containerIndex int) *ContainerStats {
	stats := &ContainerStats{ContainerIndex: containerIndex}
	return stats
}
