package util

import (
	"fmt"
	"strconv"
	"time"
)

var start time.Time
var index = 0

func StartExecution() {
	start = time.Now()
	index++
	PrintTime(start)
}

func PrintTime(time time.Time) {
	y, mon, d := time.Date()
	h, min, s := time.Clock()
	n := time.Nanosecond()
	timeString := strconv.Itoa(y) +
		"/" + mon.String() +
		"/" + strconv.Itoa(d) +
		" " + strconv.Itoa(h) +
		":" + strconv.Itoa(min) +
		":" + strconv.Itoa(s) +
		" " + strconv.Itoa(n) + "ns"
	fmt.Printf("%s\n", timeString)
}

func PrintElapsedTime(nameToMeasure string) {
	if len(nameToMeasure) == 0 {
		nameToMeasure = strconv.Itoa(index)
	}
	duration := time.Since(start)
	fmt.Printf("[%s] executed in %fsec, in %dms\n", nameToMeasure, duration.Seconds(), duration.Nanoseconds())
	if start.IsZero() {
		fmt.Println("start is zero value", start)
	}
}
