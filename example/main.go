package main

import (
	"fmt"
	"time"

	"github.com/fakeyanss/gron"
	"github.com/fakeyanss/gron/xtime"
)

type printJob struct{ Msg string }

func (p printJob) Run() {
	fmt.Println(p.Msg)
}

type canceledJob struct { // implements of JobWithCancel interface
	id string
}

func (j *canceledJob) JobID() string {
	return j.id
}

func (j *canceledJob) Run() {
	fmt.Printf("job %s run\n", j.id)
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

	c := gron.New()

	c.AddFunc(gron.Every(1*time.Hour), func() {
		fmt.Println("Every 1 hour")
	})
	c.Start()

	c.AddFunc(weekly, func() { fmt.Println("Every week") })
	c.Add(daily.At("12:30"), printFoo)
	c.Start()

	// Jobs may also be added to a running Cron
	c.Add(monthly, printBar)
	c.AddFunc(yearly, purgeTask)

	c.AddFuncWithJobID(gron.Every(1*time.Second), "job-id-1", func() {
		fmt.Println("job-id-1 runs every second")
	})
	c.AddCancelingJob(gron.Every(1*time.Second), &canceledJob{id: "job-id-2"})
	c.Start()

	time.Sleep(5 * time.Second)
	c.Cancel("job-id-1")
	time.Sleep(5 * time.Second)

	// Stop the scheduler (does not stop any jobs already running).
	defer c.Stop()
}
