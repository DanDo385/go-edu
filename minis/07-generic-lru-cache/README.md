# Project 07: generic-lru-cache

## What You're Building

A thread-safe Least Recently Used (LRU) cache with generics, optional TTL (time-to-live), and proper eviction policies. This project demonstrates advanced Go features and data structure design.

## Concepts Covered

- Generics (type parameters with constraints)
- `container/list` for doubly-linked list operations
- `sync.Mutex` for thread safety
- `time.Time` for TTL expiration
- Map + linked list combination (O(1) operations)
- Benchmarking with `testing.B`

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

## Solution Explanation

### LRU Algorithm

**Data Structure**: Map + Doubly-Linked List
- **Map**: Key → List Node (O(1) lookup)
- **List**: Ordered by recency (front = most recent, back = least recent)

**Operations**:
- **Get**: Move accessed node to front (mark as recently used)
- **Set**: Add to front; if capacity exceeded, evict from back
- **TTL**: Check expiration on Get; lazy eviction

### Why Generics?

Before Go 1.18, you'd use `interface{}` and type assertions:
```go
cache.Set("key", 42)
val := cache.Get("key").(int) // Type assertion required
```

With generics:
```go
cache := New[string, int](100, time.Minute)
cache.Set("key", 42)
val, ok := cache.Get("key") // Type-safe!
```

## Where Go Shines

**Go vs Python:**
- Python: `functools.lru_cache` decorator is convenient
- Go: Explicit cache management; no hidden side effects

**Go vs Java:**
- Java: `LinkedHashMap` with custom ordering
- Go: Simpler with generics (no type erasure)

**Go vs Rust:**
- Rust: Similar performance, but lifetime complexity
- Go: Simpler with GC; mutexes are straightforward

## Stretch Goals

1. **Add statistics**: Track hit rate, eviction count
2. **Implement LFU**: Least Frequently Used eviction policy
3. **Add persistence**: Serialize cache to disk
4. **Sharded cache**: Multiple caches to reduce lock contention
5. **Tiered eviction**: Hot vs cold items
