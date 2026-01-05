# 09: HTTP Server with Graceful Shutdown

This project addresses a critical aspect of production-ready services: **graceful shutdown**. You will learn how to build an HTTP server that doesn't just abruptly stop, but instead, intelligently finishes its work before exiting. This is essential for zero-downtime deployments, preventing data corruption, and providing a seamless user experience.

## Table of Contents

- [Learning Objectives](#learning-objectives)
- [The Big Picture: The Problem with Abrupt Shutdowns](#the-big-picture-the-problem-with-abrupt-shutdowns)
- [First Principles: Process Signals & Shutdown Phases](#first-principles-process-signals--shutdown-phases)
- [Project Structure](#project-structure)
- [Key Go Concepts in This Project](#key-go-concepts-in-this-project)
  - [The `http.Server` struct](#the-httpserver-struct)
  - [Listening for OS Signals](#listening-for-os-signals)
  - [The `server.Shutdown()` Method](#the-servershutdown-method)
  - [The Shutdown Sequence](#the-shutdown-sequence)
- [Progression: From Client to Server](#progression-from-client-to-server)
- [How to Run and Test](#how-to-run-and-test)
- [Key Takeaways](#key-takeaways)
- [Further Reading](#further-reading)

## Learning Objectives

By the end of this project, you will be able to:

-   **Build a robust HTTP server** using Go's `net/http` package.
-   **Explain why graceful shutdown is critical** for production services.
-   **Handle operating system signals** (like `SIGINT` and `SIGTERM`) to trigger a shutdown.
-   **Use `server.Shutdown()`** to allow in-flight requests to complete.
-   **Use `context.WithTimeout`** to enforce a deadline on the shutdown process.
-   **Write servers that support zero-downtime deployments**.

## The Big Picture: The Problem with Abrupt Shutdowns

Imagine a user is in the middle of submitting a form, a file upload, or a payment. A developer deploys a new version of the server. The deployment system simply kills the old server process. The user's connection is severed, their request is dropped, and they see an error. At best, this is a poor user experience; at worst, it can lead to inconsistent state or corrupted data.

**Graceful shutdown** solves this. When a shutdown is requested, the server:
1.  Immediately stops accepting *new* connections.
2.  Waits patiently for all *in-flight* requests to complete normally.
3.  Once all active requests are done, it closes its resources and exits.

This process is fundamental to modern deployment strategies like rolling updates and blue-green deployments, ensuring that updates are invisible to users.

## First Principles: Process Signals & Shutdown Phases

1.  **Process Signals**: How does an external process (like a developer's Ctrl+C or a Kubernetes orchestrator) tell the server to shut down? It sends a **signal**. The two most common signals are:
    *   `SIGINT` (Signal Interrupt): Sent when you press `Ctrl+C` in the terminal.
    *   `SIGTERM` (Signal Terminate): The standard, generic signal sent by process managers and orchestrators (like Docker, systemd, Kubernetes) to request a process to exit. Your application should always listen for this.

2.  **The Two Phases of Shutdown**:
    *   **Phase 1: Stop the World**: The server stops accepting new work. The listening port is closed. Any new connection attempts will be refused.
    *   **Phase 2: Drain and Wait**: The server enters a "lame duck" mode, waiting for all active requests to finish their lifecycle. Once the last request is served, the process can safely exit.

## Project Structure

```
.
├── cmd/
│   └── dev/
│       └── main.go       # The main entry point for the server.
└── internal/
    └── server/
        └── server.go     # The implementation of the graceful shutdown logic.
```
-   **`cmd/dev/main.go`**: Configures and runs the server.
-   **`internal/server/server.go`**: Contains the core logic for creating the server, listening for signals, and managing the graceful shutdown sequence.

## Key Go Concepts in This Project

### The `http.Server` struct
Instead of using the simple `http.ListenAndServe()`, we create an instance of `http.Server`. This gives us a server object that we can control.

```go
srv := &http.Server{
    Addr:    ":8080",
    Handler: myMux,
}
```

### Listening for OS Signals
Go's `os/signal` package allows us to receive OS signals on a channel.

```go
// Create a channel to receive signals.
quit := make(chan os.Signal, 1)

// Notify this channel for SIGINT and SIGTERM.
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

// Block until a signal is received.
<-quit
log.Println("Shutdown signal received...")
```

### The `server.Shutdown()` Method
This is the star of the show. When called, `server.Shutdown()` gracefully shuts down the server. It immediately closes the listener and then waits for active connections to drain.

It takes a `context.Context` as an argument. This is crucial for setting a **timeout** for the shutdown process itself. What if a connection is stuck and never finishes? We can't wait forever.

```go
// Create a context with a 5-second timeout for the shutdown.
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// This blocks until shutdown is complete or the context is cancelled.
if err := srv.Shutdown(ctx); err != nil {
    log.Fatalf("Server shutdown failed: %+v", err)
}
```

### The Shutdown Sequence

This visual timeline shows how the pieces fit together:

```
|
|--> Server starts, listening for requests on one goroutine
|--> Signal handler starts, blocking on `<-quit` on another goroutine
|
| (User sends requests, they are handled normally)
|
| (User presses Ctrl+C or Kubernetes sends SIGTERM)
|
|--> Signal handler unblocks
|--> Calls `srv.Shutdown(ctx)`
|--> `Shutdown` closes the server's listening port (no new connections)
|
| (Server continues processing active requests...)
|
|--> All active requests finish
|
|--> `srv.Shutdown()` returns `nil`
|--> Main function finishes, process exits cleanly
|
```

## Progression: From Client to Server

This project is the natural counterpart to **Project 08 (HTTP Client with Retries)**.
-   In Project 08, you learned to make your *client* resilient to server unavailability.
-   In Project 09, you learn to make your *server* resilient to restarts and deployments.

Together, they form a complete picture of robust client-server communication in a distributed environment. You are now thinking about the full lifecycle of a network service.

## How to Run and Test

1.  **Start the server:**
    ```bash
    go run ./cmd/dev/main.go
    ```
    The server will be running on `localhost:8080`.

2.  **Simulate a long-running request:**
    Open a new terminal and use `curl` to make a request to the `/slow` endpoint. This endpoint is designed to take 10 seconds to respond.
    ```bash
    curl localhost:8080/slow
    ```

3.  **Trigger a graceful shutdown:**
    While the `curl` command is still running, go back to the server's terminal and press `Ctrl+C`.

4.  **Observe the behavior:**
    -   The server will print "Shutdown signal received..." and "Waiting for in-flight requests to finish..."
    -   The server will **not** exit immediately.
    -   After 10 seconds, your `curl` command will complete successfully and print "Request complete!".
    -   The server will then print "Server exited gracefully" and shut down.

## Key Takeaways

-   **Production servers must shut down gracefully** to prevent data loss and user errors.
-   Use `os/signal.Notify` to listen for `SIGINT` and `SIGTERM`.
-   Use `http.Server` and its `Shutdown()` method to manage the server lifecycle.
-   **Always use a timeout context** with `Shutdown()` to prevent the server from hanging indefinitely.
-   Graceful shutdown is a key enabler of **zero-downtime deployments**.

## Further Reading

-   [**Package `net/http` (`Server` type)**](https://pkg.go.dev/net/http#Server)
-   [**Go Blog: Graceful shutdown of a server**](https://go.dev/blog/graceful-shutdown) (Note: This older post predates `server.Shutdown` but explains the core concepts).
-   [**Gist by gin-gonic Author: A complete example**](https://gist.github.com/gin-gonic/fc25f46452818b335384)
-   [**The Twelve-Factor App: Disposability**](https://12factor.net/disposability): The principle that services should be able to start and stop gracefully.