# geth-25-toolbox: Swiss Army CLI

**Goal:** build a Swiss Army CLI that combines status, block/tx lookup, and event decoding.

## Big Picture

Capstone module that stitches together previous lessons into one tool with subcommands. Reuses RPC basics, block/tx retrieval, receipts/logs decoding, and event filtering.

**Computer Science principle:** Composition—building complex systems from simple, well-tested components. Command pattern for operation dispatch.

## Learning Objectives

1. **Implement command routing** (dispatch to different handlers)
2. **Compose multiple RPC operations** into unified commands
3. **Reuse patterns from previous modules** (01-stack, 21-24)
4. **Build production-ready CLI tools** with subcommands

## Prerequisites
- **Modules 01-24:** All previous modules (this is the capstone!)
- **Comfort with:** Command-line interfaces, composing operations

## Files

- **Exercise:** `exercise/exercise.go` - TODOs guide implementation
- **Solution:** `exercise/solution.go` - Full implementation with command handlers
- **Types:** `exercise/types.go` - Unified interface combining all modules
- **Tests:** `exercise/exercise_test.go` - Test suite

## How to Run Tests

```bash
cd geth/25-toolbox
go test ./exercise/
go test -v ./exercise/
go test -tags solution -v ./exercise/
```

## Supported Commands

### status
Comprehensive node status combining modules 01, 21, 22:
- Chain ID & Network ID
- Latest block number & hash
- Sync status (syncing or synced)
- Peer count

### block <number>
Block details:
- Number, hash, parent hash
- Timestamp
- Transaction count
- Gas used/limit

### tx <hash>
Transaction details:
- Hash, nonce, value
- Gas, gas price
- To address (if not contract creation)
- Pending status

## Key Concepts

### Command Pattern

```go
switch cfg.Command {
case "status":
    return handleStatus(ctx, client)
case "block":
    return handleBlock(ctx, client, cfg.Args)
case "tx":
    return handleTx(ctx, client, cfg.Args)
default:
    return nil, errors.New("unknown command")
}
```

**Why:** Single entry point, multiple behaviors. Easy to extend with new commands.

### Composition Over Implementation

Instead of re-implementing RPC calls, we reuse client methods:
- `ChainID()`, `NetworkID()` from module 01
- `SyncProgress()` from module 21
- `PeerCount()` from module 22

**Benefits:**
- Less code to write/maintain
- Consistent behavior across commands
- Easy to test (mock the client interface)

## How This Builds on Previous Modules

**Module 01:** ChainID, NetworkID, HeaderByNumber → Used in status command
**Module 21:** SyncProgress → Used in status command
**Module 22:** PeerCount → Used in status command
**Module 24:** Block/time patterns → Used in block command

**New:** Command routing, argument parsing, composing operations

## Fun Facts

- **Similar tools:** `cast` (Foundry), `eth` (go-ethereum), `etherscan-cli`
- **Extension ideas:** Events filtering, contract calls, account balance
- **Real-world use:** Many blockchain teams build internal Swiss Army tools

## Common Pitfalls

### Pitfall 1: Not Validating Arguments
- Commands need different args (block needs number, tx needs hash)
- Solution: Validate args per command

### Pitfall 2: Inconsistent Error Messages
- Each handler returns different error formats
- Solution: Wrap errors consistently with command context

### Pitfall 3: Duplicating Logic
- Don't rewrite RPC calls, reuse client interface
- Solution: Compose from existing operations

## Next Steps

Congratulations! You've completed the geth modules series. You now understand:
- RPC client patterns and defensive programming
- Blockchain data structures (headers, blocks, transactions)
- Operational metrics (sync status, peers, mempool, health)
- Building production-ready CLI tools

**What's next?**
- Build your own tools using these patterns
- Explore advanced topics: contract calls, event filtering, indexing
- Contribute to Ethereum tooling ecosystem!
