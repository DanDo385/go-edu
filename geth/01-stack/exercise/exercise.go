//go:build !solution
// +build !solution

package exercise

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
	// TODO: Validate input parameters
	// - Check if ctx is nil and provide a default context if needed
	// - Check if client is nil and return an appropriate error

	// TODO: Retrieve Chain ID from the RPC client
	// - Call client.ChainID(ctx) to get the chain identifier
	// - Handle potential errors from the RPC call
	// - Validate that the chainID response is not nil
	// - Chain ID is critical for replay protection (EIP-155)

	// TODO: Retrieve Network ID from the RPC client
	// - Call client.NetworkID(ctx) to get the network identifier
	// - Handle potential errors from the RPC call
	// - Validate that the networkID response is not nil
	// - Network ID is a legacy identifier used for P2P networking

	// TODO: Retrieve block header from the RPC client
	// - Call client.HeaderByNumber(ctx, cfg.BlockNumber) to get the header
	// - Use cfg.BlockNumber (can be nil for latest block)
	// - Handle potential errors from the RPC call
	// - Validate that the header response is not nil

	// TODO: Construct and return the Result struct
	// - Create a new Result struct
	// - IMPORTANT: Use defensive copying for all fields:
	//   - Use new(big.Int).Set(chainID) to copy ChainID (big.Int is mutable!)
	//   - Use new(big.Int).Set(networkID) to copy NetworkID
	//   - Use types.CopyHeader(header) to copy the Header
	// - Return the result and nil error on success
	// - Why defensive copying? We don't want callers mutating our internal data

	return nil, errors.New("not implemented")
}
