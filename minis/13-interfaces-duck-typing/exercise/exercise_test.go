package exercise

import (
	"testing"
)

// ============================================================================
// EXERCISE 1: Implementing Interfaces
// ============================================================================

func TestPersonStringer(t *testing.T) {
	tests := []struct {
		name     string
		person   Person
		expected string
	}{
		{
			name:     "adult",
			person:   Person{Name: "Alice", Age: 30},
			expected: "Alice (30 years old)",
		},
		{
			name:     "child",
			person:   Person{Name: "Bob", Age: 10},
			expected: "Bob (10 years old)",
		},
		{
			name:     "elderly",
			person:   Person{Name: "Charlie", Age: 75},
			expected: "Charlie (75 years old)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.person.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}

			// Verify Person implements Stringer
			var _ Stringer = tt.person
		})
	}
}

// ============================================================================
// EXERCISE 2: Type Assertions
// ============================================================================

func TestGetAge(t *testing.T) {
	tests := []struct {
		name        string
		stringer    Stringer
		expectedAge int
		expectedOk  bool
	}{
		{
			name:        "person returns age",
			stringer:    Person{Name: "Alice", Age: 30},
			expectedAge: 30,
			expectedOk:  true,
		},
		{
			name:        "another person",
			stringer:    Person{Name: "Bob", Age: 25},
			expectedAge: 25,
			expectedOk:  true,
		},
		{
			name:        "non-person returns false",
			stringer:    customStringer("not a person"),
			expectedAge: 0,
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			age, ok := GetAge(tt.stringer)
			if age != tt.expectedAge {
				t.Errorf("Expected age %d, got %d", tt.expectedAge, age)
			}
			if ok != tt.expectedOk {
				t.Errorf("Expected ok %t, got %t", tt.expectedOk, ok)
			}
		})
	}
}

// Helper type for testing
type customStringer string

func (c customStringer) String() string {
	return string(c)
}

// ============================================================================
// EXERCISE 3: Type Switches
// ============================================================================

func TestDescribeType(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "integer",
			value:    42,
			expected: "Integer: 42",
		},
		{
			name:     "string",
			value:    "hello",
			expected: "String: hello",
		},
		{
			name:     "boolean true",
			value:    true,
			expected: "Boolean: true",
		},
		{
			name:     "boolean false",
			value:    false,
			expected: "Boolean: false",
		},
		{
			name:     "person",
			value:    Person{Name: "Alice", Age: 30},
			expected: "Person: Alice",
		},
		{
			name:     "nil",
			value:    nil,
			expected: "Nil",
		},
		{
			name:     "unknown type",
			value:    []int{1, 2, 3},
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DescribeType(tt.value)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// ============================================================================
// EXERCISE 4: Interface Nil Check
// ============================================================================

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		validator Validator
		expected bool
	}{
		{
			name:     "valid email",
			validator: &Email{Address: "test@example.com"},
			expected: true,
		},
		{
			name:     "invalid email",
			validator: &Email{Address: "invalid"},
			expected: false,
		},
		{
			name:     "nil interface",
			validator: nil,
			expected: false,
		},
		{
			name:     "nil pointer in interface",
			validator: (*Email)(nil),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidEmail(tt.validator)
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

// ============================================================================
// EXERCISE 5: Implementing Multiple Interfaces
// ============================================================================

func TestBuffer(t *testing.T) {
	t.Run("implements Reader", func(t *testing.T) {
		b := &Buffer{data: "hello"}

		// Verify it implements Reader
		var _ Reader = b

		result := b.Read()
		if result != "hello" {
			t.Errorf("Expected %q, got %q", "hello", result)
		}
	})

	t.Run("implements Writer", func(t *testing.T) {
		b := &Buffer{}

		// Verify it implements Writer
		var _ Writer = b

		err := b.Write("hello")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if b.data != "hello" {
			t.Errorf("Expected data %q, got %q", "hello", b.data)
		}
	})

	t.Run("read and write", func(t *testing.T) {
		b := &Buffer{}

		b.Write("hello ")
		b.Write("world")

		result := b.Read()
		if result != "hello world" {
			t.Errorf("Expected %q, got %q", "hello world", result)
		}
	})
}

// ============================================================================
// EXERCISE 6: Interface Composition
// ============================================================================

func TestIsReadWriter(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "buffer is readwriter",
			value:    &Buffer{},
			expected: true,
		},
		{
			name:     "string is not readwriter",
			value:    "hello",
			expected: false,
		},
		{
			name:     "int is not readwriter",
			value:    42,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsReadWriter(tt.value)
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

// ============================================================================
// EXERCISE 7: Method Sets and Receivers
// ============================================================================

func TestCounter(t *testing.T) {
	t.Run("pointer implements incrementer", func(t *testing.T) {
		c := &Counter{Value: 0}

		// Verify *Counter implements Incrementer
		var _ Incrementer = c

		c.Increment()
		if c.Value != 1 {
			t.Errorf("Expected Value 1, got %d", c.Value)
		}

		c.Increment()
		c.Increment()
		if c.Value != 3 {
			t.Errorf("Expected Value 3, got %d", c.Value)
		}
	})

	t.Run("value does not implement incrementer", func(t *testing.T) {
		// This should NOT compile if uncommented:
		// var _ Incrementer = Counter{Value: 0}

		// But we can test CanIncrement
		if CanIncrement(Counter{Value: 0}) {
			t.Error("Counter value should not implement Incrementer")
		}
	})
}

func TestCanIncrement(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "counter pointer can increment",
			value:    &Counter{Value: 0},
			expected: true,
		},
		{
			name:     "counter value cannot increment",
			value:    Counter{Value: 0},
			expected: false,
		},
		{
			name:     "int cannot increment",
			value:    42,
			expected: false,
		},
		{
			name:     "string cannot increment",
			value:    "hello",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanIncrement(tt.value)
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

// ============================================================================
// EXERCISE 8: Working with Empty Interface
// ============================================================================

func TestCountTypes(t *testing.T) {
	tests := []struct {
		name     string
		values   []interface{}
		expected map[string]int
	}{
		{
			name:   "mixed types",
			values: []interface{}{1, 2, "hello", "world", true, 3, false, "!"},
			expected: map[string]int{
				"int":    3,
				"string": 3,
				"bool":   2,
			},
		},
		{
			name:     "empty slice",
			values:   []interface{}{},
			expected: map[string]int{},
		},
		{
			name:   "single type",
			values: []interface{}{1, 2, 3, 4, 5},
			expected: map[string]int{
				"int": 5,
			},
		},
		{
			name:   "with person",
			values: []interface{}{Person{Name: "Alice", Age: 30}, 42, "hello", Person{Name: "Bob", Age: 25}},
			expected: map[string]int{
				"exercise.Person": 2,
				"int":             1,
				"string":          1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountTypes(tt.values)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d types, got %d", len(tt.expected), len(result))
			}

			for typ, count := range tt.expected {
				if result[typ] != count {
					t.Errorf("For type %s: expected count %d, got %d", typ, count, result[typ])
				}
			}
		})
	}
}

// ============================================================================
// EXERCISE 9: Error Interface
// ============================================================================

func TestValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      ValidationError
		expected string
	}{
		{
			name:     "email error",
			err:      ValidationError{Field: "email", Message: "invalid format"},
			expected: "validation error on email: invalid format",
		},
		{
			name:     "password error",
			err:      ValidationError{Field: "password", Message: "too short"},
			expected: "validation error on password: too short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}

			// Verify ValidationError implements error
			var _ error = tt.err
		})
	}
}

// ============================================================================
// EXERCISE 10: Polymorphism
// ============================================================================

func TestShapes(t *testing.T) {
	t.Run("rectangle area", func(t *testing.T) {
		r := Rectangle{Width: 4, Height: 5}

		// Verify Rectangle implements Shape
		var _ Shape = r

		area := r.Area()
		expected := 20.0
		if area != expected {
			t.Errorf("Expected area %.2f, got %.2f", expected, area)
		}
	})

	t.Run("circle area", func(t *testing.T) {
		c := Circle{Radius: 5}

		// Verify Circle implements Shape
		var _ Shape = c

		area := c.Area()
		expected := 78.53975  // π * 5²
		if area < expected-0.01 || area > expected+0.01 {
			t.Errorf("Expected area %.2f, got %.2f", expected, area)
		}
	})
}

func TestTotalArea(t *testing.T) {
	tests := []struct {
		name     string
		shapes   []Shape
		expected float64
	}{
		{
			name: "mixed shapes",
			shapes: []Shape{
				Rectangle{Width: 4, Height: 5},     // 20
				Circle{Radius: 2},                  // ~12.56636
				Rectangle{Width: 3, Height: 3},     // 9
			},
			expected: 41.56636,
		},
		{
			name:     "empty slice",
			shapes:   []Shape{},
			expected: 0,
		},
		{
			name: "only rectangles",
			shapes: []Shape{
				Rectangle{Width: 2, Height: 3},  // 6
				Rectangle{Width: 4, Height: 5},  // 20
			},
			expected: 26,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TotalArea(tt.shapes)
			if result < tt.expected-0.01 || result > tt.expected+0.01 {
				t.Errorf("Expected total area %.2f, got %.2f", tt.expected, result)
			}
		})
	}
}
