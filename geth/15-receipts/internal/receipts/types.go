package receipts

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ReceiptClient captures the ethclient calls needed for the receipts module.
type ReceiptClient interface {
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

// Config selects which transaction receipt to fetch.
type Config struct {
	TxHash common.Hash
}

// LogSummary is a minimal, defensively-copied view of a log.
type LogSummary struct {
	Address common.Address
	Topics  []common.Hash
	Data    []byte
	Index   uint
}

// Result summarizes receipt data in a stable, testable shape.
type Result struct {
	TxHash        common.Hash
	BlockNumber   *big.Int
	StatusOK      bool
	GasUsed       uint64
	CumulativeGas uint64
	Contract      common.Address
	Logs          []LogSummary
	PostStateRoot []byte
}
