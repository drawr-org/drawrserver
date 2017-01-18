package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/oklog/ulid"
	"github.com/pressly/chi"

	"github.com/drawr-team/core-server/bolt"
	"github.com/drawr-team/core-server/ulidgen"
	"github.com/drawr-team/core-server/websock"
)

type User struct {
	Name string `json:"name"`
}

type Session struct {
	ID    string `json:"id"`
	Users []User `json:"users"`
}

func sessionRouter() http.Handler {
	r := chi.NewRouter()
	// r.Get("/", listSessions)
	r.Get("/new", newSession)
	r.Route("/:sessionID", func(r chi.Router) {
		r.Use(SessionCtx)
		r.Get("/", getSession)
		r.Get("/ws", websocketHandler)
	})
	return r
}

func newSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.URL.Query().Get("")
	uid := ulid.MustNew(ulidgen.GeneratorNow())

	session := Session{
		ID: uid.String(),
	}

	if err := dbClient.Put(bolt.SessionBucket, uid.String(), session); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(session)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Status", http.StatusText(http.StatusOK))
	w.Write(b)
}

func SessionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "sessionID")
		if b, err := dbClient.Get(bolt.SessionBucket, id); err != nil {
			if err == bolt.ErrNotFound {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else {
			var session Session
			if err := json.Unmarshal(b, &session); err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

	})
}

func getSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()
	session, ok := ctx.Value("session").(Session)
	if !ok {
		log.Println("error getting session from ctx")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(session)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Status", http.StatusText(http.StatusOK))
	w.Write(b)
}

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
