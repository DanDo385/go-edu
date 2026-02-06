# 37: Http Middleware Chain

## Core Concepts

- The concrete problem in Http Middleware Chain and the correctness invariants it depends on.
- How values, pointers, slices, maps, or channels behave in this module's runtime path.
- Why this lesson's implementation pattern is the right next step in the learning arc.

## CS Connection

- Data representation and state transitions: what is copied, what is shared, and what can race.
- API boundaries: where we validate, where we propagate errors, and where we normalize output.
- Algorithmic tradeoffs in this lesson (latency, throughput, memory, and complexity).

## End-State Understanding

- Diagnose and implement Http Middleware Chain patterns without relying on hidden framework behavior.
- Explain memory and concurrency implications of the final implementation choices.
- Compare learner code and reference code by invariants, not only by syntax.

## Why This Lesson Exists Here

Problem statement:
This lesson turns the previous module's concepts into a reusable engineering pattern for http middleware chain.

At this point in the arc:
Lesson 37 introduces a sharper systems concern so later modules can assume this mental model is stable.

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


**HTTP Middleware Chain**

Build composable HTTP middleware.

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

- Middleware pattern
- Function composition
- Request/response modification
- Handler wrapping

## Functions to Implement

| Function | Description |
|----------|-------------|
| Build middleware chain | Composable handlers |

## Project Structure

```
37-http-middleware-chain/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
├── internal/httpmiddlewarechain/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/37-http-middleware-chain

# Start server with middleware
go run ./cmd/app/main.go --port 8080

# Example run
go run ./cmd/app/main.go
```

## Quick Copy & Paste

```bash
# Start server
go run ./cmd/app/main.go --port 8080

# Test with curl (see middleware effects)
curl -v http://localhost:8080/

# Example run
go run ./cmd/app/main.go
```

## Key Concepts

1. **func(http.Handler) http.Handler**: Middleware signature
2. **Chaining**: logging(auth(handler))
3. **Before/After**: Pre and post processing
4. **Context Values**: Pass data through chain

## Next Steps

After completing this exercise, proceed to `minis/38-config-loader-env-yaml`.

## Related Lessons
- Previous: `minis/36-caching-reverse-proxy`
- Next: `minis/38-config-loader-env-yaml`
