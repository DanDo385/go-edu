# 19: Channels Basics

## Core Concepts

- Unbuffered vs buffered channels.
- Blocking semantics of send/receive.
- Close semantics and `range` over channels.

## CS Connection

- Channel values are copied by value but reference shared runtime queue/state.
- Correctness depends on synchronization order, not line-by-line source order.
- Deadlocks happen when send/receive dependencies cannot be satisfied.

## End-State Understanding

- Predict when operations block.
- Structure producer/consumer flows that terminate cleanly.
- Explain why channel close is a broadcast of completion, not data.

## Why This Lesson Now

After goroutine basics, you now need deterministic communication and synchronization primitives.

Problem statement:
Coordinate concurrent work without races or ad-hoc polling loops.

## Step-by-Step Build Path

### Step 1: Problem This Step Solves
Define a minimal producer/consumer pair with explicit send/receive points.

### Step 2: Why This Approach
Channels encode synchronization in the type and operation semantics.

### Step 3: Memory / Pointer Impact
Sending values copies payloads, but sending pointers copies addresses. If you send `*T`, goroutines can alias and mutate shared pointee state.

### Step 4: What Changed
You can now choose value vs pointer payloads intentionally based on ownership needs.

## Pointer and Indirection Checklist (`*` and `&`)

- `chan T` sends copied values of `T`.
- `chan *T` sends copied addresses; pointee remains shared.
- `&x` before send means receiver can mutate caller-owned memory through that pointer.
- Recheck `docs/MEMORY_POINTERS_PRIMER.md` when channels carry references.

## Verify

```bash
go test -v ./...
go run ./cmd/app/main.go
```

## Related Lessons

- Previous: `minis/18-goroutines-1M-demo`
- Next: `minis/20-select-fanin-fanout`
