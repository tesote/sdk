package tesote

import (
	"net/http"
	"sync"
	"time"
)

// CacheBackend is the pluggable response-cache contract. Implementations must
// be safe for concurrent use.
type CacheBackend interface {
	Get(key string) (CachedResponse, bool)
	Set(key string, value CachedResponse)
}

// CachedResponse is what the transport stores in (and reads from) a CacheBackend.
type CachedResponse struct {
	Status    int
	Header    http.Header
	Body      []byte
	ExpiresAt time.Time
	RequestID string
}

// LRUCache is a small dep-free TTL LRU implementing CacheBackend.
type LRUCache struct {
	mu      sync.Mutex
	max     int
	entries map[string]*lruEntry
	order   []string
}

type lruEntry struct {
	value CachedResponse
}

// NewLRUCache returns a cache that holds up to max entries. max < 1 falls back to 256.
func NewLRUCache(max int) *LRUCache {
	if max < 1 {
		max = 256
	}
	return &LRUCache{max: max, entries: make(map[string]*lruEntry, max)}
}

// Get returns a cached entry if present and unexpired.
func (c *LRUCache) Get(key string) (CachedResponse, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		return CachedResponse{}, false
	}
	if !e.value.ExpiresAt.IsZero() && time.Now().After(e.value.ExpiresAt) {
		c.deleteLocked(key)
		return CachedResponse{}, false
	}
	c.touchLocked(key)
	return e.value, true
}

// Set inserts or updates a cache entry, evicting the LRU entry if at capacity.
func (c *LRUCache) Set(key string, value CachedResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.entries[key]; ok {
		c.entries[key].value = value
		c.touchLocked(key)
		return
	}
	c.entries[key] = &lruEntry{value: value}
	c.order = append(c.order, key)
	for len(c.order) > c.max {
		evict := c.order[0]
		c.order = c.order[1:]
		delete(c.entries, evict)
	}
}

func (c *LRUCache) touchLocked(key string) {
	for i, k := range c.order {
		if k == key {
			c.order = append(c.order[:i], c.order[i+1:]...)
			break
		}
	}
	c.order = append(c.order, key)
}

func (c *LRUCache) deleteLocked(key string) {
	delete(c.entries, key)
	for i, k := range c.order {
		if k == key {
			c.order = append(c.order[:i], c.order[i+1:]...)
			return
		}
	}
}
