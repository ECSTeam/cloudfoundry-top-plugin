# top-plugin

WORK IN PROGRESS - Currently this plugin doesn't do anything yet.

This plugin is based on the firehose plugin: https://github.com/cloudfoundry/firehose-plugin

## Installation

```bash
cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
cf install-plugin ./top-plugin-osx
```

## Usage

User must be logged in as admin

### Options

List top stats for CF.

```
NAME:
   top - Displays top stats

USAGE:
   cf nozzle

OPTIONS:
   -debug                 -d, enable debugging
   -filter                -f, specify message filter such as LogMessage, ValueMetric, CounterEvent, HttpStartStop
   -no-filter             -n, no firehose filter. Display all messages
   -subscription-id       -s, specify subscription id for distributing firehose output between clients
```


