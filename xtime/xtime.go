package xtime

import "time"

// extends time.Duration

const (
	//Second has 1 * 1e9 nanoseconds
	Second time.Duration = time.Second
	//Minute has 60 seconds
	Minute time.Duration = time.Minute
	//Hour has 60 minutes
	Hour time.Duration = time.Hour
	//Day has 24 hours
	Day time.Duration = time.Hour * 24
	//Week has 7 days
	Week time.Duration = Day * 7
)
