//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package smartcontracts

import (
	"math/big"
	"context"
)
// TODO: implement Run.
func Run(ctx context.Context, client ContractCaller, cfg Config) (*Result, error) {
	panic("TODO: implement")
}
// TODO: implement selector.
func selector(sig string) []byte { panic("TODO: implement") }
// TODO: implement decodeString.
func decodeString(data []byte) (string, error) { panic("TODO: implement") }
// TODO: implement decodeUint8.
func decodeUint8(data []byte) (uint8, error) { panic("TODO: implement") }
// TODO: implement decodeUint256.
func decodeUint256(data []byte) (*big.Int, error) { panic("TODO: implement") }
