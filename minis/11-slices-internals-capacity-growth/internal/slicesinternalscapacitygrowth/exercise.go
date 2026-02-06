//go:build !solution && !reference

package slicesinternalscapacitygrowth

/*
Problem: Understanding Go slice internals and capacity growth patterns
Requirements:
1. Track capacity changes during append operations
2. Detect when slices share backing arrays
3. Safely truncate slices to allow garbage collection
Algorithm: Slice Growth Strategy
- For cap < 256: new_cap = old_cap * 2 (doubling)
- For cap >= 256: new_cap ≈ old_cap * 1.25 + 192 (slower growth)
- Reallocation: Allocate new array, copy elements, return new slice
*/

// GrowSlice appends elem and reports the old/new capacities.
func GrowSlice(s []int, elem int) (newSlice []int, oldCap, newCap int) {
	oldCap = cap(s)
	newSlice = append(s, elem)
	newCap = cap(newSlice)
	return newSlice, oldCap, newCap
}

// SharesBackingArray returns true when a and b share underlying storage.
func SharesBackingArray(a, b []int) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	for i := range a {
		original := a[i]
		sentinel := original + 999999 + i
		a[i] = sentinel
		for j := range b {
			if b[j] == sentinel {
				a[i] = original
				return true
			}
		}
		a[i] = original
	}
	return false
}

// SafeTruncate returns an independent copy of s[:n].
func SafeTruncate(s []int, n int) []int {
	if n < 0 {
		n = 0
	}
	if n > len(s) {
		n = len(s)
	}
	if n == 0 {
		return []int{}
	}
	out := make([]int, n)
	copy(out, s[:n])
	return out
}

// PreallocateVsDynamic compares capacity growth counts.
func PreallocateVsDynamic(n int) (dynamicAllocs, preallocAllocs int) {
	var dynamic []int
	prevCap := 0
	for i := 0; i < n; i++ {
		dynamic = append(dynamic, i)
		if cap(dynamic) != prevCap {
			dynamicAllocs++
			prevCap = cap(dynamic)
		}
	}

	prealloc := make([]int, 0, n)
	prevCap = cap(prealloc)
	for i := 0; i < n; i++ {
		prealloc = append(prealloc, i)
		if cap(prealloc) != prevCap {
			preallocAllocs++
			prevCap = cap(prealloc)
		}
	}

	return dynamicAllocs, preallocAllocs
}

// ReSliceWithCapLimit uses a 3-index slice to cap capacity at len.
func ReSliceWithCapLimit(s []int, start, end int) []int {
	return s[start:end:end]
}
