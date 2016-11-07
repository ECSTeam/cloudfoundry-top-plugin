package appStats



const helpText = `
A list of application stats which include:

  APPLICATION - Application name
  SPACE - Space name
  ORG - Organization name
  2XX - HTTP(s) responses with status code 200-299
  3XX - HTTP(s) responses with status code 300-399
  4XX - HTTP(s) responses with status code 400-499
  5XX - HTTP(s) responses with status code 500-599
  TOTAL - All of the HTTP(s) responses
  L1 - Responses in last 1 second
  L10 - Responses in last 10 seconds
  L60 - Responses in last 60 seconds
  CPU% - Total CPU percent consumed by all containers
  RCR - Total reporting Containers
  RESP - Avg response time in milliseconds over last 60 seconds
  LOGS - Total number of log events

Sorting display:
Press 's' to show the sort window allowing multi-column
sorting of any column

Clear stats:
Press 'c' to clear the statistics counters

Filter display:
TODO

Select application detail:
Press UP arrow or DOWN arrow to highlight an application row.
Press ENTER to select the highlighted application and show
additional detail.
`
