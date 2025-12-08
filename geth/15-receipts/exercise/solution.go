//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Run contains the reference solution for fetching and summarizing receipts.
//
// Concepts:
// - Single-call happy path: TransactionReceipt is your anchor for execution results.
// - Defensive copies: receipts/logs contain byte slices; copy them so callers can’t mutate shared buffers.
// - Status vs postState: Byzantium introduced status (1 = success, 0 = revert). Earlier chains used PostState.
func Run(ctx context.Context, client ReceiptClient, cfg Config) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if cfg.TxHash == (common.Hash{}) {
		return nil, errors.New("tx hash required")
	}

	rcpt, err := client.TransactionReceipt(ctx, cfg.TxHash)
	if err != nil {
		return nil, fmt.Errorf("transaction receipt: %w", err)
	}
	if rcpt == nil {
		return nil, errors.New("nil receipt")
	}

	logs := make([]LogSummary, 0, len(rcpt.Logs))
	for _, lg := range rcpt.Logs {
		topicsCopy := append([]common.Hash(nil), lg.Topics...)
		dataCopy := append([]byte(nil), lg.Data...)
		logs = append(logs, LogSummary{
			Address: lg.Address,
			Topics:  topicsCopy,
			Data:    dataCopy,
			Index:   uint(lg.Index),
		})
	}

	return &Result{
		TxHash:        rcpt.TxHash,
		BlockNumber:   new(big.Int).Set(rcpt.BlockNumber),
		StatusOK:      rcpt.Status == 1,
		GasUsed:       rcpt.GasUsed,
		CumulativeGas: rcpt.CumulativeGasUsed,
		Contract:      rcpt.ContractAddress,
		Logs:          logs,
		PostStateRoot: append([]byte(nil), rcpt.PostState...),
	}, nil
}
