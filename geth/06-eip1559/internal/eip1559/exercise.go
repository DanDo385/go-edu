//go:build !solution && !reference


package eip1559

import (
	"context"
	"errors"
)

const defaultDynamicGasLimit = 21000

/*
Problem: Build and sign an EIP-1559 dynamic fee transaction with proper fee estimation.

EIP-1559 (London upgrade, August 2021) introduced a two-part fee structure:
  - Base Fee: Algorithmically determined, burned (removed from ETH supply)
  - Priority Fee (Tip): Paid to validators, incentivizes inclusion

This is more predictable than legacy transactions where users bid against each other.

Computer science principles highlighted:
  - Algorithm design: Base fee adjusts automatically based on block fullness (control theory)
  - Economic incentives: Fee burning aligns validator and user interests
  - Defensive copying: Protect mutable big.Int values from external mutation
  - Error handling: Validate all inputs and RPC responses
*/
func Run(ctx context.Context, client FeeClient, cfg Config) (*Result, error) {
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
}

