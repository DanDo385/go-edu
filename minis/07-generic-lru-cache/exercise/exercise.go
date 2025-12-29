//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "container/list" for doubly-linked list implementation
// - "sync" for Mutex to protect concurrent access
// - "time" for TTL (time-to-live) expiration support

// Doubly-linked list for maintaining recency order

// Duration and Time types for TTL expiration

// ============================================================================
// LRU CACHE WITH GENERICS: Efficient Caching with Eviction Policy
// ============================================================================
//
// LRU (Least Recently Used) is a cache eviction policy:
// - When cache is full, evict the least recently used item
// - "Recently used" = accessed (read or written)
// - Guarantees O(1) Get and Set operations
//
// Why LRU caching:
// - Bounded memory usage (capacity limit)
// - Automatic eviction (no manual cleanup needed)
// - Good hit rate for temporal locality (recent items likely accessed again)
// - Common in databases, OS page caches, CDNs, application caches
//
// Data structure design:
// - Map: key → list element (O(1) lookup)
// - Doubly-linked list: maintains recency order
//   * Front = most recently used
//   * Back = least recently used
// - Mutex: protects concurrent access
//
// Why doubly-linked list:
// - O(1) move to front (mark as recently used)
// - O(1) remove from back (evict LRU item)
// - O(1) remove arbitrary element (for updates)
// - container/list provides this built-in
//
// TTL (Time-To-Live) support:
// - Items can expire after a duration
// - Checked on Get (lazy expiration)
// - Zero TTL means no expiration
//
// Memory considerations:
// - Map overhead: ~48 bytes per entry (key, value, hash, pointers)
// - List overhead: ~48 bytes per node (prev, next, value pointers)
// - Total: ~96 bytes + size(key) + size(value) per cached item
// - Mutex: ~8 bytes (lock word)
//
// ============================================================================

// ============================================================================
// Exercise 1: Define Cache Type
// ============================================================================

// Cache is a generic LRU cache with optional TTL support.
//
// TODO: Define the Cache struct with:
// 1. Type parameters: K (key type) and V (value type)
//    - K must be comparable (can be used as map key)
//    - V can be any type
// 2. Fields:
//    - mu: sync.Mutex to protect concurrent access
//    - capacity: int (maximum number of items)
//    - defaultTTL: time.Duration (default expiration time)
//    - items: map[K]*list.Element (key to list element mapping)
//    - evictList: *list.List (doubly-linked list for recency)
//
// Generic type constraints:
// - "comparable" means the type can be compared with == and !=
// - Required for map keys in Go
// - "any" means any type (no constraints)
//
// Why pointer receiver for methods:
// - Cache contains a mutex (must not be copied)
// - Methods modify the cache state
// - All methods must operate on the same instance
//
// type Cache[K comparable, V any] struct {
//     mu         sync.Mutex       // Protects all fields below
//     capacity   int              // Maximum number of items
//     defaultTTL time.Duration    // Default expiration time (0 = no expiration)
//     items      map[K]*list.Element  // Key → list element
//     evictList  *list.List       // Doubly-linked list (front = most recent)
// }

// ============================================================================
// Exercise 2: Define entry Type
// ============================================================================

// entry holds the actual cached data within a list element.
//
// TODO: Define the entry struct with:
// 1. Type parameters: K comparable, V any
// 2. Fields:
//    - key: K (the cache key, needed for eviction)
//    - value: V (the cached value)
//    - expiresAt: time.Time (when this entry expires, zero = no expiration)
//
// Why store the key in the entry:
// - When evicting from the back of the list, we need the key to delete from map
// - List elements only contain the entry, not the key directly
//
// type entry[K comparable, V any] struct {
//     key       K
//     value     V
//     expiresAt time.Time
// }

// ============================================================================
// Exercise 3: Implement New
// ============================================================================

// New creates an LRU cache with the given capacity and default TTL.
//
// TODO: Implement cache initialization:
// 1. Create and return a pointer to Cache
// 2. Set capacity and defaultTTL
// 3. Initialize items map with make()
// 4. Initialize evictList with list.New()
//
// Memory allocation:
// - Cache struct allocated on heap (escapes because we return pointer)
// - Map allocated on heap (make(map[K]*list.Element))
// - List allocated on heap (list.New())
//
// Zero values:
// - sync.Mutex zero value is unlocked (ready to use)
// - time.Duration zero value is 0 (no expiration)
//
// func New[K comparable, V any](capacity int, defaultTTL time.Duration) *Cache[K, V] {
//     return &Cache[K, V]{
//         capacity:   capacity,
//         defaultTTL: defaultTTL,
//         items:      make(map[K]*list.Element),
//         evictList:  list.New(),
//     }
// }

// ============================================================================
// Exercise 4: Implement Get
// ============================================================================

// Get retrieves a value by key.
//
// TODO: Implement thread-safe cache lookup:
// 1. Lock the mutex (defer unlock for safety)
// 2. Look up key in items map
// 3. If not found, return zero value and false
// 4. Extract entry from list element (type assertion)
// 5. Check if entry has expired (time.Now().After(expiresAt))
// 6. If expired, remove the entry and return not found
// 7. Move element to front (mark as recently used)
// 8. Return value and true
//
// Returns:
//   - (value, true) if found and not expired
//   - (zero value, false) if not found or expired
//
// Lazy expiration:
// - We only check expiration on Get, not proactively
// - Saves resources (no background goroutine needed)
// - Trade-off: Expired items may linger until accessed
//
// Moving to front:
// - list.MoveToFront(elem) is O(1)
// - Adjusts prev/next pointers in the linked list
// - Makes this element the most recently used
//
// Type assertion:
// - elem.Value is interface{} (Go's container/list is not generic)
// - Must assert to *entry[K, V] to access fields
// - Panic if wrong type (bug in our code)
//
// func (c *Cache[K, V]) Get(key K) (V, bool) {
//     c.mu.Lock()
//     defer c.mu.Unlock()
//
//     var zero V  // Zero value for type V
//
//     // Look up in map
//     elem, ok := c.items[key]
//     if !ok {
//         return zero, false
//     }
//
//     // Extract entry
//     ent := elem.Value.(*entry[K, V])
//
//     // Check TTL expiration
//     if !ent.expiresAt.IsZero() && time.Now().After(ent.expiresAt) {
//         // Expired: remove and return not found
//         c.removeElement(elem)
//         return zero, false
//     }
//
//     // Move to front (mark as recently used)
//     c.evictList.MoveToFront(elem)
//
//     return ent.value, true
// }

// ============================================================================
// Exercise 5: Implement Set
// ============================================================================

// Set inserts or updates a key-value pair with the default TTL.
//
// TODO: Implement by delegating to SetWithTTL:
// - Call c.SetWithTTL(key, val, c.defaultTTL)
//
// This provides a convenient API (users don't need to pass TTL every time)
//
// func (c *Cache[K, V]) Set(key K, val V) {
//     c.SetWithTTL(key, val, c.defaultTTL)
// }

// ============================================================================
// Exercise 6: Implement SetWithTTL
// ============================================================================

// SetWithTTL inserts or updates a key-value pair with a custom TTL.
//
// TODO: Implement thread-safe cache insertion/update:
// 1. Lock the mutex
// 2. Calculate expiresAt:
//    - If ttl > 0: expiresAt = time.Now().Add(ttl)
//    - If ttl == 0: expiresAt = zero time (never expires)
// 3. Check if key already exists in items map
// 4. If exists (update case):
//    - Move element to front
//    - Update entry's value and expiresAt
//    - Return
// 5. If not exists (insert case):
//    - Create new entry
//    - Push to front of list
//    - Store element in map
//    - Check if over capacity
//    - If over capacity, remove back element (LRU eviction)
//
// Time handling:
// - time.Time{} is the zero value (IsZero() returns true)
// - time.Now() returns current time
// - time.Now().Add(duration) returns future time
//
// LRU eviction:
// - evictList.Back() returns least recently used element
// - removeElement() removes from both list and map
// - Eviction happens AFTER insertion (capacity exceeded by 1, then evict)
//
// func (c *Cache[K, V]) SetWithTTL(key K, val V, ttl time.Duration) {
//     c.mu.Lock()
//     defer c.mu.Unlock()
//
//     var expiresAt time.Time
//     if ttl > 0 {
//         expiresAt = time.Now().Add(ttl)
//     }
//
//     // Check if key already exists (update case)
//     if elem, ok := c.items[key]; ok {
//         c.evictList.MoveToFront(elem)
//         ent := elem.Value.(*entry[K, V])
//         ent.value = val
//         ent.expiresAt = expiresAt
//         return
//     }
//
//     // Insert new entry
//     ent := &entry[K, V]{
//         key:       key,
//         value:     val,
//         expiresAt: expiresAt,
//     }
//     elem := c.evictList.PushFront(ent)
//     c.items[key] = elem
//
//     // Evict if over capacity
//     if c.evictList.Len() > c.capacity {
//         c.removeElement(c.evictList.Back())
//     }
// }

// ============================================================================
// Exercise 7: Implement Len
// ============================================================================

// Len returns the current number of items in the cache.
//
// TODO: Implement thread-safe length check:
// 1. Lock the mutex
// 2. Return the length of evictList
//
// Why check list length instead of map length:
// - Both should be the same (invariant we maintain)
// - List is the source of truth for eviction
// - Could also use len(c.items)
//
// func (c *Cache[K, V]) Len() int {
//     c.mu.Lock()
//     defer c.mu.Unlock()
//     return c.evictList.Len()
// }

// ============================================================================
// Exercise 8: Implement removeElement (Helper)
// ============================================================================

// removeElement removes an element from both the list and map.
//
// TODO: Implement removal from both data structures:
// 1. Remove from evictList
// 2. Extract entry from element
// 3. Delete from items map using entry's key
//
// This is a helper to maintain consistency between map and list.
// Must be called while holding the lock (not thread-safe on its own).
//
// func (c *Cache[K, V]) removeElement(elem *list.Element) {
//     c.evictList.Remove(elem)
//     ent := elem.Value.(*entry[K, V])
//     delete(c.items, ent.key)
// }

// ============================================================================
// Comparison with Other Approaches
// ============================================================================
//
// Alternative 1: Use sync.RWMutex instead of sync.Mutex
//   Pros: Allows concurrent reads (multiple Get calls)
//   Cons: Get must still acquire write lock (moves element to front)
//         More complex; minimal benefit for this use case
//
// Alternative 2: Sharded cache (multiple caches with hash-based routing)
//   cache := [16]*Cache  // 16 shards
//   shard := hash(key) % 16
//   Pros: Reduces lock contention (each shard has its own lock)
//   Cons: More complex; eviction is per-shard (not global LRU)
//         Need good hash function for even distribution
//
// Alternative 3: Use map[K]*entry directly (no list)
//   Pros: Simpler (one data structure)
//   Cons: O(n) eviction (must scan all entries to find LRU)
//         Need to track access time in entry
//         Not practical for large caches
//
// Alternative 4: Active TTL cleanup (background goroutine)
//   go func() {
//       ticker := time.NewTicker(time.Minute)
//       for range ticker.C {
//           cache.removeExpired()
//       }
//   }()
//   Pros: Reclaims memory proactively
//   Cons: Adds complexity; goroutine overhead
//         Needs clean shutdown mechanism
//
// Alternative 5: Custom linked list (not container/list)
//   type node[K, V] struct {
//       key, val K, V
//       prev, next *node[K, V]
//   }
//   Pros: Avoids interface{} and type assertions (fully generic)
//   Cons: More code; error-prone (manual pointer management)
//
// ============================================================================
// Go vs Other Languages
// ============================================================================
//
// Go vs Python (functools.lru_cache):
//   from functools import lru_cache
//   @lru_cache(maxsize=128)
//   def expensive_function(arg):
//       return compute(arg)
//
//   Pros: One-liner decorator
//   Cons: Function-level only (can't cache arbitrary key-value pairs)
//         No custom TTL
//         Not thread-safe by default (need to wrap with lock)
//   Go: Explicit cache object; more control
//       Thread-safe out of the box
//       Supports TTL and custom eviction
//
// Go vs Java (LinkedHashMap):
//   Map<String, Integer> cache = Collections.synchronizedMap(
//       new LinkedHashMap<String, Integer>(capacity, 0.75f, true) {
//           protected boolean removeEldestEntry(Map.Entry eldest) {
//               return size() > capacity;
//           }
//       }
//   );
//
//   Pros: Built-in LRU support (accessOrder = true)
//   Cons: No TTL support
//         synchronizedMap has coarse-grained locking (poor concurrency)
//         Verbose syntax (anonymous inner class)
//   Go: Cleaner generics syntax
//       Explicit TTL support
//       Same locking granularity but clearer API
//
// Go vs Rust (lru crate):
//   use lru::LruCache;
//   let mut cache: LruCache<String, i32> = LruCache::new(128);
//   cache.put("key".to_string(), 42);
//   let val = cache.get(&"key".to_string());
//
//   Pros: Zero-cost abstractions
//         No garbage collection overhead
//   Cons: No built-in TTL (need external crate)
//         Arc<Mutex<LruCache>> for thread safety (more verbose)
//         Lifetimes complicate API
//   Go: Simpler API; built-in mutex patterns
//       Garbage collector handles memory (less manual management)
//
// Go vs JavaScript (lru-cache npm package):
//   const LRU = require('lru-cache');
//   const cache = new LRU({ max: 128, ttl: 1000 * 60 * 5 });
//   cache.set('key', 'value');
//   const val = cache.get('key');
//
//   Pros: Similar API; built-in TTL
//   Cons: Single-threaded (no true concurrency)
//         External package (not in stdlib)
//         No type safety without TypeScript
//   Go: Built-in concurrency support
//       Type-safe with generics
//       No external dependencies
//
// ============================================================================
// Advanced Topics
// ============================================================================
//
// 1. Eviction Policies (besides LRU):
//    - LFU (Least Frequently Used): Track access count instead of recency
//    - FIFO (First In First Out): Simple queue, no recency tracking
//    - Random: Pick random item to evict (surprisingly effective!)
//    - ARC (Adaptive Replacement Cache): Balances recency and frequency
//
// 2. Concurrency Optimizations:
//    - Sharding: Reduce lock contention with multiple cache instances
//    - Lock-free: Use atomic operations (complex, limited applicability)
//    - Copy-on-write: Readers don't need locks (but high memory overhead)
//
// 3. Memory Management:
//    - Weak references: Allow GC to collect cached items under pressure
//      (Go doesn't have weak references, but can use finalizers)
//    - Size-based eviction: Evict based on memory size, not item count
//    - Compression: Compress cached values to save memory
//
// 4. Persistence:
//    - Write-through: Write to cache and backing store synchronously
//    - Write-back: Write to cache, flush to backing store async
//    - Read-through: On cache miss, load from backing store
//
// ============================================================================
// After Implementation
// ============================================================================
//
// Testing your implementation:
// - Run: go test -v ./...
// - Run with race detector: go test -race ./...
// - Test edge cases: capacity=0, capacity=1, TTL=0
// - Test concurrency: multiple goroutines accessing cache
// - Test expiration: verify TTL eviction works
//
// Performance tips:
// - Benchmark with different capacities
// - Measure hit rate for your workload
// - Consider sharding if lock contention is high
// - Profile with pprof to find bottlenecks
//
// Common mistakes:
// - Forgetting to lock in removeElement (causes race)
// - Not checking expiration in Get (returns stale data)
// - Type assertion without checking (panic on wrong type)
// - Copying Cache struct (breaks mutex)
// - Not storing key in entry (can't evict from map)
