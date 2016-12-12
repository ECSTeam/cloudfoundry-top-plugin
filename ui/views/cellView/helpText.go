package cellView

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
 
**Cell list stats:**

  COL_A - TODO
  COL_B - TODO
  

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
on frequency the platform delivers events. Top uses passive
monitoring for stats, a faster refresh interval will not introduce
additonal load on the CF foundation.

**Select detail: **
Press UP arrow or DOWN arrow to highlight a cell row.
Press ENTER to select the highlighted cell and show
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
