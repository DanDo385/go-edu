//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
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
// TODO: implement New.
func New[K comparable, V any](capacity int, defaultTTL time.Duration) *Cache[K, V] {
	panic("TODO: implement")
}
// TODO: implement Get.
func (c *Cache[K, V]) Get(key K) (V, bool) { panic("TODO: implement") }
// TODO: implement Set.
func (c *Cache[K, V]) Set(key K, val V) { panic("TODO: implement") }
// TODO: implement SetWithTTL.
func (c *Cache[K, V]) SetWithTTL(key K, val V, ttl time.Duration) { panic("TODO: implement") }
// TODO: implement Len.
func (c *Cache[K, V]) Len() int { panic("TODO: implement") }
// TODO: implement removeElement.
func (c *Cache[K, V]) removeElement(elem *list.Element) { panic("TODO: implement") }
