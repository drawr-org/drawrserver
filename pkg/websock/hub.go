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
	broadcast chan []byte

	quit chan chan struct{}
	log  log.Logger
}

// NewHub creates a new hub
func NewHub() *Hub {
	h := &Hub{
		Verbose:       false,
		connections:   make(map[string]Connection),
		connectionsMx: sync.RWMutex{},
		broadcast:     make(chan []byte),
		quit:          make(chan chan struct{}),
		log:           *log.New(os.Stdout, "[websock]", log.LstdFlags),
	}
	go h.broadcaster()
	return h
}

func (h *Hub) BroadcastChan() chan []byte {
	return h.broadcast
}

// Close sends the quit signal to the monitor worker
func (h *Hub) Close() {
	q := make(chan struct{})
	h.quit <- q
	<-q
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

// broadcast sends a message to all connections
func (h *Hub) sendBroadcastMessage(m []byte) {
	for cID, conn := range h.connections {
		println("notified:", cID)
		conn.SendChan() <- m
	}
}

func (h *Hub) broadcaster() {
	println("starting broadcaster worker...")
	for {
		select {
		case msg := <-h.broadcast:
			println("broadcasting", string(msg))
			h.sendBroadcastMessage(msg)
		case q := <-h.quit:
			close(q)
			return
		}
	}
}
