package main

import (
	"fmt"
	"math"
	"unicode/utf8"
)

/*
Suppose Andy and Doris want to choose a restaurant for dinner, and they both have a list of favorite restaurants represented by strings.

You need to help them find out their common interest with the least list index sum.
If there is a choice tie between answers, output all of them with no order requirement.
You could assume there always exists an answer.
*/

/*
list1 = ["Shogun","Tapioca Express","Burger King","KFC"],
list2 = ["Piatti","The Grill at Torrey Pines","Hungry Hunter Steakhouse","Shogun"]
Output: ["Shogun"]

list1 = ["Shogun","Tapioca Express","Burger King","KFC"],
list2 = ["KFC","Shogun","Burger King"]
Output: ["Shogun"]

list1 = ["Shogun","Tapioca Express","Burger King","KFC"],
list2 = ["KFC","Burger King","Tapioca Express","Shogun"]
Output: ["KFC","Burger King","Tapioca Express","Shogun"]

list1 = ["Shogun","Tapioca Express","Burger King","KFC"],
list2 = ["KNN","KFC","Burger King","Tapioca Express","Shogun"]
Output: ["KFC","Burger King","Tapioca Express","Shogun"]

1 <= list1.length, list2.length <= 1000
1 <= list1[i].length, list2[i].length <= 30
list1[i] and list2[i] consist of spaces ' ' and English letters.
All the stings of list1 are unique.
All the stings of list2 are unique.

Runtime: 324 ms, faster than 6.45% of Go online submissions for Minimum Index Sum of Two Lists.
Memory Usage: 6.5 MB, less than 88.71% of Go online submissions for Minimum Index Sum of Two Lists.
*/
func findRestaurant(list1 []string, list2 []string) []string {
	list1Len := len(list1)
	ans := make([]string, 0, list1Len)
	tmp := make([]string, list1Len, list1Len)
	tmpIdx := 0
	leastIdx := 1000000
	idx := make([]int, list1Len, list1Len)
	for idx1, restaurant1 := range list1 {
		restaurant1Len := utf8.RuneCountInString(restaurant1)
		for idx2, restaurant2 := range list2 {
			restaurant2Len := utf8.RuneCountInString(restaurant2)
			if restaurant1Len != restaurant2Len {
				continue
			}
			restaurant1R := []rune(restaurant1)
			restaurant2R := []rune(restaurant2)
			for i := 0; i < restaurant1Len; i++ {
				if restaurant1R[i] != restaurant2R[i] {
					break
				} else if i == (restaurant1Len - 1) {
					currentIdx := idx1 + idx2
					leastIdx = int(math.Min(float64(currentIdx), float64(leastIdx)))
					tmp[tmpIdx] = restaurant1
					idx[tmpIdx] = currentIdx
					tmpIdx++
				}
			}
		}
	}

	for k := 0; k < tmpIdx; k++ {
		if idx[k] == leastIdx {
			ans = append(ans, tmp[k])
		}
	}

	return ans
}

/*
`rune`으로 바꾸고 비교하는 것보다 그냥 단순 문자열 비교가 훨씬 빠르다
Runtime: 36 ms, faster than 32.26% of Go online submissions for Minimum Index Sum of Two Lists.
Memory Usage: 6.5 MB, less than 79.03% of Go online submissions for Minimum Index Sum of Two Lists.
*/
func findRestaurant2(list1 []string, list2 []string) []string {
	list1Len := len(list1)
	ans := make([]string, 0, list1Len)
	tmp := make([]string, list1Len, list1Len)
	tmpIdx := 0
	leastIdx := 1000000
	idx := make([]int, list1Len, list1Len)
	for idx1, restaurant1 := range list1 {
		for idx2, restaurant2 := range list2 {
			if restaurant1 == restaurant2 {
				currentIdx := idx1 + idx2
				leastIdx = int(math.Min(float64(currentIdx), float64(leastIdx)))
				tmp[tmpIdx] = restaurant1
				idx[tmpIdx] = currentIdx
				tmpIdx++
				break
			}
		}
	}

	for k := 0; k < tmpIdx; k++ {
		if idx[k] == leastIdx {
			ans = append(ans, tmp[k])
		}
	}

	return ans
}

/*
`map`을 쓰는 게 더 빠르다. 그렇다면 `map`은 어떻게 문자열 키가 존재하는지 찾아내지?
*/

func main() {
	var list1, list2, output []string

	list1 = []string{"Shogun", "Tapioca Express", "Burger King", "KFC"}
	list2 = []string{"Piatti", "The Grill at Torrey Pines", "Hungry Hunter Steakhouse", "Shogun"}
	output = findRestaurant2(list1, list2)
	fmt.Println(output)

	list1 = []string{"Shogun", "Tapioca Express", "Burger King", "KFC"}
	list2 = []string{"KFC", "Shogun", "Burger King"}
	output = findRestaurant2(list1, list2)
	fmt.Println(output)

	list1 = []string{"Shogun", "Tapioca Express", "Burger King", "KFC"}
	list2 = []string{"KFC", "Burger King", "Tapioca Express", "Shogun"}
	output = findRestaurant2(list1, list2)
	fmt.Println(output)

	list1 = []string{"KFC"}
	list2 = []string{"KFC"}
	output = findRestaurant2(list1, list2)
	fmt.Println(output)

}
