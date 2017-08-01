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

App detail view shows details of the selected application.  

**Request Info Section**
Request Info section shows HTTP(S) request rates for the last 1/10/60
seconds. The average response time in milliseconds is also displayed 
for the same 1/10/60 second intervals.

NOTE: The HTTP(S) counters are based on traffic through the gorouter.
Applications that talk directly container-to-container will not be
reflected in the rate and response time values.

**Crash Info Section**
Crash Info section shows how many application containers have crashed
in the last 10 minutes, 1 hour, and 24 hours.  It also shows the last
time a container crashed in the previous 24 hours.
`

const HelpColumnsText = `
**Container Columns:**

  IDX - Application container index.
  STATE - Current container state: 
      DOWN, STARTING, RUNNING, CRASHED, TERM, UNKNOWN.
  STATE_DUR - Duration of time container has been in current state.
  CPU%% - CPU percent consumed by container.
  MEM_USED - Memory used by the container.
  MEM_FREE - Memory free in the container.
  DISK_USED - Disk used by container.
  DISK_FREE - Disk free in the container.
  LOG_OUT - Total number of log stdout events.  
  LOG_ERR - Total number of log stderr events.
  CELL_IP - IP address of the cell running the container.
  STRT_DUR - Start duration. Amount of time from container creation to
      healthy container. This will be less then the overall time spent
      in STARTING state.
  CCR - Create counter. Number of times the container has been created.
  STATE_TIME - Time when the container entered current state.
  CNTR_START_MSG - Last container start-up message.
  CNTR_START_MSG_TM - Time when last container start-up message occured.
`

const HelpLocalViewKeybindings = `
**Display: **
Press 'd' to show app detail view menu.
`
