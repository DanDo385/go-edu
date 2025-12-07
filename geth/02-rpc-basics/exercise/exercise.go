//go:build !solution
// +build !solution

package exercise

import "context"

// Run is the student entry point for module 02-rpc-basics.
// TODO: 1. Validate inputs (non-nil context/client) and normalize cfg.Retries.
// TODO: 2. Query the latest block number using client.BlockNumber.
// TODO: 3. Query the legacy network ID via client.NetworkID.
// TODO: 4. Fetch the latest full block using client.BlockByNumber with retry logic.
// TODO: 5. Return a Result summarizing the block/metadata.
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	panic("TODO: implement Run for 02-rpc-basics")
}
