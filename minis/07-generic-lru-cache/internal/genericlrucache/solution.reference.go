//go:build reference

package genericlrucache

/*
Reference Solution - Generic LRU Cache with Eviction
===================================================

This file implements a Least Recently Used (LRU) cache: a fixed-size cache that
evicts the least recently accessed item when full. Combines Go generics (Go 1.18+)
with the classic "map + doubly-linked list" LRU design for O(1) get/set operations.

This connects to:
- container/list: Go's doubly-linked list for O(1) reordering (MoveToFront, Remove)
- sync.Mutex: protects concurrent access for thread safety
- Generics: Cache[K comparable, V any] - K must support == (map keys), V can be anything

The LRU invariant: evictList orders items by access recency. Front = most recent,
Back = least recent. On capacity overflow, we evict from the back.

Teaching notes:
- Memory/ownership: items map and evictList share *list.Element pointers. The
  entry struct is stored in Element.Value. We must keep map and list in sync.
- Invariants: len(items) <= capacity; evictList order matches access order.
- Type assertion: elem.Value.(*entry[K,V]) - list stores interface{}, we assert
  back to our concrete type.
*/

import (
	"container/list"
	"sync"
)

// Cache holds key-value pairs with LRU eviction. K must be comparable (for map keys).
type Cache[K comparable, V any] struct {
	mu        sync.Mutex
	capacity  int
	items     map[K]*list.Element // key -> list element for O(1) lookup
	evictList *list.List          // front = most recent, back = least recent
}

// entry stores key+value in each list element (we need key for eviction lookup)
type entry[K comparable, V any] struct {
	key   K
	value V
}

/*
New - Construct LRU Cache

Creates a cache with the given capacity. Capacity is clamped to at least 1.
*/
func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity <= 0 {
		capacity = 1
	}
	return &Cache[K, V]{
		capacity:  capacity,
		items:     make(map[K]*list.Element, capacity),
		evictList: list.New(),
	}
}

/*
Get - Lookup Value by Key, Promote to Most Recent

Deep explanation of elem.Value.(*entry[K,V]) (per .cursorrules):

elem is *list.Element. elem.Value is interface{} — it can hold any type.
We stored a *entry[K,V] (pointer to our entry struct) when we did PushFront.

.(*entry[K,V]) is a type assertion:
  - "Treat the interface{} as a *entry[K,V]." We assert we know the concrete type.
  - Result is *entry[K,V] — a pointer. The .value then reads the V from that struct.
  - If we were wrong about the type, this would panic. We control what we store.

Memory: elem.Value holds a pointer to an entry. The assertion extracts that pointer;
we don't copy the entry — we follow the pointer and read .value.
*/
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.evictList.MoveToFront(elem)
		// elem.Value is interface{}; we stored *entry. Assert, then read .value.
		return elem.Value.(*entry[K, V]).value, true
	}
	var zeroV V
	return zeroV, false
}

/*
Set - Insert or Update Key-Value, Evict LRU if Full

Deep explanation of &entry[K,V]{...} and .value = value (per .cursorrules):

&entry[K,V]{key: key, value: value}:
  - Creates a new entry struct. & takes its address → *entry[K,V].
  - We store a POINTER in the list. The list Element holds our pointer in its Value (interface{}).
  - This does NOT create a copy of the entry for the list — we store the address.

elem.Value.(*entry[K,V]).value = value:
  - Assert interface{} to *entry. We get a pointer to the entry.
  - .value = value writes to the struct that pointer points to. We MUTATE the entry.
  - The entry lives inside the list element; we're updating it in place.
*/
func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		elem.Value.(*entry[K, V]).value = value
		c.evictList.MoveToFront(elem)
		return
	}

	if len(c.items) >= c.capacity {
		back := c.evictList.Back()
		if back != nil {
			ent := back.Value.(*entry[K, V])
			delete(c.items, ent.key)
			c.evictList.Remove(back)
		}
	}

	// &entry{...} creates struct, takes address. PushFront stores that pointer.
	elem := c.evictList.PushFront(&entry[K, V]{key: key, value: value})
	c.items[key] = elem
}

/*
Len - Current Number of Cached Items

Thread-safe; returns len(items) which equals evictList length by invariant.
*/
func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}
