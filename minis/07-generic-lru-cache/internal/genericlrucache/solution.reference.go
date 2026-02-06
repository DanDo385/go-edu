//go:build reference

package genericlrucache

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

func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}
