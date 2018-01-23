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

package appView

import "github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/helpView"

const HelpText = HelpOverviewText + helpView.HelpHeaderText + HelpColumnsText + HelpLocalViewKeybindings + helpView.HelpTopLevelDataViewKeybindings + helpView.HelpCommonDataViewKeybindings

const HelpOverviewText = `
**App Stats View**

App list view shows a list of all deployed applications regardless of
state. The full set of stats may not be available until the warm-up
period is complete. After the warm-up period, if any applications are
found not to be in the desired state (e.g., instances set to 4 but
only 3 are running) an alert will be displayed and the application will
be colored red.
`

const HelpColumnsText = `
**Application Columns:**

  APPLICATION - Application name. If application name has:
     () - App has been deleted. It will be removed from display shortly
     [] - App is not yet fully staged (i.e., it is still deploying)
     Asterisk (*) on the end of name means the application is being
        activily monitored. See Active Monitoring below for more info.
        See Application Name Color Key below for color information.
  SPACE - Space name.
  ORG - Organization name.
  DCR - Desired containers (instances).
  RCR - Total reporting containers (ideally should match DCR).
  CPU%% - Total CPU percent consumed by all containers.
  CRH - Crashed container count in last 24 hours.
  MEM_USED - Total memory used by all containers.
  DSK_USED - Total disk used by all containers.
  RESP - Avg response time in milliseconds over last 60 seconds.
  LOG_OUT - Total number of stdout log events for all instance of app.
  LOG_ERR - Total number of stderr log events for all instance of app.
  REQ/1 - Number of HTTP(S) request/responses in last 1 second.
  REQ/10 - Number of HTTP(S) request/responses in last 10 seconds.
  REQ/60 - Number of HTTP(S) request/responses in last 60 seconds.
  TOT_REQ - Count of all of the HTTP(S) request/responses.
  2XX - Count of HTTP(S) responses with status code 200-299.
  3XX - Count of HTTP(S) responses with status code 300-399.
  4XX - Count of HTTP(S) responses with status code 400-499.
  5XX - Count of HTTP(S) responses with status code 500-599.
  ISO_SEG - Isolation Segment assigned to space.
  STACK - The Cloud Foundry stack used by this app.

NOTE: The HTTP counters are based on traffic through the 
go-router. Applications that talk directly container-to-
container will not show up in the REQ/TOT-REQ/nXX counters.

**Application Name Color Key:**
WHITE - Normal
CYAN  - Active. HTTP(S) traffic has been recieved in the last
        10 seconds. 
RED   - Not in desired state. The application has been configured
        for DCR instances, but only RCR are running.
GRAY  - One of two possibilities:
        1. App has been deleted.
        2. Not monitored (non-privileged only). In non-privileged
           only the first 50 applications in the currently targeted
           org and space are monitored.
 
**Active Monitoring:**
For the most part cf top passively monitors a platform by analyzing
firehose events as they occur. Because of limitations of this data
a realtime display of application container state is not possible.
To work around this limitation, a call to an API is invoked when
specific events are seen on the firehose. This provides a realtime
view of an application's containers at the cost of performing API
calls. To limit the amount of additional traffic the API calls will
impose on a platform, these calls do not occur by default. Only 
when an application detail view is visited does active monitoring
for the visited app occur. This active monitoring will continue
to occur as long as the app detail view is displayed and 15 minutes
after its been displayed. This active monitoring is indicated by
an asterisk (*) displayed next to the application name on the
application view.

NOTE: Active monitoring does NOT mean polling. A call to the API
to get updated container information on a specific app is only invoked
when a container event occurs (e.g., START, STOP, CRASHED). If
the actively monitored application's containers are not changing state
then no additional API calls are made.
`

const HelpLocalViewKeybindings = `
**Clipboard menu: **
Press 'c' when a row is selected to open the clipboard menu.
This will copy to clipboard a command you can paste in 
terminal window later.
`
