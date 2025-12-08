//go:build !solution
// +build !solution

package exercise

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

// Run contains the reference solution for module 08-abigen.
func Run(ctx context.Context, backend ContractCaller, cfg Config) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if backend == nil {
		return nil, errors.New("backend is nil")
	}
	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address required")
	}

	abiJSON := cfg.ABI
	if strings.TrimSpace(abiJSON) == "" {
		abiJSON = erc20ABI
	}
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("parse ABI: %w", err)
	}

	contract := bind.NewBoundContract(cfg.Contract, parsedABI, backend, nil, nil)

	callOpts := &bind.CallOpts{
		Context:     ctx,
		BlockNumber: cfg.BlockNumber,
	}
	if cfg.Holder != nil {
		callOpts.From = *cfg.Holder
	}

	name, err := callString(contract, callOpts, "name")
	if err != nil {
		return nil, err
	}

	symbol, err := callString(contract, callOpts, "symbol")
	if err != nil {
		return nil, err
	}

	decimals, err := callUint8(contract, callOpts, "decimals")
	if err != nil {
		return nil, err
	}

	totalSupply, err := callUint256(contract, callOpts, "totalSupply")
	if err != nil {
		return nil, err
	}

	var balance *big.Int
	if cfg.Holder != nil {
		balance, err = callUint256(contract, callOpts, "balanceOf", *cfg.Holder)
		if err != nil {
			return nil, err
		}
	}

	return &Result{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
		Balance:     balance,
	}, nil
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
