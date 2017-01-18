package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/drawr-team/core-server/bolt"
	"github.com/drawr-team/core-server/message"
	"github.com/drawr-team/core-server/websock"
)

// HubProvider wraps a websock.Hub
type HubProvider struct {
	hub *websock.Hub
}

// Emit implements the MessageProvider interface
// returns a message out of the IncomingBus
func (h HubProvider) Emit() []byte {
	return <-h.hub.IncomingBus
}

// Absorb implements the Absorber interface
// pushes a message into the BroadcastBus
func (h HubProvider) Absorb(message []byte) {
	h.hub.BroadcastBus <- message
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

		if verbose {
			log.Printf("[monitor] found a `%v` message", m.Type)
		}
		switch m.Type {
		// case message.NewSessionMessageType:
		// 	if err := message.HandleNewSession(m, provider, db); err != nil {
		// 		return err
		// 	}
		// case message.JoinSessionMessageType:
		// 	if err := message.HandleJoinSession(m, provider, db); err != nil {
		// 		return err
		// 	}
		case "leave-session":
			log.Println("not yet implemented -.-")
			// TODO
		case message.UpdateCanvasMessageType:
			if err := message.HandleUpdateCanvas(m, provider, db); err != nil {
				return err
			}
		default:
			log.Println("[monitor] unknown message type")
		}

	}
}
