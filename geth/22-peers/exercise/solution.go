//go:build solution
// +build solution

package exercise

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
	// ============================================================================
	// This validation pattern is identical to modules 01 and 21. By now, you
	// should recognize this as a standard prelude to any function that accepts
	// external inputs.
	//
	// Why repeat this pattern? Consistency across the codebase makes code
	// predictable and maintainable. Every function that takes a context and
	// interface should validate both.
	//
	// Building on modules 01 and 21: Same pattern, different RPC operation.
	// This demonstrates how defensive programming applies universally.
	if ctx == nil {
		ctx = context.Background()
	}

	if client == nil {
		return nil, errors.New("client is nil")
	}

	// ============================================================================
	// STEP 2: Query Peer Count - Understanding P2P Network Health
	// ============================================================================
	// The PeerCount RPC call queries the node's P2P layer to see how many active
	// peer connections exist. This is a fundamental metric for network health.
	//
	// How P2P networking works in Ethereum:
	//   1. Discovery: Nodes find each other using DHT (Kademlia) and DNS discovery
	//   2. Connection: Nodes establish TCP connections with discovered peers
	//   3. Protocol negotiation: Peers agree on which protocols to speak (eth/67, snap, etc.)
	//   4. Gossip: Nodes exchange transactions, blocks, and state data
	//
	// Why peer count matters:
	//   - 0 peers: Node is completely isolated (firewall, network issue, or no discovery)
	//   - Low peers (1-5): Vulnerable to network partitions, slow data propagation
	//   - Moderate peers (10-50): Normal, healthy connectivity
	//   - High peers (100+): Excellent redundancy, faster block/tx propagation
	//
	// The net_peerCount RPC method:
	//   - Returns a hex string in raw JSON-RPC (e.g., "0x1a" for 26 peers)
	//   - The ethclient library automatically decodes hex → uint64
	//   - This is a lightweight call (just counts in-memory connections)
	//
	// Error handling: Network issues or node restarts can cause this call to fail.
	// We wrap the error with context for debugging.
	//
	// Building on module 01: Same error wrapping pattern (%w verb) for error chains.
	peerCount, err := client.PeerCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("peer count: %w", err)
	}

	// ============================================================================
	// STEP 3: Interpret Result - Understanding Limitations
	// ============================================================================
	// We now have the peer count, but it's important to understand what this
	// number tells us and what it doesn't.
	//
	// What peer count tells you:
	//   - Basic connectivity health (connected vs isolated)
	//   - Rough measure of network redundancy
	//   - Whether the node can gossip data at all
	//
	// What peer count DOESN'T tell you:
	//   - Peer quality (latency, bandwidth, honesty)
	//   - Which networks peers are on (mainnet vs testnet)
	//   - Which client implementations peers are running
	//   - Geographic distribution of peers
	//
	// For richer peer information, you need the admin_peers RPC endpoint:
	//   - Requires admin API enabled (typically IPC or authenticated HTTP)
	//   - Returns array of peer objects with:
	//     * Client name and version (geth, erigon, nethermind, etc.)
	//     * Network latency and bandwidth
	//     * Supported protocols and capabilities
	//     * Connection direction (inbound vs outbound)
	//
	// Why public RPCs hide this:
	//   - Privacy: Peer details could leak network topology
	//   - Security: Attackers could use peer info for targeted attacks
	//   - Most public RPC providers don't run full nodes (they proxy to backends)
	//
	// This is different from modules 01 and 21:
	//   - Module 01: Chain metadata (static, always available)
	//   - Module 21: Sync progress (dynamic, varies during sync)
	//   - Module 22: Network topology (dynamic, varies by node configuration)
	//
	// No validation needed: uint64 can't be nil or invalid. It's always a valid
	// count (even if zero means "no peers").

	// ============================================================================
	// STEP 4: Return Result - Primitive Type Handling
	// ============================================================================
	// We package the peer count into our Result struct. This is simpler than
	// modules 01 and 21 because we're dealing with a primitive type (uint64),
	// not pointers or complex structs.
	//
	// No defensive copying needed: uint64 is a primitive type that's copied by
	// value in Go. When we assign peerCount to Result.PeerCount, Go makes a
	// copy automatically. There's no shared state to worry about.
	//
	// This contrasts with module 01:
	//   - big.Int: Required defensive copying (mutable, pointer type)
	//   - Header: Required defensive copying (contains pointers/slices)
	//   - uint64: No copying needed (primitive, immutable by nature)
	//
	// API design: We use a Result struct even though it only has one field. This
	// provides consistency with other modules and room for future expansion (we
	// could add PeerQuality, RegionDistribution, etc. without breaking the API).
	//
	// Production usage:
	//   - Monitor peer count over time (sudden drops indicate issues)
	//   - Alert if peer count stays at 0 for extended periods
	//   - Compare peer count before/after network changes (firewall rules, etc.)
	//
	// Building on previous modules:
	//   - Module 01: Returned complex types (ChainID, NetworkID, Header)
	//   - Module 21: Returned boolean + optional struct (IsSyncing, Progress)
	//   - Module 22: Returns simple metric (PeerCount)
	//   - Pattern: All use Result structs for consistent API design
	return &Result{
		PeerCount: peerCount, // Primitive type, automatically copied by value
	}, nil
}
