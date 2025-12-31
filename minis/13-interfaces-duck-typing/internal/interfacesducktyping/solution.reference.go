//go:build reference
// +build reference

/*
Problem: Understanding Go's interface system and duck typing

Requirements:
1. Implement interfaces implicitly (no "implements" keyword)
2. Use type assertions to extract concrete types
3. Handle nil interface gotchas (type vs value)
4. Compose interfaces through embedding
5. Implement polymorphism with interface dispatch

Data Structure:
- Interface value: Type pointer + Data pointer (16 bytes on 64-bit)
- Type assertion: Runtime check of concrete type
- Type switch: Multi-way type-based branching

Algorithm: Dynamic Dispatch
- Interface stores concrete type metadata
- Method calls routed through virtual table
- Type assertions inspect type metadata

Why interfaces enable polymorphism:
- Decouple behavior from implementation
- Write code once, works with many types
- No inheritance needed
- Runtime flexibility with compile-time safety
*/

package interfacesducktyping

import "fmt"

// String implements the Stringer interface for Person.
// BREAKPOINT: Set breakpoint here to trace interface implementation
// DEBUG: Watch 'p' to see Person value
// DEBUG: Watch return value formatting
func (p Person) String() string {
	// BREAKPOINT: Set breakpoint here before formatting
	// DEBUG: Watch 'p.Name' and 'p.Age' fields
	// DEBUG: Watch fmt.Sprintf create formatted string
	return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
}

// GetAge extracts the age from a Stringer if it's a Person.
// BREAKPOINT: Set breakpoint here to trace type assertion
// DEBUG: Watch 's' interface value
// DEBUG: Watch 'ok' to see if assertion succeeds
func GetAge(s Stringer) (int, bool) {
	// BREAKPOINT: Set breakpoint here for type assertion
	// DEBUG: Watch 'p, ok := s.(Person)' syntax
	// DEBUG: Watch 'ok' = true if s contains Person
	// DEBUG: Watch 'ok' = false if s contains different type
	p, ok := s.(Person)
	if !ok {
		// DEBUG: Not a Person - return zero values
		return 0, false
	}

	// BREAKPOINT: Set breakpoint here when assertion succeeds
	// DEBUG: Watch 'p.Age' from extracted Person
	return p.Age, true
}

// DescribeType returns a description of the type of the value.
// BREAKPOINT: Set breakpoint here to trace type switching
// DEBUG: Watch 'i' interface value
// DEBUG: Watch type determination in switch
func DescribeType(i interface{}) string {
	// BREAKPOINT: Set breakpoint here before type switch
	// DEBUG: Watch 'v := i.(type)' syntax
	switch v := i.(type) {
	case int:
		// BREAKPOINT: Hit when i contains int
		// DEBUG: Watch 'v' has type int here
		return fmt.Sprintf("Integer: %d", v)

	case string:
		// BREAKPOINT: Hit when i contains string
		// DEBUG: Watch 'v' has type string here
		return fmt.Sprintf("String: %s", v)

	case bool:
		// BREAKPOINT: Hit when i contains bool
		// DEBUG: Watch 'v' has type bool here
		return fmt.Sprintf("Boolean: %t", v)

	case Person:
		// BREAKPOINT: Hit when i contains Person
		// DEBUG: Watch 'v' has type Person here
		// DEBUG: Watch 'v.Name' field access
		return fmt.Sprintf("Person: %s", v.Name)

	case nil:
		// BREAKPOINT: Hit when i is nil interface
		// DEBUG: Both type and value are nil
		return "Nil"

	default:
		// BREAKPOINT: Hit for unhandled types
		// DEBUG: Watch unknown type fall through
		return "Unknown"
	}
}

// IsValidEmail checks if a Validator is valid, handling nil correctly.
// BREAKPOINT: Set breakpoint here to trace nil handling
// DEBUG: Watch 'v' interface value
// DEBUG: Demonstrate the "nil interface vs nil pointer" gotcha
func IsValidEmail(v Validator) bool {
	// BREAKPOINT: Set breakpoint here for first nil check
	// DEBUG: Watch 'v == nil' checks if BOTH type and value are nil
	if v == nil {
		// DEBUG: True nil interface - no type, no value
		return false
	}

	// At this point, v != nil BUT it might contain a nil pointer!
	// Example: var e *Email = nil; var v Validator = e
	// Here v != nil (type is *Email) but the pointer is nil

	// BREAKPOINT: Set breakpoint here for type assertion
	// DEBUG: Watch 'e, ok := v.(*Email)' extract concrete type
	e, ok := v.(*Email)
	if !ok {
		// BREAKPOINT: Not an *Email
		// DEBUG: Different Validator implementation
		return v.IsValid()
	}

	// BREAKPOINT: Set breakpoint here for pointer nil check
	// DEBUG: Watch 'e == nil' checks if the *Email pointer is nil
	if e == nil {
		// DEBUG: Interface contains nil pointer!
		return false
	}

	// BREAKPOINT: Set breakpoint here when Email is valid
	// DEBUG: Watch e.IsValid() call on non-nil Email
	return e.IsValid()
}

// Read returns the current data in the buffer.
// BREAKPOINT: Set breakpoint here to trace Reader implementation
// DEBUG: Watch 'b' Buffer pointer
// DEBUG: Watch 'b.data' field access
func (b *Buffer) Read() string {
	// DEBUG: Simple accessor - returns current data
	return b.data
}

// Write appends data to the buffer.
// BREAKPOINT: Set breakpoint here to trace Writer implementation
// DEBUG: Watch 'b.data' before append
// DEBUG: Watch 'data' parameter being appended
// DEBUG: Watch 'b.data' after append
func (b *Buffer) Write(data string) error {
	// BREAKPOINT: Set breakpoint here before append
	// DEBUG: Watch string concatenation
	b.data += data
	// DEBUG: In-memory buffer never fails
	return nil
}

// IsReadWriter checks if an interface value implements ReadWriter.
// BREAKPOINT: Set breakpoint here to trace interface composition
// DEBUG: Watch 'i' interface value
// DEBUG: Watch composite interface checking
func IsReadWriter(i interface{}) bool {
	// BREAKPOINT: Set breakpoint here for type assertion
	// DEBUG: Watch '_, ok := i.(ReadWriter)' check both Read and Write
	// DEBUG: ReadWriter requires BOTH Reader and Writer methods
	_, ok := i.(ReadWriter)
	return ok
}

// Increment increments the counter value.
// BREAKPOINT: Set breakpoint here to trace pointer receiver method
// DEBUG: Watch 'c' pointer to Counter
// DEBUG: Watch 'c.Value' before and after increment
func (c *Counter) Increment() {
	// BREAKPOINT: Set breakpoint here before increment
	// DEBUG: Watch 'c.Value' increase by 1
	c.Value++
	// DEBUG: Pointer receiver required for mutation
}

// CanIncrement checks if a value can be used as an Incrementer.
// BREAKPOINT: Set breakpoint here to trace interface checking
// DEBUG: Watch 'i' interface value
// DEBUG: Only *Counter satisfies Incrementer (not Counter)
func CanIncrement(i interface{}) bool {
	// BREAKPOINT: Set breakpoint here for type assertion
	// DEBUG: Watch '_, ok := i.(Incrementer)'
	// DEBUG: Remember: pointer receiver methods only on *T, not T
	_, ok := i.(Incrementer)
	return ok
}

// CountTypes counts how many values of each type are in the slice.
// BREAKPOINT: Set breakpoint here to trace type counting
// DEBUG: Watch 'values' slice of interface{}
// DEBUG: Watch 'counts' map build up
func CountTypes(values []interface{}) map[string]int {
	// BREAKPOINT: Set breakpoint here to create map
	// DEBUG: Watch map initialization
	counts := make(map[string]int)

	// BREAKPOINT: Set breakpoint here before loop
	// DEBUG: Watch iteration through values
	for _, v := range values {
		// BREAKPOINT: Set breakpoint here for each value
		// DEBUG: Watch 'v' current value
		// DEBUG: Watch '%T' format verb extract type name
		typeName := fmt.Sprintf("%T", v)

		// BREAKPOINT: Set breakpoint here for increment
		// DEBUG: Watch 'counts[typeName]++' increment count
		// DEBUG: Zero value (0) allows safe increment without check
		counts[typeName]++
	}

	// DEBUG: Watch final 'counts' map
	return counts
}

// Error implements the error interface for ValidationError.
// BREAKPOINT: Set breakpoint here to trace error formatting
// DEBUG: Watch 'e' ValidationError value
// DEBUG: Watch error message construction
func (e ValidationError) Error() string {
	// BREAKPOINT: Set breakpoint here before formatting
	// DEBUG: Watch 'e.Field' and 'e.Message' fields
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// Area calculates the area of a rectangle.
// BREAKPOINT: Set breakpoint here to trace Rectangle.Area
// DEBUG: Watch 'r' Rectangle value
// DEBUG: Watch 'r.Width' and 'r.Height' fields
func (r Rectangle) Area() float64 {
	// BREAKPOINT: Set breakpoint here before calculation
	// DEBUG: Watch width * height calculation
	return r.Width * r.Height
}

// Area calculates the area of a circle.
// BREAKPOINT: Set breakpoint here to trace Circle.Area
// DEBUG: Watch 'c' Circle pointer
// DEBUG: Watch 'c.Radius' field
func (c Circle) Area() float64 {
	// BREAKPOINT: Set breakpoint here before calculation
	// DEBUG: Watch π * r² calculation
	// DEBUG: Using approximate π value
	return 3.14159 * c.Radius * c.Radius
}

// TotalArea calculates the total area of all shapes.
// BREAKPOINT: Set breakpoint here to trace polymorphism
// DEBUG: Watch 'shapes' slice of Shape interface
// DEBUG: Watch dynamic dispatch to correct Area() method
func TotalArea(shapes []Shape) float64 {
	// BREAKPOINT: Set breakpoint here to initialize total
	// DEBUG: Watch 'total' accumulator
	total := 0.0

	// BREAKPOINT: Set breakpoint here before loop
	// DEBUG: Watch iteration through shapes
	for _, shape := range shapes {
		// BREAKPOINT: Set breakpoint here for each shape
		// DEBUG: Watch 'shape' (could be Rectangle or Circle)
		// DEBUG: Watch 'shape.Area()' dynamically dispatch
		// DEBUG: Runtime determines Rectangle.Area or Circle.Area
		total += shape.Area()
	}

	// DEBUG: Watch final 'total' (sum of all areas)
	return total
}
