package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/oklog/ulid"
	"github.com/pressly/chi"

	"github.com/drawr-team/core-server/bolt"
	"github.com/drawr-team/core-server/ulidgen"
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
		r.Use(sessionCtx)
		r.Get("/", getSession)
		r.Get("/ws", websocketHandler)
		r.Get("/leave", leaveSession)
	})
	return r
}

func newSession(w http.ResponseWriter, r *http.Request) {
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

func getSession(w http.ResponseWriter, r *http.Request) {
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

func leaveSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, ok := ctx.Value("session").(Session)
	if !ok {
		log.Println("error getting session from ctx")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	return
}
