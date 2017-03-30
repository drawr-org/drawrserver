package websock

import "fmt"

// ListConnections returns a list of connection descriptions
func (h *Hub) ListConnections() []string {
	var ls []string

	for id, c := range h.connections {
		ls = append(ls, fmt.Sprintf("ID: %v, %v", id, c.Addr))
	}

	return ls
}
