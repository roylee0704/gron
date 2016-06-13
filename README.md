# gron
[![Build Status](https://semaphoreci.com/api/v1/roylee0704/gron/branches/master/badge.svg)](https://semaphoreci.com/roylee0704/gron)
[![Go Report Card](https://goreportcard.com/badge/github.com/roylee0704/gron)](https://goreportcard.com/report/github.com/roylee0704/gron)
[![GoDoc](https://godoc.org/github.com/roylee0704/gron?status.svg)](https://godoc.org/github.com/roylee0704/gron)

Gron provides a clear syntax for writing and deploying cron jobs.

## Installation

```sh
$ go get github.com/roylee0704/gron
```

## Hello, World!
Create `schedule.go`

```go
package main

import (
	"fmt"
	"time"
	"github.com/roylee0704/gron"
	"github.com/roylee0704/gron/xtime"
)

func main() {
	c := gron.New()
	c.AddFunc(gron.Every(3*time.Hour), func() { fmt.Print("Runs every 3 hour") })

	c.AddFunc(gron.Every(1*xtime.Day).At("04:30"), func() { fmt.Print("Runs at 4:30 in the morning")})
}
```

### Define your own job types
Gron currently ships with just 1 job type: runner. You can define your own by implementing `gron.Job` interface.


For example:

```go
type Party struct { funLevel string }

func (p *Party) Run() {
  fmt.Println(p.funLevel)
}```
