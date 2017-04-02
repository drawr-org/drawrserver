package canvas

import (
	"net/http"

	"github.com/drawr-team/drawrserver/pkg/service"
	"github.com/drawr-team/drawrserver/pkg/ulidgen"
	"github.com/drawr-team/drawrserver/pkg/websock"
)

var svc canvasService

type canvasService struct {
	verbose bool

	hubs map[string]*websock.Hub
}

// Init initializes the message service
func Init(verbose bool) {
	svc.verbose = verbose
	svc.hubs = make(map[string]*websock.Hub)
}

// Close notifies all hubs about the server going offline
func Close() {
	for _, h := range svc.hubs {
		h.Close()
	}
}

// Connect adds a new client connection to the session hub
func Connect(w http.ResponseWriter, r *http.Request, s service.Session) error {
	_, ok := svc.hubs[s.ID]
	if !ok {
		svc.hubs[s.ID] = websock.NewHub()
		svc.hubs[s.ID].Verbose = svc.verbose
	}

	c, err := websock.Upgrade(w, r, w.Header())
	if err != nil {
		return err
	}

	svc.hubs[s.ID].AddConnection(ulidgen.Now().String(), *c)
	return nil
}
