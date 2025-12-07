//go:build !solution
// +build !solution

package exercise

import "context"

// Run is the student entry point for module 10-filters.
// TODO: 1. Validate cfg/client and decide between subscription vs polling.
// TODO: 2. For subscription: SubscribeNewHead, receive headers, detect reorgs.
// TODO: 3. For polling: repeatedly call HeaderByNumber(nil) until new head arrives.
// TODO: 4. Track previous head hash to flag reorg when parent hash mismatches.
// TODO: 5. Collect up to cfg.MaxHeads heads and return them in Result.
func Run(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	panic("TODO: implement Run for 10-filters")
}
