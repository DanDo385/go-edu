//go:build !solution && !reference

package stack

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

/*
Problem: Prove RPC connectivity by reading the network identifiers and latest header.

The very first thing an Ethereum Go tool should do is dial an RPC endpoint,
retrieve the chain/network IDs (replay protection + legacy identifier), and
inspect a block header. Headers are lightweight (~500 bytes) yet contain the
state root, parent hash, and other cryptographic commitments that define the
execution stack you are about to interact with. This function mirrors the CLI
demo from module 01 but exposes it as a reusable library API.

Computer science principles highlighted:
  - Separation of configuration from code (cfg.BlockNumber allows deterministic tests)
  - Fault tolerance via context propagation—callers control cancellation/timeouts
  - Immutability via defensive copies (we never hand pointers owned by go-ethereum back to callers)
*/
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Defensive Programming Pattern
	// TODO: Implement

	// ============================================================================
	// Why validate inputs? This function is a library API that will be called by
	// TODO: Implement

	// ============================================================================
	// STEP 2: Retrieve Chain ID - Understanding Replay Protection
	// TODO: Implement

	// ============================================================================
	// Chain ID is fundamental to Ethereum's security model. Introduced in EIP-155,
	// TODO: Implement

	// ============================================================================
	// STEP 3: Retrieve Network ID - Legacy Identifier Pattern
	// TODO: Implement

	// ============================================================================
	// Network ID predates Chain ID and was used for P2P networking (identifying
	// TODO: Implement

	// ============================================================================
	// STEP 4: Retrieve Block Header - Understanding Block Structure
	// TODO: Implement

	// ============================================================================
	// Block headers are the "lightweight" representation of blocks. They contain
	// TODO: Implement

	// ============================================================================
	// STEP 5: Construct Result with Defensive Copying - Immutability Pattern
	// TODO: Implement

	// ============================================================================
	// CRITICAL CONCEPT: Defensive copying prevents data races and unexpected mutations.
	// TODO: Implement

	panic("unimplemented")
}
