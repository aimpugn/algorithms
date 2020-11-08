package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

/*
https://leetcode.com/problems/check-if-a-word-occurs-as-a-prefix-of-any-word-in-a-sentence/
*/

func isPrefixOfWord(sentence string, searchWord string) int {
	ans := -1
	sentenceSlice := strings.Split(sentence, " ")

	searchWordRuneSlice := make([]rune, 0, 10)
	searchWordSliceLen := 0

	for _, searchWordRune := range searchWord {
		searchWordRuneSlice = append(searchWordRuneSlice, searchWordRune)
		searchWordSliceLen++
	}

	for idx, word := range sentenceSlice {
		loopCnt := 0
		wordRuneSliceLen := utf8.RuneCountInString(word)

		if wordRuneSliceLen < searchWordSliceLen {
			continue
		}

		wordRuneSlice := make([]rune, 0, wordRuneSliceLen)

		for _, wordRune := range word {
			wordRuneSlice = append(wordRuneSlice, wordRune)
			wordRuneSliceLen++
		}

		for offsetOfSearchWord, searchWordRune := range searchWord {
			if searchWordRune == wordRuneSlice[offsetOfSearchWord] {
				loopCnt++
			} else {
				break
			}
			if searchWordSliceLen == loopCnt {
				return idx + 1
			}
		}
	}

	return ans
}

func main() {
	sentence := ""
	searchWord := ""
	ans := -1

	sentence = "i love eating burger"
	searchWord = "burg"
	ans = isPrefixOfWord(sentence, searchWord)
	fmt.Println(ans)

	sentence = "this problem is an easy problem"
	searchWord = "pro"
	ans = isPrefixOfWord(sentence, searchWord)
	fmt.Println(ans)

	sentence = "i am tired"
	searchWord = "you"
	ans = isPrefixOfWord(sentence, searchWord)
	fmt.Println(ans)

	sentence = "i use triple pillow"
	searchWord = "pill"
	ans = isPrefixOfWord(sentence, searchWord)
	fmt.Println(ans)

	sentence = "hello from the other side"
	searchWord = "they"
	ans = isPrefixOfWord(sentence, searchWord)
	fmt.Println(ans)
}
