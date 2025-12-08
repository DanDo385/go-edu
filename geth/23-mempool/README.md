# geth-23-mempool: Mempool Inspection

**Goal:** inspect pending transactions (where supported) and understand mempool visibility limits.

## Big Picture

Pending txs live in the txpool before inclusion. Many public RPCs do not expose full mempool; Geth's `eth_pendingTransactions` or `txpool_*` may be restricted. Visibility varies by provider and node config.

**Computer Science principle:** Queue management with priority—mempools are priority queues where higher-fee transactions get processed first. Privacy/security trade-offs limit visibility to prevent MEV exploitation.

## Learning Objectives

By the end of this module, you should be able to:

1. **Query mempool size** and understand congestion levels
2. **Understand visibility limits** of public vs private mempools
3. **Recognize MEV implications** of mempool transparency
4. **Learn transaction replacement** rules (Replace-By-Fee)

## Prerequisites
- **Modules 05-06:** Transaction basics, nonces, fees
- **Module 01:** RPC client basics
- **Comfort with:** Transaction lifecycle, gas pricing

## Files

- **Exercise:** `exercise/exercise.go` - TODOs guide implementation
- **Solution:** `exercise/solution.go` - Full implementation with educational comments
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite

## How to Run Tests

```bash
cd geth/23-mempool
go test ./exercise/
go test -v ./exercise/
go test -tags solution -v ./exercise/
```

## Key Concepts

### Mempool Basics
- **What it is:** Queue of unconfirmed transactions waiting for block inclusion
- **How it works:** Nodes receive txs via gossip, validate, store in mempool
- **Size limits:** Default ~4GB in Geth, old/low-fee txs evicted when full
- **Priority:** Higher gas price = higher priority for inclusion

### Why Mempool Visibility is Limited

**MEV (Maximal Extractable Value):**
- Bots scan mempools for profitable opportunities
- Front-running: Submit higher-fee tx to execute before target
- Sandwich attacks: Surround target tx with buy/sell orders
- Privacy concerns: Revealing pending trades enables exploitation

**Access Levels:**
- **Public RPC:** Usually only supports checking specific tx by hash
- **Private/Admin API:** Full mempool access (eth_pendingTransactions, txpool_content)
- **Paid Services:** Some MEV searchers pay for private mempool access

### Transaction Replacement (Replace-By-Fee)

```go
// Replace a pending tx with same nonce but higher fees
// Typically requires 10% fee increase minimum
```

**Use cases:**
- Speed up stuck transaction (fee too low)
- Cancel transaction (send 0 ETH to self with higher fee)

## Common Pitfalls

### Pitfall 1: Expecting Public RPC Mempool Access
- Most public RPCs hide mempool for security/privacy
- Solution: Run your own node for full mempool access

### Pitfall 2: Ignoring Mempool Differences
- Each node has slightly different mempool contents
- Different size limits, connectivity, acceptance policies

### Pitfall 3: Forgetting Transaction Lifecycle
- Transactions can be evicted from mempool if fees too low
- No confirmation guarantee until included in block

## How Concepts Build

1. **From Module 05-06:** Transaction construction → Now understanding where txs wait
2. **From Module 22:** P2P gossip → Applies to transaction propagation
3. **New:** Mempool as priority queue, MEV/privacy concerns

## Fun Facts

- **MEV extracted:** $600M+ in 2023 from front-running, arbitrage, liquidations
- **Flashbots:** Private transaction relay to reduce harmful MEV
- **PBS (Proposer-Builder Separation):** Future Ethereum upgrade to democratize MEV

## Next Steps

After completing this module, you'll move to **24-monitor** where you'll:
- Implement node health checks (block freshness, lag detection)
- Learn alerting patterns for production systems
