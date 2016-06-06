package gron

import (
	"time"

	"github.com/roylee0704/gron/xtime"
)

// Schedule is the interface that wraps the basic Next method.
//
// Next deduces next occuring time based on t and underlying states.
type Schedule interface {
	Next(t time.Time) time.Time
}

// Every returns a Schedule reoccurs every period p, p must be at least
// time.Second.
func Every(p time.Duration) Schedule {

	if p < time.Second {
		p = xtime.Second
	}

	p = p - time.Duration(p.Nanoseconds())%time.Second // round-off time.seconds

	return &periodicSchedule{
		period: p,
	}
}

type periodicSchedule struct {
	period time.Duration
}

// Next adds time t to underlying period, truncates up to unit of seconds.
func (ps periodicSchedule) Next(t time.Time) time.Time {
	return t.Truncate(time.Second).Add(ps.period)
}
