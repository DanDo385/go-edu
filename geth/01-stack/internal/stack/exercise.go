//go:build !solution && !reference

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
Problem: Prove RPC connectivity by reading the network identifiers and latest header.
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
