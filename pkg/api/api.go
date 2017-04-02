// Package api implements the HTTP API for the drawrserver
// TODO move the database client interface from pkg/bolt to this package to make migration to another DB possible
package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/drawr-team/drawrserver/pkg/bolt"
	"github.com/drawr-team/drawrserver/pkg/canvas"
)

const version string = "0.2.0"

var (
	apilog   = log.New(os.Stdout, "[api]\t", log.LstdFlags)
	dbClient *bolt.Client
)

// Configure takes a http.Server and configures it with the specified Options
func Configure(server *http.Server, opts *Options) error {
	// TODO can we move this to pkg/session?
	dbClient = bolt.NewClient(opts.Database)
	if err := dbClient.Open(); err != nil {
		apilog.Println("Error opening database:", err)
		return err
	}

	// initialize canvas service
	canvas.Init(opts.Verbose)

	// setup router
	routes, err := setupRoutes(opts)
	if err != nil {
		apilog.Println("Error setting up router:", err)
		return err
	}

	server.Addr = ":" + opts.Port
	server.ReadTimeout = time.Duration(opts.RWTimeout)
	server.WriteTimeout = time.Duration(opts.RWTimeout)
	server.Handler = routes

	apilog.Println("starting server on :" + opts.Port)
	return nil
}

// Cleanup closes the database client
func Cleanup() error {
	dbClient.Close() // close db client

	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			return fmt.Errorf("%v", r)
		}
		return err
	}
	return nil
}
