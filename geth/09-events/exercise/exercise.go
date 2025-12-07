//go:build !solution
// +build !solution

package exercise

import "context"

// Run is the student entry point for module 09-events.
// TODO: 1. Validate cfg (token, block range) plus ctx/client.
// TODO: 2. Build a FilterQuery with Transfer topic + optional from/to filters.
// TODO: 3. Call FilterLogs and iterate results.
// TODO: 4. Decode indexed params from topics (addresses).
// TODO: 5. Decode uint256 value from data via ABI rules.
// TODO: 6. Return a slice of TransferEvent structures summarizing each log.
func Run(ctx context.Context, client LogClient, cfg Config) (*Result, error) {
	panic("TODO: implement Run for 09-events")
}
