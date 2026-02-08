# 45: P2p Gossip Mock Network

## Core Concepts

- The concrete problem in P2p Gossip Mock Network and the correctness invariants it depends on.
- How values, pointers, slices, maps, or channels behave in this module's runtime path.
- Why this lesson's implementation pattern is the right next step in the learning arc.

## CS Connection

- Data representation and state transitions: what is copied, what is shared, and what can race.
- API boundaries: where we validate, where we propagate errors, and where we normalize output.
- Algorithmic tradeoffs in this lesson (latency, throughput, memory, and complexity).

## End-State Understanding

- Diagnose and implement P2p Gossip Mock Network patterns without relying on hidden framework behavior.
- Explain memory and concurrency implications of the final implementation choices.
- Compare learner code and reference code by invariants, not only by syntax.

## Why This Lesson Exists Here

Problem statement:
This lesson turns the previous module's concepts into a reusable engineering pattern for p2p gossip mock network.

At this point in the arc:
Lesson 45 introduces a sharper systems concern so later modules can assume this mental model is stable.

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


**P2P Gossip Network**

Simulate a gossip protocol for message propagation.

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

- Gossip protocol basics
- Message propagation
- Network simulation
- Deduplication

## Functions to Implement

| Function | Description |
|----------|-------------|
| Broadcast message | Gossip to peers |
| Handle receive | Process and forward |
| Deduplicate | Track seen messages |

## Project Structure

```
45-p2p-gossip-mock-network/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
├── internal/p2pgossipmocknetwork/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/45-p2p-gossip-mock-network

# Simulate network with 10 nodes
go run ./minis/45-p2p-gossip-mock-network/cmd/app/main.go --nodes 10 --fanout 3

# Example run
go run ./minis/45-p2p-gossip-mock-network/cmd/app/main.go
```

## Quick Copy & Paste

```bash
# Simulate gossip
go run ./minis/45-p2p-gossip-mock-network/cmd/app/main.go --nodes 10 --fanout 3

# Larger network
go run ./minis/45-p2p-gossip-mock-network/cmd/app/main.go --nodes 100 --fanout 5

# Example run
go run ./minis/45-p2p-gossip-mock-network/cmd/app/main.go
```

## Key Concepts

1. **Gossip**: Random peer selection
2. **Fanout**: Peers to forward to
3. **Seen Set**: Prevent re-broadcasting
4. **Exponential Spread**: O(log n) rounds

## Next Steps

After completing this exercise, proceed to `minis/46-generics-map-reduce`.

## Related Lessons
- Previous: `minis/44-mempool-in-memory`
- Next: `minis/46-generics-map-reduce`
