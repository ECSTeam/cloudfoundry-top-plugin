package appDetailView

const helpText = `
**Header information:**

  Evnts        - Total number of events received by the platform.
  Warm-up      - It can take up to 30 seconds to receive all event
                 information before stats are accurate.
  Duration     - Amount of time stats have been collected.
  Target       - The target URL of monitored foundation.
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
  Apps (Actv)  - Number of applications that have taken HTTP(S)
                 traffic through the go router in the last 60
                 seconds.
  Cntrs        - Number of reporting containers which typically
                 are app instances.
 
**Application list stats:**

  APPLICATION - Application name
  SPACE - Space name
  ORG - Organization name
  RCR - Total reporting Containers
  CPU%% - Total CPU percent consumed by all containers
  MEM - Total memory used by all containers
  DISK - Total disk used by all containers
  RESP - Avg response time in milliseconds over last 60 seconds
  LOGS - Total number of log events for all instance of app
  REQ/1 - Number of HTTP(S) request/responses in last 1 second
  REQ/10 - Number of HTTP(S) request/responses in last 10 seconds
  REQ/60 - Number of HTTP(S) request/responses in last 60 seconds
  TOT-REQ - Count of all of the HTTP(S) request/responses
  2XX - Count of HTTP(S) responses with status code 200-299
  3XX - Count of HTTP(S) responses with status code 300-399
  4XX - Count of HTTP(S) responses with status code 400-499
  5XX - Count of HTTP(S) responses with status code 500-599

NOTE: The HTTP counters are based on traffic through the 
gorouter.  Applications that talk directly container-to-
container will not show up in the REQ/nXX counters.

**Display: **
Press 'd' to select data view.

**Order / Sort display: **
Press 'o' to show the sort order window allowing multi-column
sorting of any column.

**Clear stats: **
Press shift-C to clear the statistics counters.

**Clipboard menu: **
Press 'c' when a row is selected to open the clipboard menu.
This will copy to clipboard a command you can paste in 
terminal window later.

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
on frequency the platform delivers events. Top uses passive
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
Press SPACE to force an immediate screen refresh

**Quit: **
Press 'q' to quit application

**Log Window: **
Press shift-D to open log window.  This shows internal top
logging messages.  This window will open automatically if any error
message is logged (e.g., connection timeouts)

`
