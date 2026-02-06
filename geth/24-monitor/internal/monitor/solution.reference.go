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

Structure:
- Fetch target header (latest when block number is nil).
- Compute lag against wall-clock time.
- Classify node status based on configured lag threshold.
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
