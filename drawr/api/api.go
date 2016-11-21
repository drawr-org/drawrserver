package api

import (
	"fmt"
	"net/http"
)

// API wraps http.ServeMux
type API struct {
	http.ServeMux
}

type apiv1Handler struct{}

// New returns the http.ServeMux for the API root
func New() http.Handler {
	api := &API{}
	api.init()

	return api
}

// Init sets up all server handlers
// It runs on api.New()
func (api *API) init() {
	mux := http.NewServeMux()
	mux.Handle("/api/v1/", apiv1Handler)
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	api = &API{*mux}
}
