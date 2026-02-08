//go:build reference

package trace

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var errNilClient = errors.New("nil trace client")

/*
Reference Solution - debug_traceTransaction
===========================================

This file demonstrates debug_traceTransaction: returns execution trace (call
frames, opcodes, gas) for a transaction. Requires a node with debug API enabled.

This connects to the Ethereum ecosystem by showing:
- TraceTransaction(txHash): raw trace output (format depends on tracer type)
- Use cases: debugging failed txs, gas profiling, internal call analysis
- append([]byte(nil), raw...): defensive copy of trace bytes

The exercise builds understanding of:
- Trace as opaque bytes: geth supports multiple tracer types (callTracer, etc.)
- cfg.TxHash: common.Hash, zero value = invalid
- Debug API: not all nodes expose it; often disabled on public RPCs

Teaching notes (per .cursorrules):
- Trace output format varies; we preserve raw bytes for caller to interpret.
- Memory: RPC client may reuse buffer; we copy to ensure Result outlives call.
*/
func Run(ctx context.Context, client TraceClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if cfg.TxHash == (common.Hash{}) {
		return nil, errors.New("transaction hash is required")
	}

	raw, err := client.TraceTransaction(ctx, cfg.TxHash)
	if err != nil {
		return nil, fmt.Errorf("trace transaction: %w", err)
	}

	return &Result{
		TxHash: cfg.TxHash,
		Trace:  append([]byte(nil), raw...),
	}, nil
}
