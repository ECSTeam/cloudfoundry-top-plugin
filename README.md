# top-plugin

Plugin for showing live statistics of the targeted cloud foundry foundation.  You
must be logged in as 'admin' or assign permissions to PCF user as described in the
installation below for this plugin to function.  The live statistics include
application stats and route stats. The primary source of information
that the top plugin uses is from monitoring firehose.

![Screenshot](screenshots/screencast1.gif?raw=true)

## Application Stats

Provide details about applications running on the platform including the following
stats:

* APPLICATION - Application name
* SPACE - Space name
* ORG - Organization name
* RCR - Total reporting Containers
* CPU% - Total CPU percent consumed by all containers
* MEM - Total memory used by all containers
* DISK - Total disk used by all containers
* RESP - Avg response time in milliseconds over last 60 seconds
* LOGS - Total number of log events
* L1 - HTTP(S) request/responses in last 1 second
* L10 - HTTP(S) request/responses in last 10 seconds
* L60 - HTTP(S) request/responses in last 60 seconds
* HTTP - All of the HTTP(S) responses
* 2XX - HTTP(S) responses with status code 200-299
* 3XX - HTTP(S) responses with status code 300-399
* 4XX - HTTP(S) responses with status code 400-499
* 5XX - HTTP(S) responses with status code 500-599

### Application details

A specific application can be selected to show details of a specific application including:

* CPU% - Per container (application instance) CPU percent usage
* Memory - Per container (application instance) memory used
* Disk - Per container (application instance) disk used
* Log Stdout - Per container (application instance) stdout log events
* Log Stderr - Per container (application instance) stderr log events

## Route Stats (TODO - not implemented yet)

Not yet implemented - This will show statistics based on routes (domain/host/path)

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
to an existing PCF user.  To assign needed permissions:

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

User must be logged in as admin or pcf user with permissions as described above.
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
