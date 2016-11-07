# top-plugin

Plugin for showing live statistics of the targeted cloud foundry foundation.  You
must be logged in as 'admin' for this plugin to function.  The live statistics
include application stats and route stats.  The primary source of information
that the top plugin uses is from monitoring firehose.

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


* Download the binary file for your target OS from the latest [release](https://github.com/kkellner/cloudfoundry-top-plugin/releases/latest)
* If you've already installed the plugin and are updating, you must first run cf uninstall-plugin TopPlugin
* Then install the plugin with cf install-plugin top-plugin-darwin   (or top-plugin-linux or top-plugin.exe)
* If you get a permission error run: chmod +x top-plugin-darwin (or top-plugin-linux) on the binary
* Verify the plugin installed by looking for it with cf plugins

TODO: Register plugin with the community cloud foundry plugins website (https://plugins.cloudfoundry.org/)
<!---
```bash
cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
cf install-plugin ./top-plugin-osx
```
-->

# Usage

User must be logged in as admin.
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
