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

package appDetailView

import "github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/helpView"

const HelpText = HelpOverviewText +
	helpView.HelpHeaderText +
	HelpColumnsText +
	HelpLocalViewKeybindings +
	helpView.HelpChildLevelDataViewKeybindings +
	helpView.HelpCommonDataViewKeybindings

const HelpOverviewText = `
**App Detail View**

App detail view shows details of the selected application.  In App Request
Info area shows HTTP request rates for the last 1/10/60 seconds.  The
average response time in milliseconds is also displayed for the same
1/10/60 second intervals.

NOTE: The HTTP counters are based on traffic through the 
gorouter.  Applications that talk directly container-to-
container will not show up in the REQ/nXX counters.
`

const HelpColumnsText = `
 **Container Columns:**

  IDX - Application container index
  CPU%% - CPU percent consumed by container
  MEM_USED - Memory used by the container
  MEM_FREE - Memory free in the container
  DISK_USED - Disk used by container
  DISK_FREE - Disk free in the container
  LOG_OUT - Total number of log stdout events  
  LOG_ERR - Total number of log stderr events 
  CELL_IP - IP address of the cell running the container
`

const HelpLocalViewKeybindings = `
**Info View: **
Press 'i' to display application info view.
`
