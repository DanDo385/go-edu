package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// GossipMessage represents a message being gossiped through the network
type GossipMessage struct {
	ID        string
	Payload   map[string]interface{}
	From      string
	Timestamp time.Time
	TTL       int
}

// SimNode represents a node in the simulated gossip network
type SimNode struct {
	id        string
	peers     []string
	state     map[string]interface{}
	received  map[string]bool
	fanout    int
	network   *MockNetwork
	mu        sync.RWMutex
	stats     NodeStats
}

// NodeStats tracks statistics for a node
type NodeStats struct {
	MessagesReceived int
	MessagesSent     int
	DuplicatesIgnored int
}

// NewSimNode creates a new simulated node
func NewSimNode(id string, fanout int, network *MockNetwork) *SimNode {
	return &SimNode{
		id:       id,
		peers:    make([]string, 0),
		state:    make(map[string]interface{}),
		received: make(map[string]bool),
		fanout:   fanout,
		network:  network,
	}
}

// AddPeer adds a peer to this node's peer list
func (n *SimNode) AddPeer(peerID string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Avoid duplicates
	for _, p := range n.peers {
		if p == peerID {
			return
		}
	}
	n.peers = append(n.peers, peerID)
}

// ReceiveMessage handles incoming gossip messages
func (n *SimNode) ReceiveMessage(msg GossipMessage) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.stats.MessagesReceived++

	// Deduplicate - ignore if already seen
	if n.received[msg.ID] {
		n.stats.DuplicatesIgnored++
		return
	}
	n.received[msg.ID] = true

	// Apply update to local state
	for key, value := range msg.Payload {
		n.state[key] = value
	}

	// Check TTL before forwarding
	if msg.TTL <= 0 {
		return
	}

	// Gossip to peers (async)
	go n.gossipToPeers(msg)
}

// gossipToPeers forwards a message to random peers
func (n *SimNode) gossipToPeers(msg GossipMessage) {
	n.mu.RLock()
	selectedPeers := n.selectRandomPeers(n.fanout)
	from := msg.From
	n.mu.RUnlock()

	// Decrement TTL
	msg.TTL--

	for _, peerID := range selectedPeers {
		// Don't send back to sender
		if peerID == from {
			continue
		}

		n.mu.Lock()
		n.stats.MessagesSent++
		n.mu.Unlock()

		n.network.Send(n.id, peerID, msg)
	}
}

// selectRandomPeers picks random peers up to count
func (n *SimNode) selectRandomPeers(count int) []string {
	if count >= len(n.peers) {
		return append([]string{}, n.peers...)
	}

	// Fisher-Yates shuffle and take first 'count'
	shuffled := make([]string, len(n.peers))
	copy(shuffled, n.peers)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled[:count]
}

// Broadcast initiates a new message from this node
func (n *SimNode) Broadcast(payload map[string]interface{}, ttl int) GossipMessage {
	msg := GossipMessage{
		ID:        fmt.Sprintf("%s-%d", n.id, time.Now().UnixNano()),
		Payload:   payload,
		From:      n.id,
		Timestamp: time.Now(),
		TTL:       ttl,
	}

	// Process locally first
	n.ReceiveMessage(msg)

	return msg
}

// GetState returns a copy of the node's current state
func (n *SimNode) GetState() map[string]interface{} {
	n.mu.RLock()
	defer n.mu.RUnlock()

	stateCopy := make(map[string]interface{})
	for k, v := range n.state {
		stateCopy[k] = v
	}
	return stateCopy
}

// GetStats returns a copy of node statistics
func (n *SimNode) GetStats() NodeStats {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.stats
}

// MockNetwork simulates a network with latency and packet loss
type MockNetwork struct {
	nodes       map[string]*SimNode
	latency     time.Duration
	dropRate    float64
	messageLog  []NetworkEvent
	mu          sync.RWMutex
}

// NetworkEvent represents a network event for logging
type NetworkEvent struct {
	Timestamp time.Time
	From      string
	To        string
	MessageID string
	Dropped   bool
}

// NewMockNetwork creates a new simulated network
func NewMockNetwork(latency time.Duration, dropRate float64) *MockNetwork {
	return &MockNetwork{
		nodes:      make(map[string]*SimNode),
		latency:    latency,
		dropRate:   dropRate,
		messageLog: make([]NetworkEvent, 0),
	}
}

// RegisterNode adds a node to the network
func (mn *MockNetwork) RegisterNode(node *SimNode) {
	mn.mu.Lock()
	defer mn.mu.Unlock()
	mn.nodes[node.id] = node
}

// Send simulates sending a message with latency and potential packet loss
func (mn *MockNetwork) Send(from, to string, msg GossipMessage) {
	mn.mu.Lock()

	// Log event
	event := NetworkEvent{
		Timestamp: time.Now(),
		From:      from,
		To:        to,
		MessageID: msg.ID,
		Dropped:   false,
	}

	// Simulate packet loss
	if rand.Float64() < mn.dropRate {
		event.Dropped = true
		mn.messageLog = append(mn.messageLog, event)
		mn.mu.Unlock()
		return
	}

	mn.messageLog = append(mn.messageLog, event)
	mn.mu.Unlock()

	// Simulate network latency
	time.AfterFunc(mn.latency, func() {
		mn.mu.RLock()
		node, exists := mn.nodes[to]
		mn.mu.RUnlock()

		if exists {
			node.ReceiveMessage(msg)
		}
	})
}

// GetStats returns network statistics
func (mn *MockNetwork) GetStats() (total, dropped int) {
	mn.mu.RLock()
	defer mn.mu.RUnlock()

	total = len(mn.messageLog)
	for _, event := range mn.messageLog {
		if event.Dropped {
			dropped++
		}
	}
	return
}

// Simulator orchestrates the gossip simulation
type Simulator struct {
	nodes   []*SimNode
	network *MockNetwork
	mu      sync.RWMutex
}

// NewSimulator creates a new gossip network simulator
func NewSimulator(nodeCount, fanout int, latency time.Duration, dropRate float64) *Simulator {
	network := NewMockNetwork(latency, dropRate)
	nodes := make([]*SimNode, nodeCount)

	// Create nodes
	for i := 0; i < nodeCount; i++ {
		nodeID := fmt.Sprintf("node-%d", i)
		node := NewSimNode(nodeID, fanout, network)
		nodes[i] = node
		network.RegisterNode(node)
	}

	// Build peer relationships
	// First, create a ring to ensure connectivity
	for i := 0; i < nodeCount; i++ {
		next := (i + 1) % nodeCount
		nodes[i].AddPeer(nodes[next].id)
		nodes[next].AddPeer(nodes[i].id)
	}

	// Then add random connections for redundancy
	peersPerNode := max(3, int(float64(nodeCount)*0.4))
	for _, node := range nodes {
		connectedPeers := make(map[string]bool)
		for len(connectedPeers) < peersPerNode {
			peerIdx := rand.Intn(nodeCount)
			peer := nodes[peerIdx]
			if peer.id != node.id && !connectedPeers[peer.id] {
				node.AddPeer(peer.id)
				peer.AddPeer(node.id)
				connectedPeers[peer.id] = true
			}
		}
	}

	return &Simulator{
		nodes:   nodes,
		network: network,
	}
}

// BroadcastFrom initiates a broadcast from a specific node
func (s *Simulator) BroadcastFrom(nodeID string, payload map[string]interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, node := range s.nodes {
		if node.id == nodeID {
			// TTL = network diameter estimate
			ttl := len(s.nodes)
			node.Broadcast(payload, ttl)
			return nil
		}
	}
	return fmt.Errorf("node %s not found", nodeID)
}

// CheckConvergence checks if all nodes have the same state
func (s *Simulator) CheckConvergence(key string) (bool, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.nodes) == 0 {
		return true, 0
	}

	// Get reference value from first node
	firstNode := s.nodes[0]
	referenceState := firstNode.GetState()
	referenceValue, hasKey := referenceState[key]

	if !hasKey {
		// Check if any node has the key
		for _, node := range s.nodes {
			state := node.GetState()
			if _, exists := state[key]; exists {
				return false, 0
			}
		}
		return true, 0 // No node has it yet
	}

	// Count nodes that have the correct value
	convergedCount := 0
	for _, node := range s.nodes {
		state := node.GetState()
		if val, exists := state[key]; exists && val == referenceValue {
			convergedCount++
		}
	}

	return convergedCount == len(s.nodes), convergedCount
}

// WaitForConvergence waits until all nodes converge on a key
func (s *Simulator) WaitForConvergence(key string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for time.Now().Before(deadline) {
		<-ticker.C
		if converged, _ := s.CheckConvergence(key); converged {
			return true
		}
	}
	return false
}

// PrintStats prints simulation statistics
func (s *Simulator) PrintStats() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Println("\n=== Simulation Statistics ===")

	totalMsgReceived := 0
	totalMsgSent := 0
	totalDuplicates := 0

	for _, node := range s.nodes {
		stats := node.GetStats()
		totalMsgReceived += stats.MessagesReceived
		totalMsgSent += stats.MessagesSent
		totalDuplicates += stats.DuplicatesIgnored
	}

	totalNetMsg, droppedMsg := s.network.GetStats()

	fmt.Printf("Total Nodes: %d\n", len(s.nodes))
	fmt.Printf("Messages Sent: %d\n", totalMsgSent)
	fmt.Printf("Messages Received: %d\n", totalMsgReceived)
	fmt.Printf("Duplicates Ignored: %d\n", totalDuplicates)
	fmt.Printf("Network Messages: %d\n", totalNetMsg)
	fmt.Printf("Dropped Messages: %d\n", droppedMsg)
	fmt.Printf("Efficiency: %.2f%% (received/sent)\n",
		float64(totalMsgReceived)/float64(max(1, totalMsgSent))*100)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== P2P Gossip Network Simulation ===\n")

	// Configuration
	nodeCount := 15
	fanout := 3
	latency := 10 * time.Millisecond
	dropRate := 0.0 // 0% packet loss for demo (try 0.05 for 5% loss)

	fmt.Printf("Configuration:\n")
	fmt.Printf("  Nodes: %d\n", nodeCount)
	fmt.Printf("  Fanout: %d\n", fanout)
	fmt.Printf("  Network Latency: %v\n", latency)
	fmt.Printf("  Packet Drop Rate: %.1f%%\n\n", dropRate*100)

	// Create simulator
	sim := NewSimulator(nodeCount, fanout, latency, dropRate)

	// Demonstration 1: Simple broadcast
	fmt.Println("--- Demo 1: Simple Broadcast ---")
	payload1 := map[string]interface{}{
		"message": "Hello from node-0",
		"version": 1,
	}

	start := time.Now()
	sim.BroadcastFrom("node-0", payload1)

	// Give initial time for async forwarding to start
	time.Sleep(50 * time.Millisecond)

	// Monitor convergence
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	round := 0
	for {
		<-ticker.C
		converged, count := sim.CheckConvergence("message")
		round++

		fmt.Printf("Round %d: %d/%d nodes have the message\n", round, count, nodeCount)

		if converged {
			fmt.Printf("✓ Converged in %v\n", time.Since(start))
			break
		}

		if round > 30 {
			fmt.Printf("Partial convergence: %d/%d nodes (%.1f%%)\n",
				count, nodeCount, float64(count)/float64(nodeCount)*100)
			break
		}
	}

	// Wait for messages to settle
	time.Sleep(300 * time.Millisecond)

	// Demonstration 2: Multiple concurrent broadcasts
	fmt.Println("\n--- Demo 2: Concurrent Broadcasts ---")

	payload2 := map[string]interface{}{"source": "node-5", "data": "A"}
	payload3 := map[string]interface{}{"source": "node-10", "data": "B"}
	payload4 := map[string]interface{}{"source": "node-15", "data": "C"}

	start = time.Now()
	sim.BroadcastFrom("node-5", payload2)
	sim.BroadcastFrom("node-10", payload3)
	sim.BroadcastFrom("node-15", payload4)

	// Wait for all to converge
	converged1 := sim.WaitForConvergence("source", 5*time.Second)

	if converged1 {
		fmt.Printf("✓ All broadcasts converged in %v\n", time.Since(start))
	} else {
		fmt.Println("✗ Not all broadcasts converged")
	}

	// Demonstration 3: Show state consistency
	fmt.Println("\n--- Demo 3: State Consistency Check ---")

	allStates := make(map[string]map[string]interface{})
	for _, node := range sim.nodes {
		state := node.GetState()
		allStates[node.id] = state
	}

	// Check if all nodes have the same state
	firstState := allStates["node-0"]
	allSame := true

	for nodeID, state := range allStates {
		if len(state) != len(firstState) {
			allSame = false
			fmt.Printf("✗ Node %s has different state size\n", nodeID)
			break
		}

		for k, v := range firstState {
			if state[k] != v {
				allSame = false
				fmt.Printf("✗ Node %s has different value for key %s\n", nodeID, k)
				break
			}
		}

		if !allSame {
			break
		}
	}

	if allSame {
		fmt.Println("✓ All nodes have consistent state")
		fmt.Printf("  State keys: %v\n", getKeys(firstState))
	}

	// Print final statistics
	sim.PrintStats()

	// Demonstration 4: Visualize node connectivity
	fmt.Println("\n--- Demo 4: Network Topology ---")
	for i := 0; i < min(5, len(sim.nodes)); i++ {
		node := sim.nodes[i]
		node.mu.RLock()
		peerCount := len(node.peers)
		node.mu.RUnlock()
		fmt.Printf("Node %s: %d peers\n", node.id, peerCount)
	}

	fmt.Println("\n=== Simulation Complete ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
