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
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/example/go-10x-minis/geth/12-proofs/internal/proofs"
)

/*
===================================================================
Ethereum Merkle-Patricia Trie Proof Fetcher
===================================================================

This program demonstrates how to fetch cryptographic proofs for account
state and storage slots, enabling trust-minimized verification.

USAGE:
    go run ./cmd/app/main.go <RPC_URL>
    go run ./cmd/app/main.go <RPC_URL> <DEMO>

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    [demo]          - Demo to run: all, account, storage, usdc, multi (default: "all")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com
    go run ./cmd/app/main.go https://eth.llamarpc.com account
    go run ./cmd/app/main.go https://eth.llamarpc.com storage
    go run ./cmd/app/main.go https://eth.llamarpc.com usdc
    go run ./cmd/app/main.go https://eth.llamarpc.com multi

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
- Step Into (F11) to enter proofs package functions
- Watch Variables panel to see proof structures
- Inspect proof nodes and account data

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the proofs package from internal/proofs
- exercise.go, solution.reference.go are in package proofs
- The proofs.Run() function encapsulates all proof fetching logic

KEY CONCEPTS:
- Merkle-Patricia Trie: Cryptographic tree structure for Ethereum state
- Proof Nodes: Siblings along path from leaf to root
- State Root: Root hash in block header committing to all accounts
- Storage Root: Root hash in account committing to all storage slots
- Light Clients: Verify state without full blockchain data
- Trust-Minimized Bridges: Prove state from one chain to another
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
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com account\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com storage\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, account, storage, usdc, multi\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	demo := "all"
	if len(os.Args) > 2 {
		demo = os.Args[2]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL' and 'demo' variables in Variables panel
	fmt.Printf("=== Ethereum Merkle-Patricia Trie Proof Fetcher ===\n")
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
	// the RPC methods. We wrap it with gethclient.New to access Geth-specific
	// methods like GetProof.

	ctx := context.Background()
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// BREAKPOINT: Set breakpoint here
	// DEBUG: Step Over (F10) to execute the Dial call
	// DEBUG: If it fails, inspect 'err' to see the connection error
	ethClient, err := ethclient.Dial(dialCtx, rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to RPC endpoint: %v\n", err)
	}
	defer ethClient.Close()

	// Create gethclient for GetProof access
	// BREAKPOINT: Set breakpoint here
	// DEBUG: gethclient wraps ethclient to expose debug/geth-specific methods
	client := gethclient.New(ethClient.Client())

	// BREAKPOINT: Set breakpoint here after successful connection
	// DEBUG: Connection is established—client is ready to use
	fmt.Println("✓ Connected to Ethereum RPC endpoint")
	fmt.Println()

	// ============================================================================
	// STEP 3: Example 1 - Fetch Account Proof (No Storage)
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through first example
	// DEBUG: This demonstrates fetching an account proof without storage slots
	// DEBUG: Watch the 'result' variable to see account state and proof nodes
	//
	// CONCEPT: Account proofs prove that an account exists with specific state:
	//   - Balance: How much ETH the account has
	//   - Nonce: Number of transactions sent (for EOAs) or contracts created (for contracts)
	//   - CodeHash: Hash of contract bytecode (keccak256("") for EOAs)
	//   - StorageHash: Root of the account's storage trie
	//
	// The proof nodes form a path from the state root (in block header) to
	// the account leaf in the state trie. By hashing up this path, you can
	// verify the account data is part of the committed state.

	if demo == "all" || demo == "account" {
		fmt.Println("=== Example 1: Account Proof (No Storage) ===")
		fmt.Println("Fetching proof for Vitalik's account...")
		fmt.Println()

		// Vitalik's address as an example
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice we're using a well-known address
		vitalikAddress := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

		// Create configuration for account-only proof
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice Slots is nil—we only want account proof, no storage
		// DEBUG: BlockNumber is nil—fetching proof from latest block
		cfg := proofs.Config{
			Account:     vitalikAddress,
			Slots:       nil, // No storage slots requested
			BlockNumber: nil, // Latest block
		}

		// BREAKPOINT: Set breakpoint here before calling proofs.Run
		// DEBUG: Step Into (F11) to see the implementation in solution.reference.go
		// DEBUG: You'll see how the function converts slots to strings and calls GetProof
		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		result, err := proofs.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to fetch account proof: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful fetch
		// DEBUG: Expand 'result' in Variables panel to inspect all fields
		// DEBUG: Look at Account.Balance, Account.Nonce, Account.CodeHash
		// DEBUG: Count the proof nodes—this is the depth of the trie!

		fmt.Printf("Account: %s\n", vitalikAddress.Hex())
		fmt.Printf("Balance: %s wei (%.6f ETH)\n",
			result.Account.Balance.String(),
			weiToEth(result.Account.Balance))
		fmt.Printf("Nonce: %d\n", result.Account.Nonce)
		fmt.Printf("Code Hash: %s\n", result.Account.CodeHash.Hex())
		fmt.Printf("Storage Hash: %s\n", result.Account.StorageHash.Hex())
		fmt.Printf("Proof Nodes: %d\n", len(result.Account.ProofNodes))
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: The proof nodes array contains RLP-encoded trie nodes
		// DEBUG: Each node is a step in the path from state root to account
		// DEBUG: Light clients use these to verify without full state
	}

	// ============================================================================
	// STEP 4: Example 2 - Fetch Storage Proof
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through storage proof example
	// DEBUG: This demonstrates fetching proofs for specific storage slots
	// DEBUG: Storage proofs are nested inside account proofs!
	//
	// CONCEPT: Storage proofs work in two levels:
	//   1. Account proof: Proves the account exists and has storageHash X
	//   2. Storage proof: Proves slot Y has value Z in storage trie with root X
	//
	// This two-level structure allows proving both account existence and
	// storage values in a single RPC call.

	if demo == "all" || demo == "storage" {
		fmt.Println("=== Example 2: Storage Proof ===")
		fmt.Println("Fetching proof for USDC totalSupply storage slot...")
		fmt.Println()

		// USDC contract on Ethereum mainnet
		// BREAKPOINT: Set breakpoint here
		usdcAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

		// USDC totalSupply is stored at slot 1
		// We need to provide the slot as common.Hash
		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice we're requesting storage slot 1 (totalSupply)
		slot1 := common.BigToHash(big.NewInt(1))
		slots := []common.Hash{slot1}

		cfg := proofs.Config{
			Account:     usdcAddress,
			Slots:       slots,
			BlockNumber: nil, // Latest block
		}

		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		// BREAKPOINT: Set breakpoint here before calling proofs.Run
		// DEBUG: Step Into (F11) to see how storage proofs are processed
		result, err := proofs.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to fetch storage proof: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here after successful fetch
		// DEBUG: Expand result.Account.Storage to see storage proofs
		// DEBUG: Each storage proof has Key, Value, and ProofNodes

		fmt.Printf("Contract: %s\n", usdcAddress.Hex())
		fmt.Printf("Account Proof Nodes: %d\n", len(result.Account.ProofNodes))
		fmt.Println()

		// Display storage proofs
		if len(result.Account.Storage) > 0 {
			fmt.Println("Storage Proofs:")
			for i, sp := range result.Account.Storage {
				fmt.Printf("  [%d] Slot: %s\n", i, sp.Key.Hex())
				fmt.Printf("      Value: %s\n", sp.Value.String())
				fmt.Printf("      Proof Nodes: %d\n", len(sp.ProofNodes))

				// Interpret as totalSupply (USDC has 6 decimals)
				if sp.Value.Cmp(big.NewInt(0)) > 0 {
					divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)
					usdc := new(big.Int).Div(sp.Value, divisor)
					fmt.Printf("      Total Supply: %s USDC\n", usdc.String())
				}
			}
			fmt.Println()
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Notice the two-level proof structure
		// DEBUG: Account proof proves storageHash
		// DEBUG: Storage proof proves value at slot using that storageHash
	}

	// ============================================================================
	// STEP 5: Example 3 - Multiple Storage Slots
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through multiple slots example
	// DEBUG: This demonstrates fetching proofs for multiple slots at once
	// DEBUG: More efficient than making separate calls for each slot

	if demo == "all" || demo == "multi" {
		fmt.Println("=== Example 3: Multiple Storage Slots ===")
		fmt.Println("Fetching proofs for multiple USDC storage slots...")
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		usdcAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

		// Request multiple slots at once
		// BREAKPOINT: Set breakpoint here
		// DEBUG: We're requesting 3 different slots in one RPC call
		slots := []common.Hash{
			common.BigToHash(big.NewInt(0)), // Slot 0 (often implementation address in proxies)
			common.BigToHash(big.NewInt(1)), // Slot 1 (totalSupply in USDC)
			common.BigToHash(big.NewInt(2)), // Slot 2 (name in USDC)
		}

		cfg := proofs.Config{
			Account:     usdcAddress,
			Slots:       slots,
			BlockNumber: nil,
		}

		runCtx, runCancel := context.WithTimeout(ctx, 10*time.Second)
		defer runCancel()

		result, err := proofs.Run(runCtx, client, cfg)
		if err != nil {
			log.Fatalf("Failed to fetch multi-slot proof: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: result.Account.Storage contains proofs for all requested slots
		// DEBUG: Each slot has its own independent proof

		fmt.Printf("Contract: %s\n", usdcAddress.Hex())
		fmt.Printf("Requested Slots: %d\n", len(slots))
		fmt.Printf("Storage Proofs Received: %d\n", len(result.Account.Storage))
		fmt.Println()

		for i, sp := range result.Account.Storage {
			fmt.Printf("Slot %d:\n", i)
			fmt.Printf("  Key: %s\n", sp.Key.Hex())
			fmt.Printf("  Value: %s\n", sp.Value.String())
			fmt.Printf("  Proof Nodes: %d\n", len(sp.ProofNodes))
		}
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Batching multiple slots is more efficient than separate calls
		// DEBUG: All proofs share the same account proof (no duplication)
	}

	// ============================================================================
	// STEP 6: Example 4 - Understanding Proof Structure
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand proof structure
	// DEBUG: This demonstrates the cryptographic commitment structure

	if demo == "all" {
		fmt.Println("=== Example 4: Understanding Proof Structure ===")
		fmt.Println()

		fmt.Println("Merkle-Patricia Trie Proof Concepts:")
		fmt.Println("  1. State Trie: Global trie containing all accounts")
		fmt.Println("  2. Storage Trie: Per-contract trie containing storage slots")
		fmt.Println("  3. Proof Nodes: Siblings along path from leaf to root")
		fmt.Println("  4. Root Hash: Cryptographic commitment to entire trie")
		fmt.Println()

		fmt.Println("Two-Level Proof Structure:")
		fmt.Println("  Step 1: Prove account exists in state trie")
		fmt.Println("    - Start with state root (from block header)")
		fmt.Println("    - Hash up account proof nodes")
		fmt.Println("    - Verify final hash matches state root")
		fmt.Println("    - Extract storageHash from account data")
		fmt.Println()
		fmt.Println("  Step 2: Prove storage slot in storage trie")
		fmt.Println("    - Start with storageHash (from step 1)")
		fmt.Println("    - Hash up storage proof nodes")
		fmt.Println("    - Verify final hash matches storageHash")
		fmt.Println("    - Extract slot value from storage data")
		fmt.Println()

		fmt.Println("Why This Matters:")
		fmt.Println("  - Light Clients: Verify state without full blockchain (GB → KB)")
		fmt.Println("  - Bridges: Prove state on one chain to another's smart contract")
		fmt.Println("  - Indexers: Verify indexed data matches on-chain state")
		fmt.Println("  - Security: Trust-minimized verification (cryptographic proof)")
		fmt.Println()

		fmt.Println("Proof Size:")
		fmt.Println("  - Logarithmic: O(log N) nodes for N accounts/slots")
		fmt.Println("  - Typical: 6-10 nodes per proof (~500-1500 bytes)")
		fmt.Println("  - Compare: Full state is 100+ GB, proof is <2 KB!")
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: This explains why proofs are so powerful for scaling
		// DEBUG: Logarithmic proof size enables light clients and bridges
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
	// 4. Use F11 (Step Into) to enter proofs package functions
	// 5. Watch Variables panel to see proof structures
	// 6. Inspect Account.ProofNodes and Storage[i].ProofNodes
	// 7. Use Call Stack panel to see function hierarchy
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - Merkle-Patricia Trie: How does it combine Merkle tree and Patricia trie?
	// - Proof Verification: How do you verify a proof against a root hash?
	// - Two-Level Structure: Why separate account and storage tries?
	// - Light Clients: How do they verify state with only proofs?
	// - Cross-Chain Bridges: How do they prove state from one chain to another?
	//
	// STEPPING INTO PROOFS PACKAGE:
	// - Set breakpoint at proofs.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to exercise.go or solution.reference.go
	// - Step through to see: validation, slot conversion, GetProof call
	// - Watch how proof response is processed and copied
	// - See defensive copying for mutable data (big.Int, slices)
	//
	// COMMON PROOF ERRORS TO DEBUG:
	// - Zero account address: Must provide valid account
	// - RPC endpoint doesn't support eth_getProof: Need compatible node
	// - Proof verification failure: Root hash mismatch or corrupted proof
	// - Historical block unavailable: Need archive node for old proofs
}

// weiToEth converts wei (big.Int) to ETH (float64) for display
func weiToEth(wei *big.Int) float64 {
	// 1 ETH = 10^18 wei
	ethDivisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	weiBig := new(big.Float).SetInt(wei)
	ethBig := new(big.Float).Quo(weiBig, ethDivisor)
	eth, _ := ethBig.Float64()
	return eth
}
