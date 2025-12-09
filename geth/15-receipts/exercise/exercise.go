//go:build !solution
// +build !solution

package exercise

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
	// TODO: Validate input parameters
	// - Check if ctx is nil and provide context.Background() as default
	// - Check if client is nil and return an appropriate error
	// - Validate that cfg.TxHash is not the zero hash
	// Why validate? Receipt fetching is an RPC call; validate before network operations

	// TODO: Fetch the transaction receipt
	// - Call client.TransactionReceipt(ctx, cfg.TxHash)
	// - Handle potential errors from the RPC call
	// - Validate that the receipt response is not nil
	// Receipts contain execution results: status, gas used, logs, contract address

	// TODO: Map logs to LogSummary with defensive copying
	// - Create a logs slice with capacity len(rcpt.Logs)
	// - Loop through rcpt.Logs and for each log:
	//   - Copy Topics using append([]common.Hash(nil), lg.Topics...)
	//   - Copy Data using append([]byte(nil), lg.Data...)
	//   - Create LogSummary with Address, copied Topics, copied Data, and Index
	//   - Append to logs slice
	// Why defensive copying? Topics and Data are slices (reference types)

	// TODO: Construct and return the Result
	// - Create a Result struct with:
	//   - TxHash: rcpt.TxHash
	//   - BlockNumber: new(big.Int).Set(rcpt.BlockNumber) (defensive copy!)
	//   - StatusOK: rcpt.Status == 1 (1 = success, 0 = failure)
	//   - GasUsed: rcpt.GasUsed
	//   - CumulativeGas: rcpt.CumulativeGasUsed
	//   - Contract: rcpt.ContractAddress (if contract creation, otherwise zero address)
	//   - Logs: The copied logs slice
	//   - PostStateRoot: append([]byte(nil), rcpt.PostState...) (defensive copy!)
	// - Return the result and nil error on success

	return nil, errors.New("not implemented")
}
