//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package events

import (
	"github.com/ethereum/go-ethereum/common"
	"context"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/core/types"
)

var transferSigHash = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
// TODO: implement Run.
func Run(ctx context.Context, client LogClient, cfg Config) (*Result, error) {
	panic("TODO: implement")
}
// TODO: implement addressTopic.
func addressTopic(addr common.Address) common.Hash { panic("TODO: implement") }
// TODO: implement decodeTransferLog.
func decodeTransferLog(lg types.Log) (TransferEvent, error) { panic("TODO: implement") }
