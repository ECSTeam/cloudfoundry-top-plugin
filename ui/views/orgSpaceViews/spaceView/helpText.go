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

package spaceView

import "github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/helpView"

const HelpText = HelpOverviewText +
	helpView.HelpHeaderText +
	HelpColumnsText +
	HelpLocalViewKeybindings +
	helpView.HelpTopLevelDataViewKeybindings +
	helpView.HelpCommonDataViewKeybindings

const HelpOverviewText = `
**Space View**

Space view shows a list of all spaces in the selected organization on the foundation.
`

const HelpColumnsText = `
**Space Columns:**

  SPACE - Space name
  QUOTA_NAME - Space quota name if one is assigned
  APPS - Number of apps within the space
  DCR - Number of desired containers (app instances)
  RCR - Number of reporting containers which are the the actual number of app
      instances running.  Normally DCR and RCR are equal.
  CPU%% - Total CPU used by all containers within space
  MEM_MAX - Maximum memory space can used based on quota limits
  MEM_RSVD - Total memory reserved by all desired containers
  S_MEM% - Percent of space quota consumed
  O_MEM% - Percent of org quota consumed
  MEM_USED - Memory actually in use by all containers
  DSK_RSVD - Disk reserved by all containers on cell
  DSK_USED - Disk actually in use by all containers
  LOG_OUT - Total number of stdout log events for all instance of app
  LOG_ERR - Total number of stderr log events for all instance of app
  TOT_REQ - Count of all of the HTTP(S) request/responses

`

const HelpLocalViewKeybindings = `
**Clipboard menu: **
Press 'c' when a row is selected to open the clipboard menu.
This will copy to clipboard a command you can paste in 
terminal window later.
`
