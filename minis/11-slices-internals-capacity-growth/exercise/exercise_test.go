package exercise

import (
	"testing"
)

// TestGrowSlice verifies that GrowSlice correctly tracks capacity changes.
func TestGrowSlice(t *testing.T) {
	tests := []struct {
		name        string
		initial     []int
		elem        int
		expectGrow  bool // Whether we expect capacity to grow
	}{
		{
			name:       "empty slice grows",
			initial:    []int{},
			elem:       1,
			expectGrow: true,
		},
		{
			name:       "slice with spare capacity doesn't grow",
			initial:    make([]int, 2, 5),
			elem:       999,
			expectGrow: false,
		},
		{
			name:       "full slice grows",
			initial:    []int{1, 2, 3},
			elem:       4,
			expectGrow: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newSlice, oldCap, newCap := GrowSlice(tt.initial, tt.elem)

			// Check that element was appended
			if len(newSlice) != len(tt.initial)+1 {
				t.Errorf("Expected length %d, got %d", len(tt.initial)+1, len(newSlice))
			}

			// Check that the last element is the one we appended
			if newSlice[len(newSlice)-1] != tt.elem {
				t.Errorf("Expected last element %d, got %d", tt.elem, newSlice[len(newSlice)-1])
			}

			// Check capacity behavior
			if tt.expectGrow && newCap <= oldCap {
				t.Errorf("Expected capacity to grow, but oldCap=%d, newCap=%d", oldCap, newCap)
			}
			if !tt.expectGrow && newCap != oldCap {
				t.Errorf("Expected capacity to stay same, but oldCap=%d, newCap=%d", oldCap, newCap)
			}
		})
	}
}

// TestSharesBackingArray verifies detection of shared backing arrays.
func TestSharesBackingArray(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() ([]int, []int)
		expected bool
	}{
		{
			name: "sub-slice shares backing array",
			setup: func() ([]int, []int) {
				a := []int{1, 2, 3, 4, 5}
				b := a[1:4]
				return a, b
			},
			expected: true,
		},
		{
			name: "independent slices don't share",
			setup: func() ([]int, []int) {
				a := []int{1, 2, 3}
				b := make([]int, len(a))
				copy(b, a)
				return a, b
			},
			expected: false,
		},
		{
			name: "identical slice references share",
			setup: func() ([]int, []int) {
				a := []int{1, 2, 3}
				b := a // Same slice reference
				return a, b
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, b := tt.setup()
			result := SharesBackingArray(a, b)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestSafeTruncate verifies that SafeTruncate creates an independent copy.
func TestSafeTruncate(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		n        int
		expected []int
	}{
		{
			name:     "truncate large slice",
			initial:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			n:        3,
			expected: []int{1, 2, 3},
		},
		{
			name:     "truncate to zero",
			initial:  []int{1, 2, 3},
			n:        0,
			expected: []int{},
		},
		{
			name:     "n equals length",
			initial:  []int{1, 2, 3},
			n:        3,
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeTruncate(tt.initial, tt.n)

			// Check length
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}

			// Check values
			for i, v := range tt.expected {
				if result[i] != v {
					t.Errorf("At index %d: expected %d, got %d", i, v, result[i])
				}
			}

			// CRITICAL: Verify it's NOT sharing the backing array
			if len(result) > 0 && len(tt.initial) > 0 {
				if SharesBackingArray(result, tt.initial) {
					t.Error("SafeTruncate should create an independent slice")
				}
			}
		})
	}
}

// TestPreallocateVsDynamic verifies the allocation count difference.
func TestPreallocateVsDynamic(t *testing.T) {
	tests := []struct {
		name                 string
		n                    int
		maxDynamicAllocs     int // Dynamic should have multiple allocations
		expectedPreallocated int // Pre-allocated should have zero
	}{
		{
			name:                 "small slice",
			n:                    100,
			maxDynamicAllocs:     10,
			expectedPreallocated: 0,
		},
		{
			name:                 "medium slice",
			n:                    10000,
			maxDynamicAllocs:     25,
			expectedPreallocated: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dynamicAllocs, preallocAllocs := PreallocateVsDynamic(tt.n)

			// Dynamic should have multiple reallocations
			if dynamicAllocs < 1 {
				t.Errorf("Expected dynamic approach to have multiple allocations, got %d", dynamicAllocs)
			}

			if dynamicAllocs > tt.maxDynamicAllocs {
				t.Errorf("Expected dynamic allocations <= %d, got %d", tt.maxDynamicAllocs, dynamicAllocs)
			}

			// Pre-allocated should have zero reallocations
			if preallocAllocs != tt.expectedPreallocated {
				t.Errorf("Expected pre-allocated reallocations to be %d, got %d", tt.expectedPreallocated, preallocAllocs)
			}
		})
	}
}

// TestReSliceWithCapLimit verifies 3-index slice usage.
func TestReSliceWithCapLimit(t *testing.T) {
	tests := []struct {
		name        string
		initial     []int
		start       int
		end         int
		expectedLen int
		expectedCap int
	}{
		{
			name:        "middle section",
			initial:     []int{10, 20, 30, 40, 50},
			start:       1,
			end:         3,
			expectedLen: 2,
			expectedCap: 2,
		},
		{
			name:        "from start",
			initial:     []int{10, 20, 30, 40, 50},
			start:       0,
			end:         2,
			expectedLen: 2,
			expectedCap: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReSliceWithCapLimit(tt.initial, tt.start, tt.end)

			if len(result) != tt.expectedLen {
				t.Errorf("Expected len=%d, got %d", tt.expectedLen, len(result))
			}

			if cap(result) != tt.expectedCap {
				t.Errorf("Expected cap=%d, got %d", tt.expectedCap, cap(result))
			}

			// Verify values
			for i := 0; i < len(result); i++ {
				if result[i] != tt.initial[tt.start+i] {
					t.Errorf("At index %d: expected %d, got %d", i, tt.initial[tt.start+i], result[i])
				}
			}
		})
	}
}
