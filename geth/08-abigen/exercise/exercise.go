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
	// TODO: Validate input parameters
	// - Check if ctx is nil and provide a default context if needed
	// - Check if backend is nil and return an appropriate error
	// - Check if cfg.Contract is the zero address and return an error
	// - Backend is the ContractCaller interface (RPC client)

	// TODO: Parse ABI JSON string
	// - Use cfg.ABI if provided, otherwise use erc20ABI constant (defined at top)
	// - Trim whitespace from ABI string before checking if empty
	// - Parse ABI using abi.JSON(strings.NewReader(abiJSON))
	// - Handle parsing errors (invalid JSON, malformed ABI)
	// - Why parse ABI? We need the function signatures and return types for encoding/decoding
	// - ABI is like an interface definition: it tells us what functions exist and their types

	// TODO: Create BoundContract
	// - Use bind.NewBoundContract(address, parsedABI, backend, nil, nil)
	// - Parameters explained:
	//   * address: Contract address to call
	//   * parsedABI: Parsed ABI object (from previous step)
	//   * backend: RPC client for read operations (our ContractCaller)
	//   * transactor: RPC client for write operations (nil - we're only reading)
	//   * filterer: RPC client for event filtering (nil - not needed for this module)
	// - BoundContract is the adapter pattern: it wraps low-level RPC with high-level interface
	// - This is more convenient than manual ABI encoding (module 07)

	// TODO: Create CallOpts for contract calls
	// - Build bind.CallOpts struct with:
	//   * Context: ctx (for cancellation/timeout)
	//   * BlockNumber: cfg.BlockNumber (nil for latest, or specific block)
	//   * From: cfg.Holder if provided (optional sender address for view functions)
	// - Why From? Some view functions check msg.sender (e.g., allowance checks)
	// - CallOpts is like a request context: it configures how the call should execute

	// TODO: Call and decode name() function
	// - Use callString helper (already implemented below) with:
	//   * contract (BoundContract from above)
	//   * callOpts (CallOpts from above)
	//   * method name: "name"
	// - Handle errors (network failures, contract reverts, decoding errors)
	// - callString helper wraps contract.Call and does type conversion
	// - Compare to module 07: No manual selector computation or decoding!

	// TODO: Call and decode symbol() function
	// - Follow same pattern as name()
	// - Use callString helper with method name "symbol"
	// - Handle errors
	// - Notice the pattern: same helper, just different method name
	// - This demonstrates code reuse: one helper for all string-returning functions

	// TODO: Call and decode decimals() function
	// - Follow same pattern but use callUint8 helper
	// - Method name: "decimals"
	// - Handle errors
	// - callUint8 is like callString but for uint8 return type
	// - Pattern repeats: helper function → method name → error handling

	// TODO: Call and decode totalSupply() function
	// - Follow same pattern but use callUint256 helper
	// - Method name: "totalSupply"
	// - Handle errors
	// - callUint256 returns *big.Int (for large numbers)
	// - Again, same pattern with different helper and method name

	// TODO: Optionally call balanceOf(address) if holder provided
	// - Check if cfg.Holder is not nil
	// - If provided, call callUint256 with:
	//   * method name: "balanceOf"
	//   * parameter: *cfg.Holder (dereference pointer)
	// - Handle errors
	// - balanceOf takes a parameter (address), unlike the other functions
	// - The helper supports variadic parameters for this: callUint256(contract, opts, "balanceOf", address)

	// TODO: Construct and return Result
	// - Create Result struct with:
	//   * Name: name string
	//   * Symbol: symbol string
	//   * Decimals: decimals uint8
	//   * TotalSupply: totalSupply *big.Int
	//   * Balance: balance *big.Int (or nil if holder not provided)
	// - Return result with nil error
	// - Why no defensive copying? These values are already copies from type conversion

	return nil, errors.New("not implemented")
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
