package websock

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Connection wraps a messenger connection
type Connection struct {
	messenger Messenger
	wg        *sync.WaitGroup
	done      bool

	// connection message channels
	send     chan []byte
	received chan []byte

	Addr string
}

// NewConnection constructs a new Connection from a websocket.Conn
func NewConnection(ws *websocket.Conn) *Connection {
	c := &Connection{
		messenger: WebsocketMessenger{ws},
		wg:        new(sync.WaitGroup),
		done:      false,
		send:      make(chan []byte, 256),
		received:  make(chan []byte, 256),
	}
	if ws != nil {
		c.Addr = ws.RemoteAddr().String()
	} else {
		c.Addr = "none"
	}
	return c
}

// RunWorkers starts the reader and writer in seperate goroutines
// for the connection and returns a sync.WaitGroup
func (c *Connection) RunWorkers() {
	go c.Reader()
	go c.Writer()
}

// Wait blocks until the Read and Write workers finish
func (c *Connection) Wait() {
	c.wg.Wait()
}

// StopWorkers sends the done signal to the workers
func (c *Connection) Close() error {
	close(c.send)
	c.done = true
	c.wg.Wait() // wait for Reader to finish
	close(c.received)
	return c.messenger.Close()
}

// SendChan returns the send channel
func (c *Connection) SendChan() chan []byte {
	return c.send
}

// ReceiveChan returns the received channel
func (c *Connection) ReceiveChan() chan []byte {
	return c.received
}

// Reader reads a message from the websocket connection
func (c *Connection) Reader() {
	c.wg.Add(1)
	defer c.wg.Done()

	for !c.done {
		message, err := c.messenger.ReadMessage()
		if err != nil {
			// TODO handle ReadMessage errors
			panic(err)
		} else {
			c.received <- message
		}
	}
}

// Writer writes a message to the websocket connection
func (c *Connection) Writer() {
	c.wg.Add(1)
	defer c.wg.Done()

	for message := range c.send {
		err := c.messenger.WriteMessage(message)
		if err != nil {
			// TODO handle WriteMessage errors
			panic(err)
		}
	}
}
