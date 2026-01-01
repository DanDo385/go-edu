//go:build !solution && !reference

package monitor

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
	// TODO: Implement

	// ============================================================================
	// Standard validation pattern, plus setting a sensible default for MaxLagSeconds.
	// TODO: Implement

	// ============================================================================
	// STEP 2: Fetch Latest Block Header - Timestamp is Key
	// TODO: Implement

	// ============================================================================
	// We fetch the block header to get its timestamp. The timestamp tells us
	// TODO: Implement

	// ============================================================================
	// STEP 3: Calculate Block Lag - Time-Based Health Check
	// TODO: Implement

	// ============================================================================
	// Block lag is the difference between current time and block production time.
	// TODO: Implement

	// ============================================================================
	// STEP 4: Classify Status - Threshold-Based Alerting
	// TODO: Implement

	// ============================================================================
	// We classify the node as OK or STALE based on the lag threshold.
	// TODO: Implement

	// ============================================================================
	// STEP 5: Return Comprehensive Result
	// TODO: Implement

	// ============================================================================
	// We return complete diagnostic information: status, block details, and lag.
	// TODO: Implement

	panic("unimplemented")
}
