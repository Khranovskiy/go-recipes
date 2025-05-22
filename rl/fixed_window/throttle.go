package fixed_window

import (
	"errors"
	"time"
)

var ErrCanceled error = errors.New("canceled")

// Throttle ensures that fn runs no more than limit times per second.
func Throttle(limit int, fn func()) (handle func() error, cancel func()) {
	ticker := time.NewTicker(time.Second / time.Duration(limit))
	canceled := make(chan struct{})

	handle = func() error {
		select {
		case <-ticker.C:
			go fn()
			return nil
		case <-canceled:
			return ErrCanceled
		}
	}

	cancel = func() {
		select {
		case <-canceled:
		default:
			ticker.Stop()
			close(canceled)
		}
	}

	return handle, cancel
}
