//go:build !solution
// +build !solution

package exercise

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
	// TODO: Validate input parameters
	// - Check if ctx is nil and provide a default context if needed
	// - Check if client is nil and return an appropriate error
	// - Why: Standard defensive programming pattern from previous modules

	// TODO: Query pending transaction count
	// - Call client.PendingTransactionCount(ctx) to get the count
	// - Handle potential errors from the RPC call
	// - Note: This may not be supported by all RPC providers
	// - Why: Understanding mempool size helps assess network congestion
	// - Key concepts:
	//   * Mempool (txpool) holds unconfirmed transactions
	//   * Transactions wait here until miners/validators include them
	//   * Size indicates network congestion (more pending = higher gas prices)

	// TODO: Understand mempool limitations
	// - Not all RPC endpoints expose mempool data:
	//   * eth_pendingTransactions: Often restricted or unavailable
	//   * txpool_content: Requires txpool API enabled (rare on public RPCs)
	//   * txpool_status: Basic stats (sometimes available)
	// - Why hidden?
	//   * Privacy: Prevents front-running and MEV exploitation
	//   * Security: Prevents DoS attacks targeting mempool
	//   * Resources: Full mempool dump is expensive
	// - Transaction replacement rules:
	//   * Same nonce, higher gas price can replace pending tx
	//   * Minimum 10% fee increase typically required
	//   * Useful for "stuck" transactions

	// TODO: Construct and return the Result struct
	// - Store the pending count
	// - Return the result with nil error on success
	// - Note: For full transaction details, need additional RPC methods

	return nil, errors.New("not implemented")
}
