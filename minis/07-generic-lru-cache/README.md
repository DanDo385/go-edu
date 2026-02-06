# 07: Generic Lru Cache

## Core Concepts

- The concrete problem in Generic Lru Cache and the correctness invariants it depends on.
- How values, pointers, slices, maps, or channels behave in this module's runtime path.
- Why this lesson's implementation pattern is the right next step in the learning arc.

## CS Connection

- Data representation and state transitions: what is copied, what is shared, and what can race.
- API boundaries: where we validate, where we propagate errors, and where we normalize output.
- Algorithmic tradeoffs in this lesson (latency, throughput, memory, and complexity).

## End-State Understanding

- Diagnose and implement Generic Lru Cache patterns without relying on hidden framework behavior.
- Explain memory and concurrency implications of the final implementation choices.
- Compare learner code and reference code by invariants, not only by syntax.

## Why This Lesson Exists Here

Problem statement:
This lesson turns the previous module's concepts into a reusable engineering pattern for generic lru cache.

At this point in the arc:
Lesson 07 introduces a sharper systems concern so later modules can assume this mental model is stable.

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


This project challenges you to build a high-performance, thread-safe, and reusable **Least Recently Used (LRU) Cache**. This is a classic computer science problem that combines data structures, algorithms, and concurrency. You will also use one of Go's most powerful modern features: **Generics**.

## CS Connection

- Memory layout drives behavior: variables store values, and some values are addresses into other storage.
- Go is pass-by-value, including pointers; copying a pointer value copies an address, not the pointee.
- Correctness depends on understanding copying versus aliasing (`*T`, slices, maps, and channels) and enforcing synchronization when concurrent access exists.

## End-State Understanding

- Implement the exercise with explicit reasoning about correctness, edge cases, and error paths.

## What You'll Learn

- How to implement the **LRU caching algorithm**.
- How to combine a **hash map and a doubly linked list** for O(1) performance.
- How to write **generic, type-safe** code using Go Generics.
- How to make a data structure **thread-safe** using a `sync.Mutex`.

## The Challenge: Caching and Eviction

A **cache** stores the results of expensive operations (like database queries or API calls) in memory to make future requests faster. But memory is finite. When the cache is full, you need a policy to decide which item to discard (or "evict").

**LRU (Least Recently Used)** is a popular and effective eviction policy: when you need to make space, you discard the item that hasn't been used in the longest time.

## Core Concepts

### The Optimal Data Structure for LRU

To build an efficient LRU cache, we need to perform two operations in constant time, O(1):
1.  Look up an item by its key.
2.  Move an item to the "front" of the recency list.

The solution is to combine two data structures:
1.  **A Hash Map (`map[K]*list.Element`)**: For O(1) lookups. The map's value will be a *pointer* to a node in a linked list.
2.  **A Doubly Linked List (`container/list`)**: To maintain the order of recency. The front of the list is the most recently used item. A doubly linked list allows us to move any node to the front in O(1) time.

### Go Generics: `[K comparable, V any]`

Before generics, you had to use `interface{}` to write a reusable cache, which was not type-safe. Generics solve this.
```go
type LRUCache[K comparable, V any] struct { ... }
```
- `[K comparable, V any]` is a type parameter list.
- `K` is the type for keys. The `comparable` constraint means it must be a type that can be used as a map key.
- `V` is the type for values. `any` means it can be any type.

This lets you create a type-safe cache for any key/value pair.

### Thread Safety with `sync.Mutex`

If multiple goroutines try to access the cache at the same time, you can have a "race condition," leading to data corruption. The solution is to protect the cache's internal state (the map and the list) with a `sync.Mutex` (a mutual exclusion lock).
```go
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
    c.mu.Lock()         // Lock the cache. Only one goroutine can proceed past this point.
    defer c.mu.Unlock() // Defer unlock until the function returns. This is crucial!

    // It is now safe to access the map and list.
}
```
Any method that reads or writes to the cache's internal state must be protected by the mutex.

## Your Task

Your task is to implement the `New` function and the methods for the `LRUCache` struct in `internal/genericlrucache/exercise.go`.

1.  **`New[K comparable, V any](capacity int) *LRUCache[K, V]`**: The constructor. It should initialize the `LRUCache` struct, including the map and the linked list.
2.  **`Get(key K) (V, bool)`**: This method should look up a key in the map. If found, it should move the corresponding list element to the front of the list (to mark it as most recently used) and return the value.
3.  **`Set(key K, value V)`**: This method should add a new key-value pair.
    - If the key already exists, update its value and move it to the front.
    - If the key is new, add it.
    - If the cache is at capacity, you must evict the least recently used item (the one at the back of the list).
    - **Remember to wrap your methods with a mutex lock and unlock!**

## How to Verify Your Work

Run the following command from this directory (`minis/07-generic-lru-cache`):

```bash
go test -v ./...
```
If the tests pass, you have successfully completed the lesson.

## Related Lessons
- Previous: `minis/06-worker-pool-wordcount`
- Next: `minis/08-http-client-retries`
