//go:build !solution
// +build !solution

package exercise

import "context"

// Run is the student entry point for module 08-abigen.
// TODO: 1. Validate cfg (ABI JSON + contract address + caller).
// TODO: 2. Construct a BoundContract-esque wrapper (via bind.BoundContract).
// TODO: 3. Build CallOpts (with context, optional block number, from address).
// TODO: 4. Call typed ERC20 methods (Name, Symbol, Decimals, BalanceOf) using helpers.
// TODO: 5. Aggregate and return decoded values in a Result.
func Run(ctx context.Context, backend ContractCaller, cfg Config) (*Result, error) {
	panic("TODO: implement Run for 08-abigen")
}
