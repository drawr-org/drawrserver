package api

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pressly/chi"

	"github.com/drawr-team/core-server/bolt"
)

// SessionService manages the session
type SessionService struct {
	dbClient *bolt.Client
}

// InitSessionService returns a new SessionService
func initSessionService() *SessionService {
	return &SessionService{
		dbClient: bolt.NewClient(),
	}
}

func (s *SessionService) list(w http.ResponseWriter, req *http.Request) {
	log.Println("session list")
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *SessionService) new(w http.ResponseWriter, req *http.Request) {
	log.Println("session new")

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	newSession, err := NewSession(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
	}

	if err := s.dbClient.Put(bolt.SessionBucket, newSession.ID, newSession); err != nil {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
	}

	http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
}

func (s *SessionService) context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// sessionID := chi.URLParam(req, "sessionID")
		// // get session from DB
		// session, err := s.dbClient.Get(bolt.SessionBucket, sessionID)
		// if err != nil {
		// 	http.Error(w, http.StatusText(404), 404)
		// 	return
		// }
		// // add session to context
		// ctx := context.WithValue(req.Context(), "session", session)
		// next.ServeHTTP(w, req.WithContext(ctx))
		next.ServeHTTP(w, req)
	})
}

func (s *SessionService) get(w http.ResponseWriter, req *http.Request) {
	log.Println("session get", chi.URLParam(req, "sessionID"))
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *SessionService) update(w http.ResponseWriter, req *http.Request) {
	log.Println("session update", chi.URLParam(req, "sessionID"))
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *SessionService) delete(w http.ResponseWriter, req *http.Request) {
	log.Println("session delete", chi.URLParam(req, "sessionID"))
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
