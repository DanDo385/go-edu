//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"encoding/json"
	"errors"

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
	// TODO: Validate input parameters
	// - Check if ctx is nil and provide a default context if needed
	// - Check if client is nil and return an appropriate error
	// - Validate that cfg.TxHash is not the zero hash (empty hash)
	// Why validate? Tracing is expensive; fail fast on invalid inputs

	// TODO: Call TraceTransaction to fetch execution trace
	// - Call client.TraceTransaction(ctx, cfg.TxHash)
	// - This calls debug_traceTransaction RPC method under the hood
	// - Returns json.RawMessage containing trace data
	// - The trace format depends on the tracer used (default: callTracer)
	// - Handle potential errors from the RPC call
	// - Validate that the trace response is not nil

	// TODO: Copy the trace data for safe return
	// - Create a new json.RawMessage with same length as raw
	// - Use copy() to duplicate the bytes
	// - Why? json.RawMessage is []byte (slice), which is a reference type
	// - Without copying, caller mutations would affect our internal data

	// TODO: Construct and return the Result
	// - Create a Result struct with:
	//   - TxHash: The transaction hash that was traced
	//   - Trace: The copied trace data
	// - Return the result and nil error on success

	return nil, errors.New("not implemented")
}
