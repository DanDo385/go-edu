package exercise

import (
	"testing"
)

func TestSafeDeref(t *testing.T) {
	tests := []struct {
		name         string
		ptr          *int
		defaultValue int
		expected     int
	}{
		{
			name:         "nil pointer returns default",
			ptr:          nil,
			defaultValue: 99,
			expected:     99,
		},
		{
			name:         "valid pointer returns value",
			ptr:          intPtr(42),
			defaultValue: 99,
			expected:     42,
		},
		{
			name:         "zero value pointer",
			ptr:          intPtr(0),
			defaultValue: 99,
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeDeref(tt.ptr, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestSwap(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expectA, expectB int
	}{
		{
			name:    "swap positive numbers",
			a:       10,
			b:       20,
			expectA: 20,
			expectB: 10,
		},
		{
			name:    "swap with zero",
			a:       0,
			b:       42,
			expectA: 42,
			expectB: 0,
		},
		{
			name:    "swap identical values",
			a:       5,
			b:       5,
			expectA: 5,
			expectB: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, b := tt.a, tt.b
			Swap(&a, &b)
			if a != tt.expectA {
				t.Errorf("Expected a=%d, got %d", tt.expectA, a)
			}
			if b != tt.expectB {
				t.Errorf("Expected b=%d, got %d", tt.expectB, b)
			}
		})
	}
}

func TestInitializeMap(t *testing.T) {
	t.Run("nil map is initialized", func(t *testing.T) {
		var m map[string]int
		result := InitializeMap(m)

		if result == nil {
			t.Fatal("Expected non-nil map, got nil")
		}

		// Should be able to write to it
		result["test"] = 42
		if result["test"] != 42 {
			t.Errorf("Expected 42, got %d", result["test"])
		}
	})

	t.Run("existing map is unchanged", func(t *testing.T) {
		m := map[string]int{"existing": 100}
		result := InitializeMap(m)

		if result["existing"] != 100 {
			t.Errorf("Expected 100, got %d", result["existing"])
		}
	})
}

func TestAppendNode(t *testing.T) {
	t.Run("append to empty list", func(t *testing.T) {
		var head *Node
		head = AppendNode(head, 1)

		if head == nil {
			t.Fatal("Expected non-nil head")
		}
		if head.Value != 1 {
			t.Errorf("Expected value 1, got %d", head.Value)
		}
		if head.Next != nil {
			t.Error("Expected Next to be nil")
		}
	})

	t.Run("append multiple nodes", func(t *testing.T) {
		var head *Node
		head = AppendNode(head, 1)
		head = AppendNode(head, 2)
		head = AppendNode(head, 3)

		values := []int{1, 2, 3}
		current := head
		for i, expected := range values {
			if current == nil {
				t.Fatalf("Expected node at position %d, got nil", i)
			}
			if current.Value != expected {
				t.Errorf("At position %d: expected %d, got %d", i, expected, current.Value)
			}
			current = current.Next
		}

		if current != nil {
			t.Error("Expected end of list, but found more nodes")
		}
	})
}

func TestListLength(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Node
		expected int
	}{
		{
			name:     "empty list",
			setup:    func() *Node { return nil },
			expected: 0,
		},
		{
			name: "single node",
			setup: func() *Node {
				return &Node{Value: 1}
			},
			expected: 1,
		},
		{
			name: "three nodes",
			setup: func() *Node {
				return &Node{
					Value: 1,
					Next: &Node{
						Value: 2,
						Next: &Node{
							Value: 3,
						},
					},
				}
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := tt.setup()
			result := ListLength(head)
			if result != tt.expected {
				t.Errorf("Expected length %d, got %d", tt.expected, result)
			}
		})
	}
}

// Helper function to create int pointer
func intPtr(x int) *int {
	return &x
}
