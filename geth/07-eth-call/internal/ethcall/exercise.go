//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package ethcall

import (
	"math/big"
	"context"
)

var (
	selectorName        = selector("name()")
	selectorSymbol      = selector("symbol()")
	selectorDecimals    = selector("decimals()")
	selectorTotalSupply = selector("totalSupply()")
)
// TODO: implement Run.
func Run(ctx context.Context, client CallClient, cfg Config) (*Result, error) {
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
