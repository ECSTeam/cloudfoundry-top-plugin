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

const HelpText = HelpOverviewText + helpView.HelpHeaderText + HelpColumnsText + HelpLocalViewKeybindings + helpView.HelpTopLevelDataViewKeybindings + helpView.HelpCommonDataViewKeybindings

const HelpOverviewText = `
**Space View**

Space view shows a list of all spaces in the selected organization on the foundation.
`

const HelpColumnsText = `
**Space Columns:**

  SPACE - Space name
  QUOTA_NAME - Space quota name if one is assigned
  APPS - Number of apps within the space
  DCR - Number of desired containers
  RCR - Number of reporting containers
  CPU% - Total CPU used by all containers within space
  MAX_MEM - Maximum memory space can used based on quota limits
  RSVD_MEM - Total memory reserved by all desired containers
  S_MEM% - Percent of space quota consumed
  O_MEM% - Percent of org quota consumed
  USED_MEM -
  RSVD_DSK -
  USED_DSK -
  LOG_OUT -
  LOG_ERR -
  TOT_REQ -

`

const HelpLocalViewKeybindings = `
**Clipboard menu: **
Press 'c' when a row is selected to open the clipboard menu.
This will copy to clipboard a command you can paste in 
terminal window later.
`
