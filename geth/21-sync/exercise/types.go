package exercise

import (
	"context"

	"github.com/ethereum/go-ethereum"
)

// SyncClient captures the ethclient calls needed for module 21 (sync progress).
type SyncClient interface {
	SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error)
}

// Config allows configuration for sync progress inspection.
type Config struct {
	// No config needed for basic sync check
}

// Result summarizes the sync status of the node.
type Result struct {
	// IsSyncing indicates whether the node is currently syncing (true) or fully synced (false)
	IsSyncing bool
	// Progress contains detailed sync progress information (nil if fully synced)
	Progress *ethereum.SyncProgress
}
