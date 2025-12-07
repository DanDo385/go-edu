//go:build !solution
// +build !solution

package exercise

import "context"

// Run is the student entry point for module 07-eth-call.
// TODO: 1. Validate cfg (contract address) plus ctx/client.
// TODO: 2. Encode zero-arg function selectors for name/symbol/decimals/totalSupply.
// TODO: 3. Issue eth_call via CallContract for each selector (one CallMsg per function).
// TODO: 4. Decode dynamic string outputs manually (offset, length, padded bytes).
// TODO: 5. Decode uint8/uint256 outputs for decimals/total supply.
// TODO: 6. Aggregate results (and optionally raw responses) into a Result.
func Run(ctx context.Context, client CallClient, cfg Config) (*Result, error) {
	panic("TODO: implement Run for 07-eth-call")
}
