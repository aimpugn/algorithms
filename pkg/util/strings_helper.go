package util

import (
	"strconv"
	"strings"
)

func IntSliceJoin(intSlice []int, delimiter string) string {
	intSliceLen := len(intSlice)
	stringSlice := make([]string, 0, intSliceLen)
	for _, number := range intSlice {
		stringSlice = append(stringSlice, strconv.Itoa(number))
	}

	return strings.Join(stringSlice, delimiter)
}
