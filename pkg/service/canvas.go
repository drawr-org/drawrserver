package service

import (
	"net/http"

	"github.com/drawr-team/drawrserver/pkg/model"
	"github.com/drawr-team/drawrserver/pkg/ulidgen"
	"github.com/drawr-team/drawrserver/pkg/websock"

	log "github.com/golang/glog"
)

// WSConnect adds a new client connection to the session hub
func WSConnect(w http.ResponseWriter, r *http.Request, s model.Session) error {
	log.Info("upgrade websocket connection")
	c, err := websock.Upgrade(w, r, w.Header())
	if err != nil {
		return err
	}

	_, ok := svc.hubs[s.ID]
	if !ok {
		log.Infof("creating new hub")
		svc.hubs[s.ID] = websock.NewHub()
	}

	log.Infof("joining hub: %s", s.ID)
	svc.hubs[s.ID].AddConnection(ulidgen.Now().String(), *c)
	return nil
}
