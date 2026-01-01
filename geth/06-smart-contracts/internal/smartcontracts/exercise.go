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
Reference implementation for smart contract interaction.

This is the complete, tested solution for the exercise.
Run tests with: go test -tags=reference -v ./...
*/
func Run(ctx context.Context, client ContractCaller, cfg Config) (*Result, error) {
	// TODO: Implement this function
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
