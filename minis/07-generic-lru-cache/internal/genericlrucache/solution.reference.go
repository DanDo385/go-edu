//go:build reference

package genericlrucache

/*
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
*/

import (
	"container/list"
	"sync"
)

type Cache[K comparable, V any] struct {
	mu        sync.Mutex
	capacity  int
	items     map[K]*list.Element
	evictList *list.List
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

// New implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
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

// Get implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.evictList.MoveToFront(elem)
		return elem.Value.(*entry[K, V]).value, true
	}
	var zeroV V
	return zeroV, false
}

// Set implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
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

	elem := c.evictList.PushFront(&entry[K, V]{key: key, value: value})
	c.items[key] = elem
}

// Len implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}
