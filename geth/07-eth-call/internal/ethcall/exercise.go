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
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
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
