# 11: Slices Internals & Capacity Growth

This project is a deep dive into one of Go's most important data structures: the slice. While you've used slices in previous projects, this one focuses exclusively on their internal mechanics. You will learn how `append` works under the hood, how the Go runtime grows slice capacity, and what this means for writing high-performance code.

## Table of Contents

- [Learning Objectives](#learning-objectives)
- [The Big Picture: From User to Master of Slices](#the-big-picture-from-user-to-master-of-slices)
- [First Principles: The Slice Header and the Underlying Array](#first-principles-the-slice-header-and-the-underlying-array)
- [Project Structure](#project-structure)
- [The Core Concept: The Magic of `append`](#the-core-concept-the-magic-of-append)
  - [Scenario 1: `len < cap` (The Easy Path)](#scenario-1-len--cap-the-easy-path)
  - [Scenario 2: `len == cap` (The Reallocation Path)](#scenario-2-len--cap-the-reallocation-path)
  - [The Growth Algorithm](#the-growth-algorithm)
- [Performance Implication: The Cost of Ignorance](#performance-implication-the-cost-of-ignorance)
  - [The Solution: Pre-allocation](#the-solution-pre-allocation)
- [Progression: Performance Tuning](#progression-performance-tuning)
- [How to Run](#how-to-run)
- [Key Takeaways](#key-takeaways)
- [Further Reading](#further-reading)

## Learning Objectives

By the end of this project, you will be able to:

-   **Visualize the relationship between a slice and its underlying array**.
-   **Explain the difference between length and capacity** with confidence.
-   **Predict when `append` will cause a memory allocation**.
-   **Describe Go's slice capacity growth algorithm**.
-   **Write more performant Go code** by pre-allocating slices to avoid unnecessary reallocations.

## The Big Picture: From User to Master of Slices

Slices are designed to be easy to use, and for many cases, you don't need to think about how they work. However, to write truly high-performance Go, you must understand their internal mechanics. In performance-critical loops, every memory allocation matters. An `append` that unexpectedly allocates a new, large array and copies all the existing elements can be a significant performance bottleneck.

This project moves you from being a simple *user* of slices to being a *master* of them—one who can control memory allocations, predict performance, and write highly optimized code.

## First Principles: The Slice Header and the Underlying Array

A slice is **not** the data itself. A slice is a small, 3-word struct that describes a section of an underlying array.

```
type sliceHeader struct {
    Data uintptr // Pointer to the first element of the underlying array
    Len  int     // The number of elements the slice contains
    Cap  int     // The number of elements from Data to the end of the array
}
```
-   **The Underlying Array**: A contiguous block of memory holding the actual elements.
-   **The Slice**: A lightweight "view" or "window" into that array.

Multiple slices can share the same underlying array. This is what makes sub-slicing (`s[2:5]`) so efficient—it just creates a new slice header pointing to the same array but with different `Len`, `Cap`, and `Data` pointer offsets.

## Project Structure

```
.
└── cmd/
    └── dev/
        └── main.go       # A program that demonstrates slice growth.
```
- The `main.go` file contains a loop that appends elements to a slice one by one, printing the length and capacity at each step so you can observe the growth algorithm in action.

## The Core Concept: The Magic of `append`

The built-in `append` function hides a crucial decision-making process.

### Scenario 1: `len < cap` (The Easy Path)
If there is still room in the underlying array, no allocation happens.
1.  The new element is written to the memory location just after the last element of the slice.
2.  The `Len` of the slice is incremented by 1.
3.  The `Cap` and `Data` pointer remain unchanged.
This is a very fast, constant-time operation.

### Scenario 2: `len == cap` (The Reallocation Path)
If the underlying array is full, a reallocation is triggered. This is a multi-step, more expensive process.

**The process, visualized:**
```
Before:
slice -> [ Header: ptr | len: 4 | cap: 4 ]
           |
           v
           [ Underlying Array: cap=4 | 10 | 20 | 30 | 40 ]

append(slice, 50) // Oh no, cap is full!

After:
slice -> [ Header: ptr' | len: 5 | cap: 8 ]
           |
           v
           [ New, Larger Array: cap=8 | 10 | 20 | 30 | 40 | 50 | _ | _ ]

(The old array is now unreferenced and will be garbage collected)
```
1.  **Allocate**: A new, larger underlying array is allocated.
2.  **Copy**: All the elements from the old array are copied to the new array.
3.  **Append**: The new element is added to the end of the new array.
4.  **Update Header**: The slice's header is updated. Its `Data` pointer now points to the new array, and its `Len` and `Cap` are updated.

### The Growth Algorithm
How does Go decide the size of the new array? The goal is to balance memory usage with the number of reallocations.
-   If `append` has to reallocate, it doesn't just add one slot. It adds extra capacity to accommodate future appends.
-   **Prior to Go 1.18**: The rule was simple: if the old capacity was less than 1024, it would double the capacity. Otherwise, it would grow by a factor of 1.25 (25%).
-   **Go 1.18 and later**: The growth is more gradual. The doubling strategy is used for smaller capacities, but the growth factor smoothly decreases as the slice gets larger. This change was made to reduce wasted memory for very large slices. The exact formula is an implementation detail, but you will observe this behavior in the project.

## Performance Implication: The Cost of Ignorance

The reallocation path is much more expensive than the easy path. Copying all the elements can be a significant performance hit, especially for large slices. If you are appending in a tight loop, you might be triggering many reallocations without realizing it.

### The Solution: Pre-allocation
If you have a reasonable estimate of how many elements you're going to put in a slice, you can create it with the required capacity from the start using `make`.

```go
// BAD: Potentially many reallocations
var s []int
for i := 0; i < 1000; i++ {
    s = append(s, i)
}

// GOOD: Zero reallocations in the loop!
s := make([]int, 0, 1000) // len=0, cap=1000
for i := 0; i < 1000; i++ {
    s = append(s, i)
}
```
The second version will be dramatically faster because every `append` call takes the "Easy Path".

## Progression: Performance Tuning

This project provides the "why" behind the slice best practices you've been using. It builds on your initial understanding from **Project 02 (Arrays, Slices, and Maps)** and gives you the mental model needed for performance-critical work. This knowledge is essential for later projects that involve processing large amounts of data, such as I/O buffering or high-throughput network services.

## How to Run

```bash
# Run the program to see the len/cap growth in action
go run ./cmd/dev/main.go
```
Observe the output. Notice how the capacity stays the same for a while and then suddenly jumps. That jump is the reallocation. Pay attention to the size of the jump at different slice sizes.

## Key Takeaways

-   A slice is a 3-word header (`ptr`, `len`, `cap`) that points to an underlying array.
-   `append` is fast when capacity is available and slow when it must reallocate.
-   Reallocation involves creating a new, larger array and copying all the old elements.
-   Go's capacity growth algorithm is designed to be efficient but can be a performance bottleneck if not understood.
-   **When you know the approximate size of a slice, pre-allocate it with `make`** to achieve significant performance gains.

## Further Reading

-   [**Go Blog: Go Slices: usage and internals**](https://go.dev/blog/slices-intro) - The canonical blog post on slices.
-   [**Go by Example: Slices**](https://gobyexample.com/slices)
-   [**A Tour of Go: Slices**](https://go.dev/tour/moretypes/7)
-   [**Video: GopherCon 2016: Understanding nil**](https://www.youtube.com/watch?v=ynoY2xz-F8s) by Francesc Campoy - Includes a great section on `nil` slices and how they relate to length and capacity.