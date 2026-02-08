//go:build reference

package proofs

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var errNilClient = errors.New("nil proof client")

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

func Run(ctx context.Context, client ProofClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if cfg.Account == (common.Address{}) {
		return nil, errors.New("account is required")
	}

	slots := make([]string, 0, len(cfg.Slots))
	for _, slot := range cfg.Slots {
		slots = append(slots, slot.Hex())
	}

	resp, err := client.GetProof(ctx, cfg.Account, slots, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("get proof: %w", err)
	}
	if resp == nil {
		return nil, errors.New("get proof: nil response")
	}

	balance := big.NewInt(0)
	if resp.Balance != nil {
		balance = new(big.Int).Set(resp.Balance)
	}

	storage := make([]StorageProof, 0, len(resp.StorageProof))
	for _, sp := range resp.StorageProof {
		v := big.NewInt(0)
		if sp.Value != nil {
			v = new(big.Int).Set(sp.Value)
		}
		storage = append(storage, StorageProof{
			Key:        common.HexToHash(sp.Key),
			Value:      v,
			ProofNodes: append([]string(nil), sp.Proof...),
		})
	}

	return &Result{
		Account: AccountProof{
			Balance:     balance,
			Nonce:       resp.Nonce,
			CodeHash:    resp.CodeHash,
			StorageHash: resp.StorageHash,
			ProofNodes:  append([]string(nil), resp.AccountProof...),
			Storage:     storage,
		},
	}, nil
}
