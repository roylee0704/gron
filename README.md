# gron
[![Build Status](https://semaphoreci.com/api/v1/roylee0704/gron/branches/master/badge.svg)](https://semaphoreci.com/roylee0704/gron)
[![Go Report Card](https://goreportcard.com/badge/github.com/roylee0704/gron)](https://goreportcard.com/report/github.com/roylee0704/gron)
[![GoDoc](https://godoc.org/github.com/roylee0704/gron?status.svg)](https://godoc.org/github.com/roylee0704/gron)

gron, Cron Jobs in Go. Gron provides a clear syntax for writing and deploying cron jobs.

# Features

- A simple to use API for scheduling jobs.
- Minimalist and no external dependencies.
- Thread safety.
- Excellent test coverage.
- Tested on Golang >= 1.5

# Quick Start
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
)

func main() {
	c := gron.New()
	c.AddFunc(gron.Every(30*time.Minute), func() { fmt.Print("Every half and hour") })
}
```
