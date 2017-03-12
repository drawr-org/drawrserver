package api

import (
	"net/http"

	"github.com/drawr-team/drawrserver/pkg/session"
	"github.com/drawr-team/drawrserver/pkg/ulidgen"
	"github.com/drawr-team/drawrserver/pkg/websock"
)

func websocketConnect(w http.ResponseWriter, r *http.Request, s session.Session) error {
	hub, ok := hubs[s.ID]
	if !ok {
		hubs[s.ID] = *websock.NewHub()
		hub = hubs[s.ID]
	}

	c, err := websock.Upgrade(w, r, w.Header())
	if err != nil {
		return err
	}

	hub.AddConnection(ulidgen.Now().String(), *c)
	return nil
}
