//go:build reference

package mempool

import (
	"context"
	"errors"
	"fmt"
)

var errNilClient = errors.New("nil mempool client")

/*
Reference Solution - txpool_status / PendingTransactionCount
============================================================

This file demonstrates txpool_status: count of pending transactions in the
mempool. Optional limit caps the returned count for backpressure signals.

This connects to the Ethereum ecosystem by showing:
- PendingTransactionCount(ctx): from txpool_status RPC
- cfg.Limit > 0: cap count for "is mempool overloaded?" checks
- Use case: delay sends when pending count exceeds threshold
*/
func Run(ctx context.Context, client MempoolClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}

	count, err := client.PendingTransactionCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("pending transaction count: %w", err)
	}

	if cfg.Limit > 0 && int(count) > cfg.Limit {
		count = uint(cfg.Limit)
	}

	return &Result{PendingCount: count}, nil
}
