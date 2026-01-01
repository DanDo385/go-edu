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

/*
Problem: Query ERC20 token metadata using manual ABI encoding/decoding.
*/

// Run - TODO: implement this function
func Run(ctx context.Context, client CallClient, cfg Config) (*Result, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

// selector - TODO: implement this function
func selector(sig string) []byte {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// decodeString - TODO: implement this function
func decodeString(data []byte) (string, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

// decodeUint8 - TODO: implement this function
func decodeUint8(data []byte) (uint8, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

// decodeUint256 - TODO: implement this function
func decodeUint256(data []byte) (*big.Int, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

