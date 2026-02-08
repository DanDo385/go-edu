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
Reference Solution - Block Lag Monitor
======================================

This file demonstrates a block lag monitor: compare latest block timestamp to
now. Lag > MaxLagSeconds = STALE. Used for alerting when node falls behind.

This connects to the Ethereum ecosystem by showing:
- HeaderByNumber(ctx, cfg.BlockNumber): nil = latest
- header.Time: Unix timestamp; time.Unix(sec, 0) for time.Time
- time.Since(blockTime).Seconds(): lag in seconds; negative clamped to 0

The exercise builds understanding of:
- Block freshness: stale = sync issues, clock drift, or node overload
- cfg.MaxLagSeconds: default 60; configurable threshold
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
