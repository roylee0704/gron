# gron
gron, Cron Jobs in Go.


## Design Goal

- No delay job. Next schedule for a particular job shall not be delayed for its processing time. For example: Every(1).Minute() ensures tasks run at 10:01, 10:02, 10:03, but not 10:01+processing.
- Every(second,minute,hour,day,week).At(hh:mm)

## Design Specs

### CRON
An ADT that maintains a queue of entries/jobs, sorted by time (earliest).
Cron keeps track of any number of entries, invoking the associated func as
specified by the schedule. It may also be started, stopped and the entries
may be inspected.

- **Start()**. Signals 'start' to get cron instant up & running.
- **Add(time, job)**. Signals `add` to add entry to cron instant.
- **Stop()**. Signals `stop` to halt cron instant's processing.
- **Clear()**. Clear all entries from queue.
- **run()**. Core functionality, run indefinitely(go-routine), multiplexing different channels/signals.


#### Algorithm for run()

1. Sort job entries chronologically.
2. Earliest entry be taken as the next triggering point.
3. Multiplexing of blocking channels/signals, that includes:
   - `ready`. earliest entry is ready to be run, in which subsequent entries will be measured, for which time is up and ready, it will be run as well.
   - `add`. add to entries.
   - `stop`.
4. Repeat 1. until `stop` is signaled.


### Entry
An ADT that keep tracks of the following states: `next`, `prev`, and `job`.

- **Schedule(time)**. To schedule next run, referenced from input `time`.
- **Run()**. To run the given job (go-routine), recoverable.

### Schedule
An interface which wraps `Next(time)` method.
- **Next(time)**. Deduces next occurring event w.r.t time instant t.

### periodicSchedule
A periodic schedule which occurs periodically. `t + period`
- **Every(period)**. Returns a periodicSchedule instant.
