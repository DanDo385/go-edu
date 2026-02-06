//go:build reference

package sync

import (
	"context"
	"errors"
	"fmt"
)

var errNilClient = errors.New("nil sync client")

/*
Reference Solution

Structure:
- Query `SyncProgress` once.
- `nil` means fully synced.
- Non-nil means currently syncing.

Pointer notes:
- We copy the returned struct so result ownership is explicit.
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
