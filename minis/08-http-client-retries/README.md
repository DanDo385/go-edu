# 08: Http Client Retries

## Core Concepts

- The concrete problem in Http Client Retries and the correctness invariants it depends on.
- How values, pointers, slices, maps, or channels behave in this module's runtime path.
- Why this lesson's implementation pattern is the right next step in the learning arc.

## CS Connection

- Data representation and state transitions: what is copied, what is shared, and what can race.
- API boundaries: where we validate, where we propagate errors, and where we normalize output.
- Algorithmic tradeoffs in this lesson (latency, throughput, memory, and complexity).

## End-State Understanding

- Diagnose and implement Http Client Retries patterns without relying on hidden framework behavior.
- Explain memory and concurrency implications of the final implementation choices.
- Compare learner code and reference code by invariants, not only by syntax.

## Why This Lesson Exists Here

Problem statement:
This lesson turns the previous module's concepts into a reusable engineering pattern for http client retries.

At this point in the arc:
Lesson 08 introduces a sharper systems concern so later modules can assume this mental model is stable.

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


This project tackles a fundamental challenge of distributed systems: network unreliability. You will build a robust HTTP client that intelligently handles temporary failures by automatically retrying requests using the **exponential backoff with jitter** algorithm.

## CS Connection

- Memory layout drives behavior: variables store values, and some values are addresses into other storage.
- Go is pass-by-value, including pointers; copying a pointer value copies an address, not the pointee.
- Correctness depends on understanding copying versus aliasing (`*T`, slices, maps, and channels) and enforcing synchronization when concurrent access exists.

## End-State Understanding

- Implement the exercise with explicit reasoning about correctness, edge cases, and error paths.

## What You'll Learn

- Why simple HTTP requests can fail in a distributed system.
- How to implement **exponential backoff with jitter**.
- How to customize `http.Client` using the **`http.RoundTripper` interface**.
- How to apply the **Decorator design pattern** to create HTTP middleware.
- How to **classify errors** to decide whether a request should be retried.

## The Challenge: The Unreliable Network

A core fallacy of distributed computing is assuming the network is reliable. It isn't. Requests fail for temporary reasons all the time. A naive client gives up immediately. A **resilient** client understands that some failures are temporary and tries again.

The best practice for this is **Exponential Backoff with Jitter**:
1.  **Exponential Backoff**: After each failure, double the wait time before the next retry (e.g., 1s, 2s, 4s). This gives a struggling server time to recover.
2.  **Jitter**: Add a small, random amount of time to each wait. This prevents a "thundering herd" of clients from all retrying at the exact same time.

## Core Concepts: The Decorator Pattern with `http.RoundTripper`

How can we add this retry logic without rewriting Go's `http.Client`? By using the **Decorator Pattern**.

The `http.Client` has a `Transport` field, which is an `http.RoundTripper` interface. This interface has one method: `RoundTrip(*http.Request) (*http.Response, error)`.

We can create our own `retryRoundTripper` struct that *also* implements this interface. Our struct will hold a reference to the *next* `RoundTripper` in the chain (e.g., Go's default one that actually makes the network call).

```
        Your Code             Your Wrapper              Go's Default
+---------------------+   +---------------------+   +------------------------+
| http.Client.Get()   |-->| retryRoundTripper   |-->| http.DefaultTransport  |
|                     |   |   .RoundTrip()      |   |   .RoundTrip()         |
+---------------------+   | (Contains your      |   | (Makes actual network |
                          |  retry loop)        |   |  request)              |
                          +---------------------+   +------------------------+
```

Inside our `RoundTrip` method, we will implement a `for` loop for the retries. Inside the loop, we call the `next` RoundTripper. If it fails with a *retryable* error, we wait and let the loop continue. If it succeeds, or fails with a *permanent* error, we return the result.

This is a powerful "middleware" pattern for building flexible and modular HTTP clients.

## Your Task

Your task is to implement the `NewRetryClient` function and the `RoundTrip` method for the `retryRoundTripper` in `internal/httpclientretries/exercise.go`.

1.  **`NewRetryClient(...)`**: This constructor should create and return a new `http.Client` configured to use your `retryRoundTripper` as its `Transport`.
2.  **`isRetryable(err error, resp *http.Response) bool`**: This helper function is crucial. It needs to inspect the error and the HTTP response to decide if the request should be retried.
    - Network-level errors are generally retryable.
    - HTTP status codes `429`, `502`, `503`, and `504` are retryable.
    - Other `4xx` and `5xx` codes are generally *not* retryable.
3.  **`RoundTrip(req *http.Request) (*http.Response, error)`**: This is the core of the lesson.
    - Implement a `for` loop that runs for a maximum number of retries.
    - Inside the loop, call the `next.RoundTrip(req)` method.
    - Use `isRetryable` to check the result.
    - If the request should be retried, calculate the backoff duration, add jitter, and wait (e.g., `time.Sleep`).
    - Make sure the loop respects the `context` of the request for cancellations.

## How to Verify Your Work

Build the test application from the lesson directory:
```bash
go build -o client ./cmd/app
```

Run it against a URL that will likely succeed:
```bash
./client http://google.com
```

Run it against a URL that is designed to fail, to see your retry logic in action:
```bash
./client http://example.com:81
```

Finally, run the automated tests:
```bash
go test -v ./...
```
If the tests pass, you have successfully completed the lesson.

## Related Lessons
- Previous: `minis/07-generic-lru-cache`
- Next: `minis/09-http-server-graceful`
