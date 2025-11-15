# Project 11: Slices Internals, Capacity Growth & Memory Behavior

## What Is This Project About?

This project takes you **deep into the machinery** of Go slices—one of the most important and misunderstood data structures in the language. You'll learn:

1. **What a slice REALLY is** (slice header, backing array, pointer mechanics)
2. **How capacity doubling works** (and when it doesn't)
3. **Re-slicing gotchas** (shared backing arrays, unexpected mutations)
4. **Escape analysis** (when slices stay on the stack vs heap)
5. **Performance implications** (allocations, copying, memory waste)

By the end, you'll understand exactly what happens in memory when you write `append()`, why `cap()` matters, and how to avoid common slice pitfalls that cause production bugs.

---

## The Fundamental Problem: Arrays Are Too Rigid

### First Principles: What Is an Array?

An **array** in Go is a **fixed-size** sequence of elements stored contiguously in memory:

```go
var arr [5]int  // Exactly 5 integers, no more, no less
arr[0] = 10
arr[1] = 20
// arr[5] = 30  // COMPILE ERROR: index out of bounds
```

**Memory layout:**
```
arr: [10][20][0][0][0]
     └─────────────────┘
     Exactly 5 slots, fixed forever
```

**The problem:**
- You must know the size at compile time (or declaration time)
- Cannot grow or shrink dynamically
- Passing to functions copies the entire array (expensive for large arrays)

### Real-World Scenario

Imagine building a web server that collects user IDs from requests. You don't know how many users will hit the endpoint:

```go
// BAD: What size do we choose?
var userIDs [100]int  // Too small? Too big? Unknown!

// If we get 101 users, the program crashes
// If we get 10 users, we waste 90 slots
```

This is why **slices** exist.

---

## What Is a Slice? (The Core Concept)

A **slice** is a **descriptor** (a small struct) that points to an **underlying array**. It's like a "window" into an array that can move and resize.

### The Slice Header (3 Fields)

Every slice is represented internally as a struct with **exactly 3 fields**:

```go
// This is what a slice REALLY is (simplified):
type SliceHeader struct {
    Ptr *ElementType  // Pointer to the first element in the backing array
    Len int           // Number of elements currently in the slice
    Cap int           // Maximum capacity before reallocation is needed
}
```

**Example:**
```go
s := []int{10, 20, 30}
```

**Memory layout:**
```
Slice Header (on stack):
┌──────────────┬─────┬─────┐
│ Ptr          │ Len │ Cap │
│ 0x400000     │  3  │  3  │
└──────────────┴─────┴─────┘
       │
       └──> Backing Array (on heap):
            ┌────┬────┬────┐
            │ 10 │ 20 │ 30 │
            └────┴────┴────┘
            Address: 0x400000
```

**Key insight:** The slice itself (the header) is **small and cheap to copy** (just 24 bytes on 64-bit systems). The actual data lives in the backing array.

---

## Capacity vs Length: The Critical Distinction

- **Length (`len(s)`)**: How many elements are **currently** in the slice
- **Capacity (`cap(s)`)**: How many elements the backing array **can hold** before reallocation

**Example:**
```go
s := make([]int, 3, 10)
// Len = 3 (three elements initialized to zero)
// Cap = 10 (room for 10 total before growing)

fmt.Println(len(s))  // 3
fmt.Println(cap(s))  // 10
```

**Memory layout:**
```
Slice Header:
┌──────────────┬─────┬─────┐
│ Ptr          │ Len │ Cap │
│ 0x400000     │  3  │ 10  │
└──────────────┴─────┴─────┘
       │
       └──> Backing Array (capacity 10):
            ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
            │ 0 │ 0 │ 0 │ ? │ ? │ ? │ ? │ ? │ ? │ ? │
            └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
            └───────┘ └───────────────────────────┘
              Used         Unused (reserved space)
            (len=3)          (cap-len = 7)
```

The first 3 elements are initialized to `0` (the zero value for `int`). The remaining 7 slots are **reserved but uninitialized**.

---

## The Append Operation: What Really Happens

### Case 1: Append with Spare Capacity

```go
s := make([]int, 2, 5)  // len=2, cap=5
s[0] = 10
s[1] = 20

s = append(s, 30)  // Add a third element
```

**Before append:**
```
Backing array: [10][20][?][?][?]
               └───────┘
                len=2
               └──────────────┘
                    cap=5
```

**After append:**
```
Backing array: [10][20][30][?][?]
               └──────────┘
                 len=3
               └──────────────┘
                    cap=5
```

**What happened:**
1. Go checks: `len (2) < cap (5)` → spare capacity exists
2. Writes `30` to `array[2]` (the next available slot)
3. Updates the slice header: `Len = 3` (Cap stays 5)
4. **No allocation occurred** (fast!)

### Case 2: Append Exceeding Capacity

```go
s := make([]int, 2, 2)  // len=2, cap=2 (FULL)
s[0] = 10
s[1] = 20

s = append(s, 30)  // Need to grow!
```

**What happens:**
1. Go checks: `len (2) == cap (2)` → no spare capacity!
2. **Allocates a NEW, LARGER backing array** (typically double the capacity)
3. **Copies all existing elements** to the new array
4. Appends the new element
5. Returns a NEW slice header pointing to the new array

**Before append:**
```
Old array: [10][20]  (cap=2, full)
```

**After append:**
```
New array: [10][20][30][?]  (cap=4, len=3)
Old array: [10][20]          (orphaned, will be garbage collected)
```

**Critical insight:** The old slice `s` now points to a **completely different array**. This is why `append()` **returns a new slice**—you must assign it back!

```go
// WRONG:
s := []int{1, 2}
append(s, 3)  // BUG: return value ignored, s is unchanged!

// RIGHT:
s = append(s, 3)  // Reassign to capture the potentially new slice
```

---

## Capacity Growth Strategy

Go doesn't grow capacity by a fixed amount. It uses a **smart doubling strategy** (with modifications for large slices).

### The Growth Rules (Simplified)

1. **If current capacity < 256 elements**: Double the capacity
2. **If current capacity >= 256 elements**: Grow by 1.25x plus a bit more

**Why this strategy?**
- **Doubling** (small slices): Minimizes allocations for rapidly growing data
- **Slower growth** (large slices): Prevents wasting too much memory

**Example growth sequence starting from cap=1:**
```
Cap: 1 → 2 → 4 → 8 → 16 → 32 → 64 → 128 → 256 → 320 → 400 → ...
     ×2  ×2  ×2  ×2   ×2   ×2   ×2    ×2   ×1.25 ×1.25 ...
```

### Measuring Capacity Growth

```go
s := []int{}
for i := 0; i < 1000; i++ {
    prevCap := cap(s)
    s = append(s, i)
    if cap(s) != prevCap {
        fmt.Printf("Len: %4d, Cap grew: %4d → %4d\n", len(s), prevCap, cap(s))
    }
}
```

**Output (excerpt):**
```
Len:    1, Cap grew:    0 →    1
Len:    2, Cap grew:    1 →    2
Len:    3, Cap grew:    2 →    4
Len:    5, Cap grew:    4 →    8
Len:    9, Cap grew:    8 →   16
Len:   17, Cap grew:   16 →   32
...
```

---

## Re-Slicing: Sharing Backing Arrays (The Gotcha)

When you create a slice from another slice, they **share the same backing array**. This leads to surprising mutations.

### Example 1: Unexpected Mutation

```go
original := []int{10, 20, 30, 40, 50}
slice1 := original[1:4]  // [20, 30, 40]

// Modify slice1
slice1[0] = 999

fmt.Println(original)  // [10, 999, 30, 40, 50]  ← CHANGED!
fmt.Println(slice1)    // [999, 30, 40]
```

**Why?**
Both slices point to the **same backing array**:

```
original: [10][20][30][40][50]
           └──┼──┼──┼──┘
slice1:       └──┼──┼──┘
              [20][30][40]

After slice1[0] = 999:
           [10][999][30][40][50]
```

### Example 2: Append and Re-Slice Interaction

```go
a := []int{1, 2, 3}
b := a[0:2]        // [1, 2] (shares backing array with a)

b = append(b, 99)  // What happens to a?

fmt.Println(a)     // [1, 2, 99]  ← a was modified!
fmt.Println(b)     // [1, 2, 99]
```

**Why?**
`b` has `len=2, cap=3` (the capacity comes from the original `a`). When appending `99`, there's spare capacity, so it writes to `array[2]`, which is **also part of `a`**!

### How to Avoid This (Use Full Slice Expressions)

Go provides the **3-index slice** syntax to **limit capacity**:

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3:3]  // [low:high:max]
               // Elements: a[1], a[2]
               // Cap: max - low = 3 - 1 = 2

fmt.Println(len(b))  // 2
fmt.Println(cap(b))  // 2 (not 4!)

b = append(b, 99)    // Forces a new allocation (cap exceeded)
fmt.Println(a)       // [1, 2, 3, 4, 5]  ← UNCHANGED!
```

---

## Escape Analysis: Stack vs Heap Allocation

Go's compiler decides whether a slice's backing array lives on the **stack** (fast, automatically cleaned up) or **heap** (slower, requires GC).

### When Does a Slice Escape to the Heap?

1. **Returned from a function** (outlives the function's stack frame)
2. **Stored in a struct that escapes**
3. **Passed to an interface** (size unknown at compile time)
4. **Too large** (slices > ~64KB often escape)

### Example: Analyzing Escape Behavior

```go
// ESCAPES: Returned to caller
func makeSlice() []int {
    s := []int{1, 2, 3}
    return s  // Escapes to heap (caller needs it)
}

// DOES NOT ESCAPE: Stays in function
func sumSlice() int {
    s := []int{1, 2, 3}
    sum := 0
    for _, v := range s {
        sum += v
    }
    return sum  // Only sum escapes (int), not slice
}
```

**How to check (using compiler flags):**
```bash
go build -gcflags='-m' main.go
```

**Output:**
```
./main.go:3:6: []int{...} escapes to heap
./main.go:8:6: []int{...} does not escape
```

### Why This Matters

- **Heap allocations** trigger garbage collection (GC pauses)
- **Stack allocations** are nearly free (just move a stack pointer)
- Understanding escape behavior helps optimize hot code paths

---

## Common Pitfalls and How to Avoid Them

### Pitfall 1: Ignoring append's Return Value

```go
// WRONG:
s := []int{1, 2}
append(s, 3)  // BUG: s is still [1, 2]

// RIGHT:
s = append(s, 3)  // Reassign to capture new slice
```

### Pitfall 2: Assuming Capacity After make()

```go
s := make([]int, 5)  // len=5, cap=5
s = append(s, 99)    // len=6 (new element is at index 5, not 0!)

fmt.Println(s)  // [0, 0, 0, 0, 0, 99] ← Surprise!
```

**Fix:** Use `make([]int, 0, 5)` if you want zero length but reserved capacity.

### Pitfall 3: Shared Backing Arrays

```go
a := []int{1, 2, 3}
b := a[:]           // Shares backing array
b[0] = 999
fmt.Println(a)      // [999, 2, 3] ← Modified!
```

**Fix:** Use `copy()` to create an independent slice:
```go
b := make([]int, len(a))
copy(b, a)
b[0] = 999
fmt.Println(a)  // [1, 2, 3] ← Unchanged
```

### Pitfall 4: Retaining Large Arrays via Small Slices

```go
// Load a huge file into memory
bigData := loadHugeFile()  // 1GB slice

// Extract just one line
oneLine := bigData[0:100]

// BUG: The entire 1GB backing array is kept alive in memory!
// Even though oneLine only references 100 bytes
```

**Fix:** Copy the small portion to a new slice:
```go
oneLine := make([]byte, 100)
copy(oneLine, bigData[0:100])
bigData = nil  // Allow GC to reclaim the 1GB
```

---

## Performance Implications

### Pre-Allocating vs Growing

**Scenario:** Building a slice of 10,000 elements.

**Approach 1: No pre-allocation**
```go
var s []int
for i := 0; i < 10000; i++ {
    s = append(s, i)  // Multiple reallocations!
}
```

**Approach 2: Pre-allocate**
```go
s := make([]int, 0, 10000)
for i := 0; i < 10000; i++ {
    s = append(s, i)  // No reallocations!
}
```

**Performance difference:** Approach 2 is **10-100x faster** (no allocations in the loop).

### Memory Waste from Over-Capacity

```go
s := make([]int, 10, 1000)  // len=10, cap=1000
// You're using 10 ints but reserving space for 1000
// Memory waste: 990 * 8 bytes = 7920 bytes
```

**Guideline:** Only pre-allocate if you have a good estimate of the final size.

---

## How to Run

```bash
# Run the demonstration program
cd minis/11-slices-internals-capacity-growth
go run cmd/slices-demo/main.go

# Run tests
cd exercise
go test -v

# Check escape analysis
go build -gcflags='-m' cmd/slices-demo/main.go
```

---

## Expected Output (Demo Program)

```
=== Capacity Growth Demonstration ===
Len:    1, Cap:    1 (grew from    0)
Len:    2, Cap:    2 (grew from    1)
Len:    3, Cap:    4 (grew from    2)
Len:    5, Cap:    8 (grew from    4)
...

=== Re-Slicing Gotcha ===
Original: [10 999 30 40 50]  ← Modified!
Slice1:   [999 30 40]

=== Safe Re-Slicing (3-index) ===
Original: [10 20 30 40 50]   ← Unchanged
Slice2:   [20 30 99]
```

---

## Key Takeaways

1. **Slices are headers** (pointer + len + cap), not the actual array
2. **Append can reallocate**, always assign the result: `s = append(s, x)`
3. **Capacity grows intelligently** (doubles for small slices, slower for large)
4. **Re-slicing shares backing arrays** (use 3-index slices or `copy()` to isolate)
5. **Pre-allocate if you know the size** (avoids repeated allocations)
6. **Escape analysis matters** (stack slices are much faster than heap slices)
7. **Small slices can retain large arrays** (use `copy()` to avoid memory leaks)

---

## Connections to Other Projects

- **Project 02 (arrays-maps-basics)**: You learned basic array syntax; now you see why slices are better
- **Project 06 (worker-pool-wordcount)**: Used slices for buffering; now you know the allocation costs
- **Project 17 (file-streaming-bufio)**: Will use pre-allocated byte slices to minimize allocations
- **Project 28 (pprof-cpu-mem-benchmarks)**: You'll profile slice allocations to find bottlenecks
- **Project 29 (escape-analysis-inlining)**: Deep dive into compiler optimizations for slices

---

## Stretch Goals

1. **Implement a custom `append()`** that mimics Go's growth strategy
2. **Benchmark** pre-allocated vs dynamic growth for 1M elements
3. **Write a function** that detects if two slices share a backing array
4. **Visualize** slice memory layout using `unsafe.Pointer` and reflection
