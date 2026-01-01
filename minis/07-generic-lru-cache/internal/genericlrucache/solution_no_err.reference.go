//go:build reference

package genericlrucache

/*
Core LRU Cache Algorithm - Simplified Version
==============================================

This file focuses on the core LRU algorithm without error handling or edge cases.
Use this to understand the fundamental data structure operations.

ALGORITHM STEPS:

1. Cache Lookup (Get):
   - Lookup key in map → O(1)
   - Move accessed element to front → O(1)
   - Return value

2. Cache Insert (Set):
   - Create entry with key/value
   - Push to front of list → O(1)
   - Store in map → O(1)
   - Evict back element if over capacity → O(1)

3. LRU Eviction:
   - Remove element from back of list → O(1)
   - Remove from map using key → O(1)

DEBUGGING GUIDE:
- Set breakpoint at map lookup to trace cache hits/misses
- Set breakpoint at eviction to see which items are removed
- Watch the doubly-linked list structure to see recency order
*/

// TODO: Implement SimpleCacheGet
// Core algorithm for cache lookup without TTL or error handling
// Steps:
// 1. Look up key in map
// 2. If found, move element to front
// 3. Return value
// Signature: func SimpleCacheGet[K comparable, V any](cache *Cache[K, V], key K) V

// TODO: Implement SimpleCacheSet
// Core algorithm for cache insertion without TTL
// Steps:
// 1. Check if key exists
// 2. If exists: move to front and update value
// 3. If not exists:
//    - Create new entry
//    - Push to front
//    - Add to map
//    - If over capacity, remove back element
// Signature: func SimpleCacheSet[K comparable, V any](cache *Cache[K, V], key K, val V)

// TODO: Implement SimpleLRUEvict
// Core algorithm for LRU eviction
// Steps:
// 1. Get back element from list (least recently used)
// 2. Remove from list
// 3. Extract key from entry
// 4. Delete from map
// Signature: func SimpleLRUEvict[K comparable, V any](cache *Cache[K, V])
