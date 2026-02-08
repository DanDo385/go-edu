//go:build reference

package interfacesducktyping

/*
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
*/

import (
	"fmt"
	"reflect"
)

// String implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (p Person) String() string {
	return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
}

// GetAge implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func GetAge(s Stringer) (int, bool) {
	p, ok := s.(Person)
	if !ok {
		return 0, false
	}
	return p.Age, true
}

// DescribeType implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func DescribeType(i interface{}) string {
	switch v := i.(type) {
	case int:
		return fmt.Sprintf("Integer: %d", v)
	case string:
		return fmt.Sprintf("String: %s", v)
	case bool:
		return fmt.Sprintf("Boolean: %t", v)
	case Person:
		return fmt.Sprintf("Person: %s", v.Name)
	case nil:
		return "Nil"
	default:
		return "Unknown"
	}
}

// IsValidEmail implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func IsValidEmail(v Validator) bool {
	if v == nil {
		return false
	}
	if e, ok := v.(*Email); ok && e == nil {
		return false
	}
	return v.IsValid()
}

// Read implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (b *Buffer) Read() string {
	if b == nil {
		return ""
	}
	return b.data
}

// Write implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func (b *Buffer) Write(data string) error {
	if b == nil {
		return fmt.Errorf("nil buffer")
	}
	b.data += data
	return nil
}

// IsReadWriter implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func IsReadWriter(i interface{}) bool {
	_, ok := i.(ReadWriter)
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
