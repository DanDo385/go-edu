//go:build !solution
// +build !solution

// Package exercise contains hands-on exercises for understanding Go interfaces.
//
// LEARNING OBJECTIVES:
// - Implement interfaces implicitly (duck typing)
// - Use type assertions and type switches
// - Work with empty interfaces
// - Understand nil interface behavior
// - Apply method set rules correctly

package exercise

// TODO: Implement these functions and methods according to the specifications in the tests.
// Each function tests a different aspect of interface mechanics.

// ============================================================================
// EXERCISE 1: Implementing Interfaces
// ============================================================================

// TODO: Implement the String() method for Person to satisfy the Stringer interface.
//
// REQUIREMENTS:
// - Return a string in the format: "Name (Age years old)"
// - Example: Person{Name: "Alice", Age: 30} → "Alice (30 years old)"
//
// HINT: Use fmt.Sprintf to format the string.
func (p Person) String() string {
	// TODO: Implement this method
	return ""
}

// ============================================================================
// EXERCISE 2: Type Assertions
// ============================================================================

// GetAge extracts the age from a Stringer if it's a Person.
//
// REQUIREMENTS:
// - If s is a Person, return the person's age and true
// - If s is NOT a Person, return 0 and false
//
// HINT: Use type assertion with the comma-ok idiom: p, ok := s.(Person)
func GetAge(s Stringer) (int, bool) {
	// TODO: Implement this function
	return 0, false
}

// ============================================================================
// EXERCISE 3: Type Switches
// ============================================================================

// DescribeType returns a description of the type of the value.
//
// REQUIREMENTS:
// - For int: return "Integer: <value>"
// - For string: return "String: <value>"
// - For bool: return "Boolean: <value>"
// - For Person: return "Person: <name>"
// - For nil: return "Nil"
// - For any other type: return "Unknown"
//
// HINT: Use a type switch: switch v := i.(type) { ... }
func DescribeType(i interface{}) string {
	// TODO: Implement this function
	return ""
}

// ============================================================================
// EXERCISE 4: Interface Nil Check
// ============================================================================

// IsValidEmail checks if a Validator is valid, handling nil correctly.
//
// REQUIREMENTS:
// - If v is nil (a true nil interface), return false
// - If v is not nil but contains a nil pointer (like (*Email)(nil)), return false
// - Otherwise, call v.IsValid() and return its result
//
// HINT: First check if v == nil, then try type assertion to *Email and check if that's nil.
func IsValidEmail(v Validator) bool {
	// TODO: Implement this function
	// This tests your understanding of the two-part nil problem!
	return false
}

// ============================================================================
// EXERCISE 5: Implementing Multiple Interfaces
// ============================================================================

// TODO: Implement Read() for Buffer to satisfy the Reader interface.
//
// REQUIREMENTS:
// - Return the current data in the buffer
//
// HINT: Just return b.data
func (b *Buffer) Read() string {
	// TODO: Implement this method
	return ""
}

// TODO: Implement Write() for Buffer to satisfy the Writer interface.
//
// REQUIREMENTS:
// - Append the data to the buffer's existing data
// - Return nil (no errors for in-memory buffer)
//
// HINT: b.data += data
func (b *Buffer) Write(data string) error {
	// TODO: Implement this method
	return nil
}

// ============================================================================
// EXERCISE 6: Interface Composition
// ============================================================================

// IsReadWriter checks if an interface value implements ReadWriter.
//
// REQUIREMENTS:
// - Return true if i implements both Read() and Write()
// - Return false otherwise
//
// HINT: Use type assertion to ReadWriter: _, ok := i.(ReadWriter)
func IsReadWriter(i interface{}) bool {
	// TODO: Implement this function
	return false
}

// ============================================================================
// EXERCISE 7: Method Sets and Receivers
// ============================================================================

// TODO: Implement Increment() with a POINTER receiver.
//
// REQUIREMENTS:
// - Increment the Value field by 1
// - Use a pointer receiver so the change persists
//
// HINT: func (c *Counter) Increment() { c.Value++ }
func (c *Counter) Increment() {
	// TODO: Implement this method
}

// CanIncrement checks if a value can be used as an Incrementer.
//
// REQUIREMENTS:
// - Return true if i implements the Incrementer interface
// - Return false otherwise
//
// HINT: Use type assertion to Incrementer
func CanIncrement(i interface{}) bool {
	// TODO: Implement this function
	return false
}

// ============================================================================
// EXERCISE 8: Working with Empty Interface
// ============================================================================

// CountTypes counts how many values of each type are in the slice.
//
// REQUIREMENTS:
// - Return a map where keys are type names (as strings) and values are counts
// - Use %T with fmt.Sprintf to get type names
//
// EXAMPLE:
//   values := []interface{}{1, 2, "hello", "world", 1, true}
//   CountTypes(values) → map[string]int{"int": 3, "string": 2, "bool": 1}
//
// HINT: Use fmt.Sprintf("%T", v) to get type names.
func CountTypes(values []interface{}) map[string]int {
	// TODO: Implement this function
	return nil
}

// ============================================================================
// EXERCISE 9: Error Interface
// ============================================================================

// TODO: Implement Error() for ValidationError to satisfy the error interface.
//
// REQUIREMENTS:
// - Return a string in the format: "validation error on <field>: <message>"
// - Example: ValidationError{Field: "email", Message: "invalid format"}
//   → "validation error on email: invalid format"
//
// HINT: The error interface is: type error interface { Error() string }
func (e ValidationError) Error() string {
	// TODO: Implement this method
	return ""
}

// ============================================================================
// EXERCISE 10: Polymorphism
// ============================================================================

// TODO: Implement Area() for Rectangle.
//
// REQUIREMENTS:
// - Return width * height
//
// HINT: func (r Rectangle) Area() float64 { return r.Width * r.Height }
func (r Rectangle) Area() float64 {
	// TODO: Implement this method
	return 0
}

// TODO: Implement Area() for Circle.
//
// REQUIREMENTS:
// - Return π * radius²
// - Use 3.14159 for π
//
// HINT: func (c Circle) Area() float64 { return 3.14159 * c.Radius * c.Radius }
func (c Circle) Area() float64 {
	// TODO: Implement this method
	return 0
}

// TotalArea calculates the total area of all shapes.
//
// REQUIREMENTS:
// - Sum the areas of all shapes in the slice
// - Return the total
//
// HINT: Loop through shapes and sum shape.Area()
func TotalArea(shapes []Shape) float64 {
	// TODO: Implement this function
	return 0
}
