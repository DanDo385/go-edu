//go:build !solution && !reference

package genericlrucache

import (
	"container/list"
	"sync"
	"time"
)

func New(capacity int, defaultTTL time.Duration) *interface{} {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func (c *interface{}) Get(key K) (V, bool) {
	// TODO: Implement this function
	panic("not implemented")
}

func (c *interface{}) Set(key K, val V) {
	// TODO: Implement this function
	panic("not implemented")
}

func (c *interface{}) SetWithTTL(key K, val V, ttl time.Duration) {
	// TODO: Implement this function
	panic("not implemented")
}

func (c *interface{}) Len() int {
	// TODO: Implement this function
	panic("not implemented")
}

func (c *interface{}) removeElement(elem *list.Element) {
	// TODO: Implement this function
	panic("not implemented")
}
