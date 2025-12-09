//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"time"
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
	// ============================================================================
	// STEP 1: Input Validation with Defaults
	// ============================================================================
	// Standard validation pattern, plus setting a sensible default for MaxLagSeconds.
	//
	// Building on modules 01, 21, 22, 23: Same validation pattern.
	if ctx == nil {
		ctx = context.Background()
	}

	if client == nil {
		return nil, errors.New("client is nil")
	}

	// Set default max lag if not specified
	// 60 seconds is reasonable for Ethereum mainnet (~12s block time = ~5 blocks)
	if cfg.MaxLagSeconds == 0 {
		cfg.MaxLagSeconds = 60
	}

	// ============================================================================
	// STEP 2: Fetch Latest Block Header - Timestamp is Key
	// ============================================================================
	// We fetch the block header to get its timestamp. The timestamp tells us
	// when the block was produced, which we compare to current time to detect lag.
	//
	// Why header instead of full block? Headers are lightweight (~500 bytes)
	// and contain the timestamp field we need. No need to download full block.
	//
	// Building on module 01: Same HeaderByNumber pattern, different use case.
	// Module 01 used header for chain metadata; here we use it for health monitoring.
	header, err := client.HeaderByNumber(ctx, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("header by number: %w", err)
	}

	if header == nil {
		return nil, errors.New("header response was nil")
	}

	// ============================================================================
	// STEP 3: Calculate Block Lag - Time-Based Health Check
	// ============================================================================
	// Block lag is the difference between current time and block production time.
	// High lag indicates the node isn't receiving new blocks quickly.
	//
	// How block timestamps work:
	//   - Miners/validators set timestamp when producing block
	//   - Timestamp is Unix time (seconds since Jan 1, 1970)
	//   - Network rejects blocks with timestamps too far in future
	//   - Timestamps can be slightly in the past (clock skew)
	//
	// Typical lag values:
	//   - Synced node: 0-30 seconds (just processed latest block)
	//   - Slightly behind: 30-60 seconds (caught in next few blocks)
	//   - Stale: > 60 seconds (sync issues, network problems, or stale RPC)
	//
	// Network-specific considerations:
	//   - Ethereum mainnet: ~12 second block time, 60s lag = ~5 blocks behind
	//   - Polygon: ~2 second block time, 60s lag = ~30 blocks behind
	//   - Optimism: ~2 second block time (L2), different sync model
	//
	// Why lag matters:
	//   - Applications need fresh data for accurate queries
	//   - Stale data causes failed transactions (nonce mismatches)
	//   - DeFi: Stale prices can cause liquidation/arbitrage failures
	//   - Wallets: Incorrect balance/nonce information
	//
	// Edge cases to consider:
	//   - Negative lag: Block timestamp in future (clock skew, usually OK)
	//   - Very high lag: Node stuck during initial sync
	//   - Fluctuating lag: Network congestion or connection issues
	currentTime := time.Now()
	blockTime := time.Unix(int64(header.Time), 0)
	lagSeconds := int64(currentTime.Sub(blockTime).Seconds())

	// ============================================================================
	// STEP 4: Classify Status - Threshold-Based Alerting
	// ============================================================================
	// We classify the node as OK or STALE based on the lag threshold.
	// This binary classification enables simple alerting logic.
	//
	// Alerting patterns:
	//   - Single check: Useful for one-time diagnostics
	//   - Consecutive failures: Alert only if STALE for N checks (reduce flapping)
	//   - Trending: Track lag over time, alert on increasing trend
	//   - Multi-node: Check multiple RPC endpoints, alert if all stale
	//
	// Production monitoring:
	//   - Expose as Prometheus metric (gauge for lag, counter for status)
	//   - Set up alerts in Grafana/PagerDuty
	//   - Track SLA: % time node is healthy over 7/30 days
	//   - Correlate with other metrics: sync status, peer count, RPC latency
	//
	// Why OK vs STALE (not more granular)?
	//   - Simple binary decision: "can I use this node?"
	//   - Could extend with WARNING state (30-60s lag)
	//   - Could add CRITICAL for very high lag (> 300s)
	status := "OK"
	if lagSeconds >= cfg.MaxLagSeconds {
		status = "STALE"
	}

	// Handle negative lag (clock skew - block timestamp in future)
	// This is usually harmless and shouldn't trigger STALE status
	if lagSeconds < 0 {
		status = "OK" // Accept slight clock skew
	}

	// ============================================================================
	// STEP 5: Return Comprehensive Result
	// ============================================================================
	// We return complete diagnostic information: status, block details, and lag.
	// This enables both automated alerting (status) and manual debugging (details).
	//
	// Result struct design:
	//   - Status: For programmatic alerting ("if status == STALE, alert")
	//   - BlockNumber: For debugging ("which block was checked?")
	//   - BlockTimestamp: For manual inspection ("when was block produced?")
	//   - LagSeconds: For trending/graphing ("how is lag changing over time?")
	//
	// No defensive copying needed: All fields are primitives or time.Time (immutable).
	//
	// Building on previous modules:
	//   - Module 01: Basic RPC call + result struct
	//   - Module 21: Boolean health status (IsSyncing)
	//   - Module 22: Simple metric (PeerCount)
	//   - Module 23: Simple metric (PendingCount)
	//   - Module 24: Complex health check (status + diagnostics)
	//
	// This progression shows how monitoring builds from simple metrics to
	// comprehensive health checks with actionable status and debugging data.
	return &Result{
		Status:         status,
		BlockNumber:    header.Number.Uint64(),
		BlockTimestamp: blockTime,
		LagSeconds:     lagSeconds,
	}, nil
}
