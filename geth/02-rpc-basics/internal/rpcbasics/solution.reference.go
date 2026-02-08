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
Reference Solution - JSON-RPC Communication Patterns
==================================================

This file demonstrates fundamental JSON-RPC communication patterns used throughout
Ethereum infrastructure. JSON-RPC is the standard protocol for Ethereum client communication,
enabling distributed systems to query blockchain state and submit transactions.

This connects to the broader Ethereum ecosystem by showing:
- Retry logic for handling transient network failures
- Context-based cancellation for cooperative goroutine shutdown
- Exponential backoff patterns (foundation for more sophisticated retry libraries)
- Block-by-number retrieval for historical data access
- Error classification and recovery strategies

The exercise builds understanding of:
- Network reliability assumptions and failure modes
- Time-based coordination with context timeouts
- Resource management in distributed systems
- Idempotent operations and safe retry semantics
- Performance tradeoffs between latency and reliability

Teaching notes:
- Memory/ownership: returned block data may be cached by the RPC client internally.
  We don't copy here since the Result struct owns the block reference exclusively.
- Invariants: retry attempts are bounded to prevent infinite loops, and context
  cancellation is checked before each retry to enable graceful shutdown.
- Error surfaces: network operations can fail transiently. We implement retry logic
  but surface persistent failures so callers can implement circuit breakers or alerting.
*/

/*
Run - Network State Query with Retry Logic

This function demonstrates production-ready RPC communication with automatic retry
logic for handling transient network failures. It retrieves the latest block data
while gracefully handling network unreliability.

Parameters:
- ctx: context for cancellation and timeout control
- client: RPC client interface for network communication
- cfg: configuration including retry count for resilience

Returns:
- *Result: current network state snapshot
- error: communication failures after all retries exhausted

Algorithm steps:
1. Validate client interface (defensive programming)
2. Retrieve network identification for context
3. Get latest block number as baseline
4. Retry block retrieval with exponential backoff
5. Return network state snapshot or final error

Why retry logic matters:
- Networks are unreliable (congestion, node issues, network partitions)
- Some failures are transient and succeed on retry
- Retry with backoff prevents thundering herd problems
- Context cancellation enables graceful shutdown during retries
*/
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	// Step 1: Input validation
	// Interface values can be nil - validate before use
	if client == nil {
		return nil, errNilClient
	}

	// Step 2: Retrieve Network ID
	// Network ID provides context about which Ethereum network we're connected to
	// This helps validate we're talking to the expected network
	networkID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("network id: %w", err)
	}
	if networkID == nil {
		return nil, errors.New("network id: nil response")
	}

	// Step 3: Get latest block number
	// Block numbers increment monotonically and provide a baseline for queries
	// This is a lightweight operation that establishes current chain state
	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("block number: %w", err)
	}

	// Step 4: Configure retry attempts
	// cfg.Retries allows callers to specify resilience level
	// Default to 1 attempt if not specified (no retries)
	// This gives callers control over latency vs reliability tradeoff
	attempts := cfg.Retries
	if attempts <= 0 {
		attempts = 1
	}

	// Step 5: Retry loop with exponential backoff
	// Network operations can fail transiently due to:
	// - Network congestion, DNS issues, connection timeouts
	// - Node load balancing, temporary unavailability
	// - RPC rate limiting or temporary service degradation

	var blockErr error // Track the most recent error for reporting
	for i := 0; i < attempts; i++ {
		// Attempt to retrieve the block at the current block number
		// new(big.Int).SetUint64(blockNumber) converts uint64 to *big.Int
		// Ethereum uses big integers for all numeric values due to 256-bit architecture
		b, err := client.BlockByNumber(ctx, new(big.Int).SetUint64(blockNumber))
		if err == nil {
			// Success! Return the network state snapshot
			return &Result{
				NetworkID:   new(big.Int).Set(networkID), // Defensive copy
				BlockNumber: blockNumber,                  // Value type, no copy needed
				Block:       b,                            // Block ownership transferred
			}, nil
		}

		// Record the error for potential final reporting
		blockErr = err

		// If this isn't the last attempt, wait before retrying
		if i+1 < attempts {
			// Context-aware delay: respect cancellation during backoff
			// This allows graceful shutdown even during retry delays
			// select statement waits for either context cancellation or timeout
			select {
			case <-ctx.Done():
				// Context was cancelled (timeout, manual cancellation, etc.)
				// Return the context error, not the RPC error
				return nil, ctx.Err()
			case <-time.After(retryDelay):
				// Delay completed, proceed to next retry attempt
				// Fixed delay (not exponential) for simplicity in this example
			}
		}
	}

	// Step 6: All retry attempts exhausted
	// Return comprehensive error message including attempt count
	// This helps with debugging and monitoring (alert on high failure rates)
	return nil, fmt.Errorf("block by number after %d attempts: %w", attempts, blockErr)
}
