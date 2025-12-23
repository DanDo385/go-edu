//go:build !solution
// +build !solution

package exercise

// TODO: Import required packages
// You'll need:
// - "fmt" for error formatting and string operations
// - "math/rand" for random peer selection (Fisher-Yates shuffle)
// - "sync" for RWMutex (concurrent node access)
// - "time" for message timestamps and simulated latency
//
// import (
//     "fmt"
//     "math/rand"
//     "sync"
//     "time"
// )

// ============================================================================
// P2P GOSSIP PROTOCOL: Distributed Message Propagation with Mock Network
// ============================================================================
//
// Gossip protocols (epidemic broadcast) are fundamental to distributed systems:
// - Used by: Cassandra (failure detection), Bitcoin/Ethereum (tx/block propagation)
// - Key property: Eventually consistent, probabilistically reliable
// - Scalability: O(log N) message rounds to reach all N nodes
// - Fault tolerance: Works even with node failures and network partitions
//
// How gossip works:
// 1. Node receives/creates a message
// 2. Randomly selects K peers (fanout)
// 3. Forwards message to selected peers
// 4. Each peer repeats steps 2-3 (epidemic spread)
// 5. Nodes track seen messages to avoid infinite loops
//
// Gossip vs Flooding:
// - Flooding: Send to ALL peers (inefficient, O(N²) messages)
// - Gossip: Send to K peers (efficient, O(N log N) messages)
//
// Memory Management:
// - Each message: ~200-500 bytes (ID, type, payload, timestamp)
// - Seen map per node: ~100-1000 entries = 10-50 KB per node
// - Network scales: 100 nodes × 50 KB = 5 MB total
//
// ============================================================================

import (
	"sync"
	"time"
)

// Message represents a gossip message in the network
type Message struct {
	ID        string                 // Unique message identifier (prevents duplicates)
	Type      string                 // Message type (e.g., "transaction", "block")
	From      string                 // Sender node ID
	Payload   map[string]interface{} // Message data (flexible key-value store)
	Timestamp time.Time              // Creation time (for TTL, ordering)
}

// GossipNode represents a node in the gossip network
// It maintains connections to peers and propagates messages
type GossipNode interface {
	// ID returns the unique identifier of this node
	ID() string

	// AddPeer adds a peer to this node's peer list
	AddPeer(peerID string)

	// Broadcast initiates a new message broadcast from this node
	// The message will be propagated to peers using the gossip protocol
	Broadcast(msgType string, payload map[string]interface{}) error

	// ReceiveMessage handles an incoming message from a peer
	// Returns true if the message was new (not seen before)
	ReceiveMessage(msg Message) bool

	// GetState returns a copy of the node's current state
	GetState() map[string]interface{}

	// Shutdown gracefully shuts down the node
	Shutdown()
}

// Network simulates a network that can deliver messages between nodes
// It can simulate latency, packet loss, and network partitions
type Network interface {
	// RegisterNode adds a node to the network
	RegisterNode(node GossipNode)

	// Send sends a message from one node to another
	// The delivery may be delayed or dropped based on network conditions
	Send(from, to string, msg Message)

	// SetLatency configures the network latency for message delivery
	SetLatency(latency time.Duration)

	// SetDropRate configures the packet drop rate (0.0 to 1.0)
	SetDropRate(rate float64)

	// GetMessageCount returns the total number of messages sent through the network
	GetMessageCount() int
}

// GossipProtocol defines the behavior of the gossip protocol
type GossipProtocol interface {
	// Fanout returns the number of peers to gossip to for each message
	Fanout() int

	// SelectPeers selects which peers to gossip a message to
	// It should return up to Fanout() peer IDs
	SelectPeers(allPeers []string, excludeID string) []string

	// ShouldForward determines if a message should be forwarded
	// This can be based on TTL, message age, or other criteria
	ShouldForward(msg Message) bool
}

// ============================================================================
// Exercise 1: Implement a basic gossip node
// ============================================================================

// TODO: Create a gossipNode struct that implements GossipNode
//
// Memory layout:
// - id: string (16 bytes header + len)
// - network: Network (8 bytes pointer)
// - fanout: int (8 bytes)
// - peers: []string (24 bytes slice header + peer data)
// - state: map[string]interface{} (~48 bytes per entry)
// - seen: map[string]bool (~40 bytes per entry)
// - mu: sync.RWMutex (16 bytes)
//
// Concurrency considerations:
// - Multiple goroutines may call Broadcast/ReceiveMessage simultaneously
// - Must protect peers, state, and seen maps with mutex
// - Use RWMutex for read-heavy operations (GetState, checking seen)
//
// type gossipNode struct {
//     id       string
//     network  Network
//     fanout   int
//     peers    []string
//     state    map[string]interface{}
//     seen     map[string]bool
//     mu       sync.RWMutex
//     shutdown chan struct{}
// }

// NewGossipNode creates a new gossip node.
// TODO: Implement node initialization
//
// Parameters:
// - id: Unique node identifier (e.g., "node-1", "node-2")
// - network: Network instance for sending messages
// - fanout: Number of random peers to forward each message to (typically 3-5)
//
// Algorithm:
// 1. Create gossipNode with provided id, network, fanout
// 2. Initialize empty peers slice
// 3. Initialize empty state and seen maps
// 4. Initialize shutdown channel
// 5. Return node (as GossipNode interface)
//
// Fanout selection guidelines:
// - Too low (1-2): Unreliable, messages may not reach all nodes
// - Optimal (3-5): Good balance of reliability and efficiency
// - Too high (>10): Approaches flooding, wastes bandwidth
//
// func NewGossipNode(id string, network Network, fanout int) GossipNode {
//     return &gossipNode{
//         id:       id,
//         network:  network,
//         fanout:   fanout,
//         peers:    make([]string, 0),
//         state:    make(map[string]interface{}),
//         seen:     make(map[string]bool),
//         shutdown: make(chan struct{}),
//     }
// }
func NewGossipNode(id string, network Network, fanout int) GossipNode {
	// TODO: Implement
	// - Create node struct with id, network, fanout, peers list, state map, seen map
	// - Initialize mutex for thread safety
	// - Return the node
	return nil
}

// ============================================================================
// Exercise 2: Implement a mock network
// ============================================================================

// TODO: Create a mockNetwork struct that implements Network
//
// The mock network simulates:
// - Message delivery with configurable latency (time.AfterFunc)
// - Packet loss based on drop rate (rand.Float64() < dropRate)
// - Message counting for statistics
//
// type mockNetwork struct {
//     mu           sync.RWMutex
//     nodes        map[string]GossipNode
//     latency      time.Duration
//     dropRate     float64
//     messageCount int
//     droppedCount int
// }

// NewMockNetwork creates a mock network with configurable latency and drop rate.
// TODO: Implement network initialization
//
// Parameters:
// - latency: Simulated network latency (e.g., 10ms, 100ms)
// - dropRate: Probability of packet loss (0.0 = no loss, 0.1 = 10% loss, 1.0 = all dropped)
//
// Realistic latency examples:
// - LAN: 1-5 ms
// - Same datacenter: 1-10 ms
// - Cross-datacenter (same region): 10-50 ms
// - Cross-continent: 100-300 ms
//
// func NewMockNetwork(latency time.Duration, dropRate float64) Network {
//     return &mockNetwork{
//         nodes:    make(map[string]GossipNode),
//         latency:  latency,
//         dropRate: dropRate,
//     }
// }
func NewMockNetwork(latency time.Duration, dropRate float64) Network {
	// TODO: Implement
	// - Create network struct with nodes map, latency, dropRate, message counter
	// - Initialize mutex for thread safety
	// - Return the network
	return nil
}

// ============================================================================
// Exercise 3: Implement a push-based gossip protocol
// ============================================================================

// PushGossipProtocol implements the "push" gossip variant.
// In push gossip, nodes actively send messages to peers.
//
// TODO: Add fields for fanout parameter
//
// Gossip variants:
// - Push: Sender pushes to random peers (fast initial spread)
// - Pull: Receiver pulls from random peers (slow, but catches missed messages)
// - Push-Pull: Hybrid approach (most robust)
//
// type PushGossipProtocol struct {
//     fanout int
// }

// NewPushGossipProtocol creates a new push gossip protocol.
// TODO: Implement protocol initialization
//
// func NewPushGossipProtocol(fanout int) GossipProtocol {
//     return &PushGossipProtocol{fanout: fanout}
// }
func NewPushGossipProtocol(fanout int) GossipProtocol {
	// TODO: Implement
	return nil
}

// SelectPeers selects random peers for message forwarding.
// TODO: Implement random peer selection
//
// Algorithm (Fisher-Yates shuffle):
// 1. Filter out excludeID from allPeers
// 2. If remaining peers <= fanout, return all
// 3. Otherwise, shuffle using Fisher-Yates:
//    for i from n-1 down to 1:
//        j = random integer 0 <= j <= i
//        swap shuffled[i] and shuffled[j]
// 4. Return first fanout peers
//
// Why Fisher-Yates?
// - Unbiased: Every permutation equally likely
// - Efficient: O(n) time, O(1) extra space
// - Simple: No duplicate checks needed
//
// func (p *PushGossipProtocol) SelectPeers(allPeers []string, excludeID string) []string {
//     // Filter out excluded peer
//     eligible := make([]string, 0, len(allPeers))
//     for _, peer := range allPeers {
//         if peer != excludeID {
//             eligible = append(eligible, peer)
//         }
//     }
//
//     // If we have fewer peers than fanout, return all
//     if p.fanout >= len(eligible) {
//         return eligible
//     }
//
//     // Fisher-Yates shuffle and take first fanout
//     shuffled := make([]string, len(eligible))
//     copy(shuffled, eligible)
//     for i := len(shuffled) - 1; i > 0; i-- {
//         j := rand.Intn(i + 1)
//         shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
//     }
//
//     return shuffled[:p.fanout]
// }
func (p *PushGossipProtocol) SelectPeers(allPeers []string, excludeID string) []string {
	// TODO: Implement
	// - Filter out excludeID from allPeers
	// - Shuffle the remaining peers
	// - Return up to Fanout() peers
	return nil
}

// ============================================================================
// Exercise 4: Implement convergence detection
// ============================================================================

// ConvergenceDetector checks if all nodes have converged (same state).
// Used for testing and verification.
//
// TODO: Implement convergence tracking
//
// type ConvergenceDetector struct {
//     nodes []GossipNode
//     mu    sync.RWMutex
// }

// NewConvergenceDetector creates a convergence detector for given nodes.
// TODO: Implement detector initialization
//
// func NewConvergenceDetector(nodes []GossipNode) *ConvergenceDetector {
//     return &ConvergenceDetector{nodes: nodes}
// }
func NewConvergenceDetector(nodes []GossipNode) *ConvergenceDetector {
	// TODO: Implement
	return nil
}

// IsConverged checks if all nodes have the same state for a given key.
// TODO: Implement convergence checking
//
// Algorithm:
// 1. Lock for reading
// 2. Get state from all nodes
// 3. Check if all nodes have the same value for 'key'
// 4. Return (converged bool, count of nodes with matching value)
//
// Convergence example:
// - 10 nodes, all have state["key"] = "value" → converged=true, count=10
// - 10 nodes, 8 have "value", 2 have different → converged=false, count=8
//
// func (cd *ConvergenceDetector) IsConverged(key string) (bool, int) {
//     cd.mu.RLock()
//     defer cd.mu.RUnlock()
//
//     if len(cd.nodes) == 0 {
//         return true, 0
//     }
//
//     // Get reference value from first node
//     firstState := cd.nodes[0].GetState()
//     referenceValue, hasKey := firstState[key]
//     if !hasKey {
//         // Check if any node has the key
//         for _, node := range cd.nodes {
//             state := node.GetState()
//             if _, exists := state[key]; exists {
//                 return false, 0
//             }
//         }
//         return true, 0
//     }
//
//     // Count nodes with matching value
//     convergedCount := 0
//     for _, node := range cd.nodes {
//         state := node.GetState()
//         if val, exists := state[key]; exists && val == referenceValue {
//             convergedCount++
//         }
//     }
//
//     return convergedCount == len(cd.nodes), convergedCount
// }
func (cd *ConvergenceDetector) IsConverged(key string) (bool, int) {
	// TODO: Implement
	// - Get state from all nodes
	// - Check if all have the same value for 'key'
	// - Return (converged bool, count of nodes with the value)
	return false, 0
}

// WaitForConvergence waits until convergence or timeout.
// TODO: Implement convergence waiting with polling
//
// Algorithm:
// 1. Calculate deadline (now + timeout)
// 2. Create ticker (check every 50ms)
// 3. Loop:
//    a. Wait for tick
//    b. Check if converged → return true
//    c. Check if deadline passed → return false
//
// Why polling instead of events?
// - Simpler implementation (no callback complexity)
// - Acceptable for testing (50ms overhead is fine)
// - Real systems might use event-driven approach
//
// func (cd *ConvergenceDetector) WaitForConvergence(key string, timeout time.Duration) bool {
//     deadline := time.Now().Add(timeout)
//     ticker := time.NewTicker(50 * time.Millisecond)
//     defer ticker.Stop()
//
//     for time.Now().Before(deadline) {
//         <-ticker.C
//         if converged, _ := cd.IsConverged(key); converged {
//             return true
//         }
//     }
//
//     return false
// }
func (cd *ConvergenceDetector) WaitForConvergence(key string, timeout time.Duration) bool {
	// TODO: Implement
	// - Use a ticker to check every 50ms
	// - Return true if converged
	// - Return false if timeout reached
	return false
}

// ============================================================================
// Exercise 5: Implement a gossip simulator
// ============================================================================

// Simulator orchestrates multiple gossip nodes and simulates gossip propagation.
//
// TODO: Implement simulator
//
// type Simulator struct {
//     nodes    []GossipNode
//     network  Network
//     detector *ConvergenceDetector
//     mu       sync.RWMutex
// }

// NewSimulator creates a gossip network simulator.
// TODO: Implement simulator initialization with topology building
//
// Parameters:
// - nodeCount: Number of nodes to create
// - fanout: Number of peers each node gossips to
// - latency: Network latency simulation
// - dropRate: Packet loss rate (0.0 to 1.0)
//
// Topology building strategy:
// 1. Create ring topology first (ensures connectivity)
//    - Node i connects to node (i+1) % n
//    - Guarantees all nodes reachable (no partitions)
// 2. Add random edges for redundancy
//    - Each node connects to ~30-40% of other nodes
//    - Improves fault tolerance and convergence speed
//
// Network topology comparison:
// - Ring: O(N) hops to reach all nodes (slow)
// - Fully connected: O(1) hops but O(N²) edges (doesn't scale)
// - Random (small-world): O(log N) hops, O(N log N) edges (optimal)
//
// func NewSimulator(nodeCount, fanout int, latency time.Duration, dropRate float64) *Simulator {
//     network := NewMockNetwork(latency, dropRate)
//
//     // Create nodes
//     nodes := make([]GossipNode, nodeCount)
//     for i := 0; i < nodeCount; i++ {
//         nodeID := fmt.Sprintf("node-%d", i)
//         node := NewGossipNode(nodeID, network, fanout)
//         nodes[i] = node
//         network.RegisterNode(node)
//     }
//
//     // Build topology: ring + random edges
//     for i := 0; i < nodeCount; i++ {
//         next := (i + 1) % nodeCount
//         nodes[i].AddPeer(nodes[next].ID())
//         nodes[next].AddPeer(nodes[i].ID())
//     }
//
//     // Add random connections
//     peersPerNode := max(3, int(float64(nodeCount)*0.4))
//     for _, node := range nodes {
//         connectedPeers := make(map[string]bool)
//         for len(connectedPeers) < peersPerNode && len(connectedPeers) < nodeCount-1 {
//             peerIdx := rand.Intn(nodeCount)
//             peer := nodes[peerIdx]
//             if peer.ID() != node.ID() && !connectedPeers[peer.ID()] {
//                 node.AddPeer(peer.ID())
//                 peer.AddPeer(node.ID())
//                 connectedPeers[peer.ID()] = true
//             }
//         }
//     }
//
//     detector := NewConvergenceDetector(nodes)
//     return &Simulator{
//         nodes:    nodes,
//         network:  network,
//         detector: detector,
//     }
// }
func NewSimulator(nodeCount, fanout int, latency time.Duration, dropRate float64) *Simulator {
	// TODO: Implement
	// 1. Create mock network
	// 2. Create nodeCount gossip nodes
	// 3. Register nodes with network
	// 4. Build random topology (each node connects to ~30% of other nodes)
	// 5. Create convergence detector
	// 6. Return simulator
	return nil
}

// SimulationStats contains statistics about the simulation
type SimulationStats struct {
	NodeCount    int
	MessageCount int
	DroppedCount int
}

// ============================================================================
// After implementing all functions:
// - Run: go test -v ./...
// - Run: go test -race ./... (check for data races!)
// - Experiment: Vary fanout (2, 3, 5, 10) and measure convergence time
// - Experiment: Add packet loss (10%, 20%) and observe behavior
// - Compare: Ring topology vs random graph topology
//
// Expected convergence times (100 nodes, fanout=3, latency=10ms):
// - 0% packet loss: ~5-10 rounds = 50-100ms
// - 10% packet loss: ~8-15 rounds = 80-150ms
// - 20% packet loss: ~15-25 rounds = 150-250ms
//
// Gossip protocol properties:
// - Scalability: Works with 1000+ nodes
// - Fault tolerance: Tolerates ~20-30% node failures
// - Eventual consistency: All nodes converge given enough time
// - Network efficiency: O(N log N) messages total
// ============================================================================
