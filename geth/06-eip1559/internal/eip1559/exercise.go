//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package eip1559

import "context"

const defaultDynamicGasLimit = 21000
// TODO: implement Run.
func Run(ctx context.Context, client FeeClient, cfg Config) (*Result, error) {
	panic("TODO: implement")
}
