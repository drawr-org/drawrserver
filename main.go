package main

import (
	"log"
	"net/http"
	"time"

	"github.com/drawr-team/core-server/api"
	"github.com/pressly/chi"
)

var s http.Server

func init() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("this is the backend of the drawr service"))
	})

	// api version 1
	api.Routing(r)

	s = http.Server{
		Addr:         "localhost:8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func main() {
	log.Println("Listening on...", s.Addr)
	panic(s.ListenAndServe())
}
