package main

import (
	"fmt"
	"math"

	"algorithms/pkg/util"
)

/*
Given a sorted (in ascending order) integer array nums of n elements and a target value,
write a function to search target in nums.
If target exists, then return its index, otherwise return -1.
*/

/*
You may assume that all elements in nums are unique.
n will be in the range [1, 10000].
The value of each element in nums will be in the range [-9999, 9999].
*/
/*
포인터(인덱스)를 어떻게 이동시킬 것인가
*/

var currentIdx int
var isLeftOrRight int

const (
	Left = iota
	Right
)

func BinarySearch(nums []int, target int) int {
	currentIdx = 0
	isLeftOrRight = -1

	return binarySearchBad(nums, target)
}

func binarySearchBad(nums []int, target int) int {
	ans := -1
	numsLen := len(nums)
	if numsLen == 0 {
		return ans
	}
	bi := int(math.Floor(float64(numsLen / 2)))
	switch isLeftOrRight {
	case -1:
		currentIdx = bi
	case 0:
		/* 좌로 이동 */
		currentIdx -= numsLen - bi
	case 1:
		/* 우로 이동 */
		currentIdx += bi
	}
	if nums[bi] == target {
		ans = currentIdx
	} else if numsLen > 1 {
		if nums[bi] > target {
			isLeftOrRight = Left
			ans = binarySearchBad(nums[:bi], target)
		} else {
			isLeftOrRight = Right
			ans = binarySearchBad(nums[bi:], target)
		}
	}

	return ans
}

func binarySearchGood(nums []int, target int) int {
	numsLen := len(nums)
	startIndex := 0
	endIndex := numsLen - 1
	mid := numsLen / 2

	for startIndex <= endIndex {
		if nums[mid] == target {
			return mid
		} else if nums[mid] > target {
			/* 좌로 이동 */
			endIndex = mid - 1
		} else {
			/* 우로 이동 */
			startIndex = mid + 1
		}
		mid = (endIndex + startIndex) / 2
	}

	return -1
}

func main() {
	nums := make([]int, 0, 10000)
	target := 0
	ans := -1

	/*nums = []int{-1, 0, 3, 5, 9, 12}
	target = 9
	ans = BinarySearch(nums, target)
	fmt.Println(ans)*/

	/*nums = []int{-1, 0, 3, 5, 9, 12}
	target = 2
	ans = BinarySearch(nums, target)
	fmt.Println(ans)*/

	nums, _ = util.RandRange(false, -9999, 9999, 1000, util.SortAscending)
	fmt.Println("[" + util.IntSliceJoin(nums, ",") + "]")
	target = -4294
	fmt.Println(target)
	util.StartExecution()
	ans = BinarySearch(nums, target)
	util.PrintElapsedTime("BinarySearch")
	fmt.Println(ans)

	util.StartExecution()
	ans = binarySearchGood(nums, target)
	util.PrintElapsedTime("binarySearchGood")
	fmt.Println(ans)

	/*nums = []int{5}
	target = 5
	fmt.Println(target)
	ans = BinarySearch(nums, target)
	fmt.Println(ans)*/
}
