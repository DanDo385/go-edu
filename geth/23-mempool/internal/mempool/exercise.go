//go:build !solution && !reference

package mempool

import (
	"context"
	"errors"
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
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
}

