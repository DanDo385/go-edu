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
Reference Solution - Transaction Receipt
========================================

This file demonstrates fetching a transaction receipt via eth_getTransactionReceipt.
Receipts confirm inclusion and outcome: status, gas used, logs, contract address.

This connects to the Ethereum ecosystem by showing:
- TransactionReceipt(ctx, txHash): status, gasUsed, cumulativeGas, logs
- rcpt.Status == ReceiptStatusSuccessful: 1 = success, 0 = reverted
- rcpt.Logs: address, topics, data — same structure as FilterLogs
- rcpt.ContractAddress: set for contract creation txs

The exercise builds understanding of:
- Receipt vs tx: receipt is post-execution; includes logs and status
- Defensive copies: Topics, Data, PostState — RPC may reuse buffers
- blockNumber: new(big.Int).Set for *big.Int copy

Teaching notes (per .cursorrules):
- append([]common.Hash(nil), lg.Topics...): copy slice so Result owns data.
- rcpt.PostState: legacy pre-Byzantium; we copy as bytes for consistency.
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
