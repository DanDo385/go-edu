//go:build reference

package stack

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

var errNilClient = errors.New("nil rpc client")

/*
Reference Solution - Ethereum Stack Connectivity
==============================================

This file demonstrates the fundamental connection to Ethereum networks via RPC.
Every Ethereum application starts by verifying network connectivity and retrieving
basic blockchain metadata. This establishes the foundation for all subsequent operations.

This connects to the broader Ethereum ecosystem by showing:
- JSON-RPC communication patterns used by all Ethereum clients
- Chain ID validation for replay protection (EIP-155)
- Network ID identification for peer-to-peer networking
- Block header inspection for state root verification
- Defensive copying to prevent data races with external libraries

The exercise builds understanding of:
- RPC client interfaces and dependency injection for testing
- Context propagation for cancellation and timeouts
- Big integer arithmetic for Ethereum's 256-bit numbers
- Error wrapping patterns for debugging complex call stacks
- Memory safety when interfacing with CGO-based libraries

Teaching notes:
- Memory/ownership: big.Int and types.Header are mutable pointer types from go-ethereum.
  We make defensive copies to prevent callers from accidentally modifying library internals.
- Invariants: RPC calls can fail for network, authentication, or data reasons. We validate
  all responses to establish data integrity before proceeding.
- Error surfaces: network operations are inherently unreliable. We surface all errors
  explicitly so callers can implement appropriate retry, fallback, or alerting logic.
*/

/*
Run - Ethereum Stack Connectivity Check

This function performs the essential "hello world" operation for any Ethereum application:
verify network connectivity and retrieve basic blockchain state.

Parameters:
- ctx: context for cancellation, timeouts, and request tracing
- client: RPC client interface (allows testing with mocks)
- cfg: configuration for which block to inspect (nil = latest)

Returns:
- *Result: snapshot of blockchain state at the requested block
- error: any connectivity or data retrieval failures

Algorithm steps:
1. Validate client interface (defensive programming)
2. Retrieve chain ID for replay protection validation
3. Retrieve network ID for peer networking identification
4. Retrieve block header for state root and metadata inspection
5. Return defensive copies to prevent data races

Why this pattern matters:
- Establishes baseline connectivity before complex operations
- Validates we're talking to the expected network
- Provides cryptographic commitments (state root) for data integrity
- Demonstrates proper resource cleanup and error handling
*/
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	// Step 1: Input validation
	// Interface values can be nil - check before calling methods
	// This prevents runtime panics and provides clear error messages
	if client == nil {
		return nil, errNilClient
	}

	// Step 2: Retrieve Chain ID
	// Chain ID prevents transaction replay attacks across networks (EIP-155)
	// Different networks (mainnet=1, testnets=other values) have different chain IDs
	// This is the first security check every Ethereum app performs
	chainID, err := client.ChainID(ctx)
	if err != nil {
		// Wrap the error with context about which operation failed
		// %w verb enables error chaining for debugging (Go 1.13+ feature)
		return nil, fmt.Errorf("chain id: %w", err)
	}
	// Validate response - RPC can return nil even without error
	if chainID == nil {
		return nil, errors.New("chain id: nil response")
	}

	// Step 3: Retrieve Network ID
	// Network ID is legacy identifier for P2P peer discovery
	// While Chain ID is for security, Network ID is for networking compatibility
	// Some older tools still rely on Network ID for network identification
	networkID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("network id: %w", err)
	}
	if networkID == nil {
		return nil, errors.New("network id: nil response")
	}

	// Step 4: Retrieve Block Header
	// Block headers contain cryptographic commitments to blockchain state
	// Key fields: state root (Merkle root of all accounts), transactions root,
	// receipts root, parent hash (chain integrity), timestamp, gas limits
	// cfg.BlockNumber is nil for "latest block", or specific block number
	header, err := client.HeaderByNumber(ctx, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("header by number: %w", err)
	}
	if header == nil {
		return nil, errors.New("header by number: nil response")
	}

	// Step 5: Return defensive copies
	// CRITICAL: Prevent data races by copying mutable external data
	// big.Int values from RPC client are mutable - if we return them directly,
	// callers could modify them, affecting other parts of the program
	// types.Header contains internal pointers that could be shared
	return &Result{
		// new(big.Int).Set(chainID) creates independent copy of the big integer
		// This prevents caller mutations from affecting cached client data
		ChainID: new(big.Int).Set(chainID),

		// Same defensive copy pattern for network ID
		NetworkID: new(big.Int).Set(networkID),

		// types.CopyHeader() performs deep copy of header struct and all pointers
		// This ensures complete isolation from the RPC client's internal data
		Header: types.CopyHeader(header),
	}, nil
}
