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
	// ============================================================================
	// Same validation pattern as all previous modules. Validate inputs before
	// making expensive RPC calls.
	//
	if ctx == nil {
		ctx = context.Background()
	}

	// Client validation: Same pattern from modules 01, 11, 12, 13.
	if client == nil {
		return nil, errors.New("client is nil")
	}

	// ============================================================================
	// STEP 2: Fetch Block - Understanding Block vs Header
	// ============================================================================
	// BlockByNumber fetches the full block including transactions. This is
	// heavier than HeaderByNumber (module 01) but gives us transaction data.
	//
	// When to use which:
	//   - HeaderByNumber: Lightweight, only block metadata (~500 bytes)
	//   - BlockByNumber: Full block with transactions (KBs to MBs)
	//
	// cfg.Number: nil = latest block, or specific block number for historical queries.
	//
	// Building on module 01: Same RPC call pattern, different method.
	block, err := client.BlockByNumber(ctx, cfg.Number)
	if err != nil {
		// Error message includes which block we tried to fetch (latest vs specific number).
		// This helps debugging when calls fail.
		target := "latest"
		if cfg.Number != nil {
			target = cfg.Number.String()
		}
		return nil, fmt.Errorf("fetch block %s: %w", target, err)
	}

	// Validate non-nil response: Same pattern from all previous modules.
	if block == nil {
		return nil, errors.New("nil block response")
	}

	// ============================================================================
	// STEP 3: Extract Block Metadata - Data Aggregation Pattern
	// ============================================================================
	// We extract key fields from the block header to build our explorer view.
	// This is data aggregation—selecting relevant fields for presentation.
	//
	// Why extract specific fields? blocks contain lots of data. We only expose
	// what's useful for explorer UI: number, hash, parent, gas usage.
	header := block.Header()
	res := &Result{
		Number:   header.Number.Uint64(),  // Block height
		Hash:     block.Hash(),            // Block identifier
		Parent:   header.ParentHash,       // Links to previous block (forms the chain!)
		TxCount:  len(block.Transactions()), // How many transactions in this block
		GasUsed:  header.GasUsed,          // Total gas consumed by all transactions
		GasLimit: header.GasLimit,         // Maximum gas allowed per block
	}

	// ============================================================================
	// STEP 4: Optionally Include Transaction Summaries - Controlled Expansion
	// ============================================================================
	// If cfg.IncludeTxs is true, we enumerate transactions and build summaries.
	// This is controlled expansion—callers decide how much detail they need.
	//
	// Why optional? Full transaction objects are large. Many use cases only need
	// block metadata. By making transactions optional, we optimize for common cases.
	//
	// What's a TxSummary? A lightweight view of a transaction containing only:
	//   - Hash: Transaction identifier
	//   - To: Recipient address (nil for contract creation)
	//   - Gas: Gas limit for this transaction
	//
	// Note: We don't include full transaction data (input data, signatures, etc.)
	// to keep summaries small. Callers can fetch full details if needed.
	if cfg.IncludeTxs {
		// Pre-allocate slice with exact capacity to avoid growth.
		res.Txs = make([]TxSummary, 0, len(block.Transactions()))

		for _, tx := range block.Transactions() {
			// tx.To() returns *common.Address (pointer). It can be nil for contract
			// creation transactions (where To field is not set, and the contract
			// address is derived from sender address and nonce).
			//
			// Why pointer? To distinguish between "zero address" (0x0000...0000)
			// and "no address" (nil). Contract creation uses nil.
			to := tx.To()

			res.Txs = append(res.Txs, TxSummary{
				Hash: tx.Hash(),  // Transaction identifier
				To:   to,         // Recipient (nil for contract creation)
				Gas:  tx.Gas(),   // Gas limit for this transaction
			})
		}
	}

	// ============================================================================
	// STEP 5: Return Result - Minimal Defensive Copying
	// ============================================================================
	// types.Block shares slices internally, but we only copy primitive values
	// (numbers, hashes) and create new TxSummary structs. This is safe because:
	//   - Numbers are value types (auto-copied)
	//   - Hashes are value types (arrays, not slices)
	//   - TxSummary contains only value types and a pointer (which is fine)
	//
	// We don't return pointers to Block's internal data, so callers can't
	// mutate our source data.
	//
	// Building on previous concepts:
	//   - We validated inputs (Step 1) → ensured safe operation
	//   - We handled errors (Steps 2-3) → provided helpful error messages
	//   - We used optional expansion (Step 4) → optimized for common cases
	//   - We used minimal copying (Step 5) → efficient without sacrificing safety
	return res, nil
}
