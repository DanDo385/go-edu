//go:build !solution && !reference

package ethcall

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
	// TODO: Implement

	// ============================================================================
	// Why validate inputs? This function is a library API that will be called by
	// TODO: Implement

	// ============================================================================
	// STEP 2: Create Helper Function - DRY Principle
	// TODO: Implement

	// ============================================================================
	// Why a helper function? We'll make 4 contract calls (name, symbol, decimals,
	// TODO: Implement

	// ============================================================================
	// STEP 3: Call and Decode name() - Understanding Dynamic Types
	// TODO: Implement

	// ============================================================================
	// name() is an ERC20 view function that returns a string.
	// TODO: Implement

	// ============================================================================
	// STEP 4: Call and Decode symbol() - Pattern Repetition
	// TODO: Implement

	// ============================================================================
	// symbol() is another string-returning view function.
	// TODO: Implement

	// ============================================================================
	// STEP 5: Call and Decode decimals() - Understanding Static Types
	// TODO: Implement

	// ============================================================================
	// decimals() returns uint8, which is a static type in ABI encoding.
	// TODO: Implement

	// ============================================================================
	// STEP 6: Call and Decode totalSupply() - Understanding uint256
	// TODO: Implement

	// ============================================================================
	// totalSupply() returns uint256, the most common Solidity type.
	// TODO: Implement

	// ============================================================================
	// STEP 7: Construct and Return Result - No Defensive Copying Needed
	// TODO: Implement

	// ============================================================================
	// Why no defensive copying? In modules 01 and 06, we used defensive copying
	// TODO: Implement

	panic("unimplemented")
}

func selector(sig string) []byte {
	// TODO: Implement this function
	panic("unimplemented")
}

func decodeString(data []byte) (string, error) {
	// TODO: Implement this function
	panic("unimplemented")
}

func decodeUint8(data []byte) (uint8, error) {
	// TODO: Implement this function
	panic("unimplemented")
}

func decodeUint256(data []byte) (*big.Int, error) {
	// TODO: Implement this function
	panic("unimplemented")
}
