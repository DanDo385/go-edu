//go:build reference

package accountsbalances

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var errNilClient = errors.New("nil account client")

/*
Reference Solution - Ethereum Account State Inspection
=====================================================

This file demonstrates querying Ethereum account states via JSON-RPC.
Every Ethereum address ( Externally Owned Account or Contract) has associated state:
balance (wei), code (for contracts), nonce, and storage. This forms the "state trie".

This connects to the broader Ethereum ecosystem by showing:
- How wallets query balances before transactions
- How explorers display account information
- How DeFi protocols check token balances
- How RPC clients abstract blockchain queries
- How block numbers enable historical state queries

The exercise builds understanding of:
- Account types: EOAs (users) vs contracts (programs)
- State inspection: balances, code, nonces at specific blocks
- Address validation: checksummed hex vs raw bytes
- RPC error handling: network failures, invalid addresses
- Data copying: preventing mutations of RPC client internals

Teaching notes:
- Memory/ownership: RPC client returns pointers to internal data structures that
  may be reused. We create defensive copies to prevent data races and ensure
  our Result remains valid even if client data changes.
- Invariants: zero address (0x000...000) is invalid for meaningful operations.
  Contract detection via non-empty code is reliable but may miss proxy contracts.
- Error surfaces: RPC calls can fail due to network issues, invalid addresses,
  or node synchronization problems. We surface these explicitly for caller handling.
*/

/*
Run - Query Ethereum Account States

This function inspects multiple Ethereum addresses at a specific block height,
returning their balances, code, and account types. This is the foundation for
wallet balances, contract verification, and DeFi position tracking.

Parameters:
- ctx: context for cancellation and timeout control
- client: RPC client interface for blockchain queries
- cfg: configuration with addresses and block number

Returns:
- *Result: snapshot of account states at the requested block
- error: any RPC failures or invalid inputs

Algorithm steps:
1. Validate inputs (client, addresses)
2. For each address: query balance and code
3. Determine account type (EOA vs contract)
4. Return defensive copies of all data

Why account inspection matters:
- Wallets need balances to show user funds
- DApps verify contract deployments
- Explorers display account information
- Security tools audit suspicious addresses
*/
func Run(ctx context.Context, client AccountClient, cfg Config) (*Result, error) {
	// INPUT VALIDATION
	// Interface nil check prevents runtime panics
	if client == nil {
		return nil, errNilClient
	}

	// Require at least one address to inspect
	// Empty queries waste RPC calls and provide no value
	if len(cfg.Addresses) == 0 {
		return nil, errors.New("at least one address is required")
	}

	// RESULT ACCUMULATION
	// Pre-allocate slice with exact capacity to avoid reallocations
	// This is more efficient when we know the final size upfront
	accounts := make([]AccountState, 0, len(cfg.Addresses))

	// ADDRESS INSPECTION LOOP
	// Process each address independently - failures don't stop others
	// This enables partial success for batch queries
	for i, addr := range cfg.Addresses {
		// VALIDATE ADDRESS
		// Zero address (0x000...000) is invalid for meaningful operations
		// It's the default value and represents "no address"
		// common.Address{} creates a zero address for comparison
		if addr == (common.Address{}) {
			return nil, fmt.Errorf("address %d is zero address", i)
		}

		// QUERY ACCOUNT BALANCE
		// BalanceAt retrieves wei balance at specified block
		// cfg.BlockNumber can be nil (latest) or specific block hash/number
		// nil block means "current best block"
		bal, err := client.BalanceAt(ctx, addr, cfg.BlockNumber)
		if err != nil {
			// RPC error - network issue, invalid address, or node problem
			// Include address in error for debugging
			return nil, fmt.Errorf("balance for %s: %w", addr.Hex(), err)
		}

		// HANDLE NIL BALANCE
		// Some RPC implementations return nil for zero balance
		// Normalize to explicit zero for consistent handling
		if bal == nil {
			bal = big.NewInt(0)
		}

		// QUERY ACCOUNT CODE
		// CodeAt retrieves deployed bytecode for contracts
		// EOAs (users) have empty code, contracts have bytecode
		// This is how we distinguish account types
		code, err := client.CodeAt(ctx, addr, cfg.BlockNumber)
		if err != nil {
			// RPC error similar to balance query
			return nil, fmt.Errorf("code for %s: %w", addr.Hex(), err)
		}

		// DETERMINE ACCOUNT TYPE
		// Ethereum has two account types:
		// 1. EOAs (Externally Owned Accounts): controlled by private keys
		// 2. Contracts: controlled by bytecode, no private key
		kind := AccountTypeEOA      // Assume EOA by default
		if len(code) > 0 {          // Non-empty code means contract
			kind = AccountTypeContract
		}

		// BUILD ACCOUNT STATE
		// Create AccountState with defensive copies
		// RPC client data may be reused - we need independent copies
		accounts = append(accounts, AccountState{
			Address: addr,                             // Value type, copy safe
			Balance: new(big.Int).Set(bal),           // Deep copy big integer
			Code:    append([]byte(nil), code...),    // Copy byte slice
			Type:    kind,                             // Value type, copy safe
		})
	}

	// RETURN RESULTS
	// Wrap accounts slice in Result struct
	// This provides consistent API and extension points
	return &Result{Accounts: accounts}, nil
}
