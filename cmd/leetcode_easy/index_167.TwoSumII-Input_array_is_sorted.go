package main

import (
	"fmt"

	"algorithms/pkg/util"
)

/*
오름차순으로 정렬된 상태
특정 타겟 숫자로 합할 수 있는 두 숫자 찾기
idx1 < idx2
idx1, idx2는 모두 0부터 시작(zero-based)하지 않는다
정확히 한 답만 있고, 같은 요소를 두번 쓰지 않는다
*/

/*
Constraints:
2 <= nums.length <= 3 * 10^4
-1000 <= nums[i] <= 1000
nums is sorted in increasing order.
-1000 <= target <= 1000
*/
func twoSumII(numbers []int, target int) []int {
	ans := make([]int, 0, 2)
	numbersLen := len(numbers)
	/*
		https://stackoverflow.com/a/54602693/8562273
		https://www.ardanlabs.com/blog/2013/11/label-breaks-in-go.html
	*/
loop: /* label */
	for idx, number1 := range numbers {
		for i := idx + 1; i < numbersLen; i++ {
			if tmp := number1 + numbers[i]; tmp == target {
				ans = append(ans, idx+1, i+1)
				break loop
			}
		}
	}

	return ans
}

func twoSumIIEnhanced(numbers []int, target int) []int {
	ans := make([]int, 0, 2)
	numbersLen := len(numbers)
	/*
		오름차순으로 정렬되어 있으므로, 아래에서 위로 움직는 i와 위에서 아래로 움직이는 j로 좁힐 수 있다.
	*/
	i, j := 0, numbersLen-1
loop: /* label */
	for i < j {
		tmp := numbers[i] + numbers[j]
		if tmp == target {
			ans = append(ans, i+1, j+1)
			break loop
		} else if tmp > target {
			/* 타겟보다 크면 큰 값을 줄여나간다 */
			j--
		} else {
			/* 타겟보다 작으면 작은 값을 늘려나간다 */
			i++
		}
	}

	return ans
}

func main() {
	numbers := make([]int, 0, 30000)
	target := 0
	ans := make([]int, 0, 2)
	numbers, target = append(numbers, 2, 7, 11, 15), 9
	ans = twoSumII(numbers, target)
	fmt.Println(ans)
	numbers, target, ans = nil, 0, nil

	numbers, target = append(numbers, 2, 3, 4), 6
	ans = twoSumII(numbers, target)
	fmt.Println(ans)
	numbers, target, ans = nil, 0, nil

	numbers, target = append(numbers, -1, 0), -1
	ans = twoSumII(numbers, target)
	fmt.Println(ans)
	numbers, target, ans = nil, 0, nil

	numbers, _ = util.RandRange(-1000, 1000, 1500, util.SortAscending)
	target = 10
	util.StartExecution()
	ans = twoSumII(numbers, target)
	util.PrintElapsedTime()
	fmt.Println(ans)
	ans = nil
	util.StartExecution()
	ans = twoSumIIEnhanced(numbers, target)
	util.PrintElapsedTime()
	fmt.Println(ans)
}
