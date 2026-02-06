//go:build reference

package rpcbasics

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"
)

const retryDelay = 100 * time.Millisecond

var errNilClient = errors.New("nil rpc client")

/*
Reference Solution

Structure:
- Read network id and latest block number.
- Fetch the full block with bounded retries.

Invariants:
- At most `max(1, cfg.Retries)` attempts are made.
- The returned block pointer is the pointer returned by the client.
*/
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}

	networkID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("network id: %w", err)
	}
	if networkID == nil {
		return nil, errors.New("network id: nil response")
	}

	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("block number: %w", err)
	}

	attempts := cfg.Retries
	if attempts <= 0 {
		attempts = 1
	}

	var blockErr error
	for i := 0; i < attempts; i++ {
		b, err := client.BlockByNumber(ctx, new(big.Int).SetUint64(blockNumber))
		if err == nil {
			return &Result{
				NetworkID:   new(big.Int).Set(networkID),
				BlockNumber: blockNumber,
				Block:       b,
			}, nil
		}

		blockErr = err
		if i+1 < attempts {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(retryDelay):
			}
		}
	}
	return nil, fmt.Errorf("block by number after %d attempts: %w", attempts, blockErr)
}
