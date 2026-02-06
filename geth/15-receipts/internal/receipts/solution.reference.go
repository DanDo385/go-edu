//go:build reference

package receipts

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var errNilClient = errors.New("nil receipt client")

/*
Reference Solution

Structure:
- Fetch one transaction receipt.
- Normalize into a copy-safe Result shape.

Pointer notes:
- Deep-copy slices (`PostState`, log topics/data) to avoid aliasing mutable memory.
- Copy `*big.Int` block number for ownership clarity.
*/
func Run(ctx context.Context, client ReceiptClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if cfg.TxHash == (common.Hash{}) {
		return nil, errors.New("transaction hash is required")
	}

	rcpt, err := client.TransactionReceipt(ctx, cfg.TxHash)
	if err != nil {
		return nil, fmt.Errorf("transaction receipt: %w", err)
	}
	if rcpt == nil {
		return nil, errors.New("transaction receipt: nil response")
	}

	logs := make([]LogSummary, 0, len(rcpt.Logs))
	for _, lg := range rcpt.Logs {
		if lg == nil {
			continue
		}
		logs = append(logs, LogSummary{
			Address: lg.Address,
			Topics:  append([]common.Hash(nil), lg.Topics...),
			Data:    append([]byte(nil), lg.Data...),
			Index:   lg.Index,
		})
	}

	var blockNumber *big.Int
	if rcpt.BlockNumber != nil {
		blockNumber = new(big.Int).Set(rcpt.BlockNumber)
	}

	return &Result{
		TxHash:        cfg.TxHash,
		BlockNumber:   blockNumber,
		StatusOK:      rcpt.Status == types.ReceiptStatusSuccessful,
		GasUsed:       rcpt.GasUsed,
		CumulativeGas: rcpt.CumulativeGasUsed,
		Contract:      rcpt.ContractAddress,
		Logs:          logs,
		PostStateRoot: append([]byte(nil), rcpt.PostState...),
	}, nil
}
