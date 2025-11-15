package exercise

// GrowSlice appends an element and tracks capacity changes.
func GrowSlice(s []int, elem int) (newSlice []int, oldCap, newCap int) {
	// MICRO-COMMENT: Capture the old capacity before append
	oldCap = cap(s)

	// MICRO-COMMENT: Append the element
	// If oldCap == len(s), this will trigger a reallocation
	newSlice = append(s, elem)

	// MICRO-COMMENT: Capture the new capacity after append
	newCap = cap(newSlice)

	return newSlice, oldCap, newCap
}

// SharesBackingArray detects if two slices share the same backing array.
func SharesBackingArray(a, b []int) bool {
	// MICRO-COMMENT: Edge case - if either slice is empty, they can't share
	if len(a) == 0 || len(b) == 0 {
		return false
	}

	// MACRO-COMMENT: Simple detection strategy without unsafe:
	// 1. Save the original value of a[0]
	// 2. Modify a[0] to a sentinel value
	// 3. Check if b was affected
	// 4. Restore a[0]
	//
	// This works because if they share a backing array AND overlap,
	// modifying a[0] will be visible in b.
	//
	// LIMITATION: This doesn't detect ALL cases of sharing (e.g., if
	// they share the array but don't overlap). For complete detection,
	// you'd need unsafe.Pointer to compare addresses.

	original := a[0]
	sentinel := original + 999999 // Use a value unlikely to already exist

	a[0] = sentinel
	shares := false

	// Check if any element in b has the sentinel value
	for i := range b {
		if b[i] == sentinel {
			shares = true
			break
		}
	}

	// Restore original value
	a[0] = original

	return shares
}

// SafeTruncate creates an independent truncated slice.
func SafeTruncate(s []int, n int) []int {
	// MICRO-COMMENT: Handle edge case where n > len(s)
	if n > len(s) {
		n = len(s)
	}

	// MICRO-COMMENT: Handle zero-length result
	if n == 0 {
		return []int{}
	}

	// MICRO-COMMENT: Allocate a NEW backing array with exact size needed
	result := make([]int, n)

	// MICRO-COMMENT: Copy elements from original to the new array
	// This ensures the new slice is completely independent
	copy(result, s[:n])

	return result
}

// PreallocateVsDynamic compares reallocation counts.
func PreallocateVsDynamic(n int) (dynamicAllocs, preallocAllocs int) {
	// APPROACH 1: Dynamic growth
	var s1 []int
	prevCap1 := 0

	for i := 0; i < n; i++ {
		s1 = append(s1, i)
		if cap(s1) != prevCap1 {
			dynamicAllocs++
			prevCap1 = cap(s1)
		}
	}

	// APPROACH 2: Pre-allocated
	s2 := make([]int, 0, n)
	prevCap2 := cap(s2)

	for i := 0; i < n; i++ {
		s2 = append(s2, i)
		if cap(s2) != prevCap2 {
			preallocAllocs++
			prevCap2 = cap(s2)
		}
	}

	return dynamicAllocs, preallocAllocs
}

// ReSliceWithCapLimit creates a capacity-limited sub-slice.
func ReSliceWithCapLimit(s []int, start, end int) []int {
	// MICRO-COMMENT: Use 3-index slice syntax: s[low:high:max]
	// Setting max = end means:
	// - Capacity = max - low = end - start
	// - This equals the length, so there's no spare capacity
	// - Any append will force reallocation
	return s[start:end:end]
}
