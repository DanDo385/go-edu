//go:build reference

package abigen

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

const erc20ABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

var errNilBackend = errors.New("nil contract backend")

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

func Run(ctx context.Context, backend ContractCaller, cfg Config) (*Result, error) {
	if backend == nil {
		return nil, errNilBackend
	}
	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address is required")
	}

	abiJSON := cfg.ABI
	if abiJSON == "" {
		abiJSON = erc20ABI
	}
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("parse abi: %w", err)
	}

	contract := bind.NewBoundContract(cfg.Contract, parsedABI, backend, nil, nil)
	opts := &bind.CallOpts{Context: ctx, BlockNumber: cfg.BlockNumber}

	name, err := callString(contract, opts, "name")
	if err != nil {
		return nil, fmt.Errorf("call name: %w", err)
	}
	symbol, err := callString(contract, opts, "symbol")
	if err != nil {
		return nil, fmt.Errorf("call symbol: %w", err)
	}
	decimals, err := callUint8(contract, opts, "decimals")
	if err != nil {
		return nil, fmt.Errorf("call decimals: %w", err)
	}
	totalSupply, err := callUint256(contract, opts, "totalSupply")
	if err != nil {
		return nil, fmt.Errorf("call totalSupply: %w", err)
	}

	res := &Result{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
	}

	if cfg.Holder != nil {
		balance, err := callUint256(contract, opts, "balanceOf", *cfg.Holder)
		if err != nil {
			return nil, fmt.Errorf("call balanceOf: %w", err)
		}
		res.Balance = balance
	}

	return res, nil
}

// callString implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func callString(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (string, error) {
	var out []interface{}
	if err := contract.Call(opts, &out, method, params...); err != nil {
		return "", err
	}
	if len(out) != 1 {
		return "", fmt.Errorf("method %s returned %d values", method, len(out))
	}
	s, ok := out[0].(string)
	if !ok {
		return "", fmt.Errorf("method %s returned %T, want string", method, out[0])
	}
	return s, nil
}

// callUint8 implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func callUint8(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (uint8, error) {
	var out []interface{}
	if err := contract.Call(opts, &out, method, params...); err != nil {
		return 0, err
	}
	if len(out) != 1 {
		return 0, fmt.Errorf("method %s returned %d values", method, len(out))
	}
	v, ok := out[0].(uint8)
	if !ok {
		return 0, fmt.Errorf("method %s returned %T, want uint8", method, out[0])
	}
	return v, nil
}

// callUint256 implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func callUint256(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (*big.Int, error) {
	var out []interface{}
	if err := contract.Call(opts, &out, method, params...); err != nil {
		return nil, err
	}
	if len(out) != 1 {
		return nil, fmt.Errorf("method %s returned %d values", method, len(out))
	}
	v, ok := out[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("method %s returned %T, want *big.Int", method, out[0])
	}
	return new(big.Int).Set(v), nil
}
