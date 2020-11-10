package util

import (
	"math/rand"
	"sort"
	"time"
)

const (
	SortDescending = iota
	SortAscending
)

func RandRange(min int, max int, sizeMax int, sortType int) ([]int, int) {
	randMap := make(map[int]bool, sizeMax)
	randSlice := make([]int, 0, sizeMax)
	realLen := 0
	rand.Seed(time.Now().UnixNano())

loop:
	for i := 0; i < sizeMax; i++ {
		randValue := rand.Intn(max-min) + min
		if (sizeMax == realLen) || ((max - min) == realLen) {
			break loop
		}
		if _, exist := randMap[randValue]; exist {
			i--
			continue
		}
		randMap[randValue] = true
		randSlice = append(randSlice, randValue)
		realLen++
	}

	switch sortType {
	case SortDescending:
		sort.Sort(sort.Reverse(sort.IntSlice(randSlice)))
	case SortAscending:
		sort.Ints(randSlice)
	}

	return randSlice, realLen
}
