package cellView

const helpText = `
**Header information:**

  Evnts        - Total number of events received by the platform.
  Warm-up      - It can take up to 30 seconds to receive all event
                 information before stats are accurate.
  Duration     - Amount of time stats have been collected.
  Target       - The target URL of monitored foundation.
  
 
**Cell list stats:**

  IP - IP address of Cloud Foundry cell
  CPU%% - Total CPU percent consumed by all containers on cell
  RCR - Total reporting containers
  CPUS - Number of CPUs in cell VM
  TOT_MEM - Total Memory in cell VM
  FREE_MEM - Free Memory in cell VM available for containers
  C_RSVD_MEM - Total memory reserved by all containers on cell
  C_USD_MEM - Total memory actually in use by all containers
  TOT_DISK - Total Disk space in cell VM
  FREE_DISK - Free Disk space in cell VM available for containers
  C_RSVD_DSK - Total disk reserved by all containers on cell
  C_USD_DSK - Total disk actually in use by all containers
  MAX_CNTR - Max containers cell can handle
  CNTRS - Number of containers running on cell reported by cell
  DNAME - BOSH Deployment name
  JOB_NAME - BOSH Job name
  JOB_IDX - BOSH Job index
  
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
TODO: Press ENTER to select the highlighted cell and show
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
