package drawr

import (
	"net/http"
	"time"

	"./api"
)

// Server wraps http.Server
type Server struct {
	http.Server
}

// NewServer returns a drawr.Server that serves the API of the drawr backend
func NewServer() *Server {
	return &Server{
		http.Server{
			Addr:           ":8080",
			Handler:        api.New(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}
