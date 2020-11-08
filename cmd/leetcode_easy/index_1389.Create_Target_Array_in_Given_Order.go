package main

/*
https://leetcode.com/problems/create-target-array-in-the-given-order/
*/

import (
	"fmt"
)

/*
규칙
- Initially target array is empty.
- From left to right read nums[i] and index[i], insert at index index[i] the value nums[i] in target array
- Repeat the previous step until there are no elements to read in nums and index.
*/

func createTargetArray(nums []int, index []int) []int {
	target := make([]int, 0, len(index))
	for idx, value := range index {
		if len(target) > value {
			// exist
			/*
				slice는 포인터가 가리키는 배열을 바라본다
				target := []int{1,2,3,4,5}
				former := target[:3] // {1,2,3}, 정확히는 target{1,2,3} 부분을 가리킨다
				latter := target[3:] // {4,5}, 정확히는 target{4,5} 부분을 가리킨다
				former = append(former, 10) // target{1,2,3,10}이 되지만, 사실 target[3] == 4였는데, 이 4가 10으로 치환된다
											// latter{10,5}
											// target{1,2,3,10,5} 가 된다
			*/
			former := make([]int, len(target[:value]))
			copy(former, target[:value])
			latter := target[value:]
			former = append(former, nums[idx])
			target = append(former, latter...) // ...으로 unpack slice

		} else {
			// not exist
			target = append(target, nums[idx])
		}
	}

	return target
}

func main() {
	var nums []int
	var index []int
	var target []int

	nums = make([]int, 0, 5)
	index = make([]int, 0, 5)
	nums = append(nums, 0, 1, 2, 3, 4)
	index = append(index, 0, 1, 2, 2, 1)
	target = createTargetArray(nums, index)
	fmt.Println(target)

	nums = make([]int, 0, 5)
	index = make([]int, 0, 5)
	nums = append(nums, 1, 2, 3, 4, 0)
	index = append(index, 0, 1, 2, 3, 0)
	target = createTargetArray(nums, index)
	fmt.Println(target)

	nums = make([]int, 0, 5)
	index = make([]int, 0, 5)
	nums = append(nums, 1, 2, 3, 4, 5)
	index = append(index, 0, 0, 1, 3, 1)
	target = createTargetArray(nums, index)
	fmt.Println(target)

}
