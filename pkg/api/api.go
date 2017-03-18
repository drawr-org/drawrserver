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

var (
	apilog   = log.New(os.Stdout, "[api]\t", log.LstdFlags)
	dbClient *bolt.Client
)

// Options holds the basic configuration of the http server
// TODO implement reading options from a config file
type Options struct {
	Port      string        `json:"port"`
	RWTimeout int64         `json:"timeout"`
	Verbose   bool          `json:"verbose"`
	Debug     bool          `json:"debug"`
	Database  *bolt.Options `json:"database"`
}

// Configure takes a http.Server and configures it with the specified Options
func Configure(server *http.Server, opts *Options) error {
	// setup router
	apilog.Println("starting server on :" + opts.Port)
	server.Addr = ":" + opts.Port
	server.ReadTimeout = time.Duration(opts.RWTimeout)
	server.WriteTimeout = time.Duration(opts.RWTimeout)

	// TODO can we move this to pkg/session?
	dbClient = bolt.NewClient(opts.Database)
	if err := dbClient.Open(); err != nil {
		apilog.Println("Error opening database:", err)
		return err
	}
	canvas.Init()

	routes, err := setupRoutes(opts)
	if err != nil {
		apilog.Println("Error setting up router:", err)
		return err
	}
	server.Handler = routes

	return nil
}

// Cleanup closes the database client
func Cleanup() error {
	// TODO notify websocket clients about server shutdown
	// close db
	dbClient.Close()

	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			return fmt.Errorf("%v", r)
		}
		return err
	}
	return nil
}
