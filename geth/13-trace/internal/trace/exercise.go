//go:build !solution && !reference

package trace

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
	// TODO: Implement

	// ============================================================================
	// Why validate inputs? Tracing is one of the most expensive RPC operations!
	// TODO: Implement

	// ============================================================================
	// STEP 2: Trace Transaction - Understanding debug_traceTransaction
	// TODO: Implement

	// ============================================================================
	// This calls debug_traceTransaction under the hood, which is a powerful but
	// TODO: Implement

	// ============================================================================
	// STEP 3: Defensive Copy of JSON Data
	// TODO: Implement

	// ============================================================================
	// json.RawMessage is defined as []byte, which is a slice (reference type).
	// TODO: Implement

	// ============================================================================
	// STEP 4: Return Result
	// TODO: Implement

	// ============================================================================
	// We return both the transaction hash and the trace data. Why both?
	// TODO: Implement

	panic("unimplemented")
}
