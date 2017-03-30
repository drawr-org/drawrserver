package websock

import (
	"sync"
	"testing"
)

type mockMessenger struct{ *testing.T }

const mockMessage = "hello world"

func (m mockMessenger) Close() error { return nil }
func (m mockMessenger) ReadMessage() ([]byte, error) {
	return []byte(mockMessage), nil
}
func (m mockMessenger) WriteMessage(b []byte) error {
	m.Logf("messenger: %s", string(b))
	return nil
}

func newMockConnection(t *testing.T) Connection {
	return Connection{
		messenger: mockMessenger{t},
		send:      make(chan []byte, 4096),
		received:  make(chan []byte, 4096),
		wg:        new(sync.WaitGroup),
	}
}

func TestReader(t *testing.T) {
	c := newMockConnection(t)
	go c.Reader()
	msg := string(<-c.ReceiveChan())
	if err := c.Close(); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if msg != mockMessage {
		t.Logf("msg=%v, want %v", msg, mockMessage)
		t.FailNow()
	}
}

func TestWriter(t *testing.T) {
	c := newMockConnection(t)
	go c.Writer()
	c.SendChan() <- []byte("hello world")
	if err := c.Close(); err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestMessenger(t *testing.T) {
	c := newMockConnection(t)
	c.RunWorkers()
	c.SendChan() <- []byte("test_0")
	c.SendChan() <- []byte("test_1")

	go func(t *testing.T, c *Connection) {
		c.Wait()
		t.Log("connection just closed")

	}(t, &c)

	if err := c.Close(); err != nil {
		t.Log(err)
		t.FailNow()
	}
}
