package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/example/go-10x-minis/geth/06-eip1559/internal/eip1559"
)

/*
===================================================================
EIP-1559 Dynamic Fee Transaction Builder
===================================================================

This program demonstrates how to build and sign EIP-1559 (London upgrade)
transactions with dynamic fee estimation. EIP-1559 introduced a two-part
fee structure that makes gas fees more predictable.

USAGE:
    go run ./cmd/app/main.go <RPC_URL> <PRIVATE_KEY_HEX> <TO_ADDRESS> <AMOUNT_WEI>
    go run ./cmd/app/main.go <RPC_URL> <PRIVATE_KEY_HEX> <TO_ADDRESS> <AMOUNT_WEI> <demo>

Arguments:
    <RPC_URL>         - Ethereum RPC endpoint (required)
    <PRIVATE_KEY_HEX> - Private key without 0x prefix (required)
    <TO_ADDRESS>      - Recipient address (required)
    <AMOUNT_WEI>      - Amount in wei (required)
    [demo]            - Demo to run: all, simple, custom-fees, no-send (default: "all")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com deadbeef...cafe 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb 1000000000000000 simple
    go run ./cmd/app/main.go https://sepolia.gateway.tenderly.co privatekey... 0x123... 0 custom-fees
    go run ./cmd/app/main.go https://rpc.ankr.com/eth privatekey... 0x456... 1000 no-send

Common RPC Endpoints:
    Ethereum Mainnet:
      - https://eth.llamarpc.com
      - https://rpc.ankr.com/eth
      - https://eth-mainnet.public.blastapi.io

    Sepolia Testnet:
      - https://rpc.sepolia.org
      - https://sepolia.gateway.tenderly.co

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: cmd/app")
- Step Into (F11) to enter eip1559 package functions
- Watch Variables panel to see transaction data and fee calculations
- Inspect tx fields: ChainID, Nonce, GasTipCap, GasFeeCap

EIP-1559 KEY CONCEPTS:
- Base Fee: Algorithmically determined, burned (removed from supply)
- Priority Fee (Tip): Paid to validators, incentivizes inclusion
- Max Fee Cap: Maximum total (base + tip) willing to pay
- 2x Base Fee Rule: maxFee = 2 * baseFee + tip (buffer for ~6 blocks)
- Fee Market: More predictable than legacy auction-based fees

TRANSACTION LIFECYCLE:
1. Create transaction with fee caps
2. Sign with private key
3. Broadcast to network (enters mempool)
4. Validator includes in block (if fees sufficient)
5. Transaction mined (status: pending → confirmed)
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	// DEBUG: os.Args contains all command-line arguments
	// DEBUG: os.Args[0] is the program name
	// DEBUG: os.Args[1-4] are required arguments
	// DEBUG: os.Args[5] is optional demo type

	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> <PRIVATE_KEY_HEX> <TO_ADDRESS> <AMOUNT_WEI> [demo]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com <privatekey> 0x123... 1000000000000000 simple\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://sepolia.gateway.tenderly.co <privatekey> 0x456... 0 custom-fees\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, simple, custom-fees, no-send\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	privateKeyHex := os.Args[2]
	toAddressHex := os.Args[3]
	amountWeiStr := os.Args[4]

	demo := "all"
	if len(os.Args) > 5 {
		demo = os.Args[5]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL', 'privateKeyHex', 'toAddressHex', 'amountWeiStr' variables
	fmt.Printf("=== EIP-1559 Dynamic Fee Transaction Builder ===\n")
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Demo: %s\n\n", demo)

	// ============================================================================
	// STEP 2: Parse Private Key
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before parsing private key
	// DEBUG: Private key is ECDSA secp256k1 curve (same as Bitcoin)
	// DEBUG: 32 bytes (256 bits) of entropy
	//
	// Security warning: Never hardcode private keys in production code!
	// Use environment variables, key management systems, or hardware wallets.
	//
	// crypto.HexToECDSA: Parses hex string to ECDSA private key
	// The key contains both private and public components

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v\n", err)
	}

	// BREAKPOINT: Set breakpoint here after parsing private key
	// DEBUG: Inspect 'privateKey' structure
	// DEBUG: privateKey.PublicKey is the corresponding public key
	// DEBUG: Address will be derived from public key later
	fmt.Println("✓ Parsed private key")

	// ============================================================================
	// STEP 3: Parse Recipient Address
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before parsing address
	// DEBUG: Ethereum addresses are 20 bytes (40 hex characters)
	// DEBUG: Commonly prefixed with 0x (but not required by this parser)

	if !common.IsHexAddress(toAddressHex) {
		log.Fatalf("Invalid recipient address: %s\n", toAddressHex)
	}
	toAddress := common.HexToAddress(toAddressHex)

	// BREAKPOINT: Set breakpoint here after parsing address
	// DEBUG: Inspect 'toAddress' value
	// DEBUG: common.Address is a 20-byte array
	fmt.Printf("✓ Parsed recipient address: %s\n", toAddress.Hex())

	// ============================================================================
	// STEP 4: Parse Amount
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before parsing amount
	// DEBUG: Amount is in wei (1 ETH = 10^18 wei)
	// DEBUG: big.Int is required because wei amounts can be huge

	amountWei, ok := new(big.Int).SetString(amountWeiStr, 10)
	if !ok {
		log.Fatalf("Failed to parse amount: %s\n", amountWeiStr)
	}

	// BREAKPOINT: Set breakpoint here after parsing amount
	// DEBUG: Inspect 'amountWei' value
	// DEBUG: Convert to ETH: amountWei / 10^18
	fmt.Printf("✓ Parsed amount: %s wei\n", amountWei.String())
	fmt.Println()

	// ============================================================================
	// STEP 5: Connect to Ethereum RPC Endpoint
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before connection
	// DEBUG: This establishes the RPC connection
	// DEBUG: Watch for connection errors in the error variable
	//
	// Context with timeout: Prevents hanging indefinitely if endpoint is slow
	// ethclient.Dial: Standard go-ethereum client for JSON-RPC communication

	ctx := context.Background()
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// BREAKPOINT: Set breakpoint here
	// DEBUG: Step Over (F10) to execute the Dial call
	// DEBUG: If it fails, inspect 'err' to see the connection error
	client, err := ethclient.Dial(dialCtx, rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to RPC endpoint: %v\n", err)
	}
	defer client.Close()

	// BREAKPOINT: Set breakpoint here after successful connection
	// DEBUG: Connection is established—client is ready to use
	fmt.Println("✓ Connected to Ethereum RPC endpoint")
	fmt.Println()

	// ============================================================================
	// STEP 6: Example 1 - Simple Transfer with Auto Fees
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through first example
	// DEBUG: This demonstrates the most common use case—simple ETH transfer
	// DEBUG: Auto fees: Let the network suggest priority fee, use 2x base fee rule

	if demo == "all" || demo == "simple" {
		fmt.Println("=== Example 1: Simple Transfer with Auto Fees ===")
		fmt.Println("Building transaction with automatic fee estimation...")
		fmt.Println()

		// Create configuration for simple transfer
		// BREAKPOINT: Set breakpoint here
		// DEBUG: NoSend=true means we build and sign but don't broadcast
		// DEBUG: This is useful for testing and fee estimation
		cfg := eip1559.Config{
			PrivateKey: privateKey,
			To:         toAddress,
			AmountWei:  amountWei,
			NoSend:     true, // Don't actually send (for demo safety)
		}

		// BREAKPOINT: Set breakpoint here before calling eip1559.Run
		// DEBUG: Step Into (F11) to see the implementation in solution.reference.go
		// DEBUG: You'll see fee estimation, nonce lookup, tx building, and signing
		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		result, err := eip1559.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to build transaction: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful build
		// DEBUG: Expand 'result' in Variables panel to inspect all fields
		// DEBUG: Look at Tx.GasFeeCap(), Tx.GasTipCap(), Tx.Hash()
		// DEBUG: Examine BaseFee to see current network base fee

		fmt.Println("Transaction Built Successfully:")
		fmt.Printf("  From Address:    %s\n", result.FromAddress.Hex())
		fmt.Printf("  To Address:      %s\n", toAddress.Hex())
		fmt.Printf("  Amount:          %s wei\n", amountWei.String())
		fmt.Printf("  Nonce:           %d\n", result.Nonce)
		fmt.Printf("  Gas Limit:       %d\n", result.Tx.Gas())
		fmt.Printf("  Max Priority Fee: %s wei/gas\n", result.Tx.GasTipCap().String())
		fmt.Printf("  Max Fee Cap:     %s wei/gas\n", result.Tx.GasFeeCap().String())
		fmt.Printf("  Base Fee:        %s wei/gas\n", result.BaseFee.String())
		fmt.Printf("  Tx Hash:         %s\n", result.Tx.Hash().Hex())
		fmt.Println()

		// Calculate estimated cost
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Estimated cost = gas * (baseFee + tip)
		// DEBUG: Actual cost might be lower if base fee decreases
		estimatedGasFee := new(big.Int).Mul(
			new(big.Int).Add(result.BaseFee, result.Tx.GasTipCap()),
			big.NewInt(int64(result.Tx.Gas())),
		)
		fmt.Println("Fee Breakdown:")
		fmt.Printf("  Current Base Fee:     %s wei/gas\n", result.BaseFee.String())
		fmt.Printf("  Priority Fee (Tip):   %s wei/gas\n", result.Tx.GasTipCap().String())
		fmt.Printf("  Estimated Gas Fee:    %s wei (base+tip)*gas\n", estimatedGasFee.String())
		fmt.Printf("  Maximum Gas Fee:      %s wei (if base fee spikes)\n",
			new(big.Int).Mul(result.Tx.GasFeeCap(), big.NewInt(int64(result.Tx.Gas()))).String())
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice how transaction is built and signed but not sent
		// DEBUG: NoSend=true allows testing without spending real ETH
		fmt.Println("ℹ Transaction not sent (NoSend=true)")
		fmt.Println()
	}

	// ============================================================================
	// STEP 7: Example 2 - Custom Fee Parameters
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through custom fee example
	// DEBUG: This shows how to override automatic fee estimation
	// DEBUG: Useful for urgent transactions (high fees) or patient ones (low fees)

	if demo == "all" || demo == "custom-fees" {
		fmt.Println("=== Example 2: Custom Fee Parameters ===")
		fmt.Println("Building transaction with custom priority fee and max fee...")
		fmt.Println()

		// Create configuration with custom fees
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice we specify MaxPriorityFee and MaxFee explicitly
		// DEBUG: This overrides automatic fee estimation
		customPriorityFee := big.NewInt(2_000_000_000) // 2 gwei
		customMaxFee := big.NewInt(50_000_000_000)     // 50 gwei

		cfg := eip1559.Config{
			PrivateKey:     privateKey,
			To:             toAddress,
			AmountWei:      big.NewInt(0), // Zero value transfer
			MaxPriorityFee: customPriorityFee,
			MaxFee:         customMaxFee,
			NoSend:         true,
		}

		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		// BREAKPOINT: Set breakpoint here before calling eip1559.Run
		// DEBUG: Step Into (F11) to see how custom fees are handled
		result, err := eip1559.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to build transaction with custom fees: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful build
		// DEBUG: Compare custom fees with current base fee
		// DEBUG: If maxFee < baseFee, transaction will never be included!

		fmt.Println("Transaction with Custom Fees:")
		fmt.Printf("  Max Priority Fee: %s wei/gas (custom: 2 gwei)\n", result.Tx.GasTipCap().String())
		fmt.Printf("  Max Fee Cap:     %s wei/gas (custom: 50 gwei)\n", result.Tx.GasFeeCap().String())
		fmt.Printf("  Current Base Fee: %s wei/gas\n", result.BaseFee.String())
		fmt.Println()

		// Fee analysis
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Check if custom fees are sufficient
		baseFeePlus := new(big.Int).Add(result.BaseFee, result.Tx.GasTipCap())
		if result.Tx.GasFeeCap().Cmp(baseFeePlus) < 0 {
			fmt.Println("⚠ WARNING: Max fee cap is less than base fee + tip!")
			fmt.Println("  Transaction may not be included if base fee rises.")
		} else {
			fmt.Println("✓ Custom fees are sufficient for current network conditions")
		}
		fmt.Println()

		fmt.Println("ℹ Transaction not sent (NoSend=true)")
		fmt.Println()
	}

	// ============================================================================
	// STEP 8: Example 3 - Build Without Sending
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through no-send example
	// DEBUG: This demonstrates offline signing and transaction inspection
	// DEBUG: Useful for air-gapped signing, batching, or debugging

	if demo == "all" || demo == "no-send" {
		fmt.Println("=== Example 3: Build and Sign Without Sending ===")
		fmt.Println("Creating transaction for offline signing or inspection...")
		fmt.Println()

		// Create configuration with NoSend flag
		// BREAKPOINT: Set breakpoint here
		// DEBUG: NoSend=true allows transaction creation without broadcasting
		// DEBUG: Useful for: offline signing, transaction inspection, fee estimation
		cfg := eip1559.Config{
			PrivateKey: privateKey,
			To:         toAddress,
			AmountWei:  amountWei,
			NoSend:     true, // Build and sign, but don't broadcast
		}

		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		result, err := eip1559.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to build transaction: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Transaction is fully built and signed
		// DEBUG: It can be broadcast later using SendTransaction RPC

		fmt.Println("Transaction Ready (Not Broadcasted):")
		fmt.Printf("  Tx Hash: %s\n", result.Tx.Hash().Hex())
		fmt.Printf("  Nonce:   %d\n", result.Nonce)
		fmt.Println()

		// Display raw transaction for manual broadcasting
		// BREAKPOINT: Set breakpoint here
		// DEBUG: This hex can be broadcast using eth_sendRawTransaction
		rawTx, err := result.Tx.MarshalBinary()
		if err != nil {
			log.Fatalf("Failed to marshal transaction: %v\n", err)
		}

		fmt.Println("Raw Transaction (for manual broadcasting):")
		fmt.Printf("  %x\n", rawTx)
		fmt.Println()

		fmt.Println("Use Case: Offline Signing")
		fmt.Println("  1. Build and sign transaction on air-gapped machine")
		fmt.Println("  2. Transfer signed transaction to online machine")
		fmt.Println("  3. Broadcast using eth_sendRawTransaction")
		fmt.Println("  4. This keeps private key offline for security")
		fmt.Println()
	}

	// ============================================================================
	// Understanding EIP-1559 Fee Mechanism
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand fee mechanism
	// DEBUG: This section explains how EIP-1559 fees work

	if demo == "all" {
		fmt.Println("=== Understanding EIP-1559 Fee Mechanism ===")
		fmt.Println()

		fmt.Println("Base Fee Algorithm:")
		fmt.Println("  - Adjusts based on block fullness (target: 50% full)")
		fmt.Println("  - If blocks >50% full: base fee increases 12.5%")
		fmt.Println("  - If blocks <50% full: base fee decreases 12.5%")
		fmt.Println("  - Creates feedback loop for predictable fees")
		fmt.Println()

		fmt.Println("Fee Structure:")
		fmt.Println("  - Base Fee: Burned (removed from ETH supply)")
		fmt.Println("  - Priority Fee: Paid to validators (incentive)")
		fmt.Println("  - Max Fee Cap: Your budget limit")
		fmt.Println()

		fmt.Println("Fee Calculation:")
		fmt.Println("  - Actual fee = min(maxFee, baseFee + priorityFee)")
		fmt.Println("  - Refund = maxFee - actualFee")
		fmt.Println("  - You never pay more than maxFee")
		fmt.Println()

		fmt.Println("2x Base Fee Rule of Thumb:")
		fmt.Println("  - maxFee = 2 * baseFee + tip")
		fmt.Println("  - Provides buffer for ~6 blocks of base fee increases")
		fmt.Println("  - Each block can increase base fee by max 12.5%")
		fmt.Println("  - 1.125^6 ≈ 2.03 (approximately doubles)")
		fmt.Println()

		fmt.Println("Why EIP-1559 is Better:")
		fmt.Println("  1. More predictable fees (no blind auction)")
		fmt.Println("  2. Fee burning reduces ETH supply (deflationary)")
		fmt.Println("  3. Automatic fee adjustment (no manual estimation)")
		fmt.Println("  4. Better UX (wallets can estimate accurately)")
		fmt.Println()
	}

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter eip1559 package functions
	// 5. Watch Variables panel to see transaction data
	// 6. Inspect result.Tx to see all transaction fields
	// 7. Use Call Stack panel to see function hierarchy
	//
	// EIP-1559 CONCEPTS TO EXPLORE:
	// - Base Fee vs Priority Fee: How do they differ?
	// - Max Fee Cap: Why set it higher than base fee?
	// - Fee Burning: How does it affect ETH supply?
	// - Dynamic Fee Adjustment: How does the algorithm work?
	// - Nonce Management: How to handle parallel transactions?
	//
	// STEPPING INTO EIP1559 PACKAGE:
	// - Set breakpoint at eip1559.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to solution.reference.go
	// - Step through to see: nonce lookup, fee estimation, tx building, signing
	// - Watch how errors are handled and data is validated
	// - See defensive copying in action with big.Int values
	//
	// COMMON ERRORS TO DEBUG:
	// - "nonce too low": Already sent transaction with this nonce
	// - "insufficient funds": Not enough ETH for value + fees
	// - "max fee per gas less than block base fee": Your maxFee is too low
	// - "gas limit too low": Transaction needs more gas
	// - "invalid signature": Private key doesn't match sender address
}
