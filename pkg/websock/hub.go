package websock

import (
	"errors"
	"sync"

	log "github.com/golang/glog"
	"github.com/gorilla/websocket"
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
}

// NewHub creates a new hub
func NewHub() *Hub {
	h := &Hub{
		connections:   make(map[string]Connection),
		connectionsMx: sync.RWMutex{},
		broadcast:     make(chan []byte),
		quit:          make(chan chan struct{}),
	}
	go h.broadcaster()
	return h
}

func (h *Hub) BroadcastChan() chan []byte {
	return h.broadcast
}

// Close sends the quit signal to the monitor worker
func (h *Hub) Close() {
	log.Info("closing hub")
	q := make(chan struct{})
	h.quit <- q

	for cID := range h.connections {
		h.RemoveConnection(cID)
	}
	<-q
}

// AddConnection remembers a connection
func (h *Hub) AddConnection(id string, c Connection) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()

	h.connections[id] = c

	// h.log.Printf("add connection: %v (%v)", id, c.Addr)

}

// RemoveConnection forgets a connection
func (h *Hub) RemoveConnection(id string) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()

	c, ok := h.connections[id]
	if ok {
		if err := c.Close(); err != nil {
			log.Exitf("Failed to remove connection %d from hub: %s", id, err)
		}
		delete(h.connections, id)
	}

	// h.log.Printf("remove connection: %v (%v)", id, c.Addr)

}

func (h *Hub) GetConnection(id string) (*Connection, error) {
	c, ok := h.connections[id]
	if !ok {
		return nil, ErrConnectionNotFound
	}
	return &c, nil
}

// broadcast sends a message to all connections
func (h *Hub) sendBroadcastMessage(m []byte) error {
	for cID, conn := range h.connections {
		log.Info("notified:", cID)
		pm, err := websocket.NewPreparedMessage(websocket.TextMessage, m)
		if err != nil {
			return err
		}
		if err := conn.ws.WritePreparedMessage(pm); err != nil {
			return err
		}
	}
	return nil
}

func (h *Hub) broadcaster() {
	log.Info("starting broadcaster worker...")
	for {
		select {
		case msg := <-h.broadcast:
			log.Info("broadcasting:", string(msg))
			h.sendBroadcastMessage(msg)
		case q := <-h.quit:
			close(q)
			return
		}
	}
}
