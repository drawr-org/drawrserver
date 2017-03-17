package websock

import (
	"errors"
	"log"
	"os"
	"sync"
)

var ErrConnectionNotFound = errors.New("Connection not found")

// Hub keeps track of all connections
type Hub struct {
	Verbose bool

	// TODO maybe use map[ULID]Connection
	connections   map[string]Connection
	connectionsMx sync.RWMutex

	// TODO hub-wide message channel?
	// broadcast chan []byte

	log log.Logger
}

// NewHub creates a new hub
func NewHub() *Hub {
	h := &Hub{
		Verbose:       false,
		connections:   make(map[string]Connection),
		connectionsMx: sync.RWMutex{},
		// broadcast:     make(chan []byte, 2048),
		log: *log.New(os.Stdout, "[websock]", log.LstdFlags),
	}
	return h
}

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
		if err := c.Close(); err != nil {
			panic(err)
		}
		delete(h.connections, id)
	}
}

func (h *Hub) GetConnection(id string) (*Connection, error) {
	c, ok := h.connections[id]
	if !ok {
		return nil, ErrConnectionNotFound
	}
	return &c, nil
}

// Broadcast sends a message to all connections
func (h *Hub) Broadcast(m []byte) {
	for _, conn := range h.connections {
		conn.SendChan() <- m
	}
}
