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

TODO: Put in community plugins yet
[comment]: # ```bash
[comment]: # cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
[comment]: # cf install-plugin ./top-plugin-osx
[comment]: # ```

# Usage

User must be logged in as admin

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
