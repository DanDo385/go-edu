//go:build !solution
// +build !solution

package exercise

import "context"

// Run is the student entry point for module 11-storage.
// TODO: 1. Validate cfg/client and ensure Slot is provided.
// TODO: 2. Compute the canonical 32-byte slot hash (with mapping key hashing if provided).
// TODO: 3. Invoke StorageAt for the contract/slot/block.
// TODO: 4. Return the resolved slot hash and raw 32-byte value to the caller.
func Run(ctx context.Context, client StorageClient, cfg Config) (*Result, error) {
	panic("TODO: implement Run for 11-storage")
}
