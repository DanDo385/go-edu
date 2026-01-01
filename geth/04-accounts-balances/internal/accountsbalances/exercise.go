//go:build !solution && !reference

package accountsbalances

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
	// TODO: Implement

	// ============================================================================
	// As always, we start with robust input validation. This is a recurring
	// TODO: Implement

	// ============================================================================
	// STEP 2: Initialize Result Slice
	// TODO: Implement

	// ============================================================================
	// We're going to be collecting state for multiple accounts. It's a Go
	// TODO: Implement

	// ============================================================================
	// STEP 3: Iterate and Query Accounts
	// TODO: Implement

	// ============================================================================
	// We loop through each address provided in the configuration and query its
	// TODO: Implement

	// ========================================================================
	// SUB-STEP 3a: Query Balance
	// TODO: Implement

	// ========================================================================
	// `BalanceAt` queries the ETH balance of an account at a specific block
	// TODO: Implement

	// ========================================================================
	// SUB-STEP 3b: Query Code
	// TODO: Implement

	// ========================================================================
	// `CodeAt` queries the contract bytecode at a given address. If the
	// TODO: Implement

	// ========================================================================
	// SUB-STEP 3c: Classify Account Type
	// TODO: Implement

	// ========================================================================
	// This is the core logic of the module. The defining feature of a
	// TODO: Implement

	// ========================================================================
	// SUB-STEP 3d: Append to Results
	// TODO: Implement

	// ========================================================================
	// We append the collected state for the current address to our results slice.
	// TODO: Implement

	// ============================================================================
	// STEP 4: Return Final Result
	// TODO: Implement

	// ============================================================================
	// We wrap our slice of account states in the final Result struct.
	// TODO: Implement

	panic("unimplemented")
}
