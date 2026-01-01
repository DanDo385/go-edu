//go:build !solution && !reference

package peers

import (
	"context"
	"errors"
	"fmt"
)

/*
Problem: Query the number of connected peers to assess node connectivity health.

In Ethereum's peer-to-peer network, nodes connect to other nodes (peers) to gossip
transactions and blocks. The number of connected peers is a basic health indicator:
too few peers means slow propagation of data, while zero peers means complete isolation.

The net_peerCount RPC method returns a hexadecimal string representing the count,
which the ethclient library automatically converts to uint64.

Computer science principles highlighted:
  - P2P network topology (decentralized mesh network)
  - Gossip protocols (how information spreads)
  - Health metrics and observability (monitoring system state)
*/
func Run(ctx context.Context, client PeerClient, cfg Config) (*Result, error) {
	// ============================================================================
	// STEP 1: Input Validation - Repeating Defensive Pattern
	// TODO: Implement

	// ============================================================================
	// This validation pattern is identical to modules 01 and 21. By now, you
	// TODO: Implement

	// ============================================================================
	// STEP 2: Query Peer Count - Understanding P2P Network Health
	// TODO: Implement

	// ============================================================================
	// The PeerCount RPC call queries the node's P2P layer to see how many active
	// TODO: Implement

	// ============================================================================
	// STEP 3: Interpret Result - Understanding Limitations
	// TODO: Implement

	// ============================================================================
	// We now have the peer count, but it's important to understand what this
	// TODO: Implement

	// ============================================================================
	// STEP 4: Return Result - Primitive Type Handling
	// TODO: Implement

	// ============================================================================
	// We package the peer count into our Result struct. This is simpler than
	// TODO: Implement

	panic("unimplemented")
}
