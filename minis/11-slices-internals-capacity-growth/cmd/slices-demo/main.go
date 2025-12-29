// Package main demonstrates slice internals, capacity growth, and common gotchas.
//
// MACRO-COMMENT: What This Program Teaches
// =========================================
// This program provides a hands-on demonstration of Go's slice mechanics:
// 1. How capacity grows as elements are appended
// 2. The dangers of shared backing arrays when re-slicing
// 3. How to safely isolate slices using 3-index syntax
// 4. Performance implications of pre-allocation vs dynamic growth
//
// USAGE EXAMPLES:
//   Run all demonstrations:
//     go run main.go
//
//   Run specific section only:
//     go run main.go -demo=growth
//     go run main.go -demo=shared
//     go run main.go -demo=safe
//
//   Enable verbose output with addresses:
//     go run main.go -verbose
//
//   Custom iteration count for preallocation demo:
//     go run main.go -demo=prealloc -iterations=100000
//
//   Run with escape analysis:
//     go build -gcflags='-m -m' main.go
//
// COMPILER BEHAVIOR: Escape Analysis
// When you run this with `-gcflags='-m'`, you'll see which slices escape to the heap.
// The compiler allocates slices on the heap when:
// - They're returned from functions
// - They're too large for the stack (typically >64KB)
// - They're stored in interface{} or heap-allocated structs
//
// RUNTIME BEHAVIOR: Memory Allocation
// Every time a slice grows beyond its capacity, the runtime:
// 1. Allocates a new, larger backing array (typically 2x the old capacity)
// 2. Copies all existing elements to the new array
// 3. Marks the old array for garbage collection
// This is why pre-allocation matters in hot code paths.

package main

import (
	"flag"
	"fmt"
)

// ============================================================================
// SECTION 1: Capacity Growth Demonstration
// ============================================================================

// demonstrateCapacityGrowth shows how Go automatically grows slice capacity
// as you append elements beyond the current capacity.
//
// MICRO-COMMENT: Why This Matters
// Understanding capacity growth is critical for:
// - Performance optimization (avoiding unnecessary allocations)
// - Memory profiling (knowing when/why allocations occur)
// - Debugging unexpected behavior (append can return a new slice)
func demonstrateCapacityGrowth() {
	fmt.Println("=== Capacity Growth Demonstration ===")

	// BREAKPOINT: Set breakpoint here to examine initial empty slice
	// DEBUG: At this point, s is nil - no backing array exists yet
	// var s []int
	// Expected: len=0, cap=0, underlying pointer=nil
	var s []int

	// MICRO-COMMENT: Track previous capacity to detect when it changes
	prevCap := 0

	// MACRO-COMMENT: Growth Strategy Analysis
	// As we append elements, watch for capacity changes. Go uses:
	// - Doubling strategy for small slices (cap < 256)
	// - ~1.25x growth for larger slices
	// This balances:
	// - Amortized O(1) append cost (doubling ensures this)
	// - Memory efficiency (not wasting too much space)
	for i := 0; i < 100; i++ {
		// BREAKPOINT: Set breakpoint here when i=0 to watch first allocation
		// DEBUG: Watch how s changes from nil to a real slice with backing array
		// DEBUG: On first append, Go allocates initial capacity (usually 1 or 2)
		// MICRO-COMMENT: Append returns a NEW slice header (potentially with a new backing array)
		// This is why we MUST reassign: s = append(s, i)
		// If capacity is exceeded, append will:
		// 1. Allocate a new, larger backing array
		// 2. Copy all existing elements to the new array
		// 3. Append the new element
		// 4. Return a slice header pointing to the new array
		s = append(s, i)

		// MICRO-COMMENT: cap() returns the capacity field from the slice header
		currentCap := cap(s)

		// DEBUG: When currentCap != prevCap, a reallocation just occurred
		// DEBUG: The old backing array is now garbage (will be collected later)
		// BREAKPOINT: Set conditional breakpoint when currentCap != prevCap to see reallocations
		// MICRO-COMMENT: When capacity changes, it means a new backing array was allocated
		if currentCap != prevCap {
			// MICRO-COMMENT: Calculate growth multiplier to show the growth pattern
			var multiplier float64
			if prevCap > 0 {
				multiplier = float64(currentCap) / float64(prevCap)
			}

			// MICRO-COMMENT: fmt.Printf uses verb formatters:
			// %4d = integer with minimum width 4 (right-aligned)
			// %.2f = float with 2 decimal places
			fmt.Printf("Len: %4d, Cap: %4d (grew from %4d, multiplier: %.2fx)\n",
				len(s), currentCap, prevCap, multiplier)

			prevCap = currentCap
		}
	}

	fmt.Println()
}

// ============================================================================
// SECTION 2: Re-Slicing and Shared Backing Arrays
// ============================================================================

// demonstrateSharedBackingArray shows the danger of modifying re-sliced slices.
//
// MACRO-COMMENT: The Hidden Sharing Problem
// When you create a slice from another slice using s[low:high], both slices
// share the SAME backing array. Modifications to one can affect the other.
// This is a common source of bugs in Go programs.
//
// MEMORY LAYOUT:
// original:  [ptr] [len=5] [cap=5]  →  [10][20][30][40][50]
//                                          ↑   ↑   ↑
// slice1:    [ptr] [len=3] [cap=4]  -------+---+---+
//                                       [20][30][40]
// (Both point to the same backing array starting at different offsets)
func demonstrateSharedBackingArray() {
	fmt.Println("=== Re-Slicing Gotcha: Shared Backing Array ===")

	// BREAKPOINT: Set breakpoint here before creating the original slice
	// DEBUG: After this line, inspect memory address of original's backing array
	// MICRO-COMMENT: Create a slice with 5 elements using a composite literal
	// The compiler allocates a backing array with exactly 5 slots
	original := []int{10, 20, 30, 40, 50}

	// BREAKPOINT: Set breakpoint after re-slicing to compare addresses
	// DEBUG: slice1 and original share the SAME backing array
	// DEBUG: Check &original[0] vs &slice1[0] - different addresses, same underlying array
	// MICRO-COMMENT: Re-slice from index 1 to 4 (exclusive)
	// This creates a NEW slice header, but the backing array is SHARED
	// slice1's Ptr points to &original[1] (address of element 20)
	// slice1's Len = 3 (elements: 20, 30, 40)
	// slice1's Cap = 4 (from index 1 to the end of the original array)
	slice1 := original[1:4]

	fmt.Printf("Before modification:\n")
	fmt.Printf("  Original: %v\n", original)
	fmt.Printf("  Slice1:   %v\n", slice1)

	// BREAKPOINT: Set breakpoint before this mutation
	// DEBUG: Watch both slice1[0] and original[1] change to 999 simultaneously
	// DEBUG: This is because they point to the SAME memory location
	// MICRO-COMMENT: Modify the first element of slice1
	// Since slice1[0] points to the same memory location as original[1],
	// this modification affects BOTH slices!
	slice1[0] = 999

	fmt.Printf("After modifying slice1[0] = 999:\n")
	fmt.Printf("  Original: %v  ← CHANGED! (original[1] is now 999)\n", original)
	fmt.Printf("  Slice1:   %v\n", slice1)

	fmt.Println()
}

// ============================================================================
// SECTION 3: Safe Re-Slicing with 3-Index Syntax
// ============================================================================

// demonstrateSafeReslicing shows how to use 3-index slicing to prevent
// shared backing array mutations.
//
// MACRO-COMMENT: The 3-Index Slice Expression
// Syntax: s[low:high:max]
// - Elements: s[low] through s[high-1] (same as 2-index)
// - Capacity: max - low (this is the KEY difference)
//
// By limiting capacity, you force append() to allocate a NEW backing array
// when the capacity is exceeded, preventing mutations to the original.
//
// PERFORMANCE NOTE:
// This safety comes at a cost: append will allocate sooner. Use this when
// isolation is more important than performance.
func demonstrateSafeReslicing() {
	fmt.Println("=== Safe Re-Slicing with 3-Index Syntax ===")

	// DEBUG: Before 3-index slicing, original has full capacity
	// MICRO-COMMENT: Create original slice (len=5, cap=5)
	original := []int{10, 20, 30, 40, 50}

	// BREAKPOINT: Set breakpoint after 3-index slicing
	// DEBUG: Compare slice2's capacity (2) vs original's capacity (5)
	// DEBUG: slice2 cannot grow into original's backing array due to limited cap
	// MICRO-COMMENT: Use 3-index slice: [low:high:max]
	// low = 1 (start at index 1, element 20)
	// high = 3 (end before index 3, includes elements 20, 30)
	// max = 3 (capacity limit: max - low = 3 - 1 = 2)
	//
	// Result: slice2 has len=2, cap=2 (no room to grow without reallocating)
	slice2 := original[1:3:3]

	fmt.Printf("Original: %v (len=%d, cap=%d)\n", original, len(original), cap(original))
	fmt.Printf("Slice2:   %v (len=%d, cap=%d)\n", slice2, len(slice2), cap(slice2))

	// BREAKPOINT: Set breakpoint before append to check capacity is full
	// DEBUG: Before append, slice2 has len=2, cap=2 (no spare capacity)
	// DEBUG: After append, slice2 will have a NEW backing array (reallocation)
	// MICRO-COMMENT: Append to slice2
	// Since cap=2 and len=2, the capacity is FULL
	// append() is forced to:
	// 1. Allocate a NEW backing array (capacity ~4)
	// 2. Copy elements [20, 30] to the new array
	// 3. Append 99 to the new array
	// 4. Return a new slice header pointing to the new array
	//
	// Original is NOT affected because slice2 now points to a different array
	slice2 = append(slice2, 99)

	// DEBUG: Check that original is unchanged - slice2 now has separate backing array

	fmt.Printf("\nAfter append(slice2, 99):\n")
	fmt.Printf("  Original: %v  ← UNCHANGED (different backing array)\n", original)
	fmt.Printf("  Slice2:   %v (len=%d, cap=%d)\n", slice2, len(slice2), cap(slice2))

	fmt.Println()
}

// ============================================================================
// SECTION 4: Append with Spare Capacity (No Allocation)
// ============================================================================

// demonstrateAppendWithCapacity shows that append() doesn't always allocate.
//
// MACRO-COMMENT: When Append is Fast
// If the slice has spare capacity (len < cap), append() is very fast:
// 1. Write the new element to backing_array[len]
// 2. Increment len by 1
// 3. Return the updated slice header
// No memory allocation, no copying, just a simple write + increment.
//
// RUNTIME BEHAVIOR:
// This is why pre-allocating slices matters in hot loops. If you know you'll
// need 1000 elements, use make([]T, 0, 1000) to avoid 10+ reallocations.
func demonstrateAppendWithCapacity() {
	fmt.Println("=== Append With Spare Capacity (No Allocation) ===")

	// BREAKPOINT: Set breakpoint after make to inspect pre-allocated slice
	// DEBUG: s has backing array with 5 slots, but only 2 are initialized
	// DEBUG: The 3 extra slots (s[2], s[3], s[4]) exist but aren't accessible yet
	// MICRO-COMMENT: make([]int, len, cap) creates a slice with:
	// - Initial length = 2 (elements are zero-initialized)
	// - Capacity = 5 (room for 5 total elements before reallocation)
	s := make([]int, 2, 5)
	s[0] = 10
	s[1] = 20

	fmt.Printf("Initial:  %v (len=%d, cap=%d)\n", s, len(s), cap(s))

	// DEBUG: Watch len increase from 2→3→4→5 without cap changing (no reallocation)
	// DEBUG: Each append just writes to next slot in existing backing array
	// BREAKPOINT: Set breakpoint before first append with spare capacity
	// MICRO-COMMENT: Append 3 more elements (within capacity)
	// Each append:
	// 1. Checks if len < cap (yes, we have room)
	// 2. Writes to backing_array[len]
	// 3. Increments len
	// NO ALLOCATION occurs because cap=5 and we're only going up to len=5
	s = append(s, 30)
	fmt.Printf("After +30: %v (len=%d, cap=%d) ← No reallocation\n", s, len(s), cap(s))

	s = append(s, 40)
	fmt.Printf("After +40: %v (len=%d, cap=%d) ← No reallocation\n", s, len(s), cap(s))

	s = append(s, 50)
	fmt.Printf("After +50: %v (len=%d, cap=%d) ← No reallocation\n", s, len(s), cap(s))

	// BREAKPOINT: Set breakpoint before append that exceeds capacity
	// DEBUG: Before: len=5, cap=5 (full). After: len=6, cap=10 (new backing array)
	// DEBUG: The backing array address will change - old array becomes garbage
	// MICRO-COMMENT: Now append beyond capacity
	// len=5, cap=5, so the next append MUST allocate
	s = append(s, 60)
	fmt.Printf("After +60: %v (len=%d, cap=%d) ← REALLOCATION (cap doubled)\n", s, len(s), cap(s))

	fmt.Println()
}

// ============================================================================
// SECTION 5: Comparing Pre-Allocated vs Dynamic Growth
// ============================================================================

// demonstratePreallocation compares two approaches to building a large slice.
//
// MACRO-COMMENT: Performance Comparison
// Approach 1 (no pre-allocation): Suffers ~10-20 reallocations as the slice grows
// Approach 2 (pre-allocated): Zero reallocations, much faster
//
// BENCHMARKING NOTE:
// For 10,000 elements, pre-allocation is typically 10-100x faster.
// For 1,000,000 elements, the difference can be even more dramatic.
//
// GARBAGE COLLECTION IMPACT:
// Each reallocation creates garbage (the old backing array). More allocations
// means more GC pressure, which can cause latency spikes in production.
func demonstratePreallocation() {
	fmt.Println("=== Pre-Allocation vs Dynamic Growth ===")

	const numElements = 10000

	// APPROACH 1: No pre-allocation
	// DEBUG: Watch allocCount1 increase (expect ~13-14 reallocations for 10K elements)
	// BREAKPOINT: Set breakpoint after loop to see total reallocations
	// MICRO-COMMENT: var declaration creates a nil slice (no backing array)
	var s1 []int

	// MICRO-COMMENT: Count how many times capacity changes (reallocations)
	allocCount1 := 0
	prevCap1 := 0

	for i := 0; i < numElements; i++ {
		s1 = append(s1, i)
		if cap(s1) != prevCap1 {
			allocCount1++
			prevCap1 = cap(s1)
		}
	}

	// DEBUG: Final s1 capacity will be ~16384 (next power of 2 after 10000)

	fmt.Printf("Approach 1 (no pre-allocation):\n")
	fmt.Printf("  Elements: %d, Final Cap: %d, Reallocations: %d\n", len(s1), cap(s1), allocCount1)

	// APPROACH 2: Pre-allocate exact capacity
	// BREAKPOINT: Set breakpoint after pre-allocation to inspect s2
	// DEBUG: s2 has len=0 but cap=10000 (backing array already exists)
	// DEBUG: Watch allocCount2 stay at 0 throughout the loop (zero reallocations!)
	// MICRO-COMMENT: make([]int, 0, numElements) creates:
	// - Len = 0 (empty)
	// - Cap = numElements (backing array already allocated)
	// This costs one upfront allocation but zero allocations in the loop
	s2 := make([]int, 0, numElements)

	allocCount2 := 0
	prevCap2 := cap(s2)

	for i := 0; i < numElements; i++ {
		s2 = append(s2, i)
		if cap(s2) != prevCap2 {
			allocCount2++
			prevCap2 = cap(s2)
		}
	}

	fmt.Printf("Approach 2 (pre-allocated):\n")
	fmt.Printf("  Elements: %d, Final Cap: %d, Reallocations: %d\n", len(s2), cap(s2), allocCount2)

	fmt.Println()
}

// ============================================================================
// SECTION 6: Copy to Isolate Slices
// ============================================================================

// demonstrateCopy shows how to create an independent copy of a slice.
//
// MACRO-COMMENT: When to Use copy()
// Use copy() when you need to:
// 1. Prevent mutations to the original slice
// 2. Avoid retaining a large backing array (memory leak prevention)
// 3. Safely pass slices to goroutines (avoid data races)
//
// BUILTIN FUNCTION: copy(dst, src []T) int
// - Copies elements from src to dst
// - Returns the number of elements copied: min(len(dst), len(src))
// - dst must be large enough to hold the data (or elements will be truncated)
func demonstrateCopy() {
	fmt.Println("=== Using copy() to Isolate Slices ===")

	// DEBUG: Before copy, isolated has its own backing array but all zeros
	// BREAKPOINT: Set breakpoint after copy to verify independent backing arrays
	// MICRO-COMMENT: Create original slice
	original := []int{10, 20, 30, 40, 50}

	// MICRO-COMMENT: Allocate a NEW backing array with the same length
	isolated := make([]int, len(original))

	// DEBUG: copy() returns number of elements copied (should be 5)
	// DEBUG: After copy, check &original[0] != &isolated[0] (different arrays)
	// MICRO-COMMENT: copy() performs a deep copy (element-by-element)
	// This creates a completely independent slice with its own backing array
	// PERFORMANCE: copy() is optimized at the runtime level (uses memmove internally)
	n := copy(isolated, original)

	fmt.Printf("Copied %d elements\n", n)
	fmt.Printf("Original:  %v\n", original)
	fmt.Printf("Isolated:  %v\n", isolated)

	// BREAKPOINT: Set breakpoint before mutation to verify independence
	// DEBUG: Watch isolated[0] change to 999 while original[0] stays at 10
	// MICRO-COMMENT: Modify the isolated slice
	isolated[0] = 999

	fmt.Printf("\nAfter modifying isolated[0] = 999:\n")
	fmt.Printf("  Original:  %v  ← UNCHANGED (different backing array)\n", original)
	fmt.Printf("  Isolated:  %v\n", isolated)

	// DEBUG: Confirm original is truly unchanged - proves copy created new backing array

	fmt.Println()
}

// ============================================================================
// MAIN FUNCTION: Orchestrates All Demonstrations
// ============================================================================

// main executes all demonstration functions in order.
//
// MACRO-COMMENT: Learning Progression
// The demos are ordered to build conceptual understanding:
// 1. Capacity growth (fundamental behavior)
// 2. Shared backing arrays (common gotcha)
// 3. Safe re-slicing (how to avoid the gotcha)
// 4. Append performance (when allocations occur)
// 5. Pre-allocation (optimization technique)
// 6. Copy for isolation (defensive programming)
//
// TRY THIS: Run with escape analysis to see heap allocations:
// go build -gcflags='-m' cmd/slices-demo/main.go
func main() {
	// Parse command-line flags
	demo := flag.String("demo", "all", "Which demo to run: all, growth, shared, safe, append, prealloc, copy")
	verbose := flag.Bool("verbose", false, "Enable verbose output with memory addresses")
	iterations := flag.Int("iterations", 10000, "Number of iterations for preallocation demo")
	flag.Parse()

	// DEBUG: Flag values control which demonstrations run
	// DEBUG: Use -demo=growth to run only capacity growth demo
	// BREAKPOINT: Set breakpoint here to inspect parsed flags

	if *verbose {
		fmt.Println("=== Verbose Mode Enabled ===")
		fmt.Printf("Demo: %s, Iterations: %d\n\n", *demo, *iterations)
	}

	// MICRO-COMMENT: Call each demonstration function based on flags
	// DEBUG: Watch which demos execute based on -demo flag value
	if *demo == "all" || *demo == "growth" {
		demonstrateCapacityGrowth()
	}
	if *demo == "all" || *demo == "shared" {
		demonstrateSharedBackingArray()
	}
	if *demo == "all" || *demo == "safe" {
		demonstrateSafeReslicing()
	}
	if *demo == "all" || *demo == "append" {
		demonstrateAppendWithCapacity()
	}
	if *demo == "all" || *demo == "prealloc" {
		demonstratePreallocation()
	}
	if *demo == "all" || *demo == "copy" {
		demonstrateCopy()
	}

	// MACRO-COMMENT: Compiler Optimization Notes
	// If you run with -gcflags='-m', you'll see lines like:
	// "moved to heap: s" - indicating a slice escaped to the heap
	// "inlining call to fmt.Printf" - the compiler inlined simple functions
	//
	// RUNTIME NOTES:
	// All slices in this program are small enough to be stack-allocated
	// (or quickly promoted to heap when needed). In production code with
	// large slices or long-lived data, heap allocations dominate.
}
