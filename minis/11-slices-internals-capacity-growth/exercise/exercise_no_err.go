//go:build !solution
// +build !solution

package exercise

/*
Core Slice Algorithms Without Complexity
=========================================

This file focuses on the core slice manipulation algorithms.
Understanding these patterns is essential for working with slices efficiently.

CORE CONCEPTS:
- Capacity tracking and growth patterns
- Memory sharing detection
- Independent slice creation
- Pre-allocation optimization

DEBUGGING WORKFLOW:
1. Set breakpoint at function entry
2. Step through algorithm line by line
3. Watch slice len/cap values transform
4. Observe backing array behavior

See /RUN_DEBUG.md for comprehensive debugging guide.
*/

// TODO: Implement GrowSlice
// BREAKPOINT: Set breakpoint at function start
// DEBUG: Watch oldCap = cap(s) before append
// DEBUG: Watch newCap = cap(newSlice) after append
// DEBUG: Compare oldCap vs newCap to see if reallocation occurred
//
// Algorithm:
// 1. Capture old capacity
// 2. Append element (may trigger reallocation)
// 3. Capture new capacity
// 4. Return all three values
//
// func GrowSlice(s []int, elem int) (newSlice []int, oldCap, newCap int) {
//     // BREAKPOINT: Capture state before append
//     // DEBUG: Watch cap(s) to see current capacity
//     oldCap = cap(s)
//
//     // BREAKPOINT: Trace append operation
//     // DEBUG: Watch if len(s) == cap(s) triggers reallocation
//     newSlice = append(s, elem)
//
//     // BREAKPOINT: Capture state after append
//     // DEBUG: Watch cap(newSlice) to see new capacity
//     newCap = cap(newSlice)
//
//     return newSlice, oldCap, newCap
// }

// TODO: Implement SharesBackingArray
// BREAKPOINT: Set breakpoint at function start
// DEBUG: Watch 'a' and 'b' slice states
// DEBUG: Watch sentinel value modification
// DEBUG: Watch restoration of original value
//
// Algorithm:
// 1. Check for empty slices (edge case)
// 2. Save original a[0] value
// 3. Set a[0] to sentinel value
// 4. Check if sentinel appears in b
// 5. Restore original a[0]
// 6. Return sharing result
//
// func SharesBackingArray(a, b []int) bool {
//     // BREAKPOINT: Check edge cases
//     // DEBUG: Watch len(a) and len(b)
//     if len(a) == 0 || len(b) == 0 {
//         return false
//     }
//
//     // BREAKPOINT: Save original value
//     // DEBUG: Watch original value being saved
//     original := a[0]
//     sentinel := original + 999999
//
//     // BREAKPOINT: Modify a[0]
//     // DEBUG: Watch a[0] change to sentinel
//     a[0] = sentinel
//
//     // BREAKPOINT: Check if b was affected
//     // DEBUG: Watch loop searching for sentinel in b
//     shares := false
//     for i := range b {
//         if b[i] == sentinel {
//             shares = true
//             break
//         }
//     }
//
//     // BREAKPOINT: Restore original
//     // DEBUG: Watch a[0] restored to original
//     a[0] = original
//
//     return shares
// }

// TODO: Implement SafeTruncate
// BREAKPOINT: Set breakpoint at function start
// DEBUG: Watch input slice 's' and truncation size 'n'
// DEBUG: Watch new slice allocation
// DEBUG: Watch copy operation
//
// Algorithm:
// 1. Clamp n to slice length if needed
// 2. Handle zero-length case
// 3. Allocate new independent slice
// 4. Copy first n elements
// 5. Return independent slice
//
// func SafeTruncate(s []int, n int) []int {
//     // BREAKPOINT: Validate bounds
//     // DEBUG: Watch n being clamped to len(s)
//     if n > len(s) {
//         n = len(s)
//     }
//     if n == 0 {
//         return []int{}
//     }
//
//     // BREAKPOINT: Allocate new slice
//     // DEBUG: Watch make() create independent array
//     result := make([]int, n)
//
//     // BREAKPOINT: Copy elements
//     // DEBUG: Watch copy() transfer data
//     copy(result, s[:n])
//
//     return result
// }

// TODO: Implement PreallocateVsDynamic
// BREAKPOINT: Set breakpoint at function start
// DEBUG: Watch dynamic growth allocations
// DEBUG: Watch pre-allocated behavior (should be zero allocations)
//
// Algorithm:
// 1. Test dynamic growth: append to nil slice n times
// 2. Count capacity changes (reallocations)
// 3. Test pre-allocated: make([]int, 0, n) then append n times
// 4. Count capacity changes (should be zero)
// 5. Return both counts
//
// func PreallocateVsDynamic(n int) (dynamicAllocs, preallocAllocs int) {
//     // BREAKPOINT: Dynamic growth test
//     // DEBUG: Watch s1 start as nil
//     var s1 []int
//     prevCap1 := 0
//
//     for i := 0; i < n; i++ {
//         s1 = append(s1, i)
//         // BREAKPOINT: Check for capacity change
//         // DEBUG: Watch cap(s1) change and count reallocations
//         if cap(s1) != prevCap1 {
//             dynamicAllocs++
//             prevCap1 = cap(s1)
//         }
//     }
//
//     // BREAKPOINT: Pre-allocated test
//     // DEBUG: Watch s2 start with full capacity
//     s2 := make([]int, 0, n)
//     prevCap2 := cap(s2)
//
//     for i := 0; i < n; i++ {
//         s2 = append(s2, i)
//         // BREAKPOINT: Verify no capacity change
//         // DEBUG: This should never trigger
//         if cap(s2) != prevCap2 {
//             preallocAllocs++
//             prevCap2 = cap(s2)
//         }
//     }
//
//     return dynamicAllocs, preallocAllocs
// }

// TODO: Implement ReSliceWithCapLimit
// BREAKPOINT: Set breakpoint at function start
// DEBUG: Watch 3-index slice syntax
// DEBUG: Watch result has cap == len (no spare capacity)
//
// Algorithm:
// 1. Use 3-index slice: s[start:end:end]
// 2. First end sets length
// 3. Second end sets capacity limit
// 4. Result cannot grow without reallocation
//
// func ReSliceWithCapLimit(s []int, start, end int) []int {
//     // BREAKPOINT: Create capacity-limited slice
//     // DEBUG: Watch s[start:end:end] create slice with cap == len
//     return s[start:end:end]
// }
