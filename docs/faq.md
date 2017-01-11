# top-plugin Frequently Asked Questions (FAQ)


## Why do I get dropped message error?
When running top you might see the following error in the log window:

`Upstream message indicates the nozzle or the TrafficController is not keeping up. Dropped delta: 100, total: 100`

Top plugin uses the Cloud Foundry firehose to monitor real-time stats.  In cases
where a foundation is very busy or a slow network between the client running `top`
and the foundation, firehose messages can be dropped.  Because of the volume of these
messages, particularly if you are running `top` in privileged mode, you should ensure
you are running top on a client with as few network hops as possible to the foundation.
The optimal configuration is to run `top` on a jump box VM running on the same IaaS.

There can be other factors involved in dropped events including optimizing:
load balancer, go-routers, doppler and traffic controller instances.

`top` has been tested with over 30,000 events per second without dropping events.

## Sometimes when I run top, for the first few seconds I see garbage for the application name.  Is this a bug?
This is not a bug.  When `top` starts, its goal is to start capturing events and to 
display data as quickly as possible.  What you are seeing is not "garbage" but the 
application's GUID.  In the background requests are made for metadata that help
translate internal GUIDs into human readable names (e.g., app, space and org names).

## When I run top the interface looks terrible, nothing like the screenshots I've seen.  Why?
`top` uses several characters to draw borders and other indicators that can be problematic
on some terminal emulators / font-types / character sets.  If `top` is not displaying
correctly try changing the font-type or character set within your terminal emulator.  If
this still does not work, try a different terminal emulator.  For Windows users accessing
Linux, using Putty to ssh into Linux OS seems to work well to display `top`.

## Why do I sometimes get a big red window inside my session when running top?
`top` maintains an internal log of messages.  This log is accessible by pressing ctrl-shift-D.
The log window normally has a blue background but will change to a red background if
there are any errors that have been logged.  The log view will automatically open when
an error is logged. 

## Why does it take 60 seconds to "warm-up" when top starts?
When `top` is started, it has no information or history on the foundation its monitoring.
To be friendly to the foundation, it does not submit 100s of API requests to get the
current status of the foundation and all applications.  Instead `top` monitors the 
firehose for events and "learns" what it needs to know by passively listening.  This
can take up to 60 seconds to learn all that is needed and have accurate information to
display.  You should never jump to any conclusion about the health of a foundation
until the warm-up period is complete. 

## I just added a few Diego Cells to my foundation.  Why does `top` show `UNKNOWN (cells with no containers)` in the header?
`top` learns what type of stack is running on a cell by looking at the application 
containers running on the cell.  When you have cells with no containers, `top` is unable
to determine the type of stack for the cell.  Once containers are running on the cell
`top` will refresh the header to provide accurate stack information.

## Why is the `CPU percent Used` field in the header red when its no where near the `Max` shown in the header?
The colorization of the `CPU% Used` field in the header is based on warm (yellow)
and hot (red) cell CPU maximums.  The `CPU% Used` value will be colorized if any cell's CPU
utilization is at or above 80% of the cell's capacity. Yellow (>=80%), red (>=90%). 
To determine root cause of why this field is colorized, use the Cell Stats screen to
determine which of the cells are running warm/hot on CPU resource. 

## Why is `top` still showing application instances (containers) that are no longer running?
`top` tracks containers based on events in the firehose.  A container will output health
information periodically to the firehose.  When `top` does not see this health information
from a container for 90 seconds, it assumes the container is dead.  This means that it 
can take `top` up to 90 seconds to clear old containers from the list / count.

