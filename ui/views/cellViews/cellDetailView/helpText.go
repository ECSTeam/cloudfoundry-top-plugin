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

const HelpText = `
**Cell Detail Stats View**

Cell detail view shows a list of all containers running on the selected
diego cell. The full set of stats may not be available until the warm-up
period is complete.  

**Header information:**

TODO
 
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
  
**Exit view: **
Press 'x' to exit current view

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
Press 'r' to reload metadata for app/space/org.  The metadata
is loaded at top startup but can become stale if new applications
are deployed while top is running.
TODO: Auto reload metadata upon unknown translation

**Refresh screen interval: **
Press 's' to set the sleep time between refreshes. Default
is 1 second.  Valid values are 0.1 - 60.  The refresh interval only
effects how often the client screen is refreshed, it has no effect
on frequency the foundation delivers events. Top uses passive
monitoring for stats, a faster refresh interval will not introduce
additonal load on the CF foundation.

**Scroll columns into view: **
Press RIGHT or LEFT arrow to scroll the columns into view if the
window is not wide enough to view all columns.  You can also resize
terminal window to show more columns/rows (resize of cmd.exe window
is not supported on windows while top is running).

**Refresh: **
Press SPACE to force an immediate screen refresh
`
