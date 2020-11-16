package main

import (
	"fmt"
	"math"
)

/*
candies
extraCandies
2 <= candies.length <= 100
1 <= candies[i] <= 100
1 <= extraCandies <= 50

*/
func kidsWithCandies(candies []int, extraCandies int) []bool {
	candiesLen := len(candies)
	ans := make([]bool, candiesLen, candiesLen)
	for x := 0; x < candiesLen; x++ {
		ans[x] = true
	}

	for i, candy := range candies {
		for j := 0; j < candiesLen; j++ {
			if i == j {
				continue
			}
			if candy+extraCandies < candies[j] {
				ans[i] = false
				break
			}
		}
	}

	return ans
}

func kidsWithCandies2(candies []int, extraCandies int) []bool {
	candiesLen := len(candies)
	candyMax := 0
	for _, candy := range candies {
		candyMax = int(math.Max(float64(candy), float64(candyMax)))
	}

	ans := make([]bool, candiesLen, candiesLen)
	for i, candy := range candies {
		ans[i] = candy+extraCandies >= candyMax
	}

	return ans
}

func main() {
	var candies []int
	var ans []bool
	var extraCandies int

	candies = []int{2, 3, 5, 1, 3}
	extraCandies = 3
	ans = kidsWithCandies2(candies, extraCandies)
	fmt.Println(ans)

	candies = []int{4, 2, 1, 1, 2}
	extraCandies = 1
	ans = kidsWithCandies2(candies, extraCandies)
	fmt.Println(ans)

	candies = []int{12, 1, 12}
	extraCandies = 10
	ans = kidsWithCandies2(candies, extraCandies)
	fmt.Println(ans)

	candies = []int{1, 3, 2, 10, 7, 3, 4, 8}
	extraCandies = 8
	ans = kidsWithCandies2(candies, extraCandies)
	fmt.Println(ans)

}
