package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/drawr-team/drawrserver/pkg/api"
	"github.com/drawr-team/drawrserver/pkg/bolt"
)

var (
	port         = flag.String("p", "3000", "port to run the server on")
	dbTimeout    = flag.Int64("timeout", 1, "how long until giving up on database transaction in Seconds")
	dbPath       = flag.String("database", "data.db", "location of the database file")
	verbose      = flag.Bool("verbose", false, "show log messages")
	debug        = flag.Bool("debug", false, "show log messages")
	printversion = flag.Bool("version", false, "print version number")
)

func init() {
	flag.Parse()
	if *printversion {
		fmt.Print(version)
		os.Exit(0)
	}
}

func main() {
	server := new(http.Server)

	// TODO make config loadable from JSON
	if err := api.Configure(server, &api.Options{
		Port:      *port,
		RWTimeout: int64(5 * time.Second),
		Database: &bolt.Options{
			Path:    *dbPath,
			Timeout: *dbTimeout * int64(time.Second),
			Verbose: *verbose,
		},
		Verbose: *verbose,
		Debug:   *debug,
	}); err != nil {
		log.Fatal("Unable to configure API")
	}

	go catchSignalAndCleanup(server) // handle ctrl-c, etc...

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ServerError:", err)
	}
}
