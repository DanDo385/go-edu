//go:build !solution && !reference

package pointerszerovaluesnilgotchas

// TODO: Import required packages (none needed for this exercise)

// ============================================================================
// POINTERS, ZERO VALUES, AND NIL: Understanding Go's Memory Model
// ============================================================================
//
// Go's approach to pointers, zero values, and nil:
// 1. Every type has a ZERO VALUE (the default when not explicitly initialized)
// 2. Pointers allow sharing memory between functions (pass by reference)
// 3. nil is the zero value for pointers, slices, maps, channels, interfaces
// 4. You can call methods on nil pointers (useful pattern!)
//
// Zero values by type:
// - Numbers (int, float): 0
// - Booleans: false
// - Strings: ""
// - Pointers: nil
// - Slices: nil (length 0, capacity 0)
// - Maps: nil (can read but NOT write)
// - Channels: nil (sends and receives block forever)
// - Interfaces: nil (both type and value are nil)
//
// Pointer fundamentals:
// - & operator: Takes address of a variable (creates a pointer)
// - * operator (in type): Declares a pointer type (*int)
// - * operator (before pointer): Dereferences (accesses the value)
// - Pointers are 8 bytes on 64-bit systems (just a memory address)
//
// Memory allocation:
// - var x int creates a variable on stack (fast, automatic cleanup)
// - var p *int = &x creates pointer to stack variable
// - new(int) allocates on heap, returns pointer (slower, GC cleanup)
// - Compiler decides stack vs heap (escape analysis)
//
// Why use pointers?
// 1. Modify the original value (not a copy)
// 2. Share large structs without copying (performance)
// 3. Represent optional values (nil = "not present")
// 4. Build recursive data structures (linked lists, trees)
//
// Common gotchas:
// - nil map can READ (returns zero value) but CANNOT WRITE (panic!)
// - nil slice can READ length/capacity (both 0) and APPEND (allocates)
// - nil interface is DIFFERENT from interface containing nil pointer
// - Dereferencing nil pointer = panic (always check before dereferencing)
//
// ============================================================================

// ============================================================================
// Exercise 1: Safe Pointer Dereferencing
// ============================================================================

// SafeDeref safely dereferences a pointer, returning a default value if nil.
// TODO: Implement thread-safe pointer dereferencing.
//
// This is a common pattern for handling optional values represented as pointers.
//
// Parameters:
//   - p: Pointer to an integer (may be nil)
//   - defaultValue: Value to return if p is nil
//
// Returns: Either *p (if not nil) or defaultValue (if nil)
//
// Steps to implement:
//
// 1. Check if pointer is nil
//    - Use: if p == nil { ... }
//    - nil comparison is SAFE (doesn't dereference)
//    - nil check must come BEFORE any dereference
//
// 2. Return default value if nil
//    - Use: return defaultValue
//    - This prevents dereferencing nil (which would panic)
//
// 3. Dereference and return if not nil
//    - Use: return *p
//    - The * operator accesses the value at the memory address
//    - Only safe AFTER confirming p != nil
//
// Key Go concepts:
// - Pointers can be nil (no memory address)
// - Dereferencing nil pointer = runtime panic
// - Always check before dereferencing untrusted pointers
// - This pattern is common in API design (nil = "use default")
//
// Memory behavior:
// - If p is nil: No memory access, just returns defaultValue
// - If p is valid: Reads 8 bytes from address p points to
//
// Example usage:
//   x := 42
//   SafeDeref(&x, 0)  // Returns 42 (dereferences pointer)
//   SafeDeref(nil, 0) // Returns 0 (uses default)
//
// TODO: Implement the SafeDeref function below
// func SafeDeref(p *int, defaultValue int) int {
//     if p == nil {
//         return defaultValue
//     }
//     return *p
// }

// ============================================================================
// Exercise 2: Swapping Values with Pointers
// ============================================================================

// Swap exchanges the values of two integers using pointers.
// TODO: Implement value swapping using pointer dereferencing.
//
// This demonstrates why pointers are needed for functions that modify parameters.
// Without pointers, Go passes copies (call by value).
//
// Parameters:
//   - a: Pointer to first integer
//   - b: Pointer to second integer
//
// After calling Swap(&x, &y):
//   - x contains y's old value
//   - y contains x's old value
//
// Steps to implement:
//
// 1. Save one value in a temporary variable
//    - Use: tmp := *a
//    - Dereference a to get its value
//    - Store in local variable (on stack)
//    - Necessary because we'll overwrite *a
//
// 2. Overwrite first value with second
//    - Use: *a = *b
//    - Dereference b to get its value
//    - Dereference a to SET the value at that address
//    - Now *a and *b have the same value (b's original)
//
// 3. Overwrite second value with saved temporary
//    - Use: *b = tmp
//    - Dereference b to SET the value at that address
//    - Now swap is complete
//
// Key Go concepts:
// - Functions receive COPIES of parameters
// - To modify the original, pass a pointer
// - * on left side (assignment) writes through pointer
// - * on right side (expression) reads through pointer
//
// Memory behavior:
// Before swap: a → addr1 (value: X), b → addr2 (value: Y)
// After step 1: tmp = X (stack variable)
// After step 2: a → addr1 (value: Y), b → addr2 (value: Y)
// After step 3: a → addr1 (value: Y), b → addr2 (value: X)
//
// Why we need tmp:
// - Without it: *a = *b makes both equal
// - Then *b = *a would copy the same value back
// - tmp preserves the original value of *a
//
// TODO: Implement the Swap function below
// func Swap(a, b *int) {
//     tmp := *a
//     *a = *b
//     *b = tmp
// }

// ============================================================================
// Exercise 3: Nil Map Handling
// ============================================================================

// InitializeMap creates a ready-to-use map if the input is nil.
// TODO: Implement safe map initialization.
//
// This demonstrates the nil map gotcha:
// - var m map[string]int creates a NIL map
// - Reading from nil map: OK (returns zero value)
// - Writing to nil map: PANIC!
//
// Parameters:
//   - m: Map that might be nil
//
// Returns: Either m (if not nil) or a new initialized map
//
// Steps to implement:
//
// 1. Check if map is nil
//    - Use: if m == nil { ... }
//    - Maps are reference types (pointer to hash table)
//    - nil means no hash table allocated
//
// 2. Create and return new map if nil
//    - Use: return make(map[string]int)
//    - make() allocates the hash table on heap
//    - Returns a map ready for reads AND writes
//    - Initial capacity is small (grows automatically)
//
// 3. Return existing map if not nil
//    - Use: return m
//    - Already initialized, safe to use
//
// Key Go concepts:
// - Maps are reference types (like slices, channels)
// - Zero value (nil) is NOT the same as empty map
// - nil map: Can read, CANNOT write
// - Empty map: Can read AND write
// - make() vs var:
//   * var m map[string]int → nil map (no allocation)
//   * m := make(map[string]int) → empty map (allocated)
//
// Memory behavior:
// - nil map: Just a nil pointer (0 bytes)
// - Empty map: Hash table structure (24+ bytes)
// - Writing to nil map tries to dereference nil → panic
//
// Common pattern (lazy initialization):
//   func (c *Cache) Get(key string) int {
//       c.data = InitializeMap(c.data)  // Safe to call multiple times
//       return c.data[key]
//   }
//
// TODO: Implement the InitializeMap function below
// func InitializeMap(m map[string]int) map[string]int {
//     if m == nil {
//         return make(map[string]int)
//     }
//     return m
// }

// ============================================================================
// Exercise 4: Linked List Operations
// ============================================================================

// Node represents a node in a singly-linked list.
// Each node contains:
// - Value: The data stored in this node
// - Next: Pointer to the next node (nil if this is the tail)
type Node struct {
	Value int
	Next  *Node
}

// AppendNode appends a new node with the given value to the end of the list.
// TODO: Implement linked list append operation.
//
// This demonstrates:
// - Recursive data structures (Node contains *Node)
// - Pointer manipulation
// - Nil as "end of list" marker
//
// Parameters:
//   - head: Pointer to the first node (may be nil for empty list)
//   - value: The value to append
//
// Returns: Pointer to the head of the list (might be new if input was nil)
//
// Steps to implement:
//
// 1. Create a new node
//    - Use: newNode := &Node{Value: value}
//    - &Node{...} allocates on heap, returns pointer
//    - Next is nil by default (zero value)
//    - This will become the new tail
//
// 2. Handle empty list case
//    - Use: if head == nil { return newNode }
//    - Empty list = nil pointer
//    - New node becomes the head (and tail)
//    - Return it so caller can update their reference
//
// 3. Find the last node
//    - Use: cur := head
//    - Start at the head
//    - Loop: for cur.Next != nil { cur = cur.Next }
//    - Traverse until we find a node with Next = nil
//    - That's the current tail
//
// 4. Append new node to the tail
//    - Use: cur.Next = newNode
//    - Updates the last node's Next pointer
//    - Makes newNode the new tail
//
// 5. Return the original head
//    - Use: return head
//    - Head didn't change (only tail changed)
//    - But we return it for API consistency
//
// Key Go concepts:
// - Recursive types MUST use pointers (type Node { Next Node } is illegal)
// - nil represents "no next node" (end of list)
// - Heap allocation for nodes (they outlive the function)
// - Pointer updates modify the original structure
//
// Memory behavior:
// - Each node: 16 bytes (8 for Value, 8 for Next pointer)
// - Nodes are allocated on heap (escape analysis)
// - Garbage collected when unreachable
//
// Example usage:
//   var head *Node              // Empty list
//   head = AppendNode(head, 1)  // List: 1 → nil
//   head = AppendNode(head, 2)  // List: 1 → 2 → nil
//   head = AppendNode(head, 3)  // List: 1 → 2 → 3 → nil
//
// Alternative patterns:
// - Could use value receiver and return bool (whether head changed)
// - Could keep separate head/tail pointers (faster append)
// - Could use doubly-linked list (*Prev and *Next)
//
// TODO: Implement the AppendNode function below
// func AppendNode(head *Node, value int) *Node {
//     newNode := &Node{Value: value}
//     if head == nil {
//         return newNode
//     }
//     cur := head
//     for cur.Next != nil {
//         cur = cur.Next
//     }
//     cur.Next = newNode
//     return head
// }

// ListLength returns the number of nodes in a linked list.
// TODO: Implement linked list length calculation.
//
// This demonstrates:
// - Traversing a linked list
// - Nil-safe operations
// - Iterative approach (preferred over recursion for long lists)
//
// Parameters:
//   - head: Pointer to the first node (may be nil for empty list)
//
// Returns: Number of nodes in the list (0 for nil)
//
// Steps to implement:
//
// 1. Initialize counter
//    - Use: count := 0
//    - This will track how many nodes we've seen
//
// 2. Loop while head is not nil
//    - Use: for head != nil { ... }
//    - This condition makes it nil-safe
//    - Empty list (head == nil) → loop never executes → count = 0
//
// 3. Increment counter for each node
//    - Use: count++
//    - We've seen one more node
//
// 4. Advance to next node
//    - Use: head = head.Next
//    - Move pointer to the next node
//    - Eventually head becomes nil (end of list)
//    - Then loop terminates
//
// 5. Return the count
//    - Use: return count
//    - Total number of nodes traversed
//
// Key Go concepts:
// - Pointer traversal pattern: for p != nil { ...; p = p.Next }
// - Nil-safe by design (empty list works naturally)
// - Uses the pointer parameter as loop variable (doesn't affect caller)
//
// Memory behavior:
// - Only reads pointers (no writes)
// - No allocation (just traversal)
// - Stack frame is tiny (just count and head pointer)
//
// Time complexity: O(n) - must visit every node
// Space complexity: O(1) - constant space (no recursion)
//
// Why iteration instead of recursion?
// - Recursion: One stack frame per node
//   * For 10,000 nodes → 10,000 stack frames
//   * Risk of stack overflow
// - Iteration: One stack frame total
//   * Works for any list length
//   * More efficient
//
// Recursive alternative (educational, not recommended):
//   func ListLength(head *Node) int {
//       if head == nil {
//           return 0
//       }
//       return 1 + ListLength(head.Next)
//   }
//
// TODO: Implement the ListLength function below
// func ListLength(head *Node) int {
//     count := 0
//     for head != nil {
//         count++
//         head = head.Next
//     }
//     return count
// }

// After implementing all functions:
// - Run: go test ./minis/12-pointers-zero-values-nil-gotchas/exercise/...
// - Check: go test -v for verbose output
// - Experiment with nil pointers to see panics (in a test environment)
// - Compare with solution.go to see detailed explanations
// - Try: go run -gcflags=-m to see escape analysis (stack vs heap)
