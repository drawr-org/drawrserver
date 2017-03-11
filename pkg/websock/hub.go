package websock

import (
	"log"
	"os"
	"sync"
)

// Hub keeps track of all connections
type Hub struct {
	Verbose bool

	// TODO: maybe use map[ULID]Connection
	connections   map[string]Connection
	connectionsMx sync.RWMutex

	// TODO: hub-wide message channel?

	log log.Logger
}

// NewHub creates a new hub
func NewHub() *Hub {
	return &Hub{
		Verbose:       false,
		connections:   make(map[string]Connection),
		connectionsMx: sync.RWMutex{},
		log:           *log.New(os.Stdout, "[websock]", log.LstdFlags),
	}
}

// // Run starts monitoring on the hub
// func (h *Hub) Run() {
// 	for {
// 		broadcastMessage := <-h.BroadcastBus
// 		h.connectionsMx.RLock()

// 		for c := range h.connections {
// 			select {
// 			case c.SendChan() <- broadcastMessage:

// 			// close connection after no response for Timeout
// 			case <-time.After(time.Duration(h.Timeout) * time.Second):
// 				if h.Verbose {
// 					h.log.Println("closing connection:", c)
// 				}
// 				h.RemoveConnection(c)
// 			}
// 		}

// 		h.connectionsMx.RUnlock()
// 	}
// }

// AddConnection remembers a connection
func (h *Hub) AddConnection(id string, c Connection) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()

	if h.Verbose {
		h.log.Println("new connection:", id)
	}

	h.connections[id] = c
}

// RemoveConnection forgets a connection
func (h *Hub) RemoveConnection(id string) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()

	if h.Verbose {
		h.log.Println("remove connection:", id)
	}

	if c, ok := h.connections[id]; ok {
		if err := c.SocketConn().Close(); err != nil {
			panic(err)
		}
		delete(h.connections, id)
	}
}
