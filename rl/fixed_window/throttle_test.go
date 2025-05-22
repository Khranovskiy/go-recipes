package fixed_window

import (
	"testing"
	"time"
)

func TestThrottle(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		requests int
		wantTime time.Duration
	}{
		{
			name:     "1 request per second",
			limit:    1,
			requests: 5,
			wantTime: 6 * time.Second,
		},
		{
			name:     "2 requests per second",
			limit:    2,
			requests: 6,
			wantTime: 4 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := 0
			work := func() {
				counter++
			}

			handle, cancel := Throttle(tt.limit, work)
			defer cancel()

			start := time.Now()
			for i := 0; i < tt.requests; i++ {
				if err := handle(); err != nil {
					t.Errorf("handle() error = %v", err)
				}
			}

			time.Sleep(time.Second)

			elapsed := time.Since(start)
			if counter != tt.requests {
				t.Errorf("Expected %d executions, got %d", tt.requests, counter)
			}

			if elapsed < tt.wantTime-time.Second || elapsed > tt.wantTime+time.Second {
				t.Errorf("Expected execution time around %v, got %v", tt.wantTime, elapsed)
			}
		})
	}
}

func TestThrottleCancel(t *testing.T) {
	counter := 0
	work := func() {
		counter++
	}

	handle, cancel := Throttle(1, work)

	if err := handle(); err != nil {
		t.Errorf("handle() error = %v", err)
	}

	cancel()

	if err := handle(); err != ErrCanceled {
		t.Errorf("Expected ErrCanceled, got %v", err)
	}

	time.Sleep(time.Second)
	if counter != 1 {
		t.Errorf("Expected 1 execution, got %d", counter)
	}
}
