// Package main demonstrates pointers, zero values, and nil gotchas in Go.
//
// MACRO-COMMENT: What This Program Teaches
// =========================================
// This program shows practical examples of:
// 1. Pointer basics (addressing, dereferencing, modification)
// 2. Zero values for all Go types
// 3. Nil pointer dereference dangers
// 4. Pointer vs value receivers
// 5. new, make, and composite literals
//
// COMPILER BEHAVIOR:
// When you run with -gcflags='-m', you'll see which variables escape to the heap.
// Variables that are returned as pointers or stored in heap structures escape.
//
// RUNTIME BEHAVIOR:
// Dereferencing a nil pointer causes an immediate panic (memory protection).
// This is safer than C, where nil dereferences cause undefined behavior.

package main

import (
	"fmt"
)

// ============================================================================
// SECTION 1: Pointer Basics
// ============================================================================

// demonstratePointerBasics shows the & and * operators.
//
// MACRO-COMMENT: Understanding Memory Addresses
// Every variable lives at a specific memory address. Pointers let you:
// 1. Pass large data efficiently (just pass the address)
// 2. Modify data in-place (no copying)
// 3. Share data between functions (same address, same data)
func demonstratePointerBasics() {
	fmt.Println("=== Pointer Basics ===")

	// MICRO-COMMENT: Create a regular integer variable
	// This allocates space for an int and stores 42 in it
	x := 42

	// MICRO-COMMENT: The & operator gets the memory address of x
	// p is a "pointer to int" (*int type)
	// It stores the address where x lives in memory
	p := &x

	fmt.Printf("Value of x: %d\n", x)
	fmt.Printf("Address of x: %p\n", p)  // %p formats pointers as hex addresses

	// MICRO-COMMENT: The * operator dereferences the pointer
	// *p means "go to the address in p and get the value there"
	fmt.Printf("Value at address p: %d\n", *p)

	// MACRO-COMMENT: Modifying Through Pointers
	// This is the key power of pointers: changing data indirectly
	// *p = 100 means "go to the address in p and SET the value to 100"
	// Since p points to x, this changes x
	*p = 100

	fmt.Printf("\nAfter *p = 100:\n")
	fmt.Printf("  x: %d\n", x)   // 100 (x was modified!)
	fmt.Printf("  *p: %d\n", *p) // 100 (same value, same memory location)

	fmt.Println()
}

// ============================================================================
// SECTION 2: Zero Values
// ============================================================================

// demonstrateZeroValues shows the zero value for each type.
//
// MACRO-COMMENT: Go's Safety Feature
// Unlike C/C++, Go ALWAYS initializes variables to a meaningful zero value.
// This prevents "uninitialized variable" bugs that plague C programs.
//
// MEMORY LAYOUT:
// When you declare `var x int`, the compiler:
// 1. Allocates space for an int (typically 8 bytes)
// 2. Writes zeros to all 8 bytes
// 3. This represents the integer 0
func demonstrateZeroValues() {
	fmt.Println("=== Zero Values for All Types ===")

	// MICRO-COMMENT: Numeric types zero to 0
	var i int
	var f float64
	fmt.Printf("int:     %d\n", i)  // 0
	fmt.Printf("float64: %f\n", f)  // 0.000000

	// MICRO-COMMENT: Strings zero to empty string
	var s string
	fmt.Printf("string:  %q\n", s)  // ""

	// MICRO-COMMENT: Booleans zero to false
	var b bool
	fmt.Printf("bool:    %t\n", b)  // false

	// MICRO-COMMENT: Pointers zero to nil
	// nil means "doesn't point to anything"
	var p *int
	fmt.Printf("*int:    %v\n", p)  // <nil>

	// MACRO-COMMENT: Composite Types and nil
	// Slices, maps, channels, and functions can be nil
	// But they behave DIFFERENTLY when nil!

	// MICRO-COMMENT: Nil slice (SAFE for reading)
	var slice []int
	fmt.Printf("[]int:   %v (len=%d, cap=%d)\n", slice, len(slice), cap(slice))
	// You can safely call len(), cap(), range over it
	// But it has no backing array allocated

	// MICRO-COMMENT: Nil map (UNSAFE for writing!)
	var m map[string]int
	fmt.Printf("map:     %v\n", m)  // map[]
	// Reading is OK: m["key"] returns 0
	// Writing PANICS: m["key"] = 1 would crash

	// MICRO-COMMENT: Nil channel (UNSAFE - blocks forever!)
	var ch chan int
	fmt.Printf("chan:    %v\n", ch)  // <nil>
	// Sending or receiving on nil channel blocks forever (deadlock)

	fmt.Println()
}

// ============================================================================
// SECTION 3: Nil Pointer Dereference Danger
// ============================================================================

// demonstrateNilDereference shows the classic nil pointer panic (safely).
//
// MACRO-COMMENT: The Most Common Runtime Panic
// Dereferencing a nil pointer is the #1 cause of panics in Go programs.
// ALWAYS check if a pointer is nil before dereferencing it in production code.
//
// RUNTIME BEHAVIOR:
// When you dereference nil, the CPU attempts to access memory address 0x0.
// The operating system protects this address (null page), causing a fault.
// Go's runtime catches the fault and converts it to a panic.
func demonstrateNilDereference() {
	fmt.Println("=== Nil Pointer Dereference (Recovered) ===")

	// MACRO-COMMENT: Using defer/recover to Catch Panics
	// defer schedules a function to run when the current function returns
	// recover() catches a panic and returns its value (returns nil if no panic)
	// This is Go's exception handling mechanism
	defer func() {
		// MICRO-COMMENT: recover() returns nil if no panic occurred
		// If a panic occurred, it returns the panic value and stops the panic
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	// MICRO-COMMENT: Create a nil pointer (zero value)
	var p *int
	fmt.Printf("p is nil: %t\n", p == nil)

	// MICRO-COMMENT: This will panic!
	// The runtime will try to access memory at address nil (0x0)
	// The OS will send a segmentation fault signal
	// Go's runtime will catch it and panic with a descriptive message
	*p = 42  // PANIC: invalid memory address or nil pointer dereference

	// MICRO-COMMENT: This line is never reached (panic stops execution)
	fmt.Println("This line never executes")
}

// ============================================================================
// SECTION 4: Safe Nil Checks
// ============================================================================

// User represents a user in the system.
type User struct {
	Name  string
	Email string
}

// Greet prints a greeting (nil-safe version).
//
// MACRO-COMMENT: Defensive Programming with Pointers
// Always check if a pointer receiver is nil before accessing its fields.
// This is especially important for methods that might be called on nil pointers.
//
// STYLE NOTE:
// Some Go programmers consider nil receiver checks a code smell, preferring
// to ensure pointers are always valid. However, in library code, being defensive
// is often the right choice.
func (u *User) Greet() {
	// MICRO-COMMENT: Guard clause - check for nil first
	// This prevents the panic from u.Name access
	if u == nil {
		fmt.Println("Hello, guest!")
		return
	}

	// MICRO-COMMENT: Safe to access fields now
	fmt.Printf("Hello, %s!\n", u.Name)
}

// demonstrateSafeNilChecks shows nil-safe method calls.
func demonstrateSafeNilChecks() {
	fmt.Println("\n=== Safe Nil Checks ===")

	// MICRO-COMMENT: Create a nil pointer
	var u *User
	fmt.Printf("u is nil: %t\n", u == nil)

	// MICRO-COMMENT: Call method on nil receiver
	// This is SAFE because Greet() checks for nil
	u.Greet()  // Prints "Hello, guest!"

	// MICRO-COMMENT: Create a valid user
	u = &User{Name: "Alice", Email: "alice@example.com"}
	u.Greet()  // Prints "Hello, Alice!"

	fmt.Println()
}

// ============================================================================
// SECTION 5: Value vs Pointer Receivers
// ============================================================================

// Counter demonstrates the difference between value and pointer receivers.
type Counter struct {
	count int
}

// IncrementValue uses a value receiver (receives a COPY).
//
// MICRO-COMMENT: Value Receiver Behavior
// When this method is called, Go copies the entire Counter struct.
// Modifications to c.count affect only the copy, not the original.
// The original Counter is unchanged after this method returns.
func (c Counter) IncrementValue() {
	c.count++  // Modifies the copy
	fmt.Printf("  Inside IncrementValue: %d\n", c.count)
}

// IncrementPointer uses a pointer receiver (receives the address).
//
// MICRO-COMMENT: Pointer Receiver Behavior
// When this method is called, Go passes the address of the Counter.
// Modifications to c.count affect the original Counter.
// The change persists after this method returns.
//
// PERFORMANCE NOTE:
// Even though c is a pointer, you still write c.count (not (*c).count).
// Go automatically dereferences for you (syntactic sugar).
func (c *Counter) IncrementPointer() {
	c.count++  // Modifies the original
	fmt.Printf("  Inside IncrementPointer: %d\n", c.count)
}

// demonstrateReceivers shows value vs pointer receivers.
func demonstrateReceivers() {
	fmt.Println("=== Value vs Pointer Receivers ===")

	// MICRO-COMMENT: Test value receiver
	c1 := Counter{count: 0}
	fmt.Println("Before IncrementValue:", c1.count)
	c1.IncrementValue()
	fmt.Println("After IncrementValue:", c1.count)  // Still 0 (copy was modified)

	fmt.Println()

	// MICRO-COMMENT: Test pointer receiver
	c2 := Counter{count: 0}
	fmt.Println("Before IncrementPointer:", c2.count)
	c2.IncrementPointer()  // Go automatically does (&c2).IncrementPointer()
	fmt.Println("After IncrementPointer:", c2.count)  // Now 1 (original modified)

	fmt.Println()
}

// ============================================================================
// SECTION 6: new vs make vs Composite Literals
// ============================================================================

// demonstrateAllocation shows three ways to allocate memory.
//
// MACRO-COMMENT: Allocation Strategies
// Go provides multiple ways to create values, each optimized for different use cases:
// 1. Composite literals: Most readable, best for known values
// 2. make(): Required for slices, maps, channels (initializes internal structures)
// 3. new(): Rarely used, returns pointer to zero value
func demonstrateAllocation() {
	fmt.Println("=== new vs make vs Composite Literals ===")

	// APPROACH 1: Composite literal (most common)
	// MICRO-COMMENT: Creates a User and returns a pointer to it
	// The compiler allocates memory, initializes fields, returns address
	u1 := &User{Name: "Alice", Email: "alice@example.com"}
	fmt.Printf("Composite literal: %+v\n", u1)

	// APPROACH 2: new() - allocates and zeros
	// MICRO-COMMENT: new(User) allocates a User, zeros all fields, returns *User
	// Equivalent to: var temp User; u2 := &temp
	u2 := new(User)
	u2.Name = "Bob"  // Must set fields manually
	u2.Email = "bob@example.com"
	fmt.Printf("new():             %+v\n", u2)

	// APPROACH 3: make() - for slices, maps, channels ONLY
	// MICRO-COMMENT: make() not only allocates but also initializes internal structures
	// For maps, it creates the hash table
	// For slices, it creates the backing array
	// For channels, it creates the buffer
	m := make(map[string]int)  // Creates an empty, ready-to-use map
	m["key"] = 42
	fmt.Printf("make(map):         %v\n", m)

	s := make([]int, 5)  // Creates a slice with len=5, cap=5, all zeros
	fmt.Printf("make(slice):       %v\n", s)

	// MACRO-COMMENT: Why Can't You Use new() with Maps/Slices?
	// You CAN, but you get a nil map/slice, which is dangerous:
	// badMap := new(map[string]int)  // Returns *map[string]int (pointer to nil map)
	// (*badMap)["key"] = 1  // PANIC: nil map
	// Always use make() for maps, slices, channels.

	fmt.Println()
}

// ============================================================================
// SECTION 7: Escape Analysis
// ============================================================================

// stackAllocated demonstrates a variable that stays on the stack.
//
// MICRO-COMMENT: Stack Allocation
// x is local to this function and not returned, so it lives on the stack.
// When the function returns, x is automatically destroyed (very fast).
func stackAllocated() int {
	x := 42  // Stack-allocated (compiler flag -m will confirm)
	return x  // Return the VALUE, not the address
}

// heapAllocated demonstrates a variable that escapes to the heap.
//
// MICRO-COMMENT: Heap Allocation (Escape)
// x is returned as a pointer, so it must outlive this function.
// The compiler moves x to the heap so it survives after the function returns.
// The garbage collector will eventually reclaim it.
func heapAllocated() *int {
	x := 42   // Escapes to heap (compiler flag -m will confirm)
	return &x  // Return the ADDRESS - x must survive!
}

// demonstrateEscapeAnalysis shows stack vs heap allocation.
//
// MACRO-COMMENT: To see escape analysis, run:
// go build -gcflags='-m' cmd/pointers-demo/main.go
//
// You'll see output like:
// ./main.go:X:X: moved to heap: x
// ./main.go:Y:Y: x does not escape
func demonstrateEscapeAnalysis() {
	fmt.Println("=== Escape Analysis ===")

	// MICRO-COMMENT: Call stack-allocated function
	// The return value is copied to the stack frame of main
	val := stackAllocated()
	fmt.Printf("Stack-allocated value: %d\n", val)

	// MICRO-COMMENT: Call heap-allocated function
	// The return value is a pointer to heap-allocated memory
	ptr := heapAllocated()
	fmt.Printf("Heap-allocated value: %d (address: %p)\n", *ptr, ptr)

	fmt.Println()
}

// ============================================================================
// MAIN FUNCTION
// ============================================================================

// main orchestrates all demonstrations.
//
// MACRO-COMMENT: Learning Progression
// The demos build from simple to complex:
// 1. Pointer basics (foundation)
// 2. Zero values (safety model)
// 3. Nil dangers (common pitfalls)
// 4. Safe nil checks (defensive programming)
// 5. Receivers (method design)
// 6. Allocation (memory management)
// 7. Escape analysis (performance)
func main() {
	demonstratePointerBasics()
	demonstrateZeroValues()
	demonstrateNilDereference()  // Panics and recovers
	demonstrateSafeNilChecks()
	demonstrateReceivers()
	demonstrateAllocation()
	demonstrateEscapeAnalysis()

	// MACRO-COMMENT: Try running with:
	// go build -gcflags='-m -m' cmd/pointers-demo/main.go
	// This shows detailed escape analysis decisions by the compiler
}
