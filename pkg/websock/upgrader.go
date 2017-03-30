package websock

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// WSUpgrader is our default websocket upgrader
var WSUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(*http.Request) bool { return true },
}

// Upgrade establishes a websocket connection
func Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Connection, error) {
	wsConn, err := WSUpgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	return NewConnection(wsConn), nil
}
