//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

/*
Problem: Build a Swiss Army knife CLI that combines multiple node operations.

This capstone module brings together patterns from all previous modules into a single
unified tool. Instead of separate programs for each operation, you'll have one tool
with subcommands (like git, docker, kubectl).

This demonstrates:
  - Command routing and dispatch
  - Code reuse across modules
  - Building production-ready tools
  - Composing simple operations into complex workflows

Computer science principles highlighted:
  - Command pattern (encapsulating operations)
  - Composition (building complex from simple)
  - Interface segregation (ToolboxClient combines many interfaces)
*/
func Run(ctx context.Context, client ToolboxClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Foundation for All Commands
	// ============================================================================
	// Standard validation pattern from all previous modules.
	// Building on modules 01-24: Consistent validation across the entire course.
	if ctx == nil {
		ctx = context.Background()
	}

	if client == nil {
		return nil, errors.New("client is nil")
	}

	if cfg.Command == "" {
		return nil, errors.New("command is required")
	}

	// ============================================================================
	// STEP 2: Command Routing - Dispatch Pattern
	// ============================================================================
	// The command pattern encapsulates operations as objects (or in Go, as
	// function calls). This allows us to route to different handlers based on
	// the command string.
	//
	// Why command routing?
	//   - Single entry point with multiple behaviors
	//   - Easy to add new commands without changing caller code
	//   - Similar to HTTP routing, CLI tools (git, docker), RPC methods
	//
	// Common routing patterns:
	//   - Switch statement (simple, fast, used here)
	//   - Map of command → handler function (more flexible)
	//   - Reflection-based routing (most dynamic, slower)
	//
	// This demonstrates composability: each command reuses patterns from
	// previous modules, combining them in new ways.
	switch cfg.Command {
	case "status":
		return handleStatus(ctx, client)
	case "block":
		return handleBlock(ctx, client, cfg.Args)
	case "tx":
		return handleTx(ctx, client, cfg.Args)
	default:
		return nil, fmt.Errorf("unknown command: %s (valid: status, block, tx)", cfg.Command)
	}
}

// ============================================================================
// STATUS COMMAND - Comprehensive Node Overview
// ============================================================================
// Combines patterns from modules 01, 21, 22 into single command.
// This demonstrates how to compose simple operations into complex ones.
func handleStatus(ctx context.Context, client ToolboxClient) (*Result, error) {
	// Fetch chain metadata (module 01)
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}

	networkID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("network id: %w", err)
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("header: %w", err)
	}

	// Check sync status (module 21)
	progress, err := client.SyncProgress(ctx)
	if err != nil {
		return nil, fmt.Errorf("sync progress: %w", err)
	}

	// Check peer count (module 22)
	peerCount, err := client.PeerCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("peer count: %w", err)
	}

	// Compose comprehensive status
	status := map[string]interface{}{
		"chainID":     chainID.String(),
		"networkID":   networkID.String(),
		"blockNumber": header.Number.Uint64(),
		"blockHash":   header.Hash().Hex(),
		"syncing":     progress != nil,
		"peerCount":   peerCount,
	}

	return &Result{
		Command: "status",
		Output:  status,
		Status:  "success",
	}, nil
}

// ============================================================================
// BLOCK COMMAND - Block Inspection
// ============================================================================
// Fetches and displays block details. Demonstrates parsing arguments and
// fetching blockchain data.
func handleBlock(ctx context.Context, client ToolboxClient, args []string) (*Result, error) {
	if len(args) == 0 {
		return nil, errors.New("block command requires block number argument")
	}

	// Parse block number from string
	blockNum, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid block number: %w", err)
	}

	// Fetch block
	block, err := client.BlockByNumber(ctx, big.NewInt(int64(blockNum)))
	if err != nil {
		return nil, fmt.Errorf("fetch block: %w", err)
	}

	if block == nil {
		return nil, fmt.Errorf("block %d not found", blockNum)
	}

	// Extract relevant details
	blockData := map[string]interface{}{
		"number":     block.Number().Uint64(),
		"hash":       block.Hash().Hex(),
		"parentHash": block.ParentHash().Hex(),
		"timestamp":  block.Time(),
		"txCount":    len(block.Transactions()),
		"gasUsed":    block.GasUsed(),
		"gasLimit":   block.GasLimit(),
	}

	return &Result{
		Command: "block",
		Output:  blockData,
		Status:  "success",
	}, nil
}

// ============================================================================
// TX COMMAND - Transaction Inspection
// ============================================================================
// Fetches and displays transaction details. Demonstrates tx hash parsing.
func handleTx(ctx context.Context, client ToolboxClient, args []string) (*Result, error) {
	if len(args) == 0 {
		return nil, errors.New("tx command requires transaction hash argument")
	}

	txHash := args[0]

	// Fetch transaction
	tx, pending, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("fetch transaction: %w", err)
	}

	if tx == nil {
		return nil, fmt.Errorf("transaction %s not found", txHash)
	}

	// Extract relevant details
	txData := map[string]interface{}{
		"hash":     tx.Hash().Hex(),
		"nonce":    tx.Nonce(),
		"value":    tx.Value().String(),
		"gas":      tx.Gas(),
		"gasPrice": tx.GasPrice().String(),
		"pending":  pending,
	}

	// Add To address if present (contract creation txs have nil To)
	if tx.To() != nil {
		txData["to"] = tx.To().Hex()
	}

	return &Result{
		Command: "tx",
		Output:  txData,
		Status:  "success",
	}, nil
}
