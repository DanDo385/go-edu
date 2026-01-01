package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/example/go-10x-minis/geth/08-abigen/internal/abigen"
)

/*
===================================================================
ERC20 Token Query with BoundContract (ABI Bindings)
===================================================================

This program demonstrates how to query ERC20 token metadata using
BoundContract for type-safe contract calls. This is cleaner and
safer than manual ABI encoding (module 07).

USAGE:
    go run ./cmd/app/main.go <RPC_URL> <TOKEN_ADDRESS>
    go run ./cmd/app/main.go <RPC_URL> <TOKEN_ADDRESS> <demo>
    go run ./cmd/app/main.go <RPC_URL> <TOKEN_ADDRESS> <HOLDER_ADDRESS> <demo>

Arguments:
    <RPC_URL>         - Ethereum RPC endpoint (required)
    <TOKEN_ADDRESS>   - ERC20 token contract address (required)
    [HOLDER_ADDRESS]  - Address to check balance (optional)
    [demo]            - Demo to run: all, basic, with-balance, historical (default: "all")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
    go run ./cmd/app/main.go https://eth.llamarpc.com 0x6B175474E89094C44Da98b954EedeAC495271d0F 0x123... with-balance
    go run ./cmd/app/main.go https://sepolia.gateway.tenderly.co 0x... historical

Common ERC20 Tokens (Mainnet):
    USDC:  0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
    DAI:   0x6B175474E89094C44Da98b954EedeAC495271d0F
    WETH:  0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2
    USDT:  0xdAC17F958D2ee523a2206206994597C13D831ec7

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
- Step Into (F11) to enter abigen package functions
- Watch Variables panel to see BoundContract operations
- Inspect how ABI handles encoding/decoding automatically

KEY CONCEPTS:
- BoundContract: Adapter pattern for type-safe contract calls
- ABI Package: Automatic encoding/decoding based on ABI definition
- CallOpts: Configuration for view function calls (context, block number)
- balanceOf: Example of function with parameters

COMPARISON WITH MODULE 07:
- Module 07: Manual ABI encoding/decoding (full control, more code)
- Module 08: BoundContract (cleaner, safer, less boilerplate)
- Module 09+: Will use abigen-generated typed bindings (compile-time safety)
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	// DEBUG: os.Args contains all command-line arguments
	// DEBUG: We need RPC URL, token address, optional holder address

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> <TOKEN_ADDRESS> [HOLDER_ADDRESS] [demo]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com 0x6B175474E89094C44Da98b954EedeAC495271d0F 0x123... with-balance\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, basic, with-balance, historical\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	tokenAddressHex := os.Args[2]

	var holderAddress *common.Address
	demo := "all"

	// Parse optional holder address and demo
	if len(os.Args) > 3 {
		arg3 := os.Args[3]
		if common.IsHexAddress(arg3) {
			addr := common.HexToAddress(arg3)
			holderAddress = &addr
			if len(os.Args) > 4 {
				demo = os.Args[4]
			}
		} else {
			demo = arg3
		}
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL', 'tokenAddressHex', 'holderAddress', 'demo' variables
	fmt.Printf("=== ERC20 Token Query with BoundContract ===\n")
	fmt.Printf("RPC URL: %s\n", rpcURL)
	if holderAddress != nil {
		fmt.Printf("Holder: %s\n", holderAddress.Hex())
	}
	fmt.Printf("Demo: %s\n\n", demo)

	// ============================================================================
	// STEP 2: Parse Token Contract Address
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before parsing address

	if !common.IsHexAddress(tokenAddressHex) {
		log.Fatalf("Invalid token address: %s\n", tokenAddressHex)
	}
	tokenAddress := common.HexToAddress(tokenAddressHex)

	// BREAKPOINT: Set breakpoint here after parsing address
	fmt.Printf("✓ Parsed token address: %s\n", tokenAddress.Hex())
	fmt.Println()

	// ============================================================================
	// STEP 3: Connect to Ethereum RPC Endpoint
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before connection

	ctx := context.Background()
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// BREAKPOINT: Set breakpoint here
	client, err := ethclient.Dial(dialCtx, rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to RPC endpoint: %v\n", err)
	}
	defer client.Close()

	// BREAKPOINT: Set breakpoint here after successful connection
	fmt.Println("✓ Connected to Ethereum RPC endpoint")
	fmt.Println()

	// ============================================================================
	// STEP 4: Example 1 - Basic Metadata Query
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through first example
	// DEBUG: This demonstrates basic token metadata query
	// DEBUG: No holder address, no balance query

	if demo == "all" || demo == "basic" {
		fmt.Println("=== Example 1: Basic Metadata Query ===")
		fmt.Println("Querying token metadata using BoundContract...")
		fmt.Println()

		// Create configuration without holder
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Holder is nil, so no balance query
		cfg := abigen.Config{
			Contract:    tokenAddress,
			BlockNumber: nil, // nil = latest block
			Holder:      nil, // nil = no balance query
		}

		// BREAKPOINT: Set breakpoint here before calling abigen.Run
		// DEBUG: Step Into (F11) to see BoundContract usage
		// DEBUG: You'll see ABI parsing, contract binding, and automatic encoding/decoding
		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		result, err := abigen.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to query token metadata: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful query
		// DEBUG: Expand 'result' in Variables panel
		// DEBUG: Notice Balance is nil (not queried)

		fmt.Println("Token Metadata:")
		fmt.Printf("  Name:         %s\n", result.Name)
		fmt.Printf("  Symbol:       %s\n", result.Symbol)
		fmt.Printf("  Decimals:     %d\n", result.Decimals)
		fmt.Printf("  Total Supply: %s (raw)\n", result.TotalSupply.String())
		fmt.Println()

		// Convert to human-readable format
		// BREAKPOINT: Set breakpoint here
		divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(result.Decimals)), nil)
		supplyFloat := new(big.Float).SetInt(result.TotalSupply)
		divisorFloat := new(big.Float).SetInt(divisor)
		humanReadable := new(big.Float).Quo(supplyFloat, divisorFloat)

		fmt.Println("Human-Readable Total Supply:")
		fmt.Printf("  %s %s\n", humanReadable.Text('f', int(result.Decimals)), result.Symbol)
		fmt.Println()

		fmt.Println("ℹ Compared to Module 07 (manual ABI):")
		fmt.Println("  - No manual selector computation")
		fmt.Println("  - No manual decoding logic")
		fmt.Println("  - Cleaner, more maintainable code")
		fmt.Println()
	}

	// ============================================================================
	// STEP 5: Example 2 - Query with Balance
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through balance query
	// DEBUG: This demonstrates querying balance for a specific address
	// DEBUG: Shows how BoundContract handles function parameters

	if demo == "all" || demo == "with-balance" {
		if holderAddress == nil {
			// Use a well-known address for demo (Vitalik's address)
			vitalik := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
			holderAddress = &vitalik
			fmt.Printf("No holder specified, using Vitalik's address for demo: %s\n\n", holderAddress.Hex())
		}

		fmt.Println("=== Example 2: Query with Balance ===")
		fmt.Printf("Querying balance for holder: %s...\n", holderAddress.Hex())
		fmt.Println()

		// Create configuration with holder
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Holder is set, so balanceOf will be queried
		cfg := abigen.Config{
			Contract:    tokenAddress,
			BlockNumber: nil,
			Holder:      holderAddress,
		}

		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		// BREAKPOINT: Set breakpoint here before calling abigen.Run
		// DEBUG: Step Into (F11) to see how balanceOf is called
		// DEBUG: Watch how address parameter is automatically encoded
		result, err := abigen.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to query token with balance: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful query
		// DEBUG: Notice Balance is now set (not nil)

		fmt.Println("Token Metadata:")
		fmt.Printf("  Name:         %s\n", result.Name)
		fmt.Printf("  Symbol:       %s\n", result.Symbol)
		fmt.Printf("  Decimals:     %d\n", result.Decimals)
		fmt.Printf("  Total Supply: %s\n", result.TotalSupply.String())
		fmt.Println()

		if result.Balance != nil {
			// BREAKPOINT: Set breakpoint here
			// DEBUG: Balance was successfully queried
			fmt.Println("Holder Balance:")
			fmt.Printf("  Address:      %s\n", holderAddress.Hex())
			fmt.Printf("  Raw Balance:  %s\n", result.Balance.String())

			// Convert to human-readable format
			divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(result.Decimals)), nil)
			balanceFloat := new(big.Float).SetInt(result.Balance)
			divisorFloat := new(big.Float).SetInt(divisor)
			humanReadable := new(big.Float).Quo(balanceFloat, divisorFloat)

			fmt.Printf("  Human:        %s %s\n", humanReadable.Text('f', int(result.Decimals)), result.Symbol)
			fmt.Println()
		}

		fmt.Println("ℹ balanceOf(address) demonstrates:")
		fmt.Println("  - BoundContract handling function parameters")
		fmt.Println("  - Automatic address encoding to bytes")
		fmt.Println("  - Type-safe parameter passing")
		fmt.Println()
	}

	// ============================================================================
	// STEP 6: Example 3 - Historical State with Balance
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through historical query

	if demo == "all" || demo == "historical" {
		fmt.Println("=== Example 3: Historical State Query ===")
		fmt.Println("Querying token state at historical block...")
		fmt.Println()

		// Get current block number
		header, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Printf("Could not get latest block: %v\n", err)
		} else {
			currentBlock := header.Number.Uint64()
			historicalBlock := currentBlock - 100000 // Go back 100k blocks

			fmt.Printf("Current Block:    %d\n", currentBlock)
			fmt.Printf("Historical Block: %d (current - 100000)\n", historicalBlock)
			fmt.Println()

			// Query historical state
			// BREAKPOINT: Set breakpoint here
			cfgHistorical := abigen.Config{
				Contract:    tokenAddress,
				BlockNumber: big.NewInt(int64(historicalBlock)),
				Holder:      holderAddress,
			}

			runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
			defer runCancel()

			historicalResult, err := abigen.Run(runCtx, client, cfgHistorical)
			if err != nil {
				log.Printf("Historical query failed: %v\n", err)
				fmt.Println()
				fmt.Println("ℹ Note: Some RPC providers don't support historical state")
				fmt.Println("  Use an archive node for full historical access")
			} else {
				// BREAKPOINT: Set breakpoint here
				fmt.Println("Historical Token State:")
				fmt.Printf("  Name:         %s\n", historicalResult.Name)
				fmt.Printf("  Symbol:       %s\n", historicalResult.Symbol)
				fmt.Printf("  Total Supply: %s\n", historicalResult.TotalSupply.String())

				if historicalResult.Balance != nil && holderAddress != nil {
					fmt.Println()
					fmt.Println("Historical Holder Balance:")
					fmt.Printf("  Address:      %s\n", holderAddress.Hex())
					fmt.Printf("  Balance:      %s\n", historicalResult.Balance.String())
				}
				fmt.Println()

				// Query current state for comparison
				cfgCurrent := abigen.Config{
					Contract:    tokenAddress,
					BlockNumber: nil,
					Holder:      holderAddress,
				}

				runCtx2, runCancel2 := context.WithTimeout(ctx, 10*time.Second)
				defer runCancel2()

				currentResult, err := abigen.Run(runCtx2, client, cfgCurrent)
				if err == nil {
					// Compare historical vs current
					// BREAKPOINT: Set breakpoint here
					supplyDiff := new(big.Int).Sub(currentResult.TotalSupply, historicalResult.TotalSupply)

					fmt.Println("Supply Change Analysis:")
					if supplyDiff.Sign() > 0 {
						fmt.Printf("  ✓ Supply increased by %s\n", supplyDiff.String())
					} else if supplyDiff.Sign() < 0 {
						fmt.Printf("  ✓ Supply decreased by %s\n", new(big.Int).Abs(supplyDiff).String())
					} else {
						fmt.Println("  ✓ Supply remained constant")
					}

					if currentResult.Balance != nil && historicalResult.Balance != nil {
						balanceDiff := new(big.Int).Sub(currentResult.Balance, historicalResult.Balance)
						if balanceDiff.Sign() != 0 {
							fmt.Println()
							fmt.Println("Holder Balance Change:")
							if balanceDiff.Sign() > 0 {
								fmt.Printf("  ✓ Balance increased by %s\n", balanceDiff.String())
							} else {
								fmt.Printf("  ✓ Balance decreased by %s\n", new(big.Int).Abs(balanceDiff).String())
							}
						}
					}
					fmt.Println()
				}
			}
		}
	}

	// ============================================================================
	// Understanding BoundContract Pattern
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand BoundContract

	if demo == "all" {
		fmt.Println("=== Understanding BoundContract Pattern ===")
		fmt.Println()

		fmt.Println("What is BoundContract?")
		fmt.Println("  - Adapter pattern for contract interactions")
		fmt.Println("  - Wraps low-level RPC with high-level interface")
		fmt.Println("  - Handles ABI encoding/decoding automatically")
		fmt.Println("  - Created with bind.NewBoundContract()")
		fmt.Println()

		fmt.Println("Key Components:")
		fmt.Println("  1. ABI Definition (JSON): Describes contract interface")
		fmt.Println("  2. Contract Address: Where the contract is deployed")
		fmt.Println("  3. Backend: RPC client for making calls")
		fmt.Println("  4. CallOpts: Configuration (context, block number)")
		fmt.Println()

		fmt.Println("Advantages over Manual ABI (Module 07):")
		fmt.Println("  ✓ No manual selector computation")
		fmt.Println("  ✓ No manual encoding/decoding logic")
		fmt.Println("  ✓ Type conversion handled automatically")
		fmt.Println("  ✓ Cleaner, more maintainable code")
		fmt.Println("  ✓ Easier to work with complex types")
		fmt.Println()

		fmt.Println("Comparison Chart:")
		fmt.Println("  Module 07 (Manual):    Full control, more code, error-prone")
		fmt.Println("  Module 08 (Bound):     Less code, automatic encoding, safer")
		fmt.Println("  Module 09+ (abigen):   Typed bindings, compile-time safety")
		fmt.Println()

		fmt.Println("When to use BoundContract:")
		fmt.Println("  - Dynamic contract interactions (ABI from API/database)")
		fmt.Println("  - Prototyping before generating typed bindings")
		fmt.Println("  - Educational purposes (understanding ABI layer)")
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
	// 4. Use F11 (Step Into) to enter abigen package functions
	// 5. Watch Variables panel to see BoundContract operations
	// 6. Inspect how CallOpts affects calls
	// 7. Use Call Stack panel to see function hierarchy
	//
	// BOUNDCONTRACT CONCEPTS TO EXPLORE:
	// - ABI Parsing: How is JSON ABI parsed?
	// - Contract Binding: How is contract bound to address?
	// - CallOpts: How do they configure calls?
	// - Type Conversion: How does abi.ConvertType work?
	// - Error Handling: What can go wrong?
	//
	// STEPPING INTO ABIGEN PACKAGE:
	// - Set breakpoint at abigen.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to solution.reference.go
	// - Step through to see: ABI parsing, BoundContract creation, calls
	// - Watch how helper functions (callString, callUint8, etc.) work
	// - See automatic encoding/decoding in action
	//
	// COMMON ERRORS TO DEBUG:
	// - "parse ABI": Malformed ABI JSON
	// - "call <method>": Contract call failed (wrong address, reverted)
	// - "empty result": Contract returned no data
	// - "type conversion": Unexpected return type
}
