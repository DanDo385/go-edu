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

Why Go is well-suited:
- Generics: Type-safe without code duplication
- container/list: Built-in doubly-linked list
- sync.Mutex: Simple, efficient locking
- time.Time: First-class time support

Compared to other languages:
- Python: functools.lru_cache is convenient but less flexible
- Java: LinkedHashMap similar, but type erasure complicates generics
- Rust: Similar approach, but lifetimes add complexity
*/

package exercise

import (
	"container/list"
	"sync"
	"time"
)

// Cache is a generic LRU cache with TTL support.
type Cache[K comparable, V any] struct {
	mu          sync.Mutex              // Protects all fields
	capacity    int                     // Maximum number of items
	defaultTTL  time.Duration           // Default expiration time
	items       map[K]*list.Element     // Key → list element
	evictList   *list.List              // Doubly-linked list (front = most recent)
}

// entry holds the actual cached data.
type entry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt time.Time
}

// New creates an LRU cache with the given capacity and default TTL.
func New[K comparable, V any](capacity int, defaultTTL time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		capacity:   capacity,
		defaultTTL: defaultTTL,
		items:      make(map[K]*list.Element),
		evictList:  list.New(),
	}
}

// Get retrieves a value by key.
//
// Go Concepts:
// - Generics: [K comparable, V any] type parameters
// - Zero values: var zero V returns the zero value for type V
// - Mutex locking: defer c.mu.Unlock() ensures unlock even on early return
// - Type assertions: elem.Value.(*entry[K, V]) converts interface{} to concrete type
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero V

	// Look up in map
	elem, ok := c.items[key]
	if !ok {
		return zero, false
	}

	// Extract entry
	ent := elem.Value.(*entry[K, V])

	// Check TTL expiration
	if !ent.expiresAt.IsZero() && time.Now().After(ent.expiresAt) {
		// Expired: remove and return not found
		c.removeElement(elem)
		return zero, false
	}

	// Move to front (mark as recently used)
	c.evictList.MoveToFront(elem)

	return ent.value, true
}

// Set inserts or updates a key-value pair with the default TTL.
func (c *Cache[K, V]) Set(key K, val V) {
	c.SetWithTTL(key, val, c.defaultTTL)
}

// SetWithTTL inserts or updates a key-value pair with custom TTL.
//
// Go Concepts:
// - Map insert/update: items[key] = value
// - List operations: PushFront, Remove, Back
// - Zero time: time.Time{}.IsZero() == true (no expiration)
func (c *Cache[K, V]) SetWithTTL(key K, val V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	// Check if key already exists
	if elem, ok := c.items[key]; ok {
		// Update existing entry
		c.evictList.MoveToFront(elem)
		ent := elem.Value.(*entry[K, V])
		ent.value = val
		ent.expiresAt = expiresAt
		return
	}

	// Add new entry
	ent := &entry[K, V]{
		key:       key,
		value:     val,
		expiresAt: expiresAt,
	}
	elem := c.evictList.PushFront(ent)
	c.items[key] = elem

	// Evict if over capacity
	if c.evictList.Len() > c.capacity {
		c.removeElement(c.evictList.Back())
	}
}

// Len returns the current number of items.
func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.evictList.Len()
}

// removeElement removes an element from both the list and map.
func (c *Cache[K, V]) removeElement(elem *list.Element) {
	c.evictList.Remove(elem)
	ent := elem.Value.(*entry[K, V])
	delete(c.items, ent.key)
}

/*
Alternatives & Trade-offs:

1. Use sync.RWMutex instead of sync.Mutex:
   Pros: Allows concurrent reads
   Cons: More complex; Get still needs write lock (moves element)

2. Sharded cache (multiple caches with hash-based routing):
   Pros: Reduces lock contention
   Cons: More complex; eviction is per-shard

3. Use map[K]*entry directly (no list):
   Pros: Simpler
   Cons: O(n) eviction (must find LRU item)

4. Active TTL cleanup (background goroutine):
   Pros: Reclaims memory proactively
   Cons: Adds complexity; goroutine overhead

5. Custom linked list (not container/list):
   Pros: Avoids interface{} and type assertions
   Cons: More code; error-prone

Go vs X:

Go vs Python:
  from functools import lru_cache
  @lru_cache(maxsize=128)
  def expensive_func(arg):
      return result
  Pros: One-liner decorator
  Cons: Function-level only; no custom TTL; not thread-safe by default
  Go: Explicit cache object; more control

Go vs Java:
  Map<K, V> cache = Collections.synchronizedMap(
      new LinkedHashMap<>(capacity, 0.75f, true)
  );
  Pros: Built-in LRU support
  Cons: No TTL; synchronizedMap has coarse locking
  Go: Cleaner generics; explicit TTL support

Go vs Rust:
  use lru::LruCache;
  let mut cache: LruCache<String, i32> = LruCache::new(128);
  Pros: Zero-cost abstractions
  Cons: No built-in TTL; Arc<Mutex<LruCache>> for thread safety
  Go: Simpler API; built-in mutex patterns
*/
