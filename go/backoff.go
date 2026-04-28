package tesote

import (
	"math"
	"time"
)

// backoff computes the delay before the next retry. retryAfter (seconds) wins
// when set; otherwise full-jitter exponential backoff capped at MaxDelay.
func (c *Client) backoff(attempt int, retryAfter int) time.Duration {
	if retryAfter > 0 {
		d := time.Duration(retryAfter) * time.Second
		if d > c.retry.MaxDelay {
			return c.retry.MaxDelay
		}
		return d
	}
	exp := time.Duration(math.Min(float64(c.retry.MaxDelay), float64(c.retry.BaseDelay)*math.Pow(2, float64(attempt-1))))
	if exp <= 0 {
		return 0
	}
	jit := time.Duration(c.rand63() % uint64(exp))
	return jit
}
