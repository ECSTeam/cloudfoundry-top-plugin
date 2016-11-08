package appStats



const helpText = `
** Header information:**

  Total events - Total number of events received by the platform
  Stats duration - Amount of time stats have been collected

** Application list stats:**

  APPLICATION - Application name
  SPACE - Space name
  ORG - Organization name
  RCR - Total reporting Containers
  CPU% - Total CPU percent consumed by all containers
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

** Sorting display:**
 Press 's' to show the sort window allowing multi-column
 sorting of any column.

** Clear stats:**
 Press 'c' to clear the statistics counters.

** Pause display update:**
 Press 'p' to toggle pause display update.  When display update is paused
 top will continue to capture statstics and display updated values when
 unpaused.

** Filter display:**
 TODO

** Select application detail:**
 Press UP arrow or DOWN arrow to highlight an application row.
 Press ENTER to select the highlighted application and show
 additional detail.
`
