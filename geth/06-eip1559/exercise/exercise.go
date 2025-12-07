//go:build !solution
// +build !solution

package exercise

import "context"

// Run is the student entry point for module 06-eip1559.
// TODO: 1. Validate cfg (keys, fee caps) plus ctx/client.
// TODO: 2. Resolve sender address + nonce (pending or override).
// TODO: 3. Fetch base fee (latest header) and priority tip suggestion if needed.
// TODO: 4. Compute maxFeePerGas default (e.g., base*2 + tip) when not provided.
// TODO: 5. Construct a types.DynamicFeeTx, sign with the London signer, and optionally broadcast.
// TODO: 6. Return the signed tx + metadata for upstream display/inspection.
func Run(ctx context.Context, client FeeClient, cfg Config) (*Result, error) {
	panic("TODO: implement Run for 06-eip1559")
}
