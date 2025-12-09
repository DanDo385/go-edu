//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"errors"
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
	// TODO: Validate input parameters
	// - Check if ctx is nil and provide a default context if needed
	// - Check if client is nil and return an appropriate error
	// - Why: Defensive programming prevents panics from nil pointer dereferences

	// TODO: Call SyncProgress RPC method
	// - Call client.SyncProgress(ctx) to check sync status
	// - Handle potential errors from the RPC call
	// - Key concept: nil response means the node is fully synced
	// - Non-nil response contains sync progress details (current vs highest block)
	// - Why: This is the primary way to check if a node is ready for production use

	// TODO: Interpret the SyncProgress result
	// - If progress is nil, the node is fully synced (IsSyncing = false)
	// - If progress is non-nil, the node is still syncing (IsSyncing = true)
	// - Important fields in SyncProgress:
	//   * CurrentBlock: The current block number the node has processed
	//   * HighestBlock: The highest known block number in the network
	//   * StartingBlock: The block number where syncing started
	//   * PulledStates/KnownStates: State sync counters (for snap sync mode)
	// - Why: Understanding these fields helps diagnose sync issues and estimate time remaining

	// TODO: Construct and return the Result struct
	// - Set IsSyncing based on whether progress is nil or not
	// - Include the progress object for detailed information
	// - Return the result with nil error on success
	// - Note: We can safely return the progress pointer here because SyncProgress
	//   is a read-only operation and the data represents a point-in-time snapshot

	return nil, errors.New("not implemented")
}
