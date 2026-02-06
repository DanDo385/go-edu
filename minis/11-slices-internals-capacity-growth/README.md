# Lesson: Slices Internals & Capacity Growth

## Core Concepts

- Slice header vs underlying array: `ptr`, `len`, `cap`.
- `append` growth behavior and when reallocation occurs.
- Aliasing rules when multiple slices share the same backing array.

## CS Connection

- A slice value is a small descriptor copied by value.
- Copying the descriptor does not copy the array it points to.
- Performance comes from understanding when allocation/copy happens versus in-place writes.

## End-State Understanding

- Predict when `append` mutates in-place vs allocates a new array.
- Explain why two slices can unexpectedly affect each other.
- Use preallocation intentionally to reduce allocation and copy overhead.

## Why This Lesson Now

You have used slices in earlier lessons; now you need the memory model behind them before pointer-heavy and concurrent modules.

Problem statement:
Understand slice aliasing and growth so later code does not accidentally corrupt shared state or waste allocations.

## Step-by-Step Build Path

### Step 1: Problem This Step Solves
Make slice state visible (`len`, `cap`, contents) after each append.

### Step 2: Why This Approach
Observing transitions directly builds an intuition that static rules alone do not.

### Step 3: Memory / Pointer Impact
The slice header contains a pointer-like reference to an array. Copying the header copies that reference. If capacity remains, writes alias the same array; after growth, a new array is allocated and future writes diverge.

### Step 4: What Changed
You can now reason from memory state, not guesses, when writing append-heavy code.

## Pointer and Indirection Checklist (`*` and `&`)

- Slices are not pointers syntactically, but they carry pointer semantics through their header.
- If you take `&s[i]`, that address is valid only while the element remains in the same backing array.
- After reallocation, old element addresses may point to stale arrays.
- Review `docs/MEMORY_POINTERS_PRIMER.md` before pointer-based collections.

## Verify

```bash
go run ./cmd/app/main.go
```

Watch where capacity jumps; those are allocation/copy points.

## Related Lessons

- Previous: `minis/10-grpc-telemetry-service`
- Next: `minis/12-pointers-zero-values-nil-gotchas`
