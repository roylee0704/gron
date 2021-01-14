package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fred07/gron"
	"github.com/fred07/gron/xtime"
)

type printJob struct{ Msg string }

func (p printJob) Run(wg *sync.WaitGroup) {
	fmt.Println(p.Msg)
}

func main() {

	var (
		daily     = gron.Every(1 * xtime.Day)
		weekly    = gron.Every(1 * xtime.Week)
		monthly   = gron.Every(30 * xtime.Day)
		yearly    = gron.Every(365 * xtime.Day)
		purgeTask = func() { fmt.Println("purge unwanted records") }
		printFoo  = printJob{"Foo"}
		printBar  = printJob{"Bar"}
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := gron.New()

	c.AddFunc(gron.Every(1*time.Hour), func() {
		fmt.Println("Every 1 hour")
	})
	c.Start(ctx)

	c.AddFunc(weekly, func() { fmt.Println("Every week") })
	c.Add(daily.At("12:30"), printFoo)
	c.Start(ctx)

	// Jobs may also be added to a running Cron
	c.Add(monthly, printBar)
	c.AddFunc(yearly, purgeTask)

	// Stop the scheduler (does not stop any jobs already running).
	defer c.Stop()
}
