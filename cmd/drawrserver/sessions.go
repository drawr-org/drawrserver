package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/drawr-team/core-server/bolt"
	"github.com/oklog/ulid"
)

// SessionHandler spawns websockets on /:sessionid/ws suburls
type SessionHandler struct {
	db *bolt.Client
}

// MiddlewareHandler splits the url
func (s SessionHandler) MiddlewareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context

		p := strings.Split(r.URL.Path, "/")

		// no session id in url
		if len(p) <= 1 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if p[2] == "new" {
			log.Println("[handler] got new session")
			createSessionAndHub(s.db)
			http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
			return
		}

		ctx = context.WithValue(ctx, "session", p[2])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		status = http.StatusNotImplemented
	)
	defer http.Error(w, http.StatusText(status), status)

	uid, err := ulid.Parse(r.Context().Value("session").(string))
	if err != nil {
		log.Println("something about the ULID is wrong:", uid.String(), ", error was:", err)
		status = http.StatusBadRequest
		return
	}
	log.Println("[handler] requested session id:", uid.String())

	// if err := s.db.Put(bolt.SessionBucket, id, data); err != nil {
	// 	log.Println("[handler]", err)
	// 	status = http.StatusInternalServerError
	// 	return
	// }

	return
}
