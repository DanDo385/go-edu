# 02: Arrays, Slices, and Maps

## The Big Picture: Organizing Data

In the previous project, you learned about `strings`—Go's primitive for handling text. Now, we move to the next logical step: organizing collections of data. Nearly every useful program needs to manage lists of items (like a list of users) or associate data with a specific key (like looking up a user's ID to find their profile).

This project introduces Go's three fundamental collection types: **arrays**, **slices**, and **maps**. Mastering them is essential for writing any non-trivial Go program.

## Computer Science First Principles: Data Structures

A **data structure** is a systematic way of organizing and storing data to perform operations efficiently.

1.  **Array**: The most basic data structure. It's a **fixed-size collection of elements of the same type, stored in a contiguous block of memory**.
    *   **Strengths**: Extremely fast access. If you know the index `i`, you can jump directly to the memory address of the element `i` (O(1) time complexity).
    *   **Weaknesses**: Fixed size. You must know the size of the collection when you create it. Adding or removing elements is inefficient.

2.  **Dynamic Array (Slice in Go)**: An abstraction built on top of an array. It can **grow or shrink in size**.
    *   **How it works**: It holds a pointer to an underlying array. When the dynamic array runs out of space, it creates a *new, larger* underlying array and copies the elements from the old array to the new one.
    *   **Strengths**: Flexible size. Provides the benefits of array-like access with the ability to grow.
    *   **Weaknesses**: Appending can sometimes be slow if a resize and copy is triggered.

3.  **Hash Map (Map in Go)**: A data structure that stores **key-value pairs**.
    *   **How it works**: It uses a "hash function" to compute an index (or "bucket") for a given key. This allows it to store and retrieve values in near-constant time.
    *   **Strengths**: Very fast lookups, insertions, and deletions on average (O(1)). Ideal for associating data, like `userID -> UserProfile`.
    *   **Weaknesses**: Unordered. The elements are not stored in any particular sequence. Uses more memory than an array or slice.

## Go's Implementations: A Deeper Dive

Go provides direct and powerful implementations of these concepts.

### Arrays

In Go, an array is a **value type**. When you assign or pass an array, the entire array is **copied**.

```go
// An array of 3 integers. The size `3` is part of its type!
var a [3]int // [0, 0, 0]

b := a // `b` is a complete copy of `a`.
b[0] = 100

fmt.Println(a[0]) // Prints 0. `a` is unchanged.
```

Because of this copying behavior and their fixed size, arrays are less common in Go than slices. They are primarily used when you need a very specific, fixed-size collection, often for performance-critical code or interacting with C libraries.

### Slices

The **slice** is Go's most important and versatile collection type. It's a lightweight struct that provides a "view" into an underlying array.

A slice is a 3-word header:
1.  **Pointer**: Points to the first element of the underlying array that the slice can access.
2.  **Length (`len`)**: The number of elements in the slice.
3.  **Capacity (`cap`)**: The total number of elements in the underlying array, starting from the slice's pointer.

```go
// sliceHeader = {
//     Data *T // Pointer to an element in the underlying array
//     Len  int // Number of elements in the slice
//     Cap  int // Number of elements from Data to the end of the array
// }
```

*   **Dynamic Size**: The built-in `append` function handles growth. If `len` exceeds `cap`, `append` allocates a new, larger array and copies the data over.
*   **Reference Semantics**: A slice header is a value, but it *points to* shared data. If you pass a slice and modify its elements (without appending), the changes are visible to the caller.

```go
s1 := []int{1, 2, 3}
s2 := s1 // s2 is a copy of s1's header. Both point to the same array.
s2[0] = 99

fmt.Println(s1[0]) // Prints 99. The change is visible via s1.
```

### Maps

Go's `map` is a reference type that implements a hash table.

```go
// Create a map where keys are strings and values are ints
freq := make(map[string]int)

freq["hello"] = 1
freq["world"] = 5
```

Key features:
*   **Reference Type**: Like slices, a map variable is a pointer to an underlying data structure. When you pass a map, you are passing a pointer, so modifications are visible to the caller.
*   **Zero Value `nil`**: The zero value of a map is `nil`. A `nil` map cannot have keys added to it. You must initialize it with `make()` or a map literal.
*   **"Comma Ok" Idiom**: How do you know if a key exists? Accessing a non-existent key returns the zero value for the value type (e.g., `0` for `int`, `""` for `string`). To distinguish between a stored zero and a non-existent key, use the "comma ok" idiom:
    ```go
    count, ok := freq["goodbye"]
    if !ok {
        fmt.Println("The key 'goodbye' is not in the map.")
        // count is 0 here
    }
    ```
*   **Randomized Iteration**: When you iterate over a map with a `for...range` loop, the order is **not guaranteed**. Go intentionally randomizes the starting point of iteration to prevent developers from relying on a specific order.

## Project Task: Word Frequency Counter

This project combines string processing (from Project 01) with maps. You will read text from an `io.Reader`, process it line by line, and use a `map[string]int` to count the frequency of each word.

This simple task is a classic computer science problem and a building block for many applications:
*   **Search Engines**: Use frequency counts to help rank search results (see: TF-IDF).
*   **Data Analytics**: Find the most common items in a dataset.
*   **Natural Language Processing**: Forms the basis of many language models.

## Progression: The Workhorse of Go

Slices and maps are ubiquitous. You will use them in almost every subsequent project:
*   **03-csv-stats**: You'll read CSV rows into slices of strings (`[]string`).
*   **06-worker-pool-wordcount**: You'll use maps to aggregate results from concurrent workers.
*   **10-grpc-telemetry-service**: You'll store metrics in maps.
*   **40-merkle-tree-basics**: You'll build a tree structure using slices of nodes.

Understanding the performance and memory characteristics of slices and maps is crucial for writing efficient and scalable Go code.

## How to Run and Test

```bash
# Run the development harness to see your function in action
go run ./cmd/dev

# Run the provided tests to check your implementation for correctness
go test -v ./...
```