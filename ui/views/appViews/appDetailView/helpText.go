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
  CRH - Crashed container count in last 24 hours. 
	  NOTE: The total crash count for 24 hour at top of screen can be
	  larger then the sum of the crash counts across all current 
	  instances because the app may have been scaled to a larger
	  number of instances in the past.
  MEM_USED - Memory used by the container.
  MEM_FREE - Memory free in the container.
  DISK_USED - Disk used by container.
  DISK_FREE - Disk free in the container.
  RESP - Avg response time in milliseconds over last 60 seconds.
  LOG_OUT - Total number of log stdout events.  
  LOG_ERR - Total number of log stderr events.
  REQ/1 - Number of HTTP(S) request/responses in last 1 second.
  REQ/10 - Number of HTTP(S) request/responses in last 10 seconds.
  REQ/60 - Number of HTTP(S) request/responses in last 60 seconds.
  TOT_REQ - Count of all of the HTTP(S) request/responses.
  2XX - Count of HTTP(S) responses with status code 200-299.
  3XX - Count of HTTP(S) responses with status code 300-399.
  4XX - Count of HTTP(S) responses with status code 400-499.
  5XX - Count of HTTP(S) responses with status code 500-599.
  CELL_IP - IP address of the cell running the container.
  STRT_DUR - Start duration. Amount of time from container creation to
      healthy container. This will be less then the overall time spent
      in STARTING state.
  CCR - Create counter. Number of times the container has been created.
  STATE_TIME - Time when the container entered current state.
  CNTR_START_MSG - Last container start-up message.
  CNTR_START_MSG_TM - Time when last container start-up message occured.
  
NOTE: The HTTP counters are based on traffic through the 
go-router. Applications that talk directly container-to-
container will not show up in the REQ/TOT-REQ/nXX counters.

**Application Name Color Key:**
WHITE - Normal
CYAN  - Active. HTTP(S) traffic has been recieved in the last
		10 seconds. 
YELLOW - App instance starting.
RED - App instance is crashed or down.
PURPLE - App instance is terminated.

`

const HelpLocalViewKeybindings = `
**Display: **
Press 'd' to show app detail view menu.

**Clipboard menu: **
Press 'c' to open the clipboard menu.  This will copy to clipboard a
command you can paste in terminal window later.
`
