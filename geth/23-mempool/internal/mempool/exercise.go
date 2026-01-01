//go:build !solution && !reference

package mempool

import (
	"context"
	"errors"
	"fmt"
)

/*
Problem: Inspect the mempool (transaction pool) to understand pending transactions.

The mempool contains transactions that have been broadcast to the network but not
yet included in a block. Monitoring the mempool helps you understand network congestion,
estimate gas prices, and track your own transactions.

However, mempool visibility is limited for privacy and security reasons. Many public
RPC endpoints don't expose pending transaction details.

Computer science principles highlighted:
  - Queue management (FIFO with priority)
  - Privacy/security trade-offs (transparency vs exploitation)
  - Resource management (mempool size limits)
*/
func Run(ctx context.Context, client MempoolClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Familiar Pattern
	// TODO: Implement

	// ============================================================================
	// By now, this validation pattern should be second nature. We validate
	// TODO: Implement

	// ============================================================================
	// STEP 2: Query Mempool Size - Understanding Transaction Pools
	// TODO: Implement

	// ============================================================================
	// The mempool (also called txpool) is where transactions wait before being
	// TODO: Implement

	// ============================================================================
	// STEP 3: Interpret Mempool Size - Understanding Congestion
	// TODO: Implement

	// ============================================================================
	// The pending transaction count tells us about network activity and congestion.
	// TODO: Implement

	// ============================================================================
	// STEP 4: Return Result - Simple Count with Future Extensibility
	// TODO: Implement

	// ============================================================================
	// We return the pending count in a Result struct. While we only have one
	// TODO: Implement

	panic("unimplemented")
}
