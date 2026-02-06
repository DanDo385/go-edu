# 30: Build Tags Conditional Compilation

## Core Concepts

- The concrete problem in Build Tags Conditional Compilation and the correctness invariants it depends on.
- How values, pointers, slices, maps, or channels behave in this module's runtime path.
- Why this lesson's implementation pattern is the right next step in the learning arc.

## CS Connection

- Data representation and state transitions: what is copied, what is shared, and what can race.
- API boundaries: where we validate, where we propagate errors, and where we normalize output.
- Algorithmic tradeoffs in this lesson (latency, throughput, memory, and complexity).

## End-State Understanding

- Diagnose and implement Build Tags Conditional Compilation patterns without relying on hidden framework behavior.
- Explain memory and concurrency implications of the final implementation choices.
- Compare learner code and reference code by invariants, not only by syntax.

## Why This Lesson Exists Here

Problem statement:
This lesson turns the previous module's concepts into a reusable engineering pattern for build tags conditional compilation.

At this point in the arc:
Lesson 30 introduces a sharper systems concern so later modules can assume this mental model is stable.

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


**Build Tags**

Use build tags for conditional compilation.

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

- Build constraint syntax
- Platform-specific code
- Feature flags
- Build tags vs file naming

## Functions to Implement

| Function | Description |
|----------|-------------|
| Platform-specific implementations | Conditional compilation |

## Project Structure

```
30-build-tags-conditional-compilation/
├── cmd/
│   ├── app/main.go      # CLI demonstration
├── internal/buildtagsconditionalcompilation/
│   ├── exercise.go         # Default implementation
│   ├── exercise_linux.go   # Linux-specific
│   ├── exercise_darwin.go  # macOS-specific
│   ├── exercise_test.go    # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/30-build-tags-conditional-compilation

# Build with tag
go build -tags=feature_x ./...

# Build for specific OS
GOOS=linux go build ./...
```

## Quick Copy & Paste

```bash
# Default build
go run ./cmd/app/main.go

# Build with custom tag
go build -tags=debug ./cmd/app/

# Cross-compile
GOOS=linux GOARCH=amd64 go build ./cmd/app/

# Example run
go run ./cmd/app/main.go
```

## Key Concepts

1. **//go:build tag**: New syntax (Go 1.17+)
2. **// +build tag**: Old syntax (still works)
3. **_linux.go**: File name suffix convention
4. **GOOS/GOARCH**: Cross-compilation

## Next Steps

After completing this exercise, proceed to `minis/31-static-file-server`.

## Related Lessons
- Previous: `minis/29-escape-analysis-inlining`
- Next: `minis/31-static-file-server`
