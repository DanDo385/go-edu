//go:build !solution
// +build !solution

package exercise

import "context"

// Run is the student entry point for module 05-tx-nonces.
// TODO: 1. Validate cfg (private key, recipient, amount) plus ctx/client.
// TODO: 2. Determine sender nonce (PendingNonceAt unless override provided).
// TODO: 3. Resolve chain ID + gas price (suggestion unless override provided).
// TODO: 4. Build a legacy transaction (types.NewTransaction).
// TODO: 5. Sign with EIP-155 protection and broadcast via SendTransaction.
// TODO: 6. Return the signed transaction + metadata so callers can display hashes.
func Run(ctx context.Context, client TXClient, cfg Config) (*Result, error) {
	panic("TODO: implement Run for 05-tx-nonces")
}
