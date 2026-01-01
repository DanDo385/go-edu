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
	"github.com/example/go-10x-minis/geth/11-storage/internal/storage"
)

/*
===================================================================
Ethereum Storage Reader
===================================================================

This program demonstrates how to read raw storage slots from Ethereum
smart contracts, including computing mapping slots using keccak256 hashing.

USAGE:
    go run ./cmd/app/main.go <RPC_URL>
    go run ./cmd/app/main.go <RPC_URL> <DEMO>

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    [demo]          - Demo to run: all, simple, mapping, usdc, historical (default: "all")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com
    go run ./cmd/app/main.go https://eth.llamarpc.com simple
    go run ./cmd/app/main.go https://eth.llamarpc.com mapping
    go run ./cmd/app/main.go https://eth.llamarpc.com usdc
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
- Step Into (F11) to enter storage package functions
- Watch Variables panel to see storage values
- Inspect slot calculations for mappings

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the storage package from internal/storage
- exercise.go, solution.reference.go are in package storage
- The storage.Run() function encapsulates all storage reading logic

KEY CONCEPTS:
- Storage Slots: 2^256 possible 32-byte slots per contract
- Merkle-Patricia Trie: Storage is organized as a cryptographic tree
- Storage Root: Hash committing to all storage (in account state)
- Mapping Slots: keccak256(key || baseSlot) for Solidity mappings
- Slot Layout: Sequential for simple variables, computed for mappings/arrays
- Historical Reads: Read storage from any past block (requires archive node)
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
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com simple\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com mapping\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, simple, mapping, usdc, historical\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	demo := "all"
	if len(os.Args) > 2 {
		demo = os.Args[2]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL' and 'demo' variables in Variables panel
	fmt.Printf("=== Ethereum Storage Reader ===\n")
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Demo: %s\n\n", demo)

	// ============================================================================
	// STEP 2: Connect to Ethereum RPC Endpoint
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before connection
	// DEBUG: This is where we establish the RPC connection
	// DEBUG: Watch for any connection errors in the error variable
	//
	// Context with timeout: We use context.WithTimeout to prevent hanging
	// indefinitely if the RPC endpoint is slow or unreachable. This is a
	// critical pattern for production code—always set timeouts on network calls.
	//
	// ethclient.Dial: This is the standard go-ethereum client that implements
	// the StorageClient interface. It handles JSON-RPC communication with nodes.

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
	// STEP 3: Example 1 - Read Simple Storage Slot
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through first example
	// DEBUG: This demonstrates reading a simple storage slot (slot 0)
	// DEBUG: Watch the 'result' variable to see the raw storage value
	//
	// CONCEPT: Simple storage variables in Solidity are stored sequentially:
	//   - First variable → slot 0
	//   - Second variable → slot 1
	//   - Third variable → slot 2, etc.
	//
	// Example contract:
	//   uint256 totalSupply;  // slot 0
	//   address owner;        // slot 1
	//   bool paused;          // slot 2

	if demo == "all" || demo == "simple" {
		fmt.Println("=== Example 1: Read Simple Storage Slot ===")
		fmt.Println("Reading slot 0 from a well-known contract...")
		fmt.Println()

		// USDC contract on Ethereum mainnet (proxy contract)
		// We'll read slot 0 which is often used for implementation address in proxies
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice the contract address—this is USDC on Ethereum mainnet
		usdcAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

		// Create configuration for simple slot read
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice cfg.Slot is 0—reading the first storage slot
		// DEBUG: BlockNumber is nil—reading from latest block
		// DEBUG: MappingKey is nil—this is a simple slot, not a mapping
		cfg := storage.Config{
			Contract:    usdcAddress,
			Slot:        big.NewInt(0),
			MappingKey:  nil, // Not a mapping
			BlockNumber: nil, // Latest block
		}

		// BREAKPOINT: Set breakpoint here before calling storage.Run
		// DEBUG: Step Into (F11) to see the implementation in solution.reference.go
		// DEBUG: You'll see how the function validates inputs and makes RPC calls
		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		result, err := storage.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to read storage slot 0: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful read
		// DEBUG: Expand 'result' in Variables panel to inspect fields
		// DEBUG: Look at ResolvedSlot (the actual slot that was read)
		// DEBUG: Look at Value (the raw 32-byte storage value)

		fmt.Printf("Contract: %s\n", usdcAddress.Hex())
		fmt.Printf("Slot: %s\n", cfg.Slot.String())
		fmt.Printf("Resolved Slot Hash: %s\n", result.ResolvedSlot.Hex())
		fmt.Printf("Raw Value (hex): %s\n", common.Bytes2Hex(result.Value))
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: The raw value is displayed—interpretation depends on the type
		// DEBUG: For addresses, take the rightmost 20 bytes
		// DEBUG: For uint256, interpret as big-endian number
	}

	// ============================================================================
	// STEP 4: Example 2 - Read Mapping Storage Slot
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through mapping example
	// DEBUG: This demonstrates computing mapping slots using keccak256
	// DEBUG: This is the low-level implementation of Solidity mappings!
	//
	// CONCEPT: Solidity mappings use cryptographic hashing to compute slots:
	//   mapping(address => uint256) balances;  // slot 0
	//   balances[addr] is stored at: keccak256(addr || 0x00...00)
	//
	// This distributes values across the 2^256 storage space, preventing collisions.

	if demo == "all" || demo == "mapping" {
		fmt.Println("=== Example 2: Read Mapping Storage Slot ===")
		fmt.Println("Reading a balance from USDC balances mapping...")
		fmt.Println()

		// USDC contract stores balances in a mapping at slot 9
		// mapping(address => uint256) balances; // slot 9
		// BREAKPOINT: Set breakpoint here
		usdcAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

		// Vitalik's address as an example
		vitalikAddress := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

		// Create configuration for mapping slot read
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice cfg.Slot is 9—this is the base slot for the balances mapping
		// DEBUG: MappingKey is vitalikAddress—we want balances[vitalik]
		// DEBUG: The actual slot will be computed as keccak256(vitalik || slot9)
		cfg := storage.Config{
			Contract:    usdcAddress,
			Slot:        big.NewInt(9), // balances mapping base slot
			MappingKey:  vitalikAddress.Bytes(),
			BlockNumber: nil, // Latest block
		}

		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		// BREAKPOINT: Set breakpoint here before calling storage.Run
		// DEBUG: Step Into (F11) to see how mappingSlotHash is computed
		// DEBUG: The function will hash the key and base slot to get the actual slot
		result, err := storage.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to read mapping slot: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful read
		// DEBUG: Compare ResolvedSlot with the base slot (9)
		// DEBUG: They're different! ResolvedSlot is keccak256(key || baseSlot)

		fmt.Printf("Contract: %s\n", usdcAddress.Hex())
		fmt.Printf("Mapping Base Slot: %s\n", cfg.Slot.String())
		fmt.Printf("Mapping Key (address): %s\n", vitalikAddress.Hex())
		fmt.Printf("Resolved Slot Hash: %s\n", result.ResolvedSlot.Hex())
		fmt.Printf("Raw Value (hex): %s\n", common.Bytes2Hex(result.Value))

		// Interpret the value as uint256 (USDC balance)
		// USDC has 6 decimals, so divide by 1,000,000 to get actual USDC amount
		balance := new(big.Int).SetBytes(result.Value)
		fmt.Printf("Balance (raw): %s\n", balance.String())
		if balance.Cmp(big.NewInt(0)) > 0 {
			// Convert to USDC (6 decimals)
			divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)
			usdcBalance := new(big.Int).Div(balance, divisor)
			remainder := new(big.Int).Mod(balance, divisor)
			fmt.Printf("Balance (USDC): %s.%06d USDC\n", usdcBalance.String(), remainder.Int64())
		}
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice how the mapping slot calculation enables reading any key
		// DEBUG: This is the cryptographic hash table pattern in action!
	}

	// ============================================================================
	// STEP 5: Example 3 - Read Historical Storage
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through historical read
	// DEBUG: This demonstrates reading storage from a specific past block
	// DEBUG: Requires an archive node with historical state

	if demo == "all" || demo == "historical" {
		fmt.Println("=== Example 3: Read Historical Storage ===")
		fmt.Println("Reading storage from a historical block...")
		fmt.Println()

		// Read from a specific historical block (block 18,000,000)
		// BREAKPOINT: Set breakpoint here
		usdcAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
		historicalBlock := big.NewInt(18_000_000)

		cfg := storage.Config{
			Contract:    usdcAddress,
			Slot:        big.NewInt(0),
			MappingKey:  nil,
			BlockNumber: historicalBlock,
		}

		runCtx, runCancel := context.WithTimeout(ctx, 5*time.Second)
		defer runCancel()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Note the BlockNumber parameter—this requests historical state
		// DEBUG: Many RPC providers don't support this (requires archive node)
		result, err := storage.Run(runCtx, client, cfg)
		if err != nil {
			// This is expected if using a non-archive node
			fmt.Printf("Note: Historical storage read failed (may need archive node): %v\n", err)
			fmt.Println()
		} else {
			// BREAKPOINT: Set breakpoint here if successful
			// DEBUG: You have storage from a specific block in history!
			fmt.Printf("Contract: %s\n", usdcAddress.Hex())
			fmt.Printf("Slot: %s\n", cfg.Slot.String())
			fmt.Printf("Block Number: %s\n", historicalBlock.String())
			fmt.Printf("Resolved Slot: %s\n", result.ResolvedSlot.Hex())
			fmt.Printf("Raw Value (hex): %s\n", common.Bytes2Hex(result.Value))
			fmt.Println()
		}
	}

	// ============================================================================
	// STEP 6: Example 4 - Understanding Storage Layout
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand storage layout
	// DEBUG: This demonstrates how Solidity packs multiple variables into slots

	if demo == "all" {
		fmt.Println("=== Example 4: Understanding Storage Layout ===")
		fmt.Println()

		fmt.Println("Solidity Storage Layout Principles:")
		fmt.Println("  1. Simple variables are stored sequentially starting at slot 0")
		fmt.Println("  2. Variables < 32 bytes may be packed into the same slot")
		fmt.Println("  3. Mappings don't occupy their base slot (it's empty)")
		fmt.Println("  4. Mapping values: keccak256(key || baseSlot)")
		fmt.Println("  5. Dynamic arrays: length at baseSlot, elements at keccak256(baseSlot)")
		fmt.Println()

		fmt.Println("Example Contract Storage Layout:")
		fmt.Println("  contract Example {")
		fmt.Println("    uint256 totalSupply;           // slot 0 (full slot)")
		fmt.Println("    address owner;                 // slot 1 (20 bytes)")
		fmt.Println("    bool paused;                   // slot 1 (packed, 1 byte)")
		fmt.Println("    mapping(address => uint256) balances;  // slot 2 (base)")
		fmt.Println("    uint256[] values;              // slot 3 (length)")
		fmt.Println("  }")
		fmt.Println()

		fmt.Println("Slot Calculations:")
		fmt.Println("  - totalSupply → read slot 0 directly")
		fmt.Println("  - owner → read slot 1, extract bytes 0-19")
		fmt.Println("  - paused → read slot 1, extract byte 20")
		fmt.Println("  - balances[addr] → read keccak256(addr || 0x02)")
		fmt.Println("  - values.length → read slot 3")
		fmt.Println("  - values[i] → read keccak256(0x03) + i")
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: This information helps you read any contract's storage
		// DEBUG: Understanding layout is essential for storage proofs (module 12)
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
	// 4. Use F11 (Step Into) to enter storage package functions
	// 5. Watch Variables panel to see storage values
	// 6. Inspect ResolvedSlot to see computed mapping slots
	// 7. Use Call Stack panel to see function hierarchy
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - Storage Slots: How does Solidity allocate storage?
	// - Merkle-Patricia Trie: How is storage organized?
	// - Mapping Slots: Why use keccak256 for slot calculation?
	// - Storage Root: How does it commit to all storage?
	// - Defensive Copying: Why return copies instead of references?
	//
	// STEPPING INTO STORAGE PACKAGE:
	// - Set breakpoint at storage.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to exercise.go or solution.reference.go
	// - Step through to see: validation, slot conversion, mapping calculation
	// - Watch how keccak256 is computed for mapping slots
	// - See StorageAt RPC call and response handling
	//
	// COMMON STORAGE ERRORS TO DEBUG:
	// - Zero contract address: Must provide valid contract
	// - Nil slot: Must specify which slot to read
	// - Historical block unavailable: Need archive node for old blocks
	// - Wrong mapping slot: Verify base slot number from contract source
}
