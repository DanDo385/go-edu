//go:build !solution && !reference

package websocketchatroom

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
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

// NewHub implements the exercise.
//
// TODO: Implement this function
func NewHub() *Hub {
	// TODO: Implement
	return nil
}

// Run implements the exercise.
//
// TODO: Implement this function
func (h *Hub) Run() {
	// TODO: Implement
}

// GetOrCreateRoom implements the exercise.
//
// TODO: Implement this function
func (h *Hub) GetOrCreateRoom(name string) *Room {
	// TODO: Implement
	return nil
}

// Shutdown implements the exercise.
//
// TODO: Implement this function
func (h *Hub) Shutdown() {
	// TODO: Implement
}

// Run implements the exercise.
//
// TODO: Implement this function
func (r *Room) Run() {
	// TODO: Implement
}

// broadcastToAll implements the exercise.
//
// TODO: Implement this function
func (r *Room) broadcastToAll(message []byte) {
	// TODO: Implement
}

// NewClient implements the exercise.
//
// TODO: Implement this function
func NewClient(hub *Hub, conn *websocket.Conn, username string, roomName string) *Client {
	// TODO: Implement
	return nil
}

// ReadPump implements the exercise.
//
// TODO: Implement this function
func (c *Client) ReadPump() {
	// TODO: Implement
}

// WritePump implements the exercise.
//
// TODO: Implement this function
func (c *Client) WritePump() {
	// TODO: Implement
}

// ServeWS implements the exercise.
//
// TODO: Implement this function
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
}
