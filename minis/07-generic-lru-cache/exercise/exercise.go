//go:build !solution
// +build !solution

package exercise

import (
	"time"
)

// Cache is a generic LRU cache with optional TTL support.
// K must be comparable (can be used as map key).
// V can be any type.
type Cache[K comparable, V any] struct {
	// TODO: Add fields
}

// New creates an LRU cache with the given capacity and default TTL.
// If defaultTTL is 0, items never expire.
func New[K comparable, V any](capacity int, defaultTTL time.Duration) *Cache[K, V] {
	// TODO: implement
	return nil
}

// Get retrieves a value by key.
// Returns (value, true) if found and not expired.
// Returns (zero value, false) if not found or expired.
// Moves the accessed item to the front (most recent).
func (c *Cache[K, V]) Get(key K) (V, bool) {
	// TODO: implement
	var zero V
	return zero, false
}

// Set inserts or updates a key-value pair with the default TTL.
// If the cache is at capacity, evicts the least recently used item.
func (c *Cache[K, V]) Set(key K, val V) {
	// TODO: implement
}

// SetWithTTL inserts or updates a key-value pair with a custom TTL.
// If ttl is 0, the item never expires.
func (c *Cache[K, V]) SetWithTTL(key K, val V, ttl time.Duration) {
	// TODO: implement
}

// Len returns the current number of items in the cache.
func (c *Cache[K, V]) Len() int {
	// TODO: implement
	return 0
}
