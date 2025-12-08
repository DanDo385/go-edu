//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"errors"
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
	// TODO: Validate input parameters
	// - Check if ctx is nil and provide a default context
	// - Check if client is nil and return an error
	// - Validate cfg.Command is not empty
	// - Why: Need to know which operation to perform

	// TODO: Route to appropriate command handler
	// - Implement command routing based on cfg.Command
	// - Supported commands:
	//   * "status": Show node status (chain ID, network ID, sync status, peer count)
	//   * "block": Fetch and display block details
	//   * "tx": Fetch and display transaction details
	//   * "health": Check node health (from module 24)
	// - Why: Single tool, multiple operations (like git has status, commit, push, etc.)

	// TODO: Implement "status" command
	// - Reuse patterns from modules 01, 21, 22
	// - Fetch: ChainID, NetworkID, latest header, sync progress, peer count
	// - Combine into comprehensive status report
	// - Why: Single command gives complete node overview

	// TODO: Implement "block" command
	// - Parse block number from cfg.Args
	// - Fetch block using BlockByNumber
	// - Display: number, hash, timestamp, tx count, gas used
	// - Why: Quick block inspection without full explorer

	// TODO: Implement "tx" command
	// - Parse tx hash from cfg.Args
	// - Fetch transaction using TransactionByHash
	// - Display: hash, from, to, value, gas, nonce
	// - Why: Verify transaction details from CLI

	// TODO: Handle unknown commands
	// - Return error if command not recognized
	// - Suggest valid commands in error message
	// - Why: User-friendly error handling

	// TODO: Construct and return unified Result
	// - Include command name, output data, and status
	// - Output structure varies by command (use interface{})
	// - Why: Flexible return type accommodates different command outputs

	return nil, errors.New("not implemented")
}
