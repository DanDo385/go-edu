//go:build !solution
// +build !solution

package exercise

import "context"

// Run is the student entry point for module 12-proofs.
// TODO: 1. Validate cfg (account, optional slots) and ctx/client.
// TODO: 2. Convert requested slot hashes to hex strings for eth_getProof.
// TODO: 3. Invoke client.GetProof to fetch account + storage proofs.
// TODO: 4. Transform the response into an AccountProof-friendly structure.
func Run(ctx context.Context, client ProofClient, cfg Config) (*Result, error) {
	panic("TODO: implement Run for 12-proofs")
}
