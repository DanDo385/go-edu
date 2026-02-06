# 02: Arrays Maps Basics

## Core Concepts

- The concrete problem in Arrays Maps Basics and the correctness invariants it depends on.
- How values, pointers, slices, maps, or channels behave in this module's runtime path.
- Why this lesson's implementation pattern is the right next step in the learning arc.

## CS Connection

- Data representation and state transitions: what is copied, what is shared, and what can race.
- API boundaries: where we validate, where we propagate errors, and where we normalize output.
- Algorithmic tradeoffs in this lesson (latency, throughput, memory, and complexity).

## End-State Understanding

- Diagnose and implement Arrays Maps Basics patterns without relying on hidden framework behavior.
- Explain memory and concurrency implications of the final implementation choices.
- Compare learner code and reference code by invariants, not only by syntax.

## Why This Lesson Exists Here

Problem statement:
This lesson turns the previous module's concepts into a reusable engineering pattern for arrays maps basics.

At this point in the arc:
Lesson 02 introduces a sharper systems concern so later modules can assume this mental model is stable.

## Step-by-Step Build Path

### Step 1: Problem This Step Solves
Define the smallest valid behavior and reject invalid input or impossible state early.

### Step 2: Why This Approach
Pick a direct design that keeps control flow and data flow visible for debugging and testing.

### Step 3: Memory / Pointer Impact
Call out where data is copied versus aliased, and where mutable shared state needs synchronization.

### Step 4: What Changed
Produce a stable result shape and explicit error behavior that downstream code can rely on.

## Pointer and Indirection

- Explain * and & in this module when they appear in code or docs.
- Show memory-before and memory-after when data ownership changes.
- Clarify common misconceptions: Go stays pass-by-value even when pointer values are copied.
- Primer link: docs/MEMORY_POINTERS_PRIMER.md

## Verify


a) learner path


go test -v ./...


b) reference path


go test -tags=reference -v ./...


This lesson covers Go's three fundamental collection types. You'll learn how to organize groups of data, which is a requirement for any non-trivial program. This lesson builds on your knowledge of string manipulation from the previous lesson.

## CS Connection

- Memory layout drives behavior: variables store values, and some values are addresses into other storage.
- Go is pass-by-value, including pointers; copying a pointer value copies an address, not the pointee.
- Correctness depends on understanding copying versus aliasing (`*T`, slices, maps, and channels) and enforcing synchronization when concurrent access exists.

## End-State Understanding

- Implement the exercise with explicit reasoning about correctness, edge cases, and error paths.

## What You'll Learn

- The difference between **arrays**, **slices**, and **maps**.
- The concepts of slice **length** and **capacity**.
- How to use maps for key-value associations, like frequency counters.
- The "comma ok" idiom for checking if a key exists in a map.

## Core Concepts

### Arrays vs. Slices

*   **Array**: A fixed-size collection of elements of the same type. In Go, arrays are *value types*, meaning they are copied when passed to a function. They are not used very often.
    ```go
    var a [3]int // An array of 3 integers. The size is part of the type.
    ```
*   **Slice**: A flexible, "view" into an underlying array. Slices are *reference types* (in spirit). When you pass a slice, you are copying a small header, but the header *points* to the same underlying data. This makes them efficient to pass around. They are the most common collection type in Go.
    ```go
    var s []int // A slice of integers. The size is not part of the type.
    ```

### The Slice Header

A slice is a small struct that contains three pieces of information:
1.  A **pointer** to an underlying array.
2.  The **length** of the slice (`len`).
3.  The **capacity** of the slice (`cap`), which is the total size of the underlying array.

When you pass a slice to a function, you are copying this header. Since the pointer in the copied header still points to the *same* original array, changes to the elements of the slice within the function will be visible outside the function.

### Maps

A `map` is a collection of key-value pairs. It's also a reference type, like a slice.
-   You must initialize a map before you can add to it using `make(map[keyType]valueType)` or a map literal.
-   When you iterate over a map, the order is **not** guaranteed.

A common pattern is checking if a key exists using the "comma ok" idiom:
```go
scores := map[string]int{"alice": 10, "bob": 8}

score, ok := scores["alice"] // score will be 10, ok will be true
score, ok = scores["charlie"] // score will be 0, ok will be false
```

### Interfaces

The function you will implement takes an `io.Reader` as an argument. `io.Reader` is an **interface**, which is a type that specifies a set of methods. Any type that implements those methods can be used as an `io.Reader`. For example, a file, a network connection, or even a string can be an `io.Reader`. This makes your function very flexible.

## Your Task

Your task is to implement the `FreqFromReader(r io.Reader) (map[string]int, string, error)` function in `internal/arraysmapsbasics/exercise.go`.

This function should:
1.  Read the input from the `io.Reader` line by line. A `bufio.Scanner` is great for this.
2.  For each line, split it into words. (Hint: the `strings` package is useful here).
3.  Normalize each word to lowercase.
4.  Use a map to keep track of how many times each word appears.
5.  After counting, find the most frequent word in the map.
6.  Return the frequency map, the most frequent word, and any error that occurred.

Open `internal/arraysmapsbasics/exercise.go` and fill in the `// TODO` sections.

## How to Verify Your Work

Run the following command from this directory (`minis/02-arrays-maps-basics`):

```bash
go test -v ./...
```

If the tests pass, you have successfully completed the lesson.

## Related Lessons
- Previous: `minis/01-hello-strings`
- Next: `minis/03-csv-stats`
