//go:build !solution && !reference

package txnonces

import (
	"context"
	"errors"
)

/*
Problem: Build, sign, and (optionally) send a legacy Ethereum transaction.

This module brings everything together: you'll use a private key to define a
sender, determine the correct nonce, construct a transaction, sign it with
replay protection, and broadcast it to the network. This is the complete
lifecycle of a state change on Ethereum.

Computer science principles highlighted:
  - Cryptographic signatures for authentication and integrity
  - Sequence numbers (nonces) for replay protection and ordering
  - Immutable messages (transactions) as state change proposals
  - Separation of concerns (building vs. signing vs. sending)
*/
func Run(ctx context.Context, client TXClient, cfg Config) (*Result, error) {
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
}

