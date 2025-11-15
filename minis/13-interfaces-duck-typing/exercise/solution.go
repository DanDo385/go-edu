//go:build solution
// +build solution

package exercise

import "fmt"

// ============================================================================
// SOLUTION 1: Implementing Interfaces
// ============================================================================

// String implements the Stringer interface for Person.
//
// MICRO-COMMENT: This method makes Person satisfy the Stringer interface.
// No explicit declaration needed - just implement the method with the right signature.
func (p Person) String() string {
	// MICRO-COMMENT: Format the string as specified
	// %s is for string, %d is for integer
	return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
}

// ============================================================================
// SOLUTION 2: Type Assertions
// ============================================================================

// GetAge extracts the age from a Stringer if it's a Person.
func GetAge(s Stringer) (int, bool) {
	// MICRO-COMMENT: Type assertion with comma-ok idiom
	// This checks if the concrete type inside s is Person
	// If yes: p gets the Person value, ok is true
	// If no: p gets zero value, ok is false
	p, ok := s.(Person)
	if !ok {
		return 0, false
	}

	// MICRO-COMMENT: Return the age from the Person
	return p.Age, true
}

// ============================================================================
// SOLUTION 3: Type Switches
// ============================================================================

// DescribeType returns a description of the type of the value.
func DescribeType(i interface{}) string {
	// MICRO-COMMENT: Type switch - switches on the type, not the value
	// v gets the value with the appropriate type in each case
	switch v := i.(type) {
	case int:
		// MICRO-COMMENT: v has type int here
		return fmt.Sprintf("Integer: %d", v)

	case string:
		// MICRO-COMMENT: v has type string here
		return fmt.Sprintf("String: %s", v)

	case bool:
		// MICRO-COMMENT: v has type bool here
		return fmt.Sprintf("Boolean: %t", v)

	case Person:
		// MICRO-COMMENT: v has type Person here
		// We can access Person-specific fields
		return fmt.Sprintf("Person: %s", v.Name)

	case nil:
		// MICRO-COMMENT: Special case for nil values
		return "Nil"

	default:
		// MICRO-COMMENT: Catch-all for any other type
		return "Unknown"
	}
}

// ============================================================================
// SOLUTION 4: Interface Nil Check
// ============================================================================

// IsValidEmail checks if a Validator is valid, handling nil correctly.
func IsValidEmail(v Validator) bool {
	// MICRO-COMMENT: First check - is the interface itself nil?
	// This checks if both type and value are nil
	if v == nil {
		return false
	}

	// MACRO-COMMENT: The Two-Part Nil Problem
	// At this point, we know v is not a nil interface.
	// But it could still contain a nil pointer!
	// Example: var e *Email = nil; var v Validator = e
	// In this case, v != nil (interface is not nil)
	// but e is nil (the pointer inside is nil)

	// MICRO-COMMENT: Type assert to *Email to check if the pointer is nil
	e, ok := v.(*Email)
	if !ok {
		// Not an *Email, but could be some other Validator
		// Call IsValid() and trust the implementation
		return v.IsValid()
	}

	// MICRO-COMMENT: Check if the Email pointer itself is nil
	if e == nil {
		return false
	}

	// MICRO-COMMENT: Finally, call IsValid() on the non-nil Email
	return e.IsValid()
}

// ============================================================================
// SOLUTION 5: Implementing Multiple Interfaces
// ============================================================================

// Read returns the current data in the buffer.
func (b *Buffer) Read() string {
	// MICRO-COMMENT: Just return the data
	// This makes Buffer implement the Reader interface
	return b.data
}

// Write appends data to the buffer.
func (b *Buffer) Write(data string) error {
	// MICRO-COMMENT: Append to existing data
	// This makes Buffer implement the Writer interface
	b.data += data

	// MICRO-COMMENT: In-memory buffers don't have errors
	return nil
}

// ============================================================================
// SOLUTION 6: Interface Composition
// ============================================================================

// IsReadWriter checks if an interface value implements ReadWriter.
func IsReadWriter(i interface{}) bool {
	// MICRO-COMMENT: Type assert to ReadWriter
	// ReadWriter is an interface that embeds both Reader and Writer
	// So this checks if i has both Read() and Write() methods
	_, ok := i.(ReadWriter)
	return ok
}

// ============================================================================
// SOLUTION 7: Method Sets and Receivers
// ============================================================================

// Increment increments the counter value.
//
// MACRO-COMMENT: Pointer Receiver for Mutation
// We use a pointer receiver (*Counter) because we need to modify the counter.
// This also means that only *Counter implements Incrementer, not Counter.
func (c *Counter) Increment() {
	c.Value++
}

// CanIncrement checks if a value can be used as an Incrementer.
func CanIncrement(i interface{}) bool {
	// MICRO-COMMENT: Type assert to Incrementer interface
	// This will return true only if i has an Increment() method
	// with the right signature
	_, ok := i.(Incrementer)
	return ok
}

// ============================================================================
// SOLUTION 8: Working with Empty Interface
// ============================================================================

// CountTypes counts how many values of each type are in the slice.
func CountTypes(values []interface{}) map[string]int {
	// MICRO-COMMENT: Create a map to store counts
	counts := make(map[string]int)

	// MICRO-COMMENT: Iterate through all values
	for _, v := range values {
		// MICRO-COMMENT: Get the type name using %T format verb
		// %T prints the type of the value
		typeName := fmt.Sprintf("%T", v)

		// MICRO-COMMENT: Increment the count for this type
		counts[typeName]++
	}

	return counts
}

// ============================================================================
// SOLUTION 9: Error Interface
// ============================================================================

// Error implements the error interface for ValidationError.
//
// MACRO-COMMENT: The error Interface
// The error interface is defined as:
//   type error interface {
//       Error() string
//   }
// Any type with an Error() string method automatically implements error.
func (e ValidationError) Error() string {
	// MICRO-COMMENT: Format the error message as specified
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// ============================================================================
// SOLUTION 10: Polymorphism
// ============================================================================

// Area calculates the area of a rectangle.
//
// MICRO-COMMENT: This makes Rectangle implement the Shape interface.
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Area calculates the area of a circle.
//
// MICRO-COMMENT: This makes Circle implement the Shape interface.
func (c Circle) Area() float64 {
	// MICRO-COMMENT: π * r²
	return 3.14159 * c.Radius * c.Radius
}

// TotalArea calculates the total area of all shapes.
//
// MACRO-COMMENT: Polymorphism in Action
// This function works with any type that implements Shape.
// It doesn't care if it's a Rectangle, Circle, or some future shape type.
// This is the power of interfaces: write code once, works with many types.
func TotalArea(shapes []Shape) float64 {
	// MICRO-COMMENT: Start with zero total
	total := 0.0

	// MICRO-COMMENT: Sum the areas
	// Each shape.Area() call is dynamically dispatched to the
	// appropriate implementation (Rectangle.Area or Circle.Area)
	for _, shape := range shapes {
		total += shape.Area()
	}

	return total
}
