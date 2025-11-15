// Package main demonstrates value vs pointer receivers in Go methods.
//
// MACRO-COMMENT: What This Program Teaches
// =========================================
// This program provides comprehensive demonstrations of:
// 1. Value receivers (copy semantics, read-only operations)
// 2. Pointer receivers (reference semantics, mutations)
// 3. Interface satisfaction and method sets
// 4. Performance differences between receiver types
// 5. Nil receiver handling
// 6. Common pitfalls and how to avoid them
//
// COMPILER BEHAVIOR:
// When you run with -gcflags='-m', you'll see which values escape to the heap.
// Pointer receivers may force heap allocation, while value receivers can stay on stack.
//
// RUNTIME BEHAVIOR:
// Value receivers copy the entire struct (expensive for large types).
// Pointer receivers copy only the pointer (8 bytes on 64-bit systems).

package main

import (
	"fmt"
	"time"
)

// ============================================================================
// SECTION 1: Value Receivers (Copy Semantics)
// ============================================================================

// Counter demonstrates value receiver behavior.
// MICRO-COMMENT: This is a simple struct with one field.
// We'll show why value receivers don't modify the original.
type Counter struct {
	count int
}

// IncrementValue uses a VALUE RECEIVER.
//
// MACRO-COMMENT: Copy Semantics
// When this method is called:
// 1. Go COPIES the entire Counter struct
// 2. The method modifies the COPY
// 3. The copy is discarded when the method returns
// 4. The original Counter remains unchanged
//
// MEMORY LAYOUT:
// Caller's stack: Counter{count: 0}
// Method's stack: Counter{count: 0} ← COPY
//                 After increment: {count: 1}
// After return:   Caller still has {count: 0}
func (c Counter) IncrementValue() {
	// MICRO-COMMENT: This increments the COPY, not the original
	c.count++
	// The modified copy is discarded when this function returns
}

// GetCount is a safe read-only operation with a value receiver.
// MICRO-COMMENT: Since we're only reading (not modifying), a value receiver
// is fine. The copy is cheap (just 8 bytes for one int).
func (c Counter) GetCount() int {
	return c.count
}

// demonstrateValueReceiver shows that value receivers don't modify the original.
func demonstrateValueReceiver() {
	fmt.Println("=== Value Receiver: No Modification ===")

	// MICRO-COMMENT: Create a counter with count = 0
	c := Counter{count: 0}
	fmt.Printf("Before: Counter{count: %d}\n", c.count)

	// MICRO-COMMENT: Call IncrementValue multiple times
	// Each call works on a COPY, so the original is never changed
	c.IncrementValue()
	c.IncrementValue()
	c.IncrementValue()

	// MICRO-COMMENT: The original is STILL 0 (unchanged!)
	// This is often a bug in real code when the programmer expected modification
	fmt.Printf("After:  Counter{count: %d}  ← Unchanged!\n", c.count)
	fmt.Println()
}

// ============================================================================
// SECTION 2: Pointer Receivers (Reference Semantics)
// ============================================================================

// IncrementPointer uses a POINTER RECEIVER.
//
// MACRO-COMMENT: Reference Semantics
// When this method is called:
// 1. Go passes a POINTER to the Counter (just 8 bytes)
// 2. The method modifies the original through the pointer
// 3. Changes persist after the method returns
//
// MEMORY LAYOUT:
// Caller's stack: Counter{count: 0} at address 0xc000014098
// Method's stack: *Counter = 0xc000014098 (pointer to the original)
//                 (*c).count++ modifies the original directly
func (c *Counter) IncrementPointer() {
	// MICRO-COMMENT: c is a pointer, so c.count is shorthand for (*c).count
	// Go automatically dereferences the pointer for us
	c.count++
	// The original Counter is modified
}

// Reset is another mutating method, so it needs a pointer receiver.
// MICRO-COMMENT: Consistency rule: If one method needs a pointer receiver,
// use pointer receivers for ALL methods on that type.
func (c *Counter) Reset() {
	c.count = 0
}

// demonstratePointerReceiver shows that pointer receivers modify the original.
func demonstratePointerReceiver() {
	fmt.Println("=== Pointer Receiver: Modification Works ===")

	// MICRO-COMMENT: Create a counter with count = 0
	c := Counter{count: 0}
	fmt.Printf("Before: Counter{count: %d}\n", c.count)

	// MACRO-COMMENT: Go's Syntactic Sugar
	// Even though c is a VALUE and IncrementPointer expects a POINTER,
	// Go automatically converts c.IncrementPointer() to (&c).IncrementPointer()
	c.IncrementPointer()
	c.IncrementPointer()
	c.IncrementPointer()
	c.IncrementPointer()
	c.IncrementPointer()

	// MICRO-COMMENT: The original IS modified (changed to 5!)
	fmt.Printf("After:  Counter{count: %d}  ← Changed!\n", c.count)

	// MICRO-COMMENT: Reset also uses a pointer receiver
	c.Reset()
	fmt.Printf("After Reset: Counter{count: %d}\n", c.count)
	fmt.Println()
}

// ============================================================================
// SECTION 3: Interface Satisfaction and Method Sets
// ============================================================================

// Printer is an interface with one method.
// MICRO-COMMENT: This interface will help us demonstrate method set rules.
type Printer interface {
	Print()
}

// Document has a pointer receiver method.
type Document struct {
	content string
}

// Print uses a POINTER RECEIVER.
// MICRO-COMMENT: This means only *Document satisfies the Printer interface.
// Document (the value type) does NOT satisfy Printer.
func (d *Document) Print() {
	// MICRO-COMMENT: Check for nil receiver (defensive programming)
	if d == nil {
		fmt.Println("[nil document]")
		return
	}
	fmt.Println(d.content)
}

// Book has a value receiver method.
type Book struct {
	title string
}

// Print uses a VALUE RECEIVER.
// MICRO-COMMENT: This means BOTH Book and *Book satisfy the Printer interface.
func (b Book) Print() {
	fmt.Printf("Book: %s\n", b.title)
}

// demonstrateInterfaceSatisfaction shows method set rules.
//
// MACRO-COMMENT: The Method Set Rules
// For a type T:
// - Methods with receiver T are in the method set of T and *T
// - Methods with receiver *T are ONLY in the method set of *T
//
// This affects interface satisfaction:
// - If a method has pointer receiver, only the pointer type satisfies interfaces
// - If a method has value receiver, both value and pointer types satisfy interfaces
func demonstrateInterfaceSatisfaction() {
	fmt.Println("=== Interface Satisfaction ===")

	// MICRO-COMMENT: Document with pointer receiver
	var p1 Printer

	// This FAILS to compile (commented out to let the program run):
	// p1 = Document{content: "hello"}  // ❌ ERROR: Document doesn't implement Printer

	// This WORKS:
	p1 = &Document{content: "hello"}  // ✅ *Document implements Printer
	fmt.Print("*Document implements Printer: ")
	p1.Print()

	// MACRO-COMMENT: Why Can We Call d.Print() But Not Assign d to Interface?
	// You CAN call the method on a value (Go takes the address automatically):
	d := Document{content: "world"}
	d.Print()  // ✅ Works: Go does (&d).Print()

	// But you CANNOT assign the value to an interface:
	// p1 = d  // ❌ Compile error!

	// The reason: Go can auto-address when calling because d is addressable.
	// But Go cannot auto-address when assigning to interface because the
	// interface might store an unaddressable value (like a return value).

	fmt.Println()

	// MICRO-COMMENT: Book with value receiver
	var p2 Printer

	// Both of these WORK:
	p2 = Book{title: "Go Programming"}     // ✅ Book implements Printer
	p2.Print()

	p2 = &Book{title: "Advanced Go"}       // ✅ *Book also implements Printer
	p2.Print()

	fmt.Println()
}

// ============================================================================
// SECTION 4: Performance Comparison
// ============================================================================

// SmallStruct is a small type (16 bytes).
type SmallStruct struct {
	a, b int64  // 16 bytes total
}

// LargeStruct is a large type (8000 bytes).
type LargeStruct struct {
	data [1000]int64  // 8000 bytes!
}

// ProcessValue uses a value receiver (copies the entire struct).
// MICRO-COMMENT: For SmallStruct (16 bytes), copying is cheap.
// For LargeStruct (8000 bytes), copying is VERY expensive!
func (s SmallStruct) ProcessValue() int64 {
	return s.a + s.b
}

// ProcessPointer uses a pointer receiver (copies only the pointer).
// MICRO-COMMENT: Always copies just 8 bytes (the pointer), regardless of struct size.
func (s *SmallStruct) ProcessPointer() int64 {
	return s.a + s.b
}

// ProcessValueLarge uses a value receiver on a LARGE struct.
// MICRO-COMMENT: This copies 8000 bytes EVERY time it's called!
func (l LargeStruct) ProcessValueLarge() int64 {
	return l.data[0] + l.data[999]
}

// ProcessPointerLarge uses a pointer receiver on a LARGE struct.
// MICRO-COMMENT: This copies only 8 bytes (the pointer).
func (l *LargeStruct) ProcessPointerLarge() int64 {
	return l.data[0] + l.data[999]
}

// demonstratePerformance shows the performance difference.
//
// MACRO-COMMENT: Performance Rules of Thumb
// For small types (< 64 bytes):
// - Value receiver: ~1-2 ns
// - Pointer receiver: ~1-2 ns
// - Difference: Negligible
//
// For large types (> 64 bytes):
// - Value receiver: Proportional to size (8000 bytes = ~1000 ns)
// - Pointer receiver: ~1-2 ns (constant, just the pointer)
// - Difference: Can be 10x-100x!
func demonstratePerformance() {
	fmt.Println("=== Performance: Small vs Large Structs ===")

	// MICRO-COMMENT: Small struct (16 bytes)
	small := SmallStruct{a: 10, b: 20}
	start := time.Now()
	for i := 0; i < 1000000; i++ {
		small.ProcessValue()
	}
	valueTime := time.Since(start)

	start = time.Now()
	for i := 0; i < 1000000; i++ {
		small.ProcessPointer()
	}
	pointerTime := time.Since(start)

	fmt.Println("Small struct (16 bytes), 1M iterations:")
	fmt.Printf("  Value receiver:   %v\n", valueTime)
	fmt.Printf("  Pointer receiver: %v\n", pointerTime)
	fmt.Printf("  Difference: Negligible\n")
	fmt.Println()

	// MICRO-COMMENT: Large struct (8000 bytes)
	large := LargeStruct{}
	large.data[0] = 100
	large.data[999] = 200

	start = time.Now()
	for i := 0; i < 100000; i++ {
		large.ProcessValueLarge()
	}
	valueLargeTime := time.Since(start)

	start = time.Now()
	for i := 0; i < 100000; i++ {
		large.ProcessPointerLarge()
	}
	pointerLargeTime := time.Since(start)

	fmt.Println("Large struct (8000 bytes), 100K iterations:")
	fmt.Printf("  Value receiver:   %v\n", valueLargeTime)
	fmt.Printf("  Pointer receiver: %v\n", pointerLargeTime)
	if valueLargeTime > pointerLargeTime {
		speedup := float64(valueLargeTime) / float64(pointerLargeTime)
		fmt.Printf("  Speedup: %.1fx faster with pointer receiver!\n", speedup)
	}
	fmt.Println()
}

// ============================================================================
// SECTION 5: Nil Receiver Handling
// ============================================================================

// IntList is a recursive linked list (nil-terminated).
type IntList struct {
	value int
	next  *IntList
}

// Sum calculates the sum of all elements in the list.
//
// MACRO-COMMENT: Nil-Safe Pattern
// This method works even if called on a nil pointer!
// This is a powerful pattern for:
// - Optional values (nil means "not present")
// - Recursive data structures (nil terminates recursion)
// - Default behaviors for uninitialized types
func (l *IntList) Sum() int {
	// MICRO-COMMENT: Check for nil BEFORE dereferencing any fields
	// If l is nil, we treat it as an empty list (sum = 0)
	if l == nil {
		return 0
	}

	// MICRO-COMMENT: l.next might be nil, but Sum() handles that
	// This is recursion with nil-termination
	return l.value + l.next.Sum()
}

// Len returns the length of the list.
func (l *IntList) Len() int {
	if l == nil {
		return 0
	}
	return 1 + l.next.Len()
}

// demonstrateNilReceiver shows nil-safe method patterns.
func demonstrateNilReceiver() {
	fmt.Println("=== Nil Receiver Handling ===")

	// MICRO-COMMENT: Create a list: 1 -> 2 -> 3 -> nil
	list := &IntList{
		value: 1,
		next: &IntList{
			value: 2,
			next: &IntList{
				value: 3,
				next:  nil,  // nil terminates the list
			},
		},
	}

	fmt.Printf("List: ")
	for l := list; l != nil; l = l.next {
		fmt.Printf("%d ", l.value)
	}
	fmt.Println()
	fmt.Printf("Sum: %d\n", list.Sum())
	fmt.Printf("Len: %d\n", list.Len())

	// MACRO-COMMENT: Nil Receiver Works!
	// Even though nilList is nil, we can call methods on it safely
	var nilList *IntList  // nil
	fmt.Printf("\nNil list:\n")
	fmt.Printf("Sum: %d (doesn't panic!)\n", nilList.Sum())
	fmt.Printf("Len: %d\n", nilList.Len())

	fmt.Println()
}

// ============================================================================
// SECTION 6: Common Pitfalls
// ============================================================================

// demonstratePitfalls shows common mistakes with receivers.
func demonstratePitfalls() {
	fmt.Println("=== Common Pitfalls ===")

	// PITFALL 1: Value receiver doesn't modify
	fmt.Println("Pitfall 1: Value receiver doesn't modify")
	c1 := Counter{count: 0}
	c1.IncrementValue()  // BUG: Doesn't modify c1
	fmt.Printf("  After IncrementValue(): count = %d (expected 1, got 0!)\n", c1.count)
	c1.IncrementPointer()  // FIX: Use pointer receiver
	fmt.Printf("  After IncrementPointer(): count = %d ✓\n", c1.count)
	fmt.Println()

	// PITFALL 2: Can't take address of map elements
	fmt.Println("Pitfall 2: Can't take address of map elements")
	m := map[string]Counter{
		"a": {count: 0},
	}
	// m["a"].IncrementPointer()  // ❌ COMPILE ERROR: can't take address

	// FIX 1: Extract, modify, store back
	temp := m["a"]
	temp.IncrementPointer()
	m["a"] = temp
	fmt.Printf("  After extract-modify-store: count = %d\n", m["a"].count)

	// FIX 2: Use pointers in the map
	m2 := map[string]*Counter{
		"b": &Counter{count: 0},
	}
	m2["b"].IncrementPointer()  // ✅ Works!
	fmt.Printf("  Using *Counter in map: count = %d\n", m2["b"].count)
	fmt.Println()
}

// ============================================================================
// SECTION 7: When to Use Each Receiver Type
// ============================================================================

// demonstrateDecisionCriteria shows when to use each receiver type.
func demonstrateDecisionCriteria() {
	fmt.Println("=== Decision Criteria ===")
	fmt.Println()

	fmt.Println("Use VALUE receivers when:")
	fmt.Println("  ✓ The type is small (< 64 bytes)")
	fmt.Println("  ✓ The type is naturally immutable")
	fmt.Println("  ✓ The method doesn't modify the receiver")
	fmt.Println("  ✓ The type is a map, slice, or function (already reference-like)")
	fmt.Println()

	fmt.Println("Use POINTER receivers when:")
	fmt.Println("  ✓ The method modifies the receiver")
	fmt.Println("  ✓ The type is large (> 64 bytes)")
	fmt.Println("  ✓ The type contains a mutex (sync.Mutex) [MUST use pointer!]")
	fmt.Println("  ✓ You need consistency (if one method uses *, use * for all)")
	fmt.Println("  ✓ The type needs initialization (zero value isn't useful)")
	fmt.Println()

	fmt.Println("CONSISTENCY RULE:")
	fmt.Println("  If ANY method needs a pointer receiver, use pointer receivers")
	fmt.Println("  for ALL methods on that type.")
	fmt.Println()
}

// ============================================================================
// MAIN: Run All Demonstrations
// ============================================================================

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║  Methods: Value vs Pointer Receivers                     ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	demonstrateValueReceiver()
	demonstratePointerReceiver()
	demonstrateInterfaceSatisfaction()
	demonstratePerformance()
	demonstrateNilReceiver()
	demonstratePitfalls()
	demonstrateDecisionCriteria()

	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║  Key Insights                                             ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════╣")
	fmt.Println("║  1. Value receivers COPY the struct (no modification)    ║")
	fmt.Println("║  2. Pointer receivers modify the ORIGINAL                ║")
	fmt.Println("║  3. Only *T satisfies interfaces with pointer receivers  ║")
	fmt.Println("║  4. Large structs: use pointers (avoid copy overhead)    ║")
	fmt.Println("║  5. Consistency: if one method uses *, all should        ║")
	fmt.Println("║  6. Nil receivers are valid (check before dereferencing) ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
}
