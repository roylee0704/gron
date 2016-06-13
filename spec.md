## Design Goal

- No delay job. Next schedule for a particular job shall not be delayed for its processing time. For example: Every(1).Minute() ensures tasks run at 10:01, 10:02, 10:03, but not 10:01+processing.
- Every(second,minute,hour,day,week).At(hh:mm)

## Design Specs

### CRON
An ADT that maintains a sequence of chronologically sorted entries. Cron keeps track of any number of entries, invoking the associated func as specified by the schedule. It may also be started, stopped and the entries may be inspected.

- **Start()**. Signals 'start' to get cron instant up & running.
- **Add(schedule, job)**. Signals `add` to add entry(schedule, job) to cron instant.
- **Stop()**. Signals `stop` to halt cron instant's processing. Bear in mind that (child go-routines) running jobs will be halted as well.
- **Clear()**. Clear all entries from queue.
- **run()**. Core functionality, run indefinitely(go-routine), forking out child go-routine: one for each job, multiplexing different channels/signals.

### SCHEDULE
An interface which wraps `Next(time)` method.
- **Next(time)**. Deduces next occurring schedule w.r.t time instant t.

### JOB
An interface which wraps `Run` method.
- **Run()**. To execute the underlying func.

### ENTRY
An ADT consists of a schedule and job to be run on that schedule. Also known as wrapper to `schedule` and `job`.

It keep tracks on the following states: `schedule`, `job`, `next`, `prev`.

## Detailed Design

### Cron: run() -- **core**
1. Sort entries chronologically.
2. Earliest entry be taken as the next triggering point.
3. Multiplexing of blocking channels/signals, that includes:
   - `ready`. earliest entry is ready to be run, in which subsequent entries will be measured, for which time is up and ready, it will be run as well.
   - `add`. add to entries.
   - `stop`.
4. Repeat 1. until `stop` is signaled.

### Entry: Sorted entry-slice
Entry slice implements sort package, chronologically.


### Schedule: periodicSchedule
A periodic schedule which occurs periodically: `t + period`
- **Every(period)**. Returns a periodicSchedule instant.


## Testing
Given a queue of scheduled entries. Test that:
- order of execution. add entries in reverse order, after one run, expects all entries are sorted accordingly.
- execution. ensure all entries are being executed.
- test concurrency adding.
