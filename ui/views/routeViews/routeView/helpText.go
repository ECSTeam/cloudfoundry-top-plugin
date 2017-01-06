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

const helpText = `
**Route Stats View**

Route list view shows a list of all HTTP(s) traffic flowing through
the go-router.  This view can provide different information from
the App Stats view as a single route can be assinged to multiple 
applications.  E.g., blue-green deployments.  

**Header information:**

TODO
 
**Cell list stats:**

  HOST - Host name of the URL request
  DOMAIN - Domain name of the URL request
  PATH - Path (only shown if route has path routing)
  TOT-REQ - Count of all of the HTTP(S) request/responses
  2XX - Count of HTTP(S) responses with status code 200-299
  3XX - Count of HTTP(S) responses with status code 300-399
  4XX - Count of HTTP(S) responses with status code 400-499
  5XX - Count of HTTP(S) responses with status code 500-599
  RESP_DATA - Total size of response data that has been sent to
     client.
  M_GET - Count of HTTP(S) GET method requests
  M_POST - Count of HTTP(S) POST method requests
  M_PUT - Count of HTTP(S) PUT method requests
  M_DELETE - Count of HTTP(S) DELETE method requests
  LAST_ACCESS - Last time a reponse was sent 


NOTE: The HTTP counters are based on traffic through the 
go-router.  Applications that talk directly container-to-
container will not show up in the REQ/TOT-REQ/nXX counters.
  
**Display: **
Press 'd' to select data view.

**Order / Sort display: **
Press 'o' to show the sort order window allowing multi-column
sorting of any column.

**Clear stats: **
Press shift-C to clear the statistics counters.

**Pause display update:**
Press 'p' to toggle pause display update.  When display update is
paused top will continue to capture statstics and display updated
values when unpaused.

**Filter display: **
Press 'f' to show the filter window which allows for filtering
which rows should be displayed

**Reload metadata: **
Press 'r' to force a reload of metadata for app/space/org.  The
metadata is loaded at startup and attempts to stay current by
recognizing when specific data needs to be reloaded. However there
can be circumstances were data becomes stale.

**Refresh screen interval: **
Press 's' to set the sleep time between refreshes. Default
is 1 second.  Valid values are 0.1 - 60.  The refresh interval only
effects how often the client screen is refreshed, it has no effect
on frequency the foundation delivers events. Top uses passive
monitoring for stats, a faster refresh interval will not introduce
additonal load on the CF foundation.

**Select application detail: **
Press UP arrow or DOWN arrow to highlight an application row.
Press ENTER to select the highlighted application and show
additional detail.

**Scroll columns into view: **
Press RIGHT or LEFT arrow to scroll the columns into view if the
window is not wide enough to view all columns.  You can also resize
terminal window to show more columns/rows (resize of cmd.exe window
is not supported on windows while top is running).

**Refresh: **
Press SPACE to force an immediate screen refresh.

**Quit: **
Press 'q' to quit application.

**Log Window: **
Press shift-D to open log window.  This shows internal top
logging messages.  This window will open automatically if any error
message is logged (e.g., connection timeouts).
`
