//go:build reference

package smartcontracts

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	errNilClient     = errors.New("nil contract caller")
	errShortABIValue = errors.New("abi value too short")
)

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

func Run(ctx context.Context, client ContractCaller, cfg Config) (*Result, error) {
	if client == nil {
		return nil, errNilClient
	}
	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address is required")
	}

	call := func(sel []byte) ([]byte, error) {
		payload, err := client.CallContract(ctx, ethereum.CallMsg{To: &cfg.Contract, Data: sel}, cfg.BlockNumber)
		if err != nil {
			return nil, err
		}
		return payload, nil
	}

	nameRaw, err := call(selector("name()"))
	if err != nil {
		return nil, fmt.Errorf("call name: %w", err)
	}
	name, err := decodeString(nameRaw)
	if err != nil {
		return nil, fmt.Errorf("decode name: %w", err)
	}

	symbolRaw, err := call(selector("symbol()"))
	if err != nil {
		return nil, fmt.Errorf("call symbol: %w", err)
	}
	symbol, err := decodeString(symbolRaw)
	if err != nil {
		return nil, fmt.Errorf("decode symbol: %w", err)
	}

	decimalsRaw, err := call(selector("decimals()"))
	if err != nil {
		return nil, fmt.Errorf("call decimals: %w", err)
	}
	decimals, err := decodeUint8(decimalsRaw)
	if err != nil {
		return nil, fmt.Errorf("decode decimals: %w", err)
	}

	supplyRaw, err := call(selector("totalSupply()"))
	if err != nil {
		return nil, fmt.Errorf("call totalSupply: %w", err)
	}
	totalSupply, err := decodeUint256(supplyRaw)
	if err != nil {
		return nil, fmt.Errorf("decode totalSupply: %w", err)
	}

	return &Result{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
	}, nil
}

// selector implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func selector(sig string) []byte {
	h := crypto.Keccak256([]byte(sig))
	return append([]byte(nil), h[:4]...)
}

// decodeString implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func decodeString(data []byte) (string, error) {
	if len(data) < 64 {
		return "", errShortABIValue
	}
	offset := new(big.Int).SetBytes(data[:32])
	if !offset.IsInt64() || offset.Int64() < 0 {
		return "", errors.New("invalid string offset")
	}
	off := int(offset.Int64())
	if off+32 > len(data) {
		return "", errors.New("string offset out of bounds")
	}

	length := new(big.Int).SetBytes(data[off : off+32])
	if !length.IsInt64() || length.Int64() < 0 {
		return "", errors.New("invalid string length")
	}
	n := int(length.Int64())
	start := off + 32
	end := start + n
	if end > len(data) {
		return "", errors.New("string length out of bounds")
	}

	return string(data[start:end]), nil
}

// decodeUint8 implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func decodeUint8(data []byte) (uint8, error) {
	v, err := decodeUint256(data)
	if err != nil {
		return 0, err
	}
	if !v.IsUint64() || v.Uint64() > 255 {
		return 0, errors.New("uint8 out of range")
	}
	return uint8(v.Uint64()), nil
}

// decodeUint256 implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func decodeUint256(data []byte) (*big.Int, error) {
	if len(data) < 32 {
		return nil, errShortABIValue
	}
	return new(big.Int).SetBytes(data[:32]), nil
}
