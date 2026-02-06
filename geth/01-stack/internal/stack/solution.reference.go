//go:build reference

package stack

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

var errNilClient = errors.New("nil rpc client")

/*
Reference Solution

Structure:
- Validate dependencies first.
- Fetch chain/network metadata.
- Fetch header for the requested block (nil means latest).

Invariants:
- Returned values must not alias mutable values held by the caller/client.
- A nil client is always rejected.

Pointer notes:
- `*big.Int` and `*types.Header` are mutable pointer-backed values.
- We copy them before returning so callers cannot mutate upstream state.
*/
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}
	if chainID == nil {
		return nil, errors.New("chain id: nil response")
	}

	networkID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("network id: %w", err)
	}
	if networkID == nil {
		return nil, errors.New("network id: nil response")
	}

	header, err := client.HeaderByNumber(ctx, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("header by number: %w", err)
	}
	if header == nil {
		return nil, errors.New("header by number: nil response")
	}

	return &Result{
		ChainID:   new(big.Int).Set(chainID),
		NetworkID: new(big.Int).Set(networkID),
		Header:    types.CopyHeader(header),
	}, nil
}
