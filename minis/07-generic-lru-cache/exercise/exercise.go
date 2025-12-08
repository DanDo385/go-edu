//go:build !solution
// +build !solution

package exercise

import (
	"container/list"
	"sync"
	"time"
)

// Cache is a generic LRU cache with optional TTL support.
// K must be comparable (can be used as map key).
// V can be any type.
type Cache[K comparable, V any] struct {
	mu         sync.Mutex
	capacity   int
	defaultTTL time.Duration
	items      map[K]*list.Element
	evictList  *list.List
}

// New creates an LRU cache with the given capacity and default TTL.
// If defaultTTL is 0, items never expire.
func New[K comparable, V any](capacity int, defaultTTL time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		capacity:   capacity,
		defaultTTL: defaultTTL,
		items:      make(map[K]*list.Element),
		evictList:  list.New(),
	}
}

// Get retrieves a value by key.
// Returns (value, true) if found and not expired.
// Returns (zero value, false) if not found or expired.
// Moves the accessed item to the front (most recent).
func (c *Cache[K, V]) Get(key K) (V, bool) {
	var zero V
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		return zero, false
	}

	ent := elem.Value.(*entry[K, V])
	if !ent.expiresAt.IsZero() && time.Now().After(ent.expiresAt) {
		c.removeElement(elem)
		return zero, false
	}

	c.evictList.MoveToFront(elem)
	return ent.value, true
}

// Set inserts or updates a key-value pair with the default TTL.
// If the cache is at capacity, evicts the least recently used item.
func (c *Cache[K, V]) Set(key K, val V) {
	c.SetWithTTL(key, val, c.defaultTTL)
}

// SetWithTTL inserts or updates a key-value pair with a custom TTL.
// If ttl is 0, the item never expires.
func (c *Cache[K, V]) SetWithTTL(key K, val V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	if elem, ok := c.items[key]; ok {
		c.evictList.MoveToFront(elem)
		ent := elem.Value.(*entry[K, V])
		ent.value = val
		ent.expiresAt = expiresAt
		return
	}

	ent := &entry[K, V]{key: key, value: val, expiresAt: expiresAt}
	elem := c.evictList.PushFront(ent)
	c.items[key] = elem

	if c.evictList.Len() > c.capacity {
		c.removeElement(c.evictList.Back())
	}
}

// Len returns the current number of items in the cache.
func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.evictList.Len()
}

// entry holds the cached data.
type entry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt time.Time
}

func (c *Cache[K, V]) removeElement(elem *list.Element) {
	c.evictList.Remove(elem)
	ent := elem.Value.(*entry[K, V])
	delete(c.items, ent.key)
}
