//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package websocketchatroom

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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
// TODO: implement NewHub.
func NewHub() *Hub { panic("TODO: implement") }
// TODO: implement Run.
func (h *Hub) Run() { panic("TODO: implement") }
// TODO: implement GetOrCreateRoom.
func (h *Hub) GetOrCreateRoom(name string) *Room { panic("TODO: implement") }
// TODO: implement Shutdown.
func (h *Hub) Shutdown() { panic("TODO: implement") }
// TODO: implement Run.
func (r *Room) Run() { panic("TODO: implement") }
// TODO: implement broadcastToAll.
func (r *Room) broadcastToAll(message []byte) { panic("TODO: implement") }
// TODO: implement NewClient.
func NewClient(hub *Hub, conn *websocket.Conn, username, roomName string) *Client {
	panic("TODO: implement")
}
// TODO: implement ReadPump.
func (c *Client) ReadPump() { panic("TODO: implement") }
// TODO: implement WritePump.
func (c *Client) WritePump() { panic("TODO: implement") }
// TODO: implement ServeWS.
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) { panic("TODO: implement") }
