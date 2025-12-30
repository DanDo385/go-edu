//go:build solution
// +build solution

/*
Problem: Understanding Go's pointer semantics and nil handling

Requirements:
1. Safe pointer dereferencing with nil checks
2. In-place value swapping using pointers
3. Nil map initialization and safe usage
4. Linked list operations with nil receivers
5. Pointer vs value receiver trade-offs

Data Structure:
- Pointer: 8 bytes (memory address on 64-bit systems)
- nil: Zero value for pointers, slices, maps, channels, interfaces
- Linked list: Recursive structure using pointers

Algorithm: Nil-Safe Operations
- Always check nil before dereferencing
- Methods can be called on nil receivers
- nil map: Can read but NOT write (panic!)
- nil slice: Can read, append (allocates)

Why pointers are essential:
- Modify original value (not a copy)
- Share large structs efficiently
- Represent optional values (nil = absent)
- Build recursive data structures
*/

package exercise

// SafeDeref safely dereferences a pointer.
// BREAKPOINT: Set breakpoint here to trace nil handling
// DEBUG: Watch 'p' to see if pointer is nil
// DEBUG: Watch 'defaultValue' for fallback value
func SafeDeref(p *int, defaultValue int) int {
	// BREAKPOINT: Set breakpoint here to check nil
	// DEBUG: Watch 'p == nil' comparison
	if p == nil {
		// DEBUG: Pointer is nil - returning default
		return defaultValue
	}

	// BREAKPOINT: Set breakpoint here before dereference
	// DEBUG: Watch '*p' to see dereferenced value
	// DEBUG: Dereferencing is safe here (p != nil)
	return *p
}

// Swap exchanges two integer values.
// BREAKPOINT: Set breakpoint here to trace swapping
// DEBUG: Watch 'a' and 'b' pointers
// DEBUG: Watch '*a' and '*b' values before and after swap
func Swap(a, b *int) {
	// BREAKPOINT: Set breakpoint here before saving temp
	// DEBUG: Watch 'temp' capture value from *a
	temp := *a
	// DEBUG: Save *a value to prevent losing it

	// BREAKPOINT: Set breakpoint here for first assignment
	// DEBUG: Watch '*a' receive value from *b
	*a = *b
	// DEBUG: Now *a == *b (both have original b value)

	// BREAKPOINT: Set breakpoint here for second assignment
	// DEBUG: Watch '*b' receive saved temp value
	*b = temp
	// DEBUG: Swap complete - values exchanged
}

// InitializeMap creates a map if nil.
// BREAKPOINT: Set breakpoint here to trace map initialization
// DEBUG: Watch 'm' to see if map is nil
func InitializeMap(m map[string]int) map[string]int {
	// BREAKPOINT: Set breakpoint here to check nil
	// DEBUG: Watch 'm == nil' comparison
	if m == nil {
		// BREAKPOINT: Set breakpoint here to allocate map
		// DEBUG: Watch make() create new map
		// DEBUG: New map is ready for reads AND writes
		return make(map[string]int)
	}

	// DEBUG: Map already initialized - return as-is
	return m
}

// AppendNode appends to a linked list.
// BREAKPOINT: Set breakpoint here to trace list append
// DEBUG: Watch 'head' pointer (may be nil for empty list)
// DEBUG: Watch 'value' being appended
func AppendNode(head *Node, value int) *Node {
	// BREAKPOINT: Set breakpoint here to create new node
	// DEBUG: Watch 'newNode' allocation
	// DEBUG: &Node{} allocates on heap and returns pointer
	newNode := &Node{Value: value}
	// DEBUG: newNode.Next is nil (zero value)

	// BREAKPOINT: Set breakpoint here to check empty list
	// DEBUG: Watch 'head == nil' for empty list case
	if head == nil {
		// DEBUG: Empty list - new node becomes head
		return newNode
	}

	// BREAKPOINT: Set breakpoint here to find tail
	// DEBUG: Watch 'current' traverse the list
	current := head
	// DEBUG: Start at head, move to tail

	// BREAKPOINT: Set breakpoint here in traversal loop
	// DEBUG: Watch 'current.Next' until nil (end of list)
	for current.Next != nil {
		// DEBUG: Watch 'current' advance through list
		current = current.Next
	}
	// DEBUG: current now points to tail node

	// BREAKPOINT: Set breakpoint here to link new node
	// DEBUG: Watch 'current.Next' get new node
	current.Next = newNode
	// DEBUG: Tail now points to new node

	// DEBUG: Watch 'head' return (unchanged)
	return head
}

// ListLength counts nodes in a list.
// BREAKPOINT: Set breakpoint here to trace list traversal
// DEBUG: Watch 'head' pointer (may be nil for empty list)
func ListLength(head *Node) int {
	// BREAKPOINT: Set breakpoint here to initialize counter
	// DEBUG: Watch 'count' start at 0
	count := 0

	// BREAKPOINT: Set breakpoint here to start traversal
	// DEBUG: Watch 'current' initially equals head
	current := head

	// BREAKPOINT: Set breakpoint here in loop
	// DEBUG: Watch 'current != nil' condition
	for current != nil {
		// BREAKPOINT: Set breakpoint here to increment
		// DEBUG: Watch 'count' increase for each node
		count++

		// BREAKPOINT: Set breakpoint here to advance
		// DEBUG: Watch 'current' move to next node
		current = current.Next
		// DEBUG: Eventually current becomes nil (end of list)
	}

	// DEBUG: Watch final 'count' (total nodes)
	return count
}
