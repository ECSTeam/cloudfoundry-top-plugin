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

package appCrashView

import (
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func ColumnContainerIndex() *uiCommon.ListColumn {
	defaultColSize := 4
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayContainerCrashInfo).ContainerIndex < c2.(*DisplayContainerCrashInfo).ContainerIndex)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerCrashInfo)
		display := fmt.Sprintf("%4v", stats.ContainerIndex)
		return fmt.Sprintf("%4v", display)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayContainerCrashInfo)
		return fmt.Sprintf("%v", stats.ContainerIndex)
	}
	c := uiCommon.NewListColumn("IDX", "IDX", defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func ColumnCrashTime() *uiCommon.ListColumn {
	defaultColSize := 20
	sortFunc := func(c1, c2 util.Sortable) bool {
		return (c1.(*DisplayContainerCrashInfo).CrashTime.Before(*c2.(*DisplayContainerCrashInfo).CrashTime))
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerCrashInfo)
		return fmt.Sprintf("%20v", stats.CrashTimeFormatted)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayContainerCrashInfo)
		return fmt.Sprintf("%v", stats.CrashTime.UnixNano())
	}
	c := uiCommon.NewListColumn("CRASH_TIME", "CRASH_TIME", defaultColSize,
		uiCommon.TIMESTAMP, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func ColumnExitDescription() *uiCommon.ListColumn {
	defaultColSize := 40
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayContainerCrashInfo).ExitDescription < (c2.(*DisplayContainerCrashInfo).ExitDescription)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayContainerCrashInfo)
		if stats.ExitDescription != "" {
			return stats.ExitDescription
		} else {
			return "--"
		}
	}
	rawValueFunc := func(data uiCommon.IData) string {
		appStats := data.(*DisplayContainerCrashInfo)
		return appStats.ExitDescription
	}
	c := uiCommon.NewListColumn("EXIT_DESCRIPTION", "EXIT_DESCRIPTION", defaultColSize,
		uiCommon.ALPHANUMERIC, true, sortFunc, false, displayFunc, rawValueFunc, nil)
	return c
}
