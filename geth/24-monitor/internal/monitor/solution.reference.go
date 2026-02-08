//go:build reference

package monitor

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var errNilClient = errors.New("nil monitor client")

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

func Run(ctx context.Context, client MonitorClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}

	header, err := client.HeaderByNumber(ctx, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("header by number: %w", err)
	}
	if header == nil || header.Number == nil {
		return nil, errors.New("header by number: nil or missing number")
	}

	maxLag := cfg.MaxLagSeconds
	if maxLag <= 0 {
		maxLag = 60
	}

	blockTime := time.Unix(int64(header.Time), 0).UTC()
	lag := int64(time.Since(blockTime).Seconds())
	if lag < 0 {
		lag = 0
	}

	status := "OK"
	if lag > maxLag {
		status = "STALE"
	}

	return &Result{
		Status:         status,
		BlockNumber:    header.Number.Uint64(),
		BlockTimestamp: blockTime,
		LagSeconds:     lag,
	}, nil
}
