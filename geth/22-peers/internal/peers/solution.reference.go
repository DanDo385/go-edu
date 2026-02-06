//go:build reference

package peers

import (
	"context"
	"errors"
	"fmt"
)

var errNilClient = errors.New("nil peer client")

/*
Reference Solution

Structure:
- Query peer count once.
- Return result in stable struct shape.
*/
func Run(ctx context.Context, client PeerClient, cfg Config) (*Result, error) {
	_ = cfg
	if client == nil {
		return nil, errNilClient
	}

	count, err := client.PeerCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("peer count: %w", err)
	}

	return &Result{PeerCount: count}, nil
}
