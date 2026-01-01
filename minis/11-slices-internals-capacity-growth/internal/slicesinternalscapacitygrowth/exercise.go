//go:build !solution && !reference

package slicesinternalscapacitygrowth

/*
Problem: Understanding Go slice internals and capacity growth patterns

Requirements:
1. Track capacity changes during append operations
2. Detect when slices share backing arrays
3. Safely truncate slices to allow garbage collection
4. Compare pre-allocation vs dynamic growth
5. Create capacity-limited sub-slices

Data Structure:
- Slice: Pointer to array + Length + Capacity (24 bytes on 64-bit)
- Backing array: Contiguous memory holding actual elements
- Multiple slices can reference same backing array

Algorithm: Slice Growth Strategy
- For cap < 256: new_cap = old_cap * 2 (doubling)
- For cap >= 256: new_cap ≈ old_cap * 1.25 + 192 (slower growth)
- Reallocation: Allocate new array, copy elements, return new slice

Why slices are tricky:
- Slicing creates views into same array (memory sharing)
- append() may or may not reallocate (depends on capacity)
- Holding small slice can prevent large array from being GC'd
*/

import "unsafe"

// GrowSlice appends an element and tracks capacity changes.
// BREAKPOINT: Set breakpoint here to trace append behavior
// DEBUG: Watch 's' to see input slice state (len, cap)
// DEBUG: Watch 'elem' to see value being appended
func GrowSlice(s []int, elem int) (newSlice []int, oldCap, newCap int) {
	// BREAKPOINT: Set breakpoint here to capture capacity before append
	// DEBUG: Watch 'oldCap' to see capacity before reallocation
	oldCap = cap(s)
	// DEBUG: Compare len(s) vs oldCap to predict reallocation

	// BREAKPOINT: Set breakpoint here to trace append operation
	// DEBUG: Watch 's' before append (original slice)
	newSlice = append(s, elem)
	// DEBUG: Watch 'newSlice' after append to see if address changed
	// DEBUG: If len(s) == oldCap, append will reallocate

	// BREAKPOINT: Set breakpoint here to capture new capacity
	// DEBUG: Watch 'newCap' to see capacity after append
	newCap = cap(newSlice)
	// DEBUG: If newCap > oldCap, reallocation occurred
	// DEBUG: Growth factor visible: newCap/oldCap ratio

	// DEBUG: Watch all three return values to analyze growth pattern
	return newSlice, oldCap, newCap
}

// SharesBackingArray detects if two slices share the same backing array.
// BREAKPOINT: Set breakpoint here to trace sharing detection
// DEBUG: Watch 'a' and 'b' to see input slices
func SharesBackingArray(a, b []int) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}

	elemSize := unsafe.Sizeof(a[0])

	aStart := uintptr(unsafe.Pointer(&a[0]))
	aEnd := aStart + uintptr(cap(a))*elemSize

	bStart := uintptr(unsafe.Pointer(&b[0]))
	bEnd := bStart + uintptr(cap(b))*elemSize

	// If either slice's backing region overlaps the other's start pointer,
	// they share the same backing array.
	return (bStart >= aStart && bStart < aEnd) || (aStart >= bStart && aStart < bEnd)
}

// SafeTruncate creates an independent truncated slice.
// BREAKPOINT: Set breakpoint here to trace safe truncation
// DEBUG: Watch 's' to see input slice
// DEBUG: Watch 'n' to see truncation size
func SafeTruncate(s []int, n int) []int {
	// BREAKPOINT: Set breakpoint here to check bounds
	// DEBUG: Watch comparison 'n > len(s)' to see if clamping needed
	if n > len(s) {
		// DEBUG: Clamping n to slice length
		n = len(s)
	}

	// BREAKPOINT: Set breakpoint here to handle zero case
	// DEBUG: Watch 'n' to check for empty result
	if n == 0 {
		// DEBUG: Returning empty slice (no allocation needed)
		return []int{}
	}

	// BREAKPOINT: Set breakpoint here to allocate new array
	// DEBUG: Watch 'result' to see new allocation
	result := make([]int, n)
	// DEBUG: New array with exact capacity n (not shared with s)
	// DEBUG: Watch cap(result) == n (no spare capacity)

	// BREAKPOINT: Set breakpoint here before copy
	// DEBUG: Watch 's[:n]' source slice
	// DEBUG: Watch 'result' destination slice
	copy(result, s[:n])
	// DEBUG: After copy, result is independent of s
	// DEBUG: Original large array can be GC'd when s goes out of scope

	// DEBUG: Watch 'result' to verify independent slice
	return result
}

// PreallocateVsDynamic compares reallocation counts.
// BREAKPOINT: Set breakpoint here to trace allocation comparison
// DEBUG: Watch 'n' to see test size
func PreallocateVsDynamic(n int) (dynamicAllocs, preallocAllocs int) {
	// APPROACH 1: Dynamic growth (multiple reallocations)

	// BREAKPOINT: Set breakpoint here to start dynamic test
	// DEBUG: Watch 's1' - nil slice (no allocation yet)
	var s1 []int
	// DEBUG: len=0, cap=0 initially

	prevCap1 := 0
	// DEBUG: Track previous capacity to detect changes

	// BREAKPOINT: Set breakpoint here before dynamic loop
	// DEBUG: Watch loop iterations to see reallocations
	for i := 0; i < n; i++ {
		// BREAKPOINT: Set breakpoint here before each append
		// DEBUG: Watch 'i' for current iteration
		// DEBUG: Watch 'len(s1)' and 'cap(s1)' before append
		s1 = append(s1, i)

		// BREAKPOINT: Set breakpoint here to check capacity change
		// DEBUG: Watch 'cap(s1)' after append
		if cap(s1) != prevCap1 {
			// BREAKPOINT: Set breakpoint here on reallocation
			// DEBUG: Reallocation occurred! Count it
			dynamicAllocs++
			// DEBUG: Watch 'dynamicAllocs' increment
			// DEBUG: Watch 'prevCap1' → 'cap(s1)' to see growth factor
			prevCap1 = cap(s1)
		}
	}
	// DEBUG: Watch final 'dynamicAllocs' (typically 15-20 for n=10000)

	// APPROACH 2: Pre-allocated (zero reallocations)

	// BREAKPOINT: Set breakpoint here to start pre-allocated test
	// DEBUG: Watch 's2' - allocated with exact capacity
	s2 := make([]int, 0, n)
	// DEBUG: len=0, cap=n (ready for n appends without reallocation)

	// BREAKPOINT: Set breakpoint here after allocation
	prevCap2 := cap(s2)
	// DEBUG: Watch 'prevCap2' == n (full capacity from start)

	// BREAKPOINT: Set breakpoint here before pre-allocated loop
	for i := 0; i < n; i++ {
		// BREAKPOINT: Set breakpoint here before each append
		// DEBUG: Watch 'len(s2)' grow while 'cap(s2)' stays constant
		s2 = append(s2, i)

		// BREAKPOINT: Set breakpoint here to verify no reallocation
		// DEBUG: This check should never trigger
		if cap(s2) != prevCap2 {
			// BREAKPOINT: If this hits, something went wrong
			// DEBUG: Should never increment (we pre-allocated enough)
			preallocAllocs++
			prevCap2 = cap(s2)
		}
	}
	// DEBUG: Watch final 'preallocAllocs' (should be 0)

	// DEBUG: Compare dynamicAllocs vs preallocAllocs
	// DEBUG: Pre-allocation avoids all reallocations!
	return dynamicAllocs, preallocAllocs
}

// ReSliceWithCapLimit creates a capacity-limited sub-slice.
// BREAKPOINT: Set breakpoint here to trace 3-index slicing
// DEBUG: Watch 's' to see input slice
// DEBUG: Watch 'start' and 'end' to see slice bounds
func ReSliceWithCapLimit(s []int, start, end int) []int {
	// BREAKPOINT: Set breakpoint here before 3-index slice
	// DEBUG: Watch 's[start:end:end]' syntax
	// DEBUG: First end: length = end - start
	// DEBUG: Second end: capacity = end - start (same as length!)
	// DEBUG: Result has NO spare capacity - append will reallocate
	return s[start:end:end]
	// DEBUG: Watch return value cap == len (full capacity)
	// DEBUG: This prevents accidental modification of original array
}
