package tesote

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// secureUint63 returns a non-negative 63-bit value from crypto/rand.
func secureUint63() uint64 {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		// why: rand.Read on Linux backs onto getrandom; failure is exceptional. Fall
		// back to a deterministic value rather than panicking.
		return 0
	}
	v := uint64(0)
	for _, c := range b {
		v = (v << 8) | uint64(c)
	}
	return v &^ (1 << 63)
}

// ctxSleep is the default Sleep implementation: respects context cancellation.
func ctxSleep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func atoiHeader(h http.Header, name string) (int, bool) {
	v := h.Get(name)
	if v == "" {
		return 0, false
	}
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return 0, false
	}
	return n, true
}

func parseRetryAfter(h http.Header) int {
	v := h.Get("Retry-After")
	if v == "" {
		return 0
	}
	if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
		return n
	}
	return 0
}

func isRetryableStatus(status int) bool {
	switch status {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}
	return false
}

func setAttempts(err error, n int) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		apiErr.Attempts = n
	}
}

func classifyTransportError(err error, method string, summary RequestSummary, attempt int) error {
	base := &TransportError{
		Op:             method,
		Message:        err.Error(),
		RequestSummary: summary,
		Attempts:       attempt,
		Cause:          err,
	}
	var tlsErr *tls.RecordHeaderError
	if errors.As(err, &tlsErr) {
		return &TLSError{TransportError: base}
	}
	if strings.Contains(strings.ToLower(err.Error()), "tls") || strings.Contains(strings.ToLower(err.Error()), "x509") {
		return &TLSError{TransportError: base}
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return &TimeoutError{TransportError: base}
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return &TimeoutError{TransportError: base}
	}
	return &NetworkError{TransportError: base}
}

func shouldRetryTransport(err error, method string, hasIdempotencyKey bool) bool {
	var timeoutErr *TimeoutError
	if errors.As(err, &timeoutErr) {
		// why: read-timeout on a non-idempotent mutation may have succeeded server-side.
		if _, mutating := mutatingMethods[method]; mutating && !hasIdempotencyKey {
			return false
		}
		return true
	}
	var netErr *NetworkError
	if errors.As(err, &netErr) {
		return true
	}
	return false
}

func readAll(r io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	return buf.Bytes(), err
}
