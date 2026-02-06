# 09: Http Server Graceful

## Core Concepts

- The concrete problem in Http Server Graceful and the correctness invariants it depends on.
- How values, pointers, slices, maps, or channels behave in this module's runtime path.
- Why this lesson's implementation pattern is the right next step in the learning arc.

## CS Connection

- Data representation and state transitions: what is copied, what is shared, and what can race.
- API boundaries: where we validate, where we propagate errors, and where we normalize output.
- Algorithmic tradeoffs in this lesson (latency, throughput, memory, and complexity).

## End-State Understanding

- Diagnose and implement Http Server Graceful patterns without relying on hidden framework behavior.
- Explain memory and concurrency implications of the final implementation choices.
- Compare learner code and reference code by invariants, not only by syntax.

## Why This Lesson Exists Here

Problem statement:
This lesson turns the previous module's concepts into a reusable engineering pattern for http server graceful.

At this point in the arc:
Lesson 09 introduces a sharper systems concern so later modules can assume this mental model is stable.

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


This project addresses a critical aspect of production-ready services: **graceful shutdown**. You will learn how to build an HTTP server that doesn't just abruptly stop, but instead, intelligently finishes its work before exiting. This is essential for zero-downtime deployments and preventing data corruption.

## CS Connection

- Memory layout drives behavior: variables store values, and some values are addresses into other storage.
- Go is pass-by-value, including pointers; copying a pointer value copies an address, not the pointee.
- Correctness depends on understanding copying versus aliasing (`*T`, slices, maps, and channels) and enforcing synchronization when concurrent access exists.

## End-State Understanding

- Implement the exercise with explicit reasoning about correctness, edge cases, and error paths.

## What You'll Learn

- Why graceful shutdown is critical for production services.
- How to handle operating system signals (like `Ctrl+C`) to trigger a shutdown.
- How to use Go's `http.Server` and its `Shutdown()` method.
- How to use `context.WithTimeout` to enforce a deadline on the shutdown process.

## The Challenge: Abrupt vs. Graceful Shutdown

Imagine a user is uploading a file when you deploy a new version of your server. The old server process is killed, the user's connection is severed, and their upload fails. This is an abrupt shutdown.

A **graceful shutdown** is much better:
1.  The server stops accepting *new* connections.
2.  It waits for all *in-flight* requests to complete.
3.  Once all work is done, it exits cleanly.

## Core Concepts

### Listening for OS Signals
How do we tell our running server to shut down? By listening for signals from the operating system. Go's `os/signal` package lets us receive these signals on a channel. We will listen for `SIGINT` (sent by `Ctrl+C`) and `SIGTERM` (sent by deployment systems like Kubernetes).

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

// This line will block until a signal is received on the quit channel.
<-quit
```

### The `http.Server` and `Shutdown()`
Instead of the simple `http.ListenAndServe()`, we create a `&http.Server{}` struct. This gives us an object we can control.

The most important method on this object is `Shutdown(ctx)`. When called, it:
1.  Immediately stops the server from accepting new connections.
2.  Waits for existing requests to finish.
3.  Takes a `context.Context` argument to set a timeout on the shutdown process itself, preventing it from waiting forever on a stuck request.

### 🚨 Deep Dive: `&http.Server{}`
When we create our server with `srv := &http.Server{...}`, the `&` operator gives us a **pointer** to the `http.Server` struct. We are not passing the struct itself around, but a reference to it. This is essential because methods like `ListenAndServe()` and `Shutdown()` need to modify the internal state of that single server object. If we passed a copy, these methods would be operating on a copy, and the changes would not affect our actual running server.

## Your Task

Your task is to implement the `RunGracefulServer` function in `internal/httpservergraceful/exercise.go`.

The logic should follow this sequence:
1.  Create a channel to listen for OS signals.
2.  Create a channel to receive errors from the `ListenAndServe` call.
3.  Start the server in a **goroutine**. If `srv.ListenAndServe()` returns an error (other than `http.ErrServerClosed`), it should be sent to the error channel.
4.  Block until an OS signal is received on your signal channel OR an error is received on the error channel. A `select` statement is perfect for this.
5.  Once a signal is received, create a `context.WithTimeout` to give the shutdown a deadline (e.g., 5 seconds).
6.  Call `srv.Shutdown()` with the timeout context and return its error.
7.  If a fatal server error was received from the goroutine, return that error.

## How to Verify Your Work

This lesson is interactive!

1.  **Start the server from the lesson directory:**
    ```bash
    go run ./cmd/app/main.go
    ```
    The server will be running on `localhost:8080`.

2.  **Simulate a long-running request:**
    Open a **new terminal** and use `curl` to make a request to the `/slow` endpoint. This endpoint takes 5 seconds to respond.
    ```bash
    curl localhost:8080/slow
    ```

3.  **Trigger a graceful shutdown:**
    While the `curl` command is still running, go back to the server's terminal and press **`Ctrl+C`**.

4.  **Observe the behavior:**
    - The server will print that it received a shutdown signal.
    - The server will **not** exit immediately. It will wait.
    - After 5 seconds, your `curl` command will complete successfully.
    - Only then will the server print that it has exited gracefully.

There are no automated tests for this lesson, as the goal is to observe the graceful shutdown behavior yourself.

## Related Lessons
- Previous: `minis/08-http-client-retries`
- Next: `minis/10-grpc-telemetry-service`
