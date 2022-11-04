package gron

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/wifecooky/gron/xtime"
)

// Most test jobs scheduled to run at 1 second mark.
// Test expects to fail after OneSecond: 1.01 seconds.
const OneSecond = 1*time.Second + 10*time.Millisecond

// start cron, stop cron successfully.
func TestNoEntries(t *testing.T) {
	cron := New()
	cron.Start()

	select {
	case <-time.After(OneSecond):
		t.FailNow()
	case <-stop(cron):
	}
}

// hijack time.After, start the cron, simulate a delay.
// empty entries shall not trigger to run next job in next second.
func TestNoPhantomJobs(t *testing.T) {
	entry := 0
	// overriding internal state func
	after = func(d time.Duration) <-chan time.Time {
		entry++
		return time.After(d)
	}
	defer func() {
		after = time.After
	}() // proper tear down

	cron := New()
	cron.Start()
	defer cron.Stop()

	time.Sleep(1 * time.Millisecond) // simulate a delay

	if entry > 1 {
		t.Errorf("phantom job had run %d time(s).", entry)
	}
}

// add a job, start a cron, expect it runs
func TestAddBeforeRun(t *testing.T) {
	done := make(chan struct{})
	cron := New()
	cron.AddFunc(Every(1*time.Second), func() { done <- struct{}{} })
	cron.Start()
	defer cron.Stop()

	select {
	case <-time.After(OneSecond):
		t.FailNow()
	case <-done:
	}
}

// start a cron, add a job, expect it runs
func TestAddWhileRun(t *testing.T) {
	done := make(chan struct{})
	cron := New()
	cron.Start()
	defer cron.Stop()

	cron.AddFunc(Every(1*time.Second), func() { done <- struct{}{} })

	select {
	case <-time.After(OneSecond):
		t.FailNow()
	case <-done:
	}
}

// Test that invoking stop() before start() silently returns,
// without blocking the stop channel
func TestStopWithoutStart(t *testing.T) {
	cron := New()
	cron.Stop()
}

// start cron, stop cron, add a job, verify job shouldn't run.
func TestJobsDontRunAfterStop(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	cron := New()
	cron.Start()
	cron.Stop()
	cron.AddFunc(Every(1*time.Second), func() { wg.Done() })

	select {
	case <-time.After(OneSecond):
		// no job has run
	case <-wait(wg):
		t.FailNow()

	}
}

// Test that entries are sorted correctly.
// Adds an immediate entry, make sure it runs immediately.
// Subsequent entries are checked and run at same instant, iff possessed
// same schedule as first entry.
func TestMultipleEntries(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(3)
	cron := New()

	cron.AddFunc(Every(5*time.Minute), func() {})
	cron.AddFunc(Every(1*time.Second), func() { wg.Done() })
	cron.AddFunc(Every(1*time.Second), func() { wg.Done() })
	cron.AddFunc(Every(1*time.Second), func() { wg.Done() })
	cron.AddFunc(Every(4*xtime.Week), func() {})

	cron.Start()
	defer cron.Stop()

	select {
	case <-time.After(OneSecond):
		t.FailNow()
	case <-wait(wg):
	}
}

// Test that job runs n times after (n * p) period.
func TestRunJobTwice(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	cron := New()
	cron.Start()
	defer cron.Stop()

	cron.AddFunc(Every(1*xtime.Day), func() { wg.Done() })
	cron.AddFunc(Every(1*xtime.Week), func() { wg.Done() })
	cron.AddFunc(Every(1*time.Second), func() { wg.Done() })

	select {
	case <-time.After(2 * OneSecond):
		t.FailNow()
	case <-wait(wg):
	}
}

// arbitrary job struct, with god's view enabled.
type arbitraryJob struct {
	wg *sync.WaitGroup // god's update
	id string          // god's record
}

// implements runnable
func (j arbitraryJob) Run() {
	if j.wg != nil {
		j.wg.Done()
	}
}

// simple test on job types
func TestJobImplementer(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	cron := New()
	cron.Add(Every(1*xtime.Day), arbitraryJob{wg, "job-1"}) // merely distraction
	cron.Add(Every(1*time.Second), arbitraryJob{wg, "job-2"})
	cron.Add(Every(1*xtime.Week), arbitraryJob{wg, "job-3"}) // merely distraction

	cron.Start()
	defer cron.Stop()

	select {
	case <-time.After(2 * OneSecond):
		t.FailNow()
	case <-wait(wg):
	}
}

// Test that entries are in correct sequence after n run.
func TestEntryOrdering(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	cron := New()

	cron.Add(Every(40*time.Second), arbitraryJob{wg, "job-1"})
	cron.Add(Every(1*time.Second), arbitraryJob{wg, "job-2"})
	cron.Add(Every(1*xtime.Week), arbitraryJob{wg, "job-3"})
	cron.Add(Every(12*time.Second), arbitraryJob{wg, "job-4"})
	cron.Add(Every(8*time.Second), arbitraryJob{wg, "job-5"})

	cron.Start()
	select {
	case <-time.After(2 * OneSecond):
		t.FailNow()
	case <-wait(wg):
	}
	cron.Stop()

	want := []string{"job-2", "job-5", "job-4", "job-1", "job-3"}
	var got []string
	for _, e := range cron.Entries() {
		got = append(got, e.Job.(arbitraryJob).id)
	}

	for i := range got {
		if want[i] != got[i] {
			t.Errorf("incorrect sequencing: (want) %q != %q (got)", want, got)
			break
		}
	}
}

// Test that entries are chronologically sorted
func TestByTimeSort(t *testing.T) {
	tests := []struct {
		entries string
		want    string
	}{
		// simple cases
		{"10:05, 10:04, 10:03", "10:03, 10:04, 10:05"},
		{"10:05, 10:04, 10:03", "10:03, 10:04, 10:05"},

		// hours dominate
		{"7:00, 8:00, 9:00", "7:00, 8:00, 9:00"},
		{"9:00, 8:00, 7:00", "7:00, 8:00, 9:00"},
		{"9:00, 8:49, 8:09", "8:09, 8:49, 9:00"},

		// seconds dominate
		{"00:00:01, 00:00:03, 00:00:30", "00:00:01, 00:00:03, 00:00:30"},
		{"00:00:03, 00:00:01, 00:00:30", "00:00:01, 00:00:03, 00:00:30"},
		{"00:05:10, 00:04:20, 00:03:30", "00:03:30, 00:04:20, 00:05:10"},

		// days dominate
		{
			"Wed Jun 8 9:05 2016, Tue Jun 7 8:04 2016, Wed Jun 8 9:01 2016",
			"Tue Jun 7 8:04 2016, Wed Jun 8 9:01 2016, Wed Jun 8 9:05 2016",
		},

		// months dominate
		{
			"Sun Jun 4 9:05 2016, Sun Feb 7 8:04 2016, Sun May 8 9:01 2016",
			"Sun Feb 7 8:04 2016, Sun May 8 9:01 2016, Sun Jun 4 9:05 2016",
		},

		// zero hours sort as intended
		{"00:00, 00:00, 00:10", "00:00, 00:00, 00:10"},

		// zero minutes sort as intended
		{"00:00:00, 00:00:00, 00:00:10", "00:00:00, 00:00:00, 00:00:10"},

		// zero times (uninitialised time) should push to back of queue.
		{
			"2016-01-01 00:00, 2016-01-01 00:00, 0001-01-01 00:00",
			"2016-01-01 00:00, 2016-01-01 00:00, 0001-01-01 00:00",
		},
		{
			"2016-01-01 00:01, 2016-01-01 00:00, 0001-01-01 00:00",
			"2016-01-01 00:00, 2016-01-01 00:01, 0001-01-01 00:00",
		},
		{
			"0001-01-01 00:00, 2016-01-01 00:01, 2016-01-01 00:00",
			"2016-01-01 00:00, 2016-01-01 00:01, 0001-01-01 00:00",
		},
		{
			"2016-01-01 00:01, 0001-01-01 00:00, 2016-01-01 00:00",
			"2016-01-01 00:00, 2016-01-01 00:01, 0001-01-01 00:00",
		},
		{
			"0001-01-01 00:00, 0001-01-01 00:00, 2016-01-01 00:00",
			"2016-01-01 00:00, 0001-01-01 00:00, 0001-01-01 00:00",
		},
	}

	for i, test := range tests {

		got := mockEntries(getTimes(test.entries))
		sort.Sort(byTime(got))

		want := mockEntries(getTimes(test.want))

		if !reflect.DeepEqual(got, want) {
			t.Errorf("entries[%d] out of order: \n%v (want)\n%v (got)", i, toS(want), toS(got))
		}
	}
}

func mockEntries(nexts []time.Time) []*Entry {
	var entries []*Entry

	for _, n := range nexts {
		entries = append(entries, &Entry{Next: n})
	}
	return entries
}

// getTimes splits comma-separated time.
func getTimes(s string) []time.Time {

	ts := strings.Split(s, ",")
	ret := make([]time.Time, len(ts))

	for i, t := range ts {
		ret[i] = getTime(strings.Trim(t, " "))
	}
	return ret
}

// wrapper to stringify time instant t
type toS []*Entry

func (entries toS) String() string {
	var ret string
	for _, e := range entries {
		ret += fmt.Sprintf("[%v] ", e.Next)
	}
	return ret
}

// wait signals back when WaitGroup has wait().
func wait(wg *sync.WaitGroup) <-chan bool {
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()
	return done
}

// stop signals back when cron has stop()
func stop(c *Cron) <-chan bool {
	done := make(chan bool)
	go func() {
		c.Stop()
		done <- true
	}()
	return done
}
