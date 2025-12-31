//go:build !solution
// +build !solution

package sync

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
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
}

