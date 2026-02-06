//go:build reference

package slicesinternalscapacitygrowth

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
	original := a[0]
	sentinel := original + 999999
	a[0] = sentinel

	shared := false
	for i := range b {
		if b[i] == sentinel {
			shared = true
			break
		}
	}
	a[0] = original
	return shared
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
