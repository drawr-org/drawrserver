package websock

import (
	"log"
	"sync"
	"time"
)

// Hub keeps track of all connections
type Hub struct {
	// IncomingBus is for incoming messages
	IncomingBus chan []byte
	// BroadcastBus is for outgoing messages
	BroadcastBus chan []byte
	// Timeout in seconds
	Timeout int64
	Verbose bool

	// connections is a list of registered connections
	// we have to track and sync
	connections map[Connection]struct{}
	// protect parallel handling of connections
	connectionsMx sync.RWMutex
}

// NewHub creates a new hub
func NewHub() *Hub {
	return &Hub{
		BroadcastBus:  make(chan []byte),
		IncomingBus:   make(chan []byte),
		Timeout:       10,
		connections:   make(map[Connection]struct{}),
		connectionsMx: sync.RWMutex{},
	}
}

// Run starts monitoring on the hub
// in a goroutine
func (h *Hub) Run() {
	for {
		broadcastMessage := <-h.BroadcastBus
		h.connectionsMx.RLock()

		for c := range h.connections {
			select {
			case c.SendChan() <- broadcastMessage:

			// close connection after no response for Timeout
			case <-time.After(time.Duration(h.Timeout) * time.Second):
				if h.Verbose {
					log.Println("[hub] closing connection:", c)
				}
				h.RemoveConnection(c)
			}
		}

		h.connectionsMx.RUnlock()
	}
}

// AddConnection remembers a connection
func (h *Hub) AddConnection(c Connection) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()

	if h.Verbose {
		log.Println("[hub] new connection:", c)
	}

	h.connections[c] = struct{}{}
}

// RemoveConnection forgets a connection
func (h *Hub) RemoveConnection(c Connection) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()

	if h.Verbose {
		log.Println("[hub] remove connection:", c)
	}

	if _, ok := h.connections[c]; ok {
		delete(h.connections, c)
		close(c.SendChan())
	}
}
