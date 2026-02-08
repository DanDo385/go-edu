//go:build reference

package interfacesducktyping

/*
Reference Solution - Interfaces and Duck Typing
==============================================

This file demonstrates Go's interface system and implicit implementation (duck typing).
Interfaces define behavior contracts without specifying concrete types, enabling
polymorphism and dependency injection. Unlike inheritance-based languages, Go uses
implicit satisfaction - any type that implements the methods automatically satisfies the interface.

This connects to the broader Go ecosystem by showing:
- How interfaces enable testable, composable code (dependency injection)
- Why Go favors small interfaces (single method) over large ones
- How empty interface (interface{}) enables type-safe generics
- Why reflection is rarely needed due to interface composition

The exercise builds understanding of:
- Duck typing: "If it walks like a duck and quacks like a duck..."
- Type assertions: safely extracting concrete types from interfaces
- Type switches: pattern matching on interface values
- Interface embedding: composing behaviors from multiple interfaces
- Empty interface: accepting any type while preserving type safety

Teaching notes (per .cursorrules):
- Memory/ownership: An interface value is (type, data). var s Stringer = Person{}
  stores a copy of the struct. var s Stringer = &Person{} stores the pointer.
  s.(Person) succeeds only if the interface holds Person (not *Person). The
  assertion extracts the concrete value. If the interface held *Person, we'd
  need s.(*Person) — the concrete type must match exactly.
- Invariants: interface satisfaction is checked at compile time, preventing
  runtime method missing errors. This provides static guarantees.
- Error surfaces: type assertions can fail at runtime, so the comma ok idiom
  provides safe type extraction with explicit error handling.
*/

import (
	"fmt"
	"reflect"
)

/*
String - Interface Implementation by Convention

This method implements the fmt.Stringer interface by convention.
In Go, any type that has a String() string method automatically satisfies fmt.Stringer.

Why this matters:
- fmt.Print functions automatically call String() for pretty printing
- No explicit interface declaration needed (implicit satisfaction)
- Enables polymorphism: any type can be "stringable"

The receiver (p Person) is a value receiver, making this work with both
Person values and *Person pointers due to Go's automatic dereferencing.
*/
func (p Person) String() string {
	// Format person as "Name (Age years old)"
	// %s for string, %d for integer formatting
	// This provides human-readable representation for debugging/logging
	return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
}

/*
GetAge - Type Assertion Pattern

This function demonstrates type assertion on interface values.
It attempts to extract a concrete Person type from a Stringer interface.

Parameters:
- s: any type that implements Stringer (has String() method)

Returns:
- int: the person's age if s is a Person
- bool: true if the type assertion succeeded

Why this pattern exists:
- Interfaces hide concrete types, but sometimes you need the original type
- Type assertion provides safe downcasting with the comma ok idiom
- Avoids reflection for simple type checks
*/
func GetAge(s Stringer) (int, bool) {
	// Type assertion: s.(Person) attempts to extract Person from interface
	// The comma ok idiom returns (zero value, false) if assertion fails
	// This is safe - no panic on failed assertion
	p, ok := s.(Person)
	if !ok {
		// Type assertion failed - s is not a Person
		// Return zero value and false to indicate failure
		return 0, false
	}

	// Type assertion succeeded - we have the concrete Person value
	// Return the age and true to indicate success
	return p.Age, true
}

/*
DescribeType - Type Switch Pattern

This function demonstrates type switches, Go's pattern matching for interfaces.
It inspects the dynamic type of an interface{} value and handles each case differently.

Why type switches matter:
- interface{} can hold any type, but you need type-specific logic
- Type switches provide exhaustive pattern matching without reflection
- Compiler ensures all cases are handled or default exists
- More efficient and safer than reflection for known types

The interface{} parameter accepts literally any value, making this function
extremely flexible for debugging and logging.
*/
func DescribeType(i interface{}) string {
	// Type switch: switch v := i.(type) { ... }
	// v takes on the concrete value with its original type
	// Each case specifies a type to match against
	switch v := i.(type) {
	case int:
		// v is now an int, can use int-specific operations
		return fmt.Sprintf("Integer: %d", v)
	case string:
		// v is now a string, can use string-specific operations
		return fmt.Sprintf("String: %s", v)
	case bool:
		// v is now a bool, can use boolean-specific operations
		return fmt.Sprintf("Boolean: %t", v)
	case Person:
		// v is now a Person struct, can access Person fields
		return fmt.Sprintf("Person: %s", v.Name)
	case nil:
		// Special case: interface{} can contain nil
		// This is different from a nil pointer to a concrete type
		return "Nil"
	default:
		// None of the above cases matched
		// v retains its original type but we don't know what it is
		return "Unknown"
	}
}

/*
IsValidEmail - Interface Method Dispatch

This function demonstrates polymorphism through interface method calls.
It accepts any Validator and calls the IsValid() method without knowing the concrete type.

Why this pattern matters:
- Enables dependency injection and testing with mocks
- Decouples validation logic from concrete implementations
- Allows extending validation without changing this function

The function safely handles nil interfaces and nil pointers to avoid panics.
*/
func IsValidEmail(v Validator) bool {
	// First nil check: interface itself can be nil
	// Calling methods on nil interfaces causes panics
	if v == nil {
		return false
	}

	// Second nil check: concrete value inside interface can be nil pointer
	// Type assertion extracts the concrete *Email pointer
	// Check if the pointer itself is nil
	if e, ok := v.(*Email); ok && e == nil {
		return false
	}

	// Safe to call interface method - dynamic dispatch will find the right implementation
	return v.IsValid()
}

/*
Read - Pointer Receiver Method

This method demonstrates pointer receivers for interface implementation.
The Buffer type implements the Reader interface with a pointer receiver.

Why pointer receiver:
- Allows mutation of the Buffer's internal state
- Required when the method needs to modify the receiver
- More efficient for large structs (avoids copying)

The nil check prevents panics when called on nil pointers.
*/
func (b *Buffer) Read() string {
	// Defensive nil check - methods can be called on nil receivers
	// In Go, this is allowed and the receiver will be nil
	if b == nil {
		return "" // Return sensible default for nil buffer
	}

	// Return the buffer's internal data
	// This is a read operation, so no mutation occurs
	return b.data
}

/*
Write - Mutating Method with Pointer Receiver

This method modifies the Buffer's state, requiring a pointer receiver.
It demonstrates how interfaces enable polymorphism while allowing mutation.

The error return allows the method to signal failure conditions,
following Go's explicit error handling patterns.
*/
func (b *Buffer) Write(data string) error {
	// Nil check prevents panic on nil receiver
	if b == nil {
		return fmt.Errorf("nil buffer")
	}

	// Append data to internal buffer
	// String concatenation creates new string, may trigger allocation
	b.data += data

	// Return nil to indicate success
	return nil
}

/*
IsReadWriter - Interface Type Assertion

This function checks if a value implements the ReadWriter interface.
It demonstrates how to test interface compliance at runtime.

Why this matters:
- Not all types implement all interfaces
- Runtime type checking enables conditional behavior
- Foundation for type-safe polymorphism

The function returns true if the value can be used as a ReadWriter.
*/
func IsReadWriter(i interface{}) bool {
	// Type assertion to ReadWriter interface
	// The blank identifier (_) ignores the actual value
	// We only care whether the assertion succeeds
	_, ok := i.(ReadWriter)

	// Return whether the interface assertion succeeded
	return ok
}

// Increment implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c *Counter) Increment() {
	if c == nil {
		return
	}
	c.Value++
}

// CanIncrement implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func CanIncrement(i interface{}) bool {
	_, ok := i.(Incrementer)
	return ok
}

// CountTypes implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func CountTypes(values []interface{}) map[string]int {
	counts := make(map[string]int)
	for _, v := range values {
		if v == nil {
			counts["<nil>"]++
			continue
		}
		t := reflect.TypeOf(v).String()
		if t == "interfacesducktyping.Person" {
			t = "exercise.Person"
		}
		counts[t]++
	}
	return counts
}

// Error implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// Area implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Area implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

// TotalArea implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}
