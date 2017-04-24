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

package routeView

import "github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/helpView"

const HelpText = HelpOverviewText + helpView.HelpHeaderText + HelpColumnsText + helpView.HelpTopLevelDataViewKeybindings + helpView.HelpCommonDataViewKeybindings

const HelpOverviewText = `
**Route Stats View**

Route list view shows a list of all HTTP(s) traffic flowing through
the go-router.  This view can provide different information from
the App Stats view as a single route can be assinged to multiple 
applications.  E.g., blue-green deployments.  
`

const HelpColumnsText = `
**Route Columns:**

  HOST - Host name of the URL request
  DOMAIN - Domain name of the URL request
  PATH - Path (only shown if route has path routing)
  PORT - TCP Port (if TCP route type)
  R_APPS - Reporting Apps mapped to this route (may not be current
      until route details it accessed)
  TOT_REQ - Count of all of the HTTP(S) request/responses
  2XX - Count of HTTP(S) responses with status code 200-299
  3XX - Count of HTTP(S) responses with status code 300-399
  4XX - Count of HTTP(S) responses with status code 400-499
  5XX - Count of HTTP(S) responses with status code 500-599
  RESP_DATA - Total size of response data that has been sent to
     client.
  LAST_ACCESS - Last time a reponse was sent 

NOTE: The HTTP counters are based on traffic through the 
go-router.  Applications that talk directly container-to-
container will not show up in the REQ/TOT-REQ/nXX counters.

Route List Color Key:
WHITE - Normal
CYAN  - Active.  HTTP(S) traffic has been recieved in the last
        10 seconds. 
GRAY  - Not monitored (non-privileged only). In non-privileged only
        the first 50 applications in the currently targeted org and
        space are monitored.
`
