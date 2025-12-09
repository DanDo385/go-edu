//go:build !solution
// +build !solution

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
	// TODO: Validate input parameters
	// - Check if ctx is nil and provide a default context if needed
	// - Check if client is nil and return an appropriate error
	// - Check if cfg.Contract is the zero address and return an error
	// - Contract address is required because we need to know which contract to call

	// TODO: Create a helper function for making contract calls
	// - This function should take a function selector (4 bytes) as input
	// - Build an ethereum.CallMsg with:
	//   * To: pointer to cfg.Contract address
	//   * Data: the function selector (encoded function call)
	// - Call client.CallContract(ctx, msg, cfg.BlockNumber) to execute
	// - Return the raw bytes or error
	// - Why a helper? This pattern (build CallMsg → CallContract) repeats for each function
	// - This demonstrates the DRY principle (Don't Repeat Yourself)

	// TODO: Call and decode name() function
	// - Use the selectorName constant (already defined at top of file)
	// - Call the helper function with selectorName to get raw bytes
	// - Handle errors from the call (network failures, contract reverts)
	// - Decode the raw bytes using decodeString helper (already implemented)
	// - Handle errors from decoding (invalid ABI encoding)
	// - name() returns a string, which is a dynamic type in ABI encoding
	// - String encoding: offset (32 bytes) + length (32 bytes) + data (padded)

	// TODO: Call and decode symbol() function
	// - Follow the same pattern as name()
	// - Use selectorSymbol constant
	// - Call the helper → decode with decodeString → handle errors
	// - Notice the repetition: call → check error → decode → check error
	// - This pattern is fundamental to all contract interactions

	// TODO: Call and decode decimals() function
	// - Follow the same pattern but use selectorDecimals
	// - Decode with decodeUint8 (already implemented)
	// - decimals() returns uint8, which is a static type
	// - Static types are simpler: just 32 bytes, right-aligned
	// - Why uint8? Decimals are typically 6-18, so uint8 (0-255) is sufficient

	// TODO: Call and decode totalSupply() function
	// - Follow the same pattern but use selectorTotalSupply
	// - Decode with decodeUint256 (already implemented)
	// - totalSupply() returns uint256, a static 32-byte type
	// - uint256 can represent very large numbers (up to 2^256 - 1)
	// - This is why we use *big.Int in Go (native ints would overflow)

	// TODO: Construct and return the Result struct
	// - Create a new Result struct with all decoded values:
	//   * Name: decoded name string
	//   * Symbol: decoded symbol string
	//   * Decimals: decoded decimals uint8
	//   * TotalSupply: decoded totalSupply *big.Int
	// - Return the result with nil error on success
	// - Why no defensive copying? These values are already copies from decoding

	return nil, errors.New("not implemented")
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
