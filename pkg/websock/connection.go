package websock

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Connection wraps a messenger connection
type Connection struct {
	ws *websocket.Conn

	wg *sync.WaitGroup
	// connection message channels
	send     chan []byte
	received chan []byte

	Addr string
}

// NewConnection constructs a new Connection from a websocket.Conn
func NewConnection(ws *websocket.Conn) *Connection {
	println("NewConnection:", ws.LocalAddr())
	c := &Connection{
		ws:       ws,
		wg:       new(sync.WaitGroup),
		send:     make(chan []byte),
		received: make(chan []byte),
	}
	if ws != nil {
		c.Addr = ws.RemoteAddr().String()
	} else {
		c.Addr = "none"
	}

	go c.reader()
	go c.writer()

	c.ws.SetCloseHandler(func(code int, text string) error {
		switch code {
		case websocket.CloseGoingAway:
			println("peer going away")
		case websocket.CloseNormalClosure:
			println("peer closing normally")
		}
		return nil
	})

	return c
}

// SendChan returns the send channel
func (c *Connection) SendChan() chan []byte {
	return c.send
}

// ReceiveChan returns the received channel
func (c *Connection) ReceiveChan() chan []byte {
	return c.received
}

// Close sends the done signal to the workers
func (c *Connection) Close() error {
	println("closing connection", c.Addr)

	payload := websocket.FormatCloseMessage(websocket.CloseGoingAway, "server shutting down")
	if err := c.ws.WriteMessage(websocket.CloseMessage, payload); err != nil {
		return err
	}

	println("waiting for workers to quit")
	close(c.send) // closing send channel causes writer to stop
	c.wg.Wait()

	println("closing connection was successfull")
	return c.ws.Close()
}

// Reader reads a message from the websocket connection
func (c *Connection) reader() {
	c.wg.Add(1)
	for {
		t, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				println("connection to peer closed as expected")
				break
			}
			// TODO handle ReadMessage errors
			panic("Failed to read message: " + err.Error())
		}

		switch t {
		case websocket.TextMessage:
			c.received <- msg
		case websocket.BinaryMessage:
			panic("cannot handle binary message")
		}
	}
	c.wg.Done()
	println("reader ended")
}

// writer writes a message to the websocket connection
func (c *Connection) writer() {
	c.wg.Add(1)
	for msg := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			// TODO handle WriteMessage errors
			panic("Failed to write message: " + err.Error())
		}
	}
	c.wg.Done()
	println("writer ended")
}
