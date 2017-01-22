package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/drawr-team/core-server/bolt"
	"github.com/pressly/chi"
)

var (
	dbClient  bolt.DBClient
	port      string
	dbPath    string
	dbTimeout int
	verbose   bool
	server    http.Server
	wsHubs    = make(map[string]*HubProvider)
)

func init() {
	flag.StringVar(&port, "p", "8080", "port to run the server on")
	flag.StringVar(&dbPath, "database", "data.db", "location of the database file")
	flag.IntVar(&dbTimeout, "timeout", 1, "how long until giving up on database transaction in Seconds")
	flag.BoolVar(&verbose, "verbose", false, "show log messages")

	printVersion := flag.Bool("version", false, "print version number")
	flag.Parse()

	if *printVersion {
		fmt.Print(version)
		os.Exit(0)
	}

	initDatabase()
	// initSocketHub()

	router := initRouter()

	server = http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func signalHandler() {
	// TODO: implement save shutdown
	// ELECTRON!
	// deal with external shutdown from electron
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)
	for sig := range sigChan {
		log.Println("received", sig, "...shutting down")
		dbClient.Close()
		// close listener here
		os.Exit(0)
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
		Verbose: verbose,
	}
	dbClient.Open()

}

func initRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(allowAllOrigins)
	// route: statistics about the websocket connections
	// mux.HandleFunc("/stats", func(w http.ResponseWriter, req *http.Request) {
	// 	w.Write([]byte("this is the backend of the drawr service\n"))
	// 	// list handler code here...
	// 	ls := wsHub.ListConnections()
	// 	for _, s := range ls {
	// 		w.Write([]byte(s + "\n"))
	// 	}

	// 	w.Write([]byte(fmt.Sprintf("\nfound %v connections\n", len(ls))))
	// })

	// route: validate a session id
	// r.Get("/validate", ValidateHandler)

	// route: easteregg
	r.Get("/teapot", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, http.StatusText(http.StatusTeapot), http.StatusTeapot)
		return
	})

	// route: sessions
	r.Mount("/session", sessionRouter())

	// route: web client
	r.Get("/", WebClientHandler)

	return r
}
