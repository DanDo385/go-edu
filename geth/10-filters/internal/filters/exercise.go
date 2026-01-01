//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package filters

import (
	"context"

	"time"
)

const defaultMaxHeads = 5
const defaultPollInterval = time.Second
// TODO: implement Run.
func Run(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	panic("TODO: implement")
}
// TODO: implement subscribeHeads.
func subscribeHeads(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	panic("TODO: implement")
}
// TODO: implement pollHeads.
func pollHeads(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	panic("TODO: implement")
}
