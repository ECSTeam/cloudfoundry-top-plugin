# top-plugin

This is a Cloud Foundry command-line cf interactive plugin for showing live statistics of the targeted Cloud Foundry foundation.
The live statistics include application statistics and route statistics among others.
The primary source of information that the top plugin uses is via monitoring the Cloud Foundry firehose.

The plugin will run in one of two modes, privileged or non-privileged depending on your Cloud Foundry user permission.
If you are a foundation operator you will want to use top in privileged mode.  This is done automatically if the
correct permissions are granted to your Cloud Foundry login (or if you are logged in via `admin` account).  See
[Assign Permissions](#assign-permissions-if-privileged-mode-is-needed) for more information on assigning permissions.


[Installation Instructions](#installation) 

![Screenshot](screenshots/screencast1.gif?raw=true)

## Screenshots

More [screenshots here](screenshots/screenshots.md)

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
There are two options for installation; use the plugin repo (recommended) or manual installation.

## Install from plugin repository (recommended)
NOTE: This installation method requires that your client computer has access to the internet.
If internet access is not available from client computer use the manual method.

Verify you have a repo named `CF-Community` registered in your cf client.

```
cf list-plugin-repos
```
If the above command does not show `CF-Community` you can add the repo via:

```
cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
```
Now that we have the cloud foundry community repo registered, install `top`:

```
cf install-plugin -r CF-Community "top"
```


## Manual installation method
* **Download the binary file for your target OS from the latest [release](https://github.com/ecsteam/cloudfoundry-top-plugin/releases/latest)**
* If you've already installed the plugin and are updating, you must first run `cf uninstall-plugin top`
* Then install the plugin with `cf install-plugin top-plugin-darwin`  (or `top-plugin-linux` or `top-plugin.exe`)
* If you get a permission error run: `chmod +x top-plugin-darwin` (or `top-plugin-linux`) on the binary
* Verify the plugin installed by looking for it with `cf plugins`

## Upgrade to latest version
To upgrade to the lastest version of top plugin, uninstall plugin and install again.
```
cf uninstall-plugin top
cf install-plugin -r CF-Community "top"     (or use manual install method described above)
```

## Assign permissions if privileged mode is needed

The `top` plugin will run without special permissions however it determines at runtime
what permissions you have and displays the appropriate functionality based on those
permissions.  If you are a foundation operator you will want the additional functionality
that top provides to privileged users.

If you are logged in with the Cloud Foundry `admin` account, no additional permissions
are needed, the `admin` account has everything it needs to run top with full functionality.

For non-admin accounts, to run top in privileged mode you need to assign two permissions
to an existing Cloud Foundry user.  To assign needed permissions:

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
The plugin does not require any arguments.  Simply run:
```
cf top
```

## Options

List top live statistics for CF.

```
NAME:
   top - Displays top stats - by Kurt Kellner of ECS Team

USAGE:
   cf top

OPTIONS:
   -debug              -d, enable debugging
   -no-top-check       -ntc, do not check if there are other instances of top running
   -nozzles            -n, specify the number of nozzle instances (default: 2)
   -cygwin             -c, force run under cygwin (Use this to run: 'cmd /c start cf top -cygwin' )
```
