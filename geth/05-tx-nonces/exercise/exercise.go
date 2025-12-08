//go:build !solution
// +build !solution

package exercise

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
	// TODO: Validate input parameters
	// - Check for nil context, client, and private key.
	// - Provide default values for AmountWei (0) and GasLimit (21000).

	// TODO: Determine the sender's address from the private key.
	// - Use `crypto.PubkeyToAddress(cfg.PrivateKey.PublicKey)`.

	// TODO: Determine the nonce to use for the transaction.
	// - If `cfg.Nonce` is provided, use it.
	// - Otherwise, fetch the pending nonce for the sender's address using
	//   `client.PendingNonceAt(ctx, from)`.
	// - Handle and wrap any errors.

	// TODO: Get the chain ID for signing.
	// - Use `client.ChainID(ctx)`.
	// - Handle errors and check for a nil chain ID.

	// TODO: Determine the gas price to use.
	// - If `cfg.GasPrice` is provided, use it.
	// - Otherwise, get the suggested gas price from the network using
	//   `client.SuggestGasPrice(ctx)`.
	// - Handle and wrap any errors.

	// TODO: Create the legacy transaction object.
	// - Use `types.NewTransaction` to assemble the nonce, recipient address,
	//   amount, gas limit, gas price, and data.
	// - Remember to create defensive copies of mutable types like `big.Int`
	//   and byte slices.

	// TODO: Sign the transaction.
	// - Create a signer using `types.LatestSignerForChainID(chainID)`. This
	//   signer implements EIP-155 replay protection.
	// - Use `types.SignTx` to sign the transaction with the signer and the
	//   private key.
	// - Handle and wrap any errors.

	// TODO: Send the transaction to the network.
	// - If `cfg.NoSend` is false, use `client.SendTransaction(ctx, signedTx)`
	//   to broadcast the transaction.
	// - Handle and wrap any errors.

	// TODO: Construct and return the Result struct.
	// - The result should contain the sender's address, the nonce used, and
	//   the signed transaction.

	return nil, errors.New("not implemented")
}