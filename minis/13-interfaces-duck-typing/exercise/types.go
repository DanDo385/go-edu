package exercise

// This file contains type definitions shared by both exercise.go and solution.go.
// It's always included regardless of build tags.

// ============================================================================
// Interface Definitions
// ============================================================================

// Stringer is a simplified version of fmt.Stringer.
type Stringer interface {
	String() string
}

// Validator is an interface for types that can validate themselves.
type Validator interface {
	IsValid() bool
}

// Reader defines the contract for reading.
type Reader interface {
	Read() string
}

// Writer defines the contract for writing.
type Writer interface {
	Write(data string) error
}

// ReadWriter combines Reader and Writer.
type ReadWriter interface {
	Reader
	Writer
}

// Incrementer defines a contract for types that can increment.
type Incrementer interface {
	Increment()
}

// Shape is an interface for geometric shapes.
type Shape interface {
	Area() float64
}

// ============================================================================
// Struct Definitions
// ============================================================================

// Person represents a person with a name and age.
type Person struct {
	Name string
	Age  int
}

// Email represents an email address.
type Email struct {
	Address string
}

// IsValid checks if the email address is valid.
//
// MICRO-COMMENT: This makes Email implement Validator.
func (e *Email) IsValid() bool {
	// Simple check: must contain '@'
	for _, ch := range e.Address {
		if ch == '@' {
			return true
		}
	}
	return false
}

// Buffer is a simple in-memory buffer.
type Buffer struct {
	data string
}

// Counter represents a counter.
type Counter struct {
	Value int
}

// ValidationError represents a validation error with a field name.
type ValidationError struct {
	Field   string
	Message string
}

// Rectangle represents a rectangle.
type Rectangle struct {
	Width, Height float64
}

// Circle represents a circle.
type Circle struct {
	Radius float64
}
