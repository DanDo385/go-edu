# How the Generic LRU Cache Solution Works

This document explains in plain language how the generic LRU cache solution is implemented. The solution uses Go's generics, built-in data structures, and synchronization primitives to create a thread-safe, efficient cache with automatic eviction and TTL support.

## The Big Picture

The solution implements a **Least Recently Used (LRU) cache** that automatically evicts the least recently accessed item when the cache reaches capacity. It combines two data structures to achieve O(1) operations:

- **Map (hash table)**: For O(1) key lookup
- **Doubly-linked list**: For O(1) insertion, removal, and maintaining recency order

The cache is **generic** (works with any key/value types), **thread-safe** (multiple goroutines can use it simultaneously), and supports **TTL expiration** (items automatically expire after a duration).

## Architecture Overview

The solution consists of four main components:

1. **Cache struct**: The main cache container with map and list
2. **entry struct**: Holds the actual cached data (key, value, expiration time)
3. **Get method**: Retrieves values and updates recency
4. **Set/SetWithTTL methods**: Inserts or updates items and handles eviction

## Component 1: Cache Struct - The Foundation

The `Cache` struct is the heart of the solution. It uses Go generics to work with any key and value types.

```go
type Cache[K comparable, V any] struct {
    mu          sync.Mutex
    capacity    int
    defaultTTL  time.Duration
    items       map[K]*list.Element
    evictList   *list.List
}
```

**What each field does:**

- `mu`: A mutex (mutual exclusion lock) that protects all cache operations from concurrent access. Only one goroutine can modify the cache at a time.
- `capacity`: Maximum number of items the cache can hold. When exceeded, LRU eviction occurs.
- `defaultTTL`: Default time-to-live for items. Zero means no expiration.
- `items`: A map that provides O(1) lookup by key. Maps key → list element pointer.
- `evictList`: A doubly-linked list that maintains recency order. Front = most recent, back = least recent.

**Why two data structures?**

- **Map alone**: Can't tell which item is "least recently used" (maps are unordered)
- **List alone**: Can't find items quickly (would need O(n) search)
- **Map + List together**: Get O(1) lookup (map) + O(1) eviction (list back)

**Generic type parameters:**

- `K comparable`: Key type must be comparable (can use == and !=). Required for map keys.
- `V any`: Value type can be anything (no constraints).

## Component 2: entry Struct - Storing Data

The `entry` struct holds the actual cached data within list elements:

```go
type entry[K comparable, V any] struct {
    key       K
    value     V
    expiresAt time.Time
}
```

**Why store the key?**

When evicting from the back of the list, we need the key to delete from the map. List elements only contain the entry, not the key directly, so we store it in the entry.

**expiresAt field:**

- Zero time (`time.Time{}`) means no expiration
- Non-zero time means the item expires at that moment
- Checked lazily on `Get` (not proactively)

## Component 3: New Function - Cache Creation

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

**What happens here:**

- Creates a new `Cache` instance on the heap (returns pointer)
- Initializes the map with `make()` (empty map ready to use)
- Initializes the list with `list.New()` (empty doubly-linked list)
- Mutex zero value is unlocked (ready to use)
- Returns pointer so caller can use the cache

**Memory allocation:**

- Cache struct: ~80 bytes (fields + mutex)
- Empty map: ~48 bytes (hash table header)
- Empty list: ~48 bytes (list header)
- Total: ~176 bytes for empty cache

## Component 4: Get Method - Retrieving Values

The `Get` method retrieves a value by key and updates recency order.

### Step 1: Lock and Lookup

```go
c.mu.Lock()
defer c.mu.Unlock()

var zero V
elem, ok := c.items[key]
if !ok {
    return zero, false
}
```

**What happens:**

- `c.mu.Lock()`: Acquires the mutex (blocks if another goroutine is using cache)
- `defer c.mu.Unlock()`: Schedules unlock when function exits (even on early return)
- `var zero V`: Creates zero value for type V (e.g., 0 for int, "" for string)
- Look up key in map: `c.items[key]` returns `(*list.Element, bool)`
- If not found (`!ok`), return zero value and false

**Why defer unlock?**

Ensures the mutex is always released, even if the function returns early or panics. This prevents deadlocks.

### Step 2: Extract Entry and Check Expiration

```go
ent := elem.Value.(*entry[K, V])

if !ent.expiresAt.IsZero() && time.Now().After(ent.expiresAt) {
    c.removeElement(elem)
    return zero, false
}
```

**What happens:**

- **Type assertion**: `elem.Value` is `interface{}` (Go's `container/list` isn't generic), so we assert it's `*entry[K, V]`
- **Check expiration**: 
  - `!ent.expiresAt.IsZero()`: Entry has an expiration time
  - `time.Now().After(ent.expiresAt)`: Current time is after expiration time
- If expired, remove entry and return not found

**Lazy expiration:**

We only check expiration on `Get`, not proactively. This saves resources (no background goroutine needed) but expired items may linger until accessed.

### Step 3: Move to Front and Return

```go
c.evictList.MoveToFront(elem)
return ent.value, true
```

**What happens:**

- `MoveToFront(elem)`: Moves the list element to the front (O(1) operation)
- This marks the item as "most recently used"
- Returns the value and true (found)

**Why move to front?**

LRU policy: recently accessed items should be kept. Moving to front makes this item the last to be evicted.

## Component 5: SetWithTTL Method - Inserting/Updating Values

The `SetWithTTL` method inserts or updates a key-value pair with a custom TTL.

### Step 1: Lock and Calculate Expiration

```go
c.mu.Lock()
defer c.mu.Unlock()

var expiresAt time.Time
if ttl > 0 {
    expiresAt = time.Now().Add(ttl)
}
```

**What happens:**

- Lock the mutex (same as Get)
- Calculate expiration time:
  - If `ttl > 0`: Set expiration to current time + TTL
  - If `ttl == 0`: Leave as zero time (no expiration)

### Step 2: Check if Key Exists (Update Case)

```go
if elem, ok := c.items[key]; ok {
    c.evictList.MoveToFront(elem)
    ent := elem.Value.(*entry[K, V])
    ent.value = val
    ent.expiresAt = expiresAt
    return
}
```

**What happens:**

- Check if key already exists in map
- If exists (update case):
  - Move element to front (mark as recently used)
  - Update entry's value and expiration time
  - Return early (no eviction needed)

**Why update in place?**

More efficient than removing and re-adding. We just update the existing entry.

### Step 3: Insert New Entry

```go
ent := &entry[K, V]{
    key:       key,
    value:     val,
    expiresAt: expiresAt,
}
elem := c.evictList.PushFront(ent)
c.items[key] = elem
```

**What happens:**

- Create new entry struct with key, value, and expiration
- `PushFront(ent)`: Add entry to front of list (O(1))
- `c.items[key] = elem`: Store list element pointer in map (O(1))
- Now both map and list reference the same entry

**Why push to front?**

New items are "most recently used" by definition, so they go to the front.

### Step 4: Evict if Over Capacity

```go
if c.evictList.Len() > c.capacity {
    c.removeElement(c.evictList.Back())
}
```

**What happens:**

- Check if cache size exceeds capacity
- `evictList.Back()`: Get least recently used element (back of list)
- `removeElement()`: Remove from both list and map

**Eviction strategy:**

- Evict happens AFTER insertion (cache temporarily exceeds capacity by 1)
- Always evicts the back element (least recently used)
- Maintains capacity invariant: `Len() <= capacity` after eviction

## Component 6: removeElement Helper - Maintaining Consistency

The `removeElement` helper ensures both map and list stay in sync:

```go
func (c *Cache[K, V]) removeElement(elem *list.Element) {
    c.evictList.Remove(elem)
    ent := elem.Value.(*entry[K, V])
    delete(c.items, ent.key)
}
```

**What happens:**

1. Remove from list: `evictList.Remove(elem)` (O(1))
2. Extract entry: Get the entry from the list element
3. Delete from map: `delete(c.items, ent.key)` (O(1))

**Why this helper?**

Both map and list must be updated together. This helper ensures consistency and avoids code duplication.

**Must be called while holding lock:**

This function is not thread-safe on its own. It must be called from methods that already hold the mutex.

## Component 7: Len Method - Getting Cache Size

```go
func (c *Cache[K, V]) Len() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.evictList.Len()
}
```

**What happens:**

- Lock mutex (thread-safe access)
- Return list length (same as map length - invariant we maintain)

**Why check list length?**

Both map and list should have the same length (invariant). We use list length because it's the source of truth for eviction order.

## Thread Safety

The cache is thread-safe because:

1. **Mutex protection**: All operations acquire the mutex before accessing shared state
2. **Consistent locking**: Every method that reads/writes cache state locks the mutex
3. **Defer unlock**: Ensures mutex is always released, even on early returns or panics

**Example concurrent access:**

```go
// Goroutine 1
go func() {
    cache.Set("key1", 1)
}()

// Goroutine 2
go func() {
    val, ok := cache.Get("key1")
}()
```

Both goroutines can safely access the cache simultaneously. The mutex ensures only one modifies at a time.

## Time Complexity

All operations are O(1) average case:

- **Get**: O(1) - map lookup + list move to front
- **Set**: O(1) - map insert + list push front + evict (if needed)
- **Len**: O(1) - list length check
- **Eviction**: O(1) - remove from back of list + map delete

**Why O(1)?**

- Map operations: O(1) average (hash table)
- List operations: O(1) (doubly-linked list with pointers)
- No loops or scans needed

## Space Complexity

- **Per item**: ~96 bytes overhead (map entry + list node) + size(key) + size(value)
- **Total**: O(capacity) - bounded by capacity limit
- **Empty cache**: ~176 bytes (struct + empty map + empty list)

## Example Execution Flow

For a cache with capacity 3:

1. **Set("a", 1)**: Add "a" → [a] (front)
2. **Set("b", 2)**: Add "b" → [b, a] (b is most recent)
3. **Set("c", 3)**: Add "c" → [c, b, a] (c is most recent)
4. **Get("a")**: Access "a" → [a, c, b] (a moved to front)
5. **Set("d", 4)**: Add "d" → [d, a, c] (b evicted, d is most recent)

**Key insight:** Items move to front when accessed. Items at the back get evicted first.

## TTL Expiration Flow

For items with TTL=2 seconds:

1. **SetWithTTL("key", "value", 2*time.Second)**: Item expires at `now + 2s`
2. **Get("key") within 2 seconds**: Returns value, moves to front
3. **Get("key") after 2 seconds**: Returns (zero, false), item removed

**Lazy expiration:**

- Expired items are only removed when accessed
- No background cleanup goroutine needed
- Trade-off: Expired items consume memory until accessed

## Generics in Action

The cache works with any types:

```go
// String keys, int values
cache1 := New[string, int](100, time.Minute)

// Int keys, string values
cache2 := New[int, string](50, 0)

// Custom struct keys
type UserID string
type User struct { Name string }
cache3 := New[UserID, User](1000, time.Hour)
```

**Type safety:**

- Compiler ensures type correctness
- No runtime type assertions needed (except for `container/list`)
- Zero overhead compared to `interface{}` approach

## Common Patterns

### Pattern 1: Cache-Aside

```go
func GetUser(id string) (*User, error) {
    // Try cache first
    if user, ok := cache.Get(id); ok {
        return user, nil
    }
    
    // Cache miss: load from database
    user, err := db.LoadUser(id)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    cache.Set(id, user)
    return user, nil
}
```

### Pattern 2: Write-Through

```go
func UpdateUser(id string, user *User) error {
    // Update database
    if err := db.SaveUser(id, user); err != nil {
        return err
    }
    
    // Update cache
    cache.Set(id, user)
    return nil
}
```

### Pattern 3: TTL-Based Refresh

```go
// Cache with 5-minute TTL
cache := New[string, Data](100, 5*time.Minute)

// Items automatically expire after 5 minutes
// Next access will trigger refresh from source
```

## Summary

The LRU cache solution demonstrates:

- **Generics**: Type-safe without code duplication
- **Data structures**: Map + doubly-linked list for O(1) operations
- **Thread safety**: Mutex protects concurrent access
- **TTL support**: Automatic expiration with lazy cleanup
- **Memory efficiency**: Bounded by capacity, automatic eviction

**Key patterns:**
- Two data structures working together (map + list)
- Mutex for thread safety
- Lazy expiration (check on access)
- Generic type parameters for flexibility

This pattern is widely used in production systems for caching database queries, API responses, computed results, and any expensive-to-compute data that benefits from temporal locality.
