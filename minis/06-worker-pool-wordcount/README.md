# 06: Worker Pool Wordcount

## Core Concepts

- The concrete problem in Worker Pool Wordcount and the correctness invariants it depends on.
- How values, pointers, slices, maps, or channels behave in this module's runtime path.
- Why this lesson's implementation pattern is the right next step in the learning arc.

## CS Connection

- Data representation and state transitions: what is copied, what is shared, and what can race.
- API boundaries: where we validate, where we propagate errors, and where we normalize output.
- Algorithmic tradeoffs in this lesson (latency, throughput, memory, and complexity).

## End-State Understanding

- Diagnose and implement Worker Pool Wordcount patterns without relying on hidden framework behavior.
- Explain memory and concurrency implications of the final implementation choices.
- Compare learner code and reference code by invariants, not only by syntax.

## Why This Lesson Exists Here

Problem statement:
This lesson turns the previous module's concepts into a reusable engineering pattern for worker pool wordcount.

At this point in the arc:
Lesson 06 introduces a sharper systems concern so later modules can assume this mental model is stable.

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


This project introduces you to **concurrency**, Go's most celebrated feature. You will refactor a sequential task to run in parallel using a **worker pool**, a fundamental pattern for high-performance applications.

## Core Concepts

- Value semantics in Go: what gets copied at function calls and what can still alias shared state.
- Ownership boundaries for mutation, especially when multiple code paths touch the same logical data.

## CS Connection

- Memory layout drives behavior: variables store values, and some values are addresses into other storage.
- Go is pass-by-value, including pointers; copying a pointer value copies an address, not the pointee.
- Correctness depends on understanding copying versus aliasing (`*T`, slices, maps, and channels) and enforcing synchronization when concurrent access exists.

## End-State Understanding

- Implement the exercise with explicit reasoning about correctness, edge cases, and error paths.

## What You'll Learn

- How to write concurrent code using **goroutines** and **channels**.
- How to implement the **worker pool pattern**.
- How to use `sync.WaitGroup` to wait for goroutines to finish.
- The difference between **concurrency** (structuring work) and **parallelism** (executing work).
- Go's core concurrency philosophy: "Share memory by communicating."

## The Challenge: Doing More, Faster

Imagine you need to process 100 web pages. A sequential program does this one by one, spending most of its time *waiting* for the network. Concurrency allows your program to work on other pages while it waits, dramatically improving performance.

## Core Go Primitives

### Goroutines
A goroutine is a lightweight thread. Starting one is as simple as adding the `go` keyword to a function call:
```go
go myFunction() // This function now runs concurrently.
```

### Channels
A channel is a typed "pipe" that connects concurrent goroutines. You send values into a channel from one goroutine and receive them in another. This is the primary way to communicate safely between goroutines.
```go
ch := make(chan int) // Create a channel of integers.

go func() {
    ch <- 42 // Send the value 42 into the channel.
}()

value := <-ch // Receive the value from the channel.
```
**Note:** Channels are reference types. When you pass a channel to a function, you are passing a pointer to the same underlying data structure, allowing different goroutines to communicate.

### `sync.WaitGroup`
A `WaitGroup` is a simple counter used to wait for a collection of goroutines to finish.
1.  Call `wg.Add(N)` before starting `N` goroutines.
2.  Each goroutine calls `wg.Done()` when it finishes.
3.  The main goroutine calls `wg.Wait()` to block until all goroutines are done.

## The Worker Pool Pattern

You will implement a classic worker pool to count word frequencies from a list of URLs.
```
+----------+      +--------------+      +-----------+      +-----------------+
|          |----->| jobs channel |----->| Worker 1  |----->|                 |
|          |      | (URLs)       |      +-----------+      |                 |
| Main     |----->|              |----->| Worker 2  |----->| results channel |
| Goroutine|      +--------------+      +-----------+      | (word counts)   |
|          |                            |   ...     |      |                 |
|          |      (sends URLs)           (process URLs)   (receives counts) |
+----------+                                               +-----------------+
```
The main goroutine creates a set of "worker" goroutines. It sends "jobs" (URLs) to them via a `jobs` channel. The workers process the URLs and send their results back over a `results` channel. The main goroutine then collects and aggregates all the results.

## Your Task

Your task is to implement three functions in `internal/workerpoolwordcount/exercise.go`. The logic is broken down into helper functions to make it more manageable.

1.  **`tokenizeAndCount(text string) map[string]int`**
    - This is a simple, non-concurrent function.
    - Take the input text, split it into words, and return a map of word frequencies. A `bufio.Scanner` is great for this.

2.  **`fetchAndCount(ctx context.Context, url string) (map[string]int, error)`**
    - This function fetches the content of a single URL.
    - It should use `http.NewRequestWithContext` to make the HTTP request cancellable.
    - After fetching the body, it should call `tokenizeAndCount` to get the word frequencies for that URL.

3.  **`WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error)`**
    - This is the main function where you will implement the worker pool pattern.
    - Create channels for `jobs` and `results`.
    - Start `workers` number of worker goroutines. Each worker should read URLs from the `jobs` channel, call `fetchAndCount` for each URL, and send the result to the `results` channel.
    - Send all the `urls` to the `jobs` channel.
    - Wait for all workers to finish and aggregate the results from the `results` channel into a single map.

Open `internal/workerpoolwordcount/exercise.go` and fill in the `// TODO` sections for these three functions.

## How to Verify Your Work

Run the following command from this directory (`minis/06-worker-pool-wordcount`):

```bash
go test -v ./...
```
If the tests pass, you have successfully completed the lesson.

## Related Lessons
- Previous: `minis/05-cli-todo-files`
- Next: `minis/07-generic-lru-cache`
