//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package storage

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"context"
)

var zeroHash = common.Hash{}
// TODO: implement Run.
func Run(ctx context.Context, client StorageClient, cfg Config) (*Result, error) {
	panic("TODO: implement")
}
// TODO: implement slotToHash.
func slotToHash(slot *big.Int) common.Hash { panic("TODO: implement") }
// TODO: implement mappingSlotHash.
func mappingSlotHash(key []byte, slot common.Hash) common.Hash { panic("TODO: implement") }
