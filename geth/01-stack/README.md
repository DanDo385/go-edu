# 01-stack: Understanding the Ethereum Execution Stack

**Goal:** Understand what Geth is, how it fits with consensus clients, and prove connectivity by reading chain ID + latest block.

## Big Picture: The Ethereum Stack from First Principles

Before diving into code, let's build a mental model from the ground up. Ethereum is fundamentally a **distributed state machine**—think of it like a globally synchronized database where everyone agrees on the same sequence of state transitions. But unlike a traditional database, there's no central authority. Instead, we have a **two-client architecture** that emerged from The Merge (Ethereum's transition to Proof-of-Stake):

### The Two-Client Architecture

**Execution Client (Geth)** = The CPU + Memory + Disk
- **What it does:** Executes EVM bytecode, maintains the state trie (think: Merkle-Patricia tree indexing all account balances and contract storage), and exposes JSON-RPC endpoints
- **Computer Science analogy:** Like a CPU executing instructions, Geth executes transactions. The state trie is like a hash table, but cryptographically verifiable—you can prove "account X has balance Y" without downloading the entire blockchain
- **Fun fact:** Geth stands for "Go Ethereum"—it's written in Go, but there are other execution clients: Erigon (also Go), Nethermind (C#), Besu (Java). They all implement the same EVM spec, so they're interchangeable!

**Consensus Client** (Prysm, Lighthouse, Nimbus, etc.) = The Scheduler + Validator
- **What it does:** Runs the Beacon Chain, manages validators, drives fork choice (decides which chain is canonical), and tells the execution client "execute this block"
- **Computer Science analogy:** Like an operating system scheduler deciding which process runs next, the consensus client decides which block gets appended to the chain
- **Nerdy detail:** The Beacon Chain uses a BFT-style consensus (Casper FFG + LMD GHOST). Validators stake ETH and vote on blocks. If you vote incorrectly, you get slashed (lose ETH). This is why it's called "Proof-of-Stake"

**JSON-RPC** = The API Layer
- **What it does:** Exposes a standardized interface (JSON-RPC 2.0) for querying data and submitting transactions
- **Computer Science analogy:** Like REST APIs for web services, JSON-RPC is the protocol for interacting with Ethereum nodes
- **Protocol detail:** JSON-RPC is stateless and request-response based. Methods like `eth_blockNumber` return data, while `eth_sendTransaction` submits work

### The Complete Picture

```
┌─────────────────────────────────────────────────────────┐
│                    Your Application                      │
│              (Go code using ethclient)                   │
└────────────────────┬────────────────────────────────────┘
                     │ JSON-RPC (HTTP/WebSocket)
                     ▼
┌─────────────────────────────────────────────────────────┐
│              Execution Client (Geth)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   EVM Exec   │  │  State Trie  │  │  JSON-RPC    │  │
│  │   Engine     │  │  (Merkle)    │  │  Server      │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└────────────────────┬────────────────────────────────────┘
                     │ Engine API (local IPC)
                     ▼
┌─────────────────────────────────────────────────────────┐
│           Consensus Client (Prysm/Lighthouse)           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ Beacon Chain │  │   Fork       │  │  Validator   │  │
│  │   Logic      │  │   Choice     │  │  Management  │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└────────────────────┬────────────────────────────────────┘
                     │ P2P Gossip Protocol
                     ▼
              ┌──────────────┐
              │ Other Nodes  │
              │  (Peers)     │
              └──────────────┘
```

## Learning Objectives

By the end of this module, you should be able to:

1. **Draw the high-level Ethereum stack:** execution vs consensus vs networking vs JSON-RPC
2. **Use Go + `ethclient` to dial an RPC endpoint** with proper timeout handling (critical for production!)
3. **Query `chainId`, `net_version`, and the latest block header**—these are your "hello world" operations
4. **Interpret the difference between chain ID and network ID:**
   - **Chain ID** (EIP-155): Used for replay protection in transaction signing. Mainnet = 1, Sepolia = 11155111, etc.
   - **Network ID** (legacy): Older identifier, often matches chain ID but not guaranteed. Some networks use different values.
5. **Understand why "public RPC" ≠ "running a node":**
   - Public RPCs (Infura, Alchemy) are convenient but rate-limited
   - They often disable admin/debug endpoints (`debug_traceTransaction`, `admin_*`)
   - Running your own node gives you full power but requires ~1TB disk space and sync time

## Prerequisites

- **Go basics:** modules, `go run`, flags, contexts
- **Conceptual Ethereum familiarity:** blocks, transactions, state (if you've done Solidity, you know this!)
- **From Solidity-edu:**
  - **01 Datatypes & Storage:** Block headers carry `stateRoot` that indexes the storage trie. Every storage slot you learned about in Solidity is committed to this root!
  - **03 Events & Logging:** Headers include `logsBloom` (bloom filter) and `receiptsRoot` used for efficient event queries

## Real-World Analogies

### The City Records Office Analogy
Calling the city records office: you ask "What's the latest ledger page number?" (block number) and "Which city am I talking to?" (chain ID). The clerk hands you the stamped page header (block hash/parent hash). The header is like a checksum—if someone tampered with the page, the hash would change.

### The CPU Register Snapshot Analogy
Think of a block header as a CPU register snapshot after executing a batch of instructions. The `stateRoot` is like the memory state, `transactionsRoot` is the instruction log, and `receiptsRoot` is the execution trace. The `parentHash` links to the previous snapshot, creating an immutable chain.

### The Git Commit Analogy
A block is like a Git commit:
- **Block hash** = commit SHA
- **Parent hash** = parent commit SHA
- **State root** = tree hash of the entire repository state
- **Transactions** = the diff/changes in that commit

## Fun Facts & Nerdy Details

### Chain ID History
- **EIP-155** (2016) introduced `chainId` to prevent replay attacks. Before this, a transaction signed on mainnet could be replayed on Ethereum Classic (ETC) or testnets!
- **Replay attack scenario:** You sign a transaction to send 1 ETH to Alice on mainnet. An attacker copies your signature and broadcasts it on ETC. Without chain ID, ETC would accept it, and you'd lose 1 ETC too!
- **Current chain IDs:** Mainnet = 1, Sepolia = 11155111, Holesky = 17000, Base = 8453, Arbitrum = 42161

### Network ID vs Chain ID
- **Network ID** (`net_version`) predates chain ID and was used for P2P networking (identifying which network peers belong to)
- **Chain ID** is used for transaction signing and replay protection
- On mainnet, they're both 1, but on some networks they differ (historical reasons)

### Performance Considerations
- **`eth_blockNumber`** is extremely cheap—just returns a number
- **`eth_getBlockByNumber`** with `fullTransactions=true` can return megabytes of data (each transaction includes calldata, which can be large)
- **Headers vs Full Blocks:** Headers are ~500 bytes. Full blocks can be 100KB-2MB depending on transaction count and calldata size
- **Pro tip:** If you only need the block number and hash, use `HeaderByNumber` instead of `BlockByNumber`

### Public RPC Limitations
- **Rate limits:** Free tiers often limit to 100k requests/day
- **Missing endpoints:** `debug_*`, `admin_*`, `trace_*` are usually disabled
- **Caching:** Responses may be cached, so you might not see the absolute latest state
- **Solution:** Run your own node for production applications (Geth, Erigon, or Nethermind)

## Comparisons

### Go `ethclient` vs JavaScript `ethers.js`
- **Same protocol:** Both use JSON-RPC under the hood
- **Ergonomics:** `ethers.js` has more helper methods (e.g., `provider.getNetwork()` returns chain ID + name), while `ethclient` is more low-level
- **Type safety:** Go's static typing catches errors at compile time
- **Performance:** Go is faster for heavy workloads, but JS is fine for most use cases

### Geth vs Other Execution Clients
- **Geth:** Most popular, battle-tested, written in Go
- **Erigon:** More efficient storage (stores state history), faster sync, also Go
- **Nethermind:** C#, good performance, active development
- **Besu:** Java, enterprise-friendly, Hyperledger project
- **API compatibility:** All implement the same JSON-RPC spec, so your code works with any of them!

### Mainnet vs L2 RPCs
- **L2s (Layer 2s)** like Arbitrum, Optimism, Base run their own execution clients
- **Different semantics:** Some expose L2-specific fields (e.g., `l1BlockNumber` on Optimism)
- **Gas fields:** L2s may have different gas pricing (e.g., Arbitrum uses L1 gas price + L2 fee)
- **Same JSON-RPC:** The core methods (`eth_blockNumber`, `eth_getBalance`) work the same way

## Building on Previous Concepts

This is your **first module**, so there are no previous geth modules to reference yet! But we're building on concepts from Solidity-edu:

- **From Solidity 01 (Datatypes & Storage):** You learned about storage slots. The `stateRoot` in block headers is a Merkle root committing to ALL storage slots across ALL contracts. It's like a cryptographic checksum of the entire Ethereum state!

- **From Solidity 03 (Events & Logging):** You learned about events and logs. Block headers include `logsBloom` (a bloom filter) that allows fast "does this block contain Transfer events?" queries without downloading all logs. The `receiptsRoot` commits to all transaction receipts (which contain logs).

## What You'll Build

In this module, you'll implement a library function that:
1. Validates input parameters (defensive programming)
2. Connects to an Ethereum RPC endpoint via the RPCClient interface
3. Queries the chain ID (proves you're connected to the right network)
4. Queries the network ID (legacy identifier)
5. Fetches a block header (proves connectivity and shows current state)
6. Returns a Result struct with defensive copies (immutability pattern)

This is your "hello world" for Ethereum Go development. Every subsequent module builds on these fundamentals!

## Code Structure & Patterns

### The Exercise File (`exercise/exercise.go`)

The exercise file contains TODO comments guiding you through the implementation. Each TODO represents a fundamental concept:

1. **Input Validation** - Learn defensive programming patterns
2. **RPC Calls** - Understand how to interact with Ethereum nodes
3. **Error Handling** - Master Go's idiomatic error wrapping
4. **Defensive Copying** - Learn why immutability matters in concurrent systems

### The Solution File (`exercise/solution.go`)

The solution file contains detailed educational comments explaining:
- **Why** each step is necessary (the reasoning behind the code)
- **How** concepts repeat and build on each other (pattern recognition)
- **What** fundamental principles are being demonstrated (computer science concepts)

### Key Patterns You'll Learn

#### Pattern 1: Input Validation → Early Returns
```go
if client == nil {
    return nil, errors.New("client is nil")
}
```
**Why:** Fail fast, don't continue with invalid state. This pattern appears in every function that accepts external input.

**Building on:** Go's zero values (nil for interfaces) and error handling idioms.

**Repeats in:** Every module where we validate inputs (which is all of them!).

#### Pattern 2: RPC Call → Error Check → Nil Check → Use
```go
chainID, err := client.ChainID(ctx)
if err != nil {
    return nil, fmt.Errorf("chain id: %w", err)
}
if chainID == nil {
    return nil, errors.New("chain id response was nil")
}
```
**Why:** 
- RPC calls can fail (network issues, node down, etc.)
- Even successful calls can return nil (malformed responses)
- We need to validate before using to prevent nil pointer panics

**Building on:** Go's multiple return values (value, error) pattern.

**Repeats in:** Every RPC call throughout the entire course. This is THE pattern for Ethereum Go development.

#### Pattern 3: Error Wrapping with Context
```go
return nil, fmt.Errorf("chain id: %w", err)
```
**Why:** The `%w` verb wraps errors, preserving the error chain. This allows:
- `errors.Is()` to check for specific error types
- `errors.As()` to extract underlying error details
- Better error messages that show the full call chain

**Building on:** Go 1.13's error wrapping improvements.

**Repeats in:** Every error return in production Go code.

#### Pattern 4: Defensive Copying
```go
return &Result{
    ChainID:   new(big.Int).Set(chainID),
    NetworkID: new(big.Int).Set(networkID),
    Header:    types.CopyHeader(header),
}, nil
```
**Why:** 
- The RPCClient might return pointers to internal data structures
- If we return those pointers directly, callers could mutate them
- This could cause data races in concurrent code or affect other callers
- By copying, we ensure each caller gets independent data

**Building on:** 
- Go's pointer semantics (sharing vs copying)
- Concurrent programming safety (immutability prevents races)
- Memory management (who owns what data)

**Repeats in:** Every function that returns data from external libraries or shared state.

**Deep dive:** `big.Int` is mutable. If we assigned `chainID` directly, both `Result.ChainID` and the client's internal data would point to the same `big.Int`. Mutating one would mutate both! This is a subtle but critical bug that can cause mysterious failures in production.

### Understanding the Types

#### `RPCClient` Interface
```go
type RPCClient interface {
    ChainID(ctx context.Context) (*big.Int, error)
    NetworkID(ctx context.Context) (*big.Int, error)
    HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}
```

**Why an interface?** 
- Allows testing with mock implementations (see `exercise_test.go`)
- Decouples our code from specific RPC implementations
- Follows Go's "accept interfaces, return structs" principle

**Building on:** Go's interface system (duck typing).

**Repeats in:** Every module uses interfaces for testability and flexibility.

#### `Config` Struct
```go
type Config struct {
    BlockNumber *big.Int
}
```

**Why separate config?**
- Separation of concerns: configuration vs logic
- Makes functions testable (can request specific blocks)
- Allows callers to control behavior without changing function signature

**Building on:** Functional programming principles (pure functions with explicit inputs).

**Repeats in:** Every function that needs configuration or options.

#### `Result` Struct
```go
type Result struct {
    ChainID   *big.Int
    NetworkID *big.Int
    Header    *types.Header
}
```

**Why a struct?** 
- Groups related data together
- Makes return values self-documenting
- Easier to extend (can add fields without breaking callers)

**Building on:** Go's struct types and composition.

**Repeats in:** Every function that returns multiple related values.

### Context Propagation: The Thread That Connects Everything

`context.Context` appears in every RPC call. Why?

1. **Cancellation:** Callers can cancel long-running operations
2. **Timeouts:** Prevents hanging forever on unresponsive nodes
3. **Request-scoped values:** Can pass metadata through call chains

**Example:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

chainID, err := client.ChainID(ctx)
```

If the RPC call takes longer than 5 seconds, it's automatically cancelled. This is critical for production systems!

**Building on:** Go's context package (introduced in Go 1.7).

**Repeats in:** Every RPC call, every network operation, every long-running task.

## Deep Dive: Why Defensive Copying Matters

Let's explore a concrete example of why defensive copying is critical:

### The Bug (Without Defensive Copying)
```go
// BAD: Returning pointers directly
return &Result{
    ChainID: chainID,  // Same pointer as client's internal data!
}
```

**What happens:**
1. Function returns `Result` with `ChainID` pointing to client's internal `big.Int`
2. Caller modifies `result.ChainID.SetUint64(999)`
3. Client's internal data is also modified!
4. Next caller gets wrong chain ID
5. **Data race** if multiple goroutines call the function

### The Fix (With Defensive Copying)
```go
// GOOD: Creating independent copy
return &Result{
    ChainID: new(big.Int).Set(chainID),  // New big.Int, independent copy
}
```

**What happens:**
1. Function creates a new `big.Int` and copies the value
2. Caller modifies `result.ChainID` → only affects their copy
3. Client's internal data remains unchanged
4. Safe for concurrent use

**This pattern is so important that the tests verify it:**
```go
// Mutating the result should not mutate the mock responses (defensive copy)
res.ChainID.SetUint64(99)
if mock.chainID.Uint64() != 1 {
    t.Fatalf("chain id was not copied")
}
```

## Error Handling: Building Robust Systems

Go's error handling philosophy: "Errors are values." This means:
- Errors are first-class citizens (not exceptions)
- You handle them explicitly (no hidden control flow)
- Error values can carry context and be inspected

### Error Wrapping Chain
```go
chainID, err := client.ChainID(ctx)
if err != nil {
    return nil, fmt.Errorf("chain id: %w", err)
}
```

**The error chain:**
1. `client.ChainID()` returns `err` (original error)
2. We wrap it: `fmt.Errorf("chain id: %w", err)`
3. Caller can inspect: `errors.Is(err, context.DeadlineExceeded)`
4. Caller can unwrap: `errors.Unwrap(err)` gets original error

**Why this matters:** When debugging production issues, you can trace exactly where errors originated and why they occurred.

## Testing Strategy

The test file (`exercise_test.go`) demonstrates several important patterns:

1. **Mock implementations:** `mockRPC` implements `RPCClient` interface
2. **Table-driven tests:** Multiple test cases with different scenarios
3. **Defensive copy verification:** Tests ensure immutability
4. **Error case testing:** Tests verify error handling works correctly

**Key insight:** Because we use interfaces, we can test our logic without needing a real Ethereum node. This makes tests fast, reliable, and deterministic.

## Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments explaining every concept
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns and verifying correctness

## Common Pitfalls & How to Avoid Them

### Pitfall 1: Forgetting Nil Checks
```go
// BAD: Can panic if chainID is nil
result := &Result{ChainID: chainID}
```

**Fix:** Always check for nil after RPC calls, even if error is nil.

### Pitfall 2: Not Copying Mutable Types
```go
// BAD: Shares pointer, mutations affect original
result := &Result{ChainID: chainID}
```

**Fix:** Always copy `big.Int` values: `new(big.Int).Set(chainID)`

### Pitfall 3: Ignoring Context
```go
// BAD: No timeout, can hang forever
chainID, err := client.ChainID(context.Background())
```

**Fix:** Always use context with timeouts in production code.

### Pitfall 4: Not Wrapping Errors
```go
// BAD: Loses context about where error occurred
if err != nil {
    return nil, err
}
```

**Fix:** Wrap errors with context: `fmt.Errorf("chain id: %w", err)`

## How Concepts Build on Each Other

This module introduces patterns that repeat throughout the entire course:

1. **Input validation** → Used in every function
2. **RPC call pattern** → Used in every module (02-rpc-basics, 03-keys-addresses, etc.)
3. **Error wrapping** → Used everywhere
4. **Defensive copying** → Critical for concurrent code (16-concurrency, 22-worker-pool)
5. **Context propagation** → Used in every network operation
6. **Interface-based design** → Enables testing and flexibility

**The pattern:** Learn once, apply everywhere. Each module builds on previous patterns while introducing new concepts.

## Next Steps

After completing this module, you'll move to **02-rpc-basics** where you'll:
- Fetch full blocks (not just headers)
- Understand transaction structures
- Add retry logic for resilience
- Learn about JSON-RPC method names and parameters
- **Build on:** The RPC call pattern you learned here
- **Extend:** Error handling with retries and backoff
- **Apply:** The same defensive copying patterns to transaction data
