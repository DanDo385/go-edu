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
Reference Solution

Structure:
- Validate tx hash input.
- Fetch raw trace JSON.
- Return a defensive copy of raw bytes.

Pointer notes:
- `json.RawMessage` is a `[]byte`; slicing aliases backing arrays.
- `append(nil, raw...)` forces ownership transfer to the result.
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
