package websock

// Pool is a connection pool of a session
type Pool struct {
	ID string
	// BroadcastBus is for outgoing messages
	BroadcastBus chan []byte

	// connections is a list of registered connections
	// we have to track and sync
	connections map[Connection]struct{}
}

// AddPool creates a connection pool with id
func (h *Hub) AddPool(id string) {
	p := Pool{
		ID:           id,
		BroadcastBus: make(chan []byte),
		connections:  make(map[Connection]struct{}),
	}
	h.Pools[id] = p
}
