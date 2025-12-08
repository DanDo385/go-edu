//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"errors"
)

/*
Problem: Build a resilient RPC client that can fetch full block data.

In module 01, you fetched lightweight block headers. Now, you'll fetch full
blocks, which include all transaction data. This is a more expensive operation,
so you'll also implement retry logic to handle transient network errors. This
is a fundamental pattern for building robust Ethereum applications.

Computer science principles highlighted:
  - Fault tolerance via retries and context-awareness
  - Composition (a Block contains a Header and Transactions)
  - Immutability via defensive copying (building on module 01)
*/
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	// TODO: Validate input parameters
	// - Check for nil context and provide a default
	// - Check for nil client and return an error
	// - Ensure cfg.Retries is non-negative

	// TODO: Retrieve the latest block number
	// - Call client.BlockNumber(ctx)
	// - Handle errors, wrapping them with context

	// TODO: Retrieve the network ID
	// - Call client.NetworkID(ctx)
	// - Handle errors, wrapping them with context

	// TODO: Retrieve the full block with retry logic
	// - Use a for loop to retry up to cfg.Retries times
	// - Inside the loop, call client.BlockByNumber(ctx, nil)
	// - If the call is successful, break the loop
	// - If the call fails, check if the context is canceled. If so, return
	//   the context error.
	// - Wait for a short duration (e.g., 100ms) before the next attempt.
	//   Use time.After for the delay, wrapped in a select statement with the
	//   context's Done() channel.
	// - After the loop, if the block is still nil, return the last error encountered.

	// TODO: Construct and return the Result struct
	// - Create a new Result struct
	// - IMPORTANT: Use defensive copying for all fields:
	//   - Use new(big.Int).Set(networkID) to copy NetworkID
	//   - BlockNumber is a uint64, which is a value type, so no copy needed.
	//   - Use types.CopyBlock(block) to copy the Block
	// - Why defensive copying? `types.Block` contains pointers. If we don't
	//   copy it, the caller could mutate the client's internal data.

	return nil, errors.New("not implemented")
}