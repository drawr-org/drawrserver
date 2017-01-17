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

const version = "0.1.0"

var (
	dbClient  bolt.DBClient
	port      string
	dbPath    string
	dbTimeout int
	verbose   bool
	server    http.Server
)

func init() {
	flag.StringVar(&port, "p", "8080", "port to run the server on")
	flag.StringVar(&dbPath, "database", "data.db", "location of the database file")
	flag.IntVar(&dbTimeout, "timeout", 1, "how long until giving up on database transaction in Seconds")
	flag.BoolVar(&verbose, "verbose", false, "show log messages")

	printVersion := flag.Bool("version", false, "print version number")
	flag.Parse()

	if *printVersion {
		fmt.Printf("drawr server v%v\nfrom github.com/drawr-team/core-server\ncompiled at <%s>\n", version, time.Now().Format(time.ANSIC))
		os.Exit(0)
	}

	initDatabase()
	// initSocketHub()
	handle := initHandlers()

	server = http.Server{
		Addr:         ":" + port,
		Handler:      handle,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func main() {
	defer dbClient.Close()
	go signalHandler()

	log.Println("Listening on...", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
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

func initSocketHub() {
	// initialize a new communication hub
	wsHub := websock.NewHub()
	wsHub.Verbose = verbose

	go wsHub.Run()
	go monitor(HubProvider{wsHub}, dbClient)

}

func initHandlers() *http.ServeMux {
	mux := http.NewServeMux()
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
	mux.HandleFunc("/validate", ValidateHandler)
	// route: easteregg
	mux.HandleFunc("/teapot", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, http.StatusText(http.StatusTeapot), http.StatusTeapot)
		return
	})

	// route: sessions
	sessionHandler := SessionHandler{}
	mux.Handle("/session/", sessionHandler.MiddlewareHandler(sessionHandler))

	// route: web client
	mux.HandleFunc("/", WebClientHandler)

	return mux
}
