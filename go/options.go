package tesote

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultMaxAttempts = 3
	defaultBaseDelay   = 250 * time.Millisecond
	defaultMaxDelay    = 8 * time.Second
)

// RetryPolicy controls retry behavior. Zero value yields the documented defaults.
type RetryPolicy struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

func (p RetryPolicy) normalized() RetryPolicy {
	out := p
	if out.MaxAttempts < 1 {
		out.MaxAttempts = defaultMaxAttempts
	}
	if out.BaseDelay <= 0 {
		out.BaseDelay = defaultBaseDelay
	}
	if out.MaxDelay <= 0 {
		out.MaxDelay = defaultMaxDelay
	}
	return out
}

// RateLimit captures the most recent rate-limit headers from the API.
type RateLimit struct {
	Limit     int
	Remaining int
	Reset     int
}

// Logger receives one structured event per request and per response. It must
// not block; the transport calls it inline.
type Logger func(event LogEvent)

// LogEvent is the payload passed to a Logger.
type LogEvent struct {
	Phase     string
	Method    string
	URL       string
	Status    int
	RequestID string
	Attempt   int
	Err       error
}

// Options configures a Client.
type Options struct {
	APIKey      string
	BaseURL     string
	UserAgent   string
	HTTPClient  *http.Client
	RetryPolicy RetryPolicy
	Cache       CacheBackend
	Logger      Logger
	// Now overrides time.Now for deterministic tests.
	Now func() time.Time
	// RandUint63 overrides random jitter for deterministic tests.
	RandUint63 func() uint64
	// Sleep overrides time.Sleep for deterministic tests.
	Sleep func(context.Context, time.Duration) error
}
