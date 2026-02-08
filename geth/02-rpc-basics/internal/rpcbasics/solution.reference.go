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
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
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
