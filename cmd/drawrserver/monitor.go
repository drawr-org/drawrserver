package main

import (
	"encoding/json"
	"log"

	"github.com/drawr-team/core-server/bolt"
	"github.com/drawr-team/core-server/message"
	"github.com/drawr-team/core-server/websock"
)

// HubProvider wraps a websock.Hub
type HubProvider struct {
	hub     *websock.Hub
	verbose bool
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

// SetVerbose implements the Verbosity interface
func (h HubProvider) SetVerbose(verbose bool) {
	h.verbose = verbose
}

// GetVerbose implements the Verbosity interface
func (h HubProvider) GetVerbose() bool {
	return h.verbose
}

func (h HubProvider) shutdown(msg string) {
	m := message.ShutdownMessage(msg)

	h.hub.BroadcastBus <- m
}

func monitor(provider message.Provider, db bolt.DBClient) error {
	for {
		msg := provider.Emit()
		if msg == nil {
			// time.Sleep(1 * time.Second)
			continue
		}
		var m message.GenericMessage
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Println(err)
			continue
		}

		if provider.GetVerbose() {
			log.Printf("[monitor] found <%v> message\n", m.Type)
		}
		switch m.Type {
		case message.UpdateCanvasMessageType:
			if err := message.HandleUpdateCanvas(m, provider, db); err != nil {
				return err
			}
		default:
			log.Println("[monitor] unknown message type")
		}

	}
}
