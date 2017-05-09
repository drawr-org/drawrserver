package api

import (
	"net/http"

	"github.com/drawr-team/drawrserver/pkg/service"

	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
)

func statRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", statReport)
	r.Get("/db", dbReport)
	return r
}

func statReport(w http.ResponseWriter, r *http.Request) {
	out, err := service.StatReport()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	render.PlainText(w, r, out)
}

func dbReport(w http.ResponseWriter, r *http.Request) {
	out, err := service.DatabaseReport()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	render.PlainText(w, r, out)
}
