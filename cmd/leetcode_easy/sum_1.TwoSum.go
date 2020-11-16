package main

import (
	"fmt"
)

func twoSum(nums []int, target int) []int {
	var ans []int
	if numsLen := len(nums); numsLen > 0 {
		for i := 0; i < numsLen; i++ {
			for j := i + 1; j < numsLen; j++ {
				if sum := nums[i] + nums[j]; sum == target {
					ans = append(ans, i)
					ans = append(ans, j)
					break
				}
			}
		}
	}

	return ans
}

func main() {
	var ans []int
	var arr []int
	var target int
	arr = []int{2, 7, 11, 15}
	target = 9
	ans = twoSum(arr, target)
	fmt.Println(ans)
	arr = []int{3, 2, 4}
	target = 6
	ans = twoSum(arr, target)
	fmt.Println(ans)
	arr = []int{3, 3}
	target = 6
	ans = twoSum(arr, target)
	fmt.Println(ans)
}
