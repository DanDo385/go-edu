# 07: Generic LRU Cache

This project challenges you to build a high-performance, thread-safe, and reusable **Least Recently Used (LRU) Cache**. This is a classic computer science problem that sits at the intersection of data structures, algorithms, and concurrent programming. You will also leverage one of Go's most powerful modern features: **Generics**.

## Table of Contents

- [Learning Objectives](#learning-objectives)
- [The Big Picture: The Need for Speed](#the-big-picture-the-need-for-speed)
- [First Principles: Caching Strategies](#first-principles-caching-strategies)
  - [The LRU Algorithm Explained](#the-lru-algorithm-explained)
- [Project Structure](#project-structure)
- [Key Concepts in This Project](#key-concepts-in-this-project)
  - [The Optimal Data Structure: Hash Map + Doubly Linked List](#the-optimal-data-structure-hash-map--doubly-linked-list)
  - [Go Generics: `[K comparable, V any]`](#go-generics-k-comparable-v-any)
  - [Thread Safety with `sync.Mutex`](#thread-safety-with-syncmutex)
- [Progression: Building Reusable Components](#progression-building-reusable-components)
- [How to Run and Test](#how-to-run-and-test)
- [Key Takeaways](#key-takeaways)
- [Further Reading](#further-reading)

## Learning Objectives

By the end of this project, you will be able to:

-   **Implement the LRU caching algorithm** from scratch.
-   **Combine a hash map and a doubly linked list** to create a data structure with O(1) time complexity for reads and writes.
-   **Write generic, type-safe Go code** using type parameters.
-   **Make a data structure safe for concurrent use** by protecting it with a `sync.Mutex`.
-   **Understand the trade-offs** between different caching eviction policies.

## The Big Picture: The Need for Speed

In any application, some operations are more expensive than others. Fetching data from a database, calling a remote API, or performing a complex computation all take time. If you need the same data repeatedly, fetching it over and over is wasteful and slow.

A **cache** solves this problem by storing the results of expensive operations in a faster, closer data store (usually in memory). The next time the same result is needed, it can be served directly from the cache, bypassing the slow operation entirely. This is a fundamental technique for building high-performance systems.

But memory is finite. You can't store everything forever. This leads to the central question of caching: when the cache is full, which item should you discard? This is the "eviction policy." The LRU policy is one of the most effective and widely used.

## First Principles: Caching Strategies

-   **FIFO (First-In, First-Out)**: The oldest item is discarded. Simple, but often ineffective if an old item is still very popular.
-   **LFU (Least Frequently Used)**: The item that has been accessed the fewest times is discarded. Requires tracking access counts, which can be complex.
-   **LRU (Least Recently Used)**: The item that has been accessed the longest time ago is discarded. This is a brilliant heuristic: if a piece of data hasn't been used recently, it's likely not going to be used again soon. This is the policy you will implement.

### The LRU Algorithm Explained

Imagine a cache with a capacity of 3 items.

1.  `Put("A", 1)`: Cache: `[A]`
2.  `Put("B", 2)`: Cache: `[B, A]` (B is now the most recently used)
3.  `Put("C", 3)`: Cache: `[C, B, A]` (C is the most recent)
4.  `Get("B")`: B was just used, so it moves to the front. Cache: `[B, C, A]`
5.  `Put("D", 4)`: The cache is full! The least recently used item is "A". It gets evicted. Cache: `[D, B, C]`

To implement this efficiently, we need a way to:
1.  Look up an item by its key in O(1) time. (A hash map is perfect for this).
2.  Move any item to the "front" of the recently-used list in O(1) time. (An array or slice would be O(n) for this!).
3.  Remove the "last" item in O(1) time.

## Project Structure

```
.
├── cmd/
│   └── dev/
│       └── main.go       # A simple program to demonstrate the LRU cache.
└── internal/
    └── lru/
        └── lru.go        # Your LRU cache implementation.
```

## Key Concepts in This Project

### The Optimal Data Structure: Hash Map + Doubly Linked List

The classic and most efficient way to build an LRU cache is by combining two data structures:

1.  **Hash Map (`map[K]*list.Element`)**:
    -   **Key**: The key of the cached item (`K`).
    -   **Value**: A *pointer* to a node in a doubly linked list.
    -   **Purpose**: Provides O(1) average time complexity for lookups. Given a key, we can instantly find the corresponding node in our linked list.

2.  **Doubly Linked List (`container/list`)**:
    -   Each node in the list stores a key-value pair.
    -   **Purpose**: Maintains the order of "recency." The front of the list is the most recently used item, and the back is the least recently used.
    -   A doubly linked list is crucial because it allows O(1) time complexity for moving a node to the front and removing a node from the back.

```
          Get("B") -> O(1) map lookup -> finds node B -> O(1) list move
             +-------------------------------------------------+
             |                                                 |
             v
  HashMap: { "A": -> NodeA, "B": -> NodeB, "C": -> NodeC }
                       |           |           |
                       |           |           |
                       v           v           v
LinkedList:  (MRU) Head <--> NodeC <--> NodeB <--> NodeA <--> Tail (LRU)
```

### Go Generics: `[K comparable, V any]`

Before Go 1.18, writing a reusable cache meant using `interface{}`. This was not type-safe and required clumsy type assertions. Generics solve this beautifully.

```go
type LRUCache[K comparable, V any] struct {
    // ...
}
```

-   `[K comparable, V any]`: This is a type parameter list.
-   `K`: Represents the type of the cache keys.
-   `V`: Represents the type of the cache values.
-   `comparable`: This is a constraint. It means that `K` must be a type that can be compared with `==` and `!=`, which is a requirement for map keys.
-   `any`: This is an alias for `interface{}`. It means `V` can be any type.

Now, you can create a cache for specific types with full compile-time safety:
```go
// A cache from user IDs (int) to user objects (*User)
cache1 := lru.New[int, *User](128)

// A cache from API endpoint strings to []byte responses
cache2 := lru.New[string, []byte](256)
```

### Thread Safety with `sync.Mutex`

If multiple goroutines try to access the cache at the same time, you can have a "race condition," leading to data corruption. For example, one goroutine might be moving a node while another is trying to evict it.

The solution is to protect the cache's internal state (the map and the list) with a `sync.Mutex` (a mutual exclusion lock).

```go
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
    c.mu.Lock()         // Lock the cache
    defer c.mu.Unlock() // Defer unlock until the function returns

    // ... safe to access c.items and c.list here ...
}
```
Any operation that reads or modifies the cache's internal structures must be wrapped in a `Lock()`/`Unlock()` pair.

## Progression: Building Reusable Components

This project takes the concept of data structures from **Project 02** and concurrency from **Project 06** and combines them to build a sophisticated, reusable, and generic component. An LRU cache is a high-value building block that you can drop into many other applications to improve their performance.

## How to Run and Test

```bash
# Run the development harness to see the cache in action
go run ./cmd/dev

# Run the tests to verify correctness, eviction logic, and concurrency safety
go test -v ./...
```

## Key Takeaways

-   An LRU cache evicts the least recently used item to make space for new items.
-   The combination of a **hash map and a doubly linked list** provides the optimal O(1) performance for `Get` and `Put`.
-   **Go generics** allow you to write reusable, type-safe data structures that are not tied to specific types.
-   Any data structure intended for concurrent use must be protected by a **mutex**.

## Further Reading

-   [**Go Blog: An Introduction To Generics**](https://go.dev/blog/intro-generics)
-   [**Package `container/list`**](https://pkg.go.dev/container/list): The standard library's doubly linked list implementation.
-   [**Go by Example: Mutexes**](https://gobyexample.com/mutexes)
-   [**Wikipedia: Cache replacement policies**](https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU)): A deeper look at the theory behind LRU and other policies.