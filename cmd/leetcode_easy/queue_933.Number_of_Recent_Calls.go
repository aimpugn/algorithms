package main

import (
	"fmt"
)

/*
You have a RecentCounter class which counts the number of recent requests within a certain time frame.

Implement the RecentCounter class:

RecentCounter() Initializes the counter with zero recent requests.
int ping(int t) Adds a new request at time t, where t represents some time in milliseconds,
and returns the number of requests that has happened in the past 3000 milliseconds (including the new request).

Specifically, return the number of requests that have happened in the inclusive range [t - 3000, t].

It is guaranteed that every call to ping uses a strictly larger value of t than the previous call

["RecentCounter", "ping", "ping", "ping", "ping"]
[[], [1], [100], [3001], [3002]]

Explanation
RecentCounter recentCounter = new RecentCounter();
recentCounter.ping(1);     // requests = [1], range is [-2999,1], return 1
recentCounter.ping(100);   // requests = [1, 100], range is [-2900,100], return 2
recentCounter.ping(3001);  // requests = [1, 100, 3001], range is [1,3001], return 3
recentCounter.ping(3002);  // requests = [1, 100, 3001, 3002], range is [2,3002], return 3
*/
type RecentCounter struct {
	requests []int
}

func Constructor() RecentCounter {
	return RecentCounter{requests: []int{}}
}

// t 시간에 새로운 요청을 추가
// 지난 3000 milliseconds 동안 발생한 요청 수를 반환
func (this *RecentCounter) Ping(t int) int {
	this.requests = append(this.requests, t)
	// 464 ms -> 168 ms -> 124 ms
	/*for startIdx < this.counter {
		// 시작 지점인 `start`보다 현재 값이 작으면 다음부터는 그 부분은 스킵해도 된다
		// `스킵 해도 된다` === `제거해도 된다`
		if this.requests[startIdx] < t - 3000 {
			// 제거되는 만큼 인덱스 증가
			startIdx++
			continue
		}

		break
	}*/
	/*for this.requests[startIdx] < t - 3000 {
		// 시작 지점인 `start`보다 현재 값이 작으면 다음부터는 그 부분은 스킵해도 된다
		// `스킵 해도 된다` === `제거해도 된다`
		startIdx++
	}*/
	// 464 ms -> 168 ms -> 124 ms -> 128 ms. 속도는 서버에 따라 다른 듯?
	for len(this.requests) > 0 {
		if this.requests[0] < t-3000 {
			this.requests = this.requests[1:]
		} else {
			break
		}
	}

	return len(this.requests)
}

type RecentCounter2 struct {
	time     int
	startIdx int
	requests []int
	counter  int
}

func Constructor2() RecentCounter2 {
	return RecentCounter2{time: 3000, requests: []int{}, counter: 0, startIdx: 0}
}

// t 시간에 새로운 요청을 추가
// 지난 3000 milliseconds 동안 발생한 요청 수를 반환
func (rc2 *RecentCounter2) Ping2(t int) int {
	rc2.requests = append(rc2.requests, t)
	rc2.counter++
	requestNum := 1 // t 자신은 무조건 포함된다
	start := t - rc2.time
	end := t
	// startIdx 사용 후 464 ms -> 168 ms
	for i := rc2.startIdx; i < rc2.counter-1; i++ {
		// 시작 지점인 `start`보다 현재 값이 작으면 다음부터는 그 부분은 스킵해도 된다
		if start > rc2.requests[i] {
			rc2.startIdx++
			continue
		}
		if start <= rc2.requests[i] && rc2.requests[i] <= end {
			requestNum++
		}
	}

	return requestNum
}

/**
 * Your RecentCounter object will be instantiated and called as such:
 * obj := Constructor();
 * param_1 := obj.Ping(t);
 */

func main() {
	obj := Constructor()
	param_1 := obj.Ping(1)
	fmt.Println(param_1)
	param_1 = obj.Ping(100)
	fmt.Println(param_1)
	param_1 = obj.Ping(3001)
	fmt.Println(param_1)
	param_1 = obj.Ping(3002)
	fmt.Println(param_1)
}
