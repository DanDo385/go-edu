//go:build reference

package genericlrucache

/*
Problem: Implement a thread-safe LRU cache with generics and TTL

Requirements:
1. O(1) Get and Set operations
2. Thread-safe (concurrent access from multiple goroutines)
3. LRU eviction when capacity is reached
4. Optional per-item TTL expiration
5. Generic over key and value types

Data Structure:
- Map: key → list element (O(1) lookup)
- Doubly-linked list: maintains recency order (front = most recent)
- Mutex: protects concurrent access

Time/Space Complexity:
- Get: O(1) average (map lookup + list move)
- Set: O(1) average (map insert + list append/evict)
- Space: O(capacity) for map + list

Algorithm: LRU Eviction Policy
- Track access order in doubly-linked list
- Most recent items at front
- Least recent items at back
- Evict from back when capacity exceeded

Why doubly-linked list:
- O(1) move to front (mark as recently used)
- O(1) remove from back (evict LRU item)
- O(1) remove arbitrary element (for updates)
*/

import (
	"container/list"
	"sync"
	"time"
)

// Cache is a generic LRU cache with TTL support.
// BREAKPOINT: Set breakpoint here to inspect cache initialization
type Cache[K comparable, V any] struct {
	mu         sync.Mutex          // Protects all fields
	capacity   int                 // Maximum number of items
	defaultTTL time.Duration       // Default expiration time
	items      map[K]*list.Element // Key → list element
	evictList  *list.List          // Doubly-linked list (front = most recent)
}

// entry holds the actual cached data.
// BREAKPOINT: Set breakpoint in methods that access entry to inspect cached data
type entry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt time.Time
}

// New creates an LRU cache with the given capacity and default TTL.
// BREAKPOINT: Set breakpoint here to step through cache creation
// DEBUG: Watch 'capacity' to verify initialization parameter
// DEBUG: Watch 'defaultTTL' to verify TTL configuration
func New[K comparable, V any](capacity int, defaultTTL time.Duration) *Cache[K, V] {
	// DEBUG: Inspect return value to verify all fields initialized
	return &Cache[K, V]{
		capacity:   capacity,
		defaultTTL: defaultTTL,
		items:      make(map[K]*list.Element),
		evictList:  list.New(),
	}
}

// Get retrieves a value by key.
//
// Algorithm:
// 1. Lock mutex for thread safety
// 2. Look up key in map (O(1))
// 3. Check if entry expired (lazy expiration)
// 4. Move to front of list (mark as recently used)
// 5. Return value
//
// BREAKPOINT: Set breakpoint at function entry to trace cache lookups
// DEBUG: Watch 'key' to see what's being requested
// DEBUG: Watch 'elem' to see if key exists in cache
// DEBUG: Watch 'ent.expiresAt' to check expiration time
func (c *Cache[K, V]) Get(key K) (V, bool) {
	// BREAKPOINT: Set breakpoint here before acquiring lock
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero V
	// DEBUG: 'zero' holds the zero value for type V

	// Look up in map
	// BREAKPOINT: Set breakpoint here to check map lookup
	// DEBUG: Watch 'ok' to see if key exists
	elem, ok := c.items[key]
	if !ok {
		// DEBUG: Cache miss - key not found
		return zero, false
	}

	// Extract entry
	// BREAKPOINT: Set breakpoint here to inspect entry data
	// DEBUG: Watch 'ent' to see cached value and metadata
	ent := elem.Value.(*entry[K, V])

	// Check TTL expiration
	// BREAKPOINT: Set breakpoint here to trace expiration logic
	// DEBUG: Watch 'time.Now()' vs 'ent.expiresAt' to verify expiration
	// DEBUG: Watch 'ent.expiresAt.IsZero()' to check if TTL is set
	if !ent.expiresAt.IsZero() && time.Now().After(ent.expiresAt) {
		// DEBUG: Entry expired - removing from cache
		// BREAKPOINT: Set breakpoint here to trace expiration cleanup
		c.removeElement(elem)
		return zero, false
	}

	// Move to front (mark as recently used)
	// BREAKPOINT: Set breakpoint here to trace LRU update
	// DEBUG: Watch evictList before and after to see list reordering
	c.evictList.MoveToFront(elem)

	// DEBUG: Watch 'ent.value' to see return value
	return ent.value, true
}

// Set inserts or updates a key-value pair with the default TTL.
// BREAKPOINT: Set breakpoint here to trace cache insertions
// DEBUG: Watch 'key' and 'val' to see what's being cached
func (c *Cache[K, V]) Set(key K, val V) {
	// DEBUG: Delegating to SetWithTTL with default TTL
	c.SetWithTTL(key, val, c.defaultTTL)
}

// SetWithTTL inserts or updates a key-value pair with custom TTL.
//
// Algorithm:
// 1. Calculate expiration time from TTL
// 2. Check if key already exists (update vs insert)
// 3. If exists: update value and move to front
// 4. If not exists: create entry, add to front
// 5. Evict LRU item if over capacity
//
// BREAKPOINT: Set breakpoint at function entry to trace cache updates
// DEBUG: Watch 'key', 'val', 'ttl' to see insertion parameters
// DEBUG: Watch 'expiresAt' to verify TTL calculation
func (c *Cache[K, V]) SetWithTTL(key K, val V, ttl time.Duration) {
	// BREAKPOINT: Set breakpoint here before acquiring lock
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiresAt time.Time
	// BREAKPOINT: Set breakpoint here to trace TTL calculation
	// DEBUG: Watch 'ttl' to see if expiration is enabled
	if ttl > 0 {
		// DEBUG: Watch 'expiresAt' to see calculated expiration time
		expiresAt = time.Now().Add(ttl)
	}
	// DEBUG: If ttl <= 0, expiresAt remains zero (no expiration)

	// Check if key already exists
	// BREAKPOINT: Set breakpoint here to check update vs insert
	// DEBUG: Watch 'ok' to determine update or insert path
	if elem, ok := c.items[key]; ok {
		// Update existing entry
		// BREAKPOINT: Set breakpoint here to trace cache updates
		// DEBUG: Watch 'elem' to see which entry is being updated
		c.evictList.MoveToFront(elem)
		ent := elem.Value.(*entry[K, V])
		// DEBUG: Watch 'ent' before update to see old value
		ent.value = val
		ent.expiresAt = expiresAt
		// DEBUG: Watch 'ent' after update to see new value
		return
	}

	// Add new entry
	// BREAKPOINT: Set breakpoint here to trace new insertions
	// DEBUG: Watch 'ent' to see new entry being created
	ent := &entry[K, V]{
		key:       key,
		value:     val,
		expiresAt: expiresAt,
	}
	// DEBUG: Watch 'elem' to see new list element
	elem := c.evictList.PushFront(ent)
	// DEBUG: Map now contains key → element mapping
	c.items[key] = elem

	// Evict if over capacity
	// BREAKPOINT: Set breakpoint here to trace eviction logic
	// DEBUG: Watch 'c.evictList.Len()' vs 'c.capacity'
	if c.evictList.Len() > c.capacity {
		// DEBUG: Cache over capacity - evicting LRU item
		// BREAKPOINT: Set breakpoint here to see which item is evicted
		// DEBUG: Watch 'c.evictList.Back()' to see LRU element
		c.removeElement(c.evictList.Back())
	}
}

// Len returns the current number of items.
// BREAKPOINT: Set breakpoint here to inspect cache size
// DEBUG: Watch return value to see current cache size
func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	// DEBUG: Watch 'c.evictList.Len()' for current item count
	return c.evictList.Len()
}

// removeElement removes an element from both the list and map.
//
// Algorithm:
// 1. Remove from doubly-linked list
// 2. Extract entry to get key
// 3. Delete from map using key
//
// BREAKPOINT: Set breakpoint here to trace removal logic
// DEBUG: Watch 'elem' to see which element is being removed
func (c *Cache[K, V]) removeElement(elem *list.Element) {
	// BREAKPOINT: Set breakpoint here to see list removal
	// DEBUG: Watch evictList before removal to see structure
	c.evictList.Remove(elem)

	// DEBUG: Watch 'ent' to see which entry is being removed
	ent := elem.Value.(*entry[K, V])

	// BREAKPOINT: Set breakpoint here to see map deletion
	// DEBUG: Watch 'ent.key' to see which key is being deleted
	// DEBUG: Watch 'c.items' before deletion to see map state
	delete(c.items, ent.key)
	// DEBUG: Watch 'c.items' after deletion to verify removal
}

/*
Alternatives & Trade-offs:

1. Use sync.RWMutex instead of sync.Mutex:
   Pros: Allows concurrent reads
   Cons: More complex; Get still needs write lock (moves element)
   Algorithm impact: No change to LRU logic

2. Sharded cache (multiple caches with hash-based routing):
   Pros: Reduces lock contention
   Cons: More complex; eviction is per-shard
   Algorithm: Hash key to determine shard, apply LRU per shard

3. Use map[K]*entry directly (no list):
   Pros: Simpler
   Cons: O(n) eviction (must find LRU item)
   Algorithm: Linear scan to find oldest timestamp

4. Active TTL cleanup (background goroutine):
   Pros: Reclaims memory proactively
   Cons: Adds complexity; goroutine overhead
   Algorithm: Periodic sweep to remove expired entries

5. Custom linked list (not container/list):
   Pros: Avoids interface{} and type assertions
   Cons: More code; error-prone
   Algorithm: Manual pointer management for prev/next
*/
