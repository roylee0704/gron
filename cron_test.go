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

		// wraps around hours
		{"9:05, 8:04, 7:03", "7:03, 8:04, 9:05"},
		{"23:05, 20:04, 1:03", "1:03, 20:04, 23:05"},

		// wraps around seconds
		{"9:05:10, 8:04:20, 7:03:30", "7:03:30, 8:04:20, 9:05:10"},
		{"23:05:03, 20:04:01, 1:03:30", "1:03:30, 20:04:01, 23:05:03"},
		{"00:00:03, 00:00:01, 00:00:30", "00:00:01, 00:00:03, 00:00:30"},

		// wraps around days
		{
			"Wed Jun 8 9:05 2016, Tue Jun 7 8:04 2016, Wed Jun 8 9:01 2016",
			"Tue Jun 7 8:04 2016, Wed Jun 8 9:01 2016, Wed Jun 8 9:05 2016",
		},

		// wraps around months
		{
			"Sun Jun 4 9:05 2016, Sun Feb 7 8:04 2016, Sun May 8 9:01 2016",
			"Sun Feb 7 8:04 2016, Sun May 8 9:01 2016, Sun Jun 4 9:05 2016",
		},
	}

	for _, test := range tests {

		got := mockEntries(getTimes(test.entries))
		sort.Sort(byTime(got))

		want := mockEntries(getTimes(test.want))

		if !reflect.DeepEqual(got, want) {
			t.Errorf("entries not properly sorted: (want) %v != %v (got)", want, toS(got))
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
