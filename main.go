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
	n := time.Now().Nanosecond()
	fmt.Printf("%d, %d, %d\n", 1*time.Second, n, time.Duration(n)*time.Nanosecond)
	fmt.Print(time.Now(), time.Now())
}
