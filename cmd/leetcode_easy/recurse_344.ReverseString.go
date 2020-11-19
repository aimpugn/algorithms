package main

import (
	"fmt"
)

/*
Input: ["h","e","l","l","o"]
Output: ["o","l","l","e","h"]
*/
func reverseString(s []byte) {
	// 첫 h를 끝으로, 끝의 o를 처음으로
	reverseStringHelper(s, 0, len(s)-1)
}

func reverseStringHelper(s []byte, startIdx, targetIdx int) {
	// 현재 문자의 인덱스가 문자의 인덱스보다 크거나 같은 경우, 즉 역전하려는 경우
	// 짝수/홀수를 나눌 필요가 없나?
	// 짝수인 경우 같다에 걸린다
	// 홀수인 경우 크다에 걸린다
	if startIdx >= targetIdx {
		return
	}

	s[startIdx], s[targetIdx] = s[targetIdx], s[startIdx]
	reverseStringHelper(s, startIdx+1, targetIdx-1)
}

func reverseStringHelper2(s *[]byte, startIdx int, targetIdx int) {
	lenS := len(*s)
	if lenS%2 != 0 {
		// 문자의 수가 홀수인 경우
		// 5 / 2 = 2, 7 / 2 = 3, 즉 중앙의 문자인 경우
		if startIdx == lenS/2 { // <<< 같다
			return
		}
	} else {
		// 문자의 수가 짝수인 경우
		// 4 / 2 = 2, 6 / 2 = 3
		if startIdx >= lenS/2 { // <<< 크거나 같다. 이 경우가 위의 `홀수 && 같다` 경우를 포함한다
			return
		}
	}
	(*s)[startIdx], (*s)[targetIdx] = (*s)[targetIdx], (*s)[startIdx]
	reverseStringHelper2(s, startIdx+1, targetIdx-1)
}

func main() {
	var input []byte

	input = []byte("A man, a plan, a canal: Panama")
	fmt.Print("[")
	for _, char := range input {
		fmt.Print("\"" + string(char) + "\",")
	}
	fmt.Print("]")
	fmt.Println()
	reverseString(input)
	for _, char := range input {
		fmt.Print(string(char))
	}
	fmt.Println()

}
