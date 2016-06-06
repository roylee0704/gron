package main

import (
	"fmt"
	"time"
)

func main() {

	// lastRun := time.Now()
	//
	// fmt.Println(lastRun)
	// fmt.Println(lastRun.Add(1 * time.Second * 60 * 60 * 24 * 7))

	//var t = "10"
	//n := time.Now().Nanosecond()
	//fmt.Printf("%d, %d, %d\n", 1*time.Second, n, time.Duration(n)*time.Nanosecond)

	//	every(5 * time.Milisecond)

	d := 15*time.Minute + 50*time.Nanosecond

	fmt.Println(time.Duration(d.Nanoseconds())%time.Second, d.Nanoseconds())
	//fmt.Println(At(" sun 15:00").Every("day"))
}

// func every(period time.Duration) time.Time {
//
// 	fmt.Println(period - time.Duration(period.Nanoseconds())%time.Second)
// 	return time.Now().Truncate(time.Second).Add(period)
// }

type atDelay struct {
	k string
}

func At(time string) atDelay {
	return atDelay{k: time}
}

func (at atDelay) Every(delay time.Duration) string {
	return delay + at.k
}
