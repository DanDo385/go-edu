//go:build reference

package pointerszerovaluesnilgotchas

/*
Reference Solution - Pointers, Zero Values, and nil Gotchas
=========================================================

This file demonstrates Go's pointer semantics, zero values, and nil handling.
Pointers are addresses that enable shared mutable state, but they introduce
complexity around nil checks and ownership. Understanding these concepts is
critical for memory-safe, efficient Go programs.

This connects to the broader Go ecosystem by showing:
- How Go eliminates manual memory management while preserving control
- Why nil checks are pervasive in Go code (unlike Rust's borrow checker)
- How zero values provide sensible defaults for uninitialized variables
- Why pointer receivers enable mutation in methods

The exercise builds understanding of:
- Memory layout: stack vs heap allocation decisions
- Pointer arithmetic absence: why Go forbids `p++` but allows `*p`
- nil semantics: when pointers are nil and how to handle it safely
- Zero values: Go's automatic initialization preventing undefined behavior
- Method receivers: value vs pointer receiver tradeoffs

Teaching notes:
- Memory/ownership: pointers create aliasing relationships where multiple
  variables can mutate the same memory location. Understanding ownership
  boundaries prevents data races and use-after-free bugs.
- Invariants: nil checks establish safety guarantees before dereferencing.
  Zero values provide predictable behavior for uninitialized variables.
- Error surfaces: nil pointer dereferences cause panics, so defensive nil
  checks prevent runtime crashes. This is Go's approach to memory safety.
*/

/*
SafeDeref - Null-Safe Pointer Dereference

Deep explanation of * and nil (per .cursorrules):

What is p?  p is a *int — a pointer. In memory, a pointer is an address (e.g. 0x14000123450).
That address refers to a location where an int lives. The pointer does NOT contain the int itself.

What is nil?  nil is the zero value for pointers. It means "no address" / "points nowhere".
Dereferencing nil is undefined: the CPU would try to read memory at address 0 → panic.

BEFORE dereference: p holds either (a) a valid address, or (b) nil.
AFTER *p: we read the int at that address. If p was nil, we never reach *p (we return early).

Plain English: "Go to the address stored in p; if there is no address, use the default."
*/
func SafeDeref(p *int, defaultValue int) int {
	// Nil check: does p contain a valid address?
	// p == nil means p points nowhere. *p would panic (read from address 0).
	if p == nil {
		return defaultValue
	}

	// *p = dereference. Follow the address in p, read the int there.
	// The * operator: "give me the value at the address p holds."
	// This does NOT create a copy of p — we copy the int value we read.
	return *p
}

/*
Swap - In-Place Exchange Through Pointers

Deep explanation of * for mutation (per .cursorrules):

Memory before call (caller):  x=5, y=10. &x and &y are addresses of those ints.
We pass &x and &y — we pass ADDRESSES, not copies of 5 and 10.

Inside Swap, a and b hold those addresses. *a means "the int at a's address."
*a, *b = *b, *a:
  1. Right side evaluated: read *b (10), read *a (5).
  2. Assign: write 10 to *a (memory at a), write 5 to *b (memory at b).

BEFORE: memory at a = 5, memory at b = 10.
AFTER:  memory at a = 10, memory at b = 5.

This does NOT swap the pointers a and b — we swap the VALUES in the locations
they point to. The caller's x and y are modified in place.
*/
func Swap(a, b *int) {
	if a == nil || b == nil {
		return
	}

	// *a = value at a's address; *b = value at b's address.
	// Assignment writes TO those addresses. Caller's variables change.
	*a, *b = *b, *a
}

/*
InitializeMap - Nil Map Handling Pattern

This function demonstrates how Go handles nil maps and the make() initialization pattern.
Maps in Go are reference types implemented as descriptors to hash tables.

Why this pattern exists:
- Maps can be nil (zero value) or initialized (via make)
- Operations on nil maps cause panics
- This function provides safe initialization while preserving existing maps
*/
func InitializeMap(m map[string]int) map[string]int {
	// Step 1: Check if map is nil
	// Maps have nil as their zero value (uninitialized state)
	// Operations like m["key"] = value panic on nil maps
	if m == nil {
		// Create new map with make()
		// make(map[string]int) allocates the hash table and returns descriptor
		// Without make(), the map would remain nil
		return make(map[string]int)
	}

	// Step 2: Return existing map
	// If map is already initialized, return it unchanged
	// This preserves existing data while ensuring the map is writable
	// Maps are reference types - returning m shares the same hash table
	return m
}

/*
AppendNode - Linked List Construction

Deep explanation of & and * in linked lists (per .cursorrules):

&Node{Value: value} — the & (address-of) operator:
  - Node{Value: value} creates a struct in memory (likely heap, it escapes).
  - & takes the ADDRESS of that struct. Result type: *Node.
  - BEFORE: no Node exists. AFTER: Node exists, & gives us a pointer to it.
  - Think of it like: "Where does this struct live? Give me its address."

current.Next — the . operator on a pointer:
  - current is *Node. current.Next is shorthand for (*current).Next.
  - We follow the pointer (current) to the struct, then access .Next.
  - .Next is a *Node — another pointer. nil means "end of list."

current.Next = newNode — we WRITE to the struct that current points to.
  - We are mutating the last node's Next field to point to the new node.
  - The caller's list is extended; we share the same chain of nodes.
*/
func AppendNode(head *Node, value int) *Node {
	// &Node{...}: create Node, take its address. newNode points to that Node.
	newNode := &Node{Value: value}

	if head == nil {
		return newNode
	}

	// current := head — we copy the pointer (address), not the Node.
	// current.Next follows the pointer, then reads .Next. Nil = end.
	current := head
	for current.Next != nil {
		current = current.Next
	}

	// current points to last node. current.Next = newNode mutates that node.
	current.Next = newNode

	return head
}

/*
ListLength - Linked List Traversal

This function counts nodes in a linked list by traversing pointers.
It demonstrates the fundamental pattern of walking pointer chains.

Why traversal matters:
- Shows how to visit every element in a linked structure
- Demonstrates nil as termination condition
- Illustrates pointer following in loops
- Foundation for more complex list algorithms
*/
func ListLength(head *Node) int {
	// Step 1: Initialize counter
	// Start from 0 and increment for each node visited
	length := 0

	// Step 2: Traverse the list
	// for current := head; current != nil; current = current.Next
	// This is the canonical linked list traversal pattern in Go
	//
	// Initialization: current = head (start at beginning)
	// Condition: current != nil (stop when we reach the end)
	// Increment: current = current.Next (move to next node)
	for current := head; current != nil; current = current.Next {
		// Count this node by incrementing length
		// Each iteration processes one node in the chain
		length++
	}

	// Step 3: Return total count
	// length now contains the number of nodes traversed
	return length
}
