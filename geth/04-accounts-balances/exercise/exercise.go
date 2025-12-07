//go:build !solution
// +build !solution

package exercise

import (
	"context"
)

// Run queries balances + code for each target address.
// TODO: 1. Validate inputs (ctx/client non-nil and at least one address).
// TODO: 2. Loop through cfg.Addresses and call BalanceAt/CodeAt for each.
// TODO: 3. Classify account type based on presence of bytecode.
// TODO: 4. Copy balance/code data to avoid mutating client buffers.
// TODO: 5. Aggregate results into Result.Accounts and return.
func Run(ctx context.Context, client AccountClient, cfg Config) (*Result, error) {
	panic("TODO: implement Run for 04-accounts-balances")
}
