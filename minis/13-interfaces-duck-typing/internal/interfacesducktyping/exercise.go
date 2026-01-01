//go:build !solution && !reference

package interfacesducktyping

import (
	"fmt"
	"reflect"
)

// String implements the exercise.
//
// TODO: Implement this function
func (p Person) String() string {
	// TODO: Implement
	return ""
}

// GetAge implements the exercise.
//
// TODO: Implement this function
func GetAge(s Stringer) (int, bool) {
	// TODO: Implement
	return 0, false
}

// DescribeType implements the exercise.
//
// TODO: Implement this function
func DescribeType(i interface{}) string {
	// TODO: Implement
	return ""
}

// IsValidEmail implements the exercise.
//
// TODO: Implement this function
func IsValidEmail(v Validator) bool {
	// TODO: Implement
	return false
}

// Read implements the exercise.
//
// TODO: Implement this function
func (b *Buffer) Read() string {
	// TODO: Implement
	return ""
}

// Write implements the exercise.
//
// TODO: Implement this function
func (b *Buffer) Write(data string) error {
	// TODO: Implement
	return nil
}

// IsReadWriter implements the exercise.
//
// TODO: Implement this function
func IsReadWriter(i interface{}) bool {
	// TODO: Implement
	return false
}

// Increment implements the exercise.
//
// TODO: Implement this function
func (c *Counter) Increment() {
	// TODO: Implement
}

// CanIncrement implements the exercise.
//
// TODO: Implement this function
func CanIncrement(i interface{}) bool {
	// TODO: Implement
	return false
}

// CountTypes implements the exercise.
//
// TODO: Implement this function
func CountTypes(values []interface{}) map[string]int {
	// TODO: Implement
	return nil
}

// Error implements the exercise.
//
// TODO: Implement this function
func (e ValidationError) Error() string {
	// TODO: Implement
	return ""
}

// Area implements the exercise.
//
// TODO: Implement this function
func (r Rectangle) Area() float64 {
	// TODO: Implement
	return 0
}

// Area implements the exercise.
//
// TODO: Implement this function
func (c Circle) Area() float64 {
	// TODO: Implement
	return 0
}

// TotalArea implements the exercise.
//
// TODO: Implement this function
func TotalArea(shapes []Shape) float64 {
	// TODO: Implement
	return 0
}
