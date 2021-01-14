package gron

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// Entry consists of a schedule and the job to be executed on that schedule.
type Entry struct {
	Schedule Schedule
	Job      Job

	// the next time the job will run. This is zero time if Cron has not been
	// started or invalid schedule.
	Next time.Time

	// the last time the job was run. This is zero time if the job has not been
	// run.
	Prev time.Time
}

// byTime is a handy wrapper to chronologically sort entries.
type byTime []*Entry

func (b byTime) Len() int      { return len(b) }
func (b byTime) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// Less reports `earliest` time i should sort before j.
// zero time is not `earliest` time.
func (b byTime) Less(i, j int) bool {

	if b[i].Next.IsZero() {
		return false
	}
	if b[j].Next.IsZero() {
		return true
	}

	return b[i].Next.Before(b[j].Next)
}

// Job is the interface that wraps the basic Run method.
//
// Run executes the underlying func.
type Job interface {
	Run(wg *sync.WaitGroup)
}

// Cron provides a convenient interface for scheduling job such as to clean-up
// database entry every month.
//
// Cron keeps track of any number of entries, invoking the associated func as
// specified by the schedule. It may also be started, stopped and the entries
// may be inspected.
type Cron struct {
	entries []*Entry
	running bool
	add     chan *Entry
	stop    chan struct{}
	wg      sync.WaitGroup
}

// New instantiates new Cron instant c.
func New() *Cron {
	return &Cron{
		stop: make(chan struct{}),
		add:  make(chan *Entry),
	}
}

// StartAndServe serve cron like a eternal process
func (c *Cron) StartAndServe(ctx context.Context) error {
	c.running = true
	return c.run(ctx)
}

// Start signals cron instant c to get up and running.
func (c *Cron) Start(ctx context.Context) {
	c.running = true
	go c.run(ctx)
}

// Add appends schedule, job to entries.
//
// if cron instant is not running, adding to entries is trivial.
// otherwise, to prevent data-race, adds through channel.
func (c *Cron) Add(s Schedule, j Job) {

	entry := &Entry{
		Schedule: s,
		Job:      j,
	}

	if !c.running {
		c.entries = append(c.entries, entry)
		return
	}
	c.add <- entry
}

// AddFunc registers the Job function for the given Schedule.
func (c *Cron) AddFunc(s Schedule, j func()) {
	c.Add(s, JobFunc(j))
}

// Stop halts cron instant c from running.
func (c *Cron) Stop() {

	if !c.running {
		return
	}
	c.running = false
	c.stop <- struct{}{}
}

// GracefullyStop halts cron after running jobs are finished
func (c *Cron) GracefullyStop() {

	c.stopAndWait()
	c.stop <- struct{}{}
}

func (c *Cron) stopAndWait() {
	if !c.running {
		return
	}
	c.running = false
	c.wg.Wait()
}

var after = time.After

// run the scheduler...
//
// It needs to be private as it's responsible of synchronizing a critical
// shared state: `running`.
func (c *Cron) run(ctx context.Context) error {

	var effective time.Time
	now := time.Now().Local()

	// to figure next trig time for entries, referenced from now
	for _, e := range c.entries {
		e.Next = e.Schedule.Next(now)
	}

	for {
		sort.Sort(byTime(c.entries))
		if len(c.entries) > 0 {
			effective = c.entries[0].Next
		} else {
			effective = now.AddDate(15, 0, 0) // to prevent phantom jobs.
		}

		select {
		case now = <-after(effective.Sub(now)):
			// entries with same time gets run.
			for _, entry := range c.entries {
				if entry.Next != effective {
					break
				}
				entry.Prev = now
				entry.Next = entry.Schedule.Next(now)
				c.wg.Add(1)
				go entry.Job.Run(&c.wg)
			}
		case sig := <-ctx.Done():
			c.stopAndWait()
			return fmt.Errorf("Interrupt by signal: %s", sig)
		case e := <-c.add:
			e.Next = e.Schedule.Next(time.Now())
			c.entries = append(c.entries, e)
		case <-c.stop:
			return errors.New("Cron stop") // terminate go-routine.
		}
	}
}

// Entries returns cron etn
func (c *Cron) Entries() []*Entry {
	return c.entries
}

// JobFunc is an adapter to allow the use of ordinary functions as gron.Job
// If f is a function with the appropriate signature, JobFunc(f) is a handler
// that calls f.
//
// todo: possibly func with params? maybe not needed.
type JobFunc func()

// Run calls j()
func (j JobFunc) Run(wg *sync.WaitGroup) {
	j()
	wg.Done()
}
