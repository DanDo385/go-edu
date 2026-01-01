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
	"github.com/example/go-10x-minis/geth/04-accounts-balances/internal/accountsbalances"
)

/*
===================================================================
Ethereum Accounts and Balances - Account State Queries
===================================================================

This program demonstrates how to query Ethereum account states,
retrieve balances, and classify accounts as either Externally
Owned Accounts (EOAs) or Smart Contracts.

USAGE:
    go run ./cmd/app/main.go <RPC_URL>
    go run ./cmd/app/main.go <RPC_URL> [demo]

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    [demo]          - Demo to run: all, single, multiple, classify, historical (default: "all")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com
    go run ./cmd/app/main.go https://eth.llamarpc.com single
    go run ./cmd/app/main.go https://eth.llamarpc.com multiple
    go run ./cmd/app/main.go https://eth.llamarpc.com classify
    go run ./cmd/app/main.go https://eth.llamarpc.com historical

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
- Step Into (F11) to enter accountsbalances package functions
- Watch Variables panel to see account states and balances
- Inspect AccountState structures with Balance, Code, and Type

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the accountsbalances package from internal/accountsbalances
- exercise.go, solution.reference.go are in package accountsbalances
- The accountsbalances.Run() function queries account states

KEY CONCEPTS:
- EOA (Externally Owned Account): Controlled by private key, no code
- Smart Contract: Has code, controlled by its logic, not a private key
- Balance: Amount of ETH (in wei) held by an account
- Code: Bytecode stored at a contract address (empty for EOAs)
- Account State: Snapshot of account properties at a specific block
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	// DEBUG: os.Args[0] is the program name
	// DEBUG: os.Args[1] is the RPC URL (required)
	// DEBUG: os.Args[2] is the demo type (optional)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> [demo]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com single\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com multiple\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, single, multiple, classify, historical\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	demo := "all"
	if len(os.Args) > 2 {
		demo = os.Args[2]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL' and 'demo' variables in Variables panel
	fmt.Printf("=== Ethereum Accounts and Balances ===\n")
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Demo: %s\n\n", demo)

	// ============================================================================
	// STEP 2: Connect to Ethereum RPC Endpoint
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before connection
	// DEBUG: This is where we establish the RPC connection
	// DEBUG: Watch for any connection errors in the error variable

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
	// STEP 3: Example 1 - Query Single Account Balance
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through single account query
	// DEBUG: This demonstrates querying the state of a well-known address

	if demo == "all" || demo == "single" {
		fmt.Println("=== Example 1: Query Single Account Balance ===")
		fmt.Println("Querying Vitalik Buterin's public address...")
		fmt.Println()

		// Vitalik's well-known public address
		// BREAKPOINT: Set breakpoint here
		// DEBUG: This is a real Ethereum address that often holds ETH
		vitalikAddr := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

		cfg := accountsbalances.Config{
			Addresses:   []common.Address{vitalikAddr},
			BlockNumber: nil, // nil = latest block
		}

		// BREAKPOINT: Set breakpoint here before calling accountsbalances.Run
		// DEBUG: Step Into (F11) to see the implementation
		// DEBUG: You'll see BalanceAt() and CodeAt() RPC calls
		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		result, err := accountsbalances.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to query account: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful retrieval
		// DEBUG: Expand 'result' in Variables panel to inspect all fields
		// DEBUG: Look at Accounts[0] to see Address, Balance, Code, Type

		if len(result.Accounts) > 0 {
			account := result.Accounts[0]
			fmt.Println("Account Information:")
			fmt.Printf("  Address: %s\n", account.Address.Hex())
			fmt.Printf("  Type:    %s\n", account.Type)

			// Convert balance from wei to ETH (1 ETH = 10^18 wei)
			balanceETH := new(big.Float).SetInt(account.Balance)
			balanceETH.Quo(balanceETH, big.NewFloat(1e18))

			fmt.Printf("  Balance: %s wei\n", account.Balance.String())
			fmt.Printf("           %.6f ETH\n", balanceETH)
			fmt.Printf("  Code:    %d bytes\n", len(account.Code))
			fmt.Println()

			// BREAKPOINT: Set breakpoint here
			// DEBUG: Notice that this is an EOA (no code)
			if account.Type == accountsbalances.AccountTypeEOA {
				fmt.Println("This is an EOA (Externally Owned Account):")
				fmt.Println("  - Controlled by a private key")
				fmt.Println("  - Can initiate transactions")
				fmt.Println("  - No code (contract logic)")
			}
			fmt.Println()
		}
	}

	// ============================================================================
	// STEP 4: Example 2 - Query Multiple Accounts
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through batch query
	// DEBUG: This demonstrates querying multiple accounts in one call

	if demo == "all" || demo == "multiple" {
		fmt.Println("=== Example 2: Query Multiple Accounts ===")
		fmt.Println("Querying multiple well-known addresses...")
		fmt.Println()

		// Collection of well-known addresses
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice we're querying both EOAs and contracts
		addresses := []common.Address{
			common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), // Vitalik (EOA)
			common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"), // WETH contract
			common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), // USDC contract
		}

		cfg := accountsbalances.Config{
			Addresses:   addresses,
			BlockNumber: nil, // latest block
		}

		runCtx, runCancel := context.WithTimeout(ctx, 15*time.Second)
		defer runCancel()

		// BREAKPOINT: Set breakpoint here before the call
		// DEBUG: Step Into (F11) to see the batch processing loop
		result, err := accountsbalances.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to query accounts: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Expand result.Accounts to see all queried accounts
		// DEBUG: Notice some are EOA, some are Contract

		fmt.Printf("Retrieved %d accounts:\n\n", len(result.Accounts))

		for i, account := range result.Accounts {
			// Convert balance to ETH
			balanceETH := new(big.Float).SetInt(account.Balance)
			balanceETH.Quo(balanceETH, big.NewFloat(1e18))

			// BREAKPOINT: Set breakpoint here inside the loop
			// DEBUG: Watch how each account has different properties
			fmt.Printf("[%d] %s\n", i+1, account.Address.Hex())
			fmt.Printf("    Type:    %s\n", account.Type)
			fmt.Printf("    Balance: %.6f ETH\n", balanceETH)
			fmt.Printf("    Code:    %d bytes\n", len(account.Code))
			fmt.Println()
		}
	}

	// ============================================================================
	// STEP 5: Example 3 - Classify Accounts (EOA vs Contract)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand account classification
	// DEBUG: This demonstrates the fundamental difference between account types

	if demo == "all" || demo == "classify" {
		fmt.Println("=== Example 3: Account Classification (EOA vs Contract) ===")
		fmt.Println()

		fmt.Println("Understanding Account Types:")
		fmt.Println()
		fmt.Println("1. EOA (Externally Owned Account):")
		fmt.Println("   - Controlled by a private key (not code)")
		fmt.Println("   - Can initiate transactions")
		fmt.Println("   - Has NO bytecode (len(code) == 0)")
		fmt.Println("   - Examples: User wallets, exchange hot wallets")
		fmt.Println()
		fmt.Println("2. Smart Contract:")
		fmt.Println("   - Controlled by its code logic")
		fmt.Println("   - Cannot initiate transactions (only respond to calls)")
		fmt.Println("   - HAS bytecode (len(code) > 0)")
		fmt.Println("   - Examples: ERC-20 tokens, DEX contracts, NFT contracts")
		fmt.Println()

		// Test with known examples
		// BREAKPOINT: Set breakpoint here
		// DEBUG: We'll classify both types of accounts
		fmt.Println("Testing Classification:")
		fmt.Println()

		testAddresses := []struct {
			addr        common.Address
			description string
		}{
			{common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), "Vitalik's Address"},
			{common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"), "WETH Contract"},
			{common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), "USDC Contract"},
		}

		for _, test := range testAddresses {
			cfg := accountsbalances.Config{
				Addresses:   []common.Address{test.addr},
				BlockNumber: nil,
			}

			runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
			result, err := accountsbalances.Run(runCtx, client, cfg)
			runCancel()

			if err != nil {
				log.Printf("Failed to query %s: %v\n", test.description, err)
				continue
			}

			// BREAKPOINT: Set breakpoint here
			// DEBUG: Inspect how the Type field was determined
			// DEBUG: Look at len(account.Code) to see the classification logic
			if len(result.Accounts) > 0 {
				account := result.Accounts[0]
				fmt.Printf("%s:\n", test.description)
				fmt.Printf("  Address: %s\n", account.Address.Hex())
				fmt.Printf("  Type:    %s\n", account.Type)
				fmt.Printf("  Code:    %d bytes\n", len(account.Code))

				// BREAKPOINT: Set breakpoint here
				// DEBUG: This explains the classification decision
				if account.Type == accountsbalances.AccountTypeContract {
					fmt.Printf("  ✓ Has bytecode → Classified as Contract\n")
				} else {
					fmt.Printf("  ✓ No bytecode → Classified as EOA\n")
				}
				fmt.Println()
			}
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: This summarizes the classification algorithm
		fmt.Println("Classification Algorithm:")
		fmt.Println("  if len(code) > 0:")
		fmt.Println("    account.Type = Contract")
		fmt.Println("  else:")
		fmt.Println("    account.Type = EOA")
		fmt.Println()
	}

	// ============================================================================
	// STEP 6: Example 4 - Historical Balance Queries
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand historical queries
	// DEBUG: This demonstrates querying account state at a specific block

	if demo == "all" || demo == "historical" {
		fmt.Println("=== Example 4: Historical Balance Queries ===")
		fmt.Println("Querying account balance at a specific block height...")
		fmt.Println()

		// Query the same address at different block heights
		// BREAKPOINT: Set breakpoint here
		// DEBUG: We'll see how account state changes over time
		addr := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

		// Query at block 15,000,000 (The Merge block, September 2022)
		historicalBlock := big.NewInt(15_000_000)

		fmt.Println("Comparing current vs historical state:")
		fmt.Println()

		// Current state
		cfg1 := accountsbalances.Config{
			Addresses:   []common.Address{addr},
			BlockNumber: nil, // latest
		}

		runCtx1, runCancel1 := context.WithTimeout(ctx, 10*time.Second)
		result1, err := accountsbalances.Run(runCtx1, client, cfg1)
		runCancel1()

		if err != nil {
			log.Fatalf("Failed to query current state: %v\n", err)
		}

		// Historical state
		cfg2 := accountsbalances.Config{
			Addresses:   []common.Address{addr},
			BlockNumber: historicalBlock,
		}

		runCtx2, runCancel2 := context.WithTimeout(ctx, 10*time.Second)
		result2, err := accountsbalances.Run(runCtx2, client, cfg2)
		runCancel2()

		if err != nil {
			log.Fatalf("Failed to query historical state: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Compare the two results to see balance changes
		if len(result1.Accounts) > 0 && len(result2.Accounts) > 0 {
			current := result1.Accounts[0]
			historical := result2.Accounts[0]

			currentETH := new(big.Float).SetInt(current.Balance)
			currentETH.Quo(currentETH, big.NewFloat(1e18))

			historicalETH := new(big.Float).SetInt(historical.Balance)
			historicalETH.Quo(historicalETH, big.NewFloat(1e18))

			fmt.Printf("Address: %s\n\n", addr.Hex())
			fmt.Println("Current State (Latest Block):")
			fmt.Printf("  Balance: %.6f ETH\n", currentETH)
			fmt.Printf("  Type:    %s\n\n", current.Type)

			fmt.Printf("Historical State (Block %s):\n", historicalBlock.String())
			fmt.Printf("  Balance: %.6f ETH\n", historicalETH)
			fmt.Printf("  Type:    %s\n\n", historical.Type)

			// Calculate difference
			diff := new(big.Float).Sub(currentETH, historicalETH)
			fmt.Printf("Change: %.6f ETH\n", diff)
			fmt.Println()
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: This explains the importance of historical queries
		fmt.Println("Why Historical Queries Matter:")
		fmt.Println("  - Track account balance changes over time")
		fmt.Println("  - Audit and verify transactions")
		fmt.Println("  - Analyze historical contract states")
		fmt.Println("  - Build time-series data for analytics")
		fmt.Println()
	}

	// ============================================================================
	// STEP 7: Understanding Account State and Defensive Copying
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand data safety
	// DEBUG: This demonstrates defensive copying for account data

	if demo == "all" {
		fmt.Println("=== Understanding Account State and Defensive Copying ===")
		fmt.Println()

		cfg := accountsbalances.Config{
			Addresses:   []common.Address{common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")},
			BlockNumber: nil,
		}

		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		result, err := accountsbalances.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to retrieve data: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice that result.Accounts contains defensively copied data
		// DEBUG: Modifying it won't affect the RPC client's internal state

		fmt.Println("Defensive Copying for Account Data:")
		fmt.Println("  1. Balance (*big.Int) is copied with new(big.Int).Set()")
		fmt.Println("  2. Code ([]byte) is copied with append([]byte(nil), code...)")
		fmt.Println("  3. Address (common.Address) is a value type (copied automatically)")
		fmt.Println()

		if len(result.Accounts) > 0 {
			account := result.Accounts[0]

			fmt.Println("Example: Modifying returned data is safe")
			originalBalance := new(big.Int).Set(account.Balance)
			account.Balance.Add(account.Balance, big.NewInt(1000000000000000000)) // Add 1 ETH

			fmt.Printf("  Original Balance: %s wei\n", originalBalance.String())
			fmt.Printf("  Modified Balance: %s wei\n", account.Balance.String())
			fmt.Println("  ✓ Modification only affects our copy, not RPC client's data")
			fmt.Println()
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: This explains why defensive copying prevents data corruption
		fmt.Println("Why This Matters:")
		fmt.Println("  - Prevents accidental mutations of shared state")
		fmt.Println("  - Enables safe concurrent access to data")
		fmt.Println("  - Follows Go best practices for API design")
		fmt.Println("  - Critical for multi-threaded applications")
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
	// 4. Use F11 (Step Into) to enter accountsbalances package functions
	// 5. Watch Variables panel to see account states
	// 6. Inspect AccountState structures with Balance, Code, and Type
	// 7. Use Call Stack panel to see function hierarchy
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - EOA vs Contract: What's the fundamental difference?
	// - Account State: What properties define an account?
	// - Balance Queries: How are balances stored on-chain?
	// - Code Storage: Where is contract bytecode stored?
	// - Historical Queries: How to query past states?
	//
	// STEPPING INTO ACCOUNTSBALANCES PACKAGE:
	// - Set breakpoint at accountsbalances.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to exercise.go or solution.reference.go
	// - Step through to see RPC calls: BalanceAt(), CodeAt()
	// - Watch the classification logic based on code length
	// - Observe defensive copying of balances and code
	//
	// COMMON RPC ERRORS TO DEBUG:
	// - Connection timeout: Check network connectivity
	// - Invalid address: Verify address format (0x prefix + 40 hex chars)
	// - Block not found: Historical blocks might not be available
	// - Rate limiting: Public endpoints may have rate limits
	//
	// ADVANCED DEBUGGING:
	// - Compare EOA vs Contract bytecode
	// - Track balance changes across blocks
	// - Analyze contract code size variations
	// - Test with different RPC providers
}
