package main

import (
	"fmt"
	"time"

	"github.com/roylee0704/gron"
	"github.com/roylee0704/gron/xtime"
)

func main() {
	fmt.Println(gron.Every(30 * xtime.Day).Next(time.Now()))
	fmt.Println(gron.Every(3 * xtime.Week).At("10:00").Next(time.Now()))
}
