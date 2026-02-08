//go:build reference

package peers

import (
	"context"
	"errors"
	"fmt"
)

var errNilClient = errors.New("nil peer client")

/*
Reference Solution - net_peerCount
==================================

This file demonstrates net_peerCount: number of connected P2P peers. Simple
RPC wrapper for node health and connectivity monitoring.

This connects to the Ethereum ecosystem by showing:
- PeerCount(ctx): hex-encoded count from net_peerCount RPC
- Use case: health checks, "is node connected to network?"
- cfg unused: interface placeholder for future (e.g. timeout)
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
