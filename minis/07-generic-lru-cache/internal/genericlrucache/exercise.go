//go:build !solution && !reference

package genericlrucache

import (
	"container/list"
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	mu         sync.Mutex          // Protects all fields
	capacity   int                 // Maximum number of items
	defaultTTL time.Duration       // Default expiration time
	items      map[K]*list.Element // Key → list element
	evictList  *list.List          // Doubly-linked list (front = most recent)
}

type entry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt time.Time
}

// New implements the exercise.
//
// TODO: Implement this function
func New[K comparable, V any](capacity int, defaultTTL time.Duration) *Cache[K, V] {
	// TODO: Implement
	return nil
}

// Get implements the exercise.
//
// TODO: Implement this function
func (c *Cache[K, V]) Get(key K) (V, bool) {
	// TODO: Implement
	return *new(V), false
}

// Set implements the exercise.
//
// TODO: Implement this function
func (c *Cache[K, V]) Set(key K, val V) {
	// TODO: Implement
}

// SetWithTTL implements the exercise.
//
// TODO: Implement this function
func (c *Cache[K, V]) SetWithTTL(key K, val V, ttl time.Duration) {
	// TODO: Implement
}

// Len implements the exercise.
//
// TODO: Implement this function
func (c *Cache[K, V]) Len() int {
	// TODO: Implement
	return 0
}

// removeElement implements the exercise.
//
// TODO: Implement this function
func (c *Cache[K, V]) removeElement(elem *list.Element) {
	// TODO: Implement
}
