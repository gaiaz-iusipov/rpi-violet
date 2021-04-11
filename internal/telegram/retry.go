package telegram

import (
	"fmt"
	"time"
)

func retry(attempts int, delay time.Duration, fn func() error) (err error) {
	if attempts < 1 {
		return
	}
	for i := 0; i < attempts; i++ {
		if i > 0 {
			time.Sleep(delay)
		}
		err = fn()
		if err == nil {
			return
		}
	}
	return fmt.Errorf("after %d attempts: %w", attempts, err)
}
