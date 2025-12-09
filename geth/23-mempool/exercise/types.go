package exercise

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

// MempoolClient captures the ethclient calls needed for module 23 (pending transactions).
type MempoolClient interface {
	// PendingTransactionCount returns the number of pending transactions in the mempool
	PendingTransactionCount(ctx context.Context) (uint, error)
}

// Config allows configuration for mempool inspection.
type Config struct {
	// Limit specifies the maximum number of transactions to return (0 = no limit)
	Limit int
}

// Result summarizes the mempool status of the node.
type Result struct {
	// PendingCount is the total number of pending transactions
	PendingCount uint
	// Note: Full transaction details would require additional RPC methods
	// like eth_pendingTransactions or txpool_content (often restricted)
}
