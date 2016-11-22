package api

import "net"

// User is a drawr client
// It is used to inform the host of the `Session` about
// who to update the canvas for
type User struct {
	ID              string
	IPAddress       net.IP
	Name            string
	WritePermission bool
}
