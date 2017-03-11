package main

import (
	"log"
	"net/http"

	"github.com/drawr-team/drawrserver/pkg/ulidgen"
	"github.com/drawr-team/drawrserver/pkg/websock"
)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, ok := ctx.Value("session").(Session)
	if !ok {
		log.Println("error getting session from ctx")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if hubs[session.ID] == nil {
		hubs[session.ID] = websock.NewHub()
	}

	c, err := websock.Upgrade(w, r, w.Header())
	if err != nil {
		log.Println("error upgrading connection")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	uid := ulidgen.Now()
	hubs[session.ID].AddConnection(uid.String(), c)
}
