package tesote

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

// DefaultBaseURL is the production API root.
const DefaultBaseURL = "https://equipo.tesote.com/api"

// Client is the single HTTP-touching component of the SDK. It is safe for
// concurrent use across goroutines and across all V*Client wrappers.
type Client struct {
	apiKey     string
	baseURL    string
	userAgent  string
	httpClient *http.Client
	retry      RetryPolicy
	cache      CacheBackend
	logger     Logger
	now        func() time.Time
	rand63     func() uint64
	sleep      func(context.Context, time.Duration) error

	mu            sync.RWMutex
	lastRateLimit RateLimit
	hasRateLimit  bool
}

// NewClient builds a Client. Returns *ConfigError if APIKey is empty.
func NewClient(opts Options) (*Client, error) {
	if strings.TrimSpace(opts.APIKey) == "" {
		return nil, &ConfigError{Field: "APIKey", Message: "must not be empty"}
	}
	base := opts.BaseURL
	if base == "" {
		base = DefaultBaseURL
	}
	base = strings.TrimRight(base, "/")
	hc := opts.HTTPClient
	if hc == nil {
		// why: zero Timeout -- callers control deadlines via context.Context
		hc = &http.Client{}
	}
	ua := opts.UserAgent
	if ua == "" {
		ua = fmt.Sprintf("tesote-sdk-go/%s (go/%s)", Version, runtime.Version())
	}
	now := opts.Now
	if now == nil {
		now = time.Now
	}
	r63 := opts.RandUint63
	if r63 == nil {
		r63 = secureUint63
	}
	sleeper := opts.Sleep
	if sleeper == nil {
		sleeper = ctxSleep
	}
	return &Client{
		apiKey:     opts.APIKey,
		baseURL:    base,
		userAgent:  ua,
		httpClient: hc,
		retry:      opts.RetryPolicy.normalized(),
		cache:      opts.Cache,
		logger:     opts.Logger,
		now:        now,
		rand63:     r63,
		sleep:      sleeper,
	}, nil
}

// BaseURL returns the configured base URL (without trailing slash).
func (c *Client) BaseURL() string { return c.baseURL }

// UserAgent returns the configured User-Agent string.
func (c *Client) UserAgent() string { return c.userAgent }

// LastRateLimit returns the most recently captured rate-limit headers.
// The bool is false if no response has been observed yet.
func (c *Client) LastRateLimit() (RateLimit, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastRateLimit, c.hasRateLimit
}

func (c *Client) captureRateLimit(h http.Header) {
	limit, lOK := atoiHeader(h, "X-RateLimit-Limit")
	remaining, rOK := atoiHeader(h, "X-RateLimit-Remaining")
	reset, _ := atoiHeader(h, "X-RateLimit-Reset")
	if !lOK && !rOK {
		return
	}
	c.mu.Lock()
	c.lastRateLimit = RateLimit{Limit: limit, Remaining: remaining, Reset: reset}
	c.hasRateLimit = true
	c.mu.Unlock()
}

func (c *Client) log(ev LogEvent) {
	if c.logger == nil {
		return
	}
	defer func() { _ = recover() }() // why: a misbehaving logger must not break a request
	c.logger(ev)
}

// Do issues a request and returns the response, performing retries per the
// configured RetryPolicy. Non-2xx responses become typed errors via MapAPIError.
func (c *Client) Do(ctx context.Context, method, path string, opts RequestOptions) (*Response, error) {
	if ctx == nil {
		return nil, &ConfigError{Field: "ctx", Message: "context.Context is required"}
	}
	method = strings.ToUpper(method)
	fullURL, err := c.buildURL(path, opts.Query)
	if err != nil {
		return nil, err
	}

	cacheKey := ""
	if opts.CacheTTL > 0 && method == http.MethodGet && c.cache != nil {
		cacheKey = c.cacheKey(method, fullURL)
		if cached, ok := c.cache.Get(cacheKey); ok {
			return &Response{
				Status:    cached.Status,
				Header:    cached.Header.Clone(),
				Body:      append([]byte(nil), cached.Body...),
				RequestID: cached.RequestID,
				Attempts:  1,
			}, nil
		}
	}

	var bodyBytes []byte
	bodyShape := ""
	if opts.Body != nil {
		bodyBytes, err = json.Marshal(opts.Body)
		if err != nil {
			return nil, &ConfigError{Field: "Body", Message: "cannot marshal: " + err.Error()}
		}
		bodyShape = fmt.Sprintf("json/%dB", len(bodyBytes))
	}

	idemKey := opts.IdempotencyKey
	if _, isMutating := mutatingMethods[method]; isMutating && idemKey == "" {
		idemKey, err = c.newIdempotencyKey()
		if err != nil {
			return nil, err
		}
	}

	summary := RequestSummary{
		Method:        method,
		Path:          path,
		Query:         opts.Query,
		BodyShape:     bodyShape,
		Authorization: RedactBearer(c.apiKey),
	}

	var lastErr error
	for attempt := 1; attempt <= c.retry.MaxAttempts; attempt++ {
		// why: re-create body reader each attempt; bytes.Reader is cheap and idempotent
		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}
		req, rerr := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
		if rerr != nil {
			return nil, &ConfigError{Field: "request", Message: rerr.Error()}
		}
		c.applyHeaders(req, method, idemKey, opts.ExtraHeaders, bodyBytes != nil)

		c.log(LogEvent{Phase: "request", Method: method, URL: fullURL, Attempt: attempt})

		resp, hErr := c.httpClient.Do(req)
		if hErr != nil {
			tErr := classifyTransportError(hErr, method, summary, attempt)
			lastErr = tErr
			c.log(LogEvent{Phase: "error", Method: method, URL: fullURL, Attempt: attempt, Err: tErr})
			// why: idemKey is always set for mutations (auto-generated when caller didn't pass one),
			// so retries are always safe; only non-mutating timeouts/network errors hit shouldRetry.
			if !shouldRetryTransport(tErr, method, idemKey != "") || attempt == c.retry.MaxAttempts {
				return nil, tErr
			}
			if sleepErr := c.sleep(ctx, c.backoff(attempt, 0)); sleepErr != nil {
				return nil, sleepErr
			}
			continue
		}

		body, rdErr := readAll(resp.Body)
		_ = resp.Body.Close()
		if rdErr != nil {
			tErr := &NetworkError{TransportError: &TransportError{
				Op: "read body", Message: rdErr.Error(), RequestSummary: summary, Attempts: attempt, Cause: rdErr,
			}}
			lastErr = tErr
			if attempt == c.retry.MaxAttempts {
				return nil, tErr
			}
			if sleepErr := c.sleep(ctx, c.backoff(attempt, 0)); sleepErr != nil {
				return nil, sleepErr
			}
			continue
		}

		c.captureRateLimit(resp.Header)
		requestID := resp.Header.Get("X-Request-Id")

		c.log(LogEvent{Phase: "response", Method: method, URL: fullURL, Status: resp.StatusCode, RequestID: requestID, Attempt: attempt})

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			out := &Response{Status: resp.StatusCode, Header: resp.Header, Body: body, RequestID: requestID, Attempts: attempt}
			if opts.Out != nil && len(body) > 0 {
				if jErr := json.Unmarshal(body, opts.Out); jErr != nil {
					return nil, &APIError{
						ErrorCode:      "DECODE_ERROR",
						Message:        "cannot decode response body: " + jErr.Error(),
						HTTPStatus:     resp.StatusCode,
						RequestID:      requestID,
						ResponseBody:   string(body),
						RequestSummary: summary,
						Attempts:       attempt,
					}
				}
			}
			if cacheKey != "" && opts.CacheTTL > 0 {
				c.cache.Set(cacheKey, CachedResponse{
					Status:    resp.StatusCode,
					Header:    resp.Header.Clone(),
					Body:      append([]byte(nil), body...),
					RequestID: requestID,
					ExpiresAt: c.now().Add(opts.CacheTTL),
				})
			}
			return out, nil
		}

		apiErr := MapAPIError(resp, body, summary)
		setAttempts(apiErr, attempt)
		lastErr = apiErr

		if isRetryableStatus(resp.StatusCode) && attempt < c.retry.MaxAttempts {
			retryAfter := parseRetryAfter(resp.Header)
			if sleepErr := c.sleep(ctx, c.backoff(attempt, retryAfter)); sleepErr != nil {
				return nil, sleepErr
			}
			continue
		}
		return nil, apiErr
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, errors.New("tesote: retry loop exited without result")
}
