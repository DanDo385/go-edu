//go:build !solution
// +build !solution

// Package exercise contains hands-on exercises for pointers and zero values.

package exercise

// SafeDeref safely dereferences a pointer, returning a default value if nil.
//
// REQUIREMENTS:
// - If p is nil, return defaultValue
// - If p is not nil, return *p
//
// EXAMPLE:
//   x := 42
//   SafeDeref(&x, 0)  // returns 42
//   SafeDeref(nil, 0) // returns 0
func SafeDeref(p *int, defaultValue int) int {
	// TODO: Implement this
	return 0
}

// Swap exchanges the values of two integers using pointers.
//
// REQUIREMENTS:
// - Modify the values at addresses a and b so they are swapped
// - After calling Swap(&x, &y), x should have y's old value and vice versa
//
// HINT: Use a temporary variable to hold one value during the swap
func Swap(a, b *int) {
	// TODO: Implement this
}

// InitializeMap creates a ready-to-use map if the input is nil.
//
// REQUIREMENTS:
// - If m is nil, create and return a new map
// - If m is not nil, return it unchanged
//
// This pattern is useful for lazy initialization.
func InitializeMap(m map[string]int) map[string]int {
	// TODO: Implement this
	return nil
}

// Node represents a node in a linked list.
type Node struct {
	Value int
	Next  *Node
}

// AppendNode appends a new node with the given value to the end of the list.
//
// REQUIREMENTS:
// - If the list is empty (head is nil), create a new head node
// - Otherwise, traverse to the end and append a new node
// - Return the head of the list (which may be a new node if input was nil)
//
// EXAMPLE:
//   var head *Node
//   head = AppendNode(head, 1)  // Creates: 1 -> nil
//   head = AppendNode(head, 2)  // Creates: 1 -> 2 -> nil
//   head = AppendNode(head, 3)  // Creates: 1 -> 2 -> 3 -> nil
func AppendNode(head *Node, value int) *Node {
	// TODO: Implement this
	return nil
}

// ListLength returns the number of nodes in a linked list.
//
// REQUIREMENTS:
// - If head is nil, return 0
// - Otherwise, traverse the list and count nodes
//
// HINT: Use a loop, not recursion (to avoid stack overflow for long lists)
func ListLength(head *Node) int {
	// TODO: Implement this
	return 0
}
