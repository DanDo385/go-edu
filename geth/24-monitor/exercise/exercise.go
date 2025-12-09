//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"errors"
)

/*
Problem: Implement node health monitoring by checking block freshness and detecting lag.

Monitoring nodes is critical for production systems. A stale node (not receiving new blocks)
will return outdated data, causing issues for applications. By comparing the latest block's
timestamp to the current time, we can detect if a node is lagging behind the network.

Computer science principles highlighted:
  - Time-based health checks (staleness detection)
  - Threshold-based alerting (classify OK vs STALE)
  - Observability patterns (monitoring system health)
*/
func Run(ctx context.Context, client MonitorClient, cfg Config) (*Result, error) {
	// TODO: Validate input parameters
	// - Check if ctx is nil and provide a default context
	// - Check if client is nil and return an error
	// - Validate cfg.MaxLagSeconds has a sensible default (e.g., 60 seconds)
	// - Why: Need threshold to classify health status

	// TODO: Fetch the latest block header
	// - Call client.HeaderByNumber(ctx, cfg.BlockNumber)
	// - Use cfg.BlockNumber (nil for latest block)
	// - Handle potential errors from RPC call
	// - Validate header is not nil
	// - Why: Block timestamp tells us when the last block was produced

	// TODO: Calculate block lag
	// - Get current time: time.Now()
	// - Get block time: time.Unix(header.Time, 0)
	// - Calculate lag: currentTime - blockTime
	// - Convert to seconds for threshold comparison
	// - Why: Lag indicates how far behind the node is from real-time
	// - Key concepts:
	//   * Block time is Unix timestamp (seconds since epoch)
	//   * Ethereum mainnet: ~12 second block time
	//   * Acceptable lag: < 30-60 seconds for synced node
	//   * High lag: > 60 seconds indicates sync issues or stale RPC

	// TODO: Classify node status based on lag
	// - If lag < MaxLagSeconds: Status = "OK"
	// - If lag >= MaxLagSeconds: Status = "STALE"
	// - Why: Threshold-based classification enables alerting
	// - Production considerations:
	//   * Set threshold based on block time (mainnet ~12s, testnet varies)
	//   * Alert if status == STALE for multiple consecutive checks
	//   * Track lag over time to detect degradation trends

	// TODO: Construct and return the Result struct
	// - Include Status, BlockNumber, BlockTimestamp, LagSeconds
	// - Provide complete diagnostic information for monitoring
	// - Return with nil error on success

	return nil, errors.New("not implemented")
}
