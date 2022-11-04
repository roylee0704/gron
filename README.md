# gron
[![Build Status](https://semaphoreci.com/api/v1/roylee0704/gron/branches/master/badge.svg)](https://semaphoreci.com/roylee0704/gron)
[![Go Report Card](https://goreportcard.com/badge/github.com/wifecooky/gron)](https://goreportcard.com/report/github.com/wifecooky/gron)
[![GoDoc](https://godoc.org/github.com/wifecooky/gron?status.svg)](https://godoc.org/github.com/wifecooky/gron)

Gron provides a clear syntax for writing and deploying cron jobs.

## Goals

- Minimalist APIs for scheduling jobs.
- Thread safety.
- Customizable Job Type.
- Customizable Schedule.

## Installation

```sh
$ go get github.com/wifecooky/gron
```

## Usage
Create `schedule.go`

```go
package main

import (
	"fmt"
	"time"
	"github.com/wifecooky/gron"
)

func main() {
	c := gron.New()
	c.AddFunc(gron.Every(1*time.Hour), func() {
		fmt.Println("runs every hour.")
	})
	c.Start()
}
```

#### Schedule Parameters

All scheduling is done in the machine's local time zone (as provided by the Go [time package](http://www.golang.org/pkg/time)).


Setup basic periodic schedule with `gron.Every()`.

```go
gron.Every(1*time.Second)
gron.Every(1*time.Minute)
gron.Every(1*time.Hour)
```

Also support `Day`, `Week` by importing `gron/xtime`:
```go
import "github.com/wifecooky/gron/xtime"

gron.Every(1 * xtime.Day)
gron.Every(1 * xtime.Week)
```

Schedule to run at specific time with `.At(hh:mm)`
```go
gron.Every(30 * xtime.Day).At("00:00")
gron.Every(1 * xtime.Week).At("23:59")
```

#### Custom Job Type
You may define custom job types by implementing `gron.Job` interface: `Run()`.

For example:

```go
type Reminder struct {
	Msg string
}

func (r Reminder) Run() {
  fmt.Println(r.Msg)
}
```

After job has defined, instantiate it and schedule to run in Gron.
```go
c := gron.New()
r := Reminder{ "Feed the baby!" }
c.Add(gron.Every(8*time.Hour), r)
c.Start()
```

#### Custom Job Func
You may register `Funcs` to be executed on a given schedule. Gron will run them in their own goroutines, asynchronously.

```go
c := gron.New()
c.AddFunc(gron.Every(1*time.Second), func() {
	fmt.Println("runs every second")
})
c.Start()
```


#### Custom Schedule
Schedule is the interface that wraps the basic `Next` method: `Next(p time.Duration) time.Time`

In `gron`, the interface value `Schedule` has the following concrete types:

- **periodicSchedule**. adds time instant t to underlying period p.
- **atSchedule**. reoccurs every period p, at time components(hh:mm).

For more info, checkout `schedule.go`.

### Full Example

```go
package main

import (
	"fmt"
	"github.com/wifecooky/gron"
	"github.com/wifecooky/gron/xtime"
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
