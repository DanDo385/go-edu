//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

/*
Run contains the reference solution for fetching and summarizing receipts.

Receipts are the execution results of transactions. They answer: "Did the transaction
succeed? How much gas did it use? What logs were emitted?" This is essential for
dApps, indexers, and block explorers.

Computer science principles highlighted:
  - Execution results vs execution trace (receipts vs traces from module 13)
  - Defensive copying for nested data structures (logs contain slices)
  - Status codes as execution summaries (0 = failure, 1 = success)
*/
func Run(ctx context.Context, client ReceiptClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Defensive Programming Pattern
	// ============================================================================
	// Same validation pattern as all previous modules. Validate inputs before
	// making RPC calls.
	//
	if ctx == nil {
		ctx = context.Background()
	}

	// Client validation: Same pattern from modules 01, 11, 12, 13, 14.
	if client == nil {
		return nil, errors.New("client is nil")
	}

	// Transaction hash validation: Receipts are per-transaction. We need a valid
	// transaction hash to fetch its receipt.
	if cfg.TxHash == (common.Hash{}) {
		return nil, errors.New("tx hash required")
	}

	// ============================================================================
	// STEP 2: Fetch Receipt - Understanding Execution Results
	// ============================================================================
	// TransactionReceipt fetches the execution result for a transaction. This is
	// different from tracing (module 13):
	//   - Traces: Detailed execution steps (opcode-level)
	//   - Receipts: High-level execution summary (status, gas, logs)
	//
	// What's in a receipt?
	//   - Status: Success (1) or failure (0)
	//   - GasUsed: Actual gas consumed
	//   - CumulativeGasUsed: Total gas used in block up to this transaction
	//   - Logs: Events emitted by the transaction
	//   - ContractAddress: If this was a contract creation, the new contract's address
	//   - PostState: (Pre-Byzantium) State root after transaction
	//
	// Why receipts matter:
	//   - dApps: Check if transaction succeeded
	//   - Indexers: Extract emitted events
	//   - Block explorers: Display transaction results
	//   - Gas accounting: Track actual gas consumption
	//
	// Building on previous modules: We've seen blocks (14), traces (13), now receipts.
	// Different views of the same transaction data!
	rcpt, err := client.TransactionReceipt(ctx, cfg.TxHash)
	if err != nil {
		return nil, fmt.Errorf("transaction receipt: %w", err)
	}

	// Validate non-nil response: Same pattern from all previous modules.
	if rcpt == nil {
		return nil, errors.New("nil receipt")
	}

	// ============================================================================
	// STEP 3: Process Logs with Defensive Copying
	// ============================================================================
	// Logs are events emitted by contracts during execution. Each log contains:
	//   - Address: Contract that emitted the log
	//   - Topics: Indexed event parameters (up to 4 topics, first is event signature)
	//   - Data: Non-indexed event parameters (ABI-encoded)
	//   - Index: Position of log within the transaction's logs
	//
	// Why defensive copying? Both Topics and Data are slices (reference types).
	// If we return pointers to receipt's internal data, callers could mutate them.
	// By copying, we ensure each caller gets independent, isolated data.
	//
	// Connection to Solidity-edu module 03 (Events): Logs are the low-level
	// representation of Solidity events. The Topics array contains the event
	// signature (topic[0]) and indexed parameters (topics[1-3]). The Data field
	// contains non-indexed parameters.
	//
	// Building on module 12: Same defensive copying pattern for slices.
	logs := make([]LogSummary, 0, len(rcpt.Logs))
	for _, lg := range rcpt.Logs {
		// Copy Topics: []common.Hash is a slice. Use append with nil slice to create
		// an independent copy. This prevents callers from mutating the original data.
		topicsCopy := append([]common.Hash(nil), lg.Topics...)

		// Copy Data: []byte is a slice. Same defensive copying pattern.
		dataCopy := append([]byte(nil), lg.Data...)

		// Construct LogSummary with copied data.
		logs = append(logs, LogSummary{
			Address: lg.Address,    // common.Address is value type (array), auto-copied
			Topics:  topicsCopy,    // Copied slice
			Data:    dataCopy,      // Copied slice
			Index:   uint(lg.Index), // uint64 → uint conversion, value type
		})
	}

	// ============================================================================
	// STEP 4: Construct Result with All Receipt Data
	// ============================================================================
	// We build a Result struct containing all important receipt fields. Each field
	// needs appropriate handling:
	//   - Value types (StatusOK, GasUsed, etc.): Auto-copied
	//   - big.Int: Mutable, needs defensive copy with new(big.Int).Set()
	//   - []byte: Reference type, needs defensive copy with append()
	//   - Hash/Address: Value types (arrays), auto-copied
	//
	// Status field explanation:
	//   - Post-Byzantium (EIP-658): Status field is 1 (success) or 0 (failure)
	//   - Pre-Byzantium: No status field, used PostState root instead
	//   - We convert Status (uint64) to StatusOK (bool) for clarity
	//
	// ContractAddress field:
	//   - Non-zero if transaction created a contract
	//   - Zero address if transaction was a regular call or transfer
	//
	// Building on previous concepts:
	//   - Defensive copying from modules 01, 12, 13
	//   - Value vs reference types from all modules
	//   - Receipt structure understanding
	return &Result{
		TxHash:        rcpt.TxHash,                           // Transaction identifier (Hash is value type)
		BlockNumber:   new(big.Int).Set(rcpt.BlockNumber),  // Defensive copy: big.Int is mutable
		StatusOK:      rcpt.Status == 1,                    // Convert status (1/0) to bool (true/false)
		GasUsed:       rcpt.GasUsed,                        // Actual gas consumed (uint64 is value type)
		CumulativeGas: rcpt.CumulativeGasUsed,              // Total gas used in block up to this tx
		Contract:      rcpt.ContractAddress,                // Contract address if creation, else zero
		Logs:          logs,                                // Copied logs slice
		PostStateRoot: append([]byte(nil), rcpt.PostState...), // Defensive copy: []byte is reference type
	}, nil
	// ============================================================================
	// STEP 5: Complete - Understanding the Receipt Lifecycle
	// ============================================================================
	// The progression of transaction data:
	//   1. Transaction sent (module 05-06): User creates and signs transaction
	//   2. Transaction executed (module 13): EVM processes transaction, emits logs
	//   3. Receipt generated (this module): Results stored in receipt trie
	//   4. Block finalized (module 14): Receipt committed to blockchain
	//
	// Receipts are stored in a separate Merkle-Patricia trie per block. The
	// receiptRoot in block header commits to all receipts. This allows efficient
	// verification that a receipt is part of the blockchain (similar to storage
	// proofs from module 12).
	//
	// Why separate receipt trie? Because receipts are derived data (computed
	// from transaction execution), not part of transaction data itself. This
	// separation keeps transaction data minimal while preserving execution results.
}
