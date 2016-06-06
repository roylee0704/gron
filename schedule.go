package gron

import (
	"errors"
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

	p = p - time.Duration(p.Nanoseconds())%time.Second // truncates up to seconds

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

// At returns a schedule which reoccurs every period p, at time t(hh:ss).
//
// Note: At panics when period p is less than xtime.Day
func (ps periodicSchedule) At(t string) Schedule {
	if ps.period < xtime.Day {
		panic("period must be at least in days")
	}

	// parse t naively

	return &atSchedule{}
}

// parse naively tokenises hours and seconds.
//
// returns error when input format was incorrect.
func parse(hhss string) (hh int, ss int, err error) {

	hh = int(hhss[0]-'0')*10 + int(hhss[1]-'0')
	ss = int(hhss[3]-'0')*10 + int(hhss[4]-'0')

	if hh < 0 || hh > 24 {
		hh, ss = 0, 0
		err = errors.New("invalid hh format")
	}
	if ss < 0 || ss > 59 {
		hh, ss = 0, 0
		err = errors.New("invalid ss format")
	}

	return
}

type atSchedule struct {
	period time.Duration
	hh     int
	ss     int
}

func (as atSchedule) Next(t time.Time) time.Time {
	return time.Time{}
}
