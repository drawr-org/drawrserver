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

	Pools map[string]Pool
	// protect parallel handling of connections
	connectionsMx sync.RWMutex
	poolsMx       sync.RWMutex
}

// NewHub creates a new hub
func NewHub() *Hub {
	return &Hub{
		BroadcastBus:  make(chan []byte),
		IncomingBus:   make(chan []byte),
		Timeout:       1,
		Pools:         make(map[string]Pool),
		connectionsMx: sync.RWMutex{},
		poolsMx:       sync.RWMutex{},
	}
}

// Run starts monitoring on the hub
// in a goroutine
func (h *Hub) Run() {
	for {
		for id, p := range h.Pools {
			broadcastMessage := <-p.BroadcastBus
			h.connectionsMx.RLock()

			for c := range p.connections {
				select {
				case c.SendChan() <- broadcastMessage:

				// close connection after no response for Timeout
				case <-time.After(time.Duration(h.Timeout) * time.Second):
					if h.Verbose {
						log.Println("[hub] closing connection:", c)
					}
					h.RemoveConnection(id, c)
				}
			}

			h.connectionsMx.RUnlock()
		}
	}
}

// AddConnection remembers a connection
func (h *Hub) AddConnection(id string, c Connection) {
	h.poolsMx.Lock()
	defer h.poolsMx.Unlock()
	// if we don't have a pool with id already
	// create one...
	if h.Pools[id].ID != id {
		h.AddPool(id)
	}

	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()

	if h.Verbose {
		log.Println("[hub] new connection:", c)
	}

	h.Pools[id].connections[c] = struct{}{}
}

// RemoveConnection forgets a connection
func (h *Hub) RemoveConnection(id string, c Connection) {
	h.poolsMx.Lock()
	defer h.poolsMx.Unlock()

	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()

	if h.Verbose {
		log.Println("[hub] remove connection:", c)
	}

	if _, ok := h.Pools[id].connections[c]; ok {
		delete(h.Pools[id].connections, c)
		close(c.SendChan())
	}

	if len(h.Pools[id].connections) < 1 {
		delete(h.Pools, id)
	}
}
