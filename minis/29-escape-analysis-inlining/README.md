# 29: Escape Analysis and Inlining

## Core Concepts

- Stack vs heap placement from compiler escape analysis.
- Inlining decisions and call overhead.
- How `&` and returned pointers influence allocation.

## CS Connection

- Allocation location is a compiler decision based on lifetime and aliasing.
- Returning references often extends lifetime beyond stack frame scope.
- Optimization visibility requires reading compiler diagnostics, not guessing.

## End-State Understanding

- Interpret `-gcflags=-m` output for escape and inline decisions.
- Refactor code to reduce unnecessary heap allocation when appropriate.
- Explain tradeoffs between readability and micro-optimizations.

## Why This Lesson Now

After concurrency/perf profiling modules, this lesson connects code shape to compiler/runtime cost.

Problem statement:
Understand why seemingly small API choices produce large allocation differences.

## Step-by-Step Build Path

### Step 1: Problem This Step Solves
Capture baseline compiler diagnostics for current code paths.

### Step 2: Why This Approach
Compiler output provides ground truth on escapes and inlining opportunities.

### Step 3: Memory / Pointer Impact
Using `&x`, returning `*T`, storing values in interfaces, or capturing variables in closures can force escape to heap. Copying plain values often stays stack-local when lifetime permits.

### Step 4: What Changed
You can now justify allocation/perf behavior using concrete compiler evidence.

## Pointer and Indirection Checklist (`*` and `&`)

- `&x` creates an address value; whether `x` escapes depends on usage.
- `*T` return values often imply longer lifetime and potential heap allocation.
- Pointer receiver methods do not automatically force heap allocation; context matters.
- Use `docs/MEMORY_POINTERS_PRIMER.md` plus `go build -gcflags=-m` together.

## Verify

```bash
go build -gcflags='-m' ./...
go test -v ./...
go run ./minis/29-escape-analysis-inlining/cmd/app/main.go
```

## Related Lessons

- Previous: `minis/28-pprof-cpu-mem-benchmarks`
- Next: `minis/30-build-tags-conditional-compilation`
