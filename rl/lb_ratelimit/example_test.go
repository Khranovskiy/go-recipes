package lb_ratelimit

import (
	"fmt"
	"time"
)

func Example_default() {
	rl := New(100) // per second, some slack.

	rl.Take()                         // Initialize.
	time.Sleep(time.Millisecond * 45) // Let some time pass.

	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		if i > 0 {
			fmt.Println(i, now.Sub(prev).Round(time.Millisecond*2))
		}
		prev = now
	}

	// Output:
	// 1 0s
	// 2 0s
	// 3 0s
	// 4 4ms
	// 5 10ms
	// 6 10ms
	// 7 10ms
	// 8 10ms
	// 9 10ms
}

func Example_withoutSlack() {
	rl := New(100, WithoutSlack) // per second, no slack.

	prev := time.Now()
	for i := 0; i < 6; i++ {
		now := rl.Take()
		if i > 0 {
			fmt.Println(i, now.Sub(prev))
		}
		prev = now
	}

	// Output:
	// 1 10ms
	// 2 10ms
	// 3 10ms
	// 4 10ms
	// 5 10ms
}
