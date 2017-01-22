package main

import (
	"log"
	"net/http"

	"github.com/drawr-team/core-server/websock"
)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, ok := ctx.Value("session").(Session)
	if !ok {
		log.Println("error getting session from ctx")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if wsHubs[session.ID] == nil {
		// initialize a new communication hub
		wsHub := websock.NewHub()
		wsHub.Verbose = verbose
		provHub := HubProvider{hub: wsHub}
		wsHubs[session.ID] = &provHub
	}

	go wsHubs[session.ID].hub.Run()
	go monitor(wsHubs[session.ID], dbClient)
	wsHandler := websock.Handler{Hub: wsHubs[session.ID].hub}
	wsHandler.ServeHTTP(w, r)
}
