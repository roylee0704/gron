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

- **Run()**. Core functionality, run indefinitely(go-routine), multiplex different channels/signals.
- **Add(time, job)**. Signals `add` to add entry to cron.
- **Stop()**. Signals `stop` to halt cron processing.
- **Clear()**. Clear all entries from queue.

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
