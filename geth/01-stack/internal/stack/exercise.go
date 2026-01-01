//go:build !solution && !reference

package stack

import (
	"context"
	"errors"
)

/*
Problem: Prove RPC connectivity by reading the network identifiers and latest header.

The very first thing an Ethereum Go tool should do is dial an RPC endpoint,
retrieve the chain/network IDs (replay protection + legacy identifier), and
inspect a block header. Headers are lightweight (~500 bytes) yet contain the
state root, parent hash, and other cryptographic commitments that define the
execution stack you are about to interact with. This function mirrors the CLI
demo from module 01 but exposes it as a reusable library API.

Computer science principles highlighted:
  - Separation of configuration from code (cfg.BlockNumber allows deterministic tests)
  - Fault tolerance via context propagation—callers control cancellation/timeouts
  - Immutability via defensive copies (we never hand pointers owned by go-ethereum back to callers)
*/
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
}

