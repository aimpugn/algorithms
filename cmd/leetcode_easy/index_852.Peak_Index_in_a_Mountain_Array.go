package main

/*
https://leetcode.com/problems/peak-index-in-a-mountain-array/
*/

import (
	"fmt"
	"math"
)

/*
1. 반복문 돌면서 가장 큰 값을 가진 인덱스 반환
2. Binary Search
    func peakIndexInMountainArray(A []int) int {
        lo := 0
		hi := len(A) - 1;
        while (lo < hi) {
			mi := lo + int(math.Floor(float64(hi - lo) / float64(2)))

            if A[mi] < A[mi + 1] {
                lo = mi + 1;
			} else {
                hi = mi;
			}
        }
        return lo;
    }
3. 정렬해서 가장 큰 값의 인덱스 반환

*/

type Contour struct {
	index int
	value int
}

func peakIndexInMountainArray(arr []int) int {
	cnt := 0
	mountSlice := []Contour{}
	for idx, num := range arr {
		mountSlice = append(mountSlice, Contour{index: idx, value: num})
		cnt++
	}
	if cnt == 0 {
		return 0
	}

	quickSortMountainArray(&mountSlice, 0, len(mountSlice)-1)

	return mountSlice[cnt-1].index
}

func quickSortMountainArray(arr *[]Contour, lo int, hi int) {
	if lo < hi {
		pivot := quickSortMountainArrayPartition(arr, lo, hi)
		quickSortMountainArray(arr, lo, pivot)
		quickSortMountainArray(arr, pivot+1, hi)
	}
}

func quickSortMountainArrayPartition(arr *[]Contour, lo int, hi int) int {
	a := *arr
	d := math.Floor(float64(lo+hi) / float64(2))
	pivot := a[int(d)].value
	i := lo - 1
	j := hi + 1
	for true {
		i++
		for true {
			if a[i].value < pivot {
				i++
			} else {
				break
			}
		}
		j--
		for true {
			if a[j].value > pivot {
				j--
			} else {
				break
			}
		}
		/*
			for a[i].value < pivot {
				i++
			}
			for a[j].value > pivot {
				j--
			}
		*/

		if i >= j {
			break
		}

		a[i], a[j] = a[j], a[i]
	}

	return j
}

/*
- 분할 정복 알고리즘
- 배열에서 pivot 선택 -> pivot보다 큰지 또는 작은지에 따라 다른 요소들을 두 하위 배열로 분할 -> 하위 배열을 재귀적으로 정렬
- 비교 정렬 알고리즘 -> '~보다 작음(less-than)'으로 정의되는 모든 타입의 항목 정렬 가능
- 수학적 분석에서,
	- 평균적으로 n 개의 항목을 정렬할 때 최대 O(n log n) 시간 소요되며,
	- 최악의 경우, 드물긴 하지만, 최대 O(n^2) 시간 소요
*/
func quickSortLomuto(arr *[]int, lo int, hi int) {
	/*
		- 일반적으로 배열의 마지막 요소를 pivot으로 선택
		- 방식
			- 인덱스 i를 유지,
			- 다른 인덱스 j 사용하여 배열을 스캔,
				- lo ~ i-1까지는 pivot보다 작게
				- i ~ j까지 pivot보다 크거나 같게
		- 배열이 이미 정렬되어 있다면 O(n^2)로 저하
		algorithm quicksort(A, lo, hi) is
		    if lo < hi then
				// pivot보다 작은 값/크거나 같은 값을 스왑
		        p := partition(A, lo, hi)
		        quicksort(A, lo, p - 1)
		        quicksort(A, p + 1, hi)

		algorithm partition(A, lo, hi) is
		    pivot := A[hi]
		    i := lo
		    for j := lo to hi do
		        if A[j] < pivot then
		            swap A[i] with A[j]
		            i := i + 1
		    swap A[i] with A[hi]
		    return i
	*/
	// 배열의 사이즈가 hi - lo > 0이어야 함
	if lo < hi {
		pivot := quickSortLomutoPartition(arr, lo, hi)
		quickSortLomuto(arr, lo, pivot-1)
		quickSortLomuto(arr, pivot+1, hi)
	}
}

/*
lo = 0
hi = 6
A = [4 7 8 3 6 9 5]
pivot = A[6] = 5

1.
i = 0
j = 0
i
j
4 7 8 3 6 9 5
4 7 8 3 6 9 5

2.
i = 1
j = 1
  i
  j
4 7 8 3 6 9 5

3.
i = 1
j = 2
  i
	j
4 7 8 3 6 9 5

4.
i = 1
j = 3
  i
	  j
4 7 8 3 6 9 5
4 3 8 7 6 9 5

5.
i = 2
j = 4
    i
		j
4 3 7 8 6 9 5

6.
i = 2
j = 5
	i
		  j
4 3 7 8 6 9 5

7.
i = 2
j = 6
	i
			j
4 3 7 8 6 9 5

8. swap A[i] with A[hi]
	i       hi
4 3 7 8 6 9 5
4 3 5 8 6 9 7

9. return i = 2
*/
func quickSortLomutoPartition(arr *[]int, lo int, hi int) int {
	a := *arr
	pivot := a[hi]
	i := lo
	for j := lo; j <= hi; j++ {
		// pivot보다 작은 값은 계속 좌측으로 넘긴다
		if a[j] < pivot {
			a[i], a[j] = a[j], a[i]
			// i는 점차 중앙으로 이동
			i++
		}
	}
	a[i], a[hi] = a[hi], a[i]

	return i
}

/*
Lomuto의 분할 방식보다 더 효과적
- 평균적으로 스왑을 세 배 덜한다
- 모든 값이 동일한 경우에도 효과적인 분할 생성
- 다른 것들과 마찬가지로 안정적인 정렬(stable sort)을 하지 못한다
*/
func quickSortHoare(arr *[]int, lo int, hi int) {
	/*
		algorithm quicksort(A, lo, hi) is
			if lo < hi then
				p := partition(A, lo, hi)
				quicksort(A, lo, p)
				quicksort(A, p + 1, hi)

		algorithm partition(A, lo, hi) is
			pivot := A[⌊(hi + lo) / 2⌋]
			// 분할된 배열의 끝과 끝에서 시작해서 서로 상대적으로 잘못된 순서로 되어 있는 한 쌍의 요소(X, Y에서 X <= pivot, Y >= pivot)를 발견할 때까지 다가간다
			i := lo - 1
			j := hi + 1
			loop forever
				do
			      i := i + 1
				while A[i] < pivot

				do
			    j := j - 1
				while A[j] > pivot

				if i ≥ j then
		            return j

				swap A[i] with A[j]
	*/
	if lo < hi {
		pivot := quickSortHoarePartition(arr, lo, hi)
		quickSortHoare(arr, lo, pivot)
		quickSortHoare(arr, pivot+1, hi)
	}
}

/*
lo = 0
hi = 6
A = [4 7 8 3 6 9 5]
pivot = A[6] = 5

1.
pivot = arr[((0 + 6) / 2 = 3)] = 3
i
			j
4 7 8 3 6 9 5

2.
i
		  j
4 7 8 3 6 9 5

3.
i
		j
4 7 8 3 6 9 5

4.
i
	  j
4 7 8 3 6 9 5

5. i = 0 <= j = 3, skip

6. swap
i
	  j
4 7 8 3 6 9 5
3 7 8 4 6 9 5

7.
i
	j
3 7 8 4 6 9 5

8.
i
  j
3 7 8 4 6 9 5

9.
i
j
3 7 8 4 6 9 5

10. i = 0 >= j = 0

11. return j = 0


pivot: 1
i
	j
0 1 0
  i
	j
0 1 0



*/
func quickSortHoarePartition(arr *[]int, lo int, hi int) int {
	a := *arr
	d := math.Floor(float64(lo+hi) / float64(2))
	pivot := a[int(d)]
	// leftIndex := lo
	// rightIndex := hi
	leftIndex := lo - 1
	rightIndex := hi + 1
	for true {
		// leftIndex와 rightIndex는 가장 바깥의 for문을 다시 돌 때 먼저 증/감을 하고 그 다음 단계로 넘어가야 한다

		// pivot보다 큰 값을 찾을 때까지 좌에서 우로 이동, 찾으면 다음 단계
		// golang에서는 ++idx가 안 되며, do while문도 없다
		leftIndex++
		for true {
			if a[leftIndex] < pivot {
				leftIndex++
			} else {
				break
			}
		}
		// pivot보다 작은 값을 찾을 때까지 우에서 좌로 이동, 찾으면 다음 단뎨
		rightIndex--
		for true {
			if a[rightIndex] > pivot {
				rightIndex--
			} else {
				break
			}
		}

		// 좌의 인덱스와 우의 인덱스가 만나거나 역전된 경우
		// - i와 j가 같은 값을 가리키는 경우
		// - i가 j를 지나쳐간 경우
		if leftIndex >= rightIndex {
			break
		}

		// 좌의 인덱스와 우의 인덱스가 만나지 않은 경우
		// pivot보다 큰 값과 pivot보다 작은 값을 치환
		a[leftIndex], a[rightIndex] = a[rightIndex], a[leftIndex]
	}

	return rightIndex
}

func main() {
	test := []int{}
	// quickSortLomuto(&test, 0, len(test) -1)
	// fmt.Println(test)

	// quickSortHoare(&test, 0, len(test) -1)
	// fmt.Println(test)
	// test = append(test, 3, 2, 6, 5, 8)
	// quickSortHoare(&test, 0, len(test) -1)
	// fmt.Println(test)

	// test = append(test, 0, 1, 0)
	// test = append(test, 0, 2, 1, 0)
	// test = append(test, 0, 10, 5, 2)
	test = append(test, 3, 4, 5, 1)
	ans := peakIndexInMountainArray(test)
	fmt.Println(ans)
}
