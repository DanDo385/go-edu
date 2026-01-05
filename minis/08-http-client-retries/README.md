# 08: HTTP Client with Retries & Backoff

This project tackles a fundamental challenge of distributed systems: network unreliability. You will build a robust HTTP client that can intelligently handle transient failures by automatically retrying requests. You'll learn about and implement the **exponential backoff with jitter** algorithm, a gold-standard technique for building resilient, production-grade network clients.

## Table of Contents

- [Learning Objectives](#learning-objectives)
- [The Big Picture: The Unreliable Network](#the-big-picture-the-unreliable-network)
- [First Principles: Retry Strategies](#first-principles-retry-strategies)
  - [Exponential Backoff with Jitter](#exponential-backoff-with-jitter)
- [Project Structure](#project-structure)
- [Key Concepts in This Project](#key-concepts-in-this-project)
  - [Custom `http.Client` and `http.RoundTripper`](#custom-httpclient-and-httproundtripper)
  - [The Decorator Pattern for Middleware](#the-decorator-pattern-for-middleware)
  - [Critical Error Classification](#critical-error-classification)
- [Progression: Building for Failure](#progression-building-for-failure)
- [How to Run and Test](#how-to-run-and-test)
- [Key Takeaways](#key-takeaways)
- [Further Reading](#further-reading)

## Learning Objectives

By the end of this project, you will be able to:

-   **Explain why simple HTTP requests fail** in distributed systems.
-   **Implement advanced retry logic**, including exponential backoff and jitter.
-   **Customize `http.Client`** using the `http.RoundTripper` interface to create middleware.
-   **Apply the Decorator design pattern** in a practical, real-world scenario.
-   **Classify errors** to determine whether a failed request should be retried.
-   **Use `context.Context`** to handle deadlines and cancellations across retries.

## The Big Picture: The Unreliable Network

One of the "Fallacies of Distributed Computing" is the assumption that *the network is reliable*. It is not. In any system where services communicate over a network, requests can and will fail for countless transient reasons: a temporary network glitch, a server being briefly overloaded, a load balancer restarting, etc.

A naive client will treat any failure as permanent, immediately giving up and returning an error. This creates a fragile system. A robust, **resilient** client understands that some failures are temporary and will try again. This simple act of retrying transforms a brittle component into a stable one, dramatically improving the overall reliability of the entire system.

## First Principles: Retry Strategies

-   **No Retries**: The default. Simple, but fragile.
-   **Fixed Interval**: Always wait `N` seconds between retries. Better, but can lead to a "thundering herd" problem where many clients retry in sync, overwhelming the server again.
-   **Exponential Backoff**: Increase the wait time after each failure (e.g., 1s, 2s, 4s, 8s). This gives the server time to recover. This is the core of our strategy.

### Exponential Backoff with Jitter

This is the industry-standard algorithm. It has two components:

1.  **Exponential Backoff**: The base delay between retries doubles with each failed attempt. This prevents a client from hammering a struggling server.
2.  **Jitter**: A small, random amount of time is added to (or subtracted from) the delay. This is crucial for breaking up the "thundering herd." If multiple clients experience a failure at the same time, jitter ensures their retries are staggered, preventing them from all retrying in a synchronized wave.

**Example Flow**:
- Request 1 fails. Wait `1s + random(200ms)`.
- Request 2 fails. Wait `2s + random(500ms)`.
- Request 3 fails. Wait `4s + random(1000ms)`.
- Request 4 succeeds.

## Project Structure

```
.
├── cmd/
│   └── app/
│       └── main.go       # A simple CLI to test the retry client.
└── internal/
    └── client/
        └── retry.go      # Your implementation of the retry logic.
```
-   **`cmd/app`**: A small application that takes a URL as an argument and uses your retry client to fetch it.
-   **`internal/client`**: Contains the core implementation of the custom `http.RoundTripper` and the retry algorithm.

## Key Concepts in This Project

### Custom `http.Client` and `http.RoundTripper`

The `net/http` package is incredibly flexible. The `http.Client` has a `Transport` field, which must satisfy the `http.RoundTripper` interface. This interface has just one method: `RoundTrip(*http.Request) (*http.Response, error)`.

The `http.DefaultTransport` is the standard implementation that actually makes the network request. By creating our own `http.RoundTripper`, we can wrap the default transport and add behavior before or after the request is made. This is the perfect place to implement our retry logic.

### The Decorator Pattern for Middleware

Our `retryRoundTripper` will be a **Decorator**. It holds a reference to another `http.RoundTripper` (the "next" one in the chain) and adds its own logic around the call to `next.RoundTrip()`.

```
        Your Code             Your Wrapper              Go's Default
+---------------------+   +---------------------+   +------------------------+
| http.Client.Get()   |-->| retryRoundTripper   |-->| http.DefaultTransport  |
|                     |   |   .RoundTrip()      |   |   .RoundTrip()         |
+---------------------+   | (Contains your      |   | (Makes actual network |
                          |  retry loop)        |   |  request)              |
                          +---------------------+   +------------------------+
```

This is a powerful "middleware" pattern. You can chain multiple `RoundTripper` implementations together to add logging, caching, authentication, and more, all without modifying the core application logic.

### Critical Error Classification

The most important part of a retry strategy is knowing **when not to retry**. Retrying on a permanent error wastes time and resources.

-   **Retryable Errors (Transient)**:
    -   Network errors (e.g., `net.DNSError`, `net.OpError`). The connection was dropped or couldn't be established.
    -   Server is temporarily unavailable: HTTP `502 Bad Gateway`, `503 Service Unavailable`, `504 Gateway Timeout`.
    -   Server is busy: HTTP `429 Too Many Requests`.

-   **Non-Retryable Errors (Permanent)**:
    -   Client-side errors: HTTP `4xx` status codes (e.g., `400 Bad Request`, `401 Unauthorized`, `404 Not Found`). The request itself is flawed; retrying the exact same request will always fail.
    -   Server-side bugs: HTTP `500 Internal Server Error`. While sometimes transient, it often indicates a bug that won't be fixed by a simple retry. It's safer to fail fast.

Your implementation will need a function that inspects the returned error and HTTP response code to make this critical decision.

## Progression: Building for Failure

This project builds directly on your understanding of interfaces (**Project 05**) and concurrency context (**Project 06**). You are learning to program defensively, anticipating failure and building systems that can gracefully recover from it. This mindset is essential for creating production-ready applications. The decorator pattern used here is also a general-purpose software design principle you will use again.

## How to Run and Test

```bash
# Build the test application
go build -o client ./cmd/app

# Run it against a URL that will likely succeed
./client http://google.com

# Run it against a URL that will fail, to see retries in action
# (This URL will time out, triggering the retry logic)
./client http://example.com:81

# Run the provided tests
go test -v ./...
```

## Key Takeaways

-   **The network is not reliable; design for failure.**
-   **Exponential backoff with jitter** is the gold standard for retry logic.
-   The `http.RoundTripper` interface and the **decorator pattern** provide a powerful way to create HTTP client middleware.
-   **Error classification** is the most critical part of a retry strategy; do not retry permanent errors.
-   `context.Context` is vital for ensuring long-running operations with retries can be cancelled.

## Further Reading

-   [**The Fallacies of Distributed Computing**](https://en.wikipedia.org/wiki/Fallacies_of_distributed_computing)
-   [**AWS Architecture Blog: Exponential Backoff And Jitter**](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/)
-   [**Package `net/http`**](https://pkg.go.dev/net/http): Pay special attention to the `Client` and `RoundTripper` types.
-   [**Stripe Blog: API client resilience**](https://stripe.com/blog/api-client-resilience): A great real-world article on why client-side retry logic is important.