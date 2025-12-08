//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"
)

// Run contains the reference solution for module 04-accounts-balances.
func Run(ctx context.Context, client AccountClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation
	// ============================================================================
	// As always, we start with robust input validation. This is a recurring
	// theme and a cornerstone of reliable software.
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if len(cfg.Addresses) == 0 {
		return nil, errors.New("no addresses provided")
	}

	// ============================================================================
	// STEP 2: Initialize Result Slice
	// ============================================================================
	// We're going to be collecting state for multiple accounts. It's a Go
	// performance best practice to pre-allocate slices when you know the
	// eventual size. This avoids multiple re-allocations and copies as the
	// slice grows.
	accounts := make([]AccountState, 0, len(cfg.Addresses))

	// ============================================================================
	// STEP 3: Iterate and Query Accounts
	// ============================================================================
	// We loop through each address provided in the configuration and query its
	// state. This demonstrates a common batch-processing pattern.
	for _, addr := range cfg.Addresses {
		// ========================================================================
		// SUB-STEP 3a: Query Balance
		// ========================================================================
		// `BalanceAt` queries the ETH balance of an account at a specific block
		// number (`cfg.BlockNumber`), or the latest block if `nil`.
		bal, err := client.BalanceAt(ctx, addr, cfg.BlockNumber)
		if err != nil {
			// Error wrapping provides context for which address failed.
			return nil, fmt.Errorf("balance %s: %w", addr.Hex(), err)
		}
		// Defensive Copy #1: `big.Int` is a mutable pointer type. We copy it
		// to prevent the caller from accidentally modifying the RPC client's
		// internal data. This pattern was introduced in module 01.
		if bal != nil {
			bal = new(big.Int).Set(bal)
		}

		// ========================================================================
		// SUB-STEP 3b: Query Code
		// ========================================================================
		// `CodeAt` queries the contract bytecode at a given address. If the
		// account is an EOA, this will return an empty slice.
		code, err := client.CodeAt(ctx, addr, cfg.BlockNumber)
		if err != nil {
			return nil, fmt.Errorf("code %s: %w", addr.Hex(), err)
		}
		// Defensive Copy #2: Slices are reference types. We copy the code
		// to prevent the caller from modifying the underlying array in the
		// RPC client's memory. `append([]byte(nil), code...)` is the
		// idiomatic way to create a copy of a byte slice.
		codeCopy := append([]byte(nil), code...)

		// ========================================================================
		// SUB-STEP 3c: Classify Account Type
		// ========================================================================
		// This is the core logic of the module. The defining feature of a
		// contract is the presence of code.
		//   - If `len(code) > 0`, it's a contract.
		//   - If `len(code) == 0`, it's an EOA.
		//
		// (Note: This is a reliable heuristic. Edge cases like self-destructed
		// contracts can have `nonce > 0` but `code len == 0`, but for most
		// practical purposes, checking code length is sufficient.)
		accType := AccountTypeEOA
		if len(codeCopy) > 0 {
			accType = AccountTypeContract
		}

		// ========================================================================
		// SUB-STEP 3d: Append to Results
		// ========================================================================
		// We append the collected state for the current address to our results slice.
		accounts = append(accounts, AccountState{
			Address: addr,
			Balance: bal,
			Code:    codeCopy,
			Type:    accType,
		})
	}

	// ============================================================================
	// STEP 4: Return Final Result
	// ============================================================================
	// We wrap our slice of account states in the final Result struct.
	return &Result{Accounts: accounts}, nil
}
