// Package main demonstrates Go interfaces, duck typing, and method sets.
//
// MACRO-COMMENT: What This Program Teaches
// =========================================
// This program provides hands-on examples of:
// 1. Interface basics and implicit implementation (duck typing)
// 2. Type assertions and type switches (runtime type checking)
// 3. Empty interface (interface{} and any)
// 4. Nil interface gotchas (the two-part nil problem)
// 5. Method sets and interface satisfaction rules
// 6. Polymorphism without inheritance
// 7. Real-world interface patterns (io.Writer, error, Stringer)
//
// COMPILER BEHAVIOR: Interface Values
// An interface value is a two-word struct containing:
// 1. A pointer to type information (the "itab" - interface table)
// 2. A pointer to the actual data
// This allows runtime polymorphism while maintaining type safety.
//
// RUNTIME BEHAVIOR: Interface Calls
// When you call a method on an interface:
// 1. The runtime looks up the method in the interface table (vtable)
// 2. Calls the method with the data pointer as the receiver
// This adds ~1-2ns overhead vs direct calls (negligible in most cases).

package main

import (
	"fmt"
)

// ============================================================================
// SECTION 1: Interface Basics and Implicit Implementation
// ============================================================================

// Writer defines the contract for anything that can write bytes.
//
// MACRO-COMMENT: Interface Design Principles
// A good interface:
// - Is small (1-3 methods is ideal)
// - Describes behavior, not data
// - Has a clear, single responsibility
// - Named after what it does (Writer writes, Reader reads, etc.)
//
// The Writer interface is modeled after io.Writer from the standard library.
type Writer interface {
	Write(data string) error
}

// File represents a file on disk.
//
// MICRO-COMMENT: This is a concrete type (a struct).
// It doesn't "declare" that it implements Writer.
// It just happens to have a Write method with the right signature.
type File struct {
	name string
}

// Write writes data to the file.
//
// MACRO-COMMENT: Implicit Interface Implementation
// This method makes File satisfy the Writer interface automatically.
// There's no "implements" keyword - Go uses duck typing:
// "If it has a Write(string) error method, it's a Writer."
//
// MEMORY LAYOUT when assigned to interface:
// var w Writer = &File{name: "log"}
// ┌─────────────────┬──────────────┐
// │ type: *File     │ value: 0x... │
// └─────────────────┴──────────────┘
//                        │
//                        └──> File{name: "log"}
func (f *File) Write(data string) error {
	// MICRO-COMMENT: In a real implementation, this would write to disk
	fmt.Printf("File '%s' writes: %s\n", f.name, data)
	return nil
}

// Network represents a network connection.
type Network struct {
	address string
}

// Write sends data over the network.
//
// MICRO-COMMENT: Network also implements Writer (implicitly).
// File and Network are completely unrelated types,
// but they're both Writers because they have the right method.
func (n *Network) Write(data string) error {
	// MICRO-COMMENT: In a real implementation, this would send over TCP/UDP
	fmt.Printf("Network '%s' writes: %s\n", n.address, data)
	return nil
}

// WriteData writes data using any Writer.
//
// MACRO-COMMENT: Polymorphism Without Inheritance
// This function works with ANY type that implements Writer.
// It doesn't care if it's a File, Network, or something that doesn't exist yet.
// This is the power of interfaces: decoupling from concrete types.
//
// COMPILER BEHAVIOR:
// When you pass &File{} to this function, the compiler:
// 1. Checks that *File has a Write(string) error method ✓
// 2. Creates an interface value: (type=*File, value=pointer to File)
// 3. Passes it to WriteData
func WriteData(w Writer, data string) {
	// MICRO-COMMENT: This is an interface method call
	// At runtime, Go looks up the Write method in the interface table
	// and calls the appropriate implementation (*File.Write or *Network.Write)
	err := w.Write(data)
	if err != nil {
		fmt.Printf("Error writing: %v\n", err)
	}
}

// demonstrateInterfaceBasics shows implicit interface implementation.
func demonstrateInterfaceBasics() {
	fmt.Println("=== Interface Basics: Implicit Implementation ===")

	// MICRO-COMMENT: Create concrete types
	file := &File{name: "log.txt"}
	network := &Network{address: "192.168.1.1:8080"}

	// MACRO-COMMENT: Polymorphism in Action
	// WriteData accepts a Writer interface.
	// We can pass File or Network (or any future type with a Write method).
	// This is duck typing: "If it can Write, it's a Writer."
	WriteData(file, "Hello from file!")
	WriteData(network, "Hello from network!")

	// MICRO-COMMENT: We can also assign to interface variables
	var w Writer

	// MICRO-COMMENT: Assigning *File to Writer
	// The compiler creates an interface value: (type=*File, value=file)
	w = file
	w.Write("Using interface variable with File")

	// MICRO-COMMENT: Reassigning the same interface variable to *Network
	// The interface value is updated: (type=*Network, value=network)
	w = network
	w.Write("Using interface variable with Network")

	fmt.Println()
}

// ============================================================================
// SECTION 2: Type Assertions and Type Switches
// ============================================================================

// GetFileName extracts the filename if the Writer is a File.
//
// MACRO-COMMENT: Type Assertions
// Sometimes you need to access the concrete type inside an interface.
// Type assertions let you extract the underlying value.
//
// TWO FORMS:
// 1. Unsafe: f := w.(*File)  // Panics if w is not *File
// 2. Safe:   f, ok := w.(*File)  // Returns ok=false if w is not *File
//
// ALWAYS use the safe form in production code!
func GetFileName(w Writer) string {
	// MICRO-COMMENT: Type assertion with boolean check (safe)
	// This checks: "Is the concrete type inside w a *File?"
	// If yes: f gets the *File, ok is true
	// If no: f gets nil, ok is false
	f, ok := w.(*File)
	if !ok {
		return "(not a file)"
	}

	// MICRO-COMMENT: Now we can access File-specific fields
	return f.name
}

// demonstrateTypeAssertions shows how to extract concrete types.
func demonstrateTypeAssertions() {
	fmt.Println("=== Type Assertions: Extracting Concrete Types ===")

	// MICRO-COMMENT: Create a Writer that's actually a *File
	var w Writer = &File{name: "data.txt"}

	// MICRO-COMMENT: Try to extract the *File
	filename := GetFileName(w)
	fmt.Printf("File name: %s\n", filename)

	// MICRO-COMMENT: Now assign a *Network (not a *File)
	w = &Network{address: "10.0.0.1:9000"}
	filename = GetFileName(w)
	fmt.Printf("File name: %s\n", filename)  // Will print "(not a file)"

	// MACRO-COMMENT: Unsafe Type Assertion (Demo Only!)
	// This demonstrates the panic behavior when the assertion fails.
	// NEVER do this in production without checking first!
	defer func() {
		// MICRO-COMMENT: Recover from the panic to continue the demo
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	// MICRO-COMMENT: This will panic because w is *Network, not *File
	// f := w.(*File)  // PANIC: interface holds *Network, not *File
	// fmt.Println(f.name)

	fmt.Println()
}

// DescribeValue uses a type switch to handle different types.
//
// MACRO-COMMENT: Type Switches
// When you need to check for multiple types, use a type switch.
// It's like a switch statement, but switches on the type, not a value.
//
// SYNTAX: switch v := i.(type) { ... }
// The .(type) syntax is ONLY valid in switch statements.
//
// PERFORMANCE: Type switches are fast (single comparison + jump table).
func DescribeValue(i interface{}) {
	// MICRO-COMMENT: Type switch syntax
	// v gets the concrete value with the appropriate type in each case
	switch v := i.(type) {
	case int:
		// MICRO-COMMENT: v has type int here
		fmt.Printf("Integer: %d (doubled: %d)\n", v, v*2)

	case string:
		// MICRO-COMMENT: v has type string here
		fmt.Printf("String: %q (length: %d)\n", v, len(v))

	case bool:
		// MICRO-COMMENT: v has type bool here
		fmt.Printf("Boolean: %t\n", v)

	case *File:
		// MICRO-COMMENT: v has type *File here
		// We can access File-specific fields
		fmt.Printf("File: %s\n", v.name)

	case *Network:
		// MICRO-COMMENT: v has type *Network here
		fmt.Printf("Network: %s\n", v.address)

	case nil:
		// MICRO-COMMENT: Special case for nil
		fmt.Println("Nil value")

	default:
		// MICRO-COMMENT: Fallback for any other type
		// %T prints the type, %v prints the value
		fmt.Printf("Unknown type: %T, value: %v\n", v, v)
	}
}

// demonstrateTypeSwitches shows type switching on interface values.
func demonstrateTypeSwitches() {
	fmt.Println("=== Type Switches: Pattern Matching on Types ===")

	// MICRO-COMMENT: Test the type switch with different values
	DescribeValue(42)
	DescribeValue("hello")
	DescribeValue(true)
	DescribeValue(&File{name: "config.json"})
	DescribeValue(&Network{address: "localhost:8080"})
	DescribeValue([]int{1, 2, 3})  // Unknown type (no case for []int)
	DescribeValue(nil)

	fmt.Println()
}

// ============================================================================
// SECTION 3: Empty Interface (interface{} and any)
// ============================================================================

// PrintAnything accepts any type using the empty interface.
//
// MACRO-COMMENT: The Empty Interface
// interface{} (or 'any' in Go 1.18+) is an interface with ZERO methods.
// Every type implements it (every type has at least zero methods!).
//
// USE CASES:
// - Containers that hold multiple types (maps, slices of mixed data)
// - Generic functions before Go 1.18 generics
// - JSON/XML unmarshaling (unknown structure)
//
// DOWNSIDE: You lose type safety and need runtime type checks.
//
// MEMORY LAYOUT:
// The interface still stores (type, value), so you can retrieve the type later.
func PrintAnything(value interface{}) {
	// MICRO-COMMENT: The %T verb prints the dynamic type
	// The %v verb prints the value
	fmt.Printf("Type: %T, Value: %v\n", value, value)
}

// demonstrateEmptyInterface shows the universal container.
func demonstrateEmptyInterface() {
	fmt.Println("=== Empty Interface: The Universal Type ===")

	// MICRO-COMMENT: Create a slice of interface{} (can hold ANY types)
	var items []interface{}

	// MICRO-COMMENT: Add different types to the same slice
	items = append(items, 42)
	items = append(items, "hello")
	items = append(items, true)
	items = append(items, &File{name: "mixed.txt"})
	items = append(items, []int{1, 2, 3})

	// MACRO-COMMENT: Iterating Over Mixed Types
	// We can store anything, but we need type assertions to use them.
	for _, item := range items {
		PrintAnything(item)
	}

	// MICRO-COMMENT: Go 1.18+ prefers 'any' over 'interface{}'
	// They're identical, but 'any' is more readable
	var x any = 100
	PrintAnything(x)

	fmt.Println()
}

// ============================================================================
// SECTION 4: Nil Interface Gotchas
// ============================================================================

// IsNil checks if an interface is nil.
//
// MACRO-COMMENT: The Two-Part Nil Problem
// An interface value has TWO parts: (type, value)
// An interface is nil ONLY if BOTH parts are nil.
//
// COMMON BUG:
//   var p *File = nil      // p is nil
//   var w Writer = p       // w is NOT nil! (type=*File, value=nil)
//   if w == nil { ... }    // NEVER executes (w is not nil)
//
// This is the most common interface gotcha in Go.
func IsNil(w Writer) bool {
	// MICRO-COMMENT: This checks if the interface itself is nil
	// (i.e., both type and value are nil)
	return w == nil
}

// demonstrateNilInterface shows the nil interface gotcha.
func demonstrateNilInterface() {
	fmt.Println("=== Nil Interface Gotcha: The Two-Part Nil ===")

	// MICRO-COMMENT: Create a nil pointer
	var filePtr *File = nil
	fmt.Printf("filePtr == nil: %t\n", filePtr == nil)  // true

	// MICRO-COMMENT: Assign the nil pointer to an interface
	// This creates an interface with (type=*File, value=nil)
	var w Writer = filePtr

	// MACRO-COMMENT: The Gotcha!
	// Even though filePtr is nil, the interface w is NOT nil.
	// Why? Because the type part is set (*File).
	// An interface is nil only if BOTH type and value are nil.
	fmt.Printf("w == nil: %t ← GOTCHA!\n", w == nil)  // false

	// MICRO-COMMENT: Call our IsNil function
	fmt.Printf("IsNil(w): %t\n", IsNil(w))  // false

	// MICRO-COMMENT: Create a truly nil interface
	var w2 Writer = nil  // Both type and value are nil
	fmt.Printf("w2 == nil: %t\n", w2 == nil)  // true
	fmt.Printf("IsNil(w2): %t\n", IsNil(w2))  // true

	// MACRO-COMMENT: Real-World Impact
	// This bug often appears when returning errors:
	//   func GetUser() *User { return nil }
	//   func Wrap() error { return GetUser() }  // BUG!
	//   if Wrap() != nil { ... }  // Always true!
	//
	// The fix: return typed nil directly
	//   func Wrap() error {
	//       u := GetUser()
	//       if u == nil { return nil }
	//       return u
	//   }

	fmt.Println()
}

// ============================================================================
// SECTION 5: Method Sets and Interface Satisfaction
// ============================================================================

// Incrementer defines a contract for types that can increment.
type Incrementer interface {
	Increment()
}

// Counter with a value receiver method.
type Counter struct {
	count int
}

// IncrementValue has a VALUE receiver.
//
// MICRO-COMMENT: Method Set Rule
// A method with a value receiver (c Counter) is in the method set of:
// - Counter (the value type)
// - *Counter (the pointer type)
//
// RESULT: Both Counter and *Counter satisfy interfaces with this method.
func (c Counter) IncrementValue() {
	c.count++  // Modifies a copy (doesn't change the original)
	fmt.Printf("  IncrementValue (value receiver): %d\n", c.count)
}

// Increment has a POINTER receiver.
//
// MICRO-COMMENT: Method Set Rule
// A method with a pointer receiver (c *Counter) is ONLY in the method set of:
// - *Counter (the pointer type)
//
// RESULT: Only *Counter satisfies interfaces with this method, NOT Counter.
//
// WHY? Go can't always take the address of a value (e.g., map values,
// interface values). So it can't auto-promote Counter to *Counter.
func (c *Counter) Increment() {
	c.count++  // Modifies the original
	fmt.Printf("  Increment (pointer receiver): %d\n", c.count)
}

// demonstrateMethodSets shows interface satisfaction rules.
func demonstrateMethodSets() {
	fmt.Println("=== Method Sets: Pointer vs Value Receivers ===")

	// MICRO-COMMENT: Create a Counter value and pointer
	c := Counter{count: 0}
	p := &Counter{count: 100}

	// MACRO-COMMENT: Direct Method Calls (Always Work)
	// When you call a method on a variable, Go auto-addresses:
	// c.Increment() → (&c).Increment()
	fmt.Println("Direct method calls:")
	c.IncrementValue()  // Works (value receiver on value)
	c.Increment()       // Works (Go does &c automatically)
	p.IncrementValue()  // Works (Go dereferences *p automatically)
	p.Increment()       // Works (pointer receiver on pointer)

	// MACRO-COMMENT: Interface Assignment (Strict Rules)
	// When assigning to an interface, Go does NOT auto-address.
	// You must use the correct type.

	var inc Incrementer

	// MICRO-COMMENT: This would FAIL (Counter doesn't implement Incrementer)
	// inc = Counter{count: 0}  // COMPILE ERROR
	// Why? Counter only has methods with pointer receivers in Incrementer.

	// MICRO-COMMENT: This works (*Counter implements Incrementer)
	inc = &Counter{count: 200}
	inc.Increment()

	// MACRO-COMMENT: The Rule of Thumb
	// - If you need to modify the receiver, use pointer receiver
	// - Pointer receiver methods mean only *T implements the interface
	// - Value receiver methods mean both T and *T implement the interface

	fmt.Println()
}

// ============================================================================
// SECTION 6: Polymorphism Example
// ============================================================================

// Shape is an interface for geometric shapes.
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle implements Shape.
type Circle struct {
	Radius float64
}

// Area calculates the area of a circle.
func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

// Perimeter calculates the circumference of a circle.
func (c Circle) Perimeter() float64 {
	return 2 * 3.14159 * c.Radius
}

// Rectangle implements Shape.
type Rectangle struct {
	Width, Height float64
}

// Area calculates the area of a rectangle.
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates the perimeter of a rectangle.
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// PrintShapeInfo prints information about any Shape.
//
// MACRO-COMMENT: Polymorphism Without Inheritance
// This function works with any type that implements Shape.
// No inheritance hierarchy needed - just the right methods.
func PrintShapeInfo(s Shape) {
	fmt.Printf("  Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

// demonstratePolymorphism shows polymorphism with interfaces.
func demonstratePolymorphism() {
	fmt.Println("=== Polymorphism: Multiple Types, One Interface ===")

	// MICRO-COMMENT: Create different shapes
	circle := Circle{Radius: 5.0}
	rectangle := Rectangle{Width: 4.0, Height: 6.0}

	// MICRO-COMMENT: PrintShapeInfo works with any Shape
	fmt.Println("Circle:")
	PrintShapeInfo(circle)

	fmt.Println("Rectangle:")
	PrintShapeInfo(rectangle)

	// MICRO-COMMENT: We can also store them in a slice
	shapes := []Shape{
		Circle{Radius: 3.0},
		Rectangle{Width: 2.0, Height: 5.0},
		Circle{Radius: 7.0},
	}

	fmt.Println("\nAll shapes:")
	for i, shape := range shapes {
		fmt.Printf("Shape %d: ", i+1)
		PrintShapeInfo(shape)
	}

	fmt.Println()
}

// ============================================================================
// SECTION 7: Interface Composition
// ============================================================================

// Reader defines the contract for reading.
type Reader interface {
	Read() string
}

// Closer defines the contract for closing resources.
type Closer interface {
	Close() error
}

// ReadCloser combines Reader and Closer.
//
// MACRO-COMMENT: Interface Composition
// Interfaces can embed other interfaces.
// ReadCloser requires BOTH Read() and Close() methods.
//
// This is how the standard library builds complex interfaces:
//   io.ReadWriter embeds io.Reader and io.Writer
//   io.ReadWriteCloser embeds io.ReadWriter and io.Closer
type ReadCloser interface {
	Reader
	Closer
}

// FileReader implements ReadCloser.
type FileReader struct {
	filename string
	closed   bool
}

// Read reads data from the file.
func (f *FileReader) Read() string {
	if f.closed {
		return "(file closed)"
	}
	return fmt.Sprintf("Reading from %s", f.filename)
}

// Close closes the file.
func (f *FileReader) Close() error {
	if f.closed {
		return fmt.Errorf("already closed")
	}
	f.closed = true
	fmt.Printf("Closed %s\n", f.filename)
	return nil
}

// demonstrateInterfaceComposition shows embedded interfaces.
func demonstrateInterfaceComposition() {
	fmt.Println("=== Interface Composition: Building Larger Contracts ===")

	// MICRO-COMMENT: FileReader implements ReadCloser (has both Read and Close)
	var rc ReadCloser = &FileReader{filename: "data.txt"}

	// MICRO-COMMENT: Use as Reader
	data := rc.Read()
	fmt.Println(data)

	// MICRO-COMMENT: Use as Closer
	err := rc.Close()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// MICRO-COMMENT: Try reading after close
	data = rc.Read()
	fmt.Println(data)

	fmt.Println()
}

// ============================================================================
// MAIN FUNCTION: Orchestrates All Demonstrations
// ============================================================================

// main executes all demonstration functions in order.
//
// MACRO-COMMENT: Learning Progression
// The demos are ordered to build understanding:
// 1. Interface basics (foundation)
// 2. Type assertions (extracting concrete types)
// 3. Type switches (pattern matching)
// 4. Empty interface (universal container)
// 5. Nil gotcha (common bug)
// 6. Method sets (satisfaction rules)
// 7. Polymorphism (real-world usage)
// 8. Composition (building complex interfaces)
//
// PERFORMANCE NOTE:
// Run with benchmarks to see interface call overhead:
//   go test -bench=. -benchmem
//
// ESCAPE ANALYSIS:
// Run with -gcflags='-m' to see heap allocations:
//   go build -gcflags='-m' cmd/interfaces-demo/main.go
func main() {
	demonstrateInterfaceBasics()
	demonstrateTypeAssertions()
	demonstrateTypeSwitches()
	demonstrateEmptyInterface()
	demonstrateNilInterface()
	demonstrateMethodSets()
	demonstratePolymorphism()
	demonstrateInterfaceComposition()

	// MACRO-COMMENT: Key Insights
	// After running this program, you should understand:
	// 1. Interfaces enable polymorphism without inheritance
	// 2. Duck typing makes code flexible and decoupled
	// 3. Type assertions/switches let you work with concrete types when needed
	// 4. Nil interfaces are tricky (two-part nil)
	// 5. Method sets determine which types satisfy which interfaces
	// 6. Interfaces are the foundation of Go's standard library patterns
}
