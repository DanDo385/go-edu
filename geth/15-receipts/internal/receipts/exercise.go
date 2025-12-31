//go:build !solution
// +build !solution

package receipts

import (
	"context"
	"errors"
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
func Run(ctx context.Context, client ReceiptClient, cfg Config) (*Result, error) {
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
}

