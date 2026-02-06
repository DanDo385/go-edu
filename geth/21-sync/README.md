# 21: Sync Progress Inspection

## Core Concepts

- Problem framing for Sync Progress Inspection: what state we need, what invariants we must keep.
- Value vs pointer behavior in this lesson's APIs and data structures.
- Error-path design: fail fast at boundaries, keep results deterministic.

## CS Connection

- Memory ownership: distinguish copied values from aliased references (`*T`, slices, maps).
- API contracts: define what can be mutated and by whom.
- Runtime behavior: how failures, retries, and concurrency impact correctness.

## End-State Understanding

- Explain why this lesson exists in the geth arc and what gap it closes.
- Implement `exercise.go` and justify design choices against `solution.reference.go`.
- Reason about memory/pointer effects in every non-trivial step.

## Why This Lesson Now

This starts the ops triad: sync, peers, mempool.

Problem statement:
Determine whether a node is syncing and capture progress safely.

## Step-by-Step Build Path

### Step 1: Problem This Step Solves
Establish the minimum input validation and boundary checks so invalid state fails early.

### Step 2: Why This Approach
Use small, explicit operations that map 1:1 to the underlying RPC or data-model behavior.

### Step 3: Memory / Pointer Impact
`SyncProgress` is pointer-returned; copy before exposing if callers might mutate shared state.

### Step 4: What Changed
Return a stable `Result` snapshot that callers can inspect without mutating upstream/internal state.

## Pointer and Indirection Checklist (`*` and `&`)

- `*` in a type means pointer type; it does not dereference by itself.
- `&` creates an address value; Go remains pass-by-value.
- If a pointer/slice/map is returned, document whether caller mutation is allowed.
- When mutation is not allowed, copy before return (see `docs/MEMORY_POINTERS_PRIMER.md`).

## Verify

```bash
go test -v ./internal/...
go test -tags=reference -v ./internal/...
```

## Related Lessons

- Previous: follow the preceding lesson in `geth/` ordering.
- Next: continue to the next lesson in `geth/` ordering.
