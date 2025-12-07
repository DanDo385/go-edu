//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

const retryDelay = 100 * time.Millisecond

// Run contains the reference solution for module 02-rpc-basics.
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if cfg.Retries < 0 {
		cfg.Retries = 0
	}

	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("block number: %w", err)
	}

	networkID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("network id: %w", err)
	}

	var block *types.Block
	var lastErr error
	for attempt := 0; attempt <= cfg.Retries; attempt++ {
		block, lastErr = client.BlockByNumber(ctx, nil)
		if lastErr == nil {
			break
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled while fetching block: %w", ctx.Err())
		case <-time.After(retryDelay):
		}
	}
	if lastErr != nil {
		return nil, fmt.Errorf("block fetch failed after %d attempts: %w", cfg.Retries+1, lastErr)
	}

	return &Result{
		NetworkID:   networkID,
		BlockNumber: blockNumber,
		Block:       block,
	}, nil
}
