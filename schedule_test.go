package gron

import (
	"testing"
	"time"

	"github.com/roylee0704/gron/xtime"
)

func TestPeriodicNext(t *testing.T) {

	tests := []struct {
		time   string
		period time.Duration
		want   string
	}{

		// Round to nearest second on the period
		{
			time:   "Mon Jun 9 12:45 2016",
			period: 15*xtime.Minute + 100*time.Nanosecond,
			want:   "Mon Jun 9 13:00 2016",
		},
	}

	for _, test := range tests {

		got := Every(test.period).Next(getTime(test.time))
		want := getTime(test.want)

		if got != want {
			t.Errorf("%s, \"%s\": (want) %v != %v (got)", test.time, test.period, want, got)
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
			t, err = time.Parse("2006-01-02T15:04:05-0700", value)
			if err != nil {
				panic(err)
			}
			// Daylight savings time tests require location
			if ny, err := time.LoadLocation("UTC"); err == nil {
				t = t.In(ny)
			}
		}
	}

	return t
}
