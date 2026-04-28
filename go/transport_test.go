package tesote

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// newTestClient builds a Client wired to the given test server with
// deterministic backoff/random helpers so tests stay fast and reproducible.
func newTestClient(t *testing.T, server *httptest.Server, opt ...func(*Options)) *Client {
	t.Helper()
	o := Options{
		APIKey:     "secret-abcd1234",
		BaseURL:    server.URL,
		Sleep:      func(_ context.Context, _ time.Duration) error { return nil },
		RandUint63: func() uint64 { return 0 },
	}
	for _, fn := range opt {
		fn(&o)
	}
	c, err := NewClient(o)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func TestNewClient_RequiresAPIKey(t *testing.T) {
	_, err := NewClient(Options{})
	if err == nil {
		t.Fatal("expected ConfigError, got nil")
	}
	if !errors.Is(err, ErrConfig) {
		t.Fatalf("expected ErrConfig, got %v", err)
	}
}

func TestClient_DefaultsBaseURLAndUserAgent(t *testing.T) {
	c, err := NewClient(Options{APIKey: "abcd"})
	if err != nil {
		t.Fatal(err)
	}
	if c.BaseURL() != DefaultBaseURL {
		t.Errorf("base url = %q, want %q", c.BaseURL(), DefaultBaseURL)
	}
	if !strings.HasPrefix(c.UserAgent(), "tesote-sdk-go/") {
		t.Errorf("user agent = %q, want tesote-sdk-go/ prefix", c.UserAgent())
	}
	if !strings.Contains(c.UserAgent(), "(go/") {
		t.Errorf("user agent missing runtime: %q", c.UserAgent())
	}
}

func TestDo_GETSuccess_HeadersAndDecode(t *testing.T) {
	var seenAuth, seenAccept, seenUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenAuth = r.Header.Get("Authorization")
		seenAccept = r.Header.Get("Accept")
		seenUA = r.Header.Get("User-Agent")
		w.Header().Set("X-Request-Id", "req-123")
		w.Header().Set("X-RateLimit-Limit", "200")
		w.Header().Set("X-RateLimit-Remaining", "199")
		w.Header().Set("X-RateLimit-Reset", "60")
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	var out struct {
		OK bool `json:"ok"`
	}
	resp, err := c.Do(context.Background(), "GET", "/v3/ping", RequestOptions{Out: &out})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if !out.OK {
		t.Errorf("decode failed, got %+v", out)
	}
	if resp.RequestID != "req-123" {
		t.Errorf("request id = %q, want req-123", resp.RequestID)
	}
	if seenAuth != "Bearer secret-abcd1234" {
		t.Errorf("auth header = %q", seenAuth)
	}
	if seenAccept != "application/json" {
		t.Errorf("accept = %q", seenAccept)
	}
	if !strings.HasPrefix(seenUA, "tesote-sdk-go/") {
		t.Errorf("user agent = %q", seenUA)
	}
	rl, ok := c.LastRateLimit()
	if !ok || rl.Limit != 200 || rl.Remaining != 199 || rl.Reset != 60 {
		t.Errorf("rate limit = %+v ok=%v", rl, ok)
	}
}

func TestDo_POSTAutoIdempotencyKey(t *testing.T) {
	var key string
	var contentType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key = r.Header.Get("Idempotency-Key")
		contentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, `{}`)
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	_, err := c.Do(context.Background(), "POST", "/v3/widgets", RequestOptions{Body: map[string]string{"name": "x"}})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if len(key) != 36 || strings.Count(key, "-") != 4 {
		t.Errorf("idempotency-key not UUIDv4-shaped: %q", key)
	}
	if contentType != "application/json" {
		t.Errorf("content-type = %q", contentType)
	}
}

func TestDo_POSTUserIdempotencyKeyPreserved(t *testing.T) {
	var seen string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r.Header.Get("Idempotency-Key")
		_, _ = io.WriteString(w, `{}`)
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	_, err := c.Do(context.Background(), "POST", "/v3/widgets", RequestOptions{IdempotencyKey: "user-supplied-key"})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if seen != "user-supplied-key" {
		t.Errorf("idempotency-key = %q", seen)
	}
}

func TestDo_RetryOn503ThenSuccess(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&hits, 1)
		if n < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = io.WriteString(w, `{}`)
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	resp, err := c.Do(context.Background(), "GET", "/v3/ping", RequestOptions{})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if resp.Attempts != 3 {
		t.Errorf("attempts = %d, want 3", resp.Attempts)
	}
}

func TestDo_RetryExhausted_TypedError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = io.WriteString(w, `{"error":"service paused","error_code":"SERVICE_PAUSED"}`)
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	_, err := c.Do(context.Background(), "GET", "/v3/ping", RequestOptions{})
	if err == nil {
		t.Fatal("expected error")
	}
	var sErr *ServiceUnavailableError
	if !errors.As(err, &sErr) {
		t.Fatalf("expected *ServiceUnavailableError, got %T: %v", err, err)
	}
	if sErr.Attempts != 3 {
		t.Errorf("attempts = %d, want 3", sErr.Attempts)
	}
	if sErr.HTTPStatus != 503 {
		t.Errorf("http status = %d", sErr.HTTPStatus)
	}
	if sErr.RequestSummary.Authorization != "Bearer 1234" {
		t.Errorf("redacted auth = %q", sErr.RequestSummary.Authorization)
	}
}

func TestDo_RateLimitRetryAfterUsed(t *testing.T) {
	var slept []time.Duration
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&hits, 1)
		if n == 1 {
			w.Header().Set("Retry-After", "2")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = io.WriteString(w, `{"error_code":"RATE_LIMIT_EXCEEDED"}`)
			return
		}
		_, _ = io.WriteString(w, `{}`)
	}))
	defer srv.Close()

	c := newTestClient(t, srv, func(o *Options) {
		o.Sleep = func(_ context.Context, d time.Duration) error {
			slept = append(slept, d)
			return nil
		}
	})

	_, err := c.Do(context.Background(), "GET", "/v3/ping", RequestOptions{})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if len(slept) != 1 || slept[0] != 2*time.Second {
		t.Errorf("slept = %v, want [2s]", slept)
	}
}

func TestDo_NonRetryable4xxNotRetried(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"error":"nope","error_code":"UNAUTHORIZED"}`)
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	_, err := c.Do(context.Background(), "GET", "/v3/ping", RequestOptions{})
	if err == nil {
		t.Fatal("expected error")
	}
	if hits != 1 {
		t.Errorf("hits = %d, want 1 (non-retryable)", hits)
	}
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("not unauthorized: %v", err)
	}
}

func TestDo_RequestIDIntoError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "req-xyz")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"error_code":"UNAUTHORIZED"}`)
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	_, err := c.Do(context.Background(), "GET", "/v3/ping", RequestOptions{})
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.RequestID != "req-xyz" {
		t.Errorf("request id = %q", apiErr.RequestID)
	}
}

func TestDo_TTLCacheHit(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer srv.Close()
	cache := NewLRUCache(8)
	c := newTestClient(t, srv, func(o *Options) { o.Cache = cache })

	for i := 0; i < 3; i++ {
		_, err := c.Do(context.Background(), "GET", "/v3/widgets", RequestOptions{CacheTTL: time.Minute})
		if err != nil {
			t.Fatalf("Do: %v", err)
		}
	}
	if hits != 1 {
		t.Errorf("hits = %d, want 1 (cached)", hits)
	}
}

func TestDo_CacheBypassWhenTTLZero(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		_, _ = io.WriteString(w, `{}`)
	}))
	defer srv.Close()
	cache := NewLRUCache(8)
	c := newTestClient(t, srv, func(o *Options) { o.Cache = cache })

	for i := 0; i < 2; i++ {
		_, err := c.Do(context.Background(), "GET", "/v3/widgets", RequestOptions{})
		if err != nil {
			t.Fatalf("Do: %v", err)
		}
	}
	if hits != 2 {
		t.Errorf("hits = %d, want 2 (no cache)", hits)
	}
}

func TestDo_QueryParamsSorted(t *testing.T) {
	var seen string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r.URL.RawQuery
		_, _ = io.WriteString(w, `{}`)
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	_, err := c.Do(context.Background(), "GET", "/v3/widgets", RequestOptions{
		Query: map[string]string{"z": "1", "a": "2", "m": "3"},
	})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if seen != "a=2&m=3&z=1" {
		t.Errorf("query = %q, want sorted a,m,z", seen)
	}
}

func TestDo_NetworkErrorRetried(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{}`)
	}))
	srv.Close() // why: closed server -> connection refused, classified as NetworkError

	c := newTestClient(t, srv)
	_, err := c.Do(context.Background(), "GET", "/v3/ping", RequestOptions{})
	if err == nil {
		t.Fatal("expected network error")
	}
	if !errors.Is(err, ErrNetwork) {
		t.Errorf("expected ErrNetwork, got %v", err)
	}
	var netErr *NetworkError
	if !errors.As(err, &netErr) {
		t.Fatalf("not a *NetworkError: %T", err)
	}
	if netErr.Attempts != 3 {
		t.Errorf("attempts = %d, want 3", netErr.Attempts)
	}
}

func TestDo_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := c.Do(ctx, "GET", "/v3/ping", RequestOptions{})
	if err == nil {
		t.Fatal("expected cancellation error")
	}
}

func TestLRUCache_Eviction(t *testing.T) {
	c := NewLRUCache(2)
	c.Set("a", CachedResponse{Status: 200})
	c.Set("b", CachedResponse{Status: 200})
	c.Set("c", CachedResponse{Status: 200})
	if _, ok := c.Get("a"); ok {
		t.Errorf("a should have been evicted")
	}
	if _, ok := c.Get("b"); !ok {
		t.Errorf("b should remain")
	}
}

func TestLRUCache_TTLExpiry(t *testing.T) {
	c := NewLRUCache(2)
	c.Set("a", CachedResponse{Status: 200, ExpiresAt: time.Now().Add(-time.Second)})
	if _, ok := c.Get("a"); ok {
		t.Errorf("a should have expired")
	}
}

func TestRedactBearer(t *testing.T) {
	if got := RedactBearer("ab"); got != "Bearer <redacted>" {
		t.Errorf("short = %q", got)
	}
	if got := RedactBearer("supersecretkey"); got != "Bearer tkey" {
		t.Errorf("long = %q", got)
	}
}
