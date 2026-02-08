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
Reference Solution - abigen-Style Contract Binding
==================================================

This file demonstrates contract interaction using go-ethereum's abi package
and bind.BoundContract. Unlike eth_call with manual selector/decode, we use
abi.JSON to parse the ABI, then contract.Call which packs, calls, and unpacks
for us. This is how abigen-generated bindings work under the hood.

This connects to the Ethereum ecosystem by showing:
- abi.JSON: parse JSON ABI into go-ethereum ABI struct
- bind.NewBoundContract: wraps address + ABI + caller for method invocation
- bind.CallOpts: Context, BlockNumber for eth_call parameters
- contract.Call(opts, &out, method, args...): pack, call, unpack in one step

The exercise builds understanding of:
- BoundContract: no codegen — dynamic method dispatch by name
- Type assertions: out[0].(string), out[0].(*big.Int) — ABI returns []interface{}
- cfg.Holder: optional *common.Address; nil means skip balanceOf
- Defensive copy: new(big.Int).Set(v) for *big.Int returns

Teaching notes (per .cursorrules):
- Pointer semantics: cfg.Holder is *common.Address; *cfg.Holder dereferences when
  calling balanceOf. cfg.BlockNumber can be nil (latest).
- Memory: BoundContract.Call writes into &out; we don't share those interface{}
  values with caller without copying (e.g. *big.Int).
*/
func Run(ctx context.Context, backend ContractCaller, cfg Config) (*Result, error) {
	if backend == nil {
		return nil, errNilBackend
	}
	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address is required")
	}

	// Use provided ABI or default ERC20 metadata subset
	abiJSON := cfg.ABI
	if abiJSON == "" {
		abiJSON = erc20ABI
	}
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("parse abi: %w", err)
	}

	contract := bind.NewBoundContract(cfg.Contract, parsedABI, backend, nil, nil)
	opts := &bind.CallOpts{Context: ctx, BlockNumber: cfg.BlockNumber} // nil BlockNumber = latest

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
		balance, err := callUint256(contract, opts, "balanceOf", *cfg.Holder) // *cfg.Holder = deref optional addr
		if err != nil {
			return nil, fmt.Errorf("call balanceOf: %w", err)
		}
		res.Balance = balance
	}

	return res, nil
}

// callString invokes a view method returning string; type-asserts single output.
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

// callUint8 invokes a view method returning uint8; type-asserts single output.
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

// callUint256 invokes a view method returning uint256; returns defensive copy of *big.Int.
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
