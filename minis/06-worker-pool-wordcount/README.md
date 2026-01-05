# 06: Concurrent Worker Pool

## The Big Picture: Doing More, Faster

All the applications you've built so far have been **sequential**: they do one thing at a time. This is simple, but slow. If you need to process 100 web pages, a sequential program fetches the first, processes it, then fetches the second, and so on. Most of its time is spent *waiting* for the network.

This project introduces you to **concurrency**, Go's most celebrated feature. You will refactor a word-counting task to fetch multiple URLs at the same time, dramatically reducing the total wait time. You'll achieve this by building a **worker pool**, a fundamental pattern for high-performance applications like web crawlers, video processors, and high-traffic servers.

## First Principles: Concurrency and Communication

1.  **Concurrency vs. Parallelism**:
    *   **Concurrency** is about *dealing* with many things at once. It's a way to structure your program to handle multiple tasks, like making progress on task B while task A is blocked (e.g., waiting for a network response).
    *   **Parallelism** is about *doing* many things at once. It requires multiple CPU cores to execute tasks simultaneously.
    *   Go's concurrency model makes it easy to write concurrent code that can then run in parallel on multi-core hardware.

2.  **Why Concurrency?**: It's essential for performance, especially for **I/O-bound** tasks. While one part of your program is waiting for a slow operation (like a network request or reading a file), the Go runtime can schedule another part of your program to run on the CPU. This maximizes CPU utilization and dramatically improves throughput.

3.  **Communication Models**:
    *   **Shared Memory**: Multiple concurrent tasks communicate by reading and writing to the same piece of memory. This is common in many languages but requires complex locking mechanisms (`mutexes`) to prevent data races and corruption.
    *   **Message Passing**: Concurrent tasks have no shared state. Instead, they communicate by sending messages to each other over a channel. This is Go's preferred model. The famous Go proverb is: *"Do not communicate by sharing memory; instead, share memory by communicating."*

## Key Go Primitives for Concurrency

Go provides a small, elegant set of tools to build powerful concurrent systems.

### Goroutines

A goroutine is an incredibly lightweight thread managed by the Go runtime, not the operating system. You can have hundreds of thousands, or even millions, of goroutines running in a single process. Starting one is as simple as adding the `go` keyword to a function call: `go myFunction()`.

### Channels

A channel is a typed "pipe" that connects concurrent goroutines. You can send values into a channel from one goroutine and receive those values in another. Channels are the primary way to orchestrate and communicate between goroutines safely.

```go
// Create a channel that can only carry strings.
ch := make(chan string)

// In one goroutine: send a message.
ch <- "hello"

// In another goroutine: receive the message.
msg := <-ch
```

### The `select` Statement

The `select` statement lets a goroutine wait on multiple channel operations. It's like a `switch` statement for channels. It will block until one of its cases can run, then it will execute that case. This is a powerful tool for implementing timeouts, cancellations, and complex orchestration.

### The `context` Package

How do you tell a dozen goroutines to stop what they're doing because a user cancelled a request or a deadline was exceeded? The `context` package is the answer. A `Context` is an object that carries cancellation signals, deadlines, and other request-scoped values across API boundaries and between goroutines. It's an essential tool for building robust, production-ready concurrent systems.

### `sync.WaitGroup`

A `WaitGroup` is a simple but powerful counter used to wait for a collection of goroutines to finish.
*   You call `wg.Add(1)` for each goroutine you start.
*   Each goroutine calls `wg.Done()` when it finishes.
*   The main goroutine calls `wg.Wait()` to block until the counter is zero.

## The Worker Pool Pattern

In this project, you will implement a classic worker pool to fetch and process URLs.
1.  **`jobs` Channel**: A central channel where the main goroutine sends the URLs that need to be processed.
2.  **Workers**: A fixed number of worker goroutines are started. Each one is in a loop, pulling a URL from the `jobs` channel.
3.  **`results` Channel**: After a worker fetches and processes a URL, it sends the resulting word count map to a `results` channel.
4.  **Aggregator**: The main goroutine reads all the maps from the `results` channel and merges them into a single, final word count map.
5.  **Coordination**: `WaitGroup` and `context` are used to manage cancellation and ensure a graceful shutdown.

### Simplifying with `errgroup`

After implementing the worker pool manually, you'll see a second version using the `golang.org/x/sync/errgroup` package. An `errgroup` is a higher-level abstraction that elegantly handles much of the boilerplate for you:
*   It combines a `WaitGroup` and a `Context`.
*   It automatically cancels the context for all goroutines if any single goroutine returns an error.
*   Its `Wait()` method returns the first error that occurred.
This results in cleaner, safer, and more concise concurrent code.

## Progression: From Sequential to Concurrent

This project marks a significant leap in your Go journey. You are taking the word-counting logic from **Project 02** and parallelizing it. This demonstrates how to take a single-threaded task and apply Go's concurrency model to make it dramatically faster for I/O-bound workloads.

Mastering concurrency is the key to unlocking Go's full potential for building high-performance network services, data processing systems, and more. The patterns you learn here are directly applicable to nearly every large-scale Go application.

## How to Run and Test

```bash
# The dev harness runs your function against a set of mock URLs.
go run ./cmd/dev

# The tests cover concurrent execution, error handling, and cancellation.
go test -v ./...
```