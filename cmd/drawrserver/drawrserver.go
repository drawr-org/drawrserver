package main

import (
	"log"
	"net/http"

	"github.com/drawr-team/core-server/bolt"
)

var (
	dbClient bolt.DBClient
	server   http.Server
)

func main() {
	defer dbClient.Close()
	go signalHandler()

	log.Println("Listening on...", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
