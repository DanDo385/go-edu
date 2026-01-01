//go:build !solution && !reference

package abigen

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
)

/*
Problem: Use BoundContract for type-safe contract calls with automatic ABI encoding/decoding.
*/

const erc20ABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

// Run - TODO: implement this function
func Run(ctx context.Context, backend ContractCaller, cfg Config) (*Result, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	var zero0 *Result
	var zero1 error
	return zero0, zero1
}

// callString - TODO: implement this function
func callString(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (string, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	var zero0 string
	var zero1 error
	return zero0, zero1
}

// callUint8 - TODO: implement this function
func callUint8(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (uint8, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	var zero0 uint8
	var zero1 error
	return zero0, zero1
}

// callUint256 - TODO: implement this function
func callUint256(contract *bind.BoundContract, opts *bind.CallOpts, method string, params ...interface{}) (*big.Int, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	var zero0 *big.Int
	var zero1 error
	return zero0, zero1
}
