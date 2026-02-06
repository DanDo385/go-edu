//go:build !solution && !reference

package pointerszerovaluesnilgotchas

/*
Problem: Understanding Go's pointer semantics and nil handling
Requirements:
1. Safe pointer dereferencing with nil checks
2. In-place value swapping using pointers
3. Nil map initialization and safe usage
Algorithm: Nil-Safe Operations
- Always check nil before dereferencing
- Methods can be called on nil receivers
- nil map: Can read but NOT write (panic!)
- nil slice: Can read, append (allocates)
*/

// SafeDeref - TODO: implement this function
func SafeDeref(p *int, defaultValue int) int {
	if p == nil {
		return defaultValue
	}
	return *p
}

// Swap - TODO: implement this function
func Swap(a, b *int) {
	if a == nil || b == nil {
		return
	}
	*a, *b = *b, *a
}

// InitializeMap - TODO: implement this function
func InitializeMap(m map[string]int) map[string]int {
	if m == nil {
		return make(map[string]int)
	}
	return m
}

// AppendNode - TODO: implement this function
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

// ListLength - TODO: implement this function
func ListLength(head *Node) int {
	length := 0
	for current := head; current != nil; current = current.Next {
		length++
	}
	return length
}
