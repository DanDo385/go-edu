//go:build !solution
// +build !solution

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

/*
Problem: Use BoundContract for type-safe contract calls with automatic ABI encoding/decoding.

This module teaches you how to use go-ethereum's BoundContract pattern for cleaner,
safer contract interactions. Instead of manually encoding/decoding like module 07,
you'll use the abi package to handle encoding/decoding automatically.

Computer science principles highlighted:
  - Adapter pattern: BoundContract wraps low-level RPC with high-level interface
  - Type safety: ABI definitions provide compile-time checks
  - Code reuse: Helper functions eliminate boilerplate
  - Separation of concerns: ABI encoding is separate from business logic
*/
func Run(ctx context.Context, backend ContractCaller, cfg Config) (*Result, error) {
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


func callString(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (string, error) {
	var out []interface{}
	if err := contract.Call(opts, &out, method, params...); err != nil {
		return "", fmt.Errorf("call %s: %w", method, err)
	}
	if len(out) == 0 {
		return "", fmt.Errorf("call %s: empty result", method)
	}
	return *abi.ConvertType(out[0], new(string)).(*string), nil
}

func callUint8(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (uint8, error) {
	var out []interface{}
	if err := contract.Call(opts, &out, method, params...); err != nil {
		return 0, fmt.Errorf("call %s: %w", method, err)
	}
	if len(out) == 0 {
		return 0, fmt.Errorf("call %s: empty result", method)
	}
	return *abi.ConvertType(out[0], new(uint8)).(*uint8), nil
}

func callUint256(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (*big.Int, error) {
	var out []interface{}
	if err := contract.Call(opts, &out, method, params...); err != nil {
		return nil, fmt.Errorf("call %s: %w", method, err)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("call %s: empty result", method)
	}
	return *abi.ConvertType(out[0], new(*big.Int)).(**big.Int), nil
}
