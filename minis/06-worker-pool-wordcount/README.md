# 06: Concurrent Worker Pool

This project introduces you to **concurrency**, Go's most celebrated feature. You will refactor a sequential task to run in parallel, dramatically improving its performance for I/O-bound operations. You will build a **worker pool**, a fundamental pattern for high-performance applications like web crawlers, video processors, and high-traffic servers.

## Table of Contents

- [Learning Objectives](#learning-objectives)
- [The Big Picture: Doing More, Faster](#the-big-picture-doing-more-faster)
- [First Principles: Concurrency and Communication](#first-principles-concurrency-and-communication)
- [Project Structure](#project-structure)
- [Key Go Primitives for Concurrency](#key-go-primitives-for-concurrency)
  - [Goroutines](#goroutines)
  - [Channels](#channels)
  - [The `select` Statement](#the-select-statement)
  - [`sync.WaitGroup`](#syncwaitgroup)
- [The Worker Pool Pattern](#the-worker-pool-pattern)
- [Simplifying with `errgroup`](#simplifying-with-errgroup)
- [Progression: From Sequential to Concurrent](#progression-from-sequential-to-concurrent)
- [How to Run and Test](#how-to-run-and-test)
- [Key Takeaways](#key-takeaways)
- [Further Reading](#further-reading)

## Learning Objectives

By the end of this project, you will be able to:

-   **Explain the difference between concurrency and parallelism**.
-   **Write concurrent code** using goroutines, channels, and `sync.WaitGroup`.
-   **Implement the worker pool pattern** to process tasks in parallel.
-   **Understand and apply Go's core concurrency philosophy**: "Share memory by communicating."
-   **Use the `context` package** to manage cancellation and deadlines in concurrent code.
-   **Refactor complex concurrent logic** into cleaner code using the `errgroup` package.

## The Big Picture: Doing More, Faster

All the applications you've built so far have been **sequential**: they do one thing at a time. If you need to process 100 web pages, a sequential program fetches the first, processes it, then fetches the second, and so on. Most of its time is spent *waiting* for the network. Concurrency allows the program to work on other tasks while waiting, maximizing efficiency.

## First Principles: Concurrency and Communication

1.  **Concurrency vs. Parallelism**:
    *   **Concurrency** is about *dealing* with many things at once. It's a way to structure your program to handle multiple tasks independently.
    *   **Parallelism** is about *doing* many things at once. It requires multiple CPU cores to execute tasks simultaneously.
    *   Go's concurrency model makes it easy to write concurrent code that can then run in parallel on multi-core hardware.

2.  **Communication Models**:
    *   **Shared Memory (The Old Way)**: Multiple concurrent tasks communicate by reading and writing to the same piece of memory. This requires complex locking (`mutexes`) to prevent data corruption. It is difficult and error-prone.
        ```go
        // Not the Go way
        mu.Lock()
        counter++
        mu.Unlock()
        ```
    *   **Message Passing (The Go Way)**: Concurrent tasks have no shared state. Instead, they communicate by sending messages over channels. This is Go's preferred model. The famous Go proverb is: *"Do not communicate by sharing memory; instead, share memory by communicating."*
        ```go
        // The Go way
        resultsChan <- result
        ```

## Project Structure

```
.
├── cmd/
│   └── dev/
│       └── main.go       # Entry point that runs the concurrent word counter.
└── internal/
    ├── sequential.go # The original, slow implementation.
    ├── workerpool.go # Your implementation of the manual worker pool.
    └── errgroup.go   # A cleaner implementation using the errgroup package.
```
-   **`cmd/dev`**: Runs the word count process and prints the results.
-   **`internal/`**: Contains the three different implementations, allowing you to compare and contrast the sequential, manual concurrent, and `errgroup` approaches.

## Key Go Primitives for Concurrency

### Goroutines
A goroutine is an incredibly lightweight thread managed by the Go runtime. Starting one is as simple as adding the `go` keyword to a function call: `go myFunction()`.

### Channels
A channel is a typed "pipe" that connects concurrent goroutines. You can send values into a channel from one goroutine and receive those values in another. Channels are the primary way to orchestrate and communicate between goroutines safely.

### The `select` Statement
The `select` statement lets a goroutine wait on multiple channel operations. It's like a `switch` statement for channels. It blocks until one of its cases can run, then executes that case.

### `sync.WaitGroup`
A `WaitGroup` is a simple counter used to wait for a collection of goroutines to finish.
-   Call `wg.Add(N)` before starting `N` goroutines.
-   Each goroutine calls `wg.Done()` when it finishes.
-   The main goroutine calls `wg.Wait()` to block until the counter is zero.

## The Worker Pool Pattern

In this project, you will implement a classic worker pool.

```
+----------+      +----------------+      +-----------+      +-----------------+      +-------------+
|          |----->|                |----->| Worker 1  |----->|                 |----->|             |
| Main     |      |  Jobs Channel  |      +-----------+      | Results Channel |      | Aggregator  |
| Goroutine|----->| (URLs to fetch)|----->| Worker 2  |----->| (Word Counts)   |----->| (Merges     |
|          |      |                |      +-----------+      |                 |      | results)    |
| (Sends   |      | cap: numJobs   |----->|   ...     |----->| cap: numJobs    |      |             |
| URLs)    |      +----------------+      +-----------+      +-----------------+      +-------------+
|          |                            | Worker N  |
+----------+                            +-----------+
```
1.  **`jobs` Channel**: The main goroutine sends the URLs to be processed into this channel.
2.  **Workers**: A fixed number of worker goroutines are started. Each one loops, pulling a URL from the `jobs` channel, processing it.
3.  **`results` Channel**: After a worker processes a URL, it sends the result (a word count map) to this channel.
4.  **Aggregator**: The main goroutine reads all the results from the `results` channel and merges them into a single, final map.
5.  **Coordination**: `WaitGroup` is used to ensure the aggregator waits for all workers to finish before closing the results channel.

## Simplifying with `errgroup`

After implementing the worker pool manually, you'll see a second version using the `golang.org/x/sync/errgroup` package. An `errgroup` is a higher-level abstraction that elegantly handles much of the boilerplate for you:
*   It combines a `WaitGroup` and a `context.Context`.
*   It automatically cancels the context for all goroutines if any single goroutine returns an error.
*   Its `Wait()` method blocks until all goroutines finish and returns the first error that occurred.
This results in cleaner, safer, and more concise concurrent code.

## Progression: From Sequential to Concurrent

This project marks a significant leap. You are taking the word-counting logic from **Project 02** and parallelizing it. This demonstrates how to take a single-threaded task and apply Go's concurrency model to make it dramatically faster for I/O-bound workloads. Mastering concurrency is the key to unlocking Go's full potential.

## How to Run and Test

```bash
# The dev harness runs your function against a set of mock URLs.
go run ./cmd/dev

# The tests cover concurrent execution, error handling, and cancellation.
go test -v ./...
```

## Key Takeaways

-   **Concurrency is for structuring work, Parallelism is for executing it.**
-   **Channels are the preferred way to communicate between goroutines.** Share memory by communicating.
-   A **worker pool** is a robust pattern for processing a large number of jobs concurrently.
-   **`sync.WaitGroup`** is essential for knowing when a group of goroutines has completed its work.
-   The **`errgroup` package** significantly simplifies the common pattern of a group of goroutines where one error should cancel the entire group.

## Further Reading

-   [**Go by Example: Goroutines**](https://gobyexample.com/goroutines)
-   [**Go by Example: Channels**](https://gobyexample.com/channels)
-   [**Blog: Go Concurrency Patterns: Pipelines and cancellation**](https://go.dev/blog/pipelines)
-   [**Package `errgroup`**](https://pkg.go.dev/golang.org/x/sync/errgroup)
-   [**Video: Concurrency is not Parallelism**](https://www.youtube.com/watch?v=oV9rvCAvHtQ) (A classic talk by Rob Pike)
