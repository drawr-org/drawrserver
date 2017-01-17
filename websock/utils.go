package websock

import "fmt"

// ListConnections returns a list of connection descriptions
func (h *Hub) ListConnections() []string {
	var ls []string

	for id, p := range h.Pools {
		ls = append(ls, "ID: "+id)
		for c := range p.connections {
			ls = append(ls, fmt.Sprintf("local IP: %v, remote IP: %v, subproto: %v", c.SocketConn().LocalAddr(), c.SocketConn().RemoteAddr(), c.SocketConn().Subprotocol()))
		}
	}

	return ls
}
