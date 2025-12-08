//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

const retryDelay = 100 * time.Millisecond

// Run contains the reference solution for module 02-rpc-basics.
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Defensive Programming Pattern
	// ============================================================================
	// Building on module 01, we repeat the same robust input validation.
	// This pattern is universal: always validate inputs for library functions.
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	// New in this module: validating configuration specific to this function.
	if cfg.Retries < 0 {
		cfg.Retries = 0
	}

	// ============================================================================
	// STEP 2: Retrieve Latest Block Number
	// ============================================================================
	// This is a lightweight RPC call (`eth_blockNumber`) to get the current chain
	// height. It's a cheap way to verify basic connectivity.
	// The pattern is the same as in module 01: call, check error, wrap error.
	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("block number: %w", err)
	}

	// ============================================================================
	// STEP 3: Retrieve Network ID - Legacy Identifier Pattern
	// ============================================================================
	// We retrieve the network ID, just as in module 01. This reinforces the
	// pattern of querying basic network identifiers.
	networkID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("network id: %w", err)
	}

	// ============================================================================
	// STEP 4: Retrieve Full Block with Retry Logic - Fault Tolerance Pattern
	// ============================================================================
	// This is the core of this module. Fetching a full block is a heavier
	// operation and more likely to fail due to transient network issues.
	// We introduce a retry loop to build resilience.
	//
	// Computer science principle: Fault Tolerance. We are designing our system
	// to continue operating, albeit with a short delay, in the presence of faults.
	var block *types.Block
	var lastErr error
	for attempt := 0; attempt <= cfg.Retries; attempt++ {
		// We attempt to fetch the block. `nil` for the block number argument
		// means "latest block".
		block, lastErr = client.BlockByNumber(ctx, nil)
		if lastErr == nil {
			// Success! Break the loop.
			break
		}

		// If we failed, we wait before retrying. But we must also respect
		// the context. `select` allows us to wait for one of two events:
		//   1. The context is canceled (`ctx.Done()`).
		//   2. The retry delay has passed (`time.After(retryDelay)`).
		// This is a critical pattern for any long-running or blocking
		// operation in Go. It makes our function responsive to cancellation.
		select {
		case <-ctx.Done():
			// Context was canceled, so we give up immediately.
			// We wrap the context's error to provide more information.
			return nil, fmt.Errorf("context canceled while fetching block: %w", ctx.Err())
		case <-time.After(retryDelay):
			// The delay has passed, so we'll loop for another attempt.
		}
	}
	// After the loop, we check if we were ultimately successful. If `lastErr`
	// is not nil, it means all our attempts failed.
	if lastErr != nil {
		return nil, fmt.Errorf("block fetch failed after %d attempts: %w", cfg.Retries+1, lastErr)
	}

	// ============================================================================
	// STEP 5: Construct Result with Defensive Copying - Immutability Pattern
	// ============================================================================
	// We are again using defensive copying to ensure the caller cannot mutate
	// the RPC client's internal data.
	//
	// `types.CopyBlock` is the equivalent of `types.CopyHeader` from module 01.
	// It performs a deep copy of the block, including its header and transactions.
	//
	// Note that `blockNumber` is a `uint64`, which is a value type in Go.
	// Value types are copied by default, so we don't need to do anything special.
	// `networkID` is a `*big.Int`, which is a pointer to a mutable struct, so we
	// must copy it.
	return &Result{
		NetworkID:   new(big.Int).Set(networkID),
		BlockNumber: blockNumber,
		Block:       types.CopyBlock(block),
	}, nil
}
