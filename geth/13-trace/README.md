# 13-trace: Transaction Execution Tracing

**Goal:** Trace transaction execution with `debug_traceTransaction` to see call tree, gas usage, and opcode-level details.

## Big Picture: Execution Instrumentation

Transaction tracing is **execution instrumentation**—observing program execution without modifying it. Ethereum's EVM is deterministic, so replaying a transaction always produces the same trace. This is like having a debugger attached to every smart contract!

**Computer Science principle:** Tracing is a form of **program analysis**. By instrumenting execution, we can:
- Measure performance (gas usage per operation)
- Debug behavior (see exact execution flow)
- Analyze security (identify suspicious call patterns)
- Optimize code (find expensive operations)

### The Tracing Model

```
┌─────────────────────────────────────────────────────────┐
│              Transaction Execution Trace                 │
│                                                          │
│  Transaction → EVM Replay → Trace Data                  │
│                                                          │
│  Trace Contains:                                        │
│  • Call tree (which contracts were called)              │
│  • Gas usage (per operation, cumulative)                │
│  • Storage changes (which slots were modified)          │
│  • Return values (what data was returned)               │
│  • Reverts (where and why execution failed)            │
└─────────────────────────────────────────────────────────┘
```

**Real-world analogy:** Like a **flight recorder (black box)** in an airplane:
- Records every operation during flight (transaction execution)
- Can be replayed after the fact to understand what happened
- Essential for debugging crashes (reverted transactions)
- Shows exactly what the pilot (user) did and how the plane (EVM) responded

## Learning Objectives

By the end of this module, you should be able to:

1. **Call `debug_traceTransaction`** via Go's `TraceTransaction` method
2. **Understand trace structure** (call tree, gas usage, storage changes)
3. **Debug transaction failures** using traces
4. **Analyze gas consumption** at the opcode level
5. **Understand tracer types** (default, callTracer, prestateTracer, etc.)

## Prerequisites

- **Modules 05-09:** You should understand transactions, calls, and events
- **Go basics:** Context, error handling, JSON processing
- **Access to debug API:** Local geth/anvil or debug-enabled RPC endpoint

## Building on Previous Modules

### From Module 05-06 (Transactions)
- You learned to build and send transactions
- Now you're seeing what happened during execution
- Traces show the "inside view" of transaction processing

### From Module 01-02 (RPC Basics)
- You learned to query blocks and transactions
- Traces are another view of transaction data
- Same deterministic replay principle

### Connection to Solidity-edu

**From Solidity 03 (Functions & Gas):**
- Traces show actual gas usage for each operation
- You can see which functions are expensive
- Helps optimize gas costs

**From Solidity 07 (Reentrancy):**
- Traces show the call tree, making reentrancy visible
- You can see if a contract calls back into another contract
- Essential for security audits

## Understanding Transaction Tracing

### What is a Trace?

A **trace** is a detailed record of transaction execution. It contains:

1. **Call tree:** Shows which contracts were called, in what order
2. **Gas usage:** How much gas each operation consumed
3. **Storage changes:** Which storage slots were read/written
4. **Return data:** What values were returned from calls
5. **Reverts:** Where and why execution reverted

### Trace Types

Ethereum nodes support different **tracers**:

- **Default tracer:** Full opcode-level trace (very detailed)
- **callTracer:** Simplified call tree (most common)
- **prestateTracer:** Account state before transaction
- **noopTracer:** No-op tracer (testing)
- **Custom tracers:** JavaScript-based custom analysis

### Example Trace (callTracer format)

```json
{
  "from": "0xuseral...",
  "to": "0xcontract...",
  "type": "CALL",
  "gas": "0x5208",
  "gasUsed": "0x5208",
  "input": "0x...",
  "output": "0x...",
  "calls": [
    {
      "from": "0xcontract...",
      "to": "0xother...",
      "type": "DELEGATECALL",
      "gas": "0x3000",
      "gasUsed": "0x2000"
    }
  ]
}
```

This shows the transaction called a contract, which then made a DELEGATECALL to another contract.

## Real-World Analogies

### The Flight Recorder Analogy
- **Transaction** = Flight
- **Trace** = Black box recording
- **Replay** = Investigating what happened
- **Gas usage** = Fuel consumption at each stage
- **Reverts** = Emergency events that need investigation

### The Debugger Analogy
- **Tracing** = Stepping through code with a debugger
- **Opcodes** = Assembly instructions
- **Call stack** = Function call hierarchy
- **Breakpoints** = Specific operations you want to inspect

## Fun Facts & Nerdy Details

### Tracing is Expensive

**Why?**
- Replays entire transaction execution in EVM
- Records every operation (SLOAD, SSTORE, CALL, etc.)
- Builds large JSON data structures
- Can take several seconds for complex transactions

**Impact:**
- Most public RPC providers disable `debug_traceTransaction`
- You need your own node or a debug-enabled endpoint
- Anvil (Foundry's local node) supports tracing by default

### Deterministic Replay

**Key insight:** The EVM is deterministic. Given the same:
- Transaction data
- Block state
- Block number

Replaying a transaction always produces the same result and trace!

**Why this matters:**
- Traces are reproducible (important for debugging)
- Different nodes produce identical traces
- Time-travel debugging (replay old transactions)

### Historical Tracing Requires Archive Nodes

**Problem:** Tracing a transaction requires the state at that block.

**Node types:**
- **Full node:** Keeps recent state (~128 blocks)
- **Archive node:** Keeps all historical state

**Implication:** To trace old transactions, you need an archive node (or a service that provides historical traces).

## Code Structure & Patterns

### The Exercise File (`exercise/exercise.go`)

The exercise file contains TODO comments guiding you through the implementation. Each TODO represents a fundamental concept:

1. **Input Validation** - Learn defensive programming patterns
2. **Trace Fetching** - Call debug_traceTransaction RPC method
3. **Defensive Copying** - Copy JSON data for safe return
4. **Result Construction** - Build informative responses

### The Solution File (`exercise/solution.go`)

The solution file contains detailed educational comments explaining:
- **Why** each step is necessary (the reasoning behind the code)
- **How** concepts repeat and build on each other (pattern recognition)
- **What** fundamental principles are being demonstrated (computer science concepts)

### Key Patterns You'll Learn

#### Pattern 1: Defensive Copying for Byte Slices
```go
// BAD: Shares underlying array
result.Trace = raw

// GOOD: Creates independent copy
traceCopy := make(json.RawMessage, len(raw))
copy(traceCopy, raw)
result.Trace = traceCopy
```

**Why:** json.RawMessage is []byte (slice). Slices share backing arrays, so mutations affect all references.

**Building on:** Module 01 taught defensive copying for big.Int. Module 12 extended it to string slices. Here we apply it to byte slices.

**Repeats in:** Any function returning []byte, json.RawMessage, or other slice types.

#### Pattern 2: Returning Raw JSON for Flexibility
```go
// We return json.RawMessage, not a parsed struct
type Result struct {
    TxHash common.Hash
    Trace  json.RawMessage  // Raw JSON, not parsed
}
```

**Why:** Trace format varies by tracer type and Geth version. Raw JSON allows callers to parse according to their needs.

**Building on:** Separation of concerns—we fetch data, callers interpret it.

**Repeats in:** Any API dealing with variable-format data.

## Error Handling: Building Robust Systems

### Common Tracing Errors

**1. "method not found"**
```
Cause: Node doesn't support debug_traceTransaction
Solution: Use local node (geth/anvil) or debug-enabled RPC
Prevention: Check node capabilities before attempting traces
```

**2. "transaction not found"**
```
Cause: Transaction hash doesn't exist or is too old
Solution: Verify transaction hash and check node type (archive vs full)
Prevention: Verify transaction exists with eth_getTransactionByHash first
```

**3. "missing trie node"**
```
Cause: Node doesn't have historical state for old transactions
Solution: Use archive node for historical traces
Prevention: Only trace recent transactions on full nodes
```

**4. "timeout" or "execution aborted"**
```
Cause: Transaction trace is very large or complex
Solution: Increase RPC timeout or use simpler tracer (e.g., callTracer)
Prevention: Test with simple transactions first
```

### Error Wrapping Strategy

```go
// Layer 1: RPC error
err := client.TraceTransaction(ctx, txHash)
// Error: "method not found"

// Layer 2: Add context
return fmt.Errorf("trace transaction: %w", err)
// Error: "trace transaction: method not found"

// Layer 3: Caller adds more context
return fmt.Errorf("failed to debug tx %s: %w", txHash.Hex(), err)
// Error: "failed to debug tx 0xabcd...: trace transaction: method not found"
```

## Testing Strategy

The test file (`exercise_test.go`) demonstrates several important patterns:

1. **Mock implementations:** `mockTraceClient` implements `TraceClient` interface
2. **JSON fixture data:** Tests use sample trace JSON for realistic scenarios
3. **Defensive copy verification:** Tests ensure trace data is copied
4. **Error case testing:** Tests verify error handling works correctly

**Key insight:** Because we use interfaces, we can test trace processing without needing a real Ethereum node or debug API access.

## Common Pitfalls & How to Avoid Them

### Pitfall 1: Not Copying JSON Data
```go
// BAD: Shares backing array
result.Trace = raw

// Caller mutates:
result.Trace[0] = 'X'
// This corrupts the original data!

// GOOD: Create independent copy
traceCopy := make(json.RawMessage, len(raw))
copy(traceCopy, raw)
result.Trace = traceCopy
```

**Why it's a problem:** json.RawMessage is []byte. Slices are references. Mutations affect all owners.

**Fix:** Always use `make()` + `copy()` for byte slices.

### Pitfall 2: Assuming Debug API is Available
```go
// BAD: Calling without checking availability
trace, err := client.TraceTransaction(ctx, txHash)
// Error: "method not found" (public RPC doesn't support it)

// GOOD: Document requirements and handle errors gracefully
// Document: "Requires debug API access (local node or debug-enabled RPC)"
trace, err := client.TraceTransaction(ctx, txHash)
if err != nil {
    return fmt.Errorf("trace failed (requires debug API): %w", err)
}
```

**Why it's a problem:** Most public RPC providers disable debug_* methods. Users get confusing errors.

**Fix:** Document requirements clearly and provide helpful error messages.

### Pitfall 3: Not Setting Timeouts for Large Traces
```go
// BAD: No timeout, might hang forever
ctx := context.Background()
trace, _ := client.TraceTransaction(ctx, txHash)

// GOOD: Set reasonable timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
trace, _ := client.TraceTransaction(ctx, txHash)
```

**Why it's a problem:** Complex transactions can take 10+ seconds to trace. Without timeouts, operations hang.

**Fix:** Always use context.WithTimeout for potentially slow operations.

### Pitfall 4: Tracing Old Transactions on Full Nodes
```go
// BAD: Trying to trace old transaction on pruned node
txHash := common.HexToHash("0xold...")  // Transaction from 1000 blocks ago
trace, err := client.TraceTransaction(ctx, txHash)
// Error: "missing trie node"

// GOOD: Check node type or only trace recent transactions
// Use archive node for historical traces
```

**Why it's a problem:** Full nodes only keep recent state (~128 blocks). Archive nodes keep all state.

**Fix:** Use archive nodes for historical traces, or limit tracing to recent transactions.

## How Concepts Build on Each Other

This module builds on patterns from previous modules while introducing new concepts:

1. **From Module 01-stack:**
   - Context validation → Same pattern here
   - RPC call pattern → Extended for tracing
   - Error wrapping → Consistent usage

2. **From Module 11-12 (Storage/Proofs):**
   - We read state data (storage, proofs)
   - Now we read execution data (traces)
   - Different views of the same blockchain

3. **New in this module:**
   - Execution instrumentation (observing without modifying)
   - JSON as interchange format (flexible parsing)
   - Debug API limitations (not always available)
   - Performance considerations (tracing is expensive)

4. **Patterns that repeat throughout the course:**
   - Input validation → Every function
   - Defensive copying → All mutable types
   - Error wrapping → All error returns
   - Interface-based testing → All modules

**The progression:**
- Module 01: Read block headers (metadata)
- Module 11: Read storage (state)
- Module 12: Get proofs (verification)
- Module 13: Get traces (execution)
- Future: Build debuggers, analytics, block explorers

## Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments guiding implementation
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite demonstrating patterns

## How to Run Tests

```bash
# From the project root (go-edu/)
cd geth/13-trace
go test ./exercise/

# Run with verbose output
go test -v ./exercise/

# Run solution tests
go test -tags solution -v ./exercise/

# Run specific test
go test -v ./exercise/ -run TestRun
```

## Next Steps

After completing this module, you'll move to **14-explorer** where you'll:
- Build a block/transaction explorer
- Combine blocks, transactions, and metadata
- Practice data aggregation patterns
- Build user-friendly views of blockchain data
