package gron

import (
	"errors"
	"time"

	"github.com/wifecooky/gron/xtime"
)

// Schedule is the interface that wraps the basic Next method.
//
// Next deduces next occurring time based on t and underlying states.
type Schedule interface {
	Next(t time.Time) time.Time
}

// AtSchedule extends Schedule by enabling periodic-interval & time-specific setup
type AtSchedule interface {
	At(t string) Schedule
	Schedule
}

// Every returns a Schedule reoccurs every period p, p must be at least
// time.Second.
func Every(p time.Duration) AtSchedule {

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
// Note: At panics when period p is less than xtime.Day, and error hh:ss format.
func (ps periodicSchedule) At(t string) Schedule {
	if ps.period < xtime.Day {
		panic("period must be at least in days")
	}

	// parse t naively
	h, m, err := parse(t)

	if err != nil {
		panic(err.Error())
	}

	return &atSchedule{
		period: ps.period,
		hh:     h,
		mm:     m,
	}
}

// parse naively tokenises hours and minutes.
//
// returns error when input format was incorrect.
func parse(hhmm string) (hh int, mm int, err error) {

	hh = int(hhmm[0]-'0')*10 + int(hhmm[1]-'0')
	mm = int(hhmm[3]-'0')*10 + int(hhmm[4]-'0')

	if hh < 0 || hh > 24 {
		hh, mm = 0, 0
		err = errors.New("invalid hh format")
	}
	if mm < 0 || mm > 59 {
		hh, mm = 0, 0
		err = errors.New("invalid mm format")
	}

	return
}

type atSchedule struct {
	period time.Duration
	hh     int
	mm     int
}

// reset returns new Date based on time instant t, and reconfigure its hh:ss
// according to atSchedule's hh:ss.
func (as atSchedule) reset(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), as.hh, as.mm, 0, 0, time.UTC)
}

// Next returns **next** time.
// if t passed its supposed schedule: reset(t), returns reset(t) + period,
// else returns reset(t).
func (as atSchedule) Next(t time.Time) time.Time {
	next := as.reset(t)
	if t.After(next) {
		return next.Add(as.period)
	}
	return next
}
