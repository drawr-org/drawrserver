package websock

import (
	"io"

	"github.com/gorilla/websocket"
)

// MessageReader can receive []byte messages
type MessageReader interface {
	ReadMessage() ([]byte, error)
}

// MessageWriter can send []byte messages
type MessageWriter interface {
	WriteMessage([]byte) error
}

// Messenger implements MessageReader and MessageWriter
type Messenger interface {
	MessageReader
	MessageWriter
	io.Closer
}

// WebsocketMessenger implements the the Messenger interface for websocket.Conn
type WebsocketMessenger struct{ *websocket.Conn }

// ReadMessage implements MessageReader for websocket.Conn
func (wsm WebsocketMessenger) ReadMessage() (b []byte, err error) {
	_, b, err = wsm.Conn.ReadMessage()
	return
}

// WriteMessage implements MessageWriter for websocket.Conn
func (wsm WebsocketMessenger) WriteMessage(b []byte) (err error) {
	err = wsm.Conn.WriteMessage(websocket.TextMessage, b)
	return
}

// Close implements Closer for websocket.Close
func (wsm WebsocketMessenger) Close() error {
	return wsm.Conn.Close()
}
