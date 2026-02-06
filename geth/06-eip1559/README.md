# 06: Dynamic Fee Transactions (EIP-1559)

## Core Concepts

- Problem framing for Dynamic Fee Transactions (EIP-1559): what state we need, what invariants we must keep.
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

After legacy tx flow, learners upgrade to modern fee markets and cap calculations.

Problem statement:
Construct EIP-1559 transactions from base fee + tip strategy.

## Step-by-Step Build Path

### Step 1: Problem This Step Solves
Establish the minimum input validation and boundary checks so invalid state fails early.

### Step 2: Why This Approach
Use small, explicit operations that map 1:1 to the underlying RPC or data-model behavior.

### Step 3: Memory / Pointer Impact
Base fee, tip, and fee cap are `*big.Int`; show before/after copies to separate local math from shared upstream pointers.

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
