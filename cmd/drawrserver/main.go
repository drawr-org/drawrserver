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
	"github.com/drawr-team/core-server/websock"
)

var (
	dbClient  bolt.DBClient
	s         http.Server
	port      string
	dbPath    string
	dbTimeout int
)

func init() {
	flag.StringVar(&port, "p", "8080", "port to run the server on")
	flag.StringVar(&dbPath, "db", "data.db", "location of the database file")
	flag.IntVar(&dbTimeout, "t", 1, "how long until giving up on database transaction in Seconds")
	flag.Parse()

	// initialize the database client and open the connection
	// TODO:
	// ELECTRON!
	// database path can't be next to binary
	dbClient = bolt.Client{
		Path:    dbPath,
		Timeout: time.Duration(dbTimeout) * time.Second,
	}
	dbClient.Open()

	// initialize a new communication hub
	wsHub := websock.NewHub()
	go wsHub.Run()
	go monitor(Hub{wsHub}, dbClient)

	mux := http.NewServeMux()
	// root route
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("this is the backend of the drawr service\n"))
		// list handler code here...
		var ls = wsHub.ListConnections()
		for _, s := range ls {
			w.Write([]byte(s + "\n"))
		}
		w.Write([]byte(fmt.Sprintf("\nfound %v connections\n", len(ls))))
	})
	// websocket route
	// TODO: do we want a websocket for each session?
	// like: /:session_id/ws
	mux.Handle("/ws", websock.Handler{Hub: wsHub})

	s = http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func main() {
	defer dbClient.Close()

	// TODO: implement save shutdown
	// ELECTRON!
	// code to deal with external shutdown from electron
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)
	go func() {
		for sig := range sigChan {
			log.Println("received", sig, "...shutting down")
			dbClient.Close()
			// close listener here
			os.Exit(0)
		}
	}()

	log.Println("Listening on...", s.Addr)
	panic(s.ListenAndServe())
}
