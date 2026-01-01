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

func Run(ctx context.Context, client ContractCaller, cfg Config) (*Result, error) {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func selector(sig string) []byte {
	// TODO: Implement this function
	panic("not implemented")
}

func decodeString(data []byte) (string, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func decodeUint8(data []byte) (uint8, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func decodeUint256(data []byte) (*big.Int, error) {
	// TODO: Implement this function
	panic("not implemented")
}
