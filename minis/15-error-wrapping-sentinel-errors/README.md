# 15: Error Wrapping and Sentinel Errors

This project explores the nuances of idiomatic error handling in modern Go (1.13+). You'll learn that errors are not just failures, but are values that can tell a story. We will cover the three main patterns—sentinel errors, custom error types, and error wrapping—and the tools Go provides (`errors.Is`, `errors.As`, `%w`) to work with them.

## Table of Contents

- [Learning Objectives](#learning-objectives)
- [The Big Picture: Errors as a Story](#the-big-picture-errors-as-a-story)
- [First Principles: Errors are Values](#first-principles-errors-are-values)
- [Project Structure](#project-structure)
- [The Three Patterns of Error Handling](#the-three-patterns-of-error-handling)
  - [1. Sentinel Errors](#1-sentinel-errors)
  - [2. Custom Error Types](#2-custom-error-types)
  - [3. Opaque Errors with Wrapping](#3-opaque-errors-with-wrapping)
- [The Modern Go Error Toolkit: `Is`, `As`, and `%w`](#the-modern-go-error-toolkit-is-as-and-w)
  - [`fmt.Errorf` with `%w`](#fmterrorf-with-w)
  - [`errors.Is`](#errorsis)
  - [`errors.As`](#errorsas)
- [Progression: Building Robust Functions](#progression-building-robust-functions)
- [How to Run](#how-to-run)
- [Key Takeaways](#key-takeaways)
- [Further Reading](#further-reading)

## Learning Objectives

By the end of this project, you will be able to:

-   **Explain the "errors are values" philosophy** in Go.
-   **Choose the appropriate error handling strategy**: sentinel, custom type, or opaque wrapping.
-   **Add context to errors** using `fmt.Errorf` and the `%w` verb to create error chains.
-   **Inspect error chains** using `errors.Is` to check for sentinel values.
-   **Inspect error chains** using `errors.As` to check for and extract custom error types.
-   **Write code that is more debuggable and maintainable** through better error handling.

## The Big Picture: Errors as a Story

Good error handling tells a story. When a user's request fails, a simple "something went wrong" is useless. A good error provides a trail of breadcrumbs leading back to the root cause.

-   **High Level (API Layer)**: "failed to process user checkout for user 123"
-   **Mid Level (Service Layer)**: "...because failed to debit from wallet"
-   **Low Level (Database Layer)**: "...because database connection failed"
-   **Root Cause (Network Layer)**: "...because connection refused"

Go's error wrapping mechanism is designed to build this chain of context, making debugging complex systems tractable.

## First Principles: Errors are Values

In Go, an error is any type that implements the simple `error` interface:
```go
type error interface {
    Error() string
}
```
This means errors are not special language constructs (like exceptions in Java or Python). They are just values that are returned from functions, assigned to variables, and can be inspected with normal `if` statements. This "errors are values" philosophy makes error handling explicit, clear, and robust.

## Project Structure

```
.
└── cmd/
    └── dev/
        └── main.go       # Demonstrates various error handling scenarios.
```
-   The `main.go` file contains a series of functions that simulate a multi-layered application (e.g., `apiCall` -> `serviceCall` -> `dbCall`). Each function demonstrates a different way of handling, returning, and wrapping errors.

## The Three Patterns of Error Handling

### 1. Sentinel Errors
A sentinel error is a pre-defined, public error variable that signals a specific, constant condition. The most famous example is `io.EOF`.

-   **Definition**: `var ErrNotFound = errors.New("not found")`
-   **Usage**: Callers check for it using a simple equality check (`==`) or, more robustly, `errors.Is`.
-   **When to Use**: For fixed, unambiguous conditions where the caller needs to take a specific branch of logic (e.g., "if not found, create a new one").
-   **Downside**: Creates a public dependency. If many packages depend on your sentinel error, it can become difficult to change.

### 2. Custom Error Types
For when you need to attach more context to an error than just a string. You can define your own struct that implements the `error` interface.

-   **Definition**:
    ```go
    type OpError struct {
        Op    string
        Code  int
        Cause error
    }
    func (e *OpError) Error() string { /* ... */ }
    ```
-   **Usage**: Callers inspect it with `errors.As`.
-   **When to Use**: When you need to provide structured, machine-readable context with your error (e.g., an operation name, an HTTP status code, a retry-ability flag).

### 3. Opaque Errors with Wrapping
This is the most common and generally recommended approach. You treat errors from functions you call as "opaque" (you don't inspect their type), and you simply add your own context to them as you return them up the call stack.

-   **Usage**: `return fmt.Errorf("service layer failed: %w", err)`
-   **When to Use**: Almost always. It provides a clear chain of responsibility without creating tight coupling between packages.

## The Modern Go Error Toolkit: `Is`, `As`, and `%w`

### `fmt.Errorf` with `%w`
The `%w` verb, added in Go 1.13, is the key to error wrapping. It tells `fmt.Errorf` to create a new error that wraps the original.

```
err := dbCall() // returns a low-level error
wrappedErr := fmt.Errorf("context: %w", err)

// wrappedErr now contains a reference to `err`.
```
**Error Chain Diagram:**
```
[ "failed to checkout" ] --wraps--> [ "failed to debit wallet" ] --wraps--> [ "connection refused" ]
```

### `errors.Is`
`errors.Is(err, target)` checks if any error in the chain matches a specific `target` (a sentinel error). It walks the chain, comparing each error to `target`.

```go
// db.Query() might return a wrapped error.
err := db.Query()

// errors.Is will walk the chain to find the root cause.
if errors.Is(err, sql.ErrNoRows) {
    // Handle "not found" case
}
```

### `errors.As`
`errors.As(err, &target)` checks if any error in the chain can be "assigned to" `target` (a custom error type). If it finds a match, it sets `target` to that error value and returns `true`.

```go
err := apiCall()

var opErr *OpError
if errors.As(err, &opErr) {
    // It's an OpError! We can now inspect its fields.
    fmt.Println("Operation:", opErr.Op)
    fmt.Println("Code:", opErr.Code)
}
```

## Progression: Building Robust Functions

This project builds on every previous project. Good error handling is a non-negotiable part of any function that can fail, from string parsing (**Project 01**) to HTTP requests (**Project 08**) to database calls. The patterns learned here are universally applicable and represent the current best practice for writing Go code that is easy to debug and resilient in production.

## How to Run
```bash
go run ./cmd/dev/main.go
```
The program will run several scenarios. Read the code for each scenario in `main.go`, predict what the output will be (especially which `if` blocks will be entered), and then run the program to verify your understanding.

## Key Takeaways

-   **Errors are values** that tell a story about what went wrong.
-   **Add context** to errors by wrapping them with `fmt.Errorf` and `%w`.
-   **Use `errors.Is` to check for sentinel error values** (e.g., `io.EOF`) in an error chain.
-   **Use `errors.As` to inspect and extract custom error types** from an error chain.
-   Prefer **opaque errors with wrapping** over sentinel errors to reduce coupling between packages.
-   Never just return `err`, always add context unless you are at the boundary of your system.

## Further Reading

-   [**Go Blog: Working with Errors in Go 1.13**](https://go.dev/blog/go1.13-errors) - The definitive guide to the modern error handling toolkit.
-   [**Go by Example: Errors**](https://gobyexample.com/errors)
-   [**Effective Go: Errors**](https://go.dev/doc/effective_go#errors)
-   [**Video: GopherCon 2019: Don't Just Check Errors, Handle Them Gracefully**](https://www.youtube.com/watch?v=8D3qB-dI_c8) by Dave Chaney.