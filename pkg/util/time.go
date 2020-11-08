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
}

func PrintElapsedTime() {
	if !start.IsZero() {
		fmt.Println(strconv.Itoa(index)+". executed during, ", time.Since(start))
	} else {
		fmt.Println("start is zero value", start)
	}
}
