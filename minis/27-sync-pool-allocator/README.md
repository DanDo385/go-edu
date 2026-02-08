# 27: `sync.Pool` Allocator

## Core Concepts

- Object reuse with `sync.Pool`.
- Allocation pressure and GC tradeoffs.
- Ownership rules for borrowed vs returned objects.

## CS Connection

- Pools trade allocation cost for stricter lifetime discipline.
- A pooled pointer can be reused by unrelated callers after `Put`.
- Correctness depends on resetting state before reuse.

## End-State Understanding

- Use `sync.Pool` only for short-lived, high-churn objects.
- Define clear borrow/return contracts.
- Avoid use-after-Put bugs with pointer objects.

## Why This Lesson Now

After mutex/atomic/once patterns, this module addresses memory churn in concurrent services.

Problem statement:
Reduce allocator and GC overhead without introducing shared-state corruption.

## Step-by-Step Build Path

### Step 1: Problem This Step Solves
Measure object churn and identify reusable object shape.

### Step 2: Why This Approach
`sync.Pool` gives low-friction reuse integrated with runtime GC behavior.

### Step 3: Memory / Pointer Impact
Pools usually store `*T`. `Get` returns a pointer that may reference old data; you must reset it. After `Put`, the pointer is no longer exclusively owned by the caller.

### Step 4: What Changed
Alloc-heavy paths can reuse buffers/objects with explicit ownership boundaries.

## Pointer and Indirection Checklist (`*` and `&`)

- `pool.Get()` returning `*T` means caller holds an address to mutable shared memory.
- Never keep references after `Put`; that is a use-after-release bug.
- Always zero/reset pointee fields before reuse.
- See `docs/MEMORY_POINTERS_PRIMER.md` for aliasing and lifetime rules.

## Verify

```bash
go test -v ./...
go run ./minis/27-sync-pool-allocator/cmd/app/main.go
```

## Related Lessons

- Previous: `minis/26-sync-once-singleton`
- Next: `minis/28-pprof-cpu-mem-benchmarks`
