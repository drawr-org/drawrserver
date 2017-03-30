package websock

import "github.com/gorilla/websocket"

// Connection wraps a messenger connection
type Connection struct {
	messenger Messenger

	quit chan chan struct{}
	// connection message channels
	send     chan []byte
	received chan []byte

	Addr string
}

// NewConnection constructs a new Connection from a websocket.Conn
func NewConnection(ws *websocket.Conn) *Connection {
	println("NewConnection:", ws.LocalAddr())
	c := &Connection{
		messenger: WebsocketMessenger{ws},
		quit:      make(chan chan struct{}),
		send:      make(chan []byte),
		received:  make(chan []byte),
	}
	if ws != nil {
		c.Addr = ws.RemoteAddr().String()
	} else {
		c.Addr = "none"
	}

	go c.reader()
	go c.writer()
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
	q := make(chan struct{})
	c.quit <- q
	<-q

	return c.messenger.Close()
}

// Reader reads a message from the websocket connection
func (c *Connection) reader() {
	for {
		msg, err := c.messenger.ReadMessage()
		if err != nil {
			// TODO handle ReadMessage errors
			panic(err)
		}

		select {
		case c.received <- msg:
			println("got msg:", msg)
		case q := <-c.quit:
			close(q)
			return
		}
	}
}

// writer writes a message to the websocket connection
func (c *Connection) writer() {
	for {
		select {
		case msg := <-c.send:
			println("write msg:", msg)
			err := c.messenger.WriteMessage(msg)
			if err != nil {
				// TODO handle WriteMessage errors
				panic(err)
			}
		case q := <-c.quit:
			close(q)
			return
		}
	}
}
