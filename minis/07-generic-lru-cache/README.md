# Project 07: Generic LRU Cache

## 1. What Is This About?

### Real-World Scenario

Imagine you're building a web application that fetches user profiles from a database. Each database query takes 50ms. With 10,000 requests per second, you'd be making 10,000 × 50ms = 500 seconds of DB work per second - clearly impossible!

**Solution**: Cache recently-used profiles in memory.

But you can't cache everything forever (memory is limited), so you need an **eviction policy**:

**❌ Bad policies:**
- **Random eviction**: Might evict a frequently-used item
- **FIFO (First-In-First-Out)**: Evicts oldest item, even if it's still popular

**✅ Best policy (usually):** **LRU (Least Recently Used)**
- Keep items that were used recently
- Evict items that haven't been used in a while
- Intuition: "If you used it recently, you'll probably use it again soon"

This project teaches you how to build a production-grade LRU cache with:
- **O(1) performance** for all operations (get, set, evict)
- **Generic types** (works with any key/value types)
- **Thread safety** (safe for concurrent access)
- **TTL support** (automatic expiration after a time duration)

### What You'll Learn

1. **LRU algorithm**: Map + Doubly-linked list combination
2. **Generics in Go**: Type parameters and constraints
3. **Thread safety**: Mutexes and critical sections
4. **container/list**: Go's built-in doubly-linked list
5. **Time handling**: TTL (time-to-live) expiration
6. **Zero values**: Generic zero-value handling

### The Challenge

Build a cache that:
- Stores up to N items in memory
- Retrieves items in O(1) time
- Evicts the least recently used item when capacity is exceeded
- Supports optional TTL (items expire after a duration)
- Is thread-safe (multiple goroutines can use it simultaneously)
- Works with any comparable key type and any value type

---

## 2. First Principles: Understanding LRU Caches

### What is a Cache?

A **cache** is a fast, small storage layer in front of a slow, large storage layer.

**Analogy**: Your desk vs your filing cabinet
- **Desk (cache)**: Small, only holds current documents, instant access
- **Filing cabinet (database)**: Large, holds everything, slower access
- **Strategy**: Keep frequently-used documents on your desk

### What is LRU?

**LRU (Least Recently Used)** means "throw away what you haven't used in a while."

**Example**: You have a bookshelf (capacity = 5 books)

```
Initial state: [A, B, C, D, E]

Read book C:
- C was used recently, move it to the "most recent" position
- [C, A, B, D, E]  (C is now at the front)

Add new book F:
- Shelf is full (5 books)
- Evict the least recently used book (E)
- Add F at the front
- [F, C, A, B, D]

Read book A:
- Move A to front
- [A, F, C, B, D]
```

**Key insight**: Items move to the front when accessed. Items at the back get evicted first.

### Why Map + Doubly-Linked List?

To achieve O(1) operations, we need TWO data structures working together:

#### Data Structure 1: Map (Hash Table)

```go
items map[Key]*list.Element
```

**Purpose**: O(1) lookup by key

**Limitation**: Maps don't maintain order (can't tell which item is "least recently used")

#### Data Structure 2: Doubly-Linked List

```
  ┌────┐    ┌────┐    ┌────┐    ┌────┐
  │ C  │◄──►│ A  │◄──►│ B  │◄──►│ D  │
  └────┘    └────┘    └────┘    └────┘
  Front                          Back
  (Most recent)                  (Least recent)
```

**Purpose**: Maintain recency order

**Operations**:
- Move to front: O(1) (when you know the node pointer)
- Remove from back: O(1) (always remove tail)
- Remove arbitrary node: O(1) (when you have the pointer)

#### Combining Both

```go
type Cache struct {
    items     map[K]*list.Element  // Key → Node in list
    evictList *list.List           // Ordered by recency
}
```

**Get operation**:
1. Look up key in map → O(1)
2. Get the list node from map value → O(1)
3. Move node to front of list → O(1)
4. **Total: O(1)**

**Set operation**:
1. Check if key exists in map → O(1)
2. If exists: Move to front → O(1)
3. If new and over capacity: Remove back node → O(1), delete from map → O(1)
4. Add to front of list → O(1), add to map → O(1)
5. **Total: O(1)**

### What are Generics?

**Before Go 1.18**, you'd write type-specific caches:

```go
// Cache for string keys, int values
type StringIntCache struct {
    items map[string]int
}

// Cache for int keys, string values
type IntStringCache struct {
    items map[int]string
}

// Cache for any key/value (loses type safety)
type AnyCache struct {
    items map[interface{}]interface{}
}
```

**Problem**: Code duplication or loss of type safety.

**With Go 1.18+ generics**, you write ONE cache that works for all types:

```go
type Cache[K comparable, V any] struct {
    items map[K]*list.Element
}

// Use it:
intCache := New[string, int](100, time.Minute)
userCache := New[int, User](100, time.Minute)
```

**Type parameters**:
- `K comparable`: K can be any type that supports `==` (ints, strings, pointers, etc.)
- `V any`: V can be absolutely any type

### What is TTL?

**TTL (Time To Live)** means items expire after a certain duration.

**Example**: Cache user sessions for 30 minutes

```go
cache.SetWithTTL("session-123", userData, 30*time.Minute)

// 29 minutes later:
data, ok := cache.Get("session-123")  // ok = true, data valid

// 31 minutes later:
data, ok := cache.Get("session-123")  // ok = false, expired
```

**Implementation**: Store expiration timestamp with each entry:

```go
type entry struct {
    key       K
    value     V
    expiresAt time.Time  // When this item expires
}
```

On `Get`, check if `time.Now().After(expiresAt)` → if yes, remove and return not found.

---

## 3. Breaking Down the Solution

### Step 1: Design the Data Structures

**Cache struct**:
```go
type Cache[K comparable, V any] struct {
    mu         sync.Mutex           // Thread safety
    capacity   int                  // Max items
    defaultTTL time.Duration        // Default expiration
    items      map[K]*list.Element  // Key → List node
    evictList  *list.List           // Recency order
}
```

**Entry struct** (what we store in the list):
```go
type entry[K comparable, V any] struct {
    key       K
    value     V
    expiresAt time.Time
}
```

**Why store `key` in the entry?**
When evicting from the back of the list, we need to know which key to delete from the map.

### Step 2: Get Operation

**Algorithm**:
1. Lock the mutex (thread safety)
2. Look up key in map
3. If not found → return zero value and false
4. If found → check TTL expiration
5. If expired → remove and return not found
6. If valid → move to front of list (mark as recently used)
7. Return value and true

**Visual walkthrough**:

```
Initial state:
Map: {A → node1, B → node2, C → node3}
List: [A] ← → [B] ← → [C]

Get(B):
1. Map lookup: B → node2 ✓
2. Check expiration: not expired ✓
3. Move B to front:
   List: [B] ← → [A] ← → [C]
4. Return value
```

### Step 3: Set Operation

**Algorithm**:
1. Lock the mutex
2. Calculate expiration time (if TTL > 0)
3. Check if key already exists:
   - If yes → update value, move to front
   - If no → create new entry, add to front
4. If over capacity → remove back node (LRU eviction)

**Visual walkthrough**:

```
Initial state (capacity = 3):
Map: {A → node1, B → node2, C → node3}
List: [B] ← → [A] ← → [C]

Set(D, value):
1. D doesn't exist in map
2. Create entry for D
3. Add to front:
   List: [D] ← → [B] ← → [A] ← → [C]
   (now has 4 items, over capacity!)
4. Evict back node (C):
   - Remove C from list
   - Delete C from map
   List: [D] ← → [B] ← → [A]
   Map: {A → node1, B → node2, D → node4}
```

### Step 4: Thread Safety

**Problem**: Multiple goroutines might call Get/Set simultaneously.

**Without locks**:
```go
// Goroutine 1:
cache.Set("key", 1)  // Writing to map

// Goroutine 2 (same time):
cache.Set("key", 2)  // Writing to map

// CRASH: concurrent map writes panic in Go
```

**Solution**: Protect all shared data with a mutex:
```go
func (c *Cache[K, V]) Set(key K, val V) {
    c.mu.Lock()         // Acquire lock
    defer c.mu.Unlock()  // Release lock (even if panic)

    // Safe: only one goroutine can be here at a time
    c.items[key] = val
}
```

**Mutex guarantees**: Only one goroutine can hold the lock at a time. Others wait.

---

## 4. Complete Solution Walkthrough

### Type Parameters

```go
type Cache[K comparable, V any] struct { ... }
```

**What does this mean?**

`[K comparable, V any]` declares two **type parameters**:
- `K` must be **comparable** (supports `==` operator) - required for map keys
- `V` can be **any** type (no constraints)

**Comparable types**:
- ✅ Integers, floats, strings, booleans
- ✅ Pointers, interfaces
- ✅ Structs (if all fields are comparable)
- ❌ Slices, maps, functions

**Why `K comparable`?**
Because map keys must be comparable: `map[K]` requires `K` to support `==`.

### Constructor

```go
func New[K comparable, V any](capacity int, defaultTTL time.Duration) *Cache[K, V] {
    return &Cache[K, V]{
        capacity:   capacity,
        defaultTTL: defaultTTL,
        items:      make(map[K]*list.Element),
        evictList:  list.New(),
    }
}
```

**Why explicit type parameters `[K comparable, V any]`?**
So callers can specify types:
```go
cache := New[string, int](100, time.Minute)
```

Go can't infer `K` and `V` from the parameters alone (capacity and TTL don't reveal the types).

### Get Method

```go
func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    var zero V  // Zero value for type V

    elem, ok := c.items[key]
    if !ok {
        return zero, false  // Not found
    }

    ent := elem.Value.(*entry[K, V])  // Type assertion

    // Check TTL
    if !ent.expiresAt.IsZero() && time.Now().After(ent.expiresAt) {
        c.removeElement(elem)
        return zero, false  // Expired
    }

    c.evictList.MoveToFront(elem)  // Mark as recently used
    return ent.value, true
}
```

**Line-by-line breakdown**:

1. **`c.mu.Lock()`**: Acquire lock (thread safety)
2. **`defer c.mu.Unlock()`**: Ensure lock is released (even if panic or early return)
3. **`var zero V`**: Declare a zero value of type V
   - For `int`: 0
   - For `string`: ""
   - For `*User`: nil
   - For `struct{Name string}`: `struct{Name: ""}`
4. **`elem, ok := c.items[key]`**: Look up in map
   - `ok == false` means key not found
5. **`ent := elem.Value.(*entry[K, V])`**: Type assertion
   - `list.Element.Value` has type `interface{}`
   - We assert it's actually `*entry[K, V]`
   - **Downside of container/list**: Loses type safety
6. **`ent.expiresAt.IsZero()`**: Check if expiration is set
   - `time.Time{}` (zero value) → `IsZero() == true` → no expiration
7. **`time.Now().After(ent.expiresAt)`**: Check if expired
8. **`c.evictList.MoveToFront(elem)`**: Move to front of list (O(1) operation)

**Why return `(V, bool)` instead of `(*V, error)`?**

Go convention: `(value, ok)` pattern for "value or not found" cases. Errors are for exceptional situations, not expected cases like cache misses.

### Set Method

```go
func (c *Cache[K, V]) SetWithTTL(key K, val V, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    var expiresAt time.Time
    if ttl > 0 {
        expiresAt = time.Now().Add(ttl)
    }

    // Update existing entry
    if elem, ok := c.items[key]; ok {
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
```

**Line-by-line breakdown**:

1. **Calculate expiration**: `time.Now().Add(ttl)` = current time + duration
   - `ttl <= 0` means no expiration (zero time)
2. **Update path** (key exists):
   - Move to front (mark as recently used)
   - Update value and expiration
   - Return early (don't add new entry)
3. **Insert path** (new key):
   - Create `entry` struct with key, value, expiration
   - `PushFront`: Add to front of list (most recent position)
   - Add mapping: `items[key] = elem`
4. **Eviction**: If over capacity, remove back (least recently used)
   - `Back()`: Get last element
   - `removeElement`: Remove from both list and map

**Why check `c.evictList.Len() > c.capacity` instead of `>=`?**

Because we add the new item first, then check. After adding, if len == capacity + 1, we evict one to bring it back to capacity.

### removeElement Helper

```go
func (c *Cache[K, V]) removeElement(elem *list.Element) {
    c.evictList.Remove(elem)                  // Remove from list
    ent := elem.Value.(*entry[K, V])          // Extract entry
    delete(c.items, ent.key)                  // Delete from map
}
```

**Why do we need `ent.key`?**

The list element doesn't know its key - only the entry struct inside does. We extract the entry to get the key, then delete it from the map.

**Critical**: Must keep list and map in sync. Every removal must update both.

---

## 5. Key Concepts Explained

### Concept 1: Generics in Go

**Type parameters** allow writing code that works with multiple types.

```go
// Without generics:
func PrintInt(val int) {
    fmt.Println(val)
}

func PrintString(val string) {
    fmt.Println(val)
}

// With generics:
func Print[T any](val T) {
    fmt.Println(val)
}
```

**Type constraints** limit what types are allowed:

```go
// T can be any type
func Identity[T any](val T) T {
    return val
}

// K must be comparable (can use ==, !=)
func Equal[K comparable](a, b K) bool {
    return a == b
}

// T must be a number (Go 1.21+)
func Sum[T int | float64](vals []T) T {
    var sum T
    for _, v := range vals {
        sum += v
    }
    return sum
}
```

**Common constraints**:
- `any`: No constraints (any type allowed)
- `comparable`: Types that support `==` and `!=`
- Custom: `interface{ Method() }` or unions `int | float64`

### Concept 2: container/list

Go's `container/list` is a **doubly-linked list** with these operations:

```go
l := list.New()

// Add elements
elem1 := l.PushFront(value)  // Add to front
elem2 := l.PushBack(value)   // Add to back

// Move elements
l.MoveToFront(elem1)         // Move to front
l.MoveToBack(elem1)          // Move to back

// Remove elements
l.Remove(elem1)              // Remove specific element

// Access elements
front := l.Front()           // First element
back := l.Back()             // Last element
next := elem1.Next()         // Next element
prev := elem1.Prev()         // Previous element

// Size
length := l.Len()
```

**Element structure**:
```go
type Element struct {
    Value interface{}  // The stored value (any type)
    // prev, next pointers (internal)
}
```

**Why doubly-linked?**

Allows O(1) removal of arbitrary nodes (if you have the pointer). Singly-linked would require O(n) to find predecessor.

### Concept 3: Zero Values in Generics

**Problem**: How do you return "no value" for a generic type?

```go
func Get[V any](key string) V {
    if notFound {
        return ???  // What to return?
    }
}
```

**Solution**: Use `var zero V` to get the zero value:

```go
func Get[V any](key string) (V, bool) {
    var zero V  // Zero value for type V
    if notFound {
        return zero, false
    }
    return actualValue, true
}
```

**Zero values by type**:
- `int`, `float64`: `0`
- `string`: `""`
- `bool`: `false`
- `*T` (pointer): `nil`
- `[]T` (slice): `nil`
- `map[K]V`: `nil`
- `struct{...}`: All fields are zero values

### Concept 4: sync.Mutex

**Mutex** (mutual exclusion) ensures only one goroutine accesses shared data at a time.

```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (sc *SafeCounter) Increment() {
    sc.mu.Lock()
    sc.count++
    sc.mu.Unlock()
}
```

**Better pattern**: Use `defer` to ensure unlock:

```go
func (sc *SafeCounter) Increment() {
    sc.mu.Lock()
    defer sc.mu.Unlock()  // Unlock when function returns
    sc.count++
    // Even if panic, unlock happens
}
```

**Deadlock warning**:

```go
// ❌ DEADLOCK
func (sc *SafeCounter) Bad() {
    sc.mu.Lock()
    sc.Increment()  // Tries to lock again → deadlock!
    sc.mu.Unlock()
}
```

**Solution**: Don't call locked methods from locked methods, or use helper methods without locks.

### Concept 5: LRU vs Other Eviction Policies

**LRU (Least Recently Used)**: Evict item used longest ago
- **Pros**: Good hit rate for temporal locality (recent items used again)
- **Cons**: Scan pattern (access all items once) evicts everything

**LFU (Least Frequently Used)**: Evict item accessed least often
- **Pros**: Keeps popular items even if not recent
- **Cons**: Hard to implement efficiently; old items stick around

**FIFO (First-In-First-Out)**: Evict oldest item
- **Pros**: Simple (just a queue)
- **Cons**: Ignores access patterns

**Random**: Evict random item
- **Pros**: Simplest possible; no pathological cases
- **Cons**: Can evict popular items

**ARC (Adaptive Replacement Cache)**: Balances recency and frequency
- **Pros**: Best hit rate in many workloads
- **Cons**: Complex, patented

**When to use LRU**:
- Web server caching (recent pages accessed again)
- Database query caching
- API response caching
- General-purpose caching

---

## 6. Common Patterns You Can Reuse

### Pattern 1: Generic Constructor

```go
func NewThing[T any](param int) *Thing[T] {
    return &Thing[T]{
        data: make([]T, 0, param),
    }
}

// Usage:
intThing := NewThing[int](10)
stringThing := NewThing[string](10)
```

### Pattern 2: Generic Method with Constraints

```go
type Set[T comparable] struct {
    items map[T]struct{}
}

func (s *Set[T]) Add(item T) {
    s.items[item] = struct{}{}
}

func (s *Set[T]) Contains(item T) bool {
    _, ok := s.items[item]
    return ok
}
```

### Pattern 3: Thread-Safe Collection

```go
type SafeMap[K comparable, V any] struct {
    mu   sync.RWMutex
    data map[K]V
}

func (sm *SafeMap[K, V]) Get(key K) (V, bool) {
    sm.mu.RLock()         // Read lock (allows concurrent reads)
    defer sm.mu.RUnlock()
    val, ok := sm.data[key]
    return val, ok
}

func (sm *SafeMap[K, V]) Set(key K, val V) {
    sm.mu.Lock()          // Write lock (exclusive)
    defer sm.mu.Unlock()
    sm.data[key] = val
}
```

### Pattern 4: TTL with Background Cleanup

```go
type CacheWithCleanup[K comparable, V any] struct {
    cache *Cache[K, V]
    stop  chan struct{}
}

func (c *CacheWithCleanup[K, V]) Start() {
    ticker := time.NewTicker(1 * time.Minute)
    go func() {
        for {
            select {
            case <-ticker.C:
                c.cache.RemoveExpired()  // Custom method
            case <-c.stop:
                ticker.Stop()
                return
            }
        }
    }()
}
```

### Pattern 5: Sharded Cache (Reduce Lock Contention)

```go
type ShardedCache[K comparable, V any] struct {
    shards []*Cache[K, V]
    count  int
}

func (sc *ShardedCache[K, V]) getShard(key K) *Cache[K, V] {
    hash := hashKey(key)
    return sc.shards[hash%sc.count]
}

func (sc *ShardedCache[K, V]) Get(key K) (V, bool) {
    return sc.getShard(key).Get(key)
}
```

Reduces contention: Different keys likely go to different shards with separate locks.

---

## 7. Real-World Applications

### Web Server Response Caching

**Use case**: Cache rendered HTML or API responses

```go
type ResponseCache struct {
    cache *Cache[string, []byte]
}

func (rc *ResponseCache) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cacheKey := r.Method + ":" + r.URL.Path

        if cached, ok := rc.cache.Get(cacheKey); ok {
            w.Write(cached)
            return
        }

        // Render response, cache it
        resp := render(r)
        rc.cache.Set(cacheKey, resp)
        w.Write(resp)
    })
}
```

Companies using this: Reddit, Twitter, GitHub (for static pages)

### Database Query Result Caching

**Use case**: Avoid repeated expensive queries

```go
type QueryCache struct {
    cache *Cache[string, *sql.Rows]
    db    *sql.DB
}

func (qc *QueryCache) Query(sql string, args ...interface{}) (*sql.Rows, error) {
    cacheKey := sql + fmt.Sprint(args)

    if rows, ok := qc.cache.Get(cacheKey); ok {
        return rows, nil
    }

    rows, err := qc.db.Query(sql, args...)
    if err != nil {
        return nil, err
    }

    qc.cache.SetWithTTL(cacheKey, rows, 5*time.Minute)
    return rows, nil
}
```

Companies using this: Any database-heavy application

### CDN / Image Caching

**Use case**: Cache recently-accessed images in memory

```go
type ImageCache struct {
    cache *Cache[string, []byte]
    disk  DiskStorage
}

func (ic *ImageCache) Get(imageID string) ([]byte, error) {
    // Check memory cache
    if img, ok := ic.cache.Get(imageID); ok {
        return img, nil
    }

    // Load from disk
    img, err := ic.disk.Read(imageID)
    if err != nil {
        return nil, err
    }

    // Cache in memory
    ic.cache.SetWithTTL(imageID, img, 1*time.Hour)
    return img, nil
}
```

Companies using this: Cloudflare, Fastly, AWS CloudFront

### Session Storage

**Use case**: Store user sessions with automatic expiration

```go
type SessionStore struct {
    cache *Cache[string, *Session]
}

func (ss *SessionStore) CreateSession(userID int) string {
    sessionID := generateID()
    session := &Session{UserID: userID, CreatedAt: time.Now()}

    ss.cache.SetWithTTL(sessionID, session, 30*time.Minute)
    return sessionID
}

func (ss *SessionStore) GetSession(sessionID string) (*Session, bool) {
    return ss.cache.Get(sessionID)
}
```

Companies using this: Any web application with user sessions

### DNS Caching

**Use case**: Cache DNS lookups to avoid repeated network calls

```go
type DNSCache struct {
    cache *Cache[string, net.IP]
}

func (dc *DNSCache) Resolve(hostname string) (net.IP, error) {
    if ip, ok := dc.cache.Get(hostname); ok {
        return ip, nil
    }

    ips, err := net.LookupIP(hostname)
    if err != nil {
        return nil, err
    }

    ip := ips[0]
    dc.cache.SetWithTTL(hostname, ip, 5*time.Minute)
    return ip, nil
}
```

Companies using this: Operating systems, browsers, DNS servers

---

## 8. Common Mistakes to Avoid

### Mistake 1: Not Using `defer` with Mutex

**❌ Wrong**:
```go
func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mu.Lock()
    val, ok := c.items[key]
    c.mu.Unlock()  // What if we add early return later?
    return val, ok
}
```

**✅ Correct**:
```go
func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()  // Always unlocks
    val, ok := c.items[key]
    return val, ok
}
```

**Why**: `defer` ensures unlock happens even with early returns or panics.

### Mistake 2: Forgetting to Move to Front on Get

**❌ Wrong**:
```go
func (c *Cache[K, V]) Get(key K) (V, bool) {
    elem, ok := c.items[key]
    if !ok {
        return zero, false
    }
    // Forgot to move to front!
    return elem.Value.(*entry[K, V]).value, true
}
```

**Result**: Items aren't marked as recently used → LRU logic breaks.

**✅ Correct**:
```go
func (c *Cache[K, V]) Get(key K) (V, bool) {
    elem, ok := c.items[key]
    if !ok {
        return zero, false
    }
    c.evictList.MoveToFront(elem)  // Mark as recently used
    return elem.Value.(*entry[K, V]).value, true
}
```

### Mistake 3: Map and List Out of Sync

**❌ Wrong**:
```go
func (c *Cache[K, V]) removeElement(elem *list.Element) {
    c.evictList.Remove(elem)
    // Forgot to delete from map!
}
```

**Result**: Map grows unbounded; memory leak.

**✅ Correct**:
```go
func (c *Cache[K, V]) removeElement(elem *list.Element) {
    c.evictList.Remove(elem)
    ent := elem.Value.(*entry[K, V])
    delete(c.items, ent.key)  // Keep map in sync
}
```

### Mistake 4: Using Pointer Types Incorrectly

**❌ Wrong**:
```go
func New[K comparable, V any](cap int) *Cache[K, V] {
    return &Cache[K, V]{
        items: make(map[K]V),  // Should be map[K]*list.Element
    }
}
```

**Result**: Compilation error - wrong map value type.

### Mistake 5: Not Checking Expiration on Get

**❌ Wrong**:
```go
func (c *Cache[K, V]) Get(key K) (V, bool) {
    elem, ok := c.items[key]
    if !ok {
        return zero, false
    }
    // Forgot to check if expired!
    return elem.Value.(*entry[K, V]).value, true
}
```

**Result**: Returns expired items.

**✅ Correct**: Check expiration and remove if expired (as shown in solution).

### Mistake 6: Evicting Before Adding (Off-by-One)

**❌ Wrong**:
```go
func (c *Cache[K, V]) Set(key K, val V) {
    // Evict first
    if c.evictList.Len() >= c.capacity {
        c.removeElement(c.evictList.Back())
    }
    // Then add
    c.evictList.PushFront(entry)
    c.items[key] = entry
}
```

**Result**: Cache never reaches capacity (always one less).

**✅ Correct**: Add first, then check if over capacity.

### Mistake 7: Using RWMutex for LRU

**❌ Wrong**:
```go
type Cache[K, V] struct {
    mu sync.RWMutex  // Read-write mutex
    // ...
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mu.RLock()  // Read lock
    defer c.mu.RUnlock()
    // ...
    c.evictList.MoveToFront(elem)  // WRONG: Modifies list with read lock!
    // ...
}
```

**Problem**: `MoveToFront` modifies the list, but we only have a read lock.

**Solution**: LRU cache needs write lock even for Get (because it moves elements). Use `sync.Mutex`, not `sync.RWMutex`.

---

## 9. Stretch Goals

### Goal 1: Add Statistics Tracking ⭐

Track hit rate, eviction count, and other metrics.

**Hint**:
```go
type Cache[K, V] struct {
    // ... existing fields ...
    hits      int64
    misses    int64
    evictions int64
}

func (c *Cache[K, V]) Stats() CacheStats {
    c.mu.Lock()
    defer c.mu.Unlock()
    return CacheStats{
        Hits:      c.hits,
        Misses:    c.misses,
        Evictions: c.evictions,
        HitRate:   float64(c.hits) / float64(c.hits + c.misses),
    }
}
```

### Goal 2: Implement LFU (Least Frequently Used) ⭐⭐

Instead of evicting least recently used, evict least frequently used.

**Hint**: Add frequency counter to entry:
```go
type entry[K, V] struct {
    key       K
    value     V
    frequency int
    expiresAt time.Time
}
```

Use a min-heap to track items by frequency.

### Goal 3: Add Persistence to Disk ⭐⭐⭐

Serialize cache to disk, restore on startup.

**Hint**:
```go
func (c *Cache[K, V]) Save(filename string) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    entries := make([]entry[K, V], 0, c.evictList.Len())
    for elem := c.evictList.Front(); elem != nil; elem = elem.Next() {
        entries = append(entries, *elem.Value.(*entry[K, V]))
    }

    data, err := json.Marshal(entries)
    if err != nil {
        return err
    }

    return os.WriteFile(filename, data, 0644)
}
```

### Goal 4: Sharded Cache (Reduce Lock Contention) ⭐⭐⭐

Split cache into N shards, each with its own lock. Hash keys to shards.

**Hint**:
```go
type ShardedCache[K comparable, V any] struct {
    shards    []*Cache[K, V]
    shardMask uint32
}

func (sc *ShardedCache[K, V]) getShard(key K) *Cache[K, V] {
    h := hash(key)
    return sc.shards[h&sc.shardMask]
}
```

This reduces contention: different keys go to different shards.

### Goal 5: Add Prometheus Metrics ⭐⭐⭐

Export cache metrics to Prometheus for monitoring.

**Hint**:
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    cacheHits = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "cache_hits_total",
        Help: "Total cache hits",
    })
    cacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "cache_misses_total",
        Help: "Total cache misses",
    })
)

func init() {
    prometheus.MustRegister(cacheHits, cacheMisses)
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
    val, ok := c.get(key)
    if ok {
        cacheHits.Inc()
    } else {
        cacheMisses.Inc()
    }
    return val, ok
}
```

---

## How to Run

```bash
# Run the program
make run P=07-generic-lru-cache

# Run tests
go test ./minis/07-generic-lru-cache/...

# Run benchmarks
go test -bench=. -benchmem ./minis/07-generic-lru-cache/...

# Run with race detector
go test -race ./minis/07-generic-lru-cache/...
```

---

## Summary

**What you learned**:
- ✅ LRU cache algorithm (map + doubly-linked list)
- ✅ Generics in Go (type parameters, constraints)
- ✅ Thread safety with sync.Mutex
- ✅ container/list for linked list operations
- ✅ TTL (time-to-live) expiration
- ✅ Zero values in generic code

**Why this matters**:
Caching is critical for performance in real-world systems. LRU is one of the most common caching strategies, and this project gives you a production-ready implementation you can use anywhere.

**Key insights**:
- Generics enable type-safe code reuse
- Combining data structures (map + list) achieves O(1) operations
- Thread safety requires careful lock management
- TTL is essential for keeping caches fresh

**Next steps**:
- Project 08: HTTP clients with retries and exponential backoff
- Project 09: HTTP servers with middleware and graceful shutdown
- Project 10: gRPC services with streaming

Cache on! 🚀
