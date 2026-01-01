//go:build !solution && !reference

package proofs

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
)

func Run(ctx context.Context, client ProofClient, cfg Config) (*Result, error) {
	// TODO: Input Validation
	// TODO: Convert Slots to Hex Strings
	// TODO: Fetch Proof
	// TODO: Build AccountProof with Defensive Copying
	// TODO: Process Storage Proofs
	// TODO: Return Complete Result
	panic("not implemented")
}
