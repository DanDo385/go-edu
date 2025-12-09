//go:build !solution
// +build !solution

package exercise

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
	// TODO: Validate context
	// - Check if ctx is nil and provide context.Background() as default
	// - Context is essential for cancellation and timeout propagation

	// TODO: Validate client
	// - Check if client is nil and return appropriate error
	// - Client is required for all RPC operations

	// TODO: Validate private key
	// - Check if cfg.PrivateKey is nil and return error
	// - Private key is required for transaction signing
	// - Why critical? Without key, we can't sign transactions

	// TODO: Set default AmountWei if not provided
	// - If cfg.AmountWei is nil, set to big.NewInt(0)
	// - This allows creating transactions that just call contracts without sending ETH
	// - Pattern: Use zero values as sensible defaults

	// TODO: Set default GasLimit if not provided
	// - If cfg.GasLimit is 0, set to defaultDynamicGasLimit (21000)
	// - 21000 is the cost of a basic ETH transfer (no data, no contract calls)
	// - Why 21000? 21000 gas for value transfer + 0 for empty data

	// TODO: Derive sender address from private key
	// - Use crypto.PubkeyToAddress(cfg.PrivateKey.PublicKey)
	// - This is the "from" address for the transaction
	// - Pattern repeats: Always derive address from public key, never trust user input

	// TODO: Determine transaction nonce
	// - If cfg.Nonce is provided, use it (allows manual nonce management)
	// - Otherwise, call client.PendingNonceAt(ctx, from) to get next nonce
	// - Handle errors from PendingNonceAt (network issues, invalid address)
	// - Why "pending"? Includes transactions in mempool, not just mined ones
	// - Pattern: Config allows override for testing/advanced use cases

	// TODO: Retrieve chain ID
	// - Call client.ChainID(ctx) to get the chain identifier
	// - Handle errors from ChainID call
	// - Validate chainID is not nil (defensive programming)
	// - Why needed? Chain ID is part of EIP-1559 signature (replay protection)

	// TODO: Fetch latest header to get base fee
	// - Call client.HeaderByNumber(ctx, cfg.BlockNumber) to get header
	// - Use cfg.BlockNumber (nil for latest, or specific block)
	// - Handle errors from HeaderByNumber call
	// - Validate header and header.BaseFee are not nil
	// - If BaseFee is nil, return error indicating pre-London block
	// - IMPORTANT: Make a defensive copy of baseFee (big.Int is mutable!)
	// - Why copy? header.BaseFee points to client's internal data

	// TODO: Determine max priority fee (tip cap)
	// - If cfg.MaxPriorityFee is provided, use it
	// - Otherwise, call client.SuggestGasTipCap(ctx) to get suggested tip
	// - Handle errors from SuggestGasTipCap call
	// - IMPORTANT: Make a defensive copy (big.Int is mutable!)
	// - Why suggest? Let the node/network tell us a reasonable tip

	// TODO: Determine max fee cap
	// - If cfg.MaxFee is provided, use it
	// - Otherwise, calculate: 2 * baseFee + tipCap
	// - The "2x" rule: Accounts for base fee volatility (can increase 12.5% per block)
	// - IMPORTANT: Make defensive copies when using provided maxFee
	// - Why 2x? Provides buffer against base fee spikes between submission and inclusion

	// TODO: Prepare transaction data
	// - Set gasLimit from cfg.GasLimit (already validated/defaulted above)
	// - Make defensive copy of cfg.Data (byte slices are mutable!)
	// - Why copy data? Caller might mutate it after calling us

	// TODO: Construct DynamicFeeTx struct
	// - Create types.DynamicFeeTx with all fields:
	//   - ChainID: for replay protection
	//   - Nonce: transaction sequence number
	//   - GasTipCap: maximum priority fee willing to pay
	//   - GasFeeCap: maximum total fee willing to pay (base + tip)
	//   - Gas: gas limit for execution
	//   - To: recipient address (use &cfg.To to get pointer)
	//   - Value: amount of ETH to transfer (use defensive copy!)
	//   - Data: transaction payload (use copied data)
	// - Understanding: DynamicFeeTx is EIP-1559, different from LegacyTx

	// TODO: Wrap DynamicFeeTx in transaction envelope
	// - Call types.NewTx(txData) to create transaction
	// - This wraps the DynamicFeeTx in a types.Transaction envelope
	// - Why wrap? Transaction type is polymorphic (legacy, EIP-1559, EIP-2930)

	// TODO: Sign the transaction
	// - Create signer using types.LatestSignerForChainID(chainID)
	// - Call types.SignTx(tx, signer, cfg.PrivateKey)
	// - Handle errors from SignTx (invalid key, signature failure)
	// - Understanding: Signer determines signature format (EIP-155, EIP-2930, EIP-1559)

	// TODO: Send transaction (unless NoSend is set)
	// - Check if cfg.NoSend is false
	// - If sending, call client.SendTransaction(ctx, signedTx)
	// - Handle errors from SendTransaction (network issues, nonce too low, etc.)
	// - Why NoSend option? Allows testing without broadcasting to network

	// TODO: Construct and return Result
	// - Create Result struct with:
	//   - FromAddress: sender address (already derived)
	//   - Nonce: transaction nonce (already determined)
	//   - Tx: signed transaction (returned from SignTx)
	//   - BaseFee: base fee from header (already copied)
	// - Return Result and nil error

	return nil, errors.New("not implemented")
}
