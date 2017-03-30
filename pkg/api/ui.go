package api

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
)

const assetDir string = "./assets"

func uiRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", render.NoContent) // dummy

	return r
}
