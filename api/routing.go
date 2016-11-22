package api

import (
	"context"
	"net/http"

	"github.com/pressly/chi"
)

var (
	sessionService = initSessionService()
	userService    = initUserService()
)

// Routing builds the API structure
func Routing(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {

		// pass along the api_version
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				ctx := context.WithValue(req.Context(), "api_version", 1)
				next.ServeHTTP(w, req.WithContext(ctx))
			})
		})

		r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("this is api version 1"))
		}))

		// all sessions
		r.Route("/sessions", func(r chi.Router) {
			r.Get("/", sessionService.list)
			r.Post("/", sessionService.new)

			// single session
			r.Route("/:sessionID", func(r chi.Router) {
				r.Use(sessionService.context)
				r.Get("/", sessionService.get)
				r.Put("/", sessionService.update)
				r.Delete("/", sessionService.delete)

				// all users of a session
				r.Route("/users", func(r chi.Router) {
					r.Get("/", userService.list)
					r.Post("/", userService.new)

					// single user of a session
					r.Route("/:userID", func(r chi.Router) {
						r.Use(userService.context)
						r.Get("/", userService.get)
						r.Put("/", userService.update)
						r.Delete("/", userService.delete)
					})
				})
			})
		})

	})
}
