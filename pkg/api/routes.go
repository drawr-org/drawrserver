package api

import (
	"errors"
	"net/http"

	"github.com/pressly/chi"
)

type ctxKey int

const (
	sessionCtxKey ctxKey = iota
	userCtxKey
)

var ErrSettingUpRouter = errors.New("Error setting up router")

// setupRoutes sets up the root router and calls the subroute methods
func setupRoutes() (chi.Router, error) {
	router := chi.NewRouter()
	router.Use(allowAllOriginsMiddleware)

	router.Mount("/session", sessionRouter())
	router.Mount("/stats", statRouter())
	router.Mount("/", uiRouter())

	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			return nil, ErrSettingUpRouter
		}
		return nil, err
	}
	return router, nil
}

// allowAllOriginsMiddleware allows connections from everywhere
func allowAllOriginsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
