package cellDetailView

const helpText = `
**Header information:**

 
**Container stats:**

  APPLICATION - Application name
  IDX - Application (container) index
  SPACE - Space name
  ORG - Organization name
  CPU%% - Total CPU percent consumed by all containers on cell
  RSVD_MEM - Total memory reserved by all containers on cell
  USD_MEM - Total memory actually in use by all containers
  FREE_MEM - Total memory actually in use by all containers
  
  RSVD_DSK - Total disk reserved by all containers on cell
  USD_DSK - Total disk actually in use by all containers
  DISK_FREE - Free Disk space in cell VM available for containers
  STDOUT - Number of stdout log events
  STDERR - Number of stderr log events
  
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
on frequency the platform delivers events. Top uses passive
monitoring for stats, a faster refresh interval will not introduce
additonal load on the CF foundation.

**Scroll columns into view: **
Press RIGHT or LEFT arrow to scroll the columns into view if the
window is not wide enough to view all columns.  You can also resize
terminal window to show more columns/rows (resize of cmd.exe window
is not supported on windows while top is running).

**Refresh: **
Press SPACE to force an immediate screen refresh

**Log Window: **
Press shift-D to open log window.  This shows internal top
logging messages.  This window will open automatically if any error
message is logged (e.g., connection timeouts)

`
