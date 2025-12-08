# geth-22-peers: Peer Count and P2P Health

**Goal:** query peer count and understand p2p gossip health.

## Big Picture

Peers gossip txs/blocks across the network. Peer count is a coarse signal of connectivity; richer info lives in admin APIs. Public RPCs often hide peer details.

**Computer Science principle:** Peer-to-peer mesh networks rely on redundant connections for fault tolerance and data propagation. Peer count is a basic metric for assessing network health.

## Learning Objectives

By the end of this module, you should be able to:

1. **Call net_peerCount** and understand the returned value
2. **Interpret peer count** as a connectivity health metric
3. **Recognize limitations** of public RPC vs your own node
4. **Understand P2P networking** fundamentals (discovery, gossip, protocols)

## Prerequisites
- **Module 01 (01-stack):** RPC client basics, context handling
- **Module 21 (21-sync):** Operational health checks
- **Comfort with:** P2P networking concepts (helpful but not required)

## Building on Previous Modules

### From Module 01 (01-stack)
- You learned RPC call patterns (call → error check → return)
- Same pattern applies here for PeerCount
- Input validation pattern repeats (context, client checks)

### From Module 21 (21-sync)
- You learned health checks (sync status)
- Now learning another health metric (peer connectivity)
- Both are operational concerns, not blockchain data

### New in this module
- **P2P networking concepts:** Discovery, gossip protocols
- **Primitive type handling:** uint64 (no defensive copying needed)
- **Metric interpretation:** What counts as "healthy" connectivity

### Connection to Solidity-edu
- None directly; this supports all chain interactions by ensuring healthy connectivity

## Real-World Analogy

### The Radio Station Analogy
- **Peer count:** Number of radio stations your node can hear
- **More stations:** Better news propagation (redundancy)
- **Zero stations:** Complete isolation (can't hear any broadcasts)
- **Quality matters:** Some stations have static, others are clear

### The Social Network Analogy
- **Peers:** Your friends on a social network
- **Gossip:** When you share a post, friends share it with their friends
- **0 friends:** Nobody sees your posts (isolated)
- **Many friends:** Posts spread quickly (well-connected)

### The Mesh Network Analogy
- **P2P network:** Interconnected web of nodes, no central hub
- **Each connection:** Two-way street for data exchange
- **Redundancy:** Multiple paths for information to flow
- **Resilience:** Network survives even if some nodes fail

## What You'll Build

In this module, you'll create a function that:
1. Takes an RPC client as input
2. Calls the PeerCount method
3. Returns the number of connected peers
4. Provides a basis for network health monitoring

**Key learning:** You'll understand how to assess if your node is well-connected to the Ethereum network!

## Files

- **Exercise:** `exercise/exercise.go` - Your starting point with TODO comments
- **Solution:** `exercise/solution.go` - Complete implementation with detailed educational comments
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite verifying correctness

## How to Run Tests

To run the tests for this module:

```bash
# From the project root (go-edu/)
cd geth/22-peers
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
2. **PeerCount Call** - Query P2P layer for connection count
3. **Result Interpretation** - Understanding what the count means
4. **Result Construction** - Returning simple metric data

### The Solution File (`exercise/solution.go`)

The solution file contains detailed STEP comments explaining:
- **Why** we validate inputs (prevent panics)
- **How** P2P networking works (discovery, gossip, protocols)
- **What** the peer count tells us (and doesn't tell us)
- **When** to use defensive copying (not needed for primitives)

### Key Patterns You'll Learn

#### Pattern 1: Simple Metric Retrieval

```go
// Single RPC call, single return value
peerCount, err := client.PeerCount(ctx)
if err != nil {
    return nil, fmt.Errorf("peer count: %w", err)
}
```

**Why:** Some metrics are simple scalars (uint64), not complex objects.

**Building on:** Module 01 retrieved complex objects. This module retrieves a simple count.

**Repeats in:** Many monitoring/observability operations.

#### Pattern 2: No Defensive Copying for Primitives

```go
// uint64 is copied by value automatically
return &Result{
    PeerCount: peerCount, // No need for defensive copying
}, nil
```

**Why:** Primitive types (uint64, int, bool) are immutable and copied by value in Go. No shared state exists.

**Building on:** Module 01 taught defensive copying for pointers. This module teaches when it's NOT needed.

**Repeats in:** All functions that return primitive types.

#### Pattern 3: Health Metric Interpretation

```go
// Understanding what the number means
if peerCount == 0 {
    // Node is isolated
} else if peerCount < 10 {
    // Low connectivity
} else {
    // Healthy connectivity
}
```

**Why:** Raw numbers need context to be actionable. Thresholds define "healthy" vs "unhealthy".

**Building on:** Module 21 had binary health (syncing vs synced). This module has continuous metric.

**Repeats in:** All monitoring and alerting systems.

## Understanding P2P Networking in Ethereum

### How Peer Discovery Works

1. **Bootstrap nodes:** Initial connection points (hard-coded in client)
2. **DHT (Kademlia):** Distributed hash table for finding peers
3. **DNS discovery:** Query DNS for peer lists (EIP-1459)
4. **Peer exchange:** Peers share lists of other peers they know

### How Gossip Protocols Work

- **New transaction:** Broadcast to all connected peers
- **Peers validate:** Check signature, nonce, gas limits
- **Peers propagate:** Forward to their peers (if valid)
- **Redundancy:** Multiple nodes receive the same tx (ensures reliability)

### Protocol Negotiation

When two nodes connect, they negotiate which protocols to speak:
- **eth/67:** Current Ethereum protocol version
- **snap/1:** Snap sync protocol
- **wit/0:** Witness protocol (for stateless clients)

### Peer Connection Lifecycle

1. **Discovery:** Find peer via DHT/DNS
2. **TCP connection:** Establish network connection
3. **RLPx handshake:** Encrypted communication setup
4. **Protocol negotiation:** Agree on capabilities
5. **Active peering:** Exchange blocks/transactions
6. **Disconnect:** Timeout, error, or intentional drop

## Understanding Peer Count Values

### What Different Counts Mean

- **0 peers:** Node is completely isolated
  - Possible causes: Firewall blocking P2P port (30303), network down, no discovery
  - Impact: Cannot sync, cannot broadcast transactions

- **1-5 peers:** Very low connectivity
  - Possible causes: Strict firewall, network restrictions, new node
  - Impact: Slow block propagation, vulnerable to network partitions

- **10-25 peers:** Moderate connectivity
  - Typical for nodes behind NAT or residential networks
  - Impact: Generally acceptable for personal use

- **25-50 peers:** Good connectivity
  - Typical default max peers setting in Geth
  - Impact: Reliable block/tx propagation

- **50-100+ peers:** Excellent connectivity
  - Requires explicit configuration (--maxpeers flag)
  - Impact: Best for public RPC endpoints or validators

### Peer Quality vs Quantity

**Quantity (peer count):**
- How many connections exist
- Simple metric to track
- Available via public RPC

**Quality (peer details):**
- Client version (geth, erigon, nethermind)
- Network latency (ping time)
- Bandwidth capacity
- Geographic distribution
- Only available via admin_peers API

## Error Handling

### Common Errors

**1. "method not found"**
```
Cause: RPC endpoint doesn't support net_peerCount
Solution: Use a different RPC provider or run your own node
Prevention: Check RPC capabilities documentation
```

**2. "connection refused"**
```
Cause: Node is down or unreachable
Solution: Check node status, network connectivity
Prevention: Implement retry logic with exponential backoff
```

**3. Always returns 0**
```
Cause: Public RPC provider hides peer count (privacy)
Solution: Use your own node for accurate peer metrics
Prevention: Understand provider limitations upfront
```

## Testing Strategy

The test file demonstrates:

1. **Mock implementations:** `mockPeerClient` implements `PeerClient` interface
2. **Different peer counts:** Tests with 0, 10, 50 peers
3. **Error case testing:** Verifies error handling works correctly
4. **Input validation:** Tests nil client and context handling

**Key insight:** Because we use interfaces, we can test without a real Ethereum node. Tests are fast and deterministic.

## Common Pitfalls & How to Avoid Them

### Pitfall 1: Assuming Peer Count Equals Health

```go
// BAD: Assumes any peers means healthy
if peerCount > 0 {
    return "healthy" // Too simplistic!
}

// GOOD: Use thresholds for nuanced assessment
if peerCount == 0 {
    return "critical: no peers"
} else if peerCount < 10 {
    return "warning: low connectivity"
} else {
    return "healthy"
}
```

**Why it's a problem:** A single peer is barely better than zero. Need multiple peers for redundancy.

**Fix:** Use tiered thresholds (critical, warning, healthy).

### Pitfall 2: Ignoring Peer Quality

```go
// BAD: Thinks 50 peers means 50 good connections
// Some peers might be slow, malicious, or on wrong network

// GOOD: For quality assessment, use admin_peers
// (Only available on nodes you control)
```

**Why it's a problem:** Peer count doesn't reveal quality. Could have 50 slow/malicious peers.

**Fix:** For production monitoring, also track block freshness, sync status.

### Pitfall 3: Relying on Public RPC Peer Counts

```go
// BAD: Trusting peer count from public RPC
// Many public RPCs return 0 or fake values for privacy

// GOOD: Use your own node for accurate peer metrics
```

**Why it's a problem:** Public RPCs often hide peer details for security/privacy.

**Fix:** Run your own node if you need accurate peer monitoring.

### Pitfall 4: Not Monitoring Over Time

```go
// BAD: Checking peer count once
peerCount := getPeerCount()
if peerCount > 25 { return "ok" }

// GOOD: Monitor trends over time
// Sudden drops indicate issues (network, node crash, attacks)
```

**Why it's a problem:** Single point-in-time checks miss transient issues.

**Fix:** Track peer count over time, alert on sudden changes.

## How Concepts Build on Each Other

This module builds on patterns from previous modules while introducing new concepts:

1. **From Module 01-stack:**
   - Context validation → Same pattern here
   - RPC call pattern → Applied to PeerCount
   - Error wrapping → Consistent usage
   - Result struct pattern → Simplified (one field)

2. **From Module 21-sync:**
   - Operational health checks → Extended to P2P layer
   - Understanding node readiness → Connectivity is part of readiness
   - Metric interpretation → Different metric, same concept

3. **New in this module:**
   - P2P networking concepts (discovery, gossip)
   - Primitive type handling (no defensive copying)
   - Metric thresholds and interpretation
   - Understanding API limitations (public vs admin)

4. **Patterns that repeat throughout the course:**
   - Input validation → Every function
   - Error wrapping → All error returns
   - Interface-based design → All RPC operations
   - Structured results → All returns (even simple ones)

**The progression:**
- Module 01: Read chain metadata (ChainID, NetworkID, Header)
- Module 21: Read sync status (operational health)
- Module 22: Read peer count (network health)
- Future modules: More operational metrics (mempool, monitoring)

## Fun Facts & Comparisons

### P2P Protocol Evolution

- **devp2p v4 (2014):** Original Ethereum P2P protocol
- **devp2p v5 (2019):** Improved discovery with ENR (Ethereum Node Records)
- **eth/67 (2021):** Current Ethereum wire protocol
- **Portal Network (future):** Ultra-light clients via gossip

### Peer Limits

- **Geth default:** 50 max peers
- **Erigon default:** 100 max peers
- **Archive nodes:** Often 200+ peers
- **Validators:** Sometimes limit to trusted peers only

### Network Topology Facts

- **Small world property:** Most nodes are a few hops away from each other
- **Scale-free network:** Some nodes (super-peers) have many connections
- **Resilience:** Network survives even if 50% of nodes fail

### Comparison: Ethereum vs Bitcoin P2P

| Aspect | Ethereum | Bitcoin |
|--------|----------|---------|
| Default max peers | 50 | 125 |
| Discovery mechanism | Kademlia DHT | Seed nodes + gossip |
| Protocol versions | eth/67, snap/1 | P2P version 70016 |
| Block propagation | ~1-2 seconds | ~10-30 seconds |

## Comparisons

### Public RPC vs Own Node

| Metric | Public RPC | Own Node |
|--------|-----------|----------|
| Peer count accuracy | Often hidden/fake | Accurate |
| admin_peers access | No | Yes (if enabled) |
| Cost | Free (rate limited) | Infrastructure costs |
| Trust | Must trust provider | Trustless |

### net_peerCount vs admin_peers

| Information | net_peerCount | admin_peers |
|-------------|--------------|-------------|
| Total count | Yes | Yes (array length) |
| Client versions | No | Yes |
| Latency | No | Yes |
| Capabilities | No | Yes |
| Access required | Standard RPC | Admin API |

### Primitive vs Pointer Types

| Type | Copied By | Defensive Copy Needed? | Example |
|------|-----------|------------------------|---------|
| uint64 | Value | No | PeerCount |
| *big.Int | Pointer | Yes | ChainID |
| *types.Header | Pointer | Yes | Header |
| bool | Value | No | IsSyncing |

## Next Steps

After completing this module, you'll move to **23-mempool** where you'll:
- Inspect pending transactions (where supported)
- Understand mempool visibility limits
- Learn about transaction replacement rules
- Continue building operational monitoring skills
