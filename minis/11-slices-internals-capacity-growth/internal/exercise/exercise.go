//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "unsafe" for comparing memory addresses (advanced technique)
//
// import (
//     "unsafe"
// )

// SLICE INTERNALS FUNDAMENTALS
// =============================
// A slice in Go is actually a struct with 3 fields:
//   type SliceHeader struct {
//       Data uintptr  // Pointer to the first element of the underlying array
//       Len  int      // Number of elements in the slice
//       Cap  int      // Capacity of the underlying array (from Data pointer)
//   }
//
// Key concepts:
// - Slices are REFERENCE types (they reference an underlying array)
// - When you pass a slice to a function, you copy the header (24 bytes on 64-bit)
//   but the Data pointer still points to the SAME underlying array
// - Multiple slices can share the same backing array
// - append() may or may not reallocate depending on capacity
// - Reallocation means: allocate new array, copy elements, return new slice

// GrowSlice appends an element to a slice and returns information about
// whether a reallocation occurred.
//
// TODO: Implement GrowSlice function
// Function signature: func GrowSlice(s []int, elem int) (newSlice []int, oldCap, newCap int)
//
// Steps to implement:
//
// 1. Capture the old capacity BEFORE appending
//    - Use: oldCap = cap(s)
//    - cap() returns the capacity of the underlying array
//    - This is the maximum number of elements before reallocation
//
// 2. Append the element
//    - Use: newSlice = append(s, elem)
//    - append() behavior depends on capacity:
//      * If len(s) < cap(s): Just increments length, returns same backing array
//      * If len(s) == cap(s): Allocates NEW array (usually 2x capacity), copies all elements
//    - CRITICAL: append() returns a NEW slice header (even if array didn't change)
//    - You MUST assign the result: newSlice = append(s, elem)
//    - Common mistake: append(s, elem) without assignment (loses the result!)
//
// 3. Capture the new capacity AFTER appending
//    - Use: newCap = cap(newSlice)
//    - If newCap > oldCap, a reallocation occurred
//    - If newCap == oldCap, no reallocation (used spare capacity)
//
// 4. Return all three values
//    - Use: return newSlice, oldCap, newCap
//
// Key Go concepts:
// - append() ALWAYS returns a slice (never modifies the slice header in place)
// - The slice header is passed BY VALUE (copied into the function)
// - The underlying array is shared UNTIL reallocation
// - After reallocation, old and new slices have DIFFERENT backing arrays
//
// Growth pattern (Go 1.18+):
// - For capacity < 256: new_cap = old_cap * 2
// - For capacity >= 256: new_cap = old_cap + old_cap/4 + 192 (approximately 1.25x)
//
// Memory behavior:
// - Before append: s points to array A
// - If no reallocation: newSlice also points to array A (SHARED!)
// - If reallocation: newSlice points to NEW array B, s still points to A
//
// TODO: Implement the GrowSlice function below
// func GrowSlice(s []int, elem int) (newSlice []int, oldCap, newCap int) {
//     return nil, 0, 0
// }

// SharesBackingArray returns true if slices a and b share the same backing array.
//
// TODO: Implement SharesBackingArray function
// Function signature: func SharesBackingArray(a, b []int) bool
//
// Two approaches:
//
// APPROACH 1: Unsafe pointer comparison (what the current implementation uses)
// - Use unsafe.Pointer to get the address of the first element
// - Compare memory ranges to detect overlap
// - More accurate but requires understanding of unsafe package
//
// APPROACH 2: Simple modification test (easier to understand)
// - Save original a[0]
// - Set a[0] to a sentinel value (like 999999)
// - Check if any element in b has the sentinel value
// - Restore a[0]
// - Returns true if modification was visible in b
//
// Steps for APPROACH 2 (simpler, recommended for learning):
//
// 1. Handle edge cases
//    - If len(a) == 0 || len(b) == 0, return false
//    - Can't test for sharing if there are no elements
//
// 2. Save the original value
//    - Use: original := a[0]
//    - Store this so we can restore it later
//
// 3. Modify a[0] to a sentinel value
//    - Use: a[0] = original + 999999
//    - Choose a value unlikely to already exist
//    - This modification will be visible in b IF they share the array
//
// 4. Check if b was affected
//    - Loop through b: for i := range b { ... }
//    - If b[i] == sentinel, they share the array
//    - Set shares = true and break
//
// 5. Restore the original value
//    - Use: a[0] = original
//    - CRITICAL: Always restore, even if we found sharing
//    - Otherwise we've modified the caller's data!
//
// 6. Return the result
//    - Use: return shares
//
// Steps for APPROACH 1 (current implementation, uses unsafe):
//
// 1. Handle edge cases (same as above)
//
// 2. Get memory addresses
//    - Use: unsafe.Pointer(&a[0]) to get address of first element
//    - Convert to uintptr for arithmetic
//    - Calculate start and end addresses for both slices
//
// 3. Check for overlap
//    - Two ranges overlap if:
//      * startA is within [startB, endB), OR
//      * startB is within [startA, endA)
//    - Return true if they overlap
//
// Key Go concepts:
// - Slicing (a[1:3]) creates a NEW slice header pointing to the SAME array
// - The new slice's Data pointer points somewhere inside the original array
// - append() that reallocates creates a DIFFERENT backing array
// - copy() always creates a NEW array
//
// Memory sharing examples:
//   a := []int{1, 2, 3, 4, 5}
//   b := a[1:4]               // b shares a's array
//   c := make([]int, len(a))
//   copy(c, a)                // c has its own array
//   d := a                    // d shares a's array (same slice header copy)
//   e := append(a, 6)         // e might share (if cap(a) > len(a)) or might not
//
// TODO: Implement the SharesBackingArray function below
// You can use either approach, but approach 2 is simpler for learning
// The current implementation uses approach 1 with unsafe pointers
// func SharesBackingArray(a, b []int) bool {
//     return false
// }

// SafeTruncate truncates a large slice to a smaller size, ensuring the
// large backing array can be garbage collected.
//
// TODO: Implement SafeTruncate function
// Function signature: func SafeTruncate(s []int, n int) []int
//
// THE PROBLEM:
// When you slice a large array, the entire array stays in memory:
//   huge := make([]int, 1000000)  // 8MB array
//   tiny := huge[:10]              // Only need 10 elements
//   huge = nil                     // Want to free memory...
//   // But the 8MB array is STILL in memory because tiny references it!
//
// THE SOLUTION:
// Create a NEW slice with its own backing array containing only what you need
//
// Steps to implement:
//
// 1. Handle edge case: n > len(s)
//    - Use: if n > len(s) { n = len(s) }
//    - Can't truncate to more elements than exist
//
// 2. Handle edge case: n == 0
//    - Use: if n == 0 { return []int{} }
//    - Return an empty slice (nil would also work)
//
// 3. Allocate a NEW array with the exact size needed
//    - Use: result := make([]int, n)
//    - make([]int, n) creates:
//      * A new array with capacity n
//      * A slice with length n pointing to it
//    - This array is INDEPENDENT (not shared with s)
//
// 4. Copy the first n elements from s
//    - Use: copy(result, s[:n])
//    - copy() built-in function copies elements between slices
//    - It copies min(len(dst), len(src)) elements
//    - s[:n] is a slice expression (still references s's array)
//    - result now contains the data, but in its own array
//
// 5. Return the new slice
//    - Use: return result
//    - The caller's original slice s is unchanged
//    - But result has no connection to s's backing array
//    - When s goes out of scope, its large array can be GC'd
//
// Key Go concepts:
// - Slicing (s[:n]) SHARES the backing array
// - make() ALLOCATES a new backing array
// - copy() COPIES elements (doesn't share memory)
// - GC can only free memory when there are NO references to it
//
// Memory lifecycle:
// Before SafeTruncate:
//   s: [large array with 1M elements] <- s's Data pointer
//
// After SafeTruncate:
//   s: [large array with 1M elements] <- s's Data pointer (still)
//   result: [small array with n elements] <- result's Data pointer (NEW!)
//
// After caller does "s = nil":
//   s: nil (no reference to large array)
//   result: [small array with n elements]
//   [large array] <- No more references! GC can free this!
//
// TODO: Implement the SafeTruncate function below
// func SafeTruncate(s []int, n int) []int {
//     return nil
// }

// PreallocateVsDynamic compares pre-allocated vs dynamic growth for building
// a slice of n elements.
//
// TODO: Implement PreallocateVsDynamic function
// Function signature: func PreallocateVsDynamic(n int) (dynamicAllocs, preallocAllocs int)
//
// This function demonstrates the performance difference between:
// 1. Dynamic growth: Start with nil slice, let append() handle growth
// 2. Pre-allocation: Allocate the exact capacity needed upfront
//
// Steps to implement:
//
// PART 1: Dynamic growth (multiple reallocations expected)
//
// 1. Create an empty slice
//    - Use: var s1 []int
//    - This is a nil slice (Data=nil, Len=0, Cap=0)
//    - First append will allocate initial array
//
// 2. Track the previous capacity
//    - Use: prevCap1 := 0
//    - We'll compare cap(s1) to this after each append
//
// 3. Loop n times, appending elements
//    - Use: for i := 0; i < n; i++ { ... }
//    - Append: s1 = append(s1, i)
//    - After each append, check if capacity changed
//
// 4. Count reallocations
//    - After append: if cap(s1) != prevCap1 { ... }
//    - If capacity changed, increment dynamicAllocs counter
//    - Update: prevCap1 = cap(s1)
//
// PART 2: Pre-allocated growth (zero reallocations expected)
//
// 5. Create a pre-allocated slice
//    - Use: s2 := make([]int, 0, n)
//    - This allocates an array with capacity n
//    - Length is 0 (empty), but capacity is n
//    - First n appends will NOT reallocate
//
// 6. Track the previous capacity
//    - Use: prevCap2 := cap(s2)
//    - This should be n initially
//
// 7. Loop n times, appending elements
//    - Use: for i := 0; i < n; i++ { ... }
//    - Append: s2 = append(s2, i)
//    - After each append, check if capacity changed
//
// 8. Count reallocations
//    - After append: if cap(s2) != prevCap2 { ... }
//    - This should NEVER trigger (capacity doesn't change)
//    - Increment preallocAllocs counter if it does
//    - Update: prevCap2 = cap(s2)
//
// 9. Return both counts
//    - Use: return dynamicAllocs, preallocAllocs
//
// Expected results for n=10000:
// - dynamicAllocs: ~15-20 (depends on Go version and growth algorithm)
// - preallocAllocs: 0 (no reallocations needed!)
//
// Key Go concepts:
// - append() reallocates when len == cap
// - Growth is usually 2x for small slices, ~1.25x for large slices
// - Pre-allocation avoids ALL reallocations
// - Each reallocation involves: allocate new array, copy all elements, free old array
//
// Performance implications:
// - Dynamic: O(n) total allocations, O(n log n) total copies
// - Pre-allocated: O(1) allocations, O(0) copies
// - For large n, pre-allocation is MUCH faster
//
// When to pre-allocate:
// - When you know the final size upfront
// - When building large slices (> 1000 elements)
// - When performance is critical
//
// When NOT to pre-allocate:
// - When size is unknown
// - When slice will likely be small
// - When memory is constrained (avoid over-allocation)
//
// TODO: Implement the PreallocateVsDynamic function below
// func PreallocateVsDynamic(n int) (dynamicAllocs, preallocAllocs int) {
//     return 0, 0
// }

// ReSliceWithCapLimit creates a sub-slice with limited capacity using
// the 3-index slice expression.
//
// TODO: Implement ReSliceWithCapLimit function
// Function signature: func ReSliceWithCapLimit(s []int, start, end int) []int
//
// BACKGROUND: The 3-index slice expression
// Normal slicing: s[low:high]
//   - Creates a slice from index low to high (exclusive)
//   - Length: high - low
//   - Capacity: cap(s) - low (all the way to the end of the backing array)
//
// 3-index slicing: s[low:high:max]
//   - Creates a slice from index low to high (exclusive)
//   - Length: high - low (same as normal slicing)
//   - Capacity: max - low (LIMITS the capacity!)
//   - This prevents accessing elements beyond index max
//
// WHY LIMIT CAPACITY?
// - Prevents accidental modification of array elements beyond the slice
// - Forces reallocation on append, creating independent slices
// - Useful for passing sub-slices to untrusted code
//
// Steps to implement:
//
// 1. Use the 3-index slice syntax
//    - Use: return s[start:end:end]
//    - start: First index to include
//    - end: First index to exclude (length = end - start)
//    - end (third index): Capacity limit (cap = end - start)
//
// 2. Why set max = end?
//    - This makes capacity equal to length
//    - cap(result) = end - start
//    - len(result) = end - start
//    - No spare capacity! Any append MUST reallocate
//
// Example:
//   s := []int{10, 20, 30, 40, 50}
//   // len(s)=5, cap(s)=5
//
//   // Normal 2-index slicing:
//   sub1 := s[1:3]
//   // sub1 = [20, 30], len=2, cap=4 (can access up to s[4])
//   sub1 = append(sub1, 35)
//   // Now s = [10, 20, 30, 35, 50] <- Oops! Modified s[3]!
//
//   // 3-index slicing with cap limit:
//   sub2 := s[1:3:3]
//   // sub2 = [20, 30], len=2, cap=2 (cannot access beyond s[2])
//   sub2 = append(sub2, 35)
//   // Must reallocate! s is unchanged: [10, 20, 30, 40, 50]
//
// Key Go concepts:
// - Normal slicing shares capacity with the original slice
// - 3-index slicing limits how much of the backing array is accessible
// - Setting max = high makes len == cap (full capacity)
// - This forces append() to always reallocate
// - Useful for creating independent sub-slices
//
// Memory safety:
// - Prevents buffer overrun via append
// - Ensures append creates a new array
// - Original slice remains unchanged after append on sub-slice
//
// TODO: Implement the ReSliceWithCapLimit function below
// func ReSliceWithCapLimit(s []int, start, end int) []int {
//     return nil
// }

// After implementing all functions:
// - Run: go test ./minis/11-slices-internals-capacity-growth/exercise/...
// - Check: go test -v for verbose output
// - Experiment with different slice sizes to see growth patterns
// - Compare with solution.go to see detailed explanations
