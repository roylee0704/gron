# gron

[![Build Status](https://travis-ci.com/Fred07/gron.svg?branch=master)](https://travis-ci.com/Fred07/gron)
[![Go Report Card](https://goreportcard.com/badge/github.com/fred07/gron)](https://goreportcard.com/report/github.com/fred07/gron)

Gron provides a clear syntax for writing and deploying cron jobs.

## Goals

- Minimalist APIs for scheduling jobs.
- Thread safety.
- Customizable Job Type.
- Customizable Schedule.

## Different to origin

Most features and interfaces are as same as origin repository, but also introduces one new feature to Cron, which `GracefullyStop()` make cron has ability to hold the process, stop creating new child goroutine, and ensure sub-job is finished before main goroutine close. This benefit to the scenarios that you want to handle OS signals, including `SIGINT`, `SIGTERM`. For example, after the interruption from signal `SIGINT`, stopping cron gracefully by calling `GracefullyStop()`, this could prevent bundle of processing job stop inappropriately.

Second feature is `StartAndServe()`, this feature makes gron blocking and keep go-routine alive and til it return an error, which that you don't have to create an ugly infinite loop. All you need to do is like `log.Fatal(gron.StartAndServe())`

And context is required now as a parameter when calling the `Start()` or `StartAndServe()` function, and gron would trigger familiar behavior just like `GracefullyStop()` while `<-ctx.Done()` is triggered.

### Summary of new features

- `GracefullyStop()`
- `StartAndServe()`

## Installation

```sh
$ go get github.com/fred07/gron
```

## Usage

Create `schedule.go`

```go
package main

import (
	"fmt"
	"time"
	"github.com/fred07/gron"
)

func main() {
	c := gron.New()
	c.AddFunc(gron.Every(1*time.Hour), func() {
		fmt.Println("runs every hour.")
	})
	c.Start()
}
```

### Schedule Parameters

All scheduling is done in the machine's local time zone (as provided by the Go [time package](http://www.golang.org/pkg/time)).

Setup basic periodic schedule with `gron.Every()`.

```go
gron.Every(1*time.Second)
gron.Every(1*time.Minute)
gron.Every(1*time.Hour)
```

Also support `Day`, `Week` by importing `gron/xtime`:

```go
import "github.com/fred07/gron/xtime"

gron.Every(1 * xtime.Day)
gron.Every(1 * xtime.Week)
```

Schedule to run at specific time with `.At(hh:mm)`

```go
gron.Every(30 * xtime.Day).At("00:00")
gron.Every(1 * xtime.Week).At("23:59")
```

### Custom Job Type

You may define custom job types by implementing `gron.Job` interface: `Run(wg *sync.WaitGroup)`.

For example:

```go
type Reminder struct {
	Msg string
}

func (r Reminder) Run(wg *sync.WaitGroup) {
  fmt.Println(r.Msg)
}
```

After job has defined, instantiate it and schedule to run in Gron.

```go
ctx, cancel := context.WithCancel(context.Background())
c := gron.New()
r := Reminder{ "Feed the baby!" }
c.Add(gron.Every(8*time.Hour), r)
c.Start(ctx)
```

### Custom Job Func

You may register `Funcs` to be executed on a given schedule. Gron will run them in their own goroutines, asynchronously.

```go
ctx, cancel := context.WithCancel(context.Background())
c := gron.New()
c.AddFunc(gron.Every(1*time.Second), func() {
	fmt.Println("runs every second")
})
c.Start(ctx)
```

### Custom Schedule

Schedule is the interface that wraps the basic `Next` method: `Next(p time.Duration) time.Time`

In `gron`, the interface value `Schedule` has the following concrete types:

- **periodicSchedule**. adds time instant t to underlying period p.
- **atSchedule**. reoccurs every period p, at time components(hh:mm).

For more info, checkout `schedule.go`.

### Serve like a daemon

In real case, you may need a infinite `for` loop to keep main goroutine alive, in fact, you can use `StartAndServe()` to do that.

`StartAndServe()` is the method to start cron and keeping process there like a server.

```go
ctx, cancel := context.WithCancel(context.Background())
c := gron.New()
c.AddFunc(gron.Every(1*time.Second), func() {
	fmt.Println("runs every second")
})
c.StartAndServe(ctx)
```

### Full Example

```go
package main

import (
	"fmt"
	"github.com/fred07/gron"
	"github.com/fred07/gron/xtime"
)

type PrintJob struct{ Msg string }

func (p PrintJob) Run() {
	fmt.Println(p.Msg)
}

func main() {

	var (
		// schedules
		daily     = gron.Every(1 * xtime.Day)
		weekly    = gron.Every(1 * xtime.Week)
		monthly   = gron.Every(30 * xtime.Day)
		yearly    = gron.Every(365 * xtime.Day)

		// contrived jobs
		purgeTask = func() { fmt.Println("purge aged records") }
		printFoo  = printJob{"Foo"}
		printBar  = printJob{"Bar"}
	)

	c := gron.New()

	c.Add(daily.At("12:30"), printFoo)
	c.AddFunc(weekly, func() { fmt.Println("Every week") })
	c.Start()

	// Jobs may also be added to a running Gron
	c.Add(monthly, printBar)
	c.AddFunc(yearly, purgeTask)

	// Stop Gron (running jobs are not halted).
	c.Stop()
}
```
