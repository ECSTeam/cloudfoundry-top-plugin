// Copyright (c) 2017 ECS Team, Inc. - All Rights Reserved
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

package helpView

const HelpHeaderText = `
**Header information:**

  Events       - Total number of events received by the foundation.
  Warm-up      - It can take up to 60 seconds to receive all event
                 information before stats are accurate.
  Duration     - Amount of time stats has been collecting data.
  Target       - The target URL of monitored foundation.
  IsoSeg       - Isolation Segment (shown if foundation has more then 1).
  Stack        - The Cloud Foundry stack where indented fields below
                 pertain.
  Cells        - Number of cells assigned to this IsoSeg & Stack.
    CPU (Used)   - Amount of CPU consumed by all app instances.
    CPU (Max)    - Sum of CPU capacity across all cells.
    Mem (Used)   - Amount of memory actually in use by all app
                   instances.
    Mem (Max)    - Sum of memory capacity across all cells.
    Mem (Rsrvd)  - Total amount of requested memory for all started
                   app instances.
    Dsk (Used)   - Amount of disk actually in use by all app
                   instances.
    Dsk (Max)    - Sum of disk capacity across all cells.
    Dsk (Rsrvd)  - Total amount of requested disk for all started
                   app instances.
    Apps (total) - Total number of applications deployed to
                   foundation.
    Cntrs        - Number of reporting containers
`
const HelpTopLevelDataViewKeybindings = `
**Display: **
Press 'd' to show data view menu.

**Quit: **
Press 'q' to quit application.
`

const HelpChildLevelDataViewKeybindings = `
**Exit view: **
Press 'x' to exit current view
`

const HelpCommonDataViewKeybindings = `
**Select item detail (if available): **
Press UP arrow or DOWN arrow to highlight an application row.
Press ENTER to select the highlighted application and show
additional detail.

**Order / Sort display: **
Press 'o' to show the sort order window allowing multi-column
sorting of any column.

**Filter display: **
Press 'f' to show the filter window which allows for filtering
which rows should be displayed

**Scroll columns into view: **
Press RIGHT or LEFT arrow to scroll the columns into view if the
window is not wide enough to view all columns.  You can also resize
terminal window to show more columns/rows (resize of cmd.exe window
is not supported on windows while top is running).

**Pause display update:**
Press 'p' to toggle pause display update.  When display update is
paused top will continue to capture statstics and display updated
values when unpaused.

**Refresh screen interval: **
Press 's' to set the sleep time between refreshes. Default
is 1 second.  Valid values are 0.1 - 60.  The refresh interval only
effects how often the client screen is refreshed, it has no effect
on frequency the foundation delivers events. Top uses passive
monitoring for stats, a faster refresh interval will not introduce
additonal load on the CF foundation.

**Refresh: **
Press SPACE to force an immediate screen refresh.

**Log Window: **
Press shift-D to open log window.  This shows internal top
logging messages.  This window will open automatically if any error
message is logged (e.g., connection timeouts).

**Clear stats: **
Press shift-C to clear the statistics counters.

**Reload metadata: **
Press 'r' to force a reload of metadata for app/space/org.  The
metadata is loaded at startup and attempts to stay current by
recognizing when specific data needs to be reloaded. However there
can be circumstances were data becomes stale.
`
