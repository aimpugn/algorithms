package main

import (
	"fmt"
)

func runningSum(nums []int) []int {
	ans := []int{}
	for idx, num := range nums {
		if idx == 0 {
			ans = append(ans, num)
			continue
		}
		ans = append(ans, ans[idx-1]+nums[idx])
	}

	return ans
}

func main() {
	var nums []int
	var output []int

	nums = []int{1, 2, 3, 4}
	output = runningSum(nums)
	fmt.Println(output)

	nums = []int{1, 1, 1, 1, 1}
	output = runningSum(nums)
	fmt.Println(output)
}
