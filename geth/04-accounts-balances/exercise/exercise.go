//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"errors"
)

/*
Problem: Query the state of multiple Ethereum accounts and classify them as
either Externally Owned Accounts (EOAs) or Smart Contracts.

This module dives into the fundamental difference between the two account types
on Ethereum. You'll learn how to query the balance and code of an account, and
use that information to determine its type. This is a core skill for building
any application that interacts with Ethereum accounts.

Computer science principles highlighted:
  - Type systems (classifying accounts based on their properties)
  - State machines (querying the state of an account at a specific block)
  - Defensive programming (handling nil inputs, copying reference types)
*/
func Run(ctx context.Context, client AccountClient, cfg Config) (*Result, error) {
	// TODO: Validate input parameters
	// - Check for nil context and provide a default
	// - Check for nil client and return an error
	// - Check if `cfg.Addresses` is empty and return an error if it is

	// TODO: Initialize a slice to store the account states
	// - Use `make([]AccountState, 0, len(cfg.Addresses))` to create a slice
	//   with a pre-allocated capacity.

	// TODO: Iterate over the addresses in `cfg.Addresses`
	// - For each address, you will query its balance and code.

	// TODO: Query the account balance
	// - Use `client.BalanceAt(ctx, addr, cfg.BlockNumber)`
	// - Handle and wrap any errors. Include the address in the error message
	//   for better debugging.
	// - IMPORTANT: `BalanceAt` returns a `*big.Int`, which is a pointer to a
	//   mutable struct. You must create a defensive copy of the balance to
	//   avoid accidentally mutating the RPC client's internal data.
	//   Use `new(big.Int).Set(bal)`.

	// TODO: Query the account code
	// - Use `client.CodeAt(ctx, addr, cfg.BlockNumber)`
	// - Handle and wrap any errors.
	// - IMPORTANT: `CodeAt` returns a `[]byte`, which is a slice. Slices are
	//   reference types. You must create a defensive copy of the code.
	//   The idiomatic way to copy a slice is `append([]byte(nil), code...)`.

	// TODO: Classify the account type
	// - If the length of the code slice is greater than 0, the account is a
	//   contract (`AccountTypeContract`).
	// - Otherwise, it's an EOA (`AccountTypeEOA`).

	// TODO: Append the account state to the slice
	// - Create an `AccountState` struct with the address, balance, code, and
	//   account type.
	// - Append it to the slice of account states.

	// TODO: Construct and return the Result struct
	// - The result should contain the slice of account states.

	return nil, errors.New("not implemented")
}