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
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	selectorName        = selector("name()")
	selectorSymbol      = selector("symbol()")
	selectorDecimals    = selector("decimals()")
	selectorTotalSupply = selector("totalSupply()")
)

/*
Problem: Query ERC20 token metadata using manual ABI encoding/decoding.

This module teaches you how to interact with contracts without using typed bindings.
You'll manually encode function selectors and decode return values, giving you a deep
understanding of how contract calls work at the ABI level.

Computer science principles highlighted:
  - ABI encoding/decoding: Understanding how function calls are encoded as bytes
  - Function selectors: First 4 bytes of keccak256(functionSignature)
  - eth_call: Simulating contract execution without sending transactions
  - Manual memory management: Decoding dynamic types (strings) from raw bytes
*/
func Run(ctx context.Context, client CallClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Defensive Programming Pattern
	// ============================================================================
	// Why validate inputs? This function is a library API that will be called by
	// other code. We can't trust callers to always pass valid inputs.
	//
	// Context handling: This pattern repeats from module 01-stack. If ctx is nil,
	// we provide context.Background() as a safe default. This is Go's idiomatic
	// way to handle cancellation, timeouts, and request-scoped values.
	//
	// Building on previous concepts: You saw this exact pattern in module 01-stack.
	// By module 07, this should feel automatic—validate context, validate client,
	// validate required config fields. This repetition builds muscle memory.
	if ctx == nil {
		ctx = context.Background()
	}

	// Client validation: The CallClient interface is what we depend on for making
	// RPC calls. If it's nil, we can't proceed. This check prevents nil pointer
	// panics later in the code.
	//
	// Error handling pattern: Return early with descriptive error. This is Go's
	// idiomatic error handling—fail fast, don't continue with invalid state.
	if client == nil {
		return nil, errors.New("client is nil")
	}

	// Contract address validation: For eth_call, we need a contract address to
	// query. The zero address (all zeros) is invalid for contract calls.
	//
	// Why check for zero address? In Ethereum, address(0) is often used as a
	// special value (e.g., token minting/burning). For contract calls, we need
	// a valid deployed contract address.
	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address required")
	}

	// ============================================================================
	// STEP 2: Create Helper Function - DRY Principle
	// ============================================================================
	// Why a helper function? We'll make 4 contract calls (name, symbol, decimals,
	// totalSupply), and they all follow the same pattern:
	//   1. Build CallMsg with contract address and function selector
	//   2. Execute eth_call via client.CallContract
	//   3. Return raw bytes or error
	//
	// This is the DRY (Don't Repeat Yourself) principle. By extracting common
	// logic into a helper, we reduce code duplication and make the code more
	// maintainable. If the call logic changes, we only update one place.
	//
	// Closure pattern: This helper function is a closure—it captures ctx, client,
	// and cfg from the outer scope. This is a common Go pattern for reducing
	// boilerplate while keeping code clean.
	//
	// Computer science concept: This demonstrates abstraction—hiding complex
	// details (building CallMsg, handling RPC) behind a simple interface.
	call := func(selector []byte) ([]byte, error) {
		// Build CallMsg: This struct describes the contract call we want to make.
		// - To: pointer to contract address (required for calls)
		// - Data: function selector + encoded arguments (just selector for view functions)
		//
		// Why pointer for To? The CallMsg struct uses *common.Address to allow
		// nil (for contract creation txs). For calls, we always need an address.
		msg := ethereum.CallMsg{
			To:   &cfg.Contract,
			Data: selector,
		}

		// Execute eth_call: This sends the call to the RPC node, which executes
		// the function locally (without creating a transaction) and returns the
		// result. No gas is consumed, no state is changed.
		//
		// cfg.BlockNumber: Allows querying historical state. nil means "latest".
		// This demonstrates separation of concerns—the caller controls which
		// block to query, we just execute the call.
		return client.CallContract(ctx, msg, cfg.BlockNumber)
	}

	// ============================================================================
	// STEP 3: Call and Decode name() - Understanding Dynamic Types
	// ============================================================================
	// name() is an ERC20 view function that returns a string.
	//
	// Function selector: selectorName is computed from "name()" signature:
	//   1. Hash the signature: keccak256("name()")
	//   2. Take first 4 bytes: 0x06fdde03
	//
	// Why 4 bytes? This is the ABI standard. 4 bytes gives us 4 billion possible
	// selectors, which is enough to avoid collisions in practice.
	//
	// Pattern: call → check error → decode → check error. This two-phase approach
	// (network operation + data processing) is common in distributed systems.
	// We separate concerns: networking vs parsing.
	nameBytes, err := call(selectorName)
	if err != nil {
		// Error wrapping: We add context ("call name()") to the error while
		// preserving the original error with %w. This creates a traceable error
		// chain for debugging.
		//
		// Building on previous modules: This error wrapping pattern repeats from
		// module 01-stack. By module 07, you should recognize this as standard
		// practice for all error returns.
		return nil, fmt.Errorf("call name(): %w", err)
	}

	// Decode string: Strings are dynamic types in ABI encoding. The encoding is:
	//   - Bytes 0-31: Offset to string data (usually 0x20 = 32)
	//   - Bytes 32-63: Length of string in bytes
	//   - Bytes 64+: UTF-8 string data (padded to 32-byte boundary)
	//
	// Why so complex? Dynamic types need this structure because the EVM uses
	// 32-byte words for everything. Static types (like uint256) fit in one word,
	// but strings can be any length, so they need offset + length + data.
	//
	// The decodeString helper (implemented below) handles this complexity for us.
	name, err := decodeString(nameBytes)
	if err != nil {
		return nil, fmt.Errorf("decode name(): %w", err)
	}

	// ============================================================================
	// STEP 4: Call and Decode symbol() - Pattern Repetition
	// ============================================================================
	// symbol() is another string-returning view function.
	//
	// Notice the pattern: This is EXACTLY the same as name()—call, check error,
	// decode, check error. The only differences are:
	//   1. Different selector (selectorSymbol vs selectorName)
	//   2. Different error messages for debugging
	//
	// Why repeat the pattern? Because it demonstrates consistency. Once you
	// understand one contract call, you understand them all. This predictability
	// is a key principle of good API design.
	//
	// Building on previous concepts: You've seen this "call RPC → validate →
	// process" pattern in every module so far (chain ID, network ID, headers,
	// transactions). It's the fundamental pattern of blockchain development.
	symbolBytes, err := call(selectorSymbol)
	if err != nil {
		return nil, fmt.Errorf("call symbol(): %w", err)
	}
	symbol, err := decodeString(symbolBytes)
	if err != nil {
		return nil, fmt.Errorf("decode symbol(): %w", err)
	}

	// ============================================================================
	// STEP 5: Call and Decode decimals() - Understanding Static Types
	// ============================================================================
	// decimals() returns uint8, which is a static type in ABI encoding.
	//
	// Static vs dynamic types:
	//   - Static: Fixed size (uint8, uint256, address, bool, bytes32)
	//   - Dynamic: Variable size (string, bytes, arrays, structs)
	//
	// Static types are simpler to decode: just 32 bytes, value is right-aligned.
	// For uint8, only the last byte contains the value, the rest are zero padding.
	//
	// Why uint8 for decimals? Token decimals are typically 6-18 (6 for USDC,
	// 18 for most tokens). uint8 supports 0-255, which is more than enough.
	// Using uint8 instead of uint256 saves gas when emitting events or storing.
	decBytes, err := call(selectorDecimals)
	if err != nil {
		return nil, fmt.Errorf("call decimals(): %w", err)
	}

	// Decode uint8: The decodeUint8 helper extracts the last byte from the
	// 32-byte response. This is simpler than string decoding—no offsets or
	// length fields to parse.
	decimals, err := decodeUint8(decBytes)
	if err != nil {
		return nil, fmt.Errorf("decode decimals(): %w", err)
	}

	// ============================================================================
	// STEP 6: Call and Decode totalSupply() - Understanding uint256
	// ============================================================================
	// totalSupply() returns uint256, the most common Solidity type.
	//
	// uint256 in Solidity = *big.Int in Go. Why?
	//   - Solidity uint256: 0 to 2^256-1 (huge range)
	//   - Go uint64: 0 to 2^64-1 (much smaller)
	//   - Go *big.Int: arbitrary precision (can handle uint256)
	//
	// Real-world example: Total supply of tokens can be trillions or more.
	// For instance, Shiba Inu token has a supply of 1 quadrillion (10^15).
	// This far exceeds uint64's max (~10^19), so we need big.Int.
	//
	// Pattern repetition: Again, call → check error → decode → check error.
	// By now, this should feel automatic. Every contract interaction follows
	// this pattern, regardless of the specific function being called.
	supplyBytes, err := call(selectorTotalSupply)
	if err != nil {
		return nil, fmt.Errorf("call totalSupply(): %w", err)
	}

	// Decode uint256: The decodeUint256 helper reads the last 32 bytes and
	// converts them to *big.Int. Unlike uint8, we use all 32 bytes because
	// uint256 can use the full range.
	totalSupply, err := decodeUint256(supplyBytes)
	if err != nil {
		return nil, fmt.Errorf("decode totalSupply(): %w", err)
	}

	// ============================================================================
	// STEP 7: Construct and Return Result - No Defensive Copying Needed
	// ============================================================================
	// Why no defensive copying? In modules 01 and 06, we used defensive copying
	// for big.Int values returned from the RPC client. Here, we don't need to
	// because these values are freshly created by our decoding functions.
	//
	// Ownership model:
	//   - RPC client returns: May share internal data → need defensive copy
	//   - Our decoders return: Fresh allocations → already independent copies
	//
	// This demonstrates understanding of ownership and mutation safety. Not all
	// *big.Int values need copying—only those that might be shared with other
	// code that could mutate them.
	//
	// Building on previous concepts:
	//   - Module 01: Learned RPC call pattern and defensive copying
	//   - Module 06: Applied defensive copying to transaction fees
	//   - Module 07: Recognize when defensive copying is NOT needed
	//
	// This progression shows maturity—knowing when NOT to apply a pattern is
	// as important as knowing when to apply it.
	return &Result{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
	}, nil
}

func selector(sig string) []byte {
	hash := crypto.Keccak256([]byte(sig))
	return hash[:4]
}

func decodeString(data []byte) (string, error) {
	if len(data) < 64 {
		return "", errors.New("data too short for string")
	}
	offset := new(big.Int).SetBytes(data[:32]).Int64()
	if offset < 0 || offset+32 > int64(len(data)) {
		return "", errors.New("invalid offset")
	}
	lengthStart := int(offset)
	lengthEnd := lengthStart + 32
	if lengthEnd > len(data) {
		return "", errors.New("invalid length data")
	}
	length := new(big.Int).SetBytes(data[lengthStart:lengthEnd]).Int64()
	if length < 0 {
		return "", errors.New("negative length")
	}
	dataStart := lengthEnd
	dataEnd := dataStart + int(length)
	if dataEnd > len(data) {
		return "", errors.New("string exceeds data bounds")
	}
	return string(data[dataStart:dataEnd]), nil
}

func decodeUint8(data []byte) (uint8, error) {
	if len(data) < 32 {
		return 0, errors.New("data too short for uint8")
	}
	return uint8(data[len(data)-1]), nil
}

func decodeUint256(data []byte) (*big.Int, error) {
	if len(data) < 32 {
		return nil, errors.New("data too short for uint256")
	}
	out := new(big.Int).SetBytes(data[len(data)-32:])
	return out, nil
}
