//go:build !solution && !reference

package smartcontracts

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
Problem: Demonstrate smart contract interaction concepts in Go.

This module is primarily a tutorial-based learning experience using the Geth console.
The exercises are designed to be completed in the Geth JavaScript console, not in Go.

However, this Go implementation demonstrates the same concepts you learned in the console:
  - Creating function selectors from signatures
  - Making eth_call requests to read contract state
  - Decoding ABI-encoded return values
  - Understanding the Call vs Transaction distinction

After completing the console tutorial in README.md, you'll understand:
  - Call vs Transaction distinction
  - Contract addresses and ABIs
  - Console-based contract interaction
  - Event decoding

Module 07 (eth-call) will teach you how to do the same things in Go with more depth.

Computer science principles highlighted:
  - ABI encoding: How function calls are serialized as bytes
  - Function selectors: First 4 bytes of keccak256(signature)
  - eth_call: Read-only contract execution without transactions
*/
func Run(ctx context.Context, client ContractCaller, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Defensive Programming Pattern
	// ============================================================================
	// This pattern should be familiar from modules 01-05. We always validate
	// inputs before proceeding. This is especially important for library APIs.
	//
	// BREAKPOINT: Set a breakpoint here to inspect incoming parameters
	if ctx == nil {
		ctx = context.Background()
	}

	if client == nil {
		return nil, errors.New("client is nil")
	}

	if cfg.Contract == (common.Address{}) {
		return nil, errors.New("contract address required")
	}

	// ============================================================================
	// STEP 2: Define Function Selectors
	// ============================================================================
	// In the Geth console, you created a contract object and called methods directly.
	// Under the hood, each method call is encoded as:
	//   - First 4 bytes: keccak256(signature)[:4] (the "selector")
	//   - Remaining bytes: ABI-encoded arguments
	//
	// For view functions with no arguments, we only need the selector.
	//
	// BREAKPOINT: Inspect the selector bytes to see the 4-byte encoding
	selectorName := selector("name()")
	selectorSymbol := selector("symbol()")
	selectorDecimals := selector("decimals()")
	selectorTotalSupply := selector("totalSupply()")

	// ============================================================================
	// STEP 3: Create Helper for Contract Calls
	// ============================================================================
	// This helper encapsulates the eth_call pattern:
	//   1. Build a CallMsg with contract address and function selector
	//   2. Execute via client.CallContract (equivalent to eth_call RPC)
	//   3. Return raw bytes for decoding
	//
	// In the console, this happened automatically when you called myContract.name()
	// Here, we see the underlying mechanics.
	call := func(sel []byte) ([]byte, error) {
		msg := ethereum.CallMsg{
			To:   &cfg.Contract,
			Data: sel,
		}
		return client.CallContract(ctx, msg, cfg.BlockNumber)
	}

	// ============================================================================
	// STEP 4: Query name() - String Decoding
	// ============================================================================
	// Calling name() in console: myContract.name()
	// In Go: We send the selector, get back ABI-encoded bytes, and decode.
	//
	// BREAKPOINT: Inspect nameBytes to see the raw ABI-encoded string
	nameBytes, err := call(selectorName)
	if err != nil {
		return nil, fmt.Errorf("call name(): %w", err)
	}
	name, err := decodeString(nameBytes)
	if err != nil {
		return nil, fmt.Errorf("decode name(): %w", err)
	}

	// ============================================================================
	// STEP 5: Query symbol() - Same Pattern
	// ============================================================================
	// The pattern repeats: call → decode. This consistency is intentional.
	// Once you understand one contract call, you understand them all.
	symbolBytes, err := call(selectorSymbol)
	if err != nil {
		return nil, fmt.Errorf("call symbol(): %w", err)
	}
	symbol, err := decodeString(symbolBytes)
	if err != nil {
		return nil, fmt.Errorf("decode symbol(): %w", err)
	}

	// ============================================================================
	// STEP 6: Query decimals() - uint8 Decoding
	// ============================================================================
	// decimals() returns uint8, a static type. Much simpler to decode than strings.
	//
	// BREAKPOINT: Compare decBytes (32 bytes with value in last byte) vs nameBytes
	decBytes, err := call(selectorDecimals)
	if err != nil {
		return nil, fmt.Errorf("call decimals(): %w", err)
	}
	decimals, err := decodeUint8(decBytes)
	if err != nil {
		return nil, fmt.Errorf("decode decimals(): %w", err)
	}

	// ============================================================================
	// STEP 7: Query totalSupply() - uint256 Decoding
	// ============================================================================
	// totalSupply() returns uint256, which maps to *big.Int in Go.
	// This is the most common return type for token amounts.
	supplyBytes, err := call(selectorTotalSupply)
	if err != nil {
		return nil, fmt.Errorf("call totalSupply(): %w", err)
	}
	totalSupply, err := decodeUint256(supplyBytes)
	if err != nil {
		return nil, fmt.Errorf("decode totalSupply(): %w", err)
	}

	// ============================================================================
	// STEP 8: Return Result
	// ============================================================================
	// BREAKPOINT: Inspect the final result before returning
	return &Result{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
	}, nil
}

// selector computes the 4-byte function selector from a signature.
// Example: selector("name()") returns the first 4 bytes of keccak256("name()")
//
// In the Geth console, this happens automatically when you call myContract.name().
// The console looks up "name" in the ABI, finds the signature, computes the selector,
// and includes it in the eth_call request.
func selector(sig string) []byte {
	hash := crypto.Keccak256([]byte(sig))
	return hash[:4]
}

// decodeString decodes an ABI-encoded string from raw bytes.
//
// ABI encoding for strings:
//   - Bytes 0-31: Offset to string data (usually 0x20 = 32)
//   - Bytes 32-63: Length of string in bytes
//   - Bytes 64+: UTF-8 string data (padded to 32-byte boundary)
//
// This is the same decoding that happens in the Geth console, but here we see
// the raw mechanics.
func decodeString(data []byte) (string, error) {
	if len(data) < 64 {
		return "", errors.New("data too short for string")
	}

	// Read offset (usually 32)
	offset := new(big.Int).SetBytes(data[:32]).Int64()
	if offset < 0 || offset+32 > int64(len(data)) {
		return "", errors.New("invalid offset")
	}

	// Read length
	lengthStart := int(offset)
	lengthEnd := lengthStart + 32
	if lengthEnd > len(data) {
		return "", errors.New("invalid length data")
	}
	length := new(big.Int).SetBytes(data[lengthStart:lengthEnd]).Int64()
	if length < 0 {
		return "", errors.New("negative length")
	}

	// Read string data
	dataStart := lengthEnd
	dataEnd := dataStart + int(length)
	if dataEnd > len(data) {
		return "", errors.New("string exceeds data bounds")
	}

	return string(data[dataStart:dataEnd]), nil
}

// decodeUint8 decodes an ABI-encoded uint8 from raw bytes.
// The value is in the last byte of a 32-byte word.
func decodeUint8(data []byte) (uint8, error) {
	if len(data) < 32 {
		return 0, errors.New("data too short for uint8")
	}
	return uint8(data[len(data)-1]), nil
}

// decodeUint256 decodes an ABI-encoded uint256 from raw bytes.
// Returns a *big.Int since uint256 exceeds Go's native integer types.
func decodeUint256(data []byte) (*big.Int, error) {
	if len(data) < 32 {
		return nil, errors.New("data too short for uint256")
	}
	return new(big.Int).SetBytes(data[len(data)-32:]), nil
}
