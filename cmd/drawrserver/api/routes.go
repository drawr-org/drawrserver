package api

import (
	"errors"
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/render"

	log "github.com/golang/glog"
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
	router.Use(originsMw)
	router.Use(loggerMw)

	router.Mount("/session", sessionRouter())
	router.Mount("/stats", statRouter())
	router.Get("/version", versionHandler)
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

// originsMw allows connections from everywhere
func originsMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func loggerMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if log.V(2) {
			log.Infof("[%v][%v] :: %v", r.Proto, r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}

func wrapJSON(w http.ResponseWriter, r *http.Request, fieldname string, v interface{}) {
	var wrapped = make(map[string]interface{})
	wrapped[fieldname] = v
	render.JSON(w, r, wrapped)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	wrapJSON(w, r, "version", version)
}
