//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package rpcbasics

import (
	"context"

	"time"
)

const retryDelay = 100 * time.Millisecond
// TODO: implement Run.
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	panic("TODO: implement")
}
