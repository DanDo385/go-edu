//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var transferSigHash = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

/*
Problem: Query and decode ERC20 Transfer events from blockchain logs.

This module teaches you how to work with Ethereum events/logs. Events are append-only
records emitted during contract execution. They're searchable via bloom filters and
provide an efficient way to track state changes without querying contract storage.

Computer science principles highlighted:
  - Event-driven architecture: Logs as append-only audit trail
  - Bloom filters: Probabilistic data structures for efficient searching
  - Indexed vs non-indexed parameters: Trade-off between searchability and cost
  - Log structure: Topics (indexed) vs Data (non-indexed)
*/
func Run(ctx context.Context, client LogClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Defensive Programming Pattern
	// ============================================================================
	// By module 09, input validation should be automatic. The pattern repeats:
	// validate context, validate client, validate required config fields.
	//
	// Building on previous modules: This is the same validation from modules 01-08.
	// The repetition reinforces defensive programming as a habit.
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}

	// Token address is required because we need to know which contract emitted the logs.
	// Each contract's events are separate—Transfer events from different tokens are
	// different events even though they have the same signature.
	if cfg.Token == (common.Address{}) {
		return nil, errors.New("token address required")
	}

	// ============================================================================
	// STEP 2: Build FilterQuery - Understanding Log Filtering
	// ============================================================================
	// FilterQuery tells the node which logs we want. It's like a SQL WHERE clause:
	//   - Addresses: Which contracts to query
	//   - FromBlock/ToBlock: Block range to search
	//   - Topics: Which events to match (by signature and indexed parameters)
	//
	// Why filter? Ethereum mainnet has millions of logs. Without filtering, querying
	// all logs would be impossibly slow. Filtering narrows the search to what we need.
	//
	// Bloom filters: Each block header contains a bloom filter of all log topics in
	// that block. The node uses bloom filters to quickly identify which blocks might
	// contain matching logs, then scans only those blocks. This is much faster than
	// scanning every block.
	//
	// Computer science concept: Bloom filters are probabilistic data structures.
	// They can have false positives (say a block has a log when it doesn't) but never
	// false negatives (never miss a log that exists). This is perfect for filtering—
	// we might check a few extra blocks, but we never miss relevant logs.
	query := ethereum.FilterQuery{
		// Addresses: Only query logs from this specific token contract.
		// Why? Different tokens emit different Transfer events. We only want logs
		// from the token we're interested in.
		Addresses: []common.Address{cfg.Token},

		// FromBlock/ToBlock: Define the block range to search.
		// - nil FromBlock means start from genesis (block 0)
		// - nil ToBlock means search up to latest block
		// - Specific numbers limit the range for performance
		//
		// Why block range? Ethereum mainnet has 18M+ blocks. Searching all blocks
		// takes time. Limiting the range makes queries faster and more focused.
		FromBlock: cfg.FromBlock,
		ToBlock:   cfg.ToBlock,

		// Topics: Filter by event signature and indexed parameters.
		// Topics is [][]common.Hash:
		//   - Outer slice: Topic positions (0, 1, 2, 3)
		//   - Inner slice: OR options for each position
		//
		// Topic[0] is always the event signature hash:
		//   keccak256("Transfer(address,address,uint256)") = 0xddf252ad...
		//
		// Why keccak256? Ethereum uses keccak256 for all hashing (blocks, transactions,
		// events). The event signature uniquely identifies which event this is.
		//
		// Pattern: {{eventSig}} means "only logs with this event signature in Topic[0]"
		Topics: [][]common.Hash{{transferSigHash}},
	}

	// ============================================================================
	// STEP 3: Optionally Filter by Sender Address - Understanding Indexed Parameters
	// ============================================================================
	// Transfer event signature: event Transfer(address indexed from, address indexed to, uint256 value)
	//
	// Indexed parameters go into topics:
	//   - Topic[0]: Event signature (Transfer)
	//   - Topic[1]: from address (first indexed param)
	//   - Topic[2]: to address (second indexed param)
	//
	// Non-indexed parameters go into data:
	//   - value (uint256) is in the data field
	//
	// Why indexed vs non-indexed?
	//   - Indexed: Searchable via bloom filters (can filter by value)
	//   - Non-indexed: Not searchable, but cheaper to emit (lower gas cost)
	//
	// Trade-off: Index parameters you want to search by (addresses, small values).
	// Don't index large data (strings, arrays) because it's expensive.
	//
	// FromHolder filter: If provided, only return transfers FROM this address.
	// We append addressTopic(*cfg.FromHolder) to Topic[1] position.
	if cfg.FromHolder != nil {
		// Convert address to topic: Addresses are 20 bytes, but topics are 32 bytes.
		// We left-pad the address with zeros to make it 32 bytes.
		query.Topics = append(query.Topics, []common.Hash{addressTopic(*cfg.FromHolder)})
	}

	// ============================================================================
	// STEP 4: Optionally Filter by Recipient Address - Understanding Topic Positions
	// ============================================================================
	// ToHolder filter: If provided, only return transfers TO this address.
	// We append addressTopic(*cfg.ToHolder) to Topic[2] position.
	//
	// Complexity: If filtering by `to` but not `from`, we need a nil placeholder
	// for Topic[1]. Topic positions must be sequential.
	//
	// Example scenarios:
	//   1. Filter by from only: Topics = [{eventSig}, {from}]
	//   2. Filter by to only: Topics = [{eventSig}, nil, {to}]
	//   3. Filter by both: Topics = [{eventSig}, {from}, {to}]
	//
	// Why nil? nil means "match any value in this position". We need it to preserve
	// topic position alignment when skipping Topic[1].
	if cfg.ToHolder != nil {
		// If we haven't added Topic[1] (FromHolder), add nil placeholder
		if len(query.Topics) == 1 {
			query.Topics = append(query.Topics, nil)
		}
		// Add Topic[2] (ToHolder)
		query.Topics = append(query.Topics, []common.Hash{addressTopic(*cfg.ToHolder)})
	}

	// ============================================================================
	// STEP 5: Execute Log Query - Understanding eth_getLogs RPC
	// ============================================================================
	// FilterLogs sends an eth_getLogs JSON-RPC request to the node with our filter.
	// The node:
	//   1. Uses bloom filters to identify candidate blocks
	//   2. Scans those blocks for matching logs
	//   3. Returns all matching logs as []types.Log
	//
	// Performance considerations:
	//   - Large block ranges (>10,000 blocks) can be slow
	//   - Many matching logs (>1000) can hit node limits
	//   - Public RPC providers often limit eth_getLogs queries
	//
	// Error handling: FilterLogs can fail if:
	//   - Network error (RPC call failed)
	//   - Invalid block range (ToBlock < FromBlock)
	//   - Block range too large (provider limits)
	//   - Query returned too many results (provider limits)
	//
	// Building on previous modules: This is the same error wrapping pattern from
	// modules 01-08. Always wrap errors with context for debugging.
	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("filter logs: %w", err)
	}

	// ============================================================================
	// STEP 6: Initialize Result with Preallocated Slice - Performance Optimization
	// ============================================================================
	// We preallocate the Events slice with capacity len(logs). This avoids
	// reallocations during the append loop.
	//
	// How Go slices work:
	//   - make([]T, 0, cap) creates a slice with length 0 and capacity cap
	//   - append() adds elements; if length exceeds capacity, Go allocates a new
	//     larger backing array and copies all elements (expensive!)
	//   - By preallocating, we avoid these reallocations
	//
	// Performance impact:
	//   - Without preallocation: O(n log n) due to repeated reallocations
	//   - With preallocation: O(n) single pass through logs
	//
	// When to preallocate: When you know the final size upfront. Here we know
	// we'll have exactly len(logs) events (one per log).
	//
	// Computer science concept: This is the difference between amortized O(1)
	// append (with preallocation) vs amortized O(log n) append (without).
	result := &Result{
		Events: make([]TransferEvent, 0, len(logs)),
	}

	// ============================================================================
	// STEP 7: Decode Each Log - Understanding Log Structure
	// ============================================================================
	// types.Log contains:
	//   - Address: Contract that emitted the log
	//   - Topics: []common.Hash of indexed parameters
	//   - Data: []byte of non-indexed parameters
	//   - BlockNumber: Which block contains this log
	//   - TxHash: Which transaction emitted this log
	//   - Index: Position within the block's logs
	//
	// For Transfer events:
	//   - Topics[0]: Event signature hash
	//   - Topics[1]: from address (indexed)
	//   - Topics[2]: to address (indexed)
	//   - Data: ABI-encoded value (uint256, non-indexed)
	//
	// decodeTransferLog helper: Extracts from/to/value from the raw log structure.
	// It handles:
	//   - Topic validation (correct event signature, enough topics)
	//   - Address extraction from topics (converting 32-byte hash to 20-byte address)
	//   - Value decoding from data (parsing uint256 from bytes)
	//
	// Error handling: Decoding can fail if:
	//   - Log has wrong number of topics (malformed event)
	//   - Topic[0] doesn't match Transfer signature (wrong event type)
	//   - Data is too short (incomplete value encoding)
	//
	// Building on previous modules: This is the same "process each item" loop
	// pattern from previous modules. Process, check error, accumulate results.
	for _, lg := range logs {
		event, err := decodeTransferLog(lg)
		if err != nil {
			return nil, err
		}
		result.Events = append(result.Events, event)
	}

	// ============================================================================
	// STEP 8: Return Results - Understanding Event Ordering
	// ============================================================================
	// Events are returned in the order they appear on the blockchain:
	//   1. Sorted by block number (ascending)
	//   2. Within a block, sorted by transaction index
	//   3. Within a transaction, sorted by log index
	//
	// Why this order? It reflects the chronological order of events as they
	// occurred on the blockchain. This is important for reconstructing state:
	// if you process events in order, you see state changes as they happened.
	//
	// Example: Account balance tracking
	//   - Process events in order to calculate running balance
	//   - Each Transfer updates the balance correctly
	//   - Processing out of order would give wrong intermediate balances
	//
	// Building on previous concepts:
	//   - Module 07: Queried current state (name, symbol, balance)
	//   - Module 08: Same, but with typed bindings
	//   - Module 09: Query historical state changes (all transfers over time)
	//
	// This progression shows how events complement state queries:
	//   - State queries: "What is the balance now?"
	//   - Event queries: "How did the balance change over time?"
	return result, nil
}

// addressTopic converts a 20-byte address to a 32-byte topic.
//
// Why 32 bytes? EVM uses 32-byte words for everything. Topics must be 32 bytes.
// Addresses are only 20 bytes, so we left-pad with 12 zero bytes.
//
// Padding scheme: [12 zero bytes][20 address bytes] = 32 bytes total
//
// This is the standard Ethereum encoding for address topics. All nodes and
// tools use the same padding, ensuring compatibility.
func addressTopic(addr common.Address) common.Hash {
	return common.BytesToHash(common.LeftPadBytes(addr.Bytes(), 32))
}

// decodeTransferLog extracts Transfer event data from a raw log entry.
//
// Transfer event structure:
//   - Topics[0]: keccak256("Transfer(address,address,uint256)")
//   - Topics[1]: from address (32 bytes, last 20 are address)
//   - Topics[2]: to address (32 bytes, last 20 are address)
//   - Data: value (32 bytes, uint256 in big-endian)
//
// Why Topics[1:] instead of Topics[0:]? Topic[0] is the event signature, which
// we already validated in the filter. Topics[1] and Topics[2] are the indexed
// parameters we want to extract.
//
// Error handling: This function validates:
//   1. Enough topics (need at least 3: signature + from + to)
//   2. Correct event signature (prevent decoding wrong event type)
//   3. Enough data (need at least 32 bytes for value)
//
// These checks prevent panics and catch malformed logs early.
func decodeTransferLog(lg types.Log) (TransferEvent, error) {
	// Validate topic count: Transfer has 3 topics (signature + from + to)
	if len(lg.Topics) < 3 {
		return TransferEvent{}, fmt.Errorf("log %s missing topics", lg.TxHash.Hex())
	}

	// Validate event signature: Ensure this is actually a Transfer event
	if lg.Topics[0] != transferSigHash {
		return TransferEvent{}, fmt.Errorf("unexpected topic %s", lg.Topics[0].Hex())
	}

	// Extract from address: Topics[1] is 32 bytes, last 20 bytes are the address.
	// We use [12:] to skip the 12 zero-padding bytes at the start.
	from := common.BytesToAddress(lg.Topics[1].Bytes()[12:])

	// Extract to address: Same process for Topics[2]
	to := common.BytesToAddress(lg.Topics[2].Bytes()[12:])

	// Validate data length: value is uint256 (32 bytes)
	if len(lg.Data) < 32 {
		return TransferEvent{}, fmt.Errorf("log %s data too short", lg.TxHash.Hex())
	}

	// Extract value: Last 32 bytes of Data are the uint256 value in big-endian.
	// We use [len(lg.Data)-32:] to get exactly the last 32 bytes, which handles
	// cases where Data might have extra padding (though it shouldn't for Transfer).
	value := new(big.Int).SetBytes(lg.Data[len(lg.Data)-32:])

	// Return decoded event with metadata:
	//   - BlockNumber: Which block contains this event
	//   - TxHash: Which transaction emitted this event
	//   - LogIndex: Position within block's logs (for ordering)
	//   - From/To/Value: The actual Transfer data
	return TransferEvent{
		BlockNumber: lg.BlockNumber,
		TxHash:      lg.TxHash,
		LogIndex:    lg.Index,
		From:        from,
		To:          to,
		Value:       value,
	}, nil
}
