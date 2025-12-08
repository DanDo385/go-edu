package exercise

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

// MonitorClient captures the ethclient calls needed for module 24 (node monitoring).
type MonitorClient interface {
	// HeaderByNumber fetches a block header by number (nil for latest)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

// Config allows configuration for node monitoring.
type Config struct {
	// MaxLagSeconds is the maximum acceptable lag before marking node as stale
	MaxLagSeconds int64
	// BlockNumber specifies which block to check (nil for latest)
	BlockNumber *big.Int
}

// Result summarizes the node health status.
type Result struct {
	// Status is the health classification (OK, STALE, etc.)
	Status string
	// BlockNumber is the block number checked
	BlockNumber uint64
	// BlockTimestamp is when the block was produced
	BlockTimestamp time.Time
	// LagSeconds is how far behind the block is from current time
	LagSeconds int64
}
