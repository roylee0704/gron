package gron

import "time"

// Schedule is the interface that wraps the basic Next method.
//
// Next deduces next occuring time based on t and underlying states.
type Schedule interface {
	Next(t time.Time) time.Time
}

// Every returns a Schedule that reoccurs in every period p, ranges from
// xtime.Second - xtime.Week
func Every(period time.Duration) Schedule {

	if period < time.Second {
		period = time.Second
	}

	return &periodicSchedule{
		period: period,
	}
}

type periodicSchedule struct {
	period time.Duration
}

func (ps periodicSchedule) Next(t time.Time) time.Time {

	return time.Now()
}
