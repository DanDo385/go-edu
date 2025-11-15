//go:build !solution
// +build !solution

// Package exercise contains hands-on exercises for understanding slice internals.
//
// LEARNING OBJECTIVES:
// - Implement custom capacity growth logic
// - Detect shared backing arrays
// - Safely truncate slices without memory leaks
// - Measure allocation behavior

package exercise

// TODO: Implement these functions according to the specifications in the tests.
// Each function tests a different aspect of slice mechanics.

// GrowSlice appends an element to a slice and returns information about
// whether a reallocation occurred.
//
// REQUIREMENTS:
// - Append the element to the slice
// - Return the new slice, old capacity, and new capacity
// - This helps you understand when append() causes reallocations
//
// HINT: Store cap(s) before calling append(), then compare with cap(s) after.
func GrowSlice(s []int, elem int) (newSlice []int, oldCap, newCap int) {
	// TODO: Implement this function
	return nil, 0, 0
}

// SharesBackingArray returns true if slices a and b share the same backing array.
//
// REQUIREMENTS:
// - Return true if modifying a[0] would affect b (or vice versa)
// - Return false if the slices have independent backing arrays
//
// HINT: Use unsafe.Pointer to compare the addresses of the first elements,
// or use a simpler approach: modify a[0], check if b was affected, then restore.
// For this exercise, use the simpler approach (no unsafe needed).
func SharesBackingArray(a, b []int) bool {
	// TODO: Implement this function
	// Edge cases:
	// - What if one or both slices are empty?
	// - What if the slices don't overlap?
	return false
}

// SafeTruncate truncates a large slice to a smaller size, ensuring the
// large backing array can be garbage collected.
//
// REQUIREMENTS:
// - Return a NEW slice containing only the first n elements
// - The returned slice must have its own backing array (use copy)
// - This prevents memory leaks from retaining large arrays
//
// EXAMPLE:
//   big := make([]int, 1000000)
//   small := SafeTruncate(big, 10)
//   big = nil  // Now the 1M-element array can be GC'd
//
// HINT: Allocate a new slice with make(), then use copy().
func SafeTruncate(s []int, n int) []int {
	// TODO: Implement this function
	// Edge case: What if n > len(s)?
	return nil
}

// PreallocateVsDynamic compares pre-allocated vs dynamic growth for building
// a slice of n elements.
//
// REQUIREMENTS:
// - Build two slices with n elements (0, 1, 2, ..., n-1)
// - Approach 1: Start with var s []int (no pre-allocation)
// - Approach 2: Start with s := make([]int, 0, n) (pre-allocated)
// - Return the number of reallocations for each approach
//
// HOW TO COUNT REALLOCATIONS:
// Every time cap(s) increases, a reallocation occurred.
//
// EXPECTED RESULT:
// For n=10000, approach 1 should have ~15-20 reallocations,
// while approach 2 should have 0.
func PreallocateVsDynamic(n int) (dynamicAllocs, preallocAllocs int) {
	// TODO: Implement this function
	return 0, 0
}

// ReSliceWithCapLimit creates a sub-slice with limited capacity using
// the 3-index slice expression.
//
// REQUIREMENTS:
// - Return a slice from index start to end (exclusive)
// - Set capacity such that appending even ONE element forces reallocation
// - Use the 3-index slice syntax: s[start:end:max]
//
// EXAMPLE:
//   s := []int{10, 20, 30, 40, 50}
//   sub := ReSliceWithCapLimit(s, 1, 3)
//   // sub should be [20, 30] with len=2, cap=2
//   sub = append(sub, 99)
//   // This append must allocate (cap was already full)
//
// HINT: Set max = end to make cap = end - start = len.
func ReSliceWithCapLimit(s []int, start, end int) []int {
	// TODO: Implement this function
	return nil
}
