package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/drawr-team/core-server/message"
	"github.com/drawr-team/core-server/websock"
)

// Hub wraps a websock.Hub
type Hub struct {
	*websock.Hub
}

// Emit implements the MessageProvider interface
// returns a message out of the IncomingBus
func (h Hub) Emit() []byte {
	return <-h.IncomingBus
}

// Absorb implements the Absorber interface
// pushes a message into the BroadcastBus
func (h Hub) Absorb(message []byte) {
	h.BroadcastBus <- message
}

func monitor(provider message.Provider) error {
	for {
		msg := provider.Emit()
		if msg == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		var m message.GenericMessage
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Println(err)
			continue
		}

		switch m.Type {
		case "new-session":
			log.Println("[monitor] found a `new-session` message")
			if err := message.HandleNewSession(m, provider); err != nil {
				return err
			}
		case "join-session":
			log.Println("[monitor] found a `join-session` message")
			if err := message.HandleJoinSession(m, provider); err != nil {
				return err
			}
		default:
			log.Println("[monitor] unknown message type")
		}

	}
}
