package appStats

const helpText = `
**Header information:**

  Total events - Total number of events received by the platform
  Warm-up period - It can take up to 30 seconds to receive all event
        information before stats are accurate
  Stats duration - Amount of time stats have been collected
  Target - The target URL of monitored foundation
  Rsrvd Mem - Total amount of requested memory for all started app
        instances
  Rsrvd Disk - Total amount of requested disk for all started app
        instances
  Total Apps - Total number of applications deployed to foundation
  Active Apps - Number of applications that have taken HTTP(S) traffic
        through the go router in the last 60 seconds
  Rprt Cntnrs - Number of reporting containers which typically are app
        instances
  Used Mem - Amount of memory actually in use by all app instances
  Used Disk - Amount of disk actually in use by all app instances
  Used CPU - Amount of CPU consumed by all app instances

**Application list stats:**

  APPLICATION - Application name
  SPACE - Space name
  ORG - Organization name
  RCR - Total reporting Containers
  CPU%% - Total CPU percent consumed by all containers
  MEM - Total memory used by all containers
  DISK - Total disk used by all containers
  RESP - Avg response time in milliseconds over last 60 seconds
  LOGS - Total number of log events
  L1 - HTTP(S) request/responses in last 1 second
  L10 - HTTP(S) request/responses in last 10 seconds
  L60 - HTTP(S) request/responses in last 60 seconds
  HTTP - All of the HTTP(S) responses
  2XX - HTTP(S) responses with status code 200-299
  3XX - HTTP(S) responses with status code 300-399
  4XX - HTTP(S) responses with status code 400-499
  5XX - HTTP(S) responses with status code 500-599

**Order / Sort display: **
Press 'o' to show the sort order window allowing multi-column
sorting of any column.

**Clear stats: **
Press 'c' to clear the statistics counters.

**Pause display update:**
Press 'p' to toggle pause display update.  When display update is
paused top will continue to capture statstics and display updated
values when unpaused.

**Filter display: **
Press 'f' to show the filter window which allows for filtering
which rows should be displayed

**Refresh screen interval: **
Press 's' to set the sleep time between refresh intervals. Default
is 1 second.  Valid values are 0.1 - 60.  The refresh interval only
effects how often the client screen is refreshed, it has no effect
on frequency the platform delivers events. Top uses passive
monitoring for stats and increating the refresh interval will not
introduce additonal load on the platform.

**Select application detail: **
Press UP arrow or DOWN arrow to highlight an application row.
Press ENTER to select the highlighted application and show
additional detail.

**Scroll columns into view: **
Press RIGHT arrow or LEFT arrow to scroll the columns into view if
the window is not wide enough to view all columns.  You can also
resize terminal window to show more columns/rows (resize not
supported on windows).

**Refresh: **
Press SPACE to force an immediate screen refresh

**Quit: **
Press 'q' to quit application

**Debug: **
Press 'D' (uppercase) to open debug window.  This shows internal
top logging events.  This window will open automatically if any
error message is logged (e.g., connection timeouts)

`
