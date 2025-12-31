//go:build !solution
// +build !solution

package genericlrucache

// import (
// 	"container/list"
// 	"sync"
// 	"time"
// )

// ============================================================================
// LRU CACHE: Least Recently Used Eviction Policy
// ============================================================================
//
// LRU Algorithm:
// - Track access order using doubly-linked list
// - Most recently used items at front
// - Least recently used items at back
// - Evict from back when capacity exceeded
//
// Data Structure Design:
// - Map: key → list element (O(1) lookup)
// - Doubly-linked list: maintains recency order
// - Mutex: protects concurrent access
//
// Why Doubly-Linked List:
// - O(1) move to front (mark as recently used)
// - O(1) remove from back (evict LRU item)
// - O(1) remove arbitrary element (for updates)
//
// Time Complexity:
// - Get: O(1) average
// - Set: O(1) average
// - Space: O(capacity)
//
// ============================================================================

// TODO: Define Cache struct
// Generic type parameters: [K comparable, V any]
// Fields: mu (sync.Mutex), capacity (int), defaultTTL (time.Duration),
//         items (map[K]*list.Element), evictList (*list.List)
type Cache[K comparable, V any] struct {
	// TODO: Add fields here
	mu         sync.Mutex
	capacity   int
	defaultTTL time.Duration
	items      map[K]*list.Element
	evictList  *list.List
}

// TODO: Define entry struct
// Generic type parameters: [K comparable, V any]
// Fields: key (K), value (V), expiresAt (time.Time)
type entry[K comparable, V any] struct {
	// TODO: Add fields here
	key       K
	value     V
	expiresAt time.Time
}

// TODO: Implement New[K comparable, V any](capacity int, defaultTTL time.Duration) *Cache[K, V]
// Algorithm:
// 1. Create Cache with capacity and defaultTTL
// 2. Initialize items map
// 3. Initialize evictList with list.New()
func New[K comparable, V any](capacity int, defaultTTL time.Duration) *Cache[K, V] {
	// TODO: Implement Function
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


// TODO: Implement Get(key K) (V, bool)
// Algorithm:
// 1. Lock mutex
// 2. Look up key in items map
// 3. If not found, return zero value and false
// 4. Extract entry from list element
// 5. Check if expired (time.Now().After(expiresAt))
// 6. If expired, remove entry and return false
// 7. Move element to front (MoveToFront)
// 8. Return value and true
func (c *Cache[K, V]) Get(key K) (V, bool) {
	// TODO: Implement this method
	var zero V
	return zero, false
}

// TODO: Implement Set(key K, val V)
// Delegate to SetWithTTL(key, val, c.defaultTTL)
func (c *Cache[K, V]) Set(key K, val V) {
	// TODO: Implement this method
	c.SetWithTTL(key, val, c.defaultTTL)
}

// TODO: Implement SetWithTTL(key K, val V, ttl time.Duration)
// Algorithm:
// 1. Lock mutex
// 2. Calculate expiresAt = time.Now().Add(ttl) if ttl > 0
// 3. Check if key exists in items map
// 4. If exists:
//    - Move element to front
//    - Update entry value and expiresAt
//    - Return
// 5. If not exists:
//    - Create new entry
//    - PushFront to evictList
//    - Add to items map
//    - If len > capacity, remove Back element
func (c *Cache[K, V]) SetWithTTL(key K, val V, ttl time.Duration) {
	// TODO: Implement this method
}

// TODO: Implement Len() int
// Return evictList.Len() (with mutex lock)
func (c *Cache[K, V]) Len() int {
	// TODO: Implement this method
	return 0
}

// TODO: Implement removeElement(elem *list.Element)
// Algorithm:
// 1. Remove from evictList
// 2. Extract entry from element
// 3. Delete from items map using entry.key
func (c *Cache[K, V]) removeElement(elem *list.Element) {
	// TODO: Implement this method
}
