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

package cellDetailView

import "github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/helpView"

const HelpText = HelpOverviewText +
	helpView.HelpHeaderText +
	HelpColumnsText +
	helpView.HelpChildLevelDataViewKeybindings +
	helpView.HelpCommonDataViewKeybindings

const HelpOverviewText = `
**Cell Detail Stats View**

Cell detail view shows a list of all containers running on the selected
diego cell. The full set of stats may not be available until the warm-up
period is complete.  
`

const HelpColumnsText = `
**Container stats:**

  APPLICATION - Application name
  IDX - Application (container) index
  SPACE - Space name
  ORG - Organization name
  CPU%% - Total CPU percent consumed by all containers on cell
  MEM_RSVD - Total memory reserved by all containers on cell
  MEM_USED - Total memory actually in use by all containers
  MEM_FREE - Total memory actually in use by all containers
  DISK_RSVD - Total disk reserved by all containers on cell
  DISK_USED - Total disk actually in use by all containers
  DISK_FREE - Free Disk space in cell VM available for containers
  LOG_OUT - Number of stdout log events
  LOG_ERR - Number of stderr log events
`
