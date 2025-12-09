# geth-21-sync: Sync Progress Inspection

**Goal:** inspect sync progress and understand full/snap/light modes.

## Big Picture

Sync modes: full replays all blocks, snap downloads snapshots then heals, light fetches proofs on demand. `SyncProgress` reports current vs highest block and state sync counters. Nil means the node believes it is synced.

**Computer Science principle:** Nil as a sentinel value—the absence of sync progress indicates completion. This is a classic pattern where "nothing to report" is meaningful information.

## Learning Objectives

By the end of this module, you should be able to:

1. **Call SyncProgress** and interpret nil vs non-nil responses
2. **Differentiate sync modes** conceptually (full/snap/light)
3. **Spot stale nodes** by checking progress or head lag
4. **Understand sentinel values** as a design pattern for status reporting

## Prerequisites
- **Module 01 (01-stack):** RPC client basics, context handling
- **Module 20 (node info):** Understanding node metadata (helpful but not required)
- **Comfort with:** nil semantics in Go, blockchain basics

## Building on Previous Modules

### From Module 01 (01-stack)
- You learned RPC call patterns (call → error check → validate nil)
- Same pattern applies here, but nil has special meaning
- Input validation pattern repeats (context, client checks)

### New in this module
- **Nil as sentinel value:** nil = success (synced), not absence of data
- **Progress tracking:** Understanding current vs highest block counters
- **Operational health checks:** Determining if a node is production-ready

### Connection to Solidity-edu
- Storage/gas lessons: State size directly affects sync times
- Understanding why snap sync is faster (downloads state, skips replay)

## Real-World Analogy

### The Archive Download Analogy
- **Full sync:** Reading every page of history from the beginning (complete but slow)
- **Snap sync:** Downloading a snapshot of the current state, then healing missing data (fast)
- **Light sync:** Just the index cards, requesting full pages only when needed (minimal)

### The City Building Analogy
- **Full sync:** Building the city from scratch, one building at a time from 1900 to now
- **Snap sync:** Importing a 2023 snapshot, then filling in missing permits and records
- **Light sync:** Just a map with addresses, calling city hall for details as needed

## What You'll Build

In this module, you'll create a function that:
1. Takes an RPC client as input
2. Calls the SyncProgress method
3. Interprets nil (synced) vs non-nil (syncing) responses
4. Returns structured status information

**Key learning:** You'll understand how to check if an Ethereum node is ready for production use!

## Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite verifying correctness

## How to Run Tests

To run the tests for this module:

```bash
# From the project root (go-edu/)
cd geth/21-sync
go test ./exercise/

# Run with verbose output to see test details
go test -v ./exercise/

# Run solution tests (build with solution tag)
go test -tags solution -v ./exercise/

# Run specific test
go test -v ./exercise/ -run TestRun
```

## Code Structure & Patterns

### The Exercise File (`exercise/exercise.go`)

The exercise file contains TODO comments guiding you through:

1. **Input Validation** - Defensive programming (check ctx, client)
2. **SyncProgress Call** - Understanding nil as a sentinel value
3. **Result Interpretation** - Syncing vs synced semantics
4. **Result Construction** - Providing both boolean and detailed views

### The Solution File (`exercise/solution.go`)

The solution file contains detailed STEP comments explaining:
- **Why** we validate inputs (prevent panics)
- **How** nil semantics work (nil = synced, not error)
- **What** the progress fields mean (current/highest blocks)
- **When** to use defensive copying (not needed here, unlike big.Int)

### Key Patterns You'll Learn

#### Pattern 1: Nil as Sentinel Value

```go
// Call returns (nil, nil) when fully synced
progress, err := client.SyncProgress(ctx)
if err != nil {
    return nil, fmt.Errorf("sync progress: %w", err)
}

// Nil progress means synced, non-nil means syncing
isSyncing := progress != nil
```

**Why:** Using nil to signal success is a common Go pattern. Here, "no progress to report" means "sync is complete".

**Building on:** This is different from module 01 where nil was an error. Here it's success.

**Repeats in:** Many Ethereum RPC calls use nil to indicate "use latest" or "fully synced".

#### Pattern 2: Simple Status Reporting

```go
return &Result{
    IsSyncing: isSyncing,      // Boolean for simple checks
    Progress:  progress,        // Detailed data for inspection
}, nil
```

**Why:** Provide both simple (boolean) and detailed (full object) views. Callers can choose the level of detail they need.

**Building on:** Module 01 returned structured Result. Here we add a summary boolean.

**Repeats in:** Many status/health check patterns across all systems.

#### Pattern 3: No Defensive Copying Needed

```go
// Safe to return progress directly (it's a snapshot, not mutable state)
return &Result{
    Progress: progress, // No copy needed
}, nil
```

**Why:** Unlike big.Int or Header (mutable types), SyncProgress is a point-in-time snapshot. No concurrent modification is possible.

**Building on:** Module 01 taught when to copy (big.Int, Header). This module teaches when NOT to copy.

**Repeats in:** Understanding mutability is key to knowing when defensive copying is needed.

## Understanding Sync Modes

### Full Sync (Traditional)
- Replays every transaction from genesis
- Validates every state transition
- Complete history, fully verified
- **Slow:** Can take weeks for mainnet
- **Large:** Requires significant disk space

### Snap Sync (Modern Default)
- Downloads recent state snapshot
- Heals missing data in background
- Still validates via block headers and receipts
- **Fast:** Hours to days for mainnet
- **Efficient:** Less disk I/O during initial sync

### Light Sync (Minimal)
- Downloads block headers only
- Requests proofs on demand
- Minimal storage required
- **Very Fast:** Minutes to hours
- **Limited:** Few nodes serve light clients

## Understanding SyncProgress Fields

```go
type SyncProgress struct {
    StartingBlock uint64 // Block where sync began
    CurrentBlock  uint64 // Current block being processed
    HighestBlock  uint64 // Highest known block in network

    // State sync counters (snap sync mode)
    PulledStates  uint64 // State entries downloaded
    KnownStates   uint64 // Total state entries known
}
```

**Progress calculation:**
```
percentBlocks = (CurrentBlock / HighestBlock) * 100
percentStates = (PulledStates / KnownStates) * 100
```

**Interpreting the values:**
- `CurrentBlock < HighestBlock`: Still syncing blocks
- `PulledStates < KnownStates`: Still syncing state (snap mode)
- `progress == nil`: Fully synced!

## Error Handling

### Common Errors

**1. "method not found"**
```
Cause: Some RPC endpoints don't support eth_syncing
Solution: Use a full node you control, not a public RPC
Prevention: Check RPC capabilities before deploying
```

**2. RPC timeout**
```
Cause: Node is unresponsive or network is down
Solution: Increase timeout, check node health
Prevention: Implement retry logic with backoff
```

**3. "nil client"**
```
Cause: Forgot to initialize ethclient before calling
Solution: Always initialize and validate client
Prevention: Defensive validation at function boundaries
```

## Testing Strategy

The test file demonstrates:

1. **Mock implementations:** `mockSyncClient` implements `SyncClient` interface
2. **Nil semantics testing:** Tests both synced (nil) and syncing (non-nil) cases
3. **Error case testing:** Verifies error handling works correctly
4. **Input validation:** Tests nil client and context handling

**Key insight:** Because we use interfaces, we can test without a real Ethereum node. Tests are fast and deterministic.

## Common Pitfalls & How to Avoid Them

### Pitfall 1: Treating Nil as Error

```go
// BAD: Assumes nil progress is an error
progress, err := client.SyncProgress(ctx)
if progress == nil {
    return errors.New("no progress data") // WRONG!
}

// GOOD: Recognizes nil as "synced"
progress, err := client.SyncProgress(ctx)
if err != nil {
    return err
}
isSyncing := progress != nil // Correct interpretation
```

**Why it's a problem:** Nil progress means success (fully synced), not failure.

**Fix:** Check error first, then interpret nil progress as "synced".

### Pitfall 2: Not Handling Context Properly

```go
// BAD: Passing nil context without validation
Run(nil, client, cfg) // Could panic inside

// GOOD: Validate and provide default
if ctx == nil {
    ctx = context.Background()
}
```

**Why it's a problem:** Some RPC implementations don't handle nil contexts gracefully.

**Fix:** Always validate context at function boundaries.

### Pitfall 3: Comparing Progress Without Nil Check

```go
// BAD: Accessing fields without nil check
if progress.CurrentBlock < progress.HighestBlock {
    // Panics if progress is nil!
}

// GOOD: Check nil first
if progress != nil && progress.CurrentBlock < progress.HighestBlock {
    // Safe
}
```

**Why it's a problem:** Dereferencing nil pointers causes panics.

**Fix:** Always nil-check before field access.

### Pitfall 4: Misunderstanding Sync Completion

```go
// BAD: Assumes CurrentBlock == HighestBlock means synced
if progress.CurrentBlock == progress.HighestBlock {
    // Not necessarily synced! Could be stale node.
}

// GOOD: Check if progress is nil
if progress == nil {
    // Definitely synced
}
```

**Why it's a problem:** A stale node might report equal blocks but be behind the network.

**Fix:** Only trust nil progress as confirmation of sync completion.

## How Concepts Build on Each Other

This module builds on patterns from previous modules while introducing new concepts:

1. **From Module 01-stack:**
   - Context validation → Same pattern here
   - RPC call pattern → Applied to SyncProgress
   - Error wrapping → Consistent usage
   - Result struct pattern → Extended with boolean flag

2. **New in this module:**
   - Nil as sentinel value (not just absence of data)
   - Progress tracking via counters
   - Operational health checks
   - Understanding when NOT to defensive copy

3. **Patterns that repeat throughout the course:**
   - Input validation → Every function
   - Error wrapping → All error returns
   - Interface-based design → All RPC operations
   - Structured results → All complex returns

**The progression:**
- Module 01: Read chain metadata (static identifiers)
- Module 21: Read sync status (dynamic operational state)
- Future modules: Read other operational data (peers, mempool, etc.)

## Fun Facts & Comparisons

### Sync Mode Evolution
- **Pre-2019:** Only full sync available (weeks to sync)
- **2019:** Fast sync introduced (days to sync)
- **2021:** Snap sync introduced (hours to sync)
- **Future:** "Weak subjectivity" checkpoints may enable instant sync

### Light Client Challenges
- Light clients were promised in Ethereum's original design
- Few nodes serve light clients (requires extra resources)
- Most "light" solutions now use trusted RPC providers
- The "portal network" aims to revive true light clients

### State Size Growth
- Genesis (2015): ~0 bytes
- 2018: ~100 GB
- 2020: ~300 GB
- 2023: ~900 GB
- Growth rate: ~200 GB per year

This is why snap sync matters—downloading 900 GB of state is much faster than replaying millions of blocks!

## Comparisons

### Sync Modes Comparison

| Mode | Initial Sync Time | Disk Space | Verification Level | Use Case |
|------|------------------|------------|-------------------|----------|
| Full | Weeks | ~1 TB | Complete | Archive nodes, research |
| Snap | Hours to days | ~800 GB | High (blocks + state) | Production default |
| Light | Minutes to hours | ~1 GB | Moderate (headers + proofs) | Mobile, constrained devices |

### nil vs non-nil in Go

| Return Value | Meaning in Most Functions | Meaning in SyncProgress |
|--------------|--------------------------|------------------------|
| `(nil, nil)` | Usually invalid/error | Fully synced (success!) |
| `(value, nil)` | Success with data | Syncing in progress |
| `(nil, error)` | Failed, check error | RPC call failed |

## Next Steps

After completing this module, you'll move to **22-peers** where you'll:
- Query peer count with `net_peerCount`
- Understand P2P networking health
- Learn about admin APIs for detailed peer info
- Continue building operational monitoring skills
