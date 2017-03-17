package canvas

import (
	"net/http"

	"github.com/drawr-team/drawrserver/pkg/session"
	"github.com/drawr-team/drawrserver/pkg/ulidgen"
	"github.com/drawr-team/drawrserver/pkg/websock"
)

var svc service

type service struct {
	hubs map[string]*websock.Hub
}

// Init initializes the message service
func Init() {
	svc.hubs = make(map[string]*websock.Hub)
}

// Connect adds a new client connection to the session hub
func Connect(w http.ResponseWriter, r *http.Request, s session.Session) error {
	_, ok := svc.hubs[s.ID]
	if !ok {
		svc.hubs[s.ID] = websock.NewHub()
	}

	c, err := websock.Upgrade(w, r, w.Header())
	if err != nil {
		return err
	}

	svc.hubs[s.ID].AddConnection(ulidgen.Now().String(), *c)
	return nil
}
