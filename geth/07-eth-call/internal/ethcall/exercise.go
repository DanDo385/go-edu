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

func Run(ctx context.Context, client CallClient, cfg Config) (*Result, error) {
	// TODO: Input Validation
	// TODO: Create Helper Function
	// TODO: Call and Decode name()
	// TODO: Call and Decode symbol()
	// TODO: Call and Decode decimals()
	// TODO: Call and Decode totalSupply()
	// TODO: Construct and Return Result
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
