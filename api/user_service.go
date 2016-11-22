package api

import (
	"log"
	"net/http"

	"github.com/pressly/chi"

	"../bolt"
)

// UserService manages users in a session
type UserService struct {
	dbClient *bolt.Client
}

func initUserService() *UserService {
	return &UserService{}
}

func (s *UserService) list(w http.ResponseWriter, req *http.Request) {
	log.Println("user list")
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *UserService) new(w http.ResponseWriter, req *http.Request) {
	log.Println("user new")
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *UserService) context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// userID := chi.URLParam(req, "userID")
		// // get session from DB
		// user, err := s.dbClient.Get(bolt.UserBucket, userID)
		// if err != nil {
		// 	http.Error(w, http.StatusText(404), 404)
		// 	return
		// }
		// // add session to context
		// ctx := context.WithValue(req.Context(), "user", user)
		// next.ServeHTTP(w, req.WithContext(ctx))
		next.ServeHTTP(w, req)
	})
}

func (s *UserService) get(w http.ResponseWriter, req *http.Request) {
	log.Println("user get", chi.URLParam(req, "userID"))
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *UserService) update(w http.ResponseWriter, req *http.Request) {
	log.Println("user update", chi.URLParam(req, "userID"))
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *UserService) delete(w http.ResponseWriter, req *http.Request) {
	log.Println("user delete", chi.URLParam(req, "userID"))
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
