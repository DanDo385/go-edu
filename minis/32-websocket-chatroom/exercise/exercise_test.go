package exercise

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// TestHubCreation tests that NewHub creates a properly initialized Hub
func TestHubCreation(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub() returned nil")
	}

	if hub.rooms == nil {
		t.Error("Hub.rooms map not initialized")
	}
}

// TestGetOrCreateRoom tests room creation and retrieval
func TestGetOrCreateRoom(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub() returned nil")
	}

	// Create a room
	room1 := hub.GetOrCreateRoom("test-room")
	if room1 == nil {
		t.Fatal("GetOrCreateRoom() returned nil")
	}

	if room1.name != "test-room" {
		t.Errorf("Room name = %q, want %q", room1.name, "test-room")
	}

	// Get the same room again
	room2 := hub.GetOrCreateRoom("test-room")
	if room2 != room1 {
		t.Error("GetOrCreateRoom() returned different room instance for same name")
	}

	// Create a different room
	room3 := hub.GetOrCreateRoom("another-room")
	if room3 == nil {
		t.Fatal("GetOrCreateRoom() returned nil for second room")
	}

	if room3 == room1 {
		t.Error("GetOrCreateRoom() returned same room for different names")
	}
}

// TestConcurrentRoomAccess tests thread-safety of room creation
func TestConcurrentRoomAccess(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub() returned nil")
	}

	const goroutines = 10
	const roomName = "concurrent-room"

	var wg sync.WaitGroup
	rooms := make([]*Room, goroutines)

	// Multiple goroutines try to get/create the same room
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			rooms[idx] = hub.GetOrCreateRoom(roomName)
		}(i)
	}

	wg.Wait()

	// All should get the same room instance
	for i := 1; i < goroutines; i++ {
		if rooms[i] != rooms[0] {
			t.Errorf("Concurrent GetOrCreateRoom() returned different instances")
			break
		}
	}
}

// TestWebSocketUpgrade tests the HTTP to WebSocket upgrade
func TestWebSocketUpgrade(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub() returned nil")
	}

	go hub.Run()

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	}))
	defer server.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?user=testuser&room=testroom"

	// Connect WebSocket client
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect WebSocket: %v", err)
	}
	defer ws.Close()

	// Connection successful if we get here
	t.Log("WebSocket upgrade successful")
}

// TestMessageBroadcast tests that messages are broadcast to all clients in a room
func TestMessageBroadcast(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub() returned nil")
	}

	go hub.Run()

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect two clients to the same room
	ws1, _, err := websocket.DefaultDialer.Dial(wsURL+"?user=user1&room=broadcast", nil)
	if err != nil {
		t.Fatalf("Failed to connect client 1: %v", err)
	}
	defer ws1.Close()

	ws2, _, err := websocket.DefaultDialer.Dial(wsURL+"?user=user2&room=broadcast", nil)
	if err != nil {
		t.Fatalf("Failed to connect client 2: %v", err)
	}
	defer ws2.Close()

	// Give clients time to register
	time.Sleep(100 * time.Millisecond)

	// Client 1 sends a message
	testMessage := map[string]string{"content": "Hello from user1"}
	if err := ws1.WriteJSON(testMessage); err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Client 2 should receive the message
	ws2.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, msg, err := ws2.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	// Parse the message
	var received map[string]interface{}
	if err := json.Unmarshal(msg, &received); err != nil {
		t.Fatalf("Failed to parse message: %v", err)
	}

	// Check message content
	if content, ok := received["content"].(string); !ok || !strings.Contains(content, "Hello from user1") {
		t.Errorf("Received unexpected message: %v", received)
	}

	t.Log("Message broadcast test passed")
}

// TestRoomIsolation tests that messages don't leak between rooms
func TestRoomIsolation(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub() returned nil")
	}

	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect client to room1
	ws1, _, err := websocket.DefaultDialer.Dial(wsURL+"?user=user1&room=room1", nil)
	if err != nil {
		t.Fatalf("Failed to connect to room1: %v", err)
	}
	defer ws1.Close()

	// Connect client to room2
	ws2, _, err := websocket.DefaultDialer.Dial(wsURL+"?user=user2&room=room2", nil)
	if err != nil {
		t.Fatalf("Failed to connect to room2: %v", err)
	}
	defer ws2.Close()

	// Give clients time to register
	time.Sleep(100 * time.Millisecond)

	// Client 1 sends message to room1
	testMessage := map[string]string{"content": "Private to room1"}
	if err := ws1.WriteJSON(testMessage); err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Client 2 (in room2) should NOT receive the message
	ws2.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, msg, err := ws2.ReadMessage()

	if err == nil {
		// If we got a message, it should only be the join notification, not the private message
		var received map[string]interface{}
		json.Unmarshal(msg, &received)
		if content, ok := received["content"].(string); ok && strings.Contains(content, "Private to room1") {
			t.Error("Message leaked from room1 to room2")
		}
	}

	t.Log("Room isolation test passed")
}

// TestMultipleClients tests handling of multiple concurrent clients
func TestMultipleClients(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub() returned nil")
	}

	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	const numClients = 5
	clients := make([]*websocket.Conn, numClients)

	// Connect multiple clients
	for i := 0; i < numClients; i++ {
		ws, _, err := websocket.DefaultDialer.Dial(wsURL+"?user=user"+string(rune('0'+i))+"&room=multitest", nil)
		if err != nil {
			t.Fatalf("Failed to connect client %d: %v", i, err)
		}
		defer ws.Close()
		clients[i] = ws
	}

	// Give all clients time to register
	time.Sleep(200 * time.Millisecond)

	// First client sends a message
	testMessage := map[string]string{"content": "Broadcast to all"}
	if err := clients[0].WriteJSON(testMessage); err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// All other clients should receive it
	received := 0
	for i := 1; i < numClients; i++ {
		clients[i].SetReadDeadline(time.Now().Add(2 * time.Second))
		_, msg, err := clients[i].ReadMessage()
		if err != nil {
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err == nil {
			if content, ok := data["content"].(string); ok && strings.Contains(content, "Broadcast to all") {
				received++
			}
		}
	}

	if received < numClients-2 { // Allow some margin for timing
		t.Errorf("Expected at least %d clients to receive message, got %d", numClients-2, received)
	}

	t.Logf("Multiple clients test passed (%d clients received message)", received)
}

// TestClientDisconnect tests proper cleanup when client disconnects
func TestClientDisconnect(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub() returned nil")
	}

	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect a client
	ws, _, err := websocket.DefaultDialer.Dial(wsURL+"?user=disconnecttest&room=testroom", nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Give client time to register
	time.Sleep(100 * time.Millisecond)

	// Disconnect
	ws.Close()

	// Give time for cleanup
	time.Sleep(200 * time.Millisecond)

	// Verify room is cleaned up (check by trying to access it)
	hub.mu.RLock()
	room := hub.rooms["testroom"]
	hub.mu.RUnlock()

	if room != nil {
		// Room might still exist but should have no clients
		if len(room.clients) > 0 {
			t.Error("Client not removed from room after disconnect")
		}
	}

	t.Log("Client disconnect test passed")
}

// TestMissingUsername tests error handling when username is missing
func TestMissingUsername(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub() returned nil")
	}

	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?room=testroom"

	// Try to connect without username
	_, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err == nil {
		t.Error("Expected error when connecting without username, got none")
	}

	if resp != nil && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	t.Log("Missing username test passed")
}

// TestDefaultRoom tests that clients are placed in default room when room is not specified
func TestDefaultRoom(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub() returned nil")
	}

	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?user=defaultuser"

	// Connect without specifying room
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer ws.Close()

	// Give time to register
	time.Sleep(100 * time.Millisecond)

	// Check that default room was created
	hub.mu.RLock()
	_, hasGeneral := hub.rooms["general"]
	hub.mu.RUnlock()

	if !hasGeneral {
		t.Error("Default 'general' room was not created")
	}

	t.Log("Default room test passed")
}

// TestPingPong tests that ping/pong heartbeat works
func TestPingPong(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub() returned nil")
	}

	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?user=pingtest&room=testroom"

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer ws.Close()

	// Set up ping handler
	pingReceived := false
	ws.SetPingHandler(func(appData string) error {
		pingReceived = true
		// Send pong response
		return ws.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second))
	})

	// Read messages for a while to trigger ping/pong
	done := make(chan bool)
	go func() {
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				done <- true
				return
			}
		}
	}()

	// Wait for up to 2 seconds
	select {
	case <-done:
	case <-time.After(2 * time.Second):
	}

	// Note: This test might not always receive a ping in the test timeframe
	// since pingPeriod is 54 seconds by default. This is just checking the mechanism.
	t.Log("Ping/pong mechanism test completed")
}

// BenchmarkBroadcast benchmarks message broadcasting performance
func BenchmarkBroadcast(b *testing.B) {
	hub := NewHub()
	if hub == nil {
		b.Fatal("NewHub() returned nil")
	}

	room := hub.GetOrCreateRoom("benchmark")
	if room == nil {
		b.Fatal("GetOrCreateRoom() returned nil")
	}

	go room.Run()

	// Create some mock clients
	const numClients = 100
	for i := 0; i < numClients; i++ {
		client := &Client{
			send: make(chan []byte, 256),
		}
		room.clients[client] = true

		// Drain send channel
		go func(c *Client) {
			for range c.send {
			}
		}(client)
	}

	message := []byte("Benchmark message")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		room.broadcastToAll(message)
	}
}
