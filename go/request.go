package tesote

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var mutatingMethods = map[string]struct{}{
	http.MethodPost:   {},
	http.MethodPut:    {},
	http.MethodPatch:  {},
	http.MethodDelete: {},
}

// RequestOptions tunes a single transport call.
type RequestOptions struct {
	Query          map[string]string
	Body           any
	IdempotencyKey string
	CacheTTL       time.Duration
	ExtraHeaders   map[string]string
	// Out, if non-nil and the response is 2xx, decodes the JSON response into it.
	Out any
}

// Response is what Do returns on success.
type Response struct {
	Status    int
	Header    http.Header
	Body      []byte
	RequestID string
	Attempts  int
}

func (c *Client) buildURL(path string, query map[string]string) (string, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	full := c.baseURL + path
	if len(query) == 0 {
		return full, nil
	}
	values := url.Values{}
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		values.Set(k, query[k])
	}
	if strings.Contains(full, "?") {
		return full + "&" + values.Encode(), nil
	}
	return full + "?" + values.Encode(), nil
}

func (c *Client) applyHeaders(req *http.Request, method, idemKey string, extra map[string]string, hasBody bool) {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if hasBody {
		req.Header.Set("Content-Type", "application/json")
	}
	if _, isMutating := mutatingMethods[method]; isMutating && idemKey != "" {
		req.Header.Set("Idempotency-Key", idemKey)
	}
	for k, v := range extra {
		req.Header.Set(k, v)
	}
}

func (c *Client) cacheKey(method, fullURL string) string {
	// why: include method + URL; URL already contains sorted query from buildURL.
	return method + " " + fullURL
}

func (c *Client) newIdempotencyKey() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", &ConfigError{Field: "idempotencyKey", Message: "rand.Read: " + err.Error()}
	}
	// RFC 4122 v4
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
