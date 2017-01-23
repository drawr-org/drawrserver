package websock

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Connection wraps the websocket connection
type Connection struct {
	// connection message channels
	send     chan []byte
	received chan []byte

	wsConn *websocket.Conn
}

// NewConnection constructs a new Connection
func NewConnection(ws *websocket.Conn) *Connection {
	return &Connection{
		wsConn:   ws,
		send:     make(chan []byte, 256),
		received: make(chan []byte, 256),
	}
}

// Worker starts the reader and writer
// for the connection and returns a sync.WaitGroup
func (c *Connection) Worker() {
	var wg sync.WaitGroup

	wg.Add(2)
	go c.Reader(&wg)
	go c.Writer(&wg)

	wg.Wait()
}

// SendChan returns the send channel
func (c *Connection) SendChan() chan []byte {
	return c.send
}

// ReceiveChan returns the received channel
func (c *Connection) ReceiveChan() chan []byte {
	return c.received
}

// SocketConn returns the actual websocket connection
func (c *Connection) SocketConn() *websocket.Conn {
	return c.wsConn
}

// Reader reads a message from the websocket connection
func (c *Connection) Reader(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		_, message, err := c.wsConn.ReadMessage()
		if err != nil {
			// TODO
			panic(err)
		} else {
			c.received <- message
		}
	}
}

// Writer writes a message to the websocket connection
func (c *Connection) Writer(wg *sync.WaitGroup) {
	defer wg.Done()

	for message := range c.send {
		err := c.wsConn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			// TODO
			panic(err)
		}
	}
}
