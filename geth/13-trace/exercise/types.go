package exercise

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

// TraceClient captures the minimal debug API surface we need.
// A small interface keeps the exercise easy to mock and avoids
// depending on a concrete RPC client implementation.
type TraceClient interface {
	TraceTransaction(ctx context.Context, txHash common.Hash) (json.RawMessage, error)
}

// Config controls which transaction to trace.
type Config struct {
	TxHash common.Hash
}

// Result contains the raw trace payload from the node. We intentionally
// keep it as JSON so callers can pretty-print or further decode steps.
type Result struct {
	TxHash common.Hash
	Trace  json.RawMessage
}
