package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/drawr-team/core-server/bolt"
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

func monitor(provider message.Provider, db bolt.DBClient) error {
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

		log.Printf("[monitor] found a `%v` message", m.Type)
		switch m.Type {
		case message.NewSessionMessageType:
			if err := message.HandleNewSession(m, provider, db); err != nil {
				return err
			}
		case message.JoinSessionMessageType:
			if err := message.HandleJoinSession(m, provider, db); err != nil {
				return err
			}
		case "leave-session":
			log.Println("not yet implemented -.-")
			// TODO
		case "update-canvas":
			if err := message.HandleUpdateCanvas(m, provider, db); err != nil {
				return err
			}
		default:
			log.Println("[monitor] unknown message type")
		}

	}
}
