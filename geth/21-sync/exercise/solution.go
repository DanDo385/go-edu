//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
)

/*
Problem: Inspect sync progress to determine if your Ethereum node is fully synced.

When running an Ethereum node, the first critical check is whether it's finished
syncing the blockchain. A non-synced node returns stale data and shouldn't be used
for production queries. The SyncProgress RPC call returns nil when fully synced,
or a progress object with current/highest block numbers when syncing.

Computer science principles highlighted:
  - Nil as a sentinel value (nil = fully synced, non-nil = syncing)
  - Progress tracking via counters (current vs highest block)
  - State inspection without mutation (read-only health check)
*/
func Run(ctx context.Context, client SyncClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Defensive Programming Pattern
	// ============================================================================
	// As with all production functions that accept external inputs, we validate
	// before proceeding. This pattern repeats from module 01-stack and will
	// continue throughout all modules.
	//
	// Context validation: If ctx is nil, we provide context.Background() as a
	// safe default. This ensures RPC calls won't panic and gives callers a way
	// to control cancellation/timeouts.
	//
	// Building on module 01: Same validation pattern, applied to sync operations.
	if ctx == nil {
		ctx = context.Background()
	}

	// Client validation: The SyncClient interface is our dependency. If it's nil,
	// we cannot make RPC calls. Checking for nil prevents runtime panics.
	//
	// Why this matters: In Go, interfaces have a nil zero value. Methods called
	// on nil interfaces cause panics. Always validate interface parameters.
	//
	// This pattern repeats: Every function that accepts an interface should
	// validate it's not nil before calling methods on it.
	if client == nil {
		return nil, errors.New("client is nil")
	}

	// ============================================================================
	// STEP 2: Check Sync Progress - Understanding Nil as Sentinel Value
	// ============================================================================
	// The SyncProgress RPC call is unique because it uses nil as a meaningful
	// signal: nil = fully synced, non-nil = syncing in progress.
	//
	// Why this matters: Many RPC calls return (result, error), and you check
	// error first. Here, both nil result and nil error are valid and indicate
	// "fully synced". This is a sentinel value pattern—nil carries meaning.
	//
	// How sync modes work:
	//   - Full sync: Replays every transaction from genesis (slow, complete)
	//   - Snap sync: Downloads state snapshots then heals missing data (fast)
	//   - Light sync: Fetches block headers and requests proofs on demand (minimal)
	//
	// SyncProgress tracks:
	//   - StartingBlock: Where the sync began
	//   - CurrentBlock: Where we are now
	//   - HighestBlock: The network's latest known block
	//   - PulledStates/KnownStates: State trie sync progress (snap sync)
	//
	// Error handling: If the RPC call fails (network issue, node down), we wrap
	// the error with context. This helps debugging by showing where the failure
	// occurred in the call chain.
	//
	// Building on module 01: We use the same error wrapping pattern (%w verb)
	// introduced in 01-stack. This preserves error chains for errors.Is/As.
	progress, err := client.SyncProgress(ctx)
	if err != nil {
		return nil, fmt.Errorf("sync progress: %w", err)
	}

	// ============================================================================
	// STEP 3: Interpret Result - Nil vs Non-Nil Semantics
	// ============================================================================
	// Now we interpret the progress result. This is simpler than other modules
	// because we're just checking nil vs non-nil, but the implications are
	// significant for production systems.
	//
	// If progress is nil:
	//   - The node has completed initial sync
	//   - It's receiving new blocks as they're produced (~12s on mainnet)
	//   - Data queries will return current, accurate information
	//   - The node is ready for production use
	//
	// If progress is non-nil:
	//   - The node is still syncing historical blocks/state
	//   - Data queries may return stale or incomplete results
	//   - You should wait before using the node for critical operations
	//   - Progress tracking: (CurrentBlock / HighestBlock) shows percentage
	//
	// No defensive copying needed: Unlike big.Int or Header objects that are
	// mutable, SyncProgress is a snapshot. Even if it contains pointers internally,
	// the RPC layer creates fresh instances for each call. We can safely return it.
	//
	// Why this is safe: SyncProgress represents a point-in-time snapshot. The
	// values don't change after the call completes. No concurrent goroutines
	// are mutating this data.
	//
	// How concepts build:
	//   - Module 01: We copied big.Int and Header (mutable types)
	//   - This module: We don't copy SyncProgress (immutable snapshot)
	//   - Lesson: Defensive copying depends on mutability, not just pointers
	isSyncing := progress != nil

	// ============================================================================
	// STEP 4: Return Result - Simple Status Reporting
	// ============================================================================
	// We package the sync status into our Result struct. This provides a clean
	// API: callers get both a boolean (IsSyncing) for simple checks and the full
	// progress object for detailed inspection.
	//
	// API design principle: Provide both simple and detailed views. IsSyncing
	// answers "can I use this node?" while Progress answers "how long until ready?"
	//
	// Production usage patterns:
	//   1. Health checks: Poll this endpoint, alert if IsSyncing is true too long
	//   2. Startup coordination: Wait for IsSyncing = false before serving traffic
	//   3. Monitoring: Track CurrentBlock/HighestBlock gap over time
	//
	// Building on previous modules:
	//   - Module 01: Read chain metadata (ChainID, NetworkID, Header)
	//   - This module: Read sync status (operational health)
	//   - Pattern: Both return structured Result types with multiple fields
	//
	// What's different:
	//   - Module 01: Always returns data (chain never changes identity)
	//   - This module: Result varies (syncing vs synced are different states)
	//
	// This demonstrates state inspection patterns that you'll use throughout
	// blockchain development: checking if systems are ready before using them.
	return &Result{
		IsSyncing: isSyncing,
		Progress:  progress, // Safe to return: snapshot from RPC, not shared state
	}, nil
}
