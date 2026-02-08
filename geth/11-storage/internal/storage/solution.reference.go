//go:build reference

package storage

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var errNilClient = errors.New("nil storage client")

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

func Run(ctx context.Context, client StorageClient, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address is required")
	}
	if cfg.Slot == nil {
		return nil, errors.New("slot is required")
	}

	resolved := slotToHash(cfg.Slot)
	if len(cfg.MappingKey) > 0 {
		resolved = mappingSlotHash(cfg.MappingKey, resolved)
	}

	value, err := client.StorageAt(ctx, cfg.Contract, resolved, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("storage at slot %s: %w", resolved.Hex(), err)
	}

	return &Result{
		ResolvedSlot: resolved,
		Value:        append([]byte(nil), value...),
	}, nil
}

// slotToHash implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func slotToHash(slot *big.Int) common.Hash {
	if slot == nil {
		return common.Hash{}
	}
	return common.BigToHash(slot)
}

// mappingSlotHash implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func mappingSlotHash(key []byte, slot common.Hash) common.Hash {
	buf := make([]byte, 0, 64)
	buf = append(buf, common.LeftPadBytes(key, 32)...)
	buf = append(buf, slot.Bytes()...)
	return crypto.Keccak256Hash(buf)
}
