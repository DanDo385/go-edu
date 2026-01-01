//go:build !solution && !reference

package websocketchatroom

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"github.com/gorilla/websocket"
)

func NewHub() *Hub {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func (h *Hub) Run() {
	// TODO: Implement this function
	panic("not implemented")
}

func (h *Hub) GetOrCreateRoom(name string) *Room {
	// TODO: Implement this function
	panic("not implemented")
}

func (h *Hub) Shutdown() {
	// TODO: Implement this function
	panic("not implemented")
}

func (r *Room) Run() {
	// TODO: Implement this function
	panic("not implemented")
}

func (r *Room) broadcastToAll(message []byte) {
	// TODO: Implement this function
	panic("not implemented")
}

func NewClient(hub *Hub, conn *websocket.Conn, username, roomName string) *Client {
	// TODO: Implement this function
	panic("not implemented")
}

func (c *Client) ReadPump() {
	// TODO: Implement this function
	panic("not implemented")
}

func (c *Client) WritePump() {
	// TODO: Implement this function
	panic("not implemented")
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this function
	panic("not implemented")
}
