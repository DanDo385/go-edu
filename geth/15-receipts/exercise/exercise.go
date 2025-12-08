//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ReceiptClient abstracts the single RPC we need.
type ReceiptClient interface {
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

// Config identifies which receipt to fetch.
type Config struct {
	TxHash common.Hash
}

// LogSummary is a lightweight view of a log entry.
type LogSummary struct {
	Address common.Address
	Topics  []common.Hash
	Data    []byte
	Index   uint
}

// Result captures decoded receipt data.
type Result struct {
	TxHash        common.Hash
	BlockNumber   *big.Int
	StatusOK      bool
	GasUsed       uint64
	Contract      common.Address
	Logs          []LogSummary
	CumulativeGas uint64
	PostStateRoot []byte
}

// Run is the student entry point for module 15-receipts.
//
// TODOs for students:
// 1) Validate inputs: ctx/client non-nil, cfg.TxHash set.
// 2) Call TransactionReceipt for the hash.
// 3) Copy receipt fields into Result (avoid returning go-ethereum pointers directly).
// 4) Map []*types.Log → []LogSummary (copy topics/data slices).
// 5) Return the populated Result or an error.
func Run(ctx context.Context, client ReceiptClient, cfg Config) (*Result, error) {
	panic("TODO: implement Run for 15-receipts")
}
