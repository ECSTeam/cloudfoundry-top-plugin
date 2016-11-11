# top-plugin

Plugin for showing live statistics of the targeted cloud foundry foundation.  You
must be logged in as 'admin' for this plugin to function.  The live statistics
include application stats and route stats.  The primary source of information
that the top plugin uses is from monitoring firehose.

![Screenshot](screenshots/screencast1.gif?raw=true)

## Application Stats

Provide details about applications running on the platform including the following
stats:

* HTTP(s) request/responses - Total counts, response time, status code, request rate that are routed through the gorouter
application.
* Containers - Number of reporting containers per application
* CPU% - Total CPU percent used by all instances of a given application
* Logs - Total log events that have occurred by all instances of a given application

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

* **Download the binary file for your target OS from the latest [release](https://github.com/kkellner/cloudfoundry-top-plugin/releases/latest)**
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

Login and add permission two permisson.  Note that the UAA password is NOT the
"Admin Credentials", the password is found in the ERT under Credentails tab,
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

# Background

This plugin is based on the firehose plugin: https://github.com/cloudfoundry/firehose-plugin
