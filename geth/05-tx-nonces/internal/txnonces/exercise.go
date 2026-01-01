//go:build !solution && !reference

package txnonces

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func Run(ctx context.Context, client TXClient, cfg Config) (*Result, error) {
	// TODO: Input Validation and Defaults
	// TODO: Determine Sender Address and Nonce
	// TODO: Get Network and Gas Parameters
	// TODO: Create and Sign the Transaction
	// TODO: Broadcast the Transaction
	// TODO: Return the Result
	panic("not implemented")
}
