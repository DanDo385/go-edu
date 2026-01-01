package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/example/go-10x-minis/geth/13-trace/internal/trace"
)

/*
===================================================================
Ethereum Transaction Tracer
===================================================================

This program demonstrates how to trace transaction execution to see
opcode-level details, gas usage, and call trees.

USAGE:
    go run ./cmd/app/main.go <RPC_URL> <TX_HASH>
    go run ./cmd/app/main.go <RPC_URL> <TX_HASH> [demo]

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    <TX_HASH>       - Transaction hash to trace (required)
    [demo]          - Demo mode: basic, pretty (default: "basic")

Examples:
    go run ./cmd/app/main.go https://eth.llamarpc.com 0x123...
    go run ./cmd/app/main.go https://eth.llamarpc.com 0x123... pretty

    # Example transaction (Uniswap swap):
    go run ./cmd/app/main.go https://eth.llamarpc.com \
      0x5c504ed432cb51138bcf09aa5e8a410dd4a1e204ef84bfed1be16dfba1b22060

Common RPC Endpoints:
    Ethereum Mainnet:
      - https://eth.llamarpc.com
      - https://rpc.ankr.com/eth
      - https://eth-mainnet.public.blastapi.io

    Note: Transaction tracing requires debug_traceTransaction RPC method
    Many public endpoints disable this (too expensive to run)
    You may need your own node or a debug-enabled endpoint

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: cmd/app")
- Step Into (F11) to enter trace package functions
- Watch Variables panel to see trace data
- Inspect JSON structure for execution details

PACKAGE STRUCTURE:
- This file (cmd/app/main.go) is package main
- It imports the trace package from internal/trace
- exercise.go, solution.reference.go are in package trace
- The trace.Run() function encapsulates all tracing logic

KEY CONCEPTS:
- Transaction Tracing: Replay execution and record every operation
- Opcode-Level Details: See every EVM instruction that was executed
- Gas Accounting: Track gas consumption per operation
- Call Trees: Understand which contracts called which
- Deterministic Replay: Same inputs always produce same trace
- Execution Instrumentation: Observe without changing behavior
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	// DEBUG: os.Args[0] is the program name
	// DEBUG: os.Args[1] is the RPC URL (required)
	// DEBUG: os.Args[2] is the transaction hash (required)
	// DEBUG: os.Args[3] is the demo mode (optional)

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> <TX_HASH> [demo]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com 0x123...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com 0x123... pretty\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: basic, pretty\n")
		fmt.Fprintf(os.Stderr, "\nNote: Tracing requires debug_traceTransaction (may not be available on all endpoints)\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	txHashStr := os.Args[2]
	demo := "basic"
	if len(os.Args) > 3 {
		demo = os.Args[3]
	}

	// Parse transaction hash
	// BREAKPOINT: Set breakpoint here
	// DEBUG: Check that txHashStr is a valid hex string
	txHash := common.HexToHash(txHashStr)
	if txHash == (common.Hash{}) {
		log.Fatalf("Invalid transaction hash: %s\n", txHashStr)
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL', 'txHash', and 'demo' variables in Variables panel
	fmt.Printf("=== Ethereum Transaction Tracer ===\n")
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Transaction: %s\n", txHash.Hex())
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
	// debug methods like TraceTransaction.

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

	// Create gethclient for TraceTransaction access
	// BREAKPOINT: Set breakpoint here
	// DEBUG: gethclient wraps ethclient to expose debug methods
	client := gethclient.New(ethClient.Client())

	// BREAKPOINT: Set breakpoint here after successful connection
	// DEBUG: Connection is established—client is ready to use
	fmt.Println("✓ Connected to Ethereum RPC endpoint")
	fmt.Println()

	// ============================================================================
	// STEP 3: Trace Transaction
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through tracing
	// DEBUG: This demonstrates tracing a transaction execution
	// DEBUG: Watch the 'result' variable to see the raw trace data
	//
	// CONCEPT: Transaction tracing replays the transaction in the EVM and
	// records every operation. The trace includes:
	//   - Opcode execution sequence
	//   - Gas consumption per operation
	//   - Call tree (which contracts were called)
	//   - Storage changes
	//   - Return values and revert reasons
	//
	// This is expensive! It requires re-executing the entire transaction.
	// Many public RPC providers disable this endpoint.

	fmt.Println("=== Tracing Transaction ===")
	fmt.Println("Note: This may take several seconds for complex transactions...")
	fmt.Println()

	// Create configuration for tracing
	// BREAKPOINT: Set breakpoint here
	// DEBUG: Notice cfg.TxHash is the transaction we want to trace
	cfg := trace.Config{
		TxHash: txHash,
	}

	// BREAKPOINT: Set breakpoint here before calling trace.Run
	// DEBUG: Step Into (F11) to see the implementation in solution.reference.go
	// DEBUG: You'll see how the function calls TraceTransaction RPC method
	// DEBUG: Increase timeout for tracing—it can be slow for complex transactions
	runCtx, runCancel := context.WithTimeout(ctx, 30*time.Second)
	defer runCancel()

	result, err := trace.Run(runCtx, client, cfg)
	if err != nil {
		log.Fatalf("Failed to trace transaction: %v\n", err)
	}

	// BREAKPOINT: Set breakpoint here after successful trace
	// DEBUG: Expand 'result' in Variables panel to inspect fields
	// DEBUG: Look at TxHash (confirms which transaction was traced)
	// DEBUG: Look at Trace (raw JSON containing execution details)

	fmt.Printf("Transaction: %s\n", result.TxHash.Hex())
	fmt.Printf("Trace Size: %d bytes\n", len(result.Trace))
	fmt.Println()

	// ============================================================================
	// STEP 4: Display Trace Data
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to see trace display
	// DEBUG: Different display modes show different levels of detail

	if demo == "pretty" {
		// Pretty-print the JSON for human readability
		// BREAKPOINT: Set breakpoint here
		// DEBUG: This formats the JSON with indentation
		fmt.Println("=== Trace Data (Pretty-Printed) ===")
		fmt.Println()

		var prettyTrace interface{}
		err := json.Unmarshal(result.Trace, &prettyTrace)
		if err != nil {
			log.Fatalf("Failed to parse trace JSON: %v\n", err)
		}

		prettyJSON, err := json.MarshalIndent(prettyTrace, "", "  ")
		if err != nil {
			log.Fatalf("Failed to format trace JSON: %v\n", err)
		}

		fmt.Println(string(prettyJSON))
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: The pretty-printed JSON shows the full trace structure
		// DEBUG: Look for "gas", "stack", "memory", "storage" fields
		// DEBUG: See the opcode execution sequence
	} else {
		// Basic display—just show it's valid JSON
		// BREAKPOINT: Set breakpoint here
		fmt.Println("=== Trace Data (Summary) ===")
		fmt.Println()

		var traceData map[string]interface{}
		err := json.Unmarshal(result.Trace, &traceData)
		if err != nil {
			log.Fatalf("Failed to parse trace JSON: %v\n", err)
		}

		fmt.Printf("Trace contains %d top-level fields\n", len(traceData))
		fmt.Println("Top-level fields:")
		for key := range traceData {
			fmt.Printf("  - %s\n", key)
		}
		fmt.Println()

		fmt.Println("To see full trace, run with 'pretty' mode:")
		fmt.Printf("  %s %s %s pretty\n", os.Args[0], rpcURL, txHash.Hex())
		fmt.Println()

		// BREAKPOINT: Set breakpoint here
		// DEBUG: The trace JSON structure varies by tracer type
		// DEBUG: Default tracer returns opcode-level details
	}

	// ============================================================================
	// STEP 5: Understanding Trace Data
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to understand trace concepts
	// DEBUG: This demonstrates what information is in a trace

	fmt.Println("=== Understanding Transaction Traces ===")
	fmt.Println()

	fmt.Println("What's in a Trace?")
	fmt.Println("  1. Opcode Execution: Every EVM instruction that ran")
	fmt.Println("  2. Gas Usage: How much gas each operation consumed")
	fmt.Println("  3. Call Tree: Which contracts called which (CALL, DELEGATECALL, etc.)")
	fmt.Println("  4. Stack State: EVM stack contents at each step")
	fmt.Println("  5. Memory State: EVM memory contents at each step")
	fmt.Println("  6. Storage Changes: Which storage slots were modified")
	fmt.Println("  7. Return Values: Data returned from contract calls")
	fmt.Println("  8. Revert Reasons: Why execution failed (if it did)")
	fmt.Println()

	fmt.Println("Why Trace Transactions?")
	fmt.Println("  - Debugging: See exactly where a transaction failed")
	fmt.Println("  - Gas Optimization: Identify expensive operations")
	fmt.Println("  - Security Audits: Understand call flow and state changes")
	fmt.Println("  - Block Explorers: Show users what happened in a transaction")
	fmt.Println("  - MEV Analysis: Understand arbitrage and liquidation bots")
	fmt.Println()

	fmt.Println("Trace vs Receipt:")
	fmt.Println("  - Receipt: High-level summary (status, gas used, logs)")
	fmt.Println("  - Trace: Detailed execution log (every opcode, every state change)")
	fmt.Println("  - Receipt: Always available, lightweight (<1 KB)")
	fmt.Println("  - Trace: Only on debug nodes, heavyweight (can be MBs)")
	fmt.Println()

	fmt.Println("Tracer Types:")
	fmt.Println("  - Default: Full opcode-level trace (most detailed)")
	fmt.Println("  - callTracer: Simplified call tree (just contract calls)")
	fmt.Println("  - prestateTracer: Account state before transaction")
	fmt.Println("  - 4byteTracer: Function selectors called")
	fmt.Println("  - Custom: Define your own tracer in JavaScript")
	fmt.Println()

	// BREAKPOINT: Set breakpoint here
	// DEBUG: Understanding trace data is essential for debugging contracts
	// DEBUG: Different tracers provide different views of execution

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter trace package functions
	// 5. Watch Variables panel to see trace data
	// 6. Inspect result.Trace to see raw JSON
	// 7. Use Call Stack panel to see function hierarchy
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - EVM Opcodes: What instructions can the EVM execute?
	// - Call Types: CALL vs DELEGATECALL vs STATICCALL vs CREATE
	// - Gas Accounting: How is gas consumed and refunded?
	// - Execution Context: How do contracts access msg.sender, msg.value, etc.?
	// - Deterministic Replay: Why do traces always match actual execution?
	//
	// STEPPING INTO TRACE PACKAGE:
	// - Set breakpoint at trace.Run() call
	// - Press F11 (Step Into)
	// - Debugger will jump to exercise.go or solution.reference.go
	// - Step through to see: validation, TraceTransaction call, defensive copying
	// - Watch how JSON data is handled (json.RawMessage)
	// - See why we copy the byte slice (defensive copying pattern)
	//
	// COMMON TRACING ERRORS TO DEBUG:
	// - Zero transaction hash: Must provide valid transaction
	// - RPC endpoint doesn't support debug_traceTransaction: Need debug node
	// - Transaction not found: Verify transaction exists on this network
	// - Timeout: Complex transactions take long to trace, increase timeout
	// - Archive node required: Old transactions need historical state
}
