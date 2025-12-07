//go:build !solution
// +build !solution

package exercise

import "context"

// Run is the student entry point for module 01-stack.
// TODO: 1. Validate inputs (ctx/client non-nil) and normalize cfg.
// TODO: 2. Query the chain ID (EIP-155 replay protection identifier).
// TODO: 3. Query the legacy network ID via net_version.
// TODO: 4. Fetch the latest block header using HeaderByNumber.
// TODO: 5. Return a Result with copies of the retrieved values.
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	panic("TODO: implement Run for 01-stack")
}
