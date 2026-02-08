//go:build reference

package accountsbalances

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var errNilClient = errors.New("nil account client")

/*
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
*/

func Run(ctx context.Context, client AccountClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if len(cfg.Addresses) == 0 {
		return nil, errors.New("at least one address is required")
	}

	accounts := make([]AccountState, 0, len(cfg.Addresses))
	for i, addr := range cfg.Addresses {
		if addr == (common.Address{}) {
			return nil, fmt.Errorf("address %d is zero address", i)
		}

		bal, err := client.BalanceAt(ctx, addr, cfg.BlockNumber)
		if err != nil {
			return nil, fmt.Errorf("balance for %s: %w", addr.Hex(), err)
		}
		if bal == nil {
			bal = big.NewInt(0)
		}

		code, err := client.CodeAt(ctx, addr, cfg.BlockNumber)
		if err != nil {
			return nil, fmt.Errorf("code for %s: %w", addr.Hex(), err)
		}

		kind := AccountTypeEOA
		if len(code) > 0 {
			kind = AccountTypeContract
		}

		accounts = append(accounts, AccountState{
			Address: addr,
			Balance: new(big.Int).Set(bal),
			Code:    append([]byte(nil), code...),
			Type:    kind,
		})
	}

	return &Result{Accounts: accounts}, nil
}
