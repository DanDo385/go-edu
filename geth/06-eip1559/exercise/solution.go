//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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
	// ============================================================================
	// STEP 1: Input Validation - Defensive Programming Pattern
	// ============================================================================
	// Why validate? This function is a library API. We can't trust callers to
	// always pass valid inputs. This is defensive programming - assume inputs
	// might be invalid and handle gracefully.
	//
	// Context handling: Same pattern as module 01-stack. If ctx is nil, provide
	// a safe default. Context.Background() creates a non-cancelable context that
	// never times out. This ensures RPC calls won't panic from nil context.
	//
	// This pattern repeats: Every function accepting context should validate it.
	if ctx == nil {
		ctx = context.Background()
	}

	// Client validation: Same pattern as module 01-stack. The RPCClient interface
	// is our dependency. If nil, we can't make RPC calls. Fail fast with a
	// descriptive error rather than panicking later.
	//
	// This pattern repeats: Always validate interface dependencies before use.
	if client == nil {
		return nil, errors.New("client is nil")
	}

	// Private key validation: Unlike module 01-stack (read-only operations), this
	// module creates and signs transactions. Without a private key, we can't sign.
	// This is a critical validation - transactions without signatures are rejected.
	//
	// New concept: Cryptographic signing. Every Ethereum transaction must be
	// signed with the sender's private key to prove authenticity and prevent
	// impersonation. The signature also commits to the transaction data, making
	// it tamper-proof.
	if cfg.PrivateKey == nil {
		return nil, errors.New("private key is required")
	}

	// Amount validation with sensible default: If AmountWei is nil, default to
	// zero. This allows creating transactions that interact with smart contracts
	// without transferring ETH (e.g., calling a function).
	//
	// Why zero as default? Many transactions don't transfer value - they just
	// call contract functions. Zero is a safe, non-destructive default.
	//
	// This pattern repeats: Use zero values as sensible defaults for optional
	// numeric fields. This makes the API more ergonomic - callers don't need to
	// specify zero explicitly.
	if cfg.AmountWei == nil {
		cfg.AmountWei = big.NewInt(0)
	}

	// Gas limit validation with sensible default: 21000 is the base cost of a
	// simple ETH transfer (no data, no contract calls). This is the minimum
	// gas for any valid transaction.
	//
	// Why 21000? Ethereum's gas pricing:
	//   - 21000 base transaction cost
	//   - + 68 gas per non-zero data byte
	//   - + 4 gas per zero data byte
	//   - + contract execution costs (if calling contract)
	//
	// For simple transfers (no data), 21000 is exact. For contract calls,
	// callers must specify higher limits.
	//
	// This pattern repeats: Provide sensible defaults that work for common cases.
	if cfg.GasLimit == 0 {
		cfg.GasLimit = defaultDynamicGasLimit
	}

	// ============================================================================
	// STEP 2: Derive Sender Address from Private Key
	// ============================================================================
	// Cryptographic address derivation: The sender address is derived from the
	// private key's public key using Keccak-256 hash. This is a fundamental
	// Ethereum pattern.
	//
	// The derivation: privateKey → publicKey → keccak256(publicKey) → last 20 bytes
	//
	// Why derive instead of accepting as parameter? Security! If we accepted an
	// address parameter, callers could provide any address, but they wouldn't be
	// able to sign transactions from that address (they don't have the key).
	// By deriving from the key, we ensure consistency - the signer owns the address.
	//
	// crypto.PubkeyToAddress: This function from go-ethereum performs the full
	// derivation. It's the standard way to get an address from a key.
	//
	// This pattern repeats: Always derive addresses from keys. Never trust
	// user-provided addresses for signing operations.
	from := crypto.PubkeyToAddress(cfg.PrivateKey.PublicKey)

	// ============================================================================
	// STEP 3: Determine Transaction Nonce
	// ============================================================================
	// Nonce: A sequence number for transactions from an account. Every transaction
	// from an address must have a unique, incrementing nonce. This prevents replay
	// attacks and ensures transactions are processed in order.
	//
	// Two modes:
	//   1. Manual nonce (cfg.Nonce != nil): Caller controls nonce. Useful for:
	//      - Advanced nonce management (parallel transactions)
	//      - Testing (deterministic nonces)
	//      - Replacing stuck transactions (same nonce, higher gas price)
	//
	//   2. Automatic nonce (cfg.Nonce == nil): Query network for next nonce.
	//      Most common case - let the network tell us what nonce to use.
	//
	// PendingNonceAt vs NonceAt:
	//   - NonceAt: Nonce of last MINED transaction (on-chain)
	//   - PendingNonceAt: Nonce of last SUBMITTED transaction (includes mempool)
	//
	// Why pending? If you just submitted a transaction (nonce N) and immediately
	// create another, you need nonce N+1. NonceAt would return N (last mined),
	// causing a collision. PendingNonceAt returns N+1 (correct next nonce).
	//
	// This pattern repeats: Use PendingNonceAt for transaction creation. Use
	// NonceAt only for historical queries (what was the nonce at block X?).
	var nonce uint64
	var err error
	if cfg.Nonce != nil {
		// Manual nonce: Use caller-provided value
		nonce = *cfg.Nonce
	} else {
		// Automatic nonce: Query network
		nonce, err = client.PendingNonceAt(ctx, from)
		if err != nil {
			// Error wrapping: Preserve error chain for debugging. The %w verb
			// wraps the error, allowing errors.Is() and errors.Unwrap() to work.
			return nil, fmt.Errorf("pending nonce: %w", err)
		}
	}

	// ============================================================================
	// STEP 4: Retrieve Chain ID - Replay Protection
	// ============================================================================
	// Chain ID: Unique identifier for the blockchain network. Introduced in
	// EIP-155 to prevent replay attacks (transaction signed for mainnet being
	// replayed on a testnet).
	//
	// How replay protection works: Chain ID is included in the transaction signature.
	// If someone copies your signed transaction and broadcasts it on a different
	// network, the signature verification fails because the chain ID doesn't match.
	//
	// Same pattern as module 01-stack: Call RPC method, check error, validate
	// nil response. This is the standard pattern for all RPC calls.
	//
	// Building on module 01-stack: There we fetched chain ID to prove connectivity.
	// Here we fetch it to build a transaction signature. Same RPC call, different
	// purpose. This demonstrates how fundamental operations (fetching chain ID)
	// are reused across different use cases.
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}
	if chainID == nil {
		return nil, errors.New("chain id was nil")
	}

	// ============================================================================
	// STEP 5: Fetch Block Header to Get Base Fee
	// ============================================================================
	// EIP-1559 base fee: The algorithmically determined "base price" for block
	// inclusion. This is the foundational concept of EIP-1559 - fees are no longer
	// a pure auction; there's a predictable base fee that adjusts based on demand.
	//
	// Base fee algorithm: If blocks are >50% full, base fee increases by 12.5%.
	// If blocks are <50% full, base fee decreases by 12.5%. This creates a
	// feedback loop that keeps blocks approximately 50% full.
	//
	// Why fetch header? The base fee is stored in the block header (header.BaseFee).
	// We need to know the current base fee to construct our transaction's fee caps.
	//
	// cfg.BlockNumber: If nil, fetches latest block. If specified, fetches that
	// historical block. Why allow specifying? For testing and fee estimation from
	// specific points in time.
	//
	// Defensive programming: We validate both header != nil AND header.BaseFee != nil.
	// BaseFee is nil for pre-London (pre-EIP-1559) blocks. If we try to create an
	// EIP-1559 transaction on a pre-London network, we must fail with a clear error.
	header, err := client.HeaderByNumber(ctx, cfg.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("header by number: %w", err)
	}
	if header == nil || header.BaseFee == nil {
		// Descriptive error: Tell the user exactly why we failed. "BaseFee unavailable"
		// indicates a pre-London block. "Upgrade to London block" hints at the solution.
		return nil, errors.New("base fee unavailable (upgrade to London block)")
	}

	// CRITICAL: Defensive copy of base fee!
	//
	// Why copy? header.BaseFee is a *big.Int pointer returned by the client.
	// If we return it directly in our Result, callers could mutate it:
	//   result.BaseFee.Add(result.BaseFee, big.NewInt(1000))
	// This would mutate the client's internal data, affecting future calls!
	//
	// new(big.Int).Set(header.BaseFee): Creates a new big.Int and copies the value.
	// Now our Result owns an independent copy. Callers can mutate it without
	// affecting the client.
	//
	// This pattern repeats: Always copy mutable return values (big.Int, slices,
	// maps) when returning data from external libraries or shared state.
	baseFee := new(big.Int).Set(header.BaseFee)

	// ============================================================================
	// STEP 6: Determine Max Priority Fee (Tip Cap)
	// ============================================================================
	// Priority fee (tip): The amount paid to validators for including the transaction.
	// Higher tips incentivize faster inclusion. Tips go to validators, not burned.
	//
	// Two modes:
	//   1. Manual tip (cfg.MaxPriorityFee != nil): Caller specifies exact tip.
	//      Useful for:
	//        - Testing (deterministic fees)
	//        - Advanced fee strategies (willing to wait for low fees)
	//        - Urgent transactions (high tip for immediate inclusion)
	//
	//   2. Automatic tip (cfg.MaxPriorityFee == nil): Ask network for suggested tip.
	//      Most common case - let the node tell us what tip is reasonable.
	//
	// SuggestGasTipCap: RPC method that returns a suggested priority fee based on
	// recent blocks. Different nodes may use different algorithms:
	//   - Geth: Looks at recent tips, returns percentile (e.g., 60th percentile)
	//   - Other clients: May use different heuristics
	//
	// The suggestion is just a hint - you're not required to use it. But it's
	// usually reasonable for "normal" priority.
	tipCap := cfg.MaxPriorityFee
	if tipCap == nil {
		tipCap, err = client.SuggestGasTipCap(ctx)
		if err != nil {
			return nil, fmt.Errorf("suggest gas tip cap: %w", err)
		}
	}

	// CRITICAL: Defensive copy of tip cap!
	//
	// Why copy even manual tips? If the caller passes cfg.MaxPriorityFee, it's
	// their *big.Int. If we use it directly, we might mutate their data. By
	// copying, we ensure independence.
	//
	// This pattern repeats: Always copy big.Int values, even when provided by
	// caller. Immutability prevents subtle bugs.
	tipCap = new(big.Int).Set(tipCap)

	// ============================================================================
	// STEP 7: Determine Max Fee Cap
	// ============================================================================
	// Max fee cap: The maximum total (base fee + tip) willing to pay. This is your
	// budget cap - you won't pay more than this, even if base fee spikes.
	//
	// Two modes:
	//   1. Manual cap (cfg.MaxFee != nil): Caller specifies exact budget.
	//   2. Automatic cap (cfg.MaxFee == nil): Calculate using "2x base fee + tip" rule.
	//
	// The "2x base fee" rule of thumb:
	//   - Base fee can increase max 12.5% per block (EIP-1559 design)
	//   - Over 6 blocks, base fee could approximately double (1.125^6 ≈ 2.03)
	//   - Setting maxFee = 2 * baseFee + tip gives buffer for ~6 blocks of increases
	//   - If your tx isn't included in 6 blocks, it probably won't be (too low tip)
	//
	// Why this matters: If you set maxFee too low, your transaction may never be
	// included (base fee rises above your cap). If you set it too high, you might
	// overpay (though you'll get refunded, it ties up capital).
	//
	// Actual fee paid: min(maxFee, baseFee + tip). So even if you set a high
	// maxFee, you only pay the actual baseFee + tip. The extra is refunded.
	maxFee := cfg.MaxFee
	if maxFee == nil {
		// Calculate: 2 * baseFee + tipCap
		//
		// Step 1: Multiply baseFee by 2
		twoBase := new(big.Int).Mul(baseFee, big.NewInt(2))
		//
		// Step 2: Add tipCap
		maxFee = new(big.Int).Add(twoBase, tipCap)
		//
		// Note: We use new(big.Int) for intermediate results to avoid mutating
		// baseFee or tipCap. big.Int operations often mutate the receiver!
	} else {
		// Manual maxFee: Make defensive copy
		maxFee = new(big.Int).Set(maxFee)
	}

	// ============================================================================
	// STEP 8: Prepare Transaction Data
	// ============================================================================
	// Gas limit: Already validated and defaulted in Step 1. This is the maximum
	// amount of gas (computation) the transaction is allowed to consume.
	//
	// If execution uses less gas, you get refunded for unused gas. If execution
	// needs more gas than the limit, the transaction reverts (but you still pay
	// for gas used up to the limit).
	gasLimit := cfg.GasLimit

	// Defensive copy of transaction data (payload):
	//
	// Why copy byte slices? Slices in Go are references to underlying arrays.
	// If we use cfg.Data directly, the caller could modify the slice after calling
	// us. This could cause race conditions or unexpected behavior.
	//
	// append([]byte(nil), cfg.Data...): Idiom for copying a byte slice.
	//   - []byte(nil) creates an empty slice
	//   - append(..., cfg.Data...) appends all elements, creating a new backing array
	//   - Result: Independent copy of the data
	//
	// This pattern repeats: Always copy mutable data structures (slices, maps)
	// when storing or returning them. Immutability prevents bugs.
	dataCopy := append([]byte(nil), cfg.Data...)

	// ============================================================================
	// STEP 9: Construct DynamicFeeTx Struct
	// ============================================================================
	// DynamicFeeTx: The EIP-1559 transaction type. This is different from
	// LegacyTx (pre-EIP-1559) which only had a single gasPrice field.
	//
	// Fields explained:
	//   - ChainID: Replay protection (EIP-155). Prevents transactions signed for
	//     mainnet from being replayed on other networks.
	//
	//   - Nonce: Sequence number. Each transaction from an account must have a
	//     unique, incrementing nonce. Prevents replay attacks and ensures ordering.
	//
	//   - GasTipCap: Maximum priority fee (tip) per gas willing to pay. This
	//     incentivizes validators to include your transaction. Higher tips = faster.
	//
	//   - GasFeeCap: Maximum total fee (base + tip) per gas willing to pay. This
	//     is your budget cap. You won't pay more than this.
	//
	//   - Gas: Gas limit. Maximum computation allowed. If execution exceeds this,
	//     the transaction reverts (but you still pay for gas used).
	//
	//   - To: Recipient address. For ETH transfers, this is who receives the ETH.
	//     For contract calls, this is the contract address. We use &cfg.To to get
	//     a pointer (required by the struct).
	//
	//   - Value: Amount of ETH to transfer (in wei). Can be zero for contract calls
	//     that don't transfer value. We use new(big.Int).Set() for defensive copy.
	//
	//   - Data: Transaction payload. For ETH transfers, this is empty. For contract
	//     calls, this is the encoded function call (ABI-encoded). We use our copied
	//     data to ensure immutability.
	//
	// Why all these fields? Each serves a specific purpose in Ethereum's transaction
	// model. Together they define a complete, self-contained unit of work.
	txData := &types.DynamicFeeTx{
		ChainID:   chainID,                      // Replay protection
		Nonce:     nonce,                        // Sequence number
		GasTipCap: tipCap,                       // Max priority fee
		GasFeeCap: maxFee,                       // Max total fee
		Gas:       gasLimit,                     // Computation budget
		To:        &cfg.To,                      // Recipient
		Value:     new(big.Int).Set(cfg.AmountWei), // Amount (defensive copy!)
		Data:      dataCopy,                     // Payload (already copied)
	}

	// ============================================================================
	// STEP 10: Wrap in Transaction Envelope
	// ============================================================================
	// Transaction envelope: types.Transaction is a polymorphic type that can
	// contain different transaction types:
	//   - LegacyTx (pre-EIP-155): No replay protection
	//   - AccessListTx (EIP-2930): Replay protection + access lists
	//   - DynamicFeeTx (EIP-1559): Replay protection + dynamic fees (this module)
	//
	// types.NewTx wraps our DynamicFeeTx in a Transaction envelope. The envelope
	// handles type discrimination and provides common methods (Hash, Size, etc.).
	//
	// Why polymorphism? Ethereum supports multiple transaction types for backwards
	// compatibility. Old clients understand LegacyTx, new clients understand all
	// types. The envelope pattern allows code to work with any transaction type.
	//
	// This pattern repeats: Use type-safe wrappers (envelopes) for polymorphic
	// data. Go's interfaces + type switches provide type safety while allowing
	// runtime flexibility.
	tx := types.NewTx(txData)

	// ============================================================================
	// STEP 11: Sign the Transaction
	// ============================================================================
	// Transaction signing: Cryptographically proves the sender authorized this
	// transaction. Without a valid signature, nodes reject the transaction.
	//
	// The signature commits to:
	//   - All transaction fields (nonce, to, value, data, fees, etc.)
	//   - Chain ID (prevents replay across networks)
	//
	// If any field is modified after signing, the signature becomes invalid.
	// This makes transactions tamper-proof.
	//
	// Signer: Determines signature format. Different transaction types use
	// different signature formats:
	//   - LegacyTx: EIP-155 signature format
	//   - AccessListTx: EIP-2930 signature format
	//   - DynamicFeeTx: EIP-1559 signature format (same as EIP-2930)
	//
	// types.LatestSignerForChainID: Creates a signer using the latest signature
	// format for the given chain ID. This is usually what you want - use the
	// most modern format.
	//
	// types.SignTx: Performs the actual signing. Takes unsigned tx, signer, and
	// private key. Returns signed tx (new instance, original unchanged).
	//
	// Error handling: Signing can fail if the private key is invalid or if
	// there's a cryptographic error. Always check the error.
	signer := types.LatestSignerForChainID(chainID)
	signedTx, err := types.SignTx(tx, signer, cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("sign tx: %w", err)
	}

	// ============================================================================
	// STEP 12: Send Transaction to Network (Optional)
	// ============================================================================
	// Transaction broadcasting: Submits the signed transaction to the network.
	// The node adds it to the mempool (pending transactions) and gossips it to
	// other nodes. Eventually, a validator includes it in a block.
	//
	// cfg.NoSend: Allows creating and signing transactions without broadcasting.
	// Why useful?
	//   - Testing: Verify transaction construction without actually sending
	//   - Offline signing: Sign on air-gapped machine, broadcast separately
	//   - Batching: Create multiple transactions, then send all at once
	//   - Fee estimation: Build transaction to estimate gas, don't send
	//
	// SendTransaction: RPC call that broadcasts the transaction. Returns
	// immediately - it doesn't wait for the transaction to be mined. You'll
	// need to poll for the receipt to know when it's mined.
	//
	// Common errors:
	//   - "nonce too low": You already sent a transaction with this nonce
	//   - "insufficient funds": Not enough ETH to cover value + fees
	//   - "gas limit too low": Transaction execution needs more gas
	//   - "max fee per gas less than block base fee": Your maxFee is too low
	//
	// Error wrapping: We wrap the error to add context. When debugging, seeing
	// "send tx: nonce too low" is more helpful than just "nonce too low".
	if !cfg.NoSend {
		if err := client.SendTransaction(ctx, signedTx); err != nil {
			return nil, fmt.Errorf("send tx: %w", err)
		}
	}

	// ============================================================================
	// STEP 13: Construct and Return Result
	// ============================================================================
	// Result: Package useful transaction metadata for the caller. This gives
	// them everything they need to track the transaction and display information.
	//
	// Fields:
	//   - FromAddress: Sender address (derived from private key in Step 2)
	//   - Nonce: Transaction nonce (determined in Step 3)
	//   - Tx: Signed transaction (from Step 11)
	//   - BaseFee: Base fee from header (copied in Step 5)
	//
	// Why return this info? Callers need it to:
	//   - Display "Transaction sent from X with nonce Y"
	//   - Track transaction status (poll for receipt using tx hash)
	//   - Understand fee context (baseFee at time of submission)
	//   - Manage nonces for subsequent transactions
	//
	// Pattern: Return structured data that provides context, not just success/failure.
	// This makes the API more useful and debuggable.
	//
	// Building on previous concepts:
	//   - We validated all inputs (Step 1) → now we return validated results
	//   - We handled errors consistently throughout → now we return success
	//   - We used defensive copying for mutable data → our Result is safe to use
	//   - We followed the same RPC patterns from module 01-stack → consistency!
	return &Result{
		FromAddress: from,      // Sender (derived from key)
		Nonce:       nonce,     // Sequence number
		Tx:          signedTx,  // Signed transaction (ready for tracking)
		BaseFee:     baseFee,   // Base fee context (already defensively copied)
	}, nil
}
