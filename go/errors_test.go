package tesote

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func mkResp(status int, headers map[string]string) *http.Response {
	h := http.Header{}
	for k, v := range headers {
		h.Set(k, v)
	}
	return &http.Response{StatusCode: status, Header: h}
}

func TestMapAPIError_KnownCodes(t *testing.T) {
	cases := []struct {
		name      string
		status    int
		body      string
		sentinel  error
		wantField func(t *testing.T, err error)
	}{
		{"unauthorized", 401, `{"error":"nope","error_code":"UNAUTHORIZED"}`, ErrUnauthorized, func(t *testing.T, err error) {
			var e *UnauthorizedError
			if !errors.As(err, &e) {
				t.Fatalf("not *UnauthorizedError: %T", err)
			}
		}},
		{"api key revoked", 401, `{"error_code":"API_KEY_REVOKED"}`, ErrAPIKeyRevoked, func(t *testing.T, err error) {
			var e *APIKeyRevokedError
			if !errors.As(err, &e) {
				t.Fatalf("not *APIKeyRevokedError: %T", err)
			}
		}},
		{"workspace suspended", 403, `{"error_code":"WORKSPACE_SUSPENDED"}`, ErrWorkspaceSuspended, func(t *testing.T, err error) {
			var e *WorkspaceSuspendedError
			if !errors.As(err, &e) {
				t.Fatalf("not *WorkspaceSuspendedError: %T", err)
			}
		}},
		{"account disabled", 403, `{"error_code":"ACCOUNT_DISABLED"}`, ErrAccountDisabled, func(t *testing.T, err error) {
			var e *AccountDisabledError
			if !errors.As(err, &e) {
				t.Fatalf("not *AccountDisabledError: %T", err)
			}
		}},
		{"history forbidden", 403, `{"error_code":"HISTORY_SYNC_FORBIDDEN"}`, ErrHistorySyncForbidden, func(t *testing.T, err error) {
			var e *HistorySyncForbiddenError
			if !errors.As(err, &e) {
				t.Fatalf("not *HistorySyncForbiddenError: %T", err)
			}
		}},
		{"mutation conflict", 409, `{"error_code":"MUTATION_CONFLICT"}`, ErrMutationDuringPagination, func(t *testing.T, err error) {
			var e *MutationDuringPaginationError
			if !errors.As(err, &e) {
				t.Fatalf("not *MutationDuringPaginationError: %T", err)
			}
		}},
		{"unprocessable", 422, `{"error_code":"UNPROCESSABLE_CONTENT"}`, ErrUnprocessableContent, func(t *testing.T, err error) {
			var e *UnprocessableContentError
			if !errors.As(err, &e) {
				t.Fatalf("not *UnprocessableContentError: %T", err)
			}
		}},
		{"date range", 422, `{"error_code":"INVALID_DATE_RANGE"}`, ErrInvalidDateRange, func(t *testing.T, err error) {
			var e *InvalidDateRangeError
			if !errors.As(err, &e) {
				t.Fatalf("not *InvalidDateRangeError: %T", err)
			}
		}},
		{"rate limit", 429, `{"error_code":"RATE_LIMIT_EXCEEDED","retry_after":42}`, ErrRateLimitExceeded, func(t *testing.T, err error) {
			var e *RateLimitExceededError
			if !errors.As(err, &e) {
				t.Fatalf("not *RateLimitExceededError: %T", err)
			}
			if e.RetryAfter != 42 {
				t.Errorf("retry_after = %d", e.RetryAfter)
			}
		}},
		{"service unavailable", 503, `{"error":"paused"}`, ErrServiceUnavailable, func(t *testing.T, err error) {
			var e *ServiceUnavailableError
			if !errors.As(err, &e) {
				t.Fatalf("not *ServiceUnavailableError: %T", err)
			}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp := mkResp(tc.status, map[string]string{"X-Request-Id": "req-x"})
			err := MapAPIError(resp, []byte(tc.body), RequestSummary{Method: "GET", Path: "/x", Authorization: "Bearer 1234"})
			if !errors.Is(err, tc.sentinel) {
				t.Errorf("errors.Is(%v) = false", tc.sentinel)
			}
			tc.wantField(t, err)
			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Fatal("not unwrappable into *APIError")
			}
			if apiErr.RequestID != "req-x" {
				t.Errorf("request id = %q", apiErr.RequestID)
			}
			if apiErr.HTTPStatus != tc.status {
				t.Errorf("status = %d", apiErr.HTTPStatus)
			}
			if !strings.Contains(apiErr.RequestSummary.Authorization, "1234") {
				t.Errorf("auth not redacted to last4: %q", apiErr.RequestSummary.Authorization)
			}
		})
	}
}

func TestMapAPIError_RetryAfterHeaderWins(t *testing.T) {
	resp := mkResp(429, map[string]string{"Retry-After": "7"})
	err := MapAPIError(resp, []byte(`{"error_code":"RATE_LIMIT_EXCEEDED","retry_after":1}`), RequestSummary{})
	var rl *RateLimitExceededError
	if !errors.As(err, &rl) {
		t.Fatalf("not rate limit: %T", err)
	}
	if rl.RetryAfter != 7 {
		t.Errorf("retry_after = %d, want 7 (header wins)", rl.RetryAfter)
	}
}

func TestMapAPIError_StatusFallback(t *testing.T) {
	resp := mkResp(409, nil)
	err := MapAPIError(resp, nil, RequestSummary{})
	var e *MutationDuringPaginationError
	if !errors.As(err, &e) {
		t.Fatalf("expected *MutationDuringPaginationError, got %T", err)
	}
	if e.ErrorCode != "MUTATION_CONFLICT" {
		t.Errorf("error_code = %q", e.ErrorCode)
	}
}

func TestMapAPIError_UnknownStatus(t *testing.T) {
	resp := mkResp(418, nil)
	err := MapAPIError(resp, nil, RequestSummary{})
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.ErrorCode != "HTTP_418" {
		t.Errorf("fallback code = %q", apiErr.ErrorCode)
	}
}

func TestRedactBearer_NeverLeaks(t *testing.T) {
	full := "tesote-key-abcdef123456"
	red := RedactBearer(full)
	if strings.Contains(red, "abcdef") {
		t.Errorf("leaked secret: %q", red)
	}
	if !strings.HasSuffix(red, "3456") {
		t.Errorf("missing last4: %q", red)
	}
}

func TestErrorString_Greppable(t *testing.T) {
	resp := mkResp(429, map[string]string{"X-Request-Id": "req-1"})
	err := MapAPIError(resp, []byte(`{"error":"slow down","error_code":"RATE_LIMIT_EXCEEDED"}`), RequestSummary{Method: "POST", Path: "/v3/x"})
	s := err.Error()
	if !strings.Contains(s, "RATE_LIMIT_EXCEEDED") || !strings.Contains(s, "req-1") {
		t.Errorf("error string missing context: %q", s)
	}
}

func TestConfigError_IsSentinel(t *testing.T) {
	err := &ConfigError{Field: "APIKey", Message: "missing"}
	if !errors.Is(err, ErrConfig) {
		t.Errorf("expected errors.Is(ErrConfig)")
	}
}

func TestEndpointRemovedError_IsSentinel(t *testing.T) {
	err := &EndpointRemovedError{Method: "v1.old", Replacement: "v3.new"}
	if !errors.Is(err, ErrEndpointRemoved) {
		t.Errorf("expected errors.Is(ErrEndpointRemoved)")
	}
}
