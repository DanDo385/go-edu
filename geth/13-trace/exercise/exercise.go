//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// Run contains the reference solution for module 13-trace.
// Tracing replays a transaction in the EVM and returns structured JSON
// describing every call/gas step. We keep the interface tiny so students
// can focus on the protocol mechanics rather than client plumbing.
func Run(ctx context.Context, client TraceClient, cfg Config) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if cfg.TxHash == (common.Hash{}) {
		return nil, errors.New("tx hash required")
	}

	raw, err := client.TraceTransaction(ctx, cfg.TxHash)
	if err != nil {
		return nil, fmt.Errorf("trace transaction: %w", err)
	}
	if raw == nil {
		return nil, errors.New("nil trace payload")
	}

	traceCopy := make(json.RawMessage, len(raw))
	copy(traceCopy, raw)

	return &Result{
		TxHash: cfg.TxHash,
		Trace:  traceCopy,
	}, nil
}
