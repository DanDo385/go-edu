# 16: Context Cancellation and Timeouts

## Core Concepts

- Cancellation trees with `context.Context`.
- Timeout/deadline propagation across call chains.
- Deterministic cleanup when work is canceled.

## CS Connection

- `Context` values are copied by value, but carry shared cancellation state internally.
- Concurrency correctness depends on all goroutines observing the same cancellation signal.
- Resource safety comes from bounded lifetimes and cooperative shutdown.

## End-State Understanding

- Design APIs that accept `context.Context` first.
- Enforce time bounds and return fast on cancellation.
- Explain how cancellation prevents leaked goroutines and stuck I/O.

## Why This Lesson Now

Before deeper concurrency modules, you need a standard control plane for stopping work safely.

Problem statement:
Avoid hanging operations and leaked goroutines by propagating cancellation and deadlines everywhere.

## Step-by-Step Build Path

### Step 1: Problem This Step Solves
Introduce a root context and cancellation boundary for each top-level operation.

### Step 2: Why This Approach
A single propagation mechanism keeps timeouts and cancellation semantics consistent across layers.

### Step 3: Memory / Pointer Impact
Even without explicit `*` syntax, `Context` values share internal state. Two goroutines with copied context values still observe the same cancellation event via shared internals.

### Step 4: What Changed
Long-running operations now terminate predictably, and upstream callers control lifetime.

## Pointer and Indirection Checklist (`*` and `&`)

- `context.Context` is an interface value; copying it does not duplicate cancellation internals.
- Avoid storing mutable pointer state in context values.
- If a context carries a pointer in `WithValue`, document ownership and immutability explicitly.
- Review `docs/MEMORY_POINTERS_PRIMER.md` for aliasing rules.

## Verify

```bash
go test -v ./...
go run ./minis/16-context-cancellation-timeouts/cmd/app/main.go
```

## Related Lessons

- Previous: `minis/15-error-wrapping-sentinel-errors`
- Next: `minis/17-file-streaming-bufio`
