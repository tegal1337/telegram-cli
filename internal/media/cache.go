package media

import (
	"sync"
)

// Cache is an LRU cache for rendered media (terminal strings).
type Cache struct {
	mu       sync.RWMutex
	entries  map[string]string
	order    []string
	maxSize  int
}

// NewCache creates a new media cache.
func NewCache(maxSize int) *Cache {
	return &Cache{
		entries: make(map[string]string),
		maxSize: maxSize,
	}
}

// Get returns a cached rendered image, if available.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.entries[key]
	return val, ok
}

// Set stores a rendered image in the cache.
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.entries[key]; !exists {
		c.order = append(c.order, key)
	}
	c.entries[key] = value

	// Evict oldest entries if over capacity.
	for len(c.entries) > c.maxSize && len(c.order) > 0 {
		oldest := c.order[0]
		c.order = c.order[1:]
		delete(c.entries, oldest)
	}
}

// Clear empties the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]string)
	c.order = nil
}
