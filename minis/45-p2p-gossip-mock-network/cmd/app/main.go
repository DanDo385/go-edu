package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/memberlist"
)

// EventDelegate implements memberlist.EventDelegate to handle cluster events
type EventDelegate struct {
	name string
}

func (ed *EventDelegate) NotifyJoin(node *memberlist.Node) {
	fmt.Printf("[%s] Node joined: %s (%s)\n", ed.name, node.Name, node.Addr)
}

func (ed *EventDelegate) NotifyLeave(node *memberlist.Node) {
	fmt.Printf("[%s] Node left: %s (%s)\n", ed.name, node.Name, node.Addr)
}

func (ed *EventDelegate) NotifyUpdate(node *memberlist.Node) {
	fmt.Printf("[%s] Node updated: %s (%s)\n", ed.name, node.Name, node.Addr)
}

// BroadcastDelegate implements memberlist.Delegate for custom broadcasts
type BroadcastDelegate struct {
	name     string
	messages chan []byte
}

func NewBroadcastDelegate(name string) *BroadcastDelegate {
	return &BroadcastDelegate{
		name:     name,
		messages: make(chan []byte, 100),
	}
}

// NodeMeta is used to retrieve meta-data about the current node
// when broadcasting an alive message. It's length is limited to
// the given byte size. This metadata is available in the Node structure.
func (bd *BroadcastDelegate) NodeMeta(limit int) []byte {
	meta := fmt.Sprintf("role=worker,started=%d", time.Now().Unix())
	if len(meta) > limit {
		meta = meta[:limit]
	}
	return []byte(meta)
}

// NotifyMsg is called when a user-data message is received.
func (bd *BroadcastDelegate) NotifyMsg(msg []byte) {
	select {
	case bd.messages <- msg:
		fmt.Printf("[%s] Received message: %s\n", bd.name, string(msg))
	default:
		log.Printf("[%s] Message channel full, dropping message\n", bd.name)
	}
}

// GetBroadcasts is called when user data messages can be broadcast.
func (bd *BroadcastDelegate) GetBroadcasts(overhead, limit int) [][]byte {
	// For this demo, we don't queue broadcasts
	return nil
}

// LocalState is used for a TCP Push/Pull to send the local state.
func (bd *BroadcastDelegate) LocalState(join bool) []byte {
	// Return empty state for this demo
	return []byte{}
}

// MergeRemoteState is invoked after a TCP Push/Pull to merge remote state.
func (bd *BroadcastDelegate) MergeRemoteState(buf []byte, join bool) {
	// No-op for this demo
}

// CustomBroadcast implements memberlist.Broadcast
type CustomBroadcast struct {
	msg    []byte
	notify chan struct{}
}

func (b *CustomBroadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

func (b *CustomBroadcast) Message() []byte {
	return b.msg
}

func (b *CustomBroadcast) Finished() {
	if b.notify != nil {
		close(b.notify)
	}
}

func runNode(name string, bindPort int, joinAddr string) error {
	// Create memberlist configuration
	config := memberlist.DefaultLocalConfig()
	config.Name = name
	config.BindPort = bindPort
	config.AdvertisePort = bindPort

	// Set up event delegate
	eventDelegate := &EventDelegate{name: name}
	config.Events = eventDelegate

	// Set up broadcast delegate
	broadcastDelegate := NewBroadcastDelegate(name)
	config.Delegate = broadcastDelegate

	// Create memberlist
	list, err := memberlist.Create(config)
	if err != nil {
		return fmt.Errorf("failed to create memberlist: %w", err)
	}

	fmt.Printf("[%s] Memberlist created on port %d\n", name, bindPort)
	fmt.Printf("[%s] Node address: %s\n", name, list.LocalNode().Addr)

	// Join existing cluster if join address provided
	if joinAddr != "" {
		fmt.Printf("[%s] Attempting to join cluster at %s...\n", name, joinAddr)
		_, err = list.Join([]string{joinAddr})
		if err != nil {
			log.Printf("[%s] Failed to join cluster: %v\n", name, err)
			log.Printf("[%s] Continuing as standalone node\n", name)
		} else {
			fmt.Printf("[%s] Successfully joined cluster\n", name)
		}
	}

	// Print initial members
	fmt.Printf("[%s] Initial cluster members:\n", name)
	for _, member := range list.Members() {
		fmt.Printf("  - %s (%s:%d) Meta: %s\n",
			member.Name, member.Addr, member.Port, string(member.Meta))
	}

	// Periodically print member list
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Broadcast a message every 5 seconds
	broadcastTicker := time.NewTicker(5 * time.Second)
	defer broadcastTicker.Stop()

	fmt.Printf("[%s] Node is running. Press Ctrl+C to exit.\n", name)

	for {
		select {
		case <-ticker.C:
			members := list.Members()
			fmt.Printf("\n[%s] Current cluster members (%d):\n", name, len(members))
			for _, member := range members {
				state := "alive"
				if member.Name == name {
					state = "self"
				}
				fmt.Printf("  - %s (%s:%d) [%s]\n",
					member.Name, member.Addr, member.Port, state)
			}

		case <-broadcastTicker.C:
			// Broadcast a message to the cluster
			message := fmt.Sprintf("Hello from %s at %s",
				name, time.Now().Format(time.RFC3339))

			broadcast := &CustomBroadcast{
				msg:    []byte(message),
				notify: make(chan struct{}),
			}

			list.QueueBroadcast(broadcast)
			fmt.Printf("[%s] Broadcasted message: %s\n", name, message)

			// Wait for broadcast to complete
			go func() {
				<-broadcast.notify
				fmt.Printf("[%s] Broadcast completed\n", name)
			}()

		case msg := <-broadcastDelegate.messages:
			fmt.Printf("[%s] Processed received message: %s\n", name, string(msg))

		case sig := <-sigChan:
			fmt.Printf("\n[%s] Received signal %v, shutting down...\n", name, sig)

			// Leave the cluster gracefully
			err := list.Leave(5 * time.Second)
			if err != nil {
				log.Printf("[%s] Error leaving cluster: %v\n", name, err)
			}

			// Shutdown memberlist
			err = list.Shutdown()
			if err != nil {
				log.Printf("[%s] Error shutting down: %v\n", name, err)
			}

			fmt.Printf("[%s] Goodbye!\n", name)
			return nil
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Memberlist Gossip Protocol Demo")
		fmt.Println("\nUsage:")
		fmt.Println("  Start first node:  go run main.go <node-name> <port>")
		fmt.Println("  Join cluster:      go run main.go <node-name> <port> <join-address>")
		fmt.Println("\nExample:")
		fmt.Println("  Terminal 1: go run main.go node1 7946")
		fmt.Println("  Terminal 2: go run main.go node2 7947 127.0.0.1:7946")
		fmt.Println("  Terminal 3: go run main.go node3 7948 127.0.0.1:7946")
		os.Exit(1)
	}

	nodeName := os.Args[1]
	bindPort := 7946 // Default port

	if len(os.Args) >= 3 {
		fmt.Sscanf(os.Args[2], "%d", &bindPort)
	}

	joinAddr := ""
	if len(os.Args) >= 4 {
		joinAddr = os.Args[3]
	}

	fmt.Println("=== Memberlist Gossip Protocol Demo ===\n")
	fmt.Printf("Node Name: %s\n", nodeName)
	fmt.Printf("Bind Port: %d\n", bindPort)
	if joinAddr != "" {
		fmt.Printf("Join Address: %s\n", joinAddr)
	}
	fmt.Println()

	if err := runNode(nodeName, bindPort, joinAddr); err != nil {
		log.Fatalf("Error running node: %v", err)
	}
}
