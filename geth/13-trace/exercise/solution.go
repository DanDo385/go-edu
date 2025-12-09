//go:build solution
// +build solution

package exercise

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

/*
Problem: Trace transaction execution to see opcode-level details and gas usage.

Transaction tracing replays a transaction in the EVM and returns structured data
describing every operation (call, gas usage, storage changes, etc.). This is
essential for:
  - Debugging contract behavior (why did this revert?)
  - Analyzing gas usage (which operations are expensive?)
  - Understanding internal calls (what contracts were called?)
  - Building block explorers and analytics tools

Computer science principles highlighted:
  - Deterministic replay (same inputs → same execution trace)
  - Execution instrumentation (observing without changing behavior)
  - JSON as a universal interchange format for complex data
*/
func Run(ctx context.Context, client TraceClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Defensive Programming Pattern
	// ============================================================================
	// Why validate inputs? Tracing is one of the most expensive RPC operations!
	// It requires replaying the entire transaction execution in the EVM, which
	// can take several seconds for complex transactions. Validating inputs early
	// prevents wasting computational resources.
	//
	// Context handling: Same pattern as all previous modules. If ctx is nil, we
	// provide context.Background() as a safe default. This ensures our RPC calls
	// always have a valid context for cancellation and timeout handling.
	//
	// This pattern is now familiar: validate → default → proceed.
	if ctx == nil {
		ctx = context.Background()
	}

	// Client validation: The TraceClient interface is our dependency. If it's nil,
	// we can't make RPC calls. Returning early with a descriptive error is Go's
	// idiomatic error handling pattern.
	//
	// Building on modules 01, 11, and 12: Same validation pattern, different
	// interface. Consistency across modules makes code predictable.
	if client == nil {
		return nil, errors.New("client is nil")
	}

	// Transaction hash validation: We need a valid transaction hash to trace.
	// The zero hash (0x0000...0000) is not a valid transaction hash—it's a
	// sentinel value indicating "no hash" or uninitialized data.
	//
	// Note: Unlike contract addresses (where zero address has meaning), a zero
	// transaction hash is always an error. Transactions are identified by their
	// keccak256 hash, which is never all zeros in practice.
	if cfg.TxHash == (common.Hash{}) {
		return nil, errors.New("tx hash required")
	}

	// ============================================================================
	// STEP 2: Trace Transaction - Understanding debug_traceTransaction
	// ============================================================================
	// This calls debug_traceTransaction under the hood, which is a powerful but
	// expensive RPC method. It replays the transaction and returns execution details.
	//
	// What's in a trace?
	//   - Call tree: Which contracts were called, in what order
	//   - Gas usage: How much gas each operation consumed
	//   - Storage changes: Which storage slots were modified
	//   - Return values: What data was returned from calls
	//   - Reverts: Where and why execution reverted
	//
	// Computer Science insight: This is **execution instrumentation**. We're
	// observing the execution without modifying it. The EVM is deterministic—
	// replaying a transaction always produces the same result and trace.
	//
	// Why this is powerful:
	//   - Debugging: See exactly where a transaction failed
	//   - Gas optimization: Identify expensive operations
	//   - Security audits: Understand call flow and state changes
	//   - Block explorers: Show users what happened in a transaction
	//
	// Important notes:
	//   - This is NOT part of standard eth_* RPC methods
	//   - Many public RPC providers disable it (too expensive)
	//   - You typically need your own node or a debug-enabled endpoint
	//   - Tracing old transactions requires archive node (needs historical state)
	//
	// Parameters:
	//   - ctx: Context for cancellation/timeout (important for long traces!)
	//   - cfg.TxHash: The transaction to trace
	//   - (implicit) tracer config: Default tracer returns full opcode-level trace
	//
	// Building on previous modules: We've read blocks (01), storage (11), proofs (12).
	// Now we're reading execution traces—a different view of the same blockchain data.
	raw, err := client.TraceTransaction(ctx, cfg.TxHash)
	if err != nil {
		return nil, fmt.Errorf("trace transaction: %w", err)
	}

	// Trace response validation: Even if the RPC call succeeds, the response
	// might be nil (though this shouldn't happen with well-behaved clients).
	// We validate to prevent nil pointer dereferences later.
	//
	// Building on previous patterns: Same validate-after-RPC pattern from modules
	// 01, 11, and 12. Always check both error and nil response.
	if raw == nil {
		return nil, errors.New("nil trace payload")
	}

	// ============================================================================
	// STEP 3: Defensive Copy of JSON Data
	// ============================================================================
	// json.RawMessage is defined as []byte, which is a slice (reference type).
	// Slices share their underlying array, so if we return `raw` directly,
	// callers could mutate the bytes and affect our internal state.
	//
	// Why this matters:
	//   1. The client might reuse internal buffers (performance optimization)
	//   2. If we return pointers to client's buffers, mutations cause data races
	//   3. Multiple callers might receive references to the same underlying data
	//   4. By copying, we ensure each caller gets independent, isolated data
	//
	// How to copy a slice:
	//   - make() allocates a new slice with exact capacity
	//   - copy() duplicates the bytes from source to destination
	//   - Result: traceCopy has its own backing array, independent of raw
	//
	// Building on module 01: We learned defensive copying for big.Int there.
	// Module 12 extended it to string slices. Here we apply it to byte slices.
	//
	// This pattern repeats: whenever returning data from external sources, always
	// consider whether defensive copying is needed. For reference types (slices,
	// maps, pointers to mutable structs), the answer is usually yes!
	traceCopy := make(json.RawMessage, len(raw))
	copy(traceCopy, raw)

	// ============================================================================
	// STEP 4: Return Result
	// ============================================================================
	// We return both the transaction hash and the trace data. Why both?
	//
	// TxHash: Confirms which transaction was traced. In batch processing scenarios,
	// this helps callers associate traces with transactions without manual bookkeeping.
	//
	// Trace: The raw JSON trace data. We return json.RawMessage (not parsed structs)
	// because trace format varies by tracer type and Geth version. By returning raw
	// JSON, callers can parse it according to their needs:
	//   - Default tracer: Full opcode-level trace
	//   - callTracer: Simplified call tree
	//   - prestateTracer: Account state before transaction
	//   - Custom tracers: Application-specific formats
	//
	// This flexibility is important! Different use cases need different trace details:
	//   - Debugging: Need full opcode trace
	//   - Gas analysis: Need gas-focused trace
	//   - Block explorer: Need high-level call tree
	//
	// Building on previous concepts:
	//   - We validated inputs (Step 1) → now we return validated data
	//   - We handled errors consistently (Steps 1-2) → now we return success
	//   - We used defensive copying (Step 3) → callers get isolated data
	//   - We used context for cancellation (Step 2) → operation completed successfully
	//
	// The progression across modules:
	//   - Module 01: Read block headers (lightweight metadata)
	//   - Module 11: Read storage slots (state data)
	//   - Module 12: Get proofs (cryptographic verification)
	//   - Module 13: Get traces (execution details)
	//   - Future: Combine all these to build rich analytics and debugging tools
	return &Result{
		TxHash: cfg.TxHash,
		Trace:  traceCopy,
	}, nil
}
