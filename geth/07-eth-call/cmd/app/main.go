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
	"github.com/example/go-10x-minis/geth/07-eth-call/internal/ethcall"
)

/*
===================================================================
ERC20 Token Query with Manual ABI Encoding
===================================================================

This program demonstrates how to query ERC20 token metadata using
manual ABI encoding and decoding. You'll learn how function calls
are encoded as bytecode and how return values are decoded.

USAGE:
    go run ./cmd/app/main.go <RPC_URL> <TOKEN_ADDRESS>
    go run ./cmd/app/main.go <RPC_URL> <TOKEN_ADDRESS> <demo>

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    <TOKEN_ADDRESS> - ERC20 token contract address (required)
    [demo]          - Demo to run: all, metadata, latest, historical (default: "all")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
    go run ./cmd/app/main.go https://eth.llamarpc.com 0x6B175474E89094C44Da98b954EedeAC495271d0F metadata
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
- Step Into (F11) to enter ethcall package functions
- Watch Variables panel to see ABI encoding/decoding
- Inspect raw bytes returned from contract calls

KEY CONCEPTS:
- Function Selectors: First 4 bytes of keccak256(signature)
- ABI Encoding: How parameters are encoded to bytes
- ABI Decoding: How return values are decoded from bytes
- eth_call: Read-only contract execution (no gas cost, no state change)
- Static vs Dynamic Types: uint256 vs string encoding

ERC20 STANDARD:
- name(): Returns token name (string)
- symbol(): Returns token symbol (string)
- decimals(): Returns decimal places (uint8)
- totalSupply(): Returns total token supply (uint256)
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	// DEBUG: os.Args contains all command-line arguments
	// DEBUG: We need RPC URL and token contract address

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> <TOKEN_ADDRESS> [demo]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com 0x6B175474E89094C44Da98b954EedeAC495271d0F metadata\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, metadata, latest, historical\n")
		fmt.Fprintf(os.Stderr, "\nCommon tokens (Mainnet):\n")
		fmt.Fprintf(os.Stderr, "  USDC:  0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48\n")
		fmt.Fprintf(os.Stderr, "  DAI:   0x6B175474E89094C44Da98b954EedeAC495271d0F\n")
		fmt.Fprintf(os.Stderr, "  WETH:  0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	tokenAddressHex := os.Args[2]

	demo := "all"
	if len(os.Args) > 3 {
		demo = os.Args[3]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL', 'tokenAddressHex', 'demo' variables
	fmt.Printf("=== ERC20 Token Query with Manual ABI Encoding ===\n")
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Demo: %s\n\n", demo)

	// ============================================================================
	// STEP 2: Parse Token Contract Address
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before parsing address
	// DEBUG: Token contracts are deployed at specific addresses
	// DEBUG: We need a valid address to make contract calls

	if !common.IsHexAddress(tokenAddressHex) {
		log.Fatalf("Invalid token address: %s\n", tokenAddressHex)
	}
	tokenAddress := common.HexToAddress(tokenAddressHex)

	// BREAKPOINT: Set breakpoint here after parsing address
	// DEBUG: Inspect 'tokenAddress' value
	fmt.Printf("✓ Parsed token address: %s\n", tokenAddress.Hex())
	fmt.Println()

	// ============================================================================
	// STEP 3: Connect to Ethereum RPC Endpoint
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before connection
	// DEBUG: This establishes the RPC connection
	// DEBUG: Watch for connection errors

	ctx := context.Background()
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// BREAKPOINT: Set breakpoint here
	// DEBUG: Step Over (F10) to execute the Dial call
	client, err := ethclient.Dial(dialCtx, rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to RPC endpoint: %v\n", err)
	}
	defer client.Close()

	// BREAKPOINT: Set breakpoint here after successful connection
	fmt.Println("✓ Connected to Ethereum RPC endpoint")
	fmt.Println()

	// ============================================================================
	// STEP 4: Example 1 - Query Token Metadata (Latest Block)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through first example
	// DEBUG: This demonstrates querying current token metadata
	// DEBUG: Uses latest block (nil BlockNumber)

	if demo == "all" || demo == "metadata" {
		fmt.Println("=== Example 1: Query Token Metadata (Latest Block) ===")
		fmt.Println("Querying ERC20 token metadata from latest block...")
		fmt.Println()

		// Create configuration for latest block query
		// BREAKPOINT: Set breakpoint here
		// DEBUG: BlockNumber=nil means query latest block
		// DEBUG: This is the most common use case
		cfg := ethcall.Config{
			Contract:    tokenAddress,
			BlockNumber: nil, // nil = latest block
		}

		// BREAKPOINT: Set breakpoint here before calling ethcall.Run
		// DEBUG: Step Into (F11) to see manual ABI encoding/decoding
		// DEBUG: You'll see function selectors, CallContract, and decoding logic
		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		result, err := ethcall.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to query token metadata: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful query
		// DEBUG: Expand 'result' in Variables panel to inspect all fields
		// DEBUG: Look at Name, Symbol, Decimals, TotalSupply

		fmt.Println("Token Metadata:")
		fmt.Printf("  Name:         %s\n", result.Name)
		fmt.Printf("  Symbol:       %s\n", result.Symbol)
		fmt.Printf("  Decimals:     %d\n", result.Decimals)
		fmt.Printf("  Total Supply: %s (raw)\n", result.TotalSupply.String())
		fmt.Println()

		// Convert total supply to human-readable format
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Decimals tells us how to display the number
		// DEBUG: For USDC (decimals=6): 1000000 raw units = 1.0 USDC
		// DEBUG: For DAI (decimals=18): 10^18 raw units = 1.0 DAI
		divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(result.Decimals)), nil)
		supplyFloat := new(big.Float).SetInt(result.TotalSupply)
		divisorFloat := new(big.Float).SetInt(divisor)
		humanReadable := new(big.Float).Quo(supplyFloat, divisorFloat)

		fmt.Println("Human-Readable Total Supply:")
		fmt.Printf("  %s %s\n", humanReadable.Text('f', int(result.Decimals)), result.Symbol)
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice how we decoded the data without typed bindings
		// DEBUG: Manual ABI decoding gives us full control
		fmt.Println("ℹ Data queried from latest block (current state)")
		fmt.Println()
	}

	// ============================================================================
	// STEP 5: Example 2 - Understanding Function Selectors
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand function selectors
	// DEBUG: This demonstrates how contract functions are identified
	// DEBUG: Function selectors are the first 4 bytes of keccak256(signature)

	if demo == "all" || demo == "latest" {
		fmt.Println("=== Example 2: Understanding Function Selectors ===")
		fmt.Println("Function selectors identify which function to call...")
		fmt.Println()

		// Query token metadata again to demonstrate
		cfg := ethcall.Config{
			Contract:    tokenAddress,
			BlockNumber: nil,
		}

		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		result, err := ethcall.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to query token: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Function selectors are computed from function signatures

		fmt.Println("Function Selector Calculation:")
		fmt.Println("  name() selector:        0x06fdde03 = keccak256(\"name()\")[:4]")
		fmt.Println("  symbol() selector:      0x95d89b41 = keccak256(\"symbol()\")[:4]")
		fmt.Println("  decimals() selector:    0x313ce567 = keccak256(\"decimals()\")[:4]")
		fmt.Println("  totalSupply() selector: 0x18160ddd = keccak256(\"totalSupply()\")[:4]")
		fmt.Println()

		fmt.Println("How it works:")
		fmt.Println("  1. Take function signature: \"name()\"")
		fmt.Println("  2. Hash with keccak256")
		fmt.Println("  3. Take first 4 bytes")
		fmt.Println("  4. Send as calldata to contract")
		fmt.Println()

		fmt.Println("Query Results:")
		fmt.Printf("  %s (%s) - %d decimals\n", result.Name, result.Symbol, result.Decimals)
		fmt.Println()
	}

	// ============================================================================
	// STEP 6: Example 3 - Historical State Query
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through historical query
	// DEBUG: This demonstrates querying state at a specific block
	// DEBUG: Useful for time-travel queries and historical analysis

	if demo == "all" || demo == "historical" {
		fmt.Println("=== Example 3: Historical State Query ===")
		fmt.Println("Querying token state at a specific historical block...")
		fmt.Println()

		// Query latest block first to show comparison
		// BREAKPOINT: Set breakpoint here
		cfgLatest := ethcall.Config{
			Contract:    tokenAddress,
			BlockNumber: nil,
		}

		runCtx1, runCancel1 := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel1()

		latestResult, err := ethcall.Run(runCtx1, client, cfgLatest)
		if err != nil {
			log.Fatalf("Failed to query latest state: %v\n", err)
		}

		fmt.Println("Current State (Latest Block):")
		fmt.Printf("  Name:         %s\n", latestResult.Name)
		fmt.Printf("  Symbol:       %s\n", latestResult.Symbol)
		fmt.Printf("  Total Supply: %s\n", latestResult.TotalSupply.String())
		fmt.Println()

		// Query historical block (1 million blocks ago as example)
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Get current block number, subtract offset for historical query
		header, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Printf("Could not get latest block for historical demo: %v\n", err)
		} else {
			currentBlock := header.Number.Uint64()
			historicalBlock := currentBlock - 100000 // Go back 100k blocks

			fmt.Printf("Querying Historical Block: %d (current - 100000)\n", historicalBlock)
			fmt.Println()

			cfgHistorical := ethcall.Config{
				Contract:    tokenAddress,
				BlockNumber: big.NewInt(int64(historicalBlock)),
			}

			runCtx2, runCancel2 := context.WithTimeout(ctx, 10*time.Second)
			defer runCancel2()

			historicalResult, err := ethcall.Run(runCtx2, client, cfgHistorical)
			if err != nil {
				log.Printf("Historical query failed (this is expected for some RPC providers): %v\n", err)
				fmt.Println()
				fmt.Println("ℹ Note: Some RPC providers don't support historical state queries")
				fmt.Println("  Use an archive node for full historical access")
			} else {
				// BREAKPOINT: Set breakpoint here
				// DEBUG: Compare historical vs current total supply
				// DEBUG: Total supply may have changed (minting/burning)

				fmt.Println("Historical State:")
				fmt.Printf("  Name:         %s\n", historicalResult.Name)
				fmt.Printf("  Symbol:       %s\n", historicalResult.Symbol)
				fmt.Printf("  Total Supply: %s\n", historicalResult.TotalSupply.String())
				fmt.Println()

				// Calculate supply change
				supplyDiff := new(big.Int).Sub(latestResult.TotalSupply, historicalResult.TotalSupply)
				if supplyDiff.Sign() > 0 {
					fmt.Printf("✓ Supply increased by %s (tokens were minted)\n", supplyDiff.String())
				} else if supplyDiff.Sign() < 0 {
					fmt.Printf("✓ Supply decreased by %s (tokens were burned)\n", new(big.Int).Abs(supplyDiff).String())
				} else {
					fmt.Println("✓ Supply remained constant")
				}
				fmt.Println()
			}
		}

		fmt.Println("ℹ Historical queries show time-travel capabilities")
		fmt.Println("  - Track supply changes over time")
		fmt.Println("  - Audit token economics")
		fmt.Println("  - Requires archive node for full history")
		fmt.Println()
	}

	// ============================================================================
	// Understanding ABI Encoding/Decoding
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand ABI encoding/decoding

	if demo == "all" {
		fmt.Println("=== Understanding ABI Encoding/Decoding ===")
		fmt.Println()

		fmt.Println("Static Types (Fixed Size):")
		fmt.Println("  - uint8, uint256, address, bool, bytes32")
		fmt.Println("  - Encoded as 32 bytes (right-padded or left-padded)")
		fmt.Println("  - Example: uint8(255) = 0x00...00ff (31 zeros + ff)")
		fmt.Println()

		fmt.Println("Dynamic Types (Variable Size):")
		fmt.Println("  - string, bytes, arrays")
		fmt.Println("  - Encoded with offset + length + data")
		fmt.Println("  - Offset: position where data starts (32 bytes)")
		fmt.Println("  - Length: number of bytes (32 bytes)")
		fmt.Println("  - Data: actual content (padded to 32-byte boundary)")
		fmt.Println()

		fmt.Println("Function Call Encoding:")
		fmt.Println("  CallData = [selector (4 bytes)] + [encoded parameters]")
		fmt.Println()

		fmt.Println("Example: name() call")
		fmt.Println("  - Selector: 0x06fdde03")
		fmt.Println("  - Parameters: none")
		fmt.Println("  - CallData: 0x06fdde03")
		fmt.Println()

		fmt.Println("Example: name() return value \"USDC\"")
		fmt.Println("  - Offset: 0x0020 (data starts at byte 32)")
		fmt.Println("  - Length: 0x0004 (4 characters)")
		fmt.Println("  - Data: 0x55534443 (\"USDC\" in hex)")
		fmt.Println()

		fmt.Println("Why Manual ABI Matters:")
		fmt.Println("  1. Understanding how contracts communicate")
		fmt.Println("  2. Debugging contract calls")
		fmt.Println("  3. Working with untyped contracts")
		fmt.Println("  4. Foundation for typed bindings (abigen)")
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
	// 4. Use F11 (Step Into) to enter ethcall package functions
	// 5. Watch Variables panel to see ABI encoding/decoding
	// 6. Inspect raw bytes returned from CallContract
	// 7. Use Call Stack panel to see function hierarchy
	//
	// ABI CONCEPTS TO EXPLORE:
	// - Function Selectors: How are they computed?
	// - Static vs Dynamic Types: How do they differ?
	// - Offset Encoding: Why do dynamic types need offsets?
	// - Padding: Why are values padded to 32 bytes?
	// - eth_call: How does it differ from transactions?
	//
	// STEPPING INTO ETHCALL PACKAGE:
	// - Set breakpoint at ethcall.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to solution.reference.go
	// - Step through to see: selector computation, CallContract, decoding
	// - Watch how raw bytes are decoded to strings, uint8, uint256
	// - See error handling for malformed data
	//
	// COMMON ERRORS TO DEBUG:
	// - "data too short": Contract returned less data than expected
	// - "invalid offset": Dynamic type offset points outside data bounds
	// - "execution reverted": Contract rejected the call (wrong function?)
	// - "contract not found": Address doesn't have contract code
}
