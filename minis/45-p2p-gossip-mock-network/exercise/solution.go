//go:build solution
// +build solution

package exercise

import (
	"fmt"
	"math/rand"
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
type GossipNode interface {
	ID() string
	AddPeer(peerID string)
	Broadcast(msgType string, payload map[string]interface{}) error
	ReceiveMessage(msg Message) bool
	GetState() map[string]interface{}
	Shutdown()
}

// Network simulates a network that can deliver messages between nodes
type Network interface {
	RegisterNode(node GossipNode)
	Send(from, to string, msg Message)
	SetLatency(latency time.Duration)
	SetDropRate(rate float64)
	GetMessageCount() int
}

// GossipProtocol defines the behavior of the gossip protocol
type GossipProtocol interface {
	Fanout() int
	SelectPeers(allPeers []string, excludeID string) []string
	ShouldForward(msg Message) bool
}

// gossipNode implements GossipNode
type gossipNode struct {
	id       string
	network  Network
	fanout   int
	peers    []string
	state    map[string]interface{}
	seen     map[string]bool
	mu       sync.RWMutex
	shutdown chan struct{}
}

func NewGossipNode(id string, network Network, fanout int) GossipNode {
	return &gossipNode{
		id:       id,
		network:  network,
		fanout:   fanout,
		peers:    make([]string, 0),
		state:    make(map[string]interface{}),
		seen:     make(map[string]bool),
		shutdown: make(chan struct{}),
	}
}

func (n *gossipNode) ID() string {
	return n.id
}

func (n *gossipNode) AddPeer(peerID string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Avoid duplicates and self-connections
	if peerID == n.id {
		return
	}

	for _, p := range n.peers {
		if p == peerID {
			return
		}
	}

	n.peers = append(n.peers, peerID)
}

func (n *gossipNode) Broadcast(msgType string, payload map[string]interface{}) error {
	msg := Message{
		ID:        fmt.Sprintf("%s-%d", n.id, time.Now().UnixNano()),
		Type:      msgType,
		From:      n.id,
		Payload:   payload,
		Timestamp: time.Now(),
	}

	// Process locally first
	n.ReceiveMessage(msg)

	return nil
}

func (n *gossipNode) ReceiveMessage(msg Message) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Deduplicate
	if n.seen[msg.ID] {
		return false
	}
	n.seen[msg.ID] = true

	// Apply update to local state
	for key, value := range msg.Payload {
		n.state[key] = value
	}

	// Forward to random peers asynchronously
	go n.forwardToPeers(msg)

	return true
}

func (n *gossipNode) forwardToPeers(msg Message) {
	n.mu.RLock()
	peers := n.selectRandomPeers(n.fanout, msg.From)
	n.mu.RUnlock()

	for _, peerID := range peers {
		n.network.Send(n.id, peerID, msg)
	}
}

func (n *gossipNode) selectRandomPeers(count int, excludeID string) []string {
	// Filter out the excluded ID
	eligible := make([]string, 0)
	for _, p := range n.peers {
		if p != excludeID {
			eligible = append(eligible, p)
		}
	}

	if count >= len(eligible) {
		return eligible
	}

	// Fisher-Yates shuffle and take first 'count'
	shuffled := make([]string, len(eligible))
	copy(shuffled, eligible)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled[:count]
}

func (n *gossipNode) GetState() map[string]interface{} {
	n.mu.RLock()
	defer n.mu.RUnlock()

	stateCopy := make(map[string]interface{})
	for k, v := range n.state {
		stateCopy[k] = v
	}
	return stateCopy
}

func (n *gossipNode) Shutdown() {
	close(n.shutdown)
}

// mockNetwork implements Network
type mockNetwork struct {
	nodes        map[string]GossipNode
	latency      time.Duration
	dropRate     float64
	messageCount int
	droppedCount int
	mu           sync.RWMutex
}

func NewMockNetwork(latency time.Duration, dropRate float64) Network {
	return &mockNetwork{
		nodes:    make(map[string]GossipNode),
		latency:  latency,
		dropRate: dropRate,
	}
}

func (mn *mockNetwork) RegisterNode(node GossipNode) {
	mn.mu.Lock()
	defer mn.mu.Unlock()
	mn.nodes[node.ID()] = node
}

func (mn *mockNetwork) Send(from, to string, msg Message) {
	mn.mu.Lock()
	mn.messageCount++

	// Simulate packet loss
	if rand.Float64() < mn.dropRate {
		mn.droppedCount++
		mn.mu.Unlock()
		return
	}

	latency := mn.latency
	mn.mu.Unlock()

	// Simulate network latency
	time.AfterFunc(latency, func() {
		mn.mu.RLock()
		node, exists := mn.nodes[to]
		mn.mu.RUnlock()

		if exists {
			node.ReceiveMessage(msg)
		}
	})
}

func (mn *mockNetwork) SetLatency(latency time.Duration) {
	mn.mu.Lock()
	defer mn.mu.Unlock()
	mn.latency = latency
}

func (mn *mockNetwork) SetDropRate(rate float64) {
	mn.mu.Lock()
	defer mn.mu.Unlock()
	mn.dropRate = rate
}

func (mn *mockNetwork) GetMessageCount() int {
	mn.mu.RLock()
	defer mn.mu.RUnlock()
	return mn.messageCount
}

func (mn *mockNetwork) GetDroppedCount() int {
	mn.mu.RLock()
	defer mn.mu.RUnlock()
	return mn.droppedCount
}

// PushGossipProtocol implements GossipProtocol
type PushGossipProtocol struct {
	fanout int
}

func NewPushGossipProtocol(fanout int) GossipProtocol {
	return &PushGossipProtocol{
		fanout: fanout,
	}
}

func (p *PushGossipProtocol) Fanout() int {
	return p.fanout
}

func (p *PushGossipProtocol) SelectPeers(allPeers []string, excludeID string) []string {
	// Filter out excluded ID
	eligible := make([]string, 0)
	for _, peer := range allPeers {
		if peer != excludeID {
			eligible = append(eligible, peer)
		}
	}

	if p.fanout >= len(eligible) {
		return eligible
	}

	// Shuffle and select
	shuffled := make([]string, len(eligible))
	copy(shuffled, eligible)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled[:p.fanout]
}

func (p *PushGossipProtocol) ShouldForward(msg Message) bool {
	// For basic push protocol, always forward
	return true
}

// ConvergenceDetector checks if nodes have converged
type ConvergenceDetector struct {
	nodes []GossipNode
	mu    sync.RWMutex
}

func NewConvergenceDetector(nodes []GossipNode) *ConvergenceDetector {
	return &ConvergenceDetector{
		nodes: nodes,
	}
}

func (cd *ConvergenceDetector) IsConverged(key string) (bool, int) {
	cd.mu.RLock()
	defer cd.mu.RUnlock()

	if len(cd.nodes) == 0 {
		return true, 0
	}

	// Get reference value from first node
	firstState := cd.nodes[0].GetState()
	referenceValue, hasKey := firstState[key]

	if !hasKey {
		// Check if any node has the key
		for _, node := range cd.nodes {
			state := node.GetState()
			if _, exists := state[key]; exists {
				return false, 0
			}
		}
		return true, 0
	}

	// Count nodes with matching value
	convergedCount := 0
	for _, node := range cd.nodes {
		state := node.GetState()
		if val, exists := state[key]; exists && val == referenceValue {
			convergedCount++
		}
	}

	return convergedCount == len(cd.nodes), convergedCount
}

func (cd *ConvergenceDetector) WaitForConvergence(key string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for time.Now().Before(deadline) {
		<-ticker.C
		if converged, _ := cd.IsConverged(key); converged {
			return true
		}
	}

	return false
}

// Simulator orchestrates the gossip simulation
type Simulator struct {
	nodes    []GossipNode
	network  Network
	detector *ConvergenceDetector
	mu       sync.RWMutex
}

func NewSimulator(nodeCount, fanout int, latency time.Duration, dropRate float64) *Simulator {
	// Create network
	network := NewMockNetwork(latency, dropRate)

	// Create nodes
	nodes := make([]GossipNode, nodeCount)
	for i := 0; i < nodeCount; i++ {
		nodeID := fmt.Sprintf("node-%d", i)
		node := NewGossipNode(nodeID, network, fanout)
		nodes[i] = node
		network.RegisterNode(node)
	}

	// Build topology: ring + random edges for guaranteed connectivity
	// First, create a ring to ensure all nodes are reachable
	for i := 0; i < nodeCount; i++ {
		next := (i + 1) % nodeCount
		nodes[i].AddPeer(nodes[next].ID())
		nodes[next].AddPeer(nodes[i].ID())
	}

	// Then add random connections for redundancy
	peersPerNode := max(3, int(float64(nodeCount)*0.4))
	for _, node := range nodes {
		connectedPeers := make(map[string]bool)

		for len(connectedPeers) < peersPerNode && len(connectedPeers) < nodeCount-1 {
			peerIdx := rand.Intn(nodeCount)
			peer := nodes[peerIdx]

			if peer.ID() != node.ID() && !connectedPeers[peer.ID()] {
				node.AddPeer(peer.ID())
				peer.AddPeer(node.ID()) // Bidirectional
				connectedPeers[peer.ID()] = true
			}
		}
	}

	// Create convergence detector
	detector := NewConvergenceDetector(nodes)

	return &Simulator{
		nodes:    nodes,
		network:  network,
		detector: detector,
	}
}

func (s *Simulator) BroadcastFrom(nodeID string, msgType string, payload map[string]interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, node := range s.nodes {
		if node.ID() == nodeID {
			return node.Broadcast(msgType, payload)
		}
	}

	return fmt.Errorf("node %s not found", nodeID)
}

func (s *Simulator) WaitForConvergence(key string, timeout time.Duration) bool {
	return s.detector.WaitForConvergence(key, timeout)
}

func (s *Simulator) GetStats() SimulationStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	messageCount := s.network.GetMessageCount()
	droppedCount := 0

	if mn, ok := s.network.(*mockNetwork); ok {
		droppedCount = mn.GetDroppedCount()
	}

	return SimulationStats{
		NodeCount:    len(s.nodes),
		MessageCount: messageCount,
		DroppedCount: droppedCount,
	}
}

func (s *Simulator) Shutdown() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, node := range s.nodes {
		node.Shutdown()
	}
}

// SimulationStats contains statistics about the simulation
type SimulationStats struct {
	NodeCount    int
	MessageCount int
	DroppedCount int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
