package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/drawr-team/drawrserver/pkg/bolt"
	"github.com/pressly/chi"
)

var (
	port      string
	dbPath    string
	dbTimeout int
	verbose   bool
	debug     bool
	wsHubs    = make(map[string]*HubProvider)
)

func init() {
	flag.StringVar(&port, "p", "3000", "port to run the server on")
	flag.StringVar(&dbPath, "database", "data.db", "location of the database file")
	flag.IntVar(&dbTimeout, "timeout", 1, "how long until giving up on database transaction in Seconds")
	flag.BoolVar(&verbose, "verbose", false, "show log messages")
	flag.BoolVar(&debug, "debug", false, "show log messages")

	printVersion := flag.Bool("version", false, "print version number")
	flag.Parse()

	if *printVersion {
		fmt.Print(version)
		os.Exit(0)
	}

	initDatabase()

	router := initRouter()

	server = http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func initDatabase() {
	// initialize the database client and open the connection
	// TODO:
	// ELECTRON!
	// database path can't be next to binary
	dbClient = &bolt.Client{
		Path:    dbPath,
		Timeout: time.Duration(dbTimeout) * time.Second,
		Verbose: debug,
	}
	dbClient.Open()

}

func initRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(allowAllOrigins)

	// route: server statictics
	r.Get("/stats", func(w http.ResponseWriter, r *http.Request) {
		var out string

		out = out + "drawr backend:\n"

		out = out + "connected clients:\n"
		for id, hub := range wsHubs {
			out = out + "> " + id + ":"
			ls := hub.hub.ListConnections()
			for _, s := range ls {
				out = out + s + "\n"
			}
			out = out + "\n"
		}

		w.Write([]byte(out))
	})

	// route: easteregg
	r.Get("/teapot", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusTeapot), http.StatusTeapot)
		return
	})

	// route: sessions
	r.Mount("/session", sessionRouter())

	// route: web client
	r.Get("/", WebClientHandler)

	return r
}
