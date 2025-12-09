package exercise

import (
	"context"
)

// PeerClient captures the ethclient calls needed for module 22 (peer counting).
type PeerClient interface {
	// PeerCount returns the number of connected peers
	PeerCount(ctx context.Context) (uint64, error)
}

// Config allows configuration for peer count query.
type Config struct {
	// No config needed for basic peer count
}

// Result summarizes the peer connectivity of the node.
type Result struct {
	// PeerCount is the number of currently connected peers
	PeerCount uint64
}
