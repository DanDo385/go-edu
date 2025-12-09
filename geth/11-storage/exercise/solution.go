//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var zeroHash = common.Hash{}

/*
Problem: Read raw storage slots from Ethereum contracts, including mapping slots.

Storage is the cryptographic database where contracts store their persistent state.
Every contract has 2^256 possible 32-byte slots organized as a Merkle-Patricia trie.
Understanding storage layout is essential for:
  - Debugging contract state
  - Building indexers that track specific contract data
  - Verifying storage proofs (module 12)
  - Optimizing gas costs (packed storage)

Computer science principles highlighted:
  - Cryptographic commitment via Merkle trees (storage root commits to all slots)
  - Deterministic slot calculation (mapping slots via keccak256 hash)
  - Key-value store abstraction (2^256 address space maps to 32-byte values)
*/
func Run(ctx context.Context, client StorageClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Defensive Programming Pattern
	// ============================================================================
	// Why validate inputs? Storage RPC calls are network operations that consume
	// resources. By validating early, we fail fast with clear error messages
	// rather than waiting for a network round-trip to discover bad inputs.
	//
	// Context handling: Same pattern as module 01-stack. If ctx is nil, we provide
	// context.Background() as a safe default. This ensures our RPC calls always
	// have a valid context for cancellation and timeout handling.
	//
	// This pattern repeats throughout all Ethereum Go code: validate → default → proceed.
	if ctx == nil {
		ctx = context.Background()
	}

	// Client validation: The StorageClient interface is our dependency. If it's nil,
	// we can't make RPC calls. Returning early with a descriptive error is Go's
	// idiomatic error handling—fail fast, don't continue with invalid state.
	//
	// Building on module 01: Same validation pattern, different interface.
	if client == nil {
		return nil, errors.New("client is nil")
	}

	// Contract address validation: The zero address (0x0000...0000) is not a valid
	// contract address. In Ethereum, the zero address is special—it represents "no
	// address" and is used for contract creation transactions (where To is nil).
	//
	// Why check this? StorageAt calls with zero address will fail at the RPC level.
	// Better to catch this early with a clear error message.
	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address required")
	}

	// Slot validation: Storage slots are identified by numbers (0, 1, 2, ...).
	// A nil slot means the caller didn't specify which slot to read. This is
	// an error because we can't guess what they want to read.
	//
	// Note: Slot 0 is valid (it's where the first storage variable lives in Solidity).
	// We're only checking for nil, not for zero value.
	if cfg.Slot == nil {
		return nil, errors.New("slot is required")
	}

	// ============================================================================
	// STEP 2: Slot Hash Conversion - Understanding Storage Addressing
	// ============================================================================
	// Storage slots are conceptually numbered (0, 1, 2, ...) but the RPC interface
	// requires them as 32-byte hashes. This conversion pads the slot number to
	// 32 bytes using big-endian encoding.
	//
	// Why hashes? Ethereum's storage trie uses hashes as keys. The slot number
	// is converted to a hash to navigate the Merkle-Patricia trie structure.
	//
	// Example: Slot 5 → 0x0000000000000000000000000000000000000000000000000000000000000005
	//
	// Building on previous concepts: This is similar to how addresses and block
	// hashes are represented as 32-byte values throughout Ethereum's protocol.
	slotHash := slotToHash(cfg.Slot)

	// ============================================================================
	// STEP 3: Mapping Slot Calculation - Understanding Solidity Storage Layout
	// ============================================================================
	// If a mapping key is provided, we need to compute the actual storage slot
	// where the mapped value lives. This demonstrates how Solidity's mapping
	// type works under the hood!
	//
	// How Solidity mappings work:
	//   1. The mapping variable occupies a "base slot" (e.g., slot 0)
	//   2. The base slot itself is empty (mappings don't store data there)
	//   3. For each key, the actual slot is: keccak256(key || baseSlot)
	//   4. This distributes values across storage space, preventing collisions
	//
	// Example: mapping(address => uint256) balances at slot 0
	//   - balances[0x742d35...] → keccak256(0x742d35... || 0x00...00)
	//
	// Computer Science principle: This is a cryptographic hash table! The hash
	// function (keccak256) distributes keys uniformly across the 2^256 storage space.
	//
	// Why this matters:
	//   - Understanding this is essential for reading contract state
	//   - Light clients use these calculations to verify storage proofs
	//   - Indexers need this to track specific mapping entries
	//
	// Connection to Solidity-edu 06 (Mappings): This is the low-level implementation
	// of the mapping slot calculation formula you learned in Solidity.
	if len(cfg.MappingKey) > 0 {
		slotHash = mappingSlotHash(cfg.MappingKey, slotHash)
	}

	// ============================================================================
	// STEP 4: Storage Read - RPC Call Pattern
	// ============================================================================
	// Now we make the actual RPC call to read the storage value. This calls
	// eth_getStorageAt under the hood, which returns the raw 32-byte value
	// stored at the specified slot.
	//
	// Parameters:
	//   - ctx: Context for cancellation/timeout (passed through from caller)
	//   - cfg.Contract: The contract address to read from
	//   - slotHash: The resolved storage slot (either direct or mapping-computed)
	//   - cfg.BlockNumber: Which block to read from (nil = latest)
	//
	// Why specify block number? Storage values can change over time. By specifying
	// a block, we can read historical state. This is essential for:
	//   - Time-travel queries (what was the value at block X?)
	//   - Deterministic testing (always read from same block)
	//   - Historical analysis (track how storage changed)
	//
	// Error handling: We wrap the error with context about which slot we were
	// reading. This helps debugging when calls fail (common errors: contract
	// doesn't exist, node doesn't have historical state, network timeout).
	//
	// Building on module 01: Same RPC call pattern (call → check error → validate nil).
	value, err := client.StorageAt(ctx, cfg.Contract, slotHash, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("storage at slot %s: %w", slotHash.Hex(), err)
	}

	// ============================================================================
	// STEP 5: Result Construction - Informative Response Pattern
	// ============================================================================
	// We return both the resolved slot and the value. Why both?
	//
	// ResolvedSlot: Shows the actual slot that was queried. For mapping queries,
	// this is the computed slot (after hashing), not the base slot. This helps
	// callers understand where the data actually lives in storage.
	//
	// Value: The raw 32-byte storage value. This is uninterpreted bytes—it's up
	// to the caller to decode it based on the expected type (uint256, address,
	// bytes32, etc.). Different Solidity types are encoded differently:
	//   - uint256: Direct big-endian encoding
	//   - address: Right-padded with zeros (20 bytes + 12 zero bytes)
	//   - bool: 0x00...00 (false) or 0x00...01 (true)
	//   - bytes32: Direct 32-byte value
	//
	// Note: No defensive copying needed here! The value is []byte which is already
	// a copy (StorageAt returns a fresh slice). The slotHash is a common.Hash
	// which is a value type (array, not slice), so it's automatically copied.
	//
	// Building on previous concepts:
	//   - We validated inputs (Step 1) → now we return validated data
	//   - We handled errors consistently (Steps 1-4) → now we return success
	//   - We used context for cancellation (Step 4) → operation completed successfully
	return &Result{
		ResolvedSlot: slotHash,
		Value:        value,
	}, nil
}

func slotToHash(slot *big.Int) common.Hash {
	if slot == nil {
		return zeroHash
	}
	return common.BigToHash(slot)
}

func mappingSlotHash(key []byte, slot common.Hash) common.Hash {
	keyPadded := common.LeftPadBytes(key, 32)
	data := append(keyPadded, slot.Bytes()...)
	return common.BytesToHash(crypto.Keccak256(data))
}
