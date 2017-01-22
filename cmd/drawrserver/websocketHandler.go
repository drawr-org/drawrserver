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
		wsHub.Verbose = debug
		provHub := HubProvider{hub: wsHub}
		provHub.verbose = verbose
		wsHubs[session.ID] = &provHub

		if verbose {
			log.Println("[server] open new hub")
		}
	}

	go wsHubs[session.ID].hub.Run()
	go monitor(wsHubs[session.ID], dbClient)
	wsHandler := websock.Handler{Hub: wsHubs[session.ID].hub}

	if verbose {
		log.Println("[server] join hub")
	}
	wsHandler.ServeHTTP(w, r)
}
