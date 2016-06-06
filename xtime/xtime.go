package xtime

import "time"

// extends time.Duration

const (
	//Second is a
	Second time.Duration = time.Second
	//Minute is a
	Minute time.Duration = time.Minute
	//Hour is a
	Hour time.Duration = time.Hour
	//Day is a
	Day time.Duration = time.Hour * 24
	//Week is a
	Week time.Duration = Day * 7
)
