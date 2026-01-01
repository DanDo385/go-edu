//go:build !solution && !reference

package eip1559

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func Run(ctx context.Context, client FeeClient, cfg Config) (*Result, error) {
	// TODO: Input Validation
	// TODO: Derive Sender Address from Private Key
	// TODO: Determine Transaction Nonce
	// TODO: Retrieve Chain ID
	// TODO: Fetch Block Header to Get Base Fee
	// TODO: Determine Max Priority Fee (Tip Cap)
	// TODO: Determine Max Fee Cap
	// TODO: Prepare Transaction Data
	// TODO: Construct DynamicFeeTx Struct
	// TODO: Wrap in Transaction Envelope
	// TODO: Sign the Transaction
	// TODO: Send Transaction to Network (Optional)
	// TODO: Construct and Return Result
	panic("not implemented")
}
