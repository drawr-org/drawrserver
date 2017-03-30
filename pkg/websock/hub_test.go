package websock

import (
	"reflect"
	"testing"
)

func TestGetConnection(t *testing.T) {
	h := NewHub()

	cIn := newMockConnection(t)
	h.AddConnection("test", cIn)

	cOut, err := h.GetConnection("test")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(cIn, *cOut) {
		t.Log("connection is not what I put in there")
		t.FailNow()
	}
}

func TestBroadcastMessage(t *testing.T) {
	const bMessage = "broadcast this"
	h := NewHub()

	cA := newMockConnection(t)
	cB := newMockConnection(t)

	h.AddConnection("A", cA)
	h.AddConnection("B", cB)

	go cA.Writer()
	go cB.Writer()

	h.Broadcast([]byte(bMessage))

	cA.Wait()
	cB.Wait()
	h.RemoveConnection("A")
	h.RemoveConnection("B")
}
