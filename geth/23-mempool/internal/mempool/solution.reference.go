//go:build reference

package mempool

import (
	"context"
	"errors"
	"fmt"
)

var errNilClient = errors.New("nil mempool client")

/*
Reference Solution

Structure:
- Read pending transaction count.
- Optionally clamp output using `cfg.Limit` when provided.
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
