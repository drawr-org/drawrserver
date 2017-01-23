package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/drawr-team/core-server/bolt"
	"github.com/pressly/chi"
)

func allowAllOrigins(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

// SessionCtx retrieves the sessionID from the url and adds it to the context
func sessionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "sessionID")

		// __test__ is for testing purposes
		if id == "__test__" {
			http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
			return
		}

		b, err := dbClient.Get(bolt.SessionBucket, id)
		if err != nil {
			if err == bolt.ErrNotFound {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var session Session
		if err := json.Unmarshal(b, &session); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
