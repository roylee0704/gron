package gron

// Cron provides a convenient interface for scheduling job such as to clean-up
// database entry every month.
//
// Cron keeps track of any number of entries, invoking the associated func as
// specified by the schedule. It may also be started, stopped and the entries
// may be inspected.
type Cron struct {
	running bool
	stop    chan struct{}
}

// New instantiates new Cron instant c.
func New() *Cron {
	return &Cron{
		stop: make(chan struct{}),
	}
}

// Start signals cron instant c to get up and running.
func (c *Cron) Start() {
	c.running = true
	go c.run()
}

// run the scheduler...
//
// It needs to be private as it's responsible of synchronizing a critical
// shared state: `running`.
func (c *Cron) run() {

	for {
		select {
		case <-c.stop:
			return // terminate go-routine.
		}
	}
}

// Stop halts cron instant c from running.
func (c *Cron) Stop() {

	if !c.running {
		return
	}

	c.running = false
	c.stop <- struct{}{}
}
