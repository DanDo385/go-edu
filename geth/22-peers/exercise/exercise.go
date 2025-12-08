//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"errors"
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
	// TODO: Validate input parameters
	// - Check if ctx is nil and provide a default context if needed
	// - Check if client is nil and return an appropriate error
	// - Why: Standard defensive programming pattern from module 01

	// TODO: Call PeerCount RPC method
	// - Call client.PeerCount(ctx) to get the number of connected peers
	// - Handle potential errors from the RPC call
	// - Why: This tells us how well connected the node is to the network
	// - Key concept: Peer count is a basic connectivity health indicator
	// - Typical values:
	//   * 0 peers: Node is isolated, cannot sync or propagate data
	//   * 1-10 peers: Low connectivity, vulnerable to network partitions
	//   * 25-50 peers: Good connectivity (typical default max)
	//   * 100+ peers: Excellent connectivity (if configured for high peering)

	// TODO: Interpret the peer count value
	// - Peer count alone doesn't tell quality (some peers may be slow/malicious)
	// - For richer peer information, you need admin_peers (requires admin API access)
	// - admin_peers provides: client version, latency, capabilities, protocols
	// - Most public RPCs hide peer details for privacy/security
	// - Why: Understanding limitations helps set realistic expectations

	// TODO: Construct and return the Result struct
	// - Store the peer count in the Result
	// - Return the result with nil error on success
	// - No defensive copying needed: uint64 is a primitive type (copied by value)

	return nil, errors.New("not implemented")
}
