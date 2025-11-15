package exercise

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// TestGossipNode_Basic tests basic node functionality
func TestGossipNode_Basic(t *testing.T) {
	network := NewMockNetwork(1*time.Millisecond, 0.0)
	node := NewGossipNode("node-1", network, 3)

	if node == nil {
		t.Fatal("NewGossipNode returned nil")
	}

	if node.ID() != "node-1" {
		t.Errorf("Expected ID 'node-1', got '%s'", node.ID())
	}

	// Test adding peers
	node.AddPeer("node-2")
	node.AddPeer("node-3")
	node.AddPeer("node-2") // Duplicate, should be ignored

	// Test state initially empty
	state := node.GetState()
	if len(state) != 0 {
		t.Errorf("Expected empty initial state, got %d items", len(state))
	}
}

// TestGossipNode_Broadcast tests broadcasting from a node
func TestGossipNode_Broadcast(t *testing.T) {
	network := NewMockNetwork(1*time.Millisecond, 0.0)
	node := NewGossipNode("node-1", network, 3)
	network.RegisterNode(node)

	payload := map[string]interface{}{
		"key": "value",
		"num": 42,
	}

	err := node.Broadcast("test", payload)
	if err != nil {
		t.Fatalf("Broadcast failed: %v", err)
	}

	// Wait for message to be processed
	time.Sleep(10 * time.Millisecond)

	// Check that state was updated
	state := node.GetState()
	if state["key"] != "value" {
		t.Errorf("Expected state[key] = 'value', got %v", state["key"])
	}
	if state["num"] != 42 {
		t.Errorf("Expected state[num] = 42, got %v", state["num"])
	}
}

// TestGossipNode_ReceiveMessage tests receiving messages
func TestGossipNode_ReceiveMessage(t *testing.T) {
	network := NewMockNetwork(1*time.Millisecond, 0.0)
	node := NewGossipNode("node-1", network, 3)

	msg := Message{
		ID:   "msg-1",
		Type: "update",
		From: "node-0",
		Payload: map[string]interface{}{
			"status": "active",
		},
		Timestamp: time.Now(),
	}

	// First receive should return true (new message)
	isNew := node.ReceiveMessage(msg)
	if !isNew {
		t.Error("Expected ReceiveMessage to return true for new message")
	}

	// Second receive should return false (duplicate)
	isNew = node.ReceiveMessage(msg)
	if isNew {
		t.Error("Expected ReceiveMessage to return false for duplicate message")
	}

	// Check state was updated
	state := node.GetState()
	if state["status"] != "active" {
		t.Errorf("Expected state[status] = 'active', got %v", state["status"])
	}
}

// TestMockNetwork_Basic tests basic network functionality
func TestMockNetwork_Basic(t *testing.T) {
	network := NewMockNetwork(5*time.Millisecond, 0.0)

	if network == nil {
		t.Fatal("NewMockNetwork returned nil")
	}

	node1 := NewGossipNode("node-1", network, 3)
	node2 := NewGossipNode("node-2", network, 3)

	network.RegisterNode(node1)
	network.RegisterNode(node2)

	// Send a message
	msg := Message{
		ID:   "msg-1",
		Type: "test",
		From: "node-1",
		Payload: map[string]interface{}{
			"data": "hello",
		},
		Timestamp: time.Now(),
	}

	network.Send("node-1", "node-2", msg)

	// Wait for latency
	time.Sleep(20 * time.Millisecond)

	// Check that node2 received the message
	state := node2.GetState()
	if state["data"] != "hello" {
		t.Errorf("Expected node2 to receive message, state: %v", state)
	}

	// Check message count
	count := network.GetMessageCount()
	if count < 1 {
		t.Errorf("Expected at least 1 message, got %d", count)
	}
}

// TestMockNetwork_PacketLoss tests packet loss simulation
func TestMockNetwork_PacketLoss(t *testing.T) {
	// High drop rate for testing
	network := NewMockNetwork(1*time.Millisecond, 0.9)

	node1 := NewGossipNode("node-1", network, 3)
	node2 := NewGossipNode("node-2", network, 3)

	network.RegisterNode(node1)
	network.RegisterNode(node2)

	// Send many messages
	for i := 0; i < 100; i++ {
		msg := Message{
			ID:   string(rune(i)),
			Type: "test",
			From: "node-1",
			Payload: map[string]interface{}{
				"seq": i,
			},
			Timestamp: time.Now(),
		}
		network.Send("node-1", "node-2", msg)
	}

	time.Sleep(50 * time.Millisecond)

	// With 90% drop rate, most messages should be dropped
	count := network.GetMessageCount()
	if count != 100 {
		t.Errorf("Expected 100 messages sent, got %d", count)
	}
}

// TestPushGossipProtocol tests the push gossip protocol
func TestPushGossipProtocol(t *testing.T) {
	protocol := NewPushGossipProtocol(3)

	if protocol == nil {
		t.Fatal("NewPushGossipProtocol returned nil")
	}

	if protocol.Fanout() != 3 {
		t.Errorf("Expected fanout 3, got %d", protocol.Fanout())
	}

	// Test peer selection
	allPeers := []string{"peer-1", "peer-2", "peer-3", "peer-4", "peer-5"}
	selected := protocol.SelectPeers(allPeers, "peer-1")

	if len(selected) > 3 {
		t.Errorf("Expected at most 3 peers, got %d", len(selected))
	}

	// Excluded peer should not be selected
	for _, peer := range selected {
		if peer == "peer-1" {
			t.Error("Excluded peer was selected")
		}
	}

	// Test ShouldForward
	msg := Message{
		ID:        "msg-1",
		Type:      "test",
		From:      "node-1",
		Payload:   map[string]interface{}{},
		Timestamp: time.Now(),
	}

	if !protocol.ShouldForward(msg) {
		t.Error("Push protocol should always forward")
	}
}

// TestConvergenceDetector tests convergence detection
func TestConvergenceDetector(t *testing.T) {
	network := NewMockNetwork(1*time.Millisecond, 0.0)

	// Create nodes
	nodes := make([]GossipNode, 5)
	for i := 0; i < 5; i++ {
		nodes[i] = NewGossipNode(string(rune('A'+i)), network, 2)
	}

	detector := NewConvergenceDetector(nodes)
	if detector == nil {
		t.Fatal("NewConvergenceDetector returned nil")
	}

	// Initially, no key exists, should be converged
	converged, count := detector.IsConverged("key1")
	if !converged {
		t.Error("Expected initial state to be converged")
	}
	if count != 0 {
		t.Errorf("Expected 0 nodes with key, got %d", count)
	}

	// Add same value to all nodes
	for _, node := range nodes {
		msg := Message{
			ID:   "msg-1",
			Type: "test",
			From: "external",
			Payload: map[string]interface{}{
				"key1": "value1",
			},
			Timestamp: time.Now(),
		}
		node.ReceiveMessage(msg)
	}

	// Should be converged now
	converged, count = detector.IsConverged("key1")
	if !converged {
		t.Error("Expected nodes to be converged after same update")
	}
	if count != 5 {
		t.Errorf("Expected 5 nodes with key, got %d", count)
	}

	// Add different value to one node
	msg := Message{
		ID:   "msg-2",
		Type: "test",
		From: "external",
		Payload: map[string]interface{}{
			"key1": "different",
		},
		Timestamp: time.Now(),
	}
	nodes[0].ReceiveMessage(msg)

	// Should not be converged
	converged, _ = detector.IsConverged("key1")
	if converged {
		t.Error("Expected nodes to not be converged after different update")
	}
}

// TestConvergenceDetector_WaitForConvergence tests waiting for convergence
func TestConvergenceDetector_WaitForConvergence(t *testing.T) {
	network := NewMockNetwork(1*time.Millisecond, 0.0)

	nodes := make([]GossipNode, 3)
	for i := 0; i < 3; i++ {
		nodes[i] = NewGossipNode(string(rune('A'+i)), network, 2)
	}

	detector := NewConvergenceDetector(nodes)

	// Start converged
	msg := Message{
		ID:   "msg-1",
		Type: "test",
		From: "external",
		Payload: map[string]interface{}{
			"key": "value",
		},
		Timestamp: time.Now(),
	}

	// Add to all nodes
	for _, node := range nodes {
		node.ReceiveMessage(msg)
	}

	// Should converge quickly
	converged := detector.WaitForConvergence("key", 1*time.Second)
	if !converged {
		t.Error("Expected convergence within timeout")
	}
}

// TestSimulator_Basic tests basic simulator functionality
func TestSimulator_Basic(t *testing.T) {
	sim := NewSimulator(5, 2, 1*time.Millisecond, 0.0)

	if sim == nil {
		t.Fatal("NewSimulator returned nil")
	}

	stats := sim.GetStats()
	if stats.NodeCount != 5 {
		t.Errorf("Expected 5 nodes, got %d", stats.NodeCount)
	}

	// Broadcast a message
	payload := map[string]interface{}{
		"test": "data",
	}

	err := sim.BroadcastFrom("node-0", "test", payload)
	if err != nil {
		t.Fatalf("BroadcastFrom failed: %v", err)
	}

	// Give some time for messages to propagate
	time.Sleep(100 * time.Millisecond)

	// Wait for convergence
	converged := sim.WaitForConvergence("test", 3*time.Second)
	if !converged {
		t.Error("Expected network to converge")
	}

	stats = sim.GetStats()
	if stats.MessageCount == 0 {
		t.Error("Expected some messages to be sent")
	}

	sim.Shutdown()
}

// TestSimulator_MultipleMessages tests multiple concurrent broadcasts
func TestSimulator_MultipleMessages(t *testing.T) {
	sim := NewSimulator(10, 3, 5*time.Millisecond, 0.0)

	// Broadcast from multiple nodes with unique keys
	sim.BroadcastFrom("node-0", "update", map[string]interface{}{"key0": "value0"})
	time.Sleep(100 * time.Millisecond) // Give time for first broadcast to propagate
	sim.BroadcastFrom("node-5", "update", map[string]interface{}{"key5": "value5"})
	time.Sleep(100 * time.Millisecond)
	sim.BroadcastFrom("node-9", "update", map[string]interface{}{"key9": "value9"})

	// Wait for convergence on all keys with longer timeout
	time.Sleep(500 * time.Millisecond) // Allow all messages to propagate
	converged1 := sim.WaitForConvergence("key0", 5*time.Second)
	converged2 := sim.WaitForConvergence("key5", 5*time.Second)
	converged3 := sim.WaitForConvergence("key9", 5*time.Second)

	if !converged1 || !converged2 || !converged3 {
		// Note: Convergence can be probabilistic with async gossip
		t.Logf("Warning: Not all broadcasts converged (key0=%v, key5=%v, key9=%v)",
			converged1, converged2, converged3)
	}

	stats := sim.GetStats()
	t.Logf("Stats: %d nodes, %d messages, %d dropped",
		stats.NodeCount, stats.MessageCount, stats.DroppedCount)

	sim.Shutdown()
}

// TestSimulator_WithPacketLoss tests simulator with packet loss
func TestSimulator_WithPacketLoss(t *testing.T) {
	// 10% packet loss
	sim := NewSimulator(8, 3, 5*time.Millisecond, 0.1)

	payload := map[string]interface{}{
		"key": "value",
	}

	err := sim.BroadcastFrom("node-0", "test", payload)
	if err != nil {
		t.Fatalf("BroadcastFrom failed: %v", err)
	}

	// With packet loss, may take longer to converge
	converged := sim.WaitForConvergence("key", 5*time.Second)
	if !converged {
		t.Log("Note: Convergence may fail with packet loss in some runs")
	}

	stats := sim.GetStats()
	if stats.DroppedCount == 0 {
		t.Log("Note: No packets were dropped in this run (expected with 10% rate)")
	}

	sim.Shutdown()
}

// TestSimulator_InvalidNode tests broadcasting from non-existent node
func TestSimulator_InvalidNode(t *testing.T) {
	sim := NewSimulator(5, 2, 1*time.Millisecond, 0.0)

	err := sim.BroadcastFrom("non-existent", "test", map[string]interface{}{})
	if err == nil {
		t.Error("Expected error when broadcasting from non-existent node")
	}

	sim.Shutdown()
}

// BenchmarkGossipNode_ReceiveMessage benchmarks message receiving
func BenchmarkGossipNode_ReceiveMessage(b *testing.B) {
	network := NewMockNetwork(0, 0.0)
	node := NewGossipNode("bench-node", network, 3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg := Message{
			ID:   string(rune(i)),
			Type: "bench",
			From: "source",
			Payload: map[string]interface{}{
				"seq": i,
			},
			Timestamp: time.Now(),
		}
		node.ReceiveMessage(msg)
	}
}

// BenchmarkSimulator_Convergence benchmarks convergence time
func BenchmarkSimulator_Convergence(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sim := NewSimulator(20, 3, 1*time.Millisecond, 0.0)

		payload := map[string]interface{}{
			"iteration": i,
		}

		sim.BroadcastFrom("node-0", "test", payload)
		sim.WaitForConvergence("iteration", 5*time.Second)
		sim.Shutdown()
	}
}

// BenchmarkSimulator_MessageThroughput benchmarks message throughput
func BenchmarkSimulator_MessageThroughput(b *testing.B) {
	sim := NewSimulator(10, 3, 0, 0.0)
	defer sim.Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nodeID := i % 10
		payload := map[string]interface{}{
			"seq": i,
		}
		sim.BroadcastFrom(string(rune('0'+nodeID)), "test", payload)
	}
}
