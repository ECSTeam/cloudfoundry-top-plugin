
# Cloud Foundry `top` plugin by Kurt Kellner of ECS Team (www.ECSTeam.com)

Be sure to also check [Frequently Asked Questions (FAQ)](faq.md) page
for additional information about `top`.

# Background
The `cf top` plugin was written to fill a need I encountered while performing
Cloud Foundry production support for one of our clients.  In this case the
foundation was configured to forward application `log` events from the firehose
to Splunk however due to event volumes in the firehose, the foundation events
were not saved.  This made it difficult to answer the question "What is happening
on the foundation right now?"  

What I needed was a simple tool to see what was happening.  Simple to install
and simple to run.  I wanted something that did not require configuration
or deployment into the foundation I wanted to monitor.  The idea was to create
a tool much like the UNIX `top` command which provides a wealth of information
but is simple and easily accessible.

From this need, the Cloud Foundry cf cli `top` plugin was born.

# Overview
The `top` plugin is unlike most `cf` plugins which typically run using various
command-line arguments, output information to the terminal and then exit.  
Instead, when you run `cf top` on the command-line, it will initialize a text-based 
interface that will allow user interaction with the screen.  It works much like
the UNIX top command.

There are 5 top-level data views.  When switching between data views using 'd', you
must be at a top-level view.  When in a subview, you can press 'x' to exit
back.  If you need help you can look at the bottom of the window which 
provides a quick-help of common commands or press 'h' for verbose help
which is specific to the current view.

## View / subview layout:

* App Stats (default view at startup)
    * App Details
* Cell Stats
    * Cell Details
* Route Stats
    * Route Map
* Event Stats
    * Event Origin
        * Event Detail
* Capacity Plan (memory)



## Application Stats

Provide details about applications running on the foundation including the following
stats:

* APPLICATION - Application name
* SPACE - Space name
* ORG - Organization name
* DRC - Total desired containers (configured instances of app)
* RCR - Total reporting containers 
* CPU% - Total CPU percent consumed by all containers
* MEM - Total memory used by all containers
* DISK - Total disk used by all containers
* RESP - Avg response time in milliseconds over last 60 seconds
* LOG_OUT - Total number of stdout log events
* LOG_ERR - Total number of stderr log events
* REQ/1 - HTTP(S) request/responses in last 1 second that have gone through go-router
* REQ/10 - HTTP(S) request/responses in last 10 seconds that have gone through go-router
* REQ/60 - HTTP(S) request/responses in last 60 seconds that have gone through go-router
* TOT_REQ - All of the HTTP(S) responses that have gone through go-router
* 2XX - HTTP(S) responses with status code 200-299
* 3XX - HTTP(S) responses with status code 300-399
* 4XX - HTTP(S) responses with status code 400-499
* 5XX - HTTP(S) responses with status code 500-599

### Application details

A specific application can be selected to show details of a specific application including:

* IDX - Container instance index
* CPU% - Per container (application instance) CPU percent usage
* MEM_USED - Per container (application instance) memory used
* MEM_FREE - Per container (application instance) memory free
* DISK_USED - Per container (application instance) disk used
* DISK_FREE - Per container (application instance) disk free
* LOG_OUT - Per container (application instance) stdout log events
* LOG_ERR - Per container (application instance) stderr log events

## Cell Stats

Provides details about each diego cell in the foundation.  This information can be used to find
"hot" cells -- Cells with high CPU untilization.  A cell can be selected to get details 
such as which containers are currenly running on the cell.

* IP - The IP address of the diego cell
* CPU% - Container CPU percent usage
* RCR - Number of container that have reported in
* CPUS - Number of CPUs assigned to this diego cell
* TOT_MEM -
* FREE_MEM -
* C_RSVD_MEM -
* C_USD_MEM -
* TOT_DISK -
* FREE_DISK -
* C_RSVD_DSK -
* C_USD_DSK - 
* MAX_CNTR -
* CNTRS - Number of containers cell thinks it has running
* DNAME - BOSH deployment name
* JOB_NAME - BOSH job name
* JOB_IDX - BOSH deployment index