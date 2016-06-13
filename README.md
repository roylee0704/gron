# gron
[![Build Status](https://semaphoreci.com/api/v1/roylee0704/gron/branches/master/badge.svg)](https://semaphoreci.com/roylee0704/gron)
[![Go Report Card](https://goreportcard.com/badge/github.com/roylee0704/gron)](https://goreportcard.com/report/github.com/roylee0704/gron)
[![GoDoc](https://godoc.org/github.com/roylee0704/gron?status.svg)](https://godoc.org/github.com/roylee0704/gron)

Gron provides a clear syntax for writing and deploying cron jobs.

## Goals

- Minimalist APIs for scheduling jobs.
- Thread safety.
- Ability to define own job type.

## Installation

```sh
$ go get github.com/roylee0704/gron
```

## Usage
Create `schedule.go`

```go
package main

import (
	"fmt"
	"time"
	"github.com/roylee0704/gron"
)

func main() {
	c := gron.New()
	c.AddFunc(gron.Every(3 * time.Hour), func() { fmt.Print("Runs every 3 hour") })
}
```

### Define your own job types
Gron currently ships with just 1 job type: runner. You can define your own by implementing `gron.Job` interface.


For example:

```go
type Reminder struct { ... }

func (r *Reminder) Run() {
  r.send()
}

c := gron.New()
c.Add(gron.Every(5 * time.Hour), Reminder{ ... })
```
