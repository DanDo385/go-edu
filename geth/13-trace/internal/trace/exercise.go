//go:build !solution
// +build !solution

package trace

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
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
}

