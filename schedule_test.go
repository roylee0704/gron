package gron

import (
	"errors"
	"testing"
	"time"

	"github.com/wifecooky/gron/xtime"
)

func TestPeriodicAtNext(t *testing.T) {

	tests := []struct {
		time   string
		period time.Duration
		hhss   string
		want   string
	}{
		// simple cases
		{"Mon Jun 6 20:00 2016", 1 * xtime.Day, "23:00", "Tue Jun 6 23:00 2016"},
		{"Mon Jun 6 23:01 2016", 1 * xtime.Day, "23:00", "Tue Jun 7 23:00 2016"},
	}

	for _, test := range tests {
		got := Every(test.period).At(test.hhss).Next(getTime(test.time))
		want := getTime(test.want)
		if got != want {
			t.Errorf("(want) %v != %v (got)", want, got)
		}
	}
}

func TestPeriodicNext(t *testing.T) {

	tests := []struct {
		time   string
		period time.Duration
		want   string
	}{

		//Wraps around days
		{
			time:   "Mon Jun 6 23:46 2016",
			period: 15 * time.Minute,
			want:   "Mon Jun 7 00:01 2016",
		},
		{
			time:   "Mon Jun 6 23:59:40 2016",
			period: 21 * time.Second,
			want:   "Mon Jun 7 00:00:01 2016",
		},
		{
			time:   "Mon Jun 6 23:30:19 2016",
			period: 29*time.Minute + 41*time.Second,
			want:   "Mon Jun 7 00:00:00 2016",
		},
		{
			time:   "Mon Jun 6 23:46:20 2016",
			period: 15*time.Minute + 40*time.Second,
			want:   "Mon Jun 7 00:02:00 2016",
		},

		// Wrap around months
		{
			time:   "Mon Jun 6 16:49 2016",
			period: 30 * xtime.Day, // adds 30 days, equates to #days in June.
			want:   "Mon Jul 6 16:49 2016",
		},

		// Wrap around minute, hour, day, month, and year
		{

			time:   "Sat Dec 31 23:59:59 2016",
			period: 1 * time.Second,
			want:   "Sun Jan 1 00:00:00 2017",
		},

		// Round to nearest second on the period
		{
			time:   "Mon Jun 6 12:45 2016",
			period: 15*time.Minute + 100*time.Nanosecond,
			want:   "Mon Jun 6 13:00 2016",
		},

		// Round to 1 second if the duration is less
		{
			time:   "Mon Jun 6 17:38:01 2016",
			period: 59 * time.Millisecond,
			want:   "Mon Jun 6 17:38:02 2016",
		},

		// Round to nearest second when calculating the next time.
		{"Mon Jun 6 17:38:01.009 2016", 15 * time.Minute, "Mon Jun 6 17:53:01 2016"},

		// Round to nearest second for both
		{"Mon Jun 6 17:38:01.009 2016", 15*time.Minute + 50*time.Nanosecond, "Mon Jun 6 17:53:01 2016"},
	}

	for i, test := range tests {

		got := Every(test.period).Next(getTime(test.time))
		want := getTime(test.want)

		if got != want {
			t.Errorf("case[%d], %s, \"%s\": (want) %v != %v (got)", i, test.time, test.period, want, got)
		}
	}

}

func TestNaiveParse(t *testing.T) {

	tests := []struct {
		hhmm   string
		hh, mm int
		err    error
	}{
		// simple cases
		{"10:11", 10, 11, nil},
		{"24:44", 24, 44, nil},

		// single digit parsing
		{"01:44", 1, 44, nil},
		{"01:06", 1, 6, nil},

		// edge test
		{"25:11", 0, 0, errors.New("invalid hh format")},
		{"25:70", 0, 0, errors.New("invalid hh format")},
		{"24:70", 0, 0, errors.New("invalid mm format")},
	}

	for _, test := range tests {
		h, m, err := parse(test.hhmm)
		if h != test.hh || m != test.mm {
			t.Errorf("invalid input (%s):  (want) %d:%d != %d:%d (got)", test.hhmm, test.hh, test.mm, h, m)
		}
		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("failed to catch formatting error: %s != %s", test.err, err)
		}
	}

}

func TestAtScheduleReset(t *testing.T) {
	tests := []struct {
		time string
		hhmm string
		want string
	}{
		// simple cases
		{"Mon Jun 6 22:13 2016", "10:11", "Mon Jun 6 10:11 2016"},

		// reset seconds
		{"Mon Jun 6 22:13:03 2016", "12:11", "Mon Jun 6 12:11:00 2016"},
		{"Mon Jun 6 22:13:03.009 2016", "12:11", "Mon Jun 6 12:11:00 2016"},
	}
	for _, test := range tests {
		hh, mm, _ := parse(test.hhmm)
		as := atSchedule{
			hh: hh,
			mm: mm,
		}

		got := as.reset(getTime(test.time))
		want := getTime(test.want)

		if got != want {
			t.Errorf(" (want) %v != %v (got)", want, got)
		}
	}
}

func getTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	t, err := time.Parse("Mon Jan 2 15:04 2006", value)
	if err != nil {
		t, err = time.Parse("Mon Jan 2 15:04:05 2006", value)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04", value)
			if err != nil {
				t, err = time.Parse("2006-01-02 15:04:05", value)
				if err != nil {
					t, err = time.Parse("15:04", value)
					if err != nil {
						t, err = time.Parse("15:04:05", value)
						if err != nil {
							t, err = time.Parse("2006-01-02T15:04:05-0700", value)
							if err != nil {
								panic(err)
							}
							if loc, err := time.LoadLocation("UTC"); err == nil {
								t = t.In(loc)
							}
						}
					}
				}
			}
		}
	}

	return t
}
