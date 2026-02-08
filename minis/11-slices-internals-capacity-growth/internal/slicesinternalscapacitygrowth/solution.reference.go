//go:build reference

package slicesinternalscapacitygrowth

/*
Reference Solution - Slice Internals: Capacity, Growth, Backing Arrays
=====================================================================

First principles (per .cursorrules):

A slice is a header: (ptr, len, cap). ptr points to the backing array.
Copying a slice copies the header — both slices share the SAME backing array.
Mutating s1[i] changes s2[i] if they share. "Think of the slice as a business
card that says 'data starts here, length X, capacity Y' — the card is cheap to
copy, but the building (array) is shared."

append(s, x): If cap permits, writes x into the existing array, increments len.
If len==cap, append allocates a NEW array (typically 2x), copies old data, then
appends. The returned slice may have a different ptr. Old slice still points to
old array — they no longer share.

s[low:high:max] — 3-index slice: cap = max-low. Future append cannot overwrite
parent's elements beyond high, because cap limits growth.
*/

// GrowSlice - Append Element and Report Capacity Change
//
// When len(s) < cap(s), append reuses the backing array; oldCap == newCap.
// When len(s) == cap(s), append allocates a new larger array (Go typically doubles);
// the returned slice has a different backing array.
func GrowSlice(s []int, elem int) (newSlice []int, oldCap, newCap int) {
	oldCap = cap(s)
	newSlice = append(s, elem)
	newCap = cap(newSlice)
	return newSlice, oldCap, newCap
}

// SharesBackingArray - Detect If Two Slices Share Storage
//
// Mutates each a[i] to a sentinel, checks if b contains it, restores a[i].
// We must try all indices in a because b might be a sub-slice (e.g. b = a[1:4])
// that doesn't include a[0]. Empty slices can't share. Teaching trick; production
// would use reflect or unsafe.
func SharesBackingArray(a, b []int) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	for i := range a {
		original := a[i]
		sentinel := original + 999999
		a[i] = sentinel

		shared := false
		for j := range b {
			if b[j] == sentinel {
				shared = true
				break
			}
		}
		a[i] = original
		if shared {
			return true
		}
	}
	return false
}

// SafeTruncate - Truncate Without Aliasing
//
// s[:n] shares the backing array with s. Caller could append to the result and
// overwrite s's elements. SafeTruncate returns a copy so the result is independent.
// Clamps n to [0, len(s)] for safety.
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

// PreallocateVsDynamic - Count Allocations
//
// Dynamic: append to nil slice; capacity grows as needed (1,2,4,8,...).
// Prealloc: make([]int, 0, n) reserves capacity upfront; append never reallocates
// until n elements. Returns number of times capacity changed in each case.
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

// ReSliceWithCapLimit - 3-Index Slice Expression
//
// s[low:high:max] creates a slice with len=high-low, cap=max-low.
// By using end for both high and max: s[start:end:end], we set cap = end-start = len.
// Future append must allocate; the new slice won't overwrite s's elements beyond end.
func ReSliceWithCapLimit(s []int, start, end int) []int {
	return s[start:end:end]
}
