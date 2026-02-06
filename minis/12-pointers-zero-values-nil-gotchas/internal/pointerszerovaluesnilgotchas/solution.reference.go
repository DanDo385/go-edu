//go:build reference

package pointerszerovaluesnilgotchas

/*
Reference Solution
==================

This file is intentionally direct: each function demonstrates one memory rule.
Use it as the behavior contract for exercise.go.
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
