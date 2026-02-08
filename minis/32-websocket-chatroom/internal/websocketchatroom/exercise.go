//go:build !solution && !reference

package websocketchatroom

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
	roomName string
}

type Room struct {
	name       string
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Hub struct {
	rooms      map[string]*Room
	mu         sync.RWMutex
	register   chan *Client
	unregister chan *Client
}

// NewHub - TODO: implement this function
func NewHub() *Hub {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// Run - TODO: implement this function
func (h *Hub) Run() {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return
}

// GetOrCreateRoom - TODO: implement this function
func (h *Hub) GetOrCreateRoom(name string) *Room {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// Shutdown - TODO: implement this function
func (h *Hub) Shutdown() {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return
}

// Run - TODO: implement this function
func (r *Room) Run() {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return
}

// broadcastToAll - TODO: implement this function
func (r *Room) broadcastToAll(message []byte) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return
}

// NewClient - TODO: implement this function
func NewClient(hub *Hub, conn *websocket.Conn, username, roomName string) *Client {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// ReadPump - TODO: implement this function
func (c *Client) ReadPump() {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return
}

// WritePump - TODO: implement this function
func (c *Client) WritePump() {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return
}

// ServeWS - TODO: implement this function
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
}

