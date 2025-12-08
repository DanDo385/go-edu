//go:build solution
// +build solution

package exercise

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
	// ============================================================================
	// Why validate inputs? This function is a library API that will be called by
	// other code. We can't trust callers to always pass valid inputs. This is a
	// fundamental principle of robust software design: "be liberal in what you
	// accept, but conservative in what you return" (Postel's Law).
	//
	// Context handling: In Go, context.Context is the idiomatic way to handle
	// cancellation, timeouts, and request-scoped values. If ctx is nil, we provide
	// context.Background() as a safe default. This pattern repeats throughout
	// Ethereum Go codebases—always check for nil contexts and provide defaults.
	//
	// This concept builds on: Go's context package fundamentals, which you'll see
	// repeated in every RPC call throughout this course.
	if ctx == nil {
		ctx = context.Background()
	}
	
	// Client validation: The RPCClient interface is what we depend on. If it's nil,
	// we can't proceed. This is a critical check because Go's zero value for interfaces
	// is nil, and calling methods on nil interfaces causes panics.
	//
	// Error handling pattern: We return early with a descriptive error. This is
	// Go's idiomatic error handling—fail fast, don't continue with invalid state.
	// This pattern (validate → return error) repeats in every function we'll write.
	if client == nil {
		return nil, errors.New("client is nil")
	}

	// ============================================================================
	// STEP 2: Retrieve Chain ID - Understanding Replay Protection
	// ============================================================================
	// Chain ID is fundamental to Ethereum's security model. Introduced in EIP-155,
	// it prevents "replay attacks" where a transaction signed on one network
	// (e.g., mainnet) could be replayed on another network (e.g., a testnet).
	//
	// How it works: When you sign a transaction, the chain ID is included in the
	// signature. If someone tries to replay that transaction on a different chain,
	// the signature verification will fail because the chain ID won't match.
	//
	// Why we check for nil: The RPCClient interface returns (*big.Int, error).
	// In Go, a pointer can be nil even if error is nil (though this shouldn't happen
	// with well-behaved clients). We validate to prevent nil pointer dereferences
	// later. This defensive programming pattern is critical in production code.
	//
	// Error wrapping: We use fmt.Errorf with %w verb to wrap the original error.
	// This preserves the error chain, allowing callers to use errors.Is() and
	// errors.As() to inspect underlying errors. This is Go's idiomatic error
	// handling pattern introduced in Go 1.13.
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}
	if chainID == nil {
		return nil, errors.New("chain id response was nil")
	}

	// ============================================================================
	// STEP 3: Retrieve Network ID - Legacy Identifier Pattern
	// ============================================================================
	// Network ID predates Chain ID and was used for P2P networking (identifying
	// which network peers belong to). While Chain ID is now the standard for
	// transaction signing, Network ID is still useful for compatibility and
	// network identification.
	//
	// Why retrieve both? Some older tools and libraries still rely on Network ID.
	// By retrieving both, we provide maximum compatibility. This is a common
	// pattern in Ethereum tooling: support both old and new standards during
	// transition periods.
	//
	// Notice the pattern: We're following the exact same structure as Chain ID
	// retrieval. This repetition is intentional—it demonstrates how consistent
	// patterns make code predictable and easier to understand. You'll see this
	// "call → check error → validate nil → use" pattern repeated throughout
	// all Ethereum Go code.
	networkID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("network id: %w", err)
	}
	if networkID == nil {
		return nil, errors.New("network id response was nil")
	}

	// ============================================================================
	// STEP 4: Retrieve Block Header - Understanding Block Structure
	// ============================================================================
	// Block headers are the "lightweight" representation of blocks. They contain
	// cryptographic commitments (hashes) to the full block data without including
	// the data itself. This is a Merkle tree pattern—you can verify data integrity
	// without downloading everything.
	//
	// Key fields in headers:
	//   - stateRoot: Merkle root of the entire Ethereum state (all accounts, balances, storage)
	//   - transactionsRoot: Merkle root of all transactions in the block
	//   - receiptsRoot: Merkle root of all transaction receipts (logs, gas used, etc.)
	//   - parentHash: Links to the previous block, creating an immutable chain
	//
	// cfg.BlockNumber: This is nil for "latest" block, or a specific block number.
	// This demonstrates separation of concerns—configuration (which block to fetch)
	// is separate from logic (how to fetch it). This makes the function testable
	// and deterministic (tests can request specific blocks).
	//
	// Why headers instead of full blocks? Headers are ~500 bytes vs full blocks
	// which can be 100KB-2MB. If you only need metadata (block number, hash, timestamp),
	// headers are much more efficient. This is a performance optimization pattern
	// you'll use throughout Ethereum development.
	header, err := client.HeaderByNumber(ctx, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("header by number: %w", err)
	}
	if header == nil {
		return nil, errors.New("header response was nil")
	}

	// ============================================================================
	// STEP 5: Construct Result with Defensive Copying - Immutability Pattern
	// ============================================================================
	// CRITICAL CONCEPT: Defensive copying prevents data races and unexpected mutations.
	//
	// Why defensive copying?
	//   1. The RPCClient might return pointers to internal data structures
	//   2. If we return those pointers directly, callers could mutate them
	//   3. This could affect other callers or cause data races in concurrent code
	//   4. By copying, we ensure each caller gets their own independent copy
	//
	// big.Int copying: big.Int is a mutable type (unlike Go's primitive types).
	// If we just assigned chainID directly, both Result.ChainID and the internal
	// client data would point to the same big.Int. Mutating one would mutate both!
	// new(big.Int).Set(chainID) creates a new big.Int and copies the value.
	//
	// Header copying: types.Header is a struct, but it contains pointers and slices
	// internally. A shallow copy (just assigning header) would share those internal
	// pointers. types.CopyHeader() performs a deep copy, ensuring complete isolation.
	//
	// This immutability pattern is fundamental to safe concurrent programming in Go.
	// You'll see it repeated whenever we return data from external libraries or
	// when data might be shared across goroutines.
	//
	// Building on previous concepts:
	//   - We validated inputs (Step 1) → now we validate outputs
	//   - We handled errors consistently (Steps 2-4) → now we return success
	//   - We used context for cancellation → now we return data safely
	return &Result{
		ChainID:   new(big.Int).Set(chainID),   // Defensive copy: prevents mutation of client's internal data
		NetworkID: new(big.Int).Set(networkID), // Defensive copy: ensures each caller gets independent data
		Header:    types.CopyHeader(header),     // Deep copy: header contains pointers/slices that must be copied
	}, nil
}
