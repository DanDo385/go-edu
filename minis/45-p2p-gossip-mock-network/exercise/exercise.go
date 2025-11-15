//go:build !solution
// +build !solution

package exercise

import (
	"sync"
	"time"
)

// Message represents a gossip message in the network
type Message struct {
	ID        string
	Type      string
	From      string
	Payload   map[string]interface{}
	Timestamp time.Time
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

// Exercise 1: Implement a basic gossip node
//
// Create a gossipNode that:
// - Maintains a list of peer IDs
// - Tracks seen messages to avoid duplicates
// - Forwards new messages to random peers
// - Stores received data in local state
//
// Hints:
// - Use a map to track seen message IDs
// - Use a mutex to protect concurrent access
// - Implement random peer selection (fanout = 3)
// - Don't forward messages back to the sender
func NewGossipNode(id string, network Network, fanout int) GossipNode {
	// TODO: Implement
	// - Create node struct with id, network, fanout, peers list, state map, seen map
	// - Initialize mutex for thread safety
	// - Return the node
	return nil
}

// Exercise 2: Implement a mock network
//
// Create a mockNetwork that:
// - Maintains a registry of nodes
// - Simulates message delivery with configurable latency
// - Simulates packet loss based on drop rate
// - Tracks message count for statistics
//
// Hints:
// - Use time.AfterFunc for simulated latency
// - Use rand.Float64() < dropRate to simulate packet loss
// - Use mutex to protect the nodes map
// - Keep a counter of all messages (successful and dropped)
func NewMockNetwork(latency time.Duration, dropRate float64) Network {
	// TODO: Implement
	// - Create network struct with nodes map, latency, dropRate, message counter
	// - Initialize mutex for thread safety
	// - Return the network
	return nil
}

// Exercise 3: Implement a push-based gossip protocol
//
// Create a pushGossipProtocol that:
// - Uses a fixed fanout (number of peers to forward to)
// - Randomly selects peers for message forwarding
// - Continues forwarding messages (no TTL limit for now)
//
// Hints:
// - Use Fisher-Yates shuffle for random peer selection
// - Exclude the sender from peer selection
// - Return up to 'fanout' peers
type PushGossipProtocol struct {
	// TODO: Add fields
	// - fanout: number of peers to gossip to
}

func NewPushGossipProtocol(fanout int) GossipProtocol {
	// TODO: Implement
	return nil
}

func (p *PushGossipProtocol) Fanout() int {
	// TODO: Implement
	return 0
}

func (p *PushGossipProtocol) SelectPeers(allPeers []string, excludeID string) []string {
	// TODO: Implement
	// - Filter out excludeID from allPeers
	// - Shuffle the remaining peers
	// - Return up to Fanout() peers
	return nil
}

func (p *PushGossipProtocol) ShouldForward(msg Message) bool {
	// TODO: Implement
	// - For basic push protocol, always forward
	// - Can add TTL or age-based filtering later
	return false
}

// Exercise 4: Implement convergence detection
//
// Create a ConvergenceDetector that:
// - Checks if all nodes in a network have the same state
// - Counts how many nodes have converged
// - Can wait for convergence with a timeout
//
// Hints:
// - Compare state maps from all nodes
// - Use a ticker to periodically check convergence
// - Return early if converged before timeout
type ConvergenceDetector struct {
	nodes []GossipNode
	mu    sync.RWMutex
}

func NewConvergenceDetector(nodes []GossipNode) *ConvergenceDetector {
	// TODO: Implement
	return nil
}

// IsConverged checks if all nodes have the same state for a given key
func (cd *ConvergenceDetector) IsConverged(key string) (bool, int) {
	// TODO: Implement
	// - Get state from all nodes
	// - Check if all have the same value for 'key'
	// - Return (converged bool, count of nodes with the value)
	return false, 0
}

// WaitForConvergence waits until convergence or timeout
func (cd *ConvergenceDetector) WaitForConvergence(key string, timeout time.Duration) bool {
	// TODO: Implement
	// - Use a ticker to check every 50ms
	// - Return true if converged
	// - Return false if timeout reached
	return false
}

// Exercise 5: Implement a gossip simulator
//
// Create a Simulator that:
// - Creates multiple gossip nodes
// - Connects them in a random network topology
// - Can initiate broadcasts from any node
// - Tracks convergence and statistics
//
// Hints:
// - Each node should connect to ~log(N) peers for good connectivity
// - Use bidirectional connections (if A→B, then B→A)
// - Avoid self-connections
type Simulator struct {
	nodes            []GossipNode
	network          Network
	detector         *ConvergenceDetector
	mu               sync.RWMutex
}

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

// BroadcastFrom initiates a broadcast from a specific node
func (s *Simulator) BroadcastFrom(nodeID string, msgType string, payload map[string]interface{}) error {
	// TODO: Implement
	// - Find node with matching ID
	// - Call node.Broadcast()
	// - Return error if node not found
	return nil
}

// WaitForConvergence waits for all nodes to converge on a key
func (s *Simulator) WaitForConvergence(key string, timeout time.Duration) bool {
	// TODO: Implement
	// - Use convergence detector
	return false
}

// GetStats returns simulation statistics
func (s *Simulator) GetStats() SimulationStats {
	// TODO: Implement
	// - Collect message count from network
	// - Count nodes
	// - Return stats
	return SimulationStats{}
}

// Shutdown gracefully shuts down all nodes
func (s *Simulator) Shutdown() {
	// TODO: Implement
	// - Call Shutdown() on all nodes
}

// SimulationStats contains statistics about the simulation
type SimulationStats struct {
	NodeCount     int
	MessageCount  int
	DroppedCount  int
}
