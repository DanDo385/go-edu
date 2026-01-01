//go:build !solution && !reference

package rpcbasics

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
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
}

