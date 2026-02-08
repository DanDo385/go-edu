//go:build reference

package sync

import (
	"context"
	"errors"
	"fmt"
)

var errNilClient = errors.New("nil sync client")

/*
Reference Solution - eth_syncing / SyncProgress
================================================

This file demonstrates eth_syncing: check if node is syncing and optionally
return progress (CurrentBlock, HighestBlock). nil progress = fully synced.

This connects to the Ethereum ecosystem by showing:
- SyncProgress(ctx): nil when synced; non-nil with CurrentBlock, HighestBlock
- copyProgress := *progress: defensive copy so caller can't mutate our Result
- IsSyncing: progress != nil

The exercise builds understanding of:
- Defensive copy: *progress dereferences; copyProgress is a value copy. We
  store &copyProgress so Result owns independent data.
- progress may be reused by RPC client; we don't share the pointer.
*/
func Run(ctx context.Context, client SyncClient, cfg Config) (*Result, error) {
	_ = cfg
	if client == nil {
		return nil, errNilClient
	}

	progress, err := client.SyncProgress(ctx)
	if err != nil {
		return nil, fmt.Errorf("sync progress: %w", err)
	}
	if progress == nil {
		return &Result{IsSyncing: false, Progress: nil}, nil
	}

	copyProgress := *progress
	return &Result{IsSyncing: true, Progress: &copyProgress}, nil
}
