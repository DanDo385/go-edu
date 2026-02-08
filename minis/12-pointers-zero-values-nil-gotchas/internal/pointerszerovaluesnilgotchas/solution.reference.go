//go:build reference

package pointerszerovaluesnilgotchas

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

// SafeDeref returns a fallback when the pointer is nil.
func SafeDeref(p *int, defaultValue int) int {
	if p == nil {
		return defaultValue
	}
	// *p follows the address stored in p and reads that int cell.
	return *p
}

// Swap exchanges two ints in place through their addresses.
func Swap(a, b *int) {
	if a == nil || b == nil {
		return
	}
	// Tuple assignment evaluates right side first, then writes left side.
	*a, *b = *b, *a
}

// InitializeMap makes nil maps writable while preserving existing maps.
func InitializeMap(m map[string]int) map[string]int {
	if m == nil {
		return make(map[string]int)
	}
	// Maps are descriptors; returning m keeps sharing the same underlying table.
	return m
}

// AppendNode appends at tail and returns the (possibly new) head.
func AppendNode(head *Node, value int) *Node {
	newNode := &Node{Value: value}
	if head == nil {
		return newNode
	}

	current := head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode
	return head
}

// ListLength counts nodes by walking pointers until nil.
func ListLength(head *Node) int {
	length := 0
	for current := head; current != nil; current = current.Next {
		length++
	}
	return length
}
