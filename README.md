# top-plugin

Cloud Foundry command-line cf plugin for showing live statistics of the targeted cloud foundry foundation.
You must be logged in as 'admin' user or assign permissions to the cloud foundry user as 
described in the installation instructions below for this plugin to function.  
The live statistics include application statistics and diego cell statistics amoung others.
The primary source of information that the top plugin uses is from monitoring the cloud foundry firehose.

![Screenshot](screenshots/screencast1.gif?raw=true)

## Screenshots

More [screenshot here](screenshots/screenshots.md)

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
* TOT-REQ - All of the HTTP(S) responses that have gone through go-router
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


# Installation

* **Download the binary file for your target OS from the latest [release](https://github.com/ecsteam/cloudfoundry-top-plugin/releases/latest)**
* If you've already installed the plugin and are updating, you must first run `cf uninstall-plugin TopPlugin`
* Then install the plugin with `cf install-plugin top-plugin-darwin`  (or `top-plugin-linux` or `top-plugin.exe`)
* If you get a permission error run: `chmod +x top-plugin-darwin` (or `top-plugin-linux`) on the binary
* Verify the plugin installed by looking for it with `cf plugins`

TODO: Register plugin with the community cloud foundry plugins website (https://plugins.cloudfoundry.org/)
<!---
```bash
cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
cf install-plugin ./top-plugin-osx
```
-->
## Assign needed permissions

To use this plugin you must be logged in as 'admin' or assign two permissions
to an existing cloud foundry user.  To assign needed permissions:

Install the uaac client CLI if you do not already have it:
```
gem install cf-uaac
```

Login and add two permission.  Note that the UAA password is NOT the
"Admin Credentials", the password is found in the ERT under Credentials tab,
look for password for "Admin Client Credentials".

```
uaac token client get admin -s [UAA Admin Client Credentials]  
uaac member add cloud_controller.admin [username]
uaac member add doppler.firehose [username]
```
Note: The change in permissions does not take effect until user username performs
a logout and login.


# Usage

User must be logged in as admin or cloud foundry user with permissions as described above.
```
cf top
```

## Options

List top live statistics for CF.

```
NAME:
   top - Displays top live statistics

USAGE:
   cf top

OPTIONS:
   -debug                 -d, enable debugging
   -cygwin                -c, force run under cygwin (Use this to run: 'cmd /c start cf top -cygwin' )
```
