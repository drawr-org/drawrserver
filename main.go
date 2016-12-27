package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/drawr-team/core-server/bolt"
	"github.com/drawr-team/core-server/websock"
)

var (
	dbClient bolt.Client
	s        http.Server
	port     string
)

func init() {
	flag.StringVar(&port, "p", "8080", "port to run the server on")
	flag.Parse()

	// initialize the database client and open the connection
	bolt.NewClient()

	// initialize a new communication hub
	wsHub := websock.NewHub()
	go wsHub.Run()
	go monitor(Hub{wsHub})

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
	mux.Handle("/ws", websock.Handler{Hub: wsHub})

	s = http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func main() {
	defer bolt.Close()

	log.Println("Listening on...", s.Addr)
	panic(s.ListenAndServe())
}
