package gron

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

// Test that invoking stop() before start() silently returns,
// without blocking the stop channel
func TestStopWithoutStart(t *testing.T) {
	cron := New()
	cron.Stop()
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

func mockEntries(nexts []time.Time) []Entry {
	var entries []Entry

	for _, n := range nexts {
		entries = append(entries, Entry{Next: n})
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
type toS []Entry

func (entries toS) String() string {
	var ret string
	for _, e := range entries {
		ret += fmt.Sprintf("[%v] ", e.Next)
	}
	return ret
}
