// Package exercise provides type definitions shared between exercises and tests.
package exercise

import "fmt"

// ============================================================================
// Shared Types for Exercises
// ============================================================================

// User represents a user with JSON and validation tags.
type User struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"email"`
	Age   int    `json:"age" validate:"min=0,max=150"`
}

// Product represents a product with database mapping tags.
type Product struct {
	ID    int     `db:"id" json:"id"`
	Name  string  `db:"product_name" json:"name"`
	Price float64 `db:"price" json:"price"`
}

// Point represents a 2D point.
type Point struct {
	X int
	Y int
}

// Rectangle represents a rectangle shape.
type Rectangle struct {
	Width  float64
	Height float64
}

// Circle represents a circle shape.
type Circle struct {
	Radius float64
}

// Calculator provides arithmetic operations.
type Calculator struct{}

// Add returns the sum of two numbers.
func (c Calculator) Add(a, b int) int {
	return a + b
}

// Multiply returns the product of two numbers.
func (c Calculator) Multiply(a, b int) int {
	return a * b
}

// Counter represents a simple counter.
type Counter struct {
	Value int
}

// Increment increases the counter by 1.
func (c *Counter) Increment() {
	c.Value++
}

// Decrement decreases the counter by 1.
func (c *Counter) Decrement() {
	c.Value--
}

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}
