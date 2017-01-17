package websock

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Connection wraps the websocket connection
type Connection struct {
	// Hub is the hub this connection belongs to
	Hub *Hub

	wsConn *websocket.Conn
	// send is an outbound message channel
	send chan []byte
}

// NewConnection constructs a new Connection
func NewConnection(h *Hub, ws *websocket.Conn) Connection {
	return Connection{
		Hub:    h,
		wsConn: ws,
		send:   make(chan []byte, 256),
	}
}

// SendChan implements the Conn interface for websock.Hub
// Returns the send channel
func (c *Connection) SendChan() chan []byte {
	return c.send
}

// SocketConn implements the Conn interface for websock.Hub
// Returns the actual websocket connection
func (c *Connection) SocketConn() *websocket.Conn {
	return c.wsConn
}

// Reader reads a message from the websocket connection
func (c *Connection) Reader(wg *sync.WaitGroup, ws *websocket.Conn) {
	defer wg.Done()

	for {
		_, message, err := ws.ReadMessage()
		log.Println("[hub] received:", string(message))
		if err != nil {
			// TODO
			log.Println(err)
			break
		} else {
			c.Hub.IncomingBus <- message
		}
	}
}

// Writer writes a message to the websocket connection
func (c *Connection) Writer(wg *sync.WaitGroup, ws *websocket.Conn) {
	defer wg.Done()

	for message := range c.send {
		err := ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			// TODO
			log.Println(err)
			break
		}
		log.Println("[hub] sent:", string(message))
	}
}
