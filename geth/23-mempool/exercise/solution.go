//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
)

/*
Problem: Inspect the mempool (transaction pool) to understand pending transactions.

The mempool contains transactions that have been broadcast to the network but not
yet included in a block. Monitoring the mempool helps you understand network congestion,
estimate gas prices, and track your own transactions.

However, mempool visibility is limited for privacy and security reasons. Many public
RPC endpoints don't expose pending transaction details.

Computer science principles highlighted:
  - Queue management (FIFO with priority)
  - Privacy/security trade-offs (transparency vs exploitation)
  - Resource management (mempool size limits)
*/
func Run(ctx context.Context, client MempoolClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Familiar Pattern
	// ============================================================================
	// By now, this validation pattern should be second nature. We validate
	// context and client for defensive programming, preventing panics and
	// providing safe defaults.
	//
	// Building on modules 01, 21, 22: Same pattern, consistently applied.
	if ctx == nil {
		ctx = context.Background()
	}

	if client == nil {
		return nil, errors.New("client is nil")
	}

	// ============================================================================
	// STEP 2: Query Mempool Size - Understanding Transaction Pools
	// ============================================================================
	// The mempool (also called txpool) is where transactions wait before being
	// included in blocks. Understanding the mempool is critical for:
	//   1. Gas price estimation (high mempool = higher prices)
	//   2. Transaction tracking (is my tx still pending?)
	//   3. Network health monitoring (mempool backlog = congestion)
	//
	// How mempools work:
	//   - Transactions are broadcast to nodes via P2P gossip
	//   - Each node validates and adds tx to its local mempool
	//   - Miners/validators select txs from mempool for inclusion
	//   - Higher gas price = higher priority (in most cases)
	//   - Mempool has size limits (default ~4GB in Geth)
	//   - Old/low-fee txs can be evicted if mempool is full
	//
	// Why mempool visibility matters:
	//   - MEV (Maximal Extractable Value): Bots scan mempool for profitable txs
	//   - Front-running: See pending trades, submit higher-fee tx to execute first
	//   - Sandwich attacks: Place txs before and after a target tx
	//   - Gas price estimation: More pending txs = higher competition = higher prices
	//
	// Why mempools are often hidden:
	//   - Privacy: Revealing pending txs enables front-running/MEV exploitation
	//   - Security: Full mempool dumps could enable DoS attacks
	//   - Resources: Returning thousands of pending txs is expensive
	//   - Business: Some MEV searchers pay for private mempool access
	//
	// Available RPC methods (by access level):
	//   - Public (usually available):
	//     * eth_getTransactionByHash: Check if specific tx is pending
	//     * Sometimes: basic mempool stats
	//   - Restricted (rarely available on public RPCs):
	//     * eth_pendingTransactions: Full list of pending txs
	//     * txpool_content: Detailed mempool contents
	//     * txpool_status: Mempool statistics
	//     * txpool_inspect: Transaction summaries
	//
	// Error handling: If the RPC endpoint doesn't support mempool queries,
	// this call will fail. We wrap the error with context for debugging.
	//
	// Building on module 01: Same error wrapping pattern (%w verb).
	pendingCount, err := client.PendingTransactionCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("pending transaction count: %w", err)
	}

	// ============================================================================
	// STEP 3: Interpret Mempool Size - Understanding Congestion
	// ============================================================================
	// The pending transaction count tells us about network activity and congestion.
	//
	// Typical mempool sizes:
	//   - Low activity (< 1000 pending): Base fee decreases, txs confirm quickly
	//   - Moderate (1000-10000): Normal operations, predictable confirmation times
	//   - High (10000-50000): Network congestion, higher gas prices needed
	//   - Extreme (> 50000): Severe congestion, only high-fee txs confirm quickly
	//
	// Transaction lifecycle in mempool:
	//   1. Broadcast: Tx sent to network via P2P
	//   2. Validation: Nodes check signature, nonce, gas limits
	//   3. Propagation: Valid txs gossiped to other nodes
	//   4. Inclusion: Miner/validator includes tx in block
	//   5. Confirmation: Block is finalized on chain
	//
	// Transaction replacement (RBF - Replace By Fee):
	//   - Same sender, same nonce, higher gas price
	//   - Typically requires 10% fee increase minimum
	//   - Useful for "stuck" transactions (fee too low)
	//   - Can also cancel by sending 0 ETH to self with higher fee
	//
	// Mempool eviction rules:
	//   - When mempool is full, lowest-fee txs are evicted
	//   - Very old txs may be dropped (default timeout varies)
	//   - Invalid txs (bad nonce, insufficient balance) are rejected
	//
	// Why mempool sizes differ across nodes:
	//   - Different mempool size limits
	//   - Different P2P connectivity (not all txs reach all nodes)
	//   - Different acceptance policies (some nodes filter by min gas price)
	//   - Some nodes run private/restricted mempools
	//
	// This contrasts with modules 01, 21, 22:
	//   - Module 01: Static chain metadata (never changes)
	//   - Module 21: Sync status (changes during initial sync)
	//   - Module 22: Peer count (changes as connections fluctuate)
	//   - Module 23: Mempool size (changes every second as txs arrive/confirm)
	//
	// No validation needed: uint is always valid (even if zero).

	// ============================================================================
	// STEP 4: Return Result - Simple Count with Future Extensibility
	// ============================================================================
	// We return the pending count in a Result struct. While we only have one
	// field now, the struct design allows future expansion without breaking the API.
	//
	// Future extensions could include:
	//   - Queued vs pending breakdown
	//   - Gas price statistics (min/median/max)
	//   - Sample of recent transactions
	//   - Mempool fee histogram
	//
	// No defensive copying needed: uint is a primitive type (copied by value).
	//
	// Production usage patterns:
	//   1. Gas price estimation: High pending count → increase gas price
	//   2. Transaction monitoring: Check if your tx is still pending
	//   3. Network health: Sustained high mempool → network congestion
	//   4. MEV opportunities: Scan mempool for arbitrage (controversial)
	//
	// Building on previous modules:
	//   - Module 01: Complex types (ChainID, NetworkID, Header)
	//   - Module 21: Boolean + optional struct (IsSyncing, Progress)
	//   - Module 22: Simple metric (PeerCount)
	//   - Module 23: Simple metric with context (PendingCount)
	//   - Pattern: All use Result structs for consistent API design
	//
	// This demonstrates how different types of blockchain data have different
	// access patterns and visibility constraints. Not all data is equally
	// available or trustworthy.
	return &Result{
		PendingCount: pendingCount, // Primitive type, automatically copied by value
	}, nil
}
