package main

import (
	"testing"
	"time"
)

type MockProvider struct {
	mChan chan []byte
	aChan chan []byte
	T     *testing.T
}

func (p MockProvider) Emit() []byte {
	select {
	case incoming := <-p.mChan:
		return incoming
	case <-time.After(1 * time.Second):
		return nil
	}
}

func (p MockProvider) Absorb(message []byte) {
	// p.T.Logf("got a new message:\n%v\n", string(message))
	p.aChan <- message
}

func TestMonitorForNewSession(t *testing.T) {
	// var p = MockProvider{
	// 	mChan: make(chan []byte, 1),
	// 	aChan: make(chan []byte, 1),
	// 	T:     t,
	// }

	// // create a message
	// mockNewSession, err := message.CreateMessage("new-session", message.NewSessionData{
	// 	Username:    "johndoe",
	// 	SessionName: "test session",
	// })
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// // marshal it
	// rawMessage, err := json.Marshal(mockNewSession)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// p.mChan <- rawMessage

	// if err := monitor(p); err != nil {
	// 	t.Fatal(err)
	// }

	// rawAckMsg := <-p.aChan
	// var ackMsg message.GenericMessage
	// if err := json.Unmarshal(rawAckMsg, &ackMsg); err != nil {
	// 	t.Fatal(err)
	// }
	// var data map[string]string
	// if err := json.Unmarshal(ackMsg.Data, &data); err != nil {
	// 	t.Fatal(err)
	// }

	// if status := data["Status"]; status != "new-session-success" {
	// 	t.Fatal("status not as expected:", status)
	// }
	// if msg := data["Message"]; msg != "Session successfully created." {
	// 	t.Fatal("message not as expected:", msg)
	// }
	// if id := data["SessionID"]; id == "" {
	// 	t.Fatal("no session id:", id)
	// }
}
