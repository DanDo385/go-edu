//go:build !solution && !reference

package explorer

import (
	"context"
	"errors"
	"fmt"
)

/*
Run contains the reference implementation for the tiny block explorer.

This demonstrates how block explorers work: fetch a block, extract metadata,
and optionally enumerate transactions. Block explorers are just structured
views of blockchain data.

Computer science principles highlighted:
  - Data aggregation (combining block and transaction data)
  - Optional expansion (cfg.IncludeTxs controls detail level)
  - Minimal data transfer (only fetch what's needed)
*/
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Defensive Programming Pattern
	// TODO: Implement

	// ============================================================================
	// Same validation pattern as all previous modules. Validate inputs before
	// TODO: Implement

	// ============================================================================
	// STEP 2: Fetch Block - Understanding Block vs Header
	// TODO: Implement

	// ============================================================================
	// BlockByNumber fetches the full block including transactions. This is
	// TODO: Implement

	// ============================================================================
	// STEP 3: Extract Block Metadata - Data Aggregation Pattern
	// TODO: Implement

	// ============================================================================
	// We extract key fields from the block header to build our explorer view.
	// TODO: Implement

	// ============================================================================
	// STEP 4: Optionally Include Transaction Summaries - Controlled Expansion
	// TODO: Implement

	// ============================================================================
	// If cfg.IncludeTxs is true, we enumerate transactions and build summaries.
	// TODO: Implement

	// ============================================================================
	// STEP 5: Return Result - Minimal Defensive Copying
	// TODO: Implement

	// ============================================================================
	// types.Block shares slices internally, but we only copy primitive values
	// TODO: Implement

	panic("unimplemented")
}
